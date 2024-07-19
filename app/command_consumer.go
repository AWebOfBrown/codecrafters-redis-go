package main

func CommandConsumerController(queue <-chan RedisCommandQueueMessage, dict map[string]string, transactionContext *TransactionContext) {
	parser := NewRESPParser(dict, transactionContext)
	for {
		redisCommand := <-queue
		conn := *redisCommand.connection

		// Could be singular list, or list of list of token responses in case of a TX
		var responseTokens RESPResponse
		var shouldCloseConnection bool

		clientID := conn.RemoteAddr().String()
		parser.SetClientConnection(clientID)
		isActiveTx := transactionContext.CheckActiveTX(clientID)

		result, err := parser.Parse(redisCommand.command, isActiveTx)
		if err != nil {
			errorToken, _ := NewRESPToken(Error, err.Error())
			responseTokens = NewIndividualRESPResponse([]*RESPToken{errorToken})
			shouldCloseConnection = true
		} else {
			responseTokens = result
			shouldCloseConnection = false
		}

		go responseWriter(responseTokens, shouldCloseConnection, &conn)
	}
}
