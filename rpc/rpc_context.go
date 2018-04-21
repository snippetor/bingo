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
	"github.com/snippetor/bingo/log/fwlogger"
	"github.com/snippetor/bingo/utils"
	"reflect"
)

const (
	TYPE_INT           = iota
	TYPE_INT8
	TYPE_INT16
	TYPE_INT32
	TYPE_INT64
	TYPE_UINT
	TYPE_UINT8
	TYPE_UINT16
	TYPE_UINT32
	TYPE_UINT64
	TYPE_FLOAT32
	TYPE_FLOAT64
	TYPE_STRING
	TYPE_BOOL
	TYPE_INT32_ARRAY
	TYPE_INT64_ARRAY
	TYPE_UINT32_ARRAY
	TYPE_UINT64_ARRAY
	TYPE_FLOAT32_ARRAY
	TYPE_FLOAT64_ARRAY
	TYPE_STRING_ARRAY
	TYPE_BOOL_ARRAY
	TYPE_BYTE_ARRAY
)

type Args struct {
	utils.ValueMap
}

func (a *Args) ToRPCMap(m *map[string]*RPCValue) {
	a.Range(func(k string, v *utils.Value) bool {
		kind := reflect.TypeOf(v.Get()).Kind()
		switch kind {
		case reflect.Int:
			(*m)[k] = &RPCValue{Kind: TYPE_INT, I32: int32(v.GetInt())}
		case reflect.Int8:
			(*m)[k] = &RPCValue{Kind: TYPE_INT8, I32: int32(v.GetInt8())}
		case reflect.Int16:
			(*m)[k] = &RPCValue{Kind: TYPE_INT16, I32: int32(v.GetInt16())}
		case reflect.Int32:
			(*m)[k] = &RPCValue{Kind: TYPE_INT32, I32: v.GetInt32()}
		case reflect.Int64:
			(*m)[k] = &RPCValue{Kind: TYPE_INT64, I64: v.GetInt64()}
		case reflect.Uint:
			(*m)[k] = &RPCValue{Kind: TYPE_UINT, U32: uint32(v.GetUint())}
		case reflect.Uint8:
			(*m)[k] = &RPCValue{Kind: TYPE_UINT8, U32: uint32(v.GetUint8())}
		case reflect.Uint16:
			(*m)[k] = &RPCValue{Kind: TYPE_UINT16, U32: uint32(v.GetUint16())}
		case reflect.Uint32:
			(*m)[k] = &RPCValue{Kind: TYPE_UINT32, U32: v.GetUint32()}
		case reflect.Uint64:
			(*m)[k] = &RPCValue{Kind: TYPE_UINT64, U64: v.GetUint64()}
		case reflect.Float32:
			(*m)[k] = &RPCValue{Kind: TYPE_FLOAT32, F32: v.GetFloat32()}
		case reflect.Float64:
			(*m)[k] = &RPCValue{Kind: TYPE_FLOAT64, F64: v.GetFloat64()}
		case reflect.String:
			(*m)[k] = &RPCValue{Kind: TYPE_STRING, S: v.GetString()}
		case reflect.Bool:
			(*m)[k] = &RPCValue{Kind: TYPE_BOOL, B: v.GetBool()}
		case reflect.Array, reflect.Slice:
			value := reflect.ValueOf(v.Get())
			l := value.Len()
			if l > 0 {
				t := value.Index(0).Type()
				switch t.Kind() {
				case reflect.Int32:
					(*m)[k] = &RPCValue{Kind: TYPE_INT32_ARRAY, I32A: v.GetInt32Array()}
				case reflect.Int64:
					(*m)[k] = &RPCValue{Kind: TYPE_INT64_ARRAY, I64A: v.GetInt64Array()}
				case reflect.Uint32:
					(*m)[k] = &RPCValue{Kind: TYPE_UINT32_ARRAY, U32A: v.GetUint32Array()}
				case reflect.Uint64:
					(*m)[k] = &RPCValue{Kind: TYPE_UINT64_ARRAY, U64A: v.GetUint64Array()}
				case reflect.Float32:
					(*m)[k] = &RPCValue{Kind: TYPE_FLOAT32_ARRAY, F32A: v.GetFloat32Array()}
				case reflect.Float64:
					(*m)[k] = &RPCValue{Kind: TYPE_FLOAT64_ARRAY, F64A: v.GetFloat64Array()}
				case reflect.String:
					(*m)[k] = &RPCValue{Kind: TYPE_STRING_ARRAY, Sa: v.GetStringArray()}
				case reflect.Bool:
					(*m)[k] = &RPCValue{Kind: TYPE_BOOL_ARRAY, Ba: v.GetBoolArray()}
				}
			}
		}
		return true
	})
}

func (a *Args) FromRPCMap(m *map[string]*RPCValue) {
	var res *utils.Value
	for k, v := range *m {
		res = &utils.Value{}
		switch v.Kind {
		case TYPE_INT:
			res.Set(int(v.I32))
		case TYPE_INT8:
			res.Set(int8(v.I32))
		case TYPE_INT16:
			res.Set(int16(v.I32))
		case TYPE_INT32:
			res.Set(v.I32)
		case TYPE_INT64:
			res.Set(v.I64)
		case TYPE_UINT:
			res.Set(uint(v.U32))
		case TYPE_UINT8:
			res.Set(uint8(v.U32))
		case TYPE_UINT16:
			res.Set(uint16(v.U32))
		case TYPE_UINT32:
			res.Set(v.U32)
		case TYPE_UINT64:
			res.Set(v.U64)
		case TYPE_FLOAT32:
			res.Set(v.F32)
		case TYPE_FLOAT64:
			res.Set(v.F64)
		case TYPE_STRING:
			res.Set(v.S)
		case TYPE_BOOL:
			res.Set(v.B)
		case TYPE_INT32_ARRAY:
			res.Set(v.I32A)
		case TYPE_INT64_ARRAY:
			res.Set(v.I64A)
		case TYPE_UINT32_ARRAY:
			res.Set(v.U32A)
		case TYPE_UINT64_ARRAY:
			res.Set(v.U64A)
		case TYPE_FLOAT32_ARRAY:
			res.Set(v.F32A)
		case TYPE_FLOAT64_ARRAY:
			res.Set(v.F64A)
		case TYPE_STRING_ARRAY:
			res.Set(v.Sa)
		case TYPE_BOOL_ARRAY:
			res.Set(v.Ba)
		case TYPE_BYTE_ARRAY:
			res.Set(v.S)
		}
		a.Put(k, res)
	}
}

type Context struct {
	callSeq        int32
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
