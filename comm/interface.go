package comm

import "sync"

// 可配置接口
type Configable struct {
	config map[string]string
}

func (c *Configable) SetConfig(key, value string) {
	if c.config == nil {
		c.config = make(map[string]string)
	}
	c.config[key] = value
}

func (c *Configable) GetConfig(key string) (string, bool) {
	v, ok := c.config[key]
	return v, ok
}