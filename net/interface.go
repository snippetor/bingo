// Copyright 2017 bingo Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package net

import (
	errors1 "errors"
	"github.com/snippetor/bingo/utils"
	"github.com/snippetor/bingo/codec"
)

/**
消息ID组成：
123456
1 表示消息类别REQ、ACK、NTF、CMD
2 表示消息组，例如gate, hall, game
34 在game消息组中表示游戏ID，其他消息组暂时为0
56 表示具体消息
 */
const (
	MsgTypeReq int32 = 1
	MsgTypeAck int32 = 2
	MsgTypeNtf int32 = 3
	MsgTypeCmd int32 = 4
)

type MessageId int32
type MessageBody []byte

func (i MessageId) Int32() int32 {
	return int32(i)
}

func (i MessageId) Int() int {
	return int(i)
}

func (i MessageId) Type() int32 {
	return i.Int32() / 100000
}

func (i MessageId) Group() int32 {
	return (i.Int32() % 100000) / 10000
}

func (i MessageId) Extra() int32 {
	return (i.Int32() % 10000) / 100
}

func (i MessageId) MsgId() int32 {
	return i.Int32() % 100
}

func PackId(idType, group, extra, id int32) MessageId {
	return MessageId(idType*100000 + group*10000 + extra*100 + id)
}

const (
	MSGID_CONNECT_DISCONNECT = -1
	MSGID_CONNECT_CONNECTED  = -2
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
	Identity() uint32
	GetState() ConnState
	setState(ConnState)
}

// 服务器接口
type IServer interface {
	listen(int, IMessageCallback) bool
	GetConnection(uint32) (IConn, bool)
	Close()
}

// 客户端接口
type IClient interface {
	connect(string, IMessageCallback) bool
	Send(msgId MessageId, body MessageBody) bool
	Close()
	Reconnect()
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
	identity uint32
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
	return -1, errors1.New("-- not implements --")
}
func (c *absConn) GetNetProtocol() NetProtocol {
	return -1
}
func (c *absConn) Identity() uint32 {
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
