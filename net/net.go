package net

import (
	"github.com/snippetor/bingo"
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

type NetProtocol byte

const (
	Tcp       NetProtocol = iota
	WebSocket
)

// 同步执行网络监听
func Listen(net NetProtocol, port int, callback IMessageCallback) bool {
	var server iServer
	switch net {
	case Tcp:
		server = iServer(&tcpServer{})
	case WebSocket:
		server = iServer(&wsServer{})
	default:
		bingo.E("-- error net type '%d', must be 'ws' or 'tcp' --", net)
		return false
	}
	return server.listen(port, callback)
}

// 同步连接服务器
func Connect(net NetProtocol, serverAddr string, callback IMessageCallback) bool {
	var client iClient
	switch net {
	case Tcp:
		client = iClient(&tcpClient{})
	case WebSocket:
		client = iClient(&wsClient{})
	default:
		bingo.E("-- error net type '%d', must be 'ws' or 'tcp' --", net)
		return false
	}
	return client.connect(serverAddr, callback)
}
