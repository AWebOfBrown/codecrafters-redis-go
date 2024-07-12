package main

import (
	"strconv"
	"strings"
	"time"
)

type RESPParser struct {
	lexer   *RESPLexer
	encoder *RESPEncoder
	dict    map[string]string
}

func NewRESPParser(mp map[string]string) *RESPParser {
	return &RESPParser{
		dict: mp,
	}
}

func (p *RESPParser) Parse(tokens []*RESPToken) []*RESPToken {
	var response []*RESPToken

	command := tokens[1].Value
	if str, ok := command.(string); ok {
		switch strings.ToLower(str) {
		case "echo":
			response = []*RESPToken{tokens[2]}
		case "ping":
			response = []*RESPToken{
				{
					Type:   "$",
					Value:  "PONG",
					length: 4,
				},
			}
		case "set":
			p.parseSet(tokens)
			//todo: handle returning set value if requested
			response = []*RESPToken{{Type: "+", Value: "OK"}}
		case "get":
			key := tokens[2].Value.(string)
			value := p.dict[key]
			response = []*RESPToken{{Type: "$", Value: value}}
		case "incr":
			response = p.parseIncr(tokens)
		case "multi":
			response = []*RESPToken{{Type: "$", Value: "OK"}}
		}
	}

	return response
}

func (p *RESPParser) parseIncr(tokens []*RESPToken) []*RESPToken {
	key := tokens[2].Value.(string)

	currVal := p.dict[key]

	if currVal == "" {
		p.dict[key] = strconv.Itoa(1)
		return []*RESPToken{{Value: 1, Type: Integer}}
	}

	currValAsInteger, err := strconv.Atoi(currVal)
	if err != nil {
		return []*RESPToken{{
			Value: "value is not an integer or out of range",
			Type:  Error,
		}}
	}
	//todo: handle incrementing strings (error)
	currValAsInteger = 1 + currValAsInteger
	p.dict[key] = strconv.Itoa(currValAsInteger)

	return []*RESPToken{{Value: currValAsInteger, Type: Integer}}
}

func (p *RESPParser) parseSet(tokens []*RESPToken) []*RESPToken {
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

	return []*RESPToken{
		{
			Type:   "$",
			Value:  p.dict[key],
			length: len(value),
		},
	}
}
