package main

import (
	"bufio"
	"fmt"
	"net"
	"testing"
)

func Test_CommandConsumer(t *testing.T) {
	t.Run("Test a transaction", func(t *testing.T) {

		dict := make(map[string]string)
		commandQueue := make(chan RedisCommandQueueMessage)
		mc := NewMultiContext()

		client, server := net.Pipe()

		go CommandConsumerController(commandQueue, dict, &mc)
		go commandProducerController(&server, commandQueue)

		client.Write([]byte("*1\r\n$5\r\nMULTI\r\n"))
		client.Write([]byte("*1\r\n$4\r\nEXEC\r\n"))

		reader := bufio.NewReader(client)
		b, _ := reader.ReadBytes('\n')
		b2, _ := reader.ReadBytes('\n')
		b3, _ := reader.ReadBytes('\n')
		fmt.Printf("%s, %s, %s", b, b2, b3)
	})
}
