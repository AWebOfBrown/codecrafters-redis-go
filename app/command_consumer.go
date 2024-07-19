package main

import (
	"fmt"
	"net"
	"strings"
)

func commandConsumer(command []*RESPToken, parser *RESPParser) ([]*RESPToken, error) {
	response, error := parser.Parse(command)
	return response, error
}

func handleTransaction(command []*RESPToken, conn *net.Conn, mc *MultiContext, parser *RESPParser) (interface{}, error) {
	cmdToken := command[1].Value.(string)

	var response interface{}

	switch strings.ToLower(cmdToken) {
	case "multi":
		token, _ := NewRESPToken(String, "OK")
		response = []*RESPToken{token}
	case "exec":
		execResponses := make([][]*RESPToken, 0)
		queuedCommands := mc.GetQueuedCommands(conn)
		qtyQueuedCommands := len(queuedCommands)
		if qtyQueuedCommands == 0 {
			emptyArray, _ := NewRESPToken(Array, "0")
			execResponses = append(execResponses, []*RESPToken{emptyArray})
			response = execResponses
		} else {
			for _, queuedCommand := range queuedCommands {
				if _, ok := response.([][]*RESPToken); ok {
					result, err := parser.Parse(queuedCommand)
					if err != nil {
						//todo: Handle better.
						errToken, _ := NewRESPToken(Error, "oops")
						execResponses = append(execResponses, []*RESPToken{errToken})
					}
					execResponses = append(execResponses, result)
				}
				response = execResponses
			}
		}
		mc.RemoveTxConnection(conn)
	default:
		//todo: Validate command before enqueue
		mc.EnqueueCommand(conn, command)
		token, _ := NewRESPToken(String, "QUEUED")
		response = []*RESPToken{token}
	}

	return response, nil
}

func CommandConsumerController(queue <-chan RedisCommandQueueMessage, dict map[string]string, multiContext *MultiContext) {
	parser := NewRESPParser(dict, multiContext)
	for {
		redisCommand := <-queue
		conn := *redisCommand.connection
		parser.SetClientConnection(redisCommand.connection)

		// Could be singular, or multi-responses in case of a TX
		var response any

		isActiveTx := multiContext.CheckActiveTX(&conn)

		if isActiveTx {
			result, err := handleTransaction(redisCommand.command, redisCommand.connection, multiContext, parser)
			if err != nil {
				conn.Write([]byte("+ERR\r\n"))
				conn.Close()
				return
			}
			response = result
		} else {
			result, err := commandConsumer(redisCommand.command, parser)
			if err != nil {
				responseToken, _ := NewRESPToken(Error, err.Error())
				conn.Write([]byte(responseToken.Value.([]byte)))
				conn.Close()
				return
			}
			response = result
		}

		responseBytes := serialiseTokens(response)
		conn.Write(responseBytes)
	}
}

func serialiseTokens(tokens any) []byte {
	responseData := make([]byte, 0)

	switch v := (tokens).(type) {
	case []*RESPToken:
		if len(v) >= 1 {
			for _, tok := range v {
				strBytes, ok := tok.Value.([]byte)
				if ok {
					responseData = append(responseData, strBytes...)
				}
			}
		}
	case [][]*RESPToken:
		if len(v) >= 1 {
			for _, tokenList := range v {
				for _, token := range tokenList {
					strBytes, ok := token.Value.([]byte)
					if ok {
						responseData = append(responseData, strBytes...)
					}
				}
			}
		}
	default:
		panic(fmt.Sprintf("encoding unsupported type: %T", tokens))
	}

	return responseData
}
