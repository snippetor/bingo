package rpc

import (
	"github.com/snippetor/bingo/net"
	"github.com/snippetor/bingo/protocol"
	"github.com/snippetor/bingo"
	"sync"
	"time"
	"math/rand"
)

type RPCEnd struct {
	connId net.Identity
	name   string
}

type Server struct {
	serv net.ILCServer
	end  map[string][]*RPCEnd
	l    *sync.RWMutex
}

func (s *Server) Listen(port int) {
	s.l = &sync.RWMutex{}
	s.end = make(map[string][]*RPCEnd)
	if s.serv = net.GoListen(net.Tcp, port, s.handleMessage); s.serv == nil {
		bingo.E("-- start rpc server failed! --")
	}
}

func (s *Server) Call(endName, target, method string, args Args) {
	if s.serv == nil {
		bingo.E("-- call rpc method %s.%s failed! rpc server no startup --", target, method)
		return
	}
	s.l.RLock()
	defer s.l.RUnlock()
	if ends, ok := s.end[endName]; ok {
		size := len(ends)
		if size == 0 {
			bingo.E("-- call rpc method failed! no end connected with name is %s --", endName)
			return
		}
		var i int = 0
		if size > 1 {
			i = rand.Intn(len(ends))
		}
		if conn, ok := s.serv.GetConnection(ends[i].connId); ok {
			if body, ok := protocol.Marshal(&RPCMethodCall{Target: target, Method:method, Args:args}); ok {
				if !conn.Send(net.MessageId(RPC_MSGID_CALL), body) {
					bingo.E("-- call rpc method %s.%s failed! send message failed --", target, method)
				}
			} else {
				bingo.E("-- call rpc method %s.%s failed! marshal message failed --", target, method)
			}
		} else {
			bingo.E("-- call rpc method %s.%s failed! no connection for call --", target, method)
		}
	} else {
		bingo.E("-- call rpc method failed! no end connected with name is %s --", endName)
	}
}

func (s *Server) handleMessage(conn net.ILongConn, msgId net.MessageId, body net.MessageBody) {
	switch msgId {
	case net.MSGID_CONNECT_DISCONNECT:
		for name, arr := range s.end {
			for i, e := range arr {
				if e.connId == conn.Identity() {
					var newEnds []*RPCEnd = make([]*RPCEnd, 0)
					newEnds = append(newEnds, arr[:i]...)
					newEnds = append(newEnds, arr[i + 1:]...)
					s.l.Lock()
					s.end[name] = newEnds
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
		end := &RPCEnd{connId:conn.Identity(), name:handshake.EndName}
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
		bingo.D("@call method %s.%s with args %s", call.Target, call.Method, call.Args)
		ctx := &Context{conn: conn, target: call.Target, method: call.Method, args: call.Args}
		callMethod(call.Target, call.Method, ctx)
	}
}

type Client struct {
	conn  net.ILongConn
	l     *sync.RWMutex
	addr  string
	state net.LongConnState
}

func (c *Client) Connect(serverAddress string) {
	c.state = net.STATE_CONNECTING
	c.addr = serverAddress
	c.l = &sync.RWMutex{}
	if net.GoConnect(net.Tcp, serverAddress, c.handleMessage) == nil {
		c.state = net.STATE_CLOSED
		c.reconnect()
	} else {
		c.state = net.STATE_CONNECTED
	}
}

func (c *Client) Call(target, method string, args Args) {
	if c.conn == nil {
		bingo.E("-- call rpc method %s.%s failed! rpc client not connect to server --", target, method)
		return
	}
	c.l.RLock()
	defer c.l.RUnlock()
	if body, ok := protocol.Marshal(&RPCMethodCall{Target: target, Method:method, Args:args}); ok {
		if !c.conn.Send(net.MessageId(RPC_MSGID_CALL), body) {
			bingo.E("-- call rpc method %s.%s failed! send message failed --", target, method)
		}
	} else {
		bingo.E("-- call rpc method %s.%s failed! marshal message failed --", target, method)
	}
}

func (c *Client) handleMessage(conn net.ILongConn, msgId net.MessageId, body net.MessageBody) {
	switch msgId {
	case net.MSGID_CONNECT_DISCONNECT:
		c.conn = nil
		c.state = net.STATE_CLOSED
		c.reconnect()
	case RPC_MSGID_CALL:
		call := &RPCMethodCall{}
		if err := protocol.Unmarshal(body, call); err != nil {
			bingo.E("-- RPC call failed! -- ")
			return
		}
		bingo.D("@call method %s.%s with args %s", call.Target, call.Method, call.Args)
		ctx := &Context{conn: conn, target:call.Target, method: call.Method, args: call.Args}
		callMethod(call.Target, call.Method, ctx)
	}
}

func (c *Client) reconnect() {
	if net.STATE_CLOSED == c.state {
		if net.GoConnect(net.Tcp, c.addr, c.handleMessage) == nil {
			time.Sleep(1 * time.Second)
			c.state = net.STATE_CLOSED
			c.reconnect()
		} else {
			c.state = net.STATE_CONNECTED
		}
	}
}
