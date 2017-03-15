package net

import (
	"sync"
)

type MessageId int32
type MessageBody []byte

const (
	MSGID_CONNECT_DISCONNECT = -1
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

// 连接接口
type IConn interface {
	Send(msgId MessageId, body MessageBody) bool
	Close()
	Address() string
	read(*[]byte) (int, error)
	GetNetProtocol() NetProtocol
	Identity() Identity
}

// 服务器接口
type IServer interface {
	listen(int, IMessageCallback) bool
	GetConnection(Identity) (IConn, bool)
	Close()
}

// 客户端接口
type IClient interface {
	connect(string, IMessageCallback) bool
	Send(msgId MessageId, body MessageBody) bool
	Close()
}

// ID生成
type Identity int

const (
	MINMUM_IDENTIFY = 1000000
	MAXMUM_IDENTIFY = 9999999
)

var (
	_identify_ Identity = MINMUM_IDENTIFY
	l          *sync.Mutex
)

func init() {
	l = &sync.Mutex{}
}

func genIdentity() Identity {
	l.Lock()
	_identify_++
	if _identify_ > MAXMUM_IDENTIFY {
		_identify_ = MINMUM_IDENTIFY
	}
	l.Unlock()
	return _identify_
}

func isValidIdentity(id Identity) bool {
	return id <= MAXMUM_IDENTIFY && id >= MINMUM_IDENTIFY
}
