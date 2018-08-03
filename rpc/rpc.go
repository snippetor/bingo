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

package rpc

import (
	"github.com/snippetor/bingo/net"
	"github.com/snippetor/bingo/codec"
	"sync"
	"time"
	"github.com/snippetor/bingo/utils"
	"github.com/snippetor/bingo/log/fwlogger"
)

var (
	defaultCodec codec.ICodec
)

func init() {
	defaultCodec = codec.NewCodec(codec.Protobuf)
}

type IEndStub interface {
	SetRouter(router *Router)
	Call(method string, args *Args) (*Args, bool)
	CallNoReturn(method string, args *Args) bool
}

type Server struct {
	name       string
	modelName  string
	serv       net.IServer
	clients    []*Client
	l          *sync.RWMutex
	identifier *utils.Identifier
	callSyncWorker
	r          *Router
}

func (s *Server) Listen(name, modelName string, port int) {
	s.name = name
	s.modelName = modelName
	s.l = &sync.RWMutex{}
	s.identifier = utils.NewIdentifier(2)
	if s.serv = net.GoListen(net.Tcp, port, s.handleMessage); s.serv == nil {
		fwlogger.E("-- start rpc server failed! --")
	}
}

func (s *Server) Close() {
	s.serv.Close()
}

func (s *Server) handleMessage(conn net.IConn, msgId net.MessageId, body net.MessageBody) {
	switch RPC_MSGID(msgId) {
	case net.MSGID_CONNECT_DISCONNECT:
		for i, c := range s.clients {
			if c.conn.Identity() == conn.Identity() {
				s.l.Lock()
				s.clients = append(s.clients[:i], s.clients[i+1:]...)
				s.l.Unlock()
				break
			}
		}
	case RPC_MSGID_HANDSHAKE:
		// ack handshake to RPC server
		body := defaultCodec.Marshal(&RPCHandShake{EndName: s.name, EndModelName: s.modelName})
		if !conn.Send(net.MessageId(RPC_MSGID_HANDSHAKE), body) {
			fwlogger.E("-- send handshake %s failed! send message failed --", s.name)
		}
		handshake := &RPCHandShake{}
		defaultCodec.Unmarshal(body, handshake)
		c := &Client{}
		c.conn = conn
		c.EndName = handshake.EndName
		c.EndModelName = handshake.EndModelName
		c.state = net.STATE_CONNECTED
		c.addr = conn.Address()
		c.identifier = utils.NewIdentifier(3)
		c.l = &sync.RWMutex{}
		s.l.Lock()
		s.clients = append(s.clients, c)
		s.l.Unlock()
	case RPC_MSGID_CALL:
		call := &RPCMethodCall{}
		defaultCodec.Unmarshal(body, call)
		fwlogger.D("@call noreturn method %s(%d) with args %s", call.Method, call.CallSeq, call.Args)
		args := Args{}
		args.FromRPCMap(&call.Args)
		ctx := &Context{conn: conn, callSeq: call.CallSeq, Method: call.Method, Args: args}
		s.r.Invoke(s.name, call.Method, ctx)
	case RPC_MSGID_RETURN:
		ret := &RPCMethodReturn{}
		defaultCodec.Unmarshal(body, ret)
		fwlogger.D("@receive return from RPC method %s(%d) with result %s", ret.Method, ret.CallSeq, ret.Returns)
		args := Args{}
		args.FromRPCMap(&ret.Returns)
		s.receiveResult(ret.CallSeq, &args)
	}
}

func (s *Server) GetClients() []*Client {
	return s.clients
}

func (s *Server) GetClient(name string) (*Client, bool) {
	for _, c := range s.clients {
		if name == c.EndName {
			return c, true
		}
	}
	return nil, false
}

func (s *Server) SetRouter(router *Router) {
	s.r = router
}

type Client struct {
	name         string
	modelName    string
	EndName      string // 远端节点名称
	EndModelName string // 远端节点模型名
	conn         net.IConn
	l            *sync.RWMutex
	addr         string
	state        net.ConnState
	identifier   *utils.Identifier
	callSyncWorker
	forceClose   bool
	r            *Router
	tcpClient    net.IClient
}

