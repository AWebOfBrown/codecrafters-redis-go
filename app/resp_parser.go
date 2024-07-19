package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type RESPParser struct {
	dict          map[string]string
	multiContext  *MultiContext
	currentClient *net.Conn
}

func NewRESPParser(mp map[string]string, mc *MultiContext) *RESPParser {
	return &RESPParser{
		dict:         mp,
		multiContext: mc,
	}
}

func (p *RESPParser) SetClientConnection(conn *net.Conn) {
	p.currentClient = conn
}

func (p *RESPParser) Parse(tokens []*RESPToken) ([]*RESPToken, error) {
	var response []*RESPToken

	command := tokens[1].Value

	var parsingError error

	if str, ok := command.(string); ok {
		switch strings.ToLower(str) {
		case "echo":
			token, err := NewRESPToken(BulkString, tokens[2].Value.(string))
			parsingError = err
			response = []*RESPToken{token}
		case "ping":
			token, err := NewRESPToken(BulkString, "PONG")
			parsingError = err
			response = []*RESPToken{token}
		case "set":
			p.parseSet(tokens)
			//todo: handle returning set value if requested
			token, e := NewRESPToken(BulkString, "OK")
			parsingError = e
			response = []*RESPToken{token}
		case "get":
			key := tokens[2].Value.(string)
			// todo: Stop assuming this is a string
			value := p.dict[key]
			token, err := NewRESPToken(BulkString, value)
			parsingError = err
			response = []*RESPToken{token}
		case "incr":
			i, err := p.parseIncr(tokens)
			parsingError = err
			token, _ := NewRESPToken(Integer, strconv.Itoa(i))
			response = []*RESPToken{token}
		case "multi":
			token, _ := NewRESPToken(BulkString, "OK")
			p.multiContext.AddTxConnection(p.currentClient)
			response = []*RESPToken{token}
			// Should be handled elsewhere, this case is when exec is called w/o multi first.
		case "exec":
			token, _ := NewRESPToken(Error, "EXEC without MULTI")
			response = []*RESPToken{token}
		default:
			panic(fmt.Errorf("encountered unhandled/unsupported command %s", command))
		}
	}

	return response, parsingError
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
