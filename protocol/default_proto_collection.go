package protocol

import (
	"sync"
	"github.com/snippetor/bingo/net"
	"reflect"
	"strconv"
)

// 消息ID和消息结构体对应集合
type DefaultProtoCollection struct {
	m map[string]reflect.Type
	sync.RWMutex
}

func (c *DefaultProtoCollection) PutDefault(id net.MessageId, v interface{}) {
	c.Put(id, v, "")
}

func (c *DefaultProtoCollection) Put(id net.MessageId, v interface{}, protoVersion string) {
	c.Lock()
	defer c.Unlock()
	if c.m == nil {
		c.m = make(map[string]reflect.Type, 0)
	}
	var t reflect.Type
	if t = reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	c.m[c.makeKey(id, protoVersion)] = t
}

func (c *DefaultProtoCollection) GetDefault(id net.MessageId) (interface{}, bool) {
	return c.Get(id, "")
}

func (c *DefaultProtoCollection) Get(id net.MessageId, protoVersion string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	if c.m == nil {
		return nil, false
	}
	t, ok := c.m[c.makeKey(id, protoVersion)]
	if ok {
		return reflect.New(t).Interface(), ok
	} else {
		t, ok = c.m[c.makeKey(id, "")]
		if ok {
			return reflect.New(t).Interface(), ok
		}
		return nil, false
	}
}

func (c *DefaultProtoCollection) RemoveDefault(id net.MessageId) {
	c.Remove(id, "")
}

func (c *DefaultProtoCollection) Remove(id net.MessageId, protoVersion string) {
	c.Lock()
	defer c.Unlock()
	if c.m != nil {
		delete(c.m, c.makeKey(id, protoVersion))
	}
}

func (c *DefaultProtoCollection) Clear() {
	c.Lock()
	defer c.Unlock()
	if c.m != nil {
		c.m = nil
	}
}

func (c *DefaultProtoCollection) Size() int {
	c.RLock()
	defer c.RUnlock()
	if c.m != nil {
		return len(c.m)
	}
	return 0
}

func (c *DefaultProtoCollection) makeKey(id net.MessageId, protoVersion string) string {
	if protoVersion == "" {
		return strconv.Itoa(int(id))
	}
	return strconv.Itoa(int(id)) + "-" + protoVersion
}
