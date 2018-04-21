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
	"github.com/snippetor/bingo/codec"
)

var (
	globalPacker IMessagePacker
	JsonCodec codec.ICodec
	ProtobufCodec codec.ICodec
)

func init() {
	globalPacker = IMessagePacker(&DefaultMessagePacker{})
	JsonCodec = codec.NewCodec(codec.Json)
	ProtobufCodec = codec.NewCodec(codec.Protobuf)
}

func SetDefaultMessagePacker(packer IMessagePacker) {
	globalPacker = packer
}

func GetDefaultMessagePacker() IMessagePacker {
	return globalPacker
}

type NetProtocol int

const (
	Tcp       NetProtocol = iota
	WebSocket
)

// 同步执行网络监听
func Listen(net NetProtocol, port int, callback IMessageCallback) (IServer, bool) {
	var server IServer
	switch net {
	case Tcp:
		server = IServer(&tcpServer{})
	case WebSocket:
		server = IServer(&wsServer{})
	default:
		fwlogger.E("-- errors net type '%d', must be 'ws' or 'tcp' --", net)
		return nil, false
	}
	return server, server.listen(port, callback)
}

// 异步执行网络监听
func GoListen(net NetProtocol, port int, callback IMessageCallback) IServer {
	var server IServer
	switch net {
	case Tcp:
		server = IServer(&tcpServer{})
	case WebSocket:
		server = IServer(&wsServer{})
	default:
		fwlogger.E("-- errors net type '%d', must be 'ws' or 'tcp' --", net)
		return nil
	}
	go server.listen(port, callback)
	return server
}

// 同步连接服务器
func Connect(net NetProtocol, serverAddr string, callback IMessageCallback) (IClient, bool) {
	var client IClient
	switch net {
	case Tcp:
		client = IClient(&tcpClient{})
	case WebSocket:
		client = IClient(&wsClient{})
	default:
		fwlogger.E("-- errors net type '%d', must be 'ws' or 'tcp' --", net)
		return nil, false
	}
	return client, client.connect(serverAddr, callback)
}

// 异步连接服务器
func GoConnect(net NetProtocol, serverAddr string, callback IMessageCallback) IClient {
	var client IClient
	switch net {
	case Tcp:
		client = IClient(&tcpClient{})
	case WebSocket:
		client = IClient(&wsClient{})
	default:
		fwlogger.E("-- errors net type '%d', must be 'ws' or 'tcp' --", net)
		return nil
	}
	go client.connect(serverAddr, callback)
	return client
}
