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
	"github.com/snippetor/bingo/log/fwlogger"
	"github.com/snippetor/bingo/log"
	"strings"
	"strconv"
	"github.com/valyala/fasthttp"
	"github.com/snippetor/bingo/middleware/latency"
	"os"
	"path"
	"github.com/snippetor/bingo/middleware/recover"
	"github.com/snippetor/bingo/mvc"
	"github.com/snippetor/bingo/rpc"
	"github.com/snippetor/bingo/net"
	"github.com/snippetor/bingo/codec"
	context2 "context"
	"fmt"
)

type Application interface {
	Name() string
	Run()
	ShutDown()

	// mvc objects
	RegisterMVCObjects(...interface{})

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
	name          string
	config        utils.ValueMap
	modules       map[string]module.Module
	mvcObjs       []interface{}
	defaultRouter Router

	rpcCtxPool     *Pool
	serviceCtxPool *Pool
	webApiCtxPool  *Pool

	globalMiddleWares Handlers
	endRunning        chan bool
}

func New() Application {
	wd, _ := os.Getwd()
	a := &application{
		name:          appName,
		modules:       make(map[string]module.Module),
		defaultRouter: NewRouter(),
		endRunning:    make(chan bool, 1),
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

	parse(configFilePath)
	return NewWithConfig(appName, path.Join(wd, "bingo.yaml"))
}

func (a *application) Run() {
	appConfig := findApp(a.name)
	if appConfig == nil {
		fmt.Errorf("run app failed! not found app config by name %s ", a.name)
		return
	}
	// log
	loggers := module.Loggers{}
	for name, c := range appConfig.Logs {
		loggers[name] = log.NewLogger(c)
	}
	a.AddModule(module.NewLogModule(loggers))
	// default logger
	if defLogger, ok := loggers["default"]; ok {
		log.DefaultLogger = defLogger
	}
	// config
	a.config = utils.NewValueMap()
	for k, v := range appConfig.Config {
		a.config.Put(k, v)
	}
	// middleware
	a.Use(recover.New(), latency.New())
	// db
	for _, c := range appConfig.Db {
		if c.Type == "mongo" {
			a.AddModule(module.NewMongoModule(c.Addr, c.User, c.Pwd, c.Db))
		} else if c.Type == "mysql" {
			a.AddModule(module.NewMysqlModule(c.Addr, c.User, c.Pwd, c.Db, c.TbPrefix))
		}
	}
	// init mvc objects
	for _, obj := range a.mvcObjs {
		if mvc.IsController(obj) {
			builder := newRouterBuild(a.defaultRouter, obj)
			obj.(mvc.Controller).Route(builder)
			builder.Build()
		} else if mvc.IsOrmModel(obj) {
			a.MySql().AutoMigrate(obj.(mvc.MysqlOrmModel))
		}
	}
	// rpc
	var rpcServer *rpc.Server
	if appConfig.Rpc.Port > 0 {
		rpcServer = &rpc.Server{}
		rpcServer.Listen(appConfig.Name, appConfig.ModelName, appConfig.Rpc.Port, []string{}, func(server *rpc.Server) {
			for key := range a.defaultRouter.Handlers("RPC") {
				server.RegisterFunction(key, func(c context2.Context, args []byte, reply *[]byte) error {
					ctx := a.RpcCtxPool().Acquire().(*RpcContext)
					defer a.RpcCtxPool().Release(ctx)
					ctx.args = args
					ctx.reply = reply
					a.defaultRouter.OnHandleRequest(ctx)
					return ctx.error
				})
			}
		})
	}
	var rpcClients []*rpc.Client
	for _, serverName := range appConfig.Rpc.To {
		c := &rpc.Client{}
		serverApp := findApp(serverName)
		if serverApp != nil {
			c.Connect(appConfig.Name, appConfig.ModelName, BingoConfiguration().Domains[serverApp.Domain]+":"+strconv.Itoa(serverApp.Rpc.Port), []string{})
			rpcClients = append(rpcClients, c)
		}
	}
	a.AddModule(module.NewRPCModule(appConfig.Name, rpcClients, rpcServer))
	// service
	services := module.Services{}
	for _, s := range appConfig.Services {
		var c codec.Codec
		if strings.ToLower(s.Codec) == "json" {
			c = codec.NewCodec(codec.Json)
		} else if strings.ToLower(s.Codec) == "protobuf" {
			c = codec.NewCodec(codec.Protobuf)
		} else {
			c = codec.NewCodec(codec.Json)
		}
		switch strings.ToLower(s.Type) {
		case "tcp":
			serv := net.GoListen(net.Tcp, s.Port, func(conn net.Conn, msgId net.MessageId, body net.MessageBody) {
				ctx := a.ServiceCtxPool().Acquire().(*ServiceContext)
				defer a.ServiceCtxPool().Release(ctx)
				ctx.Conn = conn
				ctx.MessageId = msgId.MsgId()
				ctx.MessageType = msgId.Type()
				ctx.MessageGroup = msgId.Group()
				ctx.MessageExtra = msgId.Extra()
				ctx.MessageBody = &MessageBodyWrapper{RawContent: body, Codec: c}
				ctx.Codec = c
				a.defaultRouter.OnHandleRequest(ctx)
			})
			services[s.Name] = serv
		case "ws":
			serv := net.GoListen(net.WebSocket, s.Port, func(conn net.Conn, msgId net.MessageId, body net.MessageBody) {
				ctx := a.ServiceCtxPool().Acquire().(*ServiceContext)
				defer a.ServiceCtxPool().Release(ctx)
				ctx.Conn = conn
				ctx.MessageId = msgId.MsgId()
				ctx.MessageType = msgId.Type()
				ctx.MessageGroup = msgId.Group()
				ctx.MessageExtra = msgId.Extra()
				ctx.MessageBody = &MessageBodyWrapper{RawContent: body, Codec: c}
				ctx.Codec = c
				a.defaultRouter.OnHandleRequest(ctx)
			})
			services[s.Name] = serv
		case "http":
			go func() {
				fwlogger.D("-- http service start on %s --", strconv.Itoa(s.Port))
				if err := fasthttp.ListenAndServe(":"+strconv.Itoa(s.Port), func(req *fasthttp.RequestCtx) {
					fwlogger.D("====> %s %s", string(req.Path()), string(req.Request.Body()))
					ctx := a.WebApiCtxPool().Acquire().(*WebApiContext)
					defer a.WebApiCtxPool().Release(ctx)
					ctx.RequestCtx = req
					ctx.Codec = c
					a.defaultRouter.OnHandleRequest(ctx)
				}); err != nil {
					fwlogger.E("-- startup http service failed! %s --", err.Error())
				}
			}()
		}
	}
	a.AddModule(module.NewServiceModule(services))

	select {
	case <-a.endRunning:
		a.Destroy()
	}
}

func (a *application) ShutDown() {
	a.endRunning <- true
}

func (a *application) Name() string {
	return a.name
}

func (a *application) RegisterMVCObjects(objs ...interface{}) {
	a.mvcObjs = append(a.mvcObjs, objs)
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
