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
	"github.com/snippetor/bingo/module"
	"reflect"
	"github.com/snippetor/bingo/route"
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
	Use(middleware ...route.Handler)
	GlobalMiddleWares() route.Handlers

	RpcCtxPool() *route.Pool
	ServiceCtxPool() *route.Pool
	WebApiCtxPool() *route.Pool
}

var _ Application = (*application)(nil)

// model
type application struct {
	name    string
	config  utils.ValueMap
	modules map[string]module.Module

	rpcCtxPool     *route.Pool
	serviceCtxPool *route.Pool
	webApiCtxPool  *route.Pool

	globalMiddleWares route.Handlers
}

func New(name string, config utils.ValueMap) Application {
	a := &application{
		name:    name,
		config:  config,
		modules: make(map[string]module.Module),
	}
	a.rpcCtxPool = route.NewPool(func() route.Context {
		return &route.RpcContext{
			Context: route.NewContext(a),
		}
	})
	a.serviceCtxPool = route.NewPool(func() route.Context {
		return &route.ServiceContext{
			Context: route.NewContext(a),
		}
	})
	a.webApiCtxPool = route.NewPool(func() route.Context {
		return &route.WebApiContext{
			Context: route.NewContext(a),
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

func (a *application) Use(middleware ...route.Handler) {
	a.globalMiddleWares = append(a.globalMiddleWares, middleware...)
}

func (a *application) GlobalMiddleWares() route.Handlers {
	clone := make(route.Handlers, len(a.globalMiddleWares))
	copy(clone, a.globalMiddleWares)
	return clone
}

func (a *application) RPC() module.RPCModule {
	return a.modules["*module.RPCModule"].(module.RPCModule)
}

func (a *application) Service() module.ServiceModule {
	return a.modules["*module.ServiceModule"].(module.ServiceModule)
}

func (a *application) Log() module.LogModule {
	return a.modules["*module.LogModule"].(module.LogModule)
}

func (a *application) MySql() module.MySqlModule {
	return a.modules["*module.MySqlModule"].(module.MySqlModule)
}

func (a *application) Mongo() module.MongoModule {
	return a.modules["*module.MongoModule"].(module.MongoModule)
}

func (a *application) Config() utils.ValueMap {
	return a.config
}

func (a *application) RpcCtxPool() *route.Pool {
	return a.rpcCtxPool
}

func (a *application) ServiceCtxPool() *route.Pool {
	return a.serviceCtxPool
}

func (a *application) WebApiCtxPool() *route.Pool {
	return a.webApiCtxPool
}

func (a *application) Destroy() {
	for _, m := range a.modules {
		if m != nil {
			m.Close()
		}
	}
}
