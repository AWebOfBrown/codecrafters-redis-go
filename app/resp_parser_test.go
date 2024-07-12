package main

import (
	"testing"
)

func Test_ParserTest(t *testing.T) {
	t.Run("Parse INCR", func(t *testing.T) {
		dict := make(map[string]string)
		parser := NewRESPParser(dict)

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

		parser.Parse(setCommand)
		incResult := parser.Parse(incrCommand)

		if incResult[0].Value != 6 {
			t.Error("Whoops")
		}
	})
}
