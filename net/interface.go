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
	"encoding/json"
	"github.com/snippetor/bingo/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/snippetor/bingo/utils"
)

type MessageId int32
type MessageBody []byte

func (b MessageBody) ToJson(v interface{}) bool {
	return json.Unmarshal(b, v) == nil
}

func (b MessageBody) ToProtobuf(v interface{}) bool {
	return proto.Unmarshal(b, v.(proto.Message)) == nil
}

func (b MessageBody) FromJson(v interface{}) {
	res, err := json.Marshal(v)
	errors.Check(err)
	copy(b, res)
}

func (b MessageBody) FromProtobuf(v interface{}) {
	res, err := proto.Marshal(v.(proto.Message))
	errors.Check(err)
	copy(b, res)
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
	return -1, errors1.New("-- not implements --")
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
