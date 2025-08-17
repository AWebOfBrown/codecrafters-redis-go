package commands

import (
	"bufio"
	"net"
	"testing"

	"github.com/AWebOfBrown/codecrafters-http-server-go/internal/resp"
)

func makeServer() (net.Conn, net.Conn) {
	dict := make(map[string]interface{})
	commandQueue := make(chan RedisCommandQueueMessage)
	mc := resp.NewTransactionContext()

	client, server := net.Pipe()

	go CommandConsumerController(commandQueue, dict, &mc)
	go CommandProducerController(&server, commandQueue)
	return client, server
}

func Test_CommandConsumer(t *testing.T) {
	t.Run("Test a transaction", func(t *testing.T) {

		client, server := makeServer()
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

		client.Close()
		server.Close()
	})
}
