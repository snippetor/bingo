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

func GetDefaultMessagePacker() IMessagePacker {
	return globalPacker
}

type LCNetProtocol byte
type SCNetProtocol byte

const (
	Tcp       LCNetProtocol = iota
	WebSocket
)

const (
	Http SCNetProtocol = iota
)

// 创建长连接服务器
func NewLCServer(net LCNetProtocol) (ILCServer, bool) {
	var server ILCServer
	switch net {
	case Tcp:
		server = ILCServer(&tcpServer{})
	case WebSocket:
		server = ILCServer(&wsServer{})
	default:
		bingo.E("-- error net type '%d', must be 'ws' or 'tcp' --", net)
		return nil, false
	}
	return server, true
}

// 创建长连接客户端
func NewLCClient(net LCNetProtocol) (ILCClient, bool) {
	var client ILCClient
	switch net {
	case Tcp:
		client = ILCClient(&tcpClient{})
	case WebSocket:
		client = ILCClient(&wsClient{})
	default:
		bingo.E("-- error net type '%d', must be 'ws' or 'tcp' --", net)
		return nil, false
	}
	return client, true
}
