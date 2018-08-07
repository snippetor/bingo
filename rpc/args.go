package rpc

import (
	"github.com/snippetor/bingo/utils"
	"go/types"
	"reflect"
)

const (
	ArgTypeInt          = iota
	ArgTypeInt8
	ArgTypeInt16
	ArgTypeInt32
	ArgTypeInt64
	ArgTypeUint
	ArgTypeUint8
	ArgTypeUint16
	ArgTypeUint32
	ArgTypeUint64
	ArgTypeFloat32
	ArgTypeFloat64
	ArgTypeString
	ArgTypeBool
	ArgType32Array
	ArgType64Array
	ArgTypeUint32Array
	ArgTypeUint64Array
	ArgTypeFloat32Array
	ArgTypeFloat64Array
	ArgTypeStringArray
	ArgTypeBoolArray
)

type Args struct {
	inner utils.ValueMap
}

func (a *Args) Put(key string, value interface{}) {
	a.inner.Put(key, value)
}

func (a *Args) Get(key string) utils.Value {
	return a.inner.Get(key)
}

func (a *Args) Has(key string) bool {
	return a.inner.Has(key)
}

func (a *Args) Range(f func(k string, v utils.Value) bool) {
	a.inner.Range(f)
}

func (a *Args) ToRPCMap(m *map[string]*RPCValue) {
	if a.inner == nil {
		return
	}
	a.inner.Range(func(k string, v utils.Value) bool {
		switch v.Get().(type) {
		case int:
			(*m)[k] = &RPCValue{Kind: ArgTypeInt, I32: int32(v.GetInt())}
		case int8:
			(*m)[k] = &RPCValue{Kind: ArgTypeInt8, I32: int32(v.GetInt8())}
		case int16:
			(*m)[k] = &RPCValue{Kind: ArgTypeInt16, I32: int32(v.GetInt16())}
		case int32:
			(*m)[k] = &RPCValue{Kind: ArgTypeInt32, I32: v.GetInt32()}
		case int64:
			(*m)[k] = &RPCValue{Kind: ArgTypeInt64, I64: v.GetInt64()}
		case uint:
			(*m)[k] = &RPCValue{Kind: ArgTypeUint, U32: uint32(v.GetUint())}
		case uint8:
			(*m)[k] = &RPCValue{Kind: ArgTypeUint8, U32: uint32(v.GetUint8())}
		case uint16:
			(*m)[k] = &RPCValue{Kind: ArgTypeUint16, U32: uint32(v.GetUint16())}
		case uint32:
			(*m)[k] = &RPCValue{Kind: ArgTypeUint32, U32: v.GetUint32()}
		case uint64:
			(*m)[k] = &RPCValue{Kind: ArgTypeUint64, U64: v.GetUint64()}
		case float32:
			(*m)[k] = &RPCValue{Kind: ArgTypeFloat32, F32: v.GetFloat32()}
		case float64:
			(*m)[k] = &RPCValue{Kind: ArgTypeFloat64, F64: v.GetFloat64()}
		case string:
			(*m)[k] = &RPCValue{Kind: ArgTypeString, S: v.GetString()}
		case bool:
			(*m)[k] = &RPCValue{Kind: ArgTypeBool, B: v.GetBool()}
		case types.Array, types.Slice:
			value := reflect.ValueOf(v.Get())
			l := value.Len()
			if l > 0 {
				t := value.Index(0).Type()
				switch t.Kind() {
				case reflect.Int32:
					(*m)[k] = &RPCValue{Kind: ArgType32Array, I32A: v.GetInt32Array()}
				case reflect.Int64:
					(*m)[k] = &RPCValue{Kind: ArgType64Array, I64A: v.GetInt64Array()}
				case reflect.Uint32:
					(*m)[k] = &RPCValue{Kind: ArgTypeUint32Array, U32A: v.GetUint32Array()}
				case reflect.Uint64:
					(*m)[k] = &RPCValue{Kind: ArgTypeUint64Array, U64A: v.GetUint64Array()}
				case reflect.Float32:
					(*m)[k] = &RPCValue{Kind: ArgTypeFloat32Array, F32A: v.GetFloat32Array()}
				case reflect.Float64:
					(*m)[k] = &RPCValue{Kind: ArgTypeFloat64Array, F64A: v.GetFloat64Array()}
				case reflect.String:
					(*m)[k] = &RPCValue{Kind: ArgTypeStringArray, Sa: v.GetStringArray()}
				case reflect.Bool:
					(*m)[k] = &RPCValue{Kind: ArgTypeBoolArray, Ba: v.GetBoolArray()}
				}
			}
		}
		return true
	})
}

func (a *Args) FromRPCMap(m *map[string]*RPCValue) {
	a.inner = utils.NewValueMap()
	var res utils.Value
	for k, v := range *m {
		res = utils.NewValue()
		switch v.Kind {
		case ArgTypeInt:
			res.Set(int(v.I32))
		case ArgTypeInt8:
			res.Set(int8(v.I32))
		case ArgTypeInt16:
			res.Set(int16(v.I32))
		case ArgTypeInt32:
			res.Set(v.I32)
		case ArgTypeInt64:
			res.Set(v.I64)
		case ArgTypeUint:
			res.Set(uint(v.U32))
		case ArgTypeUint8:
			res.Set(uint8(v.U32))
		case ArgTypeUint16:
			res.Set(uint16(v.U32))
		case ArgTypeUint32:
			res.Set(v.U32)
		case ArgTypeUint64:
			res.Set(v.U64)
		case ArgTypeFloat32:
			res.Set(v.F32)
		case ArgTypeFloat64:
			res.Set(v.F64)
		case ArgTypeString:
			res.Set(v.S)
		case ArgTypeBool:
			res.Set(v.B)
		case ArgType32Array:
			res.Set(v.I32A)
		case ArgType64Array:
			res.Set(v.I64A)
		case ArgTypeUint32Array:
			res.Set(v.U32A)
		case ArgTypeUint64Array:
			res.Set(v.U64A)
		case ArgTypeFloat32Array:
			res.Set(v.F32A)
		case ArgTypeFloat64Array:
			res.Set(v.F64A)
		case ArgTypeStringArray:
			res.Set(v.Sa)
		case ArgTypeBoolArray:
			res.Set(v.Ba)
		}
		a.inner.Put(k, res)
	}
}
