package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type RedisCommandQueueMessage struct {
	command    []*resp.RESPToken
	connection *net.Conn
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	dict := make(map[string]interface{})
	commandQueue := make(chan RedisCommandQueueMessage)

	tc := resp.NewTransactionContext()
	// single thread for handling writes/reads to dict
	go CommandConsumerController(commandQueue, dict, &tc)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go commandProducerController(&conn, commandQueue)
	}
}
