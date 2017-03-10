package net

import (
	"sync"
)

// 消息ID和消息结构体对应集合
type protoCollection struct {
	m map[MessageId]interface{}
	sync.RWMutex
}

func (c *protoCollection) put(id MessageId, v interface{}) {
	c.RLock()
	defer c.RUnlock()
	if c.m == nil {
		c.m = make(map[MessageId]interface{}, 0)
	}
	c.m[id] = v
}

func (c *protoCollection) get(id MessageId) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	if c.m == nil {
		return nil, false
	}
	v, ok := c.m[id]
	return v, ok
}

func (c *protoCollection) del(id MessageId) {
	c.RLock()
	defer c.RUnlock()
	if c.m != nil {
		delete(c.m, id)
	}
}

func (c *protoCollection) clear() {
	c.RLock()
	defer c.RUnlock()
	if c.m != nil {
		c.m = nil
	}
}
