package utils

import "sync"

type ValueMap struct {
	inner map[string]*Value
}

func (m *ValueMap) Put(key string, value interface{}) {
	if m.inner == nil {
		m.inner = make(map[string]*Value)
	}
	m.inner[key] = &Value{inner: value}
}

func (m *ValueMap) Get(key string) *Value {
	if m.inner == nil {
		return nil
	}
	if v, ok := m.inner[key]; ok {
		return v
	}
	return nil
}

func (m *ValueMap) Has(key string) bool {
	if m.inner == nil {
		return false
	}
	if v, ok := m.inner[key]; ok {
		return true
	}
	return false
}

func (m *ValueMap) Del(key string) {
	if m.inner == nil {
		return
	}
	delete(m.inner, key)
}

func (m *ValueMap) Range(f func(k string, v *Value) bool) {
	for key, value := range m.inner {
		if !f(key, value) {
			break
		}
	}
}

type ConcurrentValueMap struct {
	inner *sync.Map
}

func (m *ConcurrentValueMap) Put(key string, value interface{}) {
	if m.inner == nil {
		m.inner = &sync.Map{}
	}
	m.inner.Store(key, value)
}

func (m *ConcurrentValueMap) Get(key string) *Value {
	if m.inner == nil {
		return nil
	}
	if v, ok := m.inner.Load(key); ok {
		return v.(*Value)
	}
	return nil
}

func (m *ConcurrentValueMap) Del(key string) {
	if m.inner == nil {
		return
	}
	m.inner.Delete(key)
}

func (m *ConcurrentValueMap) Range(f func(k string, v *Value) bool) {
	m.inner.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(*Value))
	})
}
