package net

import (
	"sync"
	"errors"
)

// ------------------------------------ 长连接 -------------------------------------//

type MessageId int32
type MessageBody []byte

const (
	MSGID_CONNECT_DISCONNECT = -1
)

// 消息回调
type IMessageCallback func(conn ILongConn, msgId MessageId, body MessageBody)

// 消息封装器接口
type IMessagePacker interface {
	// 封包，传入消息ID和包体，返回字节集
	Pack(MessageId, MessageBody) []byte
	// 解包，传入符合包结构的字节集，返回消息ID，包体，剩余内容
	Unpack([]byte) (MessageId, MessageBody, []byte)
}

// 长连接接口
type ILongConn interface {
	Send(msgId MessageId, body MessageBody) bool
	Close()
	Address() string
	read(*[]byte) (int, error)
	GetNetProtocol() LCNetProtocol
	Identity() Identity
	GetState() LongConnState
	setState(LongConnState)
}

// 服务器接口
type ILCServer interface {
	Listen(int, IMessageCallback) bool
	GetConnection(Identity) (ILongConn, bool)
	Close()
}

// 客户端接口
type ILCClient interface {
	Connect(string, IMessageCallback) bool
	Send(msgId MessageId, body MessageBody) bool
	Close()
}

type LongConnState int

const (
	STATE_CLOSED     LongConnState = iota
	STATE_CONNECTING
	STATE_CONNECTED
)

type absLongConn struct {
	identity Identity
	state    LongConnState
}

func (c *absLongConn) Send(msgId MessageId, body MessageBody) bool {
	return false
}
func (c *absLongConn) Close() {
}
func (c *absLongConn) Address() string {
	return "0.0.0.0"
}
func (c *absLongConn) read(*[]byte) (int, error) {
	return -1, errors.New("-- not implements --")
}
func (c *absLongConn) GetNetProtocol() LCNetProtocol {
	return -1
}
func (c *absLongConn) Identity() Identity {
	if !isValidIdentity(c.identity) {
		c.identity = genIdentity()
	}
	return c.identity
}
func (c *absLongConn) GetState() LongConnState {
	return c.state
}

func (c *absLongConn) setState(state LongConnState) {
	c.state = state
}

// ------------------------------------ 短连接 -------------------------------------//
