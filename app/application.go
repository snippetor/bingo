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
)

type Application interface {
	Name() string
	AddModule(module.Module)
	GetModule(module.Module) module.Module
	RPC() module.RPCModule
	Service() module.ServiceModule
	Log() module.LogModule
	Config() utils.ValueMap
}

var _ Application = (*application)(nil)

// model
type application struct {
	name    string
	config  utils.ValueMap
	modules map[string]module.Module
}

func New(name string, config utils.ValueMap) Application {
	a := &application{name, config, make(map[string]module.Module)}
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

func (a *application) RPC() module.RPCModule {
	return a.modules["*module.RPCModule"].(module.RPCModule)
}

func (a *application) Service() module.ServiceModule {
	return a.modules["*module.ServiceModule"].(module.ServiceModule)
}

func (a *application) Log() module.LogModule {
	return a.modules["*module.LogModule"].(module.LogModule)
}

func (a *application) Config() utils.ValueMap {
	return a.config
}

func (a *application) Destroy() {
	for _, m := range a.modules {
		if m != nil {
			m.Close()
		}
	}
}
