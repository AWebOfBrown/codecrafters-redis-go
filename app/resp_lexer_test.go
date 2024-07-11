package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestRESPLexer_ProduceTokens(t *testing.T) {
	server, client := net.Pipe()

	input := "*2\r\n$4\r\nECHO\r\n$5\r\nmango\r\n"
	encodedInput := []byte(input)

	var wg sync.WaitGroup

	// Client
	go func(conn net.Conn) {
		defer wg.Done()
		r := bufio.NewReader(conn)
		num, err := conn.Write(encodedInput)
		fmt.Printf("num %d, err %s", num, err)
		msg, err := r.ReadBytes('\n')
		msg2, err := r.ReadBytes('\n')
		fmt.Printf("HEREHERHEHREHR %s, %s, err %s", string(msg), string(msg2), err)
		conn.Close()
	}(client)
	wg.Add(1)

	handleConnection(server)
	wg.Wait()
}
