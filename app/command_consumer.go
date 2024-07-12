package main

import (
	"fmt"
	"net"
)

func commandConsumer(channel <-chan Message, conn net.Conn, dict map[string]string) {
	for {
		msg := <-channel

		if msg.command == nil {
			break
		}

		encoder := NewRESPEncoder()
		parser := NewRESPParser(encoder, dict)
		response := parser.Parse(msg.command)

		if len(response) >= 1 {
			encodedResponse := make([]byte, 0)
			for _, tok := range response {
				strBytes, ok := tok.Value.([]byte)
				if ok {
					encodedResponse = append(encodedResponse, strBytes...)
				}
			}
			conn.Write(encodedResponse)
		} else {
			fmt.Errorf("Did not generate response for command")
			conn.Write([]byte("+OK\r\n"))
		}
	}

}
