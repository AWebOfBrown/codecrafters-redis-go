package main

import (
	"io"
	"net"
)

type MultiContext struct {
	dict map[*net.Conn][][]*RESPToken
}

func NewMultiContext() MultiContext {
	dict := make(map[*net.Conn][][]*RESPToken)
	return MultiContext{
		dict: dict,
	}
}

func (mc *MultiContext) CheckActiveTX(conn *net.Conn) bool {
	isTx := mc.dict[conn] != nil
	return isTx
}

func (mc *MultiContext) AddTxConnection(conn *net.Conn) {
	mc.dict[conn] = make([][]*RESPToken, 0)
}

func (mc *MultiContext) EnqueueCommand(conn *net.Conn, cmd []*RESPToken) error {
	val := mc.dict[conn]
	if val == nil {
		return io.EOF
	}
	mc.dict[conn] = append(val, cmd)
	return nil
}

func (mc *MultiContext) GetQueuedCommands(conn *net.Conn) [][]*RESPToken {
	return mc.dict[conn]
}

func (mc *MultiContext) RemoveTxConnection(c *net.Conn) {
	delete(mc.dict, c)
}
