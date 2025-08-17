package commands

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/AWebOfBrown/codecrafters-http-server-go/internal/resp"
)

type RedisCommandQueueMessage struct {
	command    []*resp.RESPToken
	connection *net.Conn
}

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

func CommandProducerController(conn *net.Conn, queue chan<- RedisCommandQueueMessage) {
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
