package main

import (
	"testing"
)

func Test_ParserTest(t *testing.T) {
	t.Run("Parse INCR of a previously SET value", func(t *testing.T) {
		dict := make(map[string]string)
		mc := NewTransactionContext()
		parser := NewRESPParser(dict, &mc)

		setCommand := []*RESPToken{{
			Type:   Array,
			length: 3,
		}, {
			Type:   BulkString,
			Value:  "SET",
			length: 3,
		},
			{
				Type:   BulkString,
				Value:  "foo",
				length: 3,
			},
			{
				Type:  Integer,
				Value: 5,
			}}

		incrCommand := []*RESPToken{{
			Type:   Array,
			length: 2,
		}, {
			Type:   BulkString,
			Value:  "INCR",
			length: 4,
		},
			{
				Type:   BulkString,
				Value:  "foo",
				length: 3,
			}}

		parser.Parse(setCommand, false)
		incResult, _ := parser.Parse(incrCommand, false)

		str := string(incResult.serialiseRESPTokens())
		if str != ":6\r\n" {
			t.Errorf("want: :6\r\n, got %s", str)
		}
	})
}
