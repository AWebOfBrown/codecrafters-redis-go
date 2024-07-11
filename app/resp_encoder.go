package main

import "fmt"

type RESPEncoder struct {
}

func NewRESPEncoder() *RESPEncoder {
	return &RESPEncoder{}
}

func (re *RESPEncoder) Encode(tokens []*RESPToken) []*RESPToken {
	for _, tok := range tokens {
		switch tok.Type {
		case BulkString:
			re.encodeBulkString(tok)
		case String:
			re.encodeString(tok)
		default:
			panic(fmt.Sprintf("Unknown token type: %s", tok.Type))
		}
	}

	return tokens
}

func (re *RESPEncoder) encodeBulkString(token *RESPToken) {
	len := token.length
	token.Value = []byte(fmt.Sprintf("$%d\r\n%s\r\n", len, token.Value))
}

func (re *RESPEncoder) encodeString(token *RESPToken) {
	token.Value = []byte("+OK\r\n")
}
