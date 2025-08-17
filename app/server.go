package main

import (
	"fmt"
	"net"
	"os"

	"github.com/AWebOfBrown/codecrafters-http-server-go/internal/commands"
	"github.com/AWebOfBrown/codecrafters-http-server-go/internal/resp"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	fmt.Println("listening on 6379")

	dict := make(map[string]interface{})
	commandQueue := make(chan commands.RedisCommandQueueMessage)

	tc := resp.NewTransactionContext()
	// single thread for handling writes/reads to dict
	go commands.CommandConsumerController(commandQueue, dict, &tc)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go commands.CommandProducerController(&conn, commandQueue)
	}
}
