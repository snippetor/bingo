package net

import (
	"errors"
	"github.com/snippetor/bingo/utils"
)

type MessageId int32
type MessageBody []byte

const (
	MSGID_CONNECT_DISCONNECT = -1
	MSGID_CONNECT_CONNECTED = -2
)

// 消息回调
type IMessageCallback func(conn IConn, msgId MessageId, body MessageBody)

// 消息封装器接口
type IMessagePacker interface {
	// 封包，传入消息ID和包体，返回字节集
	Pack(MessageId, MessageBody) []byte
	// 解包，传入符合包结构的字节集，返回消息ID，包体，剩余内容
	Unpack([]byte) (MessageId, MessageBody, []byte)
}

// 长连接接口
type IConn interface {
	Send(msgId MessageId, body MessageBody) bool
	Close()
	Address() string
	read(*[]byte) (int, error)
	GetNetProtocol() NetProtocol
	Identity() utils.Identity
	GetState() ConnState
	setState(ConnState)
}

// 服务器接口
type IServer interface {
	listen(int, IMessageCallback) bool
	GetConnection(utils.Identity) (IConn, bool)
	Close()
}

// 客户端接口
type IClient interface {
	connect(string, IMessageCallback) bool
	Send(msgId MessageId, body MessageBody) bool
	Close()
}

type ConnState int

const (
	STATE_CLOSED     ConnState = iota
	STATE_CONNECTING
	STATE_CONNECTED
)

var (
	identifier *utils.Identifier
)

func init() {
	identifier = utils.NewIdentifier(1)
}

type absConn struct {
	identity utils.Identity
	state    ConnState
}

func (c *absConn) Send(msgId MessageId, body MessageBody) bool {
	return false
}
func (c *absConn) Close() {
}
func (c *absConn) Address() string {
	return "0.0.0.0"
}
func (c *absConn) read(*[]byte) (int, error) {
	return -1, errors.New("-- not implements --")
}
func (c *absConn) GetNetProtocol() NetProtocol {
	return -1
}
func (c *absConn) Identity() utils.Identity {
	if !identifier.IsValidIdentity(c.identity) {
		c.identity = identifier.GenIdentity()
	}
	return c.identity
}
func (c *absConn) GetState() ConnState {
	return c.state
}

func (c *absConn) setState(state ConnState) {
	c.state = state
}
