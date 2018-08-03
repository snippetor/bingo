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
	"github.com/snippetor/bingo/net"
	"github.com/valyala/fasthttp"
	"github.com/snippetor/bingo/utils"
	"github.com/snippetor/bingo/module"
	"reflect"
)

type Application interface {
	Name() string
	AddModule(module.Module)
	RPCModule() module.RPCModule
	ServiceModule() module.ServiceModule
	LogModule() module.LogModule
	Config() utils.ValueMap

	destroy()
}

var _ Application = (*app)(nil)

// model
type app struct {
	name    string
	config  utils.ValueMap
	modules map[string]module.Module
}

func New(name string, config utils.ValueMap) Application {
	a := &app{name, config, make(map[string]module.Module)}
	return a
}

func (a *app) Name() string {
	return a.name
}

func (a *app) AddModule(module module.Module) {
	a.modules[reflect.TypeOf(module).String()] = module
}

func (a *app) RPCModule() module.RPCModule {
	return a.modules["*module.RPCModule"].(module.RPCModule)
}

func (a *app) ServiceModule() module.ServiceModule {
	return a.modules["*module.ServiceModule"].(module.ServiceModule)
}

func (a *app) LogModule() module.LogModule {
	return a.modules["*module.LogModule"].(module.LogModule)
}

func (a *app) Config() utils.ValueMap {
	return a.config
}

func (a *app) destroy() {
	for _, m := range a.modules {
		if m != nil {
			m.Close()
		}
	}
}

func (a *app) OnServiceClientConnected(serviceName string, conn net.IConn) {
}

func (a *app) OnServiceClientDisconnected(serviceName string, conn net.IConn) {
}

func (a *app) OnReceiveServiceMessage(conn net.IConn, msgId net.MessageId, msgBody net.MessageBody) {
}

func (a *app) OnReceiveHttpServiceRequest(ctx *fasthttp.RequestCtx) {
}
