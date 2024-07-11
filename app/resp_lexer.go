package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type RedisDataType string

const (
	String         = "+"
	Error          = "-"
	Integer        = ":"
	BulkString     = "$"
	Array          = "*"
	Null           = "_"
	Boolean        = "#"
	Double         = ","
	BigNumber      = "("
	BulkError      = "!"
	VerbatimString = "="
	Map            = "%"
	Set            = "~"
	Push           = ">"
)

type RESPToken struct {
	Type   RedisDataType
	Value  interface{}
	length int
}

type RESPLexer struct {
	reader    *bufio.Reader
	readIndex int
}

func NewRESPLexer(reader *bufio.Reader) *RESPLexer {
	return &RESPLexer{
		reader:    reader,
		readIndex: 0,
	}
}

func (rl *RESPLexer) nextToken() ([]*RESPToken, error) {
	tokenType, _ := rl.reader.Peek(1)

	switch string(tokenType[0]) {
	case Array:
		return rl.parseArray()
	case BulkString:
		return rl.parseBulkString()
	// case String:
	// 	token = p.parseString()
	// case Error:
	// 	token = p.parseError()
	case Integer:
		return rl.parseInteger()
		// case Null:
		// 	token = p.parseNull()
		// case Boolean:
		// 	token = p.parseBoolean()
		// case Double:
		// 	token = p.parseDouble()
		// case BigNumber:
		// 	token = p.parseBigNumber()
		// case BulkError:
		// 	token = p.parseBulkError()
		// case VerbatimString:
		// 	token = p.VerbatimString()
		// case Map:
		// 	token = p.parseMap()
		// case Set:
		// 	token = p.parseSet()
		// case Push:
	}
	return nil, fmt.Errorf("No token match")
}

func (rl *RESPLexer) ProduceTokens(ch chan<- Message) error {
	for {
		firstByte, err := rl.reader.Peek(1)
		if err != nil {
			return rl.handleReadError(err)
		}
		if string(firstByte) != Array {
			return fmt.Errorf("first byte must be array")
		}

		tokens, err := rl.parseArray()
		if err != nil {
			return rl.handleReadError(err)
		}

		ch <- Message{command: tokens}
	}
}

func (p *RESPLexer) handleReadError(err error) error {
	if err == io.EOF {
		return err
	}
	fmt.Errorf("Error encountered %s", err)
	panic(err)
}

func (p *RESPLexer) readTillCLRF() ([]byte, error) {
	var b []byte
	for {
		chars, err := p.reader.ReadBytes('\r')
		if err != nil {
			return nil, p.handleReadError(err)
		}
		b = append(b, chars...)
		nextChar, err := p.reader.Peek(1)
		if err != nil {
			return nil, p.handleReadError(err)
		}
		if nextChar[0] == '\n' {
			newLine, err := p.reader.ReadByte()
			if err != nil {
				return nil, p.handleReadError(err)
			}
			b = append(b, newLine)
			break
		}
	}
	return b, nil
}

func (rl *RESPLexer) parseInteger() ([]*RESPToken, error) {
	integer, err := rl.readTillCLRF()
	if err != nil {
		return nil, rl.handleReadError(err)
	}
	numberStr := string(integer[1:])
	number, _ := strconv.Atoi(numberStr)
	token := make([]*RESPToken, 0, 1)
	token = append(token, &RESPToken{Type: Integer,
		Value: number})
	return token, nil
}

func (rl *RESPLexer) parseBulkString() ([]*RESPToken, error) {
	bulkStr, err := rl.readTillCLRF()
	if err != nil {
		return nil, rl.handleReadError(err)
	}
	lengthStr := string(bulkStr[1 : len(bulkStr)-2])
	length, _ := strconv.Atoi(lengthStr)

	content, err := rl.readTillCLRF()
	if err != nil {
		return nil, rl.handleReadError(err)
	}
	removeCLRF := content[:len(content)-2]
	return []*RESPToken{
		{
			Type:   BulkString,
			Value:  string(removeCLRF),
			length: length,
		},
	}, nil
}

func (rl *RESPLexer) parseArray() ([]*RESPToken, error) {
	tokens := make([]*RESPToken, 0)

	arrayLength, err := rl.readTillCLRF()
	if err != nil {
		return nil, rl.handleReadError(err)
	}

	lengthStr := string(arrayLength[1 : len(arrayLength)-2])
	length, _ := strconv.Atoi(lengthStr)

	tokens = append(tokens, &RESPToken{
		Type:   Array,
		length: length,
	})

	var parseError error
	for i := 0; i < length; i++ {
		nextToken, err := rl.nextToken()
		if err != nil {
			parseError = err
			break
		}
		tokens = append(tokens, nextToken...)
	}

	if parseError != nil {
		return nil, parseError
	}

	return tokens, nil
}
