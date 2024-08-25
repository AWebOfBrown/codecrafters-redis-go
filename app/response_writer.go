package main

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func responseWriter(responseTokens resp.RESPResponse, closeConn bool, conn *net.Conn) {
	responseBytes := responseTokens.SerialiseRESPTokens()
	(*conn).Write(responseBytes)
	if closeConn {
		(*conn).Close()
	}
}
