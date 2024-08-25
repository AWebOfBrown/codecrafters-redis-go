package main

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func CommandConsumerController(queue <-chan RedisCommandQueueMessage, dict map[string]interface{}, transactionContext *resp.TransactionContext) {
	parser := resp.NewRESPParser(dict, transactionContext)
	for {
		redisCommand := <-queue
		conn := *redisCommand.connection

		// Could be singular list, or list of list of token responses in case of a TX
		var responseTokens resp.RESPResponse
		var shouldCloseConnection bool

		clientID := conn.RemoteAddr().String()
		parser.SetClientConnection(clientID)
		isActiveTx := transactionContext.CheckActiveTX(clientID)

		result, err := parser.Parse(redisCommand.command, isActiveTx)
		if err != nil {
			errorToken, _ := resp.NewRESPToken(resp.Error, err.Error())
			responseTokens = resp.NewIndividualRESPResponse([]*resp.RESPToken{errorToken})
			shouldCloseConnection = true
		} else {
			responseTokens = result
			shouldCloseConnection = false
		}

		go responseWriter(responseTokens, shouldCloseConnection, &conn)
	}
}
