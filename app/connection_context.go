package main

import (
	"io"
)

type TransactionContext struct {
	dict map[string][][]*RESPToken
}

func NewTransactionContext() TransactionContext {
	dict := make(map[string][][]*RESPToken)
	return TransactionContext{
		dict: dict,
	}
}

func (tc *TransactionContext) CheckActiveTX(id string) bool {
	_, ok := tc.dict[id]
	return ok
}

func (tc *TransactionContext) RegisterActiveClientTX(id string) {
	tc.dict[id] = make([][]*RESPToken, 0)
}

func (tc *TransactionContext) EnqueueCommand(id string, cmd []*RESPToken) error {
	val := tc.dict[id]
	if val == nil {
		return io.EOF
	}
	tc.dict[id] = append(val, cmd)
	return nil
}

func (tc *TransactionContext) GetQueuedCommands(id string) [][]*RESPToken {
	return tc.dict[id]
}

func (mc *TransactionContext) RemoveTxConnection(id string) {
	delete(mc.dict, id)
}
