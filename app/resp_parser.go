package main

import "strings"

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
			key := tokens[2].Value.(string)
			value := tokens[3].Value.(string)
			p.dict[key] = value
			response = p.encoder.Encode([]*RESPToken{{Type: "+", Value: "OK"}})
		case "get":
			key := tokens[2].Value.(string)
			value := p.dict[key]
			response = p.encoder.Encode([]*RESPToken{
				{
					Type:   "$",
					Value:  p.dict[key],
					length: len(value),
				},
			})
		}
	}

	return response
}

func (p *RESPParser) parseEcho(echoArg *RESPToken) []*RESPToken {
	echoVal := echoArg.Value
	response := append(make([]*RESPToken, 0), echoArg)

	if _, ok := echoVal.(string); ok {
		response = p.encoder.Encode(response)
	}

	return response
}
