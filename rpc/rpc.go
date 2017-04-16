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
	"strings"
)

var (
	defaultCodec codec.ICodec
)

func init() {
	defaultCodec = codec.NewCodec(codec.Protobuf)
}

type RPCEnd struct {
	connId utils.Identity
	name   string
}

type Server struct {
	endName    string
	serv       net.IServer
	end        map[string]*RPCEnd
	l          *sync.RWMutex
	identifier *utils.Identifier
	callSyncWorker
}

func (s *Server) Listen(endName string, port int) {
	s.endName = endName
	s.l = &sync.RWMutex{}
	s.end = make(map[string]*RPCEnd)
	s.identifier = utils.NewIdentifier(2)
	if s.serv = net.GoListen(net.Tcp, port, s.handleMessage); s.serv == nil {
		fwlogger.E("-- start rpc server failed! --")
	}
}

func (s *Server) Close() {
	s.serv.Close()
}

func (s *Server) Call(endName, method string, args *Args, callback RPCCallback) bool {
	fwlogger.D("Call: %s, %s, %s", endName, method, args)
	if s.serv == nil {
		fwlogger.E("-- call rpc method %s failed! rpc server no startup --", method)
		return false
	}
	s.l.RLock()
	defer s.l.RUnlock()
	if end, ok := s.end[endName]; ok {
		if conn, ok := s.serv.GetConnection(end.connId); ok {
			seq := s.identifier.GenIdentity()
			if body, err := defaultCodec.Marshal(&RPCMethodCall{CallSeq: int32(seq), Method: method, Args: *args}); err == nil {
				if !conn.Send(net.MessageId(RPC_MSGID_CALL), body) {
					fwlogger.E("-- call rpc method %s failed! send message failed --", method)
				} else {
					s.waitingResult(&callTask{seq: seq, conn: conn, msg: body, c: callback, t: time.Now().UnixNano()})
					return true
				}
			} else {
				fwlogger.E("-- call rpc method %s failed! marshal message failed --", method)
			}
		} else {
			fwlogger.E("-- call rpc method %s failed! no connection for call --", method)
		}
	} else {
		fwlogger.E("-- call rpc method failed! no end connected with name is %s --", endName)
	}
	return false
}

func (s *Server) CallNoReturn(endName, method string, args *Args) bool {
	fwlogger.D("CallNoReturn: %s, %s, %s", endName, method, args)
	if s.serv == nil {
		fwlogger.E("-- call rpc noreturn method %s failed! rpc server no startup --", method)
		return false
	}
	s.l.RLock()
	defer s.l.RUnlock()
	if end, ok := s.end[endName]; ok {
		if conn, ok := s.serv.GetConnection(end.connId); ok {
			if body, err := defaultCodec.Marshal(&RPCMethodCall{CallSeq: int32(s.identifier.GenIdentity()), Method: method, Args: *args}); err == nil {
				if !conn.Send(net.MessageId(RPC_MSGID_CALL), body) {
					fwlogger.E("-- call rpc noreturn method %s failed! send message failed --", method)
				} else {
					return true
				}
			} else {
				fwlogger.E("-- call rpc noreturn method %s failed! marshal message failed --", method)
			}
		} else {
			fwlogger.E("-- call rpc noreturn method %s failed! no connection for call --", method)
		}
	} else {
		fwlogger.E("-- call rpc noreturn method failed! no end connected with name is %s --", endName)
	}
	return false
}

func (s *Server) handleMessage(conn net.IConn, msgId net.MessageId, body net.MessageBody) {
	switch RPC_MSGID(msgId) {
	case net.MSGID_CONNECT_DISCONNECT:
		for name, e := range s.end {
			if e.connId == conn.Identity() {
				s.l.Lock()
				delete(s.end, name)
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
		end := &RPCEnd{connId: conn.Identity(), name: handshake.EndName}
		s.l.Lock()
		s.end[handshake.EndName] = end
		s.l.Unlock()
	case RPC_MSGID_CALL:
		call := &RPCMethodCall{}
		if err := defaultCodec.Unmarshal(body, call); err != nil {
			fwlogger.E("-- RPC call failed! -- ")
			return
		}
		ctx := &Context{conn: conn, Method: call.Method, Args: call.Args}
		r := callMethod(s.endName, call.Method, ctx)
		if body, err := defaultCodec.Marshal(&RPCMethodReturn{CallSeq: call.CallSeq, Method: call.Method, Returns: r.Args}); err == nil {
			if !conn.Send(net.MessageId(RPC_MSGID_RETURN), body) {
				fwlogger.E("-- return rpc method %s failed! send message failed --", call.Method)
			}
		} else {
			fwlogger.E("-- return rpc method %s failed! marshal message failed --", call.Method)
		}
	case RPC_MSGID_CALL_NORETURN:
		call := &RPCMethodCall{}
		if err := defaultCodec.Unmarshal(body, call); err != nil {
			fwlogger.E("-- RPC noreturn call failed! -- ")
			return
		}
		fwlogger.D("@call noreturn method %s(%d) with args %s", call.Method, call.CallSeq, call.Args)
		ctx := &Context{conn: conn, Method: call.Method, Args: call.Args}
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

func (s *Server) HasEndName(endName string) bool {
	if s.end == nil {
		return false
	}
	if _, ok := s.end[endName]; ok {
		return true
	}
	return false
}

func (s *Server) GetEndNamesWithPrefix(endNamePrefix string) []string {
	names := make([]string, 0)
	for name := range s.end {
		if strings.HasPrefix(name, endNamePrefix) {
			names = append(names, name)
		}
	}
	return names
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
			fwlogger.E("-- RPC call failed! -- ")
			return
		}
		ctx := &Context{conn: conn, Method: call.Method, Args: call.Args}
		r := callMethod(c.endName, call.Method, ctx)
		if body, err := defaultCodec.Marshal(&RPCMethodReturn{CallSeq: call.CallSeq, Method: call.Method, Returns: r.Args}); err == nil {
			if !conn.Send(net.MessageId(RPC_MSGID_RETURN), body) {
				fwlogger.E("-- return rpc method %s failed! send message failed --", call.Method)
			}
		} else {
			fwlogger.E("-- return rpc method %s failed! marshal message failed --", call.Method)
		}
	case RPC_MSGID_CALL_NORETURN:
		call := &RPCMethodCall{}
		if err := defaultCodec.Unmarshal(body, call); err != nil {
			fwlogger.E("-- RPC noreturn call failed! -- ")
			return
		}
		fwlogger.D("@call noreturn method %s(%d) with args %s", call.Method, call.CallSeq, call.Args)
		ctx := &Context{conn: conn, Method: call.Method, Args: call.Args}
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
