package rpc

import (
	"github.com/snippetor/bingo/net"
	"github.com/snippetor/bingo/protocol"
	"github.com/snippetor/bingo"
	"sync"
	"time"
	"reflect"
	"strings"
	"math/rand"
)

type Args map[string]string

type Context struct {
	conn   net.IConn
	method string
	args   Args
}

var (
	methods map[string]reflect.Value
)

// v：必须是指针
func RegisterMethods(target string, v interface{}) {
	if methods == nil {
		methods = make(map[string]reflect.Value)
	}
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		bingo.E("-- register methods failed! v must be a pointer. --")
		return
	}
	t := reflect.TypeOf(v)
	for i := 0; i < t.NumMethod(); i++ {
		methods[makeKey(target, t.Method(i).Name)] = reflect.ValueOf(v).Method(i)
	}
}

func callMethod(target, method string, ctx *Context) {
	if v, ok := methods[makeKey(target, method)]; ok {
		v.Call([]reflect.Value{reflect.ValueOf(ctx)})
	}
}

func makeKey(target, method string) string {
	return target + "." + method
}

type RPCEnd struct {
	connId net.Identity
	name   string
}

type Server struct {
	serv net.IServer
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

func (s *Server) Call(endName, methodName string, args Args) {
	if s.serv == nil {
		bingo.E("-- call rpc method %s failed! rpc server no startup --", methodName)
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
			if body, ok := protocol.Marshal(&RPCMethodCall{MethodName:methodName, Args:args}); ok {
				if !conn.Send(net.MessageId(RPC_MSGID_CALL), body) {
					bingo.E("-- call rpc method %s failed! send message failed --", methodName)
				}
			} else {
				bingo.E("-- call rpc method %s failed! marshal message failed --", methodName)
			}
		} else {
			bingo.E("-- call rpc method %s failed! no connection for call --", methodName)
		}
	} else {
		bingo.E("-- call rpc method failed! no end connected with name is %s --", endName)
	}
}

func (s *Server) handleMessage(conn net.IConn, msgId net.MessageId, body net.MessageBody) {
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
		bingo.D("@call method %s with args %s", call.MethodName, call.Args)
		ctx := &Context{conn: conn, method: call.MethodName, args: call.Args}
		strs := strings.Split(call.MethodName, ".")
		if len(strs) >= 2 {
			callMethod(strs[0], strs[1], ctx)
		} else {
			bingo.E("-- RPC call failed! MethodName is not in format {Target}.{Method} -- ")
		}
	}
}

const (
	STATE_CLOSED     = 0
	STATE_CONNECTING = 1
	STATE_CONNECTED  = 2
)

type Client struct {
	conn  net.IConn
	l     *sync.RWMutex
	addr  string
	state int
}

func (c *Client) Connect(serverAddress string) {
	c.state = STATE_CONNECTING
	c.addr = serverAddress
	c.l = &sync.RWMutex{}
	if net.GoConnect(net.Tcp, serverAddress, c.handleMessage) == nil {
		c.state = STATE_CLOSED
		c.reconnect()
	} else {
		c.state = STATE_CONNECTED
	}
}

func (c *Client) Call(methodName string, args Args) {
	if c.conn == nil {
		bingo.E("-- call rpc method %s failed! rpc client not connect to server --", methodName)
		return
	}
	c.l.RLock()
	defer c.l.RUnlock()
	if body, ok := protocol.Marshal(&RPCMethodCall{MethodName:methodName, Args:args}); ok {
		if !c.conn.Send(net.MessageId(RPC_MSGID_CALL), body) {
			bingo.E("-- call rpc method %s failed! send message failed --", methodName)
		}
	} else {
		bingo.E("-- call rpc method %s failed! marshal message failed --", methodName)
	}
}

func (c *Client) handleMessage(conn net.IConn, msgId net.MessageId, body net.MessageBody) {
	switch msgId {
	case net.MSGID_CONNECT_DISCONNECT:
		c.conn = nil
		c.state = STATE_CLOSED
		c.reconnect()
	case RPC_MSGID_CALL:
		call := &RPCMethodCall{}
		if err := protocol.Unmarshal(body, call); err != nil {
			bingo.E("-- RPC call failed! -- ")
			return
		}
		bingo.D("@call method %s with args %s", call.MethodName, call.Args)
		ctx := &Context{conn: conn, method: call.MethodName, args: call.Args}
		strs := strings.Split(call.MethodName, ".")
		if len(strs) >= 2 {
			callMethod(strs[0], strs[1], ctx)
		} else {
			bingo.E("-- RPC call failed! MethodName is not in format {Target}.{Method} -- ")
		}
	}
}

func (c *Client) reconnect() {
	if STATE_CLOSED == c.state {
		if net.GoConnect(net.Tcp, c.addr, c.handleMessage) == nil {
			time.Sleep(1 * time.Second)
			c.state = STATE_CLOSED
			c.reconnect()
		} else {
			c.state = STATE_CONNECTED
		}
	}
}
