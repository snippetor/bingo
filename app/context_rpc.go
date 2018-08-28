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
	"github.com/snippetor/bingo/codec"
	"github.com/gogo/protobuf/proto"
)

type RpcContext struct {
	Context
	Method string
	args   []byte
	reply  *[]byte
	error  error
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

func (c *RpcContext) Args(i interface{}) {
	err := codec.ProtobufCodec.Unmarshal(c.args, i)
	if err != nil {
		panic(err)
	}
}

func (c *RpcContext) Return(i interface{}) {
	msg, ok := i.(proto.Message)
	if !ok {
		panic("must return proto.Message for RPC reply.")
	}
	body, err := codec.ProtobufCodec.Marshal(msg)
	if err != nil {
		panic("must return proto.Message for RPC reply.")
	}
	*c.reply = body
	c.StopExecution()
}

func (c *RpcContext) ReturnNil() {
	c.StopExecution()
}

func (c *RpcContext) CheckError(e error) {
	if e != nil {
		c.error = e
		c.StopExecution()
		panic(e)
	}
}
