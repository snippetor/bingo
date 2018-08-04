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

package route

import (
	"github.com/snippetor/bingo/net"
	"github.com/snippetor/bingo/log/fwlogger"
	"github.com/snippetor/bingo/utils"
	"reflect"
	"go/types"
)



type Context struct {
	callSeq        uint32
	RemoteNodeName string
	conn           net.IConn
	Method         string
	Args
}

func (c *Context) Return(r *Args) {
	res := make(map[string]*RPCValue)
	r.ToRPCMap(&res)
	body := defaultCodec.Marshal(&RPCMethodReturn{CallSeq: c.callSeq, Method: c.Method, Returns: res})
	if !c.conn.Send(net.MessageId(RPC_MSGID_RETURN), body) {
		fwlogger.E("-- return rpc method %s failed! send message failed --", c.Method)
	}
}

func (c *Context) ReturnNil() {
	body := defaultCodec.Marshal(&RPCMethodReturn{CallSeq: c.callSeq, Method: c.Method, Returns: nil})
	if !c.conn.Send(net.MessageId(RPC_MSGID_RETURN), body) {
		fwlogger.E("-- return rpc method %s failed! send message failed --", c.Method)
	}
}

func (c *Context) GetConnectionIdentify() uint32 {
	return c.conn.Identity()
}
