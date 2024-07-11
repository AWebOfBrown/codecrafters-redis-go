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
}

func NewRESPParser(encoder *RESPEncoder) *RESPParser {
	return &RESPParser{
		encoder: encoder,
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
