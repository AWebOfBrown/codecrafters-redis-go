package main

import (
	"bufio"
	"net"
	"testing"
)

func Test_CommandConsumer(t *testing.T) {
	t.Run("Test a transaction", func(t *testing.T) {

		dict := make(map[string]string)
		commandQueue := make(chan RedisCommandQueueMessage)
		mc := NewTransactionContext()

		client, server := net.Pipe()

		go CommandConsumerController(commandQueue, dict, &mc)
		go commandProducerController(&server, commandQueue)

		reader := bufio.NewReader(client)

		client.Write([]byte("*1\r\n$5\r\nMULTI\r\n"))
		reader.ReadBytes('\n')
		client.Write([]byte("*3\r\n$3\r\nSET\r\n$5\r\napple\r\n$2\r\n67\r\n"))
		reader.ReadBytes('\n')

		client.Write([]byte("*2\r\n$4\r\nINCR\r\n$5\r\napple\r\n"))
		reader.ReadBytes('\n')

		client.Write([]byte("*2\r\n$4\r\nINCR\r\n$9\r\nblueberry\r\n"))
		reader.ReadBytes('\n')

		client.Write([]byte("*2\r\n$3\r\nGET\r\n$9\r\nblueberry\r\n"))
		reader.ReadBytes('\n')

		client.Write([]byte("*1\r\n$4\r\nEXEC\r\n"))
		reader.ReadBytes('\n')
	})
}
