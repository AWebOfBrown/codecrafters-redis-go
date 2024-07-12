package main

import (
	"fmt"
	"net"
	"os"
)

type RedisCommandQueueMessage struct {
	command    []*RESPToken
	connection net.Conn
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	dict := make(map[string]string)
	commandQueue := make(chan RedisCommandQueueMessage)

	// single thread for handling writes/reads to dict
	go commandConsumerController(commandQueue, dict)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go commandProducerController(conn, commandQueue)
	}
}
