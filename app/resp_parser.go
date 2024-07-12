package main

import (
	"strconv"
	"strings"
	"time"
)

type RedisCommand string

const (
	ECHO = "ECHO"
	SET  = "SET"
)

type RESPParser struct {
	lexer   *RESPLexer
	encoder *RESPEncoder
	dict    map[string]string
}

func NewRESPParser(encoder *RESPEncoder, mp map[string]string) *RESPParser {
	return &RESPParser{
		encoder: encoder,
		dict:    mp,
	}
}

func (p *RESPParser) Parse(tokens []*RESPToken) []*RESPToken {
	var response []*RESPToken

	command := tokens[1].Value
	if str, ok := command.(string); ok {
		switch strings.ToLower(str) {
		case "echo":
			response = p.parseEcho(tokens[2])
		case "ping":
			response = p.encoder.Encode([]*RESPToken{
				{
					Type:   "$",
					Value:  "PONG",
					length: 4,
				},
			})
		case "set":
			p.parseSet(tokens)
			//todo: handle returning set value if requested
			response = p.encoder.Encode([]*RESPToken{{Type: "+", Value: "OK"}})
		case "get":
			key := tokens[2].Value.(string)
			value := p.dict[key]
			response = p.encoder.Encode([]*RESPToken{{Type: "$", Value: value}})
		}
	}

	return response
}

func (p *RESPParser) parseSet(tokens []*RESPToken) []*RESPToken {
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
		go func(dict *map[string]string, key string) {
			timer := time.NewTimer(time.Millisecond * time.Duration(pxValue))
			<-timer.C
			(*dict)[key] = ""
			timer.Stop()
		}(&p.dict, key)
	}

	return p.encoder.Encode([]*RESPToken{
		{
			Type:   "$",
			Value:  p.dict[key],
			length: len(value),
		},
	})
}

func (p *RESPParser) parseEcho(echoArg *RESPToken) []*RESPToken {
	echoVal := echoArg.Value
	response := append(make([]*RESPToken, 0), echoArg)

	if _, ok := echoVal.(string); ok {
		response = p.encoder.Encode(response)
	}

	return response
}
