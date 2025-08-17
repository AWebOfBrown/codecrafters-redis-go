package resp

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/AWebOfBrown/codecrafters-http-server-go/internal/stream"
)

type RESPParser struct {
	dict            map[string]interface{}
	txContext       *TransactionContext
	currentClientID string
}

func NewRESPParser(mp map[string]interface{}, mc *TransactionContext) *RESPParser {
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
			token, err := NewRESPToken(String, "PONG")
			parsingError = err
			response = NewIndividualRESPResponse([]*RESPToken{token})
		case "set":
			p.parseSet(tokens)
			//todo: handle returning set value if requested
			token, e := NewRESPToken(BulkString, "OK")
			parsingError = e
			response = NewIndividualRESPResponse([]*RESPToken{token})
		case "get":
			key := tokens[2].Value
			k, ok := key.(string)
			if !ok {
				return nil, nil
			}
			value := p.dict[k]

			if value != nil {
				token, err := NewRESPToken(BulkString, value.(string))
				parsingError = err
				response = NewIndividualRESPResponse([]*RESPToken{token})
			} else {
				nullBulkString, _ := NewRESPToken(BulkString, "")
				return NewIndividualRESPResponse([]*RESPToken{nullBulkString}), nil
			}
		case "type":
			key := tokens[2].Value.(string)
			value, ok := p.dict[key]
			var token RESPToken
			if ok {
				_, isStream := value.(*stream.Stream)
				if isStream {
					t, _ := NewRESPToken(String, "stream")
					token = *t
				} else {
					t, _ := NewRESPToken(String, "string")
					token = *t
				}
			} else {
				t, _ := NewRESPToken(String, "none")
				token = *t
			}
			response = NewIndividualRESPResponse([]*RESPToken{&token})
		case "incr":
			i, err := p.parseIncr(tokens)
			parsingError = err
			token, _ := NewRESPToken(Integer, i)
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
		case "xadd":
			key := tokens[2].Value.(string)
			streamId := tokens[3].Value.(string)
			mapOfValuesToInsert := make(map[string]interface{}, 0)
			for i := 4; i < len(tokens); i += 2 {
				key := tokens[i].Value.(string)
				value := tokens[i+1].Value.(string)
				mapOfValuesToInsert[key] = value
			}

			var targetStream *stream.Stream
			if p.dict[key] == nil {
				targetStream = stream.NewStream()
				p.dict[key] = targetStream
			} else {
				s, ok := p.dict[key].(stream.Stream)
				if !ok {
					return nil, fmt.Errorf("tried to add to non-stream value")
				}
				targetStream = &s
			}
			targetStream.Insert(streamId, mapOfValuesToInsert)
			token, _ := NewRESPToken(String, streamId)
			return NewIndividualRESPResponse([]*RESPToken{token}), nil

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

func (p *RESPParser) parseIncr(tokens []*RESPToken) (string, error) {
	key := tokens[2].Value.(string)

	currVal := p.dict[key]

	if currVal == nil {
		p.dict[key] = "1"
		return "1", nil
	}

	currValAsStr, ok := currVal.(string)
	if !ok {
		//todo: error
		return "", nil
	}

	currentNumber, error := strconv.Atoi(currValAsStr)
	if error == nil {
		newVal := currentNumber + 1
		newValString := strconv.Itoa(newVal)
		p.dict[key] = newValString
		return newValString, nil
	} else {
		return "", fmt.Errorf("value is not an integer or out of range")
	}
}

func (p *RESPParser) parseSet(tokens []*RESPToken) string {
	key := tokens[2].Value.(string)

	value := tokens[3].Value.(string)

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
		go func(dict *map[string]interface{}, key string) {
			timer := time.NewTimer(time.Millisecond * time.Duration(pxValue))
			<-timer.C
			(*dict)[key] = ""
			timer.Stop()
		}(&p.dict, key)
	}

	return value
}
