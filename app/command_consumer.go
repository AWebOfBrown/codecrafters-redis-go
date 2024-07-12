package main

import "fmt"

func commandConsumer(command []*RESPToken, parser *RESPParser) []*RESPToken {
	response := parser.Parse(command)
	return response
}

func commandConsumerController(queue <-chan RedisCommandQueueMessage, dict map[string]string) {
	encoder := NewRESPEncoder()
	parser := NewRESPParser(encoder, dict)
	for {
		redisCommand := <-queue
		response := commandConsumer(redisCommand.command, parser)
		if len(response) >= 1 {
			encodedResponse := make([]byte, 0)
			for _, tok := range response {
				strBytes, ok := tok.Value.([]byte)
				if ok {
					encodedResponse = append(encodedResponse, strBytes...)
				}
			}
			redisCommand.connection.Write(encodedResponse)
		} else {
			fmt.Errorf("Did not generate response for command")
			redisCommand.connection.Write([]byte("+OK\r\n"))
		}
	}
}
