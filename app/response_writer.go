package main

import (
	"net"
)

func responseWriter(responseTokens RESPResponse, closeConn bool, conn *net.Conn) {
	responseBytes := responseTokens.serialiseRESPTokens()
	(*conn).Write(responseBytes)
	if closeConn {
		(*conn).Close()
	}
}
