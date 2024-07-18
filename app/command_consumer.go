package main

import (
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
		response = []*RESPToken{{Type: String, Value: "+OK\r\n"}}
	case "exec":
		execResponses := make([][]*RESPToken, 0)
		for _, queuedCommand := range mc.GetQueuedCommands(conn) {
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
		mc.RemoveTxConnection(conn)
	default:
		//todo: Validate command before enqueue
		mc.EnqueueCommand(conn, command)
		response = []*RESPToken{{Type: String, Value: "+QUEUED\r\n"}}
	}

	return response, nil
}

func commandConsumerController(queue <-chan RedisCommandQueueMessage, dict map[string]string, multiContext *MultiContext) {
	parser := NewRESPParser(dict)
	for {
		redisCommand := <-queue

		// Could be singular, or multi-responses in case of a TX
		var response any

		if multiContext.CheckActiveTX(&redisCommand.connection) {
			result, err := handleTransaction(redisCommand.command, &redisCommand.connection, multiContext, parser)
			if err != nil {
				redisCommand.connection.Write([]byte("+ERR\r\n"))
				redisCommand.connection.Close()
				return
			}
			response = result
		} else {
			result, err := commandConsumer(redisCommand.command, parser)
			if err != nil {
				responseToken, _ := NewRESPToken(Error, err.Error())
				redisCommand.connection.Write([]byte(responseToken.Value.([]byte)))
				redisCommand.connection.Close()
				return
			}
			response = result
		}

		responseBytes := serialiseTokens(response)
		redisCommand.connection.Write(responseBytes)
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
		panic("encoding unsupported type")
	}

	return responseData
}
