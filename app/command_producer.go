package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func commandProducer(reader *bufio.Reader, ch chan<- Message, conn net.Conn) error {
	lexer := NewRESPLexer(reader)
	for {
		err := lexer.ProduceTokens(ch)
		if err != nil {
			ch <- Message{command: nil}
			if err == io.EOF {
				conn.Close()
				return err
			}
			fmt.Errorf("Error %s", err)
			conn.Close()
			return err
		}
	}
}
