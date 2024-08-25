package resp

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
		case Integer:
			re.encodeInteger(tok)
		case Error:
			re.encodeError(tok)
		default:
			panic(fmt.Sprintf("Unhandled token type: %s", tok.Type))
		}
	}

	return tokens
}

func (re *RESPEncoder) encodeError(token *RESPToken) {
	token.Value = []byte(fmt.Sprintf("-ERR %s\r\n", token.Value))
}

func (re *RESPEncoder) encodeInteger(token *RESPToken) {
	token.Value = []byte(fmt.Sprintf(":%d\r\n", token.Value))
}

func (re *RESPEncoder) encodeBulkString(token *RESPToken) {
	len := len(token.Value.(string))

	if len == 0 {
		token.Value = []byte("$-1\r\n")
		return
	}

	token.Value = []byte(fmt.Sprintf("$%d\r\n%s\r\n", len, token.Value))
}

func (re *RESPEncoder) encodeString(token *RESPToken) {
	token.Value = []byte("+OK\r\n")
}
