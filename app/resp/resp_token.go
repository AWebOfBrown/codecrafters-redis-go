package resp

import (
	"fmt"
	"strconv"
)

type RedisDataType string

const (
	String         = "+"
	Error          = "-"
	Integer        = ":"
	BulkString     = "$"
	Array          = "*"
	Null           = "_"
	Boolean        = "#"
	Double         = ","
	BigNumber      = "("
	BulkError      = "!"
	VerbatimString = "="
	Map            = "%"
	Set            = "~"
	Push           = ">"
)

type RESPToken struct {
	Type   RedisDataType
	Value  interface{}
	Length int
}

func NewRESPToken(rdType RedisDataType, value string) (*RESPToken, error) {
	var token *RESPToken

	switch rdType {
	case Null:
		v := encodeNullValue()
		token = &RESPToken{
			Type:  Null,
			Value: v,
		}
	case BulkString:
		v := encodeBulkStringValue(value)
		token = &RESPToken{
			Type:  BulkString,
			Value: v,
		}
	case String:
		token = &RESPToken{
			Type:  String,
			Value: encodeStringValue(value),
		}
	case Integer:
		int, e := strconv.Atoi(value)
		if e != nil {
			return nil, e
		}
		token = &RESPToken{
			Type:  Integer,
			Value: encodeIntegerValue(int),
		}
	case Error:
		token = &RESPToken{
			Type:  Error,
			Value: encodeErrorValue(value),
		}
	case Array:
		token = &RESPToken{
			Type:  Array,
			Value: []byte(fmt.Sprintf("*%s\r\n", value)),
		}
	default:
		panic(fmt.Sprintf("Unhandled token type: %s", rdType))
	}

	return token, nil
}

func encodeNullValue() []byte {
	return []byte("_\r\n")
}

func encodeErrorValue(value string) []byte {
	return []byte(fmt.Sprintf("-ERR %s\r\n", value))
}

func encodeIntegerValue(value int) []byte {
	return []byte(fmt.Sprintf(":%d\r\n", value))
}

func encodeBulkStringValue(value string) []byte {
	len := len(value)

	if len == 0 {
		return []byte("$-1\r\n")
	}

	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len, value))
}

func encodeStringValue(value string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", value))
}
