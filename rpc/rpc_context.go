package rpc

import (
	"github.com/snippetor/bingo/net"
	"github.com/snippetor/bingo"
	"strconv"
)

type Args map[string]string

func (a Args) Put(key, value string) {
	a[key] = value
}

func (a Args) Get(key string) (string, bool) {
	v, ok := a[key]
	return v, ok
}

func (a Args) MustGet(key, def string) string {
	if v, ok := a[key]; ok {
		return v
	}
	return def
}

func (a Args) PutInt(key string, value int) {
	a[key] = strconv.Itoa(value)
}

func (a Args) GetInt(key string) (int, bool) {
	if v, ok := a[key]; ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i, true
		}
	}
	return 0, false
}

func (a Args) MustGetInt(key string, def int) int {
	if v, ok := a.GetInt(key); ok {
		return v
	}
	return def
}

func (a Args) PutInt64(key string, value int64) {
	a[key] = strconv.FormatInt(value, 10)
}

func (a Args) GetInt64(key string) (int64, bool) {
	if v, ok := a[key]; ok {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i, true
		}
	}
	return 0, false
}

func (a Args) MustGetInt64(key string, def int64) int64 {
	if v, ok := a.GetInt64(key); ok {
		return v
	}
	return def
}

func (a Args) PutFloat32(key string, value float32) {
	a[key] = strconv.FormatFloat(float64(value), 'f', -1, 32)
}

func (a Args) GetFloat32(key string) (float32, bool) {
	if v, ok := a[key]; ok {
		if i, err := strconv.ParseFloat(v, 32); err == nil {
			return float32(i), true
		}
	}
	return 0.0, false
}

func (a Args) MustGetFloat32(key string, def float32) float32 {
	if v, ok := a.GetFloat32(key); ok {
		return v
	}
	return def
}

func (a Args) PutBool(key string, value bool) {
	a[key] = strconv.FormatBool(value)
}

func (a Args) GetBool(key string) (bool, bool) {
	if v, ok := a[key]; ok {
		if i, err := strconv.ParseBool(v); err == nil {
			return i, true
		}
	}
	return false, false
}

func (a Args) MustGetBool(key string, def bool) bool {
	if v, ok := a.GetBool(key); ok {
		return v
	}
	return def
}

type Context struct {
	conn   net.IConn
	method string
	Args
}

func (c *Context) Callback(method string, args Args) {
	if c.conn == nil {
		bingo.E("-- call rpc method %s failed! Context#conn is nil --", method)
	} else {
		if body, err := defaultCodec.Marshal(&RPCMethodCall{Method: method, Args: args}); err == nil {
			if !c.conn.Send(net.MessageId(RPC_MSGID_CALL), body) {
				bingo.E("-- call rpc method %s.%s failed! send message failed --", method)
			}
		} else {
			bingo.E("-- call rpc method %s failed! marshal message failed --", method)
		}
	}
}
