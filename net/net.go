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
	"github.com/snippetor/bingo/log/fwlogger"
)

const (
	MsgTypeReq int32 = 1
	MsgTypeAck int32 = 2
	MsgTypeNtf int32 = 3
	MsgTypeCmd int32 = 4
)

const (
	Tcp       Protocol = iota
	WebSocket
	Kcp
)

const (
	MsgIdConnDisconnect = -1
	MsgIdConnConnect    = -2
)

// 消息回调
type MessageCallback func(conn Conn, msgId MessageId, body MessageBody)

// 服务器接口
type Server interface {
	listen(int, MessageCallback) bool
	GetConnection(uint32) (Conn, bool)
	Close()
}

// 客户端接口
type Client interface {
	connect(string, MessageCallback) bool
	Send(msgId MessageId, body MessageBody) bool
	Close()
	Reconnect()
}

// 网络协议定义
type Protocol int

// 消息ID
/**
	消息ID组成：
	123456
	1 表示消息类别REQ、ACK、NTF、CMD
	2 表示消息组，例如gate, hall, game
	34 在game消息组中表示游戏ID，其他消息组暂时为0
	56 表示具体消息
 */
type MessageId int32

// 消息体
type MessageBody []byte

// 默认消息封装器
var (
	globalPacker MessagePacker
)

func init() {
	globalPacker = &messagePacker{}
}

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

func (i MessageBody) Len() int {
	return len(i)
}

// 封装消息ID
func PackId(idType, group, extra, id int32) MessageId {
	return MessageId(idType*100000 + group*10000 + extra*100 + id)
}

// 同步执行网络监听
func Listen(net Protocol, port int, callback MessageCallback) (Server, bool) {
	var server Server
	switch net {
	case Tcp:
		server = &tcpServer{}
	case WebSocket:
		server = &wsServer{}
	default:
		fwlogger.E("-- errors net type '%d', must be 'ws' or 'tcp' --", net)
		return nil, false
	}
	return server, server.listen(port, callback)
}

// 异步执行网络监听
func GoListen(net Protocol, port int, callback MessageCallback) Server {
	var server Server
	switch net {
	case Tcp:
		server = &tcpServer{}
	case WebSocket:
		server = &wsServer{}
	default:
		fwlogger.E("-- errors net type '%d', must be 'ws' or 'tcp' --", net)
		return nil
	}
	go server.listen(port, callback)
	return server
}

// 同步连接服务器
func Connect(net Protocol, serverAddr string, callback MessageCallback) (Client, bool) {
	var client Client
	switch net {
	case Tcp:
		client = &tcpClient{}
	case WebSocket:
		client = &wsClient{}
	default:
		fwlogger.E("-- errors net type '%d', must be 'ws' or 'tcp' --", net)
		return nil, false
	}
	return client, client.connect(serverAddr, callback)
}

// 异步连接服务器
func GoConnect(net Protocol, serverAddr string, callback MessageCallback) Client {
	var client Client
	switch net {
	case Tcp:
		client = &tcpClient{}
	case WebSocket:
		client = &wsClient{}
	default:
		fwlogger.E("-- errors net type '%d', must be 'ws' or 'tcp' --", net)
		return nil
	}
	go client.connect(serverAddr, callback)
	return client
}
