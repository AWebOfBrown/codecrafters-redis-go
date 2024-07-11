package main

import "net"

func commandConsumer(channel <-chan Message, conn net.Conn) {
	for {
		msg := <-channel

		if msg.command == nil {
			break
		}

		encoder := NewRESPEncoder()
		parser := NewRESPParser(encoder)
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
		}
	}

}
