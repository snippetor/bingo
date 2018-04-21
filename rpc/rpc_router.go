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
package rpc

import (
	"github.com/snippetor/bingo/log/fwlogger"
	"reflect"
)

type Router struct {
	routes map[string]*reflect.Method
}

func (r *Router) RegisterController(target string, c IController) {
	c.makeRoutes(target, r)
}

func (r *Router) Route(target, methodName string, routeFunc *reflect.Method) {
	if r.routes == nil {
		r.routes = make(map[string]*reflect.Method)
	}
	r.routes[r.makeKey(target, methodName)] = routeFunc
}

func (r *Router) Invoke(target, methodName string, ctx *Context) {
	if r.routes == nil {
		return
	}
	k := r.makeKey(target, methodName)
	if r.routes[k] != nil {
		in := []reflect.Value{reflect.ValueOf(ctx)}
		r.routes[k].Func.Call(in)
	}
}

func (r *Router) makeKey(target, methodName string) string {
	return target + "." + methodName
}

func (r *Router) Dump() {
	fwlogger.I("########## RPC Router Dump ##########")
	if r.routes != nil {
		for k, v := range r.routes {
			fwlogger.I("## uri=%s, func=%#v", k, v)
		}
	}
	fwlogger.I("########## END ##########")
}
