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

package app

import (
	"github.com/snippetor/bingo/net"
	"github.com/snippetor/bingo/log/fwlogger"
	"github.com/snippetor/bingo/rpc"
)

type RpcContext struct {
	Context
	Conn    net.IConn
	CallSeq uint32
	Method  string
	Args    *rpc.Args
	Caller  string
}

// The only one important if you will override the Context
// with an embedded context.Context inside it.
// Required in order to run the handlers via this "*MyContext".
func (c *RpcContext) Do(handlers Handlers) {
	Do(c, handlers)
}

// The second one important if you will override the Context
// with an embedded context.Context inside it.
// Required in order to run the chain of handlers via this "*MyContext".
func (c *RpcContext) Next() {
	Next(c)
}

func (c *RpcContext) Return(r *rpc.Args) {
	res := make(map[string]*rpc.RPCValue)
	r.ToRPCMap(&res)
	body := rpc.DefaultCodec.Marshal(&rpc.RPCMethodReturn{CallSeq: c.CallSeq, Method: c.Method, Returns: res})
	if !c.Conn.Send(net.MessageId(rpc.RPC_MSGID_RETURN), body) {
		fwlogger.E("-- return rpc method %s failed! send message failed --", c.Method)
	}
}

func (c *RpcContext) ReturnNil() {
	body := rpc.DefaultCodec.Marshal(&rpc.RPCMethodReturn{CallSeq: c.CallSeq, Method: c.Method, Returns: nil})
	if !c.Conn.Send(net.MessageId(rpc.RPC_MSGID_RETURN), body) {
		fwlogger.E("-- return rpc method %s failed! send message failed --", c.Method)
	}
}
