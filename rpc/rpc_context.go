package rpc

import (
	"github.com/snippetor/bingo/net"
	"github.com/snippetor/bingo/protocol"
	"github.com/snippetor/bingo"
)

type Args map[string]string

type Context struct {
	conn   net.ILongConn
	target string
	method string
	args   Args
}

func (c *Context) Callback(target, method string, args Args) {
	if c.conn == nil {
		bingo.E("-- call rpc method %s.%s failed! Context#conn is nil --", target, method)
	} else {
		if body, ok := protocol.Marshal(&RPCMethodCall{Target: target, Method:method, Args:args}); ok {
			if !c.conn.Send(net.MessageId(RPC_MSGID_CALL), body) {
				bingo.E("-- call rpc method %s.%s failed! send message failed --", target, method)
			}
		} else {
			bingo.E("-- call rpc method %s.%s failed! marshal message failed --", target, method)
		}
	}
}
