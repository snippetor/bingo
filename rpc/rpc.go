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
	Call(method string, args *Args, callback RPCCallback) bool
	CallNoReturn(method string, args *Args) bool
}

type Server struct {
	endName    string
	serv       net.IServer
	clients    map[string]*Client
	l          *sync.RWMutex
	identifier *utils.Identifier
	callSyncWorker
}

func (s *Server) Listen(endName string, port int) {
	s.endName = endName
	s.l = &sync.RWMutex{}
	s.clients = make(map[string]*Client)
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
		for name, c := range s.clients {
			if c.conn.Identity() == conn.Identity() {
				s.l.Lock()
				delete(s.clients, name)
				s.l.Unlock()
				break
			}
		}
	case RPC_MSGID_HANDSHAKE:
		handshake := &RPCHandShake{}
		if err := defaultCodec.Unmarshal(body, handshake); err != nil {
			fwlogger.E("-- RPC handshake failed! -- ")
			return
		}
		c := &Client{}
		c.conn = conn
		c.endName = handshake.EndName
		c.state = net.STATE_CONNECTED
		c.addr = conn.Address()
		c.identifier = utils.NewIdentifier(3)
		c.l = &sync.RWMutex{}
		s.l.Lock()
		s.clients[handshake.EndName] = c
		s.l.Unlock()
	case RPC_MSGID_CALL:
		call := &RPCMethodCall{}
		if err := defaultCodec.Unmarshal(body, call); err != nil {
			fwlogger.E("-- RPC noreturn call failed! -- ")
			return
		}
		fwlogger.D("@call noreturn method %s(%d) with args %s", call.Method, call.CallSeq, call.Args)
		ctx := &Context{conn: conn, callSeq: call.CallSeq, Method: call.Method, Args: call.Args}
		callMethod(s.endName, call.Method, ctx)
	case RPC_MSGID_RETURN:
		ret := &RPCMethodReturn{}
		if err := defaultCodec.Unmarshal(body, ret); err != nil {
			fwlogger.E("-- RPC return failed! -- ")
			return
		}
		fwlogger.D("@receive return from RPC method %s(%d) with result %s", ret.Method, ret.CallSeq, ret.Returns)
		s.receiveResult(utils.Identity(ret.CallSeq), &Result{Args: ret.Returns})
	}
}

func (s *Server) GetClients() *map[string]*Client {
	return &s.clients
}

func (s *Server) GetClient(name string) (*Client, bool) {
	c, ok := s.clients[name]
	return c, ok
}

type Client struct {
	endName    string
	conn       net.IConn
	l          *sync.RWMutex
	addr       string
	state      net.ConnState
	identifier *utils.Identifier
	callSyncWorker
	forceClose bool
}

func (c *Client) Connect(endName, serverAddress string) {
	c.endName = endName
	c.state = net.STATE_CONNECTING
	c.addr = serverAddress
	c.identifier = utils.NewIdentifier(3)
	c.l = &sync.RWMutex{}
	if net.GoConnect(net.Tcp, serverAddress, c.handleMessage) == nil {
		c.state = net.STATE_CLOSED
		c.reconnect()
	} else {
		c.state = net.STATE_CONNECTED
	}
}

func (c *Client) Close() {
	c.forceClose = true
	c.conn.Close()
}

func (c *Client) Call(method string, args *Args, callback RPCCallback) bool {
	if c.conn == nil {
		fwlogger.E("-- call rpc method %s failed! rpc client not connect to server --", method)
		return false
	}
	c.l.RLock()
	defer c.l.RUnlock()
	seq := c.identifier.GenIdentity()
	if body, err := defaultCodec.Marshal(&RPCMethodCall{CallSeq: int32(seq), Method: method, Args: *args}); err == nil {
		if !c.conn.Send(net.MessageId(RPC_MSGID_CALL), body) {
			fwlogger.E("-- call rpc method %s failed! send message failed --", method)
		} else {
			c.waitingResult(&callTask{seq: seq, conn: c.conn, msg: body, c: callback, t: time.Now().UnixNano()})
			return true
		}
	} else {
		fwlogger.E("-- call rpc method %s failed! marshal message failed --", method)
	}
	return false
}

func (c *Client) CallNoReturn(method string, args *Args) bool {
	if c.conn == nil {
		fwlogger.E("-- call rpc method %s failed! rpc client not connect to server --", method)
		return false
	}
	c.l.RLock()
	defer c.l.RUnlock()
	if body, err := defaultCodec.Marshal(&RPCMethodCall{CallSeq: int32(c.identifier.GenIdentity()), Method: method, Args: *args}); err == nil {
		if !c.conn.Send(net.MessageId(RPC_MSGID_CALL), body) {
			fwlogger.E("-- call rpc method %s failed! send message failed --", method)
		} else {
			return true
		}
	} else {
		fwlogger.E("-- call rpc method %s failed! marshal message failed --", method)
	}
	return false
}

func (c *Client) handleMessage(conn net.IConn, msgId net.MessageId, body net.MessageBody) {
	switch RPC_MSGID(msgId) {
	case RPC_MSGID(net.MSGID_CONNECT_CONNECTED):
		fwlogger.D("-- connect success %s --", c.endName)
		// send handshake to RPC server
		if body, err := defaultCodec.Marshal(&RPCHandShake{EndName: c.endName}); err == nil {
			if !conn.Send(net.MessageId(RPC_MSGID_HANDSHAKE), body) {
				fwlogger.E("-- send handshake %s failed! send message failed --", c.endName)
			}
		} else {
			fwlogger.E("-- send handshake %s failed! marshal message failed --", c.endName)
		}
	case RPC_MSGID(net.MSGID_CONNECT_DISCONNECT):
		c.conn = nil
		c.state = net.STATE_CLOSED
		if !c.forceClose {
			c.reconnect()
		}
	case RPC_MSGID_CALL:
		call := &RPCMethodCall{}
		if err := defaultCodec.Unmarshal(body, call); err != nil {
			fwlogger.E("-- RPC noreturn call failed! -- ")
			return
		}
		fwlogger.D("@call method %s(%d) with args %s", call.Method, call.CallSeq, call.Args)
		ctx := &Context{conn: conn, callSeq: call.CallSeq, Method: call.Method, Args: call.Args}
		callMethod(c.endName, call.Method, ctx)
	case RPC_MSGID_RETURN:
		ret := &RPCMethodReturn{}
		if err := defaultCodec.Unmarshal(body, ret); err != nil {
			fwlogger.E("-- RPC return failed! -- ")
			return
		}
		fwlogger.D("@receive return from RPC method %s(%d) with result %s", ret.Method, ret.CallSeq, ret.Returns)
		c.receiveResult(utils.Identity(ret.CallSeq), &Result{Args: ret.Returns})
	}
}

func (c *Client) reconnect() {
	if net.STATE_CLOSED == c.state {
		fwlogger.D("@reconnect RPC, remote address=%s", c.addr)
		if net.GoConnect(net.Tcp, c.addr, c.handleMessage) == nil {
			time.Sleep(1 * time.Second)
			c.state = net.STATE_CLOSED
			c.reconnect()
		} else {
			c.state = net.STATE_CONNECTED
		}
	}
}
