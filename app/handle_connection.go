package main

import (
	"bufio"
	"net"
)

type Message struct {
	command []*RESPToken
}

func handleConnection(conn net.Conn, mp *map[string]string) {
	messageChan := make(chan Message)
	reader := bufio.NewReader(conn)

	go commandProducer(reader, messageChan, conn)
	go commandConsumer(messageChan, conn, *mp)
}
