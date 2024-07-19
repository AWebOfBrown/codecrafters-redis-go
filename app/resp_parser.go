package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type RESPParser struct {
	dict            map[string]string
	txContext       *TransactionContext
	currentClientID string
}

func NewRESPParser(mp map[string]string, mc *TransactionContext) *RESPParser {
	return &RESPParser{
		dict:      mp,
		txContext: mc,
	}
}

func (p *RESPParser) SetClientConnection(id string) {
	p.currentClientID = id
}

func (p *RESPParser) Parse(tokens []*RESPToken, isTransactional bool) (RESPResponse, error) {
	var response RESPResponse

	if isTransactional {
		return p.parseTransaction(tokens)
	}

	command := tokens[1].Value

	var parsingError error

	if str, ok := command.(string); ok {
		switch strings.ToLower(str) {
		case "echo":
			token, err := NewRESPToken(BulkString, tokens[2].Value.(string))
			parsingError = err
			response = NewIndividualRESPResponse([]*RESPToken{token})
		case "ping":
			token, err := NewRESPToken(BulkString, "PONG")
			parsingError = err
			response = NewIndividualRESPResponse([]*RESPToken{token})
		case "set":
			p.parseSet(tokens)
			//todo: handle returning set value if requested
			token, e := NewRESPToken(BulkString, "OK")
			parsingError = e
			response = NewIndividualRESPResponse([]*RESPToken{token})
		case "get":
			key := tokens[2].Value.(string)
			// todo: Stop assuming this is a string
			value := p.dict[key]
			token, err := NewRESPToken(BulkString, value)
			parsingError = err
			response = NewIndividualRESPResponse([]*RESPToken{token})
		case "incr":
			i, err := p.parseIncr(tokens)
			parsingError = err
			token, _ := NewRESPToken(Integer, strconv.Itoa(i))
			response = NewIndividualRESPResponse([]*RESPToken{token})
		case "multi":
			token, _ := NewRESPToken(BulkString, "OK")
			p.txContext.RegisterActiveClientTX(p.currentClientID)
			response = NewIndividualRESPResponse([]*RESPToken{token})
			// Should be handled elsewhere, this case is when exec is called w/o multi first.
		case "exec":
			token, _ := NewRESPToken(Error, "EXEC without MULTI")
			response = NewIndividualRESPResponse([]*RESPToken{token})
		case "discard":
			token, _ := NewRESPToken(Error, "DISCARD without MULTI")
			response = NewIndividualRESPResponse([]*RESPToken{token})
		default:
			panic(fmt.Errorf("encountered unhandled/unsupported command %s", command))
		}
	}

	p.currentClientID = ""
	return response, parsingError
}

func (p *RESPParser) parseTransaction(tokens []*RESPToken) (RESPResponse, error) {
	cmdToken := tokens[1].Value.(string)

	var response RESPResponse

	switch strings.ToLower(cmdToken) {
	case "discard":
		p.txContext.RemoveClientTX(p.currentClientID)
		token, _ := NewRESPToken(String, "OK")
		response = NewIndividualRESPResponse([]*RESPToken{token})
	case "multi":
		token, _ := NewRESPToken(String, "OK")
		response = *NewRESPResponseList([][]*RESPToken{{token}})
	case "exec":
		execResponses := make([][]*RESPToken, 0)
		queuedCommands := p.txContext.GetQueuedCommands(p.currentClientID)
		qtyQueuedCommands := len(queuedCommands)
		if qtyQueuedCommands == 0 {
			emptyArray, _ := NewRESPToken(Array, "0")
			execResponses = append(execResponses, []*RESPToken{emptyArray})
			response = *NewRESPResponseList(execResponses)
		} else {
			for _, queuedCommand := range queuedCommands {
				result, err := p.Parse(queuedCommand, false)
				if err != nil {
					//todo: Handle better.
					errToken, _ := NewRESPToken(Error, err.Error())
					execResponses = append(execResponses, []*RESPToken{errToken})
				} else {
					if singleResponse, ok := result.(*IndividualRESPResponse); ok {
						execResponses = append(execResponses, singleResponse.tokens)
					}
				}
			}
			response = *NewRESPResponseList(execResponses)
			if responseList, ok := response.(RESPResponseList); ok {
				lengthOfResponse := len(responseList.tokens)
				leadingArrayToken, _ := NewRESPToken(Array, strconv.Itoa(lengthOfResponse))
				responseList.tokens = append([][]*RESPToken{{leadingArrayToken}}, responseList.tokens...)
				response = responseList
			} else {
				panic("Should not be returning an exec with list of commands that is not a RESPResponseList")
			}
		}
		p.txContext.RemoveClientTX(p.currentClientID)
	default:
		//todo: Validate command before enqueue
		p.txContext.EnqueueCommand(p.currentClientID, tokens)
		token, err := NewRESPToken(String, "QUEUED")
		if err != nil {
			fmt.Printf("%s", err)
		}
		tokenList := [][]*RESPToken{{token}}
		response = *NewRESPResponseList(tokenList)
	}

	return response, nil
}

func (p *RESPParser) parseIncr(tokens []*RESPToken) (int, error) {
	key := tokens[2].Value.(string)

	currVal := p.dict[key]

	if currVal == "" {
		p.dict[key] = strconv.Itoa(1)
		return 1, nil
	}

	currValAsInteger, err := strconv.Atoi(currVal)
	if err != nil {
		return -1, fmt.Errorf("value is not an integer or out of range")
	}
	//todo: handle incrementing strings (error)
	currValAsInteger = 1 + currValAsInteger
	p.dict[key] = strconv.Itoa(currValAsInteger)

	return currValAsInteger, nil
}

func (p *RESPParser) parseSet(tokens []*RESPToken) string {
	key := tokens[2].Value.(string)

	var value string

	switch v := tokens[3].Value.(type) {
	case int:
		value = strconv.Itoa(v)
	case string:
		value = v
	default:
		//todo: handle
		panic("Unrecognised type")
	}

	p.dict[key] = value

	includesPx := false
	var pxValue int
	//todo: extract all args
	for i := 2; i < len(tokens); i++ {
		if str, ok := tokens[i].Value.(string); ok {
			if strings.ToLower(str) == "px" {
				includesPx = true
				if num, ok := tokens[i+1].Value.(string); ok {
					strNum, err := strconv.Atoi(num)
					if err != nil {
						includesPx = false
						break
					}
					pxValue = strNum
				}

				break
			}
		}
	}

	if includesPx && pxValue >= 1 {
		go func(dict *map[string]string, key string) {
			timer := time.NewTimer(time.Millisecond * time.Duration(pxValue))
			<-timer.C
			(*dict)[key] = ""
			timer.Stop()
		}(&p.dict, key)
	}

	return value
}
