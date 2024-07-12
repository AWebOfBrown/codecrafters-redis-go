package main

import (
	"bufio"
	"net"
	"sync"
	"testing"
)

func TestLexer_Test(t *testing.T) {
	t.Run("Test parsing an array of bulk strings", func(t *testing.T) {
		server, client := net.Pipe()
		reader := bufio.NewReader(server)
		lexer := NewRESPLexer(reader)

		var wg sync.WaitGroup

		bulkString := "*1\r\n$5\r\nhello\r\n"
		wg.Add(1)
		go func() {
			defer client.Close()
			client.Write([]byte(bulkString))
			wg.Done()
		}()
		tokens, err := lexer.ProduceTokens()

		server.Write([]byte("+OK\r\n"))
		wg.Wait()

		if err != nil {
			t.Errorf("Failed with err %s", err)
		}

		if len(tokens) != 2 {
			t.Errorf("Expected 2 tokens, got %d", len(tokens))
		}

		tokenOneType := tokens[0].Type

		tokenTwoType := tokens[1].Type
		tokenTwoValue := tokens[1].Value

		if tokenOneType != Array {
			t.Errorf("Expected array, got %s", tokenOneType)
		}
		if tokenTwoType != BulkString {
			t.Errorf("Expected bulk string, got %s", tokenTwoType)
		}

		if tokenTwoValue != "hello" {
			t.Errorf("Expected hello, got %s", tokenTwoValue)
		}
	})

}
