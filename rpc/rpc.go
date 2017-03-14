package rpc

import (
	"github.com/snippetor/bingo/net"
	"github.com/snippetor/bingo/protocol"
	"github.com/snippetor/bingo"
	"encoding/gob"
	"encoding/json"
)

type server struct {
	port int
}

func RunAsServer(port int) {
	go net.Listen(net.Tcp, port, handleMessage)
}

func RunAsClient(serverAddr string) {
	go net.Connect(net.Tcp, serverAddr, handleMessage)
}

func handleMessage(conn net.IConn, msgId net.MessageId, body net.MessageBody) {
	switch msgId {
	case RPC_MSGID_CALL:
		call := &MethodCall{}
		if err := protocol.Unmarshal(body, call); err != nil {
			bingo.E("-- RPC call failed! -- ")
			return
		}
		bingo.D("@call method %s (%s) with args %s", call.MethodName, call.Version, call.Args)
		ctx := &Context{conn: conn, method: call.MethodName, args: call.Args, version: call.Version}

	}
}
