package utils

import "sync"

type ValueMap interface {
	Put(key string, value interface{})
	Get(key string) Value
	Has(key string) bool
	Del(key string)
	Range(f func(k string, v Value) bool)
}

type vm struct {
	inner map[string]Value
}

func NewValueMap() ValueMap {
	return &vm{make(map[string]Value)}
}

func (m *vm) Put(key string, value interface{}) {
	v := NewValue()
	v.Set(value)
	m.inner[key] = v
}

func (m *vm) Get(key string) Value {
	if v, ok := m.inner[key]; ok {
		return v
	}
	return nil
}

func (m *vm) Has(key string) bool {
	if _, ok := m.inner[key]; ok {
		return true
	}
	return false
}

func (m *vm) Del(key string) {
	delete(m.inner, key)
}

func (m *vm) Range(f func(k string, v Value) bool) {
	for key, value := range m.inner {
		if !f(key, value) {
			break
		}
	}
}

// Concurrent ValueMap
type concurrent struct {
	inner *sync.Map
}

func NewConcurrentValueMap() ValueMap {
	return &concurrent{new(sync.Map)}
}

func (m *concurrent) Put(key string, value interface{}) {
	m.inner.Store(key, value)
}

func (m *concurrent) Get(key string) Value {
	if v, ok := m.inner.Load(key); ok {
		return v.(Value)
	}
	return nil
}

func (m *concurrent) Has(key string) bool {
	if _, ok := m.inner.Load(key); ok {
		return true
	}
	return false
}

func (m *concurrent) Del(key string) {
	m.inner.Delete(key)
}

func (m *concurrent) Range(f func(k string, v Value) bool) {
	m.inner.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(Value))
	})
}
