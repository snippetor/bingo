package net

import (
	"github.com/snippetor/bingo/log/fwlogger"
)

var (
	globalPacker IMessagePacker
)

func init() {
	globalPacker = IMessagePacker(&DefaultMessagePacker{})
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
		fwlogger.E("-- error net type '%d', must be 'ws' or 'tcp' --", net)
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
		fwlogger.E("-- error net type '%d', must be 'ws' or 'tcp' --", net)
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
		fwlogger.E("-- error net type '%d', must be 'ws' or 'tcp' --", net)
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
		fwlogger.E("-- error net type '%d', must be 'ws' or 'tcp' --", net)
		return nil
	}
	go client.connect(serverAddr, callback)
	return client
}