func (c *Client) Connect(name, modelName, serverAddress string) {
	c.name = name
	c.modelName = modelName
	c.state = net.STATE_CONNECTING
	c.addr = serverAddress
	c.identifier = utils.NewIdentifier(3)
	c.l = &sync.RWMutex{}
	if client := net.GoConnect(net.Tcp, serverAddress, c.handleMessage); client == nil {
		c.state = net.STATE_CLOSED
	} else {
		c.state = net.STATE_CONNECTED
		c.tcpClient = client
	}
}

func (c *Client) Close() {
	c.forceClose = true
	c.conn.Close()
}

func (c *Client) Call(method string, args *Args) (*Args, bool) {
	if c.conn == nil {
		fwlogger.E("-- call rpc method %s failed! rpc client not connect to server --", method)
		return nil, false
	}
	c.l.RLock()
	defer c.l.RUnlock()
	seq := c.identifier.GenIdentity()
	res := make(map[string]*RPCValue)
	args.ToRPCMap(&res)
	body := defaultCodec.Marshal(&RPCMethodCall{CallSeq: int32(seq), Method: method, Args: res})
	if !c.conn.Send(net.MessageId(RPC_MSGID_CALL), body) {
		fwlogger.E("-- call rpc method %s failed! send message failed --", method)
	} else {
		ch := make(chan *Args)
		c.waitingResult(&callTask{seq: seq, conn: c.conn, msg: body, c: func(result *Args) {
			ch <- result
		}, t: time.Now().UnixNano()})
		res := <-ch
		return res, true
	}
	return nil, false
}

func (c *Client) CallNoReturn(method string, args *Args) bool {
	if c.conn == nil {
		fwlogger.E("-- call rpc method %s failed! rpc client not connect to server --", method)
		return false
	}
	c.l.RLock()
	defer c.l.RUnlock()
	res := make(map[string]*RPCValue)
	args.ToRPCMap(&res)
	body := defaultCodec.Marshal(&RPCMethodCall{CallSeq: c.identifier.GenIdentity(), Method: method, Args: res})
	if !c.conn.Send(net.MessageId(RPC_MSGID_CALL), body) {
		fwlogger.E("-- call rpc method %s failed! send message failed --", method)
	} else {
		return true
	}
	return false
}

func (c *Client) SetRouter(router *Router) {
	c.r = router
}

func (c *Client) handleMessage(conn net.IConn, msgId net.MessageId, body net.MessageBody) {
	switch RPC_MSGID(msgId) {
	case RPC_MSGID(net.MSGID_CONNECT_CONNECTED):
		fwlogger.D("-- %s connect RPC server success  --", c.EndName)
		// send handshake to RPC server
		body := defaultCodec.Marshal(&RPCHandShake{EndName: c.name, EndModelName: c.modelName})
		if !conn.Send(net.MessageId(RPC_MSGID_HANDSHAKE), body) {
			fwlogger.E("-- send handshake %s failed! send message failed --", c.EndName)
		}
	case RPC_MSGID(net.MSGID_CONNECT_DISCONNECT):
		c.conn = nil
		c.state = net.STATE_CLOSED
		if !c.forceClose {
			c.reconnect()
		}
	case RPC_MSGID_HANDSHAKE:
		handshake := &RPCHandShake{}
		defaultCodec.Unmarshal(body, handshake)
		c.EndName = handshake.EndName
		c.EndModelName = handshake.EndModelName
	case RPC_MSGID_CALL:
		call := &RPCMethodCall{}
		defaultCodec.Unmarshal(body, call)
		fwlogger.D("@call method %s(%d) with args %s", call.Method, call.CallSeq, call.Args)
		args := Args{}
		args.FromRPCMap(&call.Args)
		ctx := &Context{conn: conn, callSeq: call.CallSeq, Method: call.Method, Args: args}
		c.r.Invoke(c.EndName, call.Method, ctx)
	case RPC_MSGID_RETURN:
		ret := &RPCMethodReturn{}
		defaultCodec.Unmarshal(body, ret)
		fwlogger.D("@receive return from RPC method %s(%d) with result %s", ret.Method, ret.CallSeq, ret.Returns)
		args := Args{}
		args.FromRPCMap(&ret.Returns)
		c.receiveResult(ret.CallSeq, &args)
	}
}

func (c *Client) reconnect() {
	if net.STATE_CLOSED == c.state {
		time.Sleep(5 * time.Second)
		fwlogger.D("@reconnect RPC, remote address=%s", c.addr)
		c.tcpClient.Reconnect()
	}
}
