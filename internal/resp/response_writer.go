package resp

import (
	"net"
)

func ResponseWriter(responseTokens RESPResponse, closeConn bool, conn *net.Conn) {
	responseBytes := responseTokens.SerialiseRESPTokens()
	(*conn).Write(responseBytes)
	if closeConn {
		(*conn).Close()
	}
}
