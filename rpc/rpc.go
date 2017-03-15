package rpc

import (
	"github.com/snippetor/bingo/net"
	"github.com/snippetor/bingo/protocol"
	"github.com/snippetor/bingo"
	"sync"
)

type RPCEnd struct {
	connId  net.Identity
	ability []string
	name    string
}

type Server struct {
	end map[string][]*RPCEnd
	l   *sync.RWMutex
}

func (s *Server) Listen(port int) {
	s.l = &sync.RWMutex{}
	s.end = make(map[string][]*RPCEnd)
	net.GoListen(net.Tcp, port, s.handleMessage)
}

func (s *Server) handleMessage(conn net.IConn, msgId net.MessageId, body net.MessageBody) {
	switch msgId {
	case net.MSGID_CONNECT_DISCONNECT:
		for name, arr := range s.end {
			for i, e := range arr {
				if e.connId == conn.Identity() {
					var new []*RPCEnd = make([]*RPCEnd, 0)
					new = append(new, arr[:i]...)
					new = append(new, arr[i + 1:]...)
					s.l.Lock()
					s.end[name] = new
					s.l.Unlock()
					break
				}
			}
		}
	case RPC_MSGID_HANDSHAKE:
		handshake := &RPCHandShake{}
		if err := protocol.Unmarshal(body, handshake); err != nil {
			bingo.E("-- RPC handshake failed! -- ")
			return
		}
		end := &RPCEnd{connId:conn.Identity(), ability:handshake.EndAbility, name:handshake.EndName}
		ends, ok := s.end[handshake.EndName]
		if !ok {
			ends = make([]*RPCEnd, 0)
		}
		ends = append(ends, end)
		s.l.Lock()
		s.end[handshake.EndName] = ends
		s.l.Unlock()
	case RPC_MSGID_CALL:
		call := &RPCMethodCall{}
		if err := protocol.Unmarshal(body, call); err != nil {
			bingo.E("-- RPC call failed! -- ")
			return
		}
		bingo.D("@call method %s (%s) with args %s", call.MethodName, call.Version, call.Args)
		ctx := &Context{conn: conn, method: call.MethodName, args: call.Args, version: call.Version}

	}
}

type Client struct {
	conn net.IConn
	l    *sync.RWMutex
	addr string
}

func (c *Client) Connect(serverAddress string) {
	c.addr = serverAddress
	c.l = &sync.RWMutex{}
	if net.GoConnect(net.Tcp, serverAddress, c.handleMessage) == nil {
		c.reconnect()
	}
}

func (c *Client) handleMessage(conn net.IConn, msgId net.MessageId, body net.MessageBody) {
	switch msgId {
	case net.MSGID_CONNECT_DISCONNECT:
		c.conn = nil
		c.reconnect()
	case RPC_MSGID_CALL:
		call := &RPCMethodCall{}
		if err := protocol.Unmarshal(body, call); err != nil {
			bingo.E("-- RPC call failed! -- ")
			return
		}
		bingo.D("@call method %s (%s) with args %s", call.MethodName, call.Version, call.Args)
		ctx := &Context{conn: conn, method: call.MethodName, args: call.Args, version: call.Version}

	}
}

func (c *Client) reconnect() {

}
