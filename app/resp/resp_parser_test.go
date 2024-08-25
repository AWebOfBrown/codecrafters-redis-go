package resp

import (
	"testing"
)

func Test_ParserTest(t *testing.T) {
	t.Run("Parse INCR of a previously SET value", func(t *testing.T) {
		dict := make(map[string]interface{})
		mc := NewTransactionContext()
		parser := NewRESPParser(dict, &mc)

		setCommand := []*RESPToken{{
			Type:   Array,
			Length: 3,
		}, {
			Type:   BulkString,
			Value:  "SET",
			Length: 3,
		},
			{
				Type:   BulkString,
				Value:  "foo",
				Length: 3,
			},
			{
				Type:  Integer,
				Value: 5,
			}}

		incrCommand := []*RESPToken{{
			Type:   Array,
			Length: 2,
		}, {
			Type:   BulkString,
			Value:  "INCR",
			Length: 4,
		},
			{
				Type:   BulkString,
				Value:  "foo",
				Length: 3,
			}}

		parser.Parse(setCommand, false)
		incResult, _ := parser.Parse(incrCommand, false)

		str := string(incResult.SerialiseRESPTokens())
		if str != ":6\r\n" {
			t.Errorf("want: :6\r\n, got %s", str)
		}
	})
}
