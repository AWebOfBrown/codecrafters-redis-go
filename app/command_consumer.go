package main

import "fmt"

func commandConsumer(command []*RESPToken, parser *RESPParser) []*RESPToken {
	response := parser.Parse(command)
	return response
}

func commandConsumerController(queue <-chan RedisCommandQueueMessage, dict map[string]string) {
	encoder := NewRESPEncoder()
	parser := NewRESPParser(dict)
	for {
		redisCommand := <-queue
		response := commandConsumer(redisCommand.command, parser)
		tokensWithEncodedValues := encoder.Encode(response)

		if len(tokensWithEncodedValues) >= 1 {
			encodedAggregatedResponse := make([]byte, 0)
			for _, tok := range tokensWithEncodedValues {
				strBytes, ok := tok.Value.([]byte)
				if ok {
					encodedAggregatedResponse = append(encodedAggregatedResponse, strBytes...)
				}
			}
			redisCommand.connection.Write(encodedAggregatedResponse)
		} else {
			fmt.Errorf("Did not generate response for command")
			redisCommand.connection.Write([]byte("+OK\r\n"))
		}
	}
}
