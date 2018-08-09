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
	"github.com/snippetor/bingo/utils"
	"reflect"
	"github.com/snippetor/bingo/module"
)

type Application interface {
	Name() string

	// modules
	AddModule(module.Module)
	GetModule(module.Module) module.Module
	RPC() module.RPCModule
	Service() module.ServiceModule
	Log() module.LogModule
	MySql() module.MySqlModule
	Mongo() module.MongoModule
	Config() utils.ValueMap

	// 使用中间件，中间件将在其他Handler之前执行
	Use(middleware ...Handler)
	GlobalMiddleWares() Handlers

	RpcCtxPool() *Pool
	ServiceCtxPool() *Pool
	WebApiCtxPool() *Pool
}

var _ Application = (*application)(nil)

// model
type application struct {
	name    string
	config  utils.ValueMap
	modules map[string]module.Module

	rpcCtxPool     *Pool
	serviceCtxPool *Pool
	webApiCtxPool  *Pool

	globalMiddleWares Handlers
}

func New(name string, config utils.ValueMap) Application {
	a := &application{
		name:    name,
		config:  config,
		modules: make(map[string]module.Module),
	}
	a.rpcCtxPool = NewPool(func() Context {
		return &RpcContext{
			Context: NewContext(a),
		}
	})
	a.serviceCtxPool = NewPool(func() Context {
		return &ServiceContext{
			Context: NewContext(a),
		}
	})
	a.webApiCtxPool = NewPool(func() Context {
		return &WebApiContext{
			Context: NewContext(a),
		}
	})
	return a
}

func (a *application) Name() string {
	return a.name
}

func (a *application) AddModule(module module.Module) {
	a.modules[reflect.TypeOf(module).String()] = module
}

func (a *application) GetModule(module module.Module) module.Module {
	if m, ok := a.modules[reflect.TypeOf(module).String()]; ok {
		return m
	}
	return nil
}

func (a *application) Use(middleware ...Handler) {
	a.globalMiddleWares = append(a.globalMiddleWares, middleware...)
}

func (a *application) GlobalMiddleWares() Handlers {
	clone := make(Handlers, len(a.globalMiddleWares))
	copy(clone, a.globalMiddleWares)
	return clone
}

func (a *application) RPC() module.RPCModule {
	m, ok := a.modules["*RPCModule"]
	if ok {
		return m.(module.RPCModule)
	}
	return nil
}

func (a *application) Service() module.ServiceModule {
	m, ok := a.modules["*ServiceModule"]
	if ok {
		return m.(module.ServiceModule)
	}
	return nil
}

func (a *application) Log() module.LogModule {
	m, ok := a.modules["*LogModule"]
	if ok {
		return m.(module.LogModule)
	}
	return nil
}

func (a *application) MySql() module.MySqlModule {
	m, ok := a.modules["*MySqlModule"]
	if ok {
		return m.(module.MySqlModule)
	}
	return nil
}

func (a *application) Mongo() module.MongoModule {
	m, ok := a.modules["*MongoModule"]
	if ok {
		return m.(module.MongoModule)
	}
	return nil
}

func (a *application) Config() utils.ValueMap {
	return a.config
}

func (a *application) RpcCtxPool() *Pool {
	return a.rpcCtxPool
}

func (a *application) ServiceCtxPool() *Pool {
	return a.serviceCtxPool
}

func (a *application) WebApiCtxPool() *Pool {
	return a.webApiCtxPool
}

func (a *application) Destroy() {
	for _, m := range a.modules {
		if m != nil {
			m.Close()
		}
	}
}
