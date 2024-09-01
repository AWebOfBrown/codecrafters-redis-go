package main

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func commandProducer(lexer *resp.RESPLexer) ([]*resp.RESPToken, error) {
	tokens, err := lexer.ProduceTokens()

	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		fmt.Printf("Error %s", err)
		return nil, err
	}

	return tokens, nil
}

func commandProducerController(conn *net.Conn, queue chan<- RedisCommandQueueMessage) {
	reader := bufio.NewReader(*conn)
	lexer := resp.NewRESPLexer(reader)

	for {
		command, err := commandProducer(lexer)

		if err != nil {
			if err == io.EOF {
				(*conn).Close()
				return
			}
			fmt.Printf("Unknown error %s", err)
			(*conn).Close()
			return
		}

		queue <- RedisCommandQueueMessage{
			command:    command,
			connection: conn,
		}
	}
}
