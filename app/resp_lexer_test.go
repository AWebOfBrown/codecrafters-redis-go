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
	//"*2\r\n$4\r\nECHO\r\n$5\r\nmango\r\n"
	// input := "*3\r\n$3\r\nset\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	// input := "*3\r\n$3\r\nSET\r\n$6\r\norange\r\n$5\r\napple\r\n"
	input := "*5\r\n$3\r\nSET\r\n$5\r\ngrape\r\n$6\r\nbanana\r\n$2\r\npx\r\n$3\r\n100\r\n"
	inputGet := "*2\r\n$3\r\nGET\r\n$5\r\ngrape\r\n"
	encodedInput := []byte(input)
	encodedGet := []byte(inputGet)
	dict := make(map[string]string)

	var wg sync.WaitGroup

	// Client
	go func(conn net.Conn) {
		defer wg.Done()
		r := bufio.NewReader(conn)
		num, err := conn.Write(encodedInput)
		fmt.Printf("num %d, err %s", num, err)
		msg, err := r.ReadBytes('\n')
		fmt.Printf("HEREHERHEHREHR %s, err %s", string(msg), err)
		conn.Write(encodedGet)
		msg3, _ := r.ReadBytes('\n')
		msg4, _ := r.ReadBytes('\n')
		fmt.Printf("HEREHERHEHREHR %s, %s, err %s", string(msg3), string(msg4), err)
		conn.Close()
	}(client)
	wg.Add(1)

	handleConnection(server, &dict)
	wg.Wait()
}
