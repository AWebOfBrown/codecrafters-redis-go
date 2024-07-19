package main

import (
	"io"
	"net"
)

type MultiContext struct {
	dict map[string][][]*RESPToken
}

func NewMultiContext() MultiContext {
	dict := make(map[string][][]*RESPToken)
	return MultiContext{
		dict: dict,
	}
}

func (mc *MultiContext) CheckActiveTX(conn *net.Conn) bool {
	key := (*conn).RemoteAddr().String()
	_, ok := mc.dict[key]
	return ok
}

func (mc *MultiContext) AddTxConnection(conn *net.Conn) {
	key := (*conn).RemoteAddr().String()
	mc.dict[key] = make([][]*RESPToken, 0)
}

func (mc *MultiContext) EnqueueCommand(conn *net.Conn, cmd []*RESPToken) error {
	key := (*conn).RemoteAddr().String()
	val := mc.dict[key]
	if val == nil {
		return io.EOF
	}
	mc.dict[key] = append(val, cmd)
	return nil
}

func (mc *MultiContext) GetQueuedCommands(conn *net.Conn) [][]*RESPToken {
	key := (*conn).RemoteAddr().String()
	return mc.dict[key]
}

func (mc *MultiContext) RemoveTxConnection(conn *net.Conn) {
	key := (*conn).RemoteAddr().String()
	delete(mc.dict, key)
}
