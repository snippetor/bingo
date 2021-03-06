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
	"github.com/snippetor/bingo/log"
	"strings"
	"strconv"
	"github.com/valyala/fasthttp"
	"github.com/snippetor/bingo/rpc"
	"github.com/snippetor/bingo/net"
	"github.com/snippetor/bingo/codec"
	context2 "context"
	"github.com/snippetor/bingo/config"
	"flag"
	"os"
	"path/filepath"
	"fmt"
)

type Application interface {
	Name() string
	Run()
	ShutDown()

	// mvc objects
	RegisterController(...interface{})

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

	// 系统配置
	BingoConfig() *config.BingoConfig
	// 获取运行环境
	Env() string
}

var _ Application = (*application)(nil)

// model
type application struct {
	name          string
	config        utils.ValueMap
	modules       map[string]module.Module
	controllers   []interface{}
	defaultRouter Router

	rpcCtxPool     *Pool
	serviceCtxPool *Pool
	webApiCtxPool  *Pool

	globalMiddleWares Handlers
	endRunning        chan bool

	bingoConfig    *config.BingoConfig
	env            string
	configFilePath string
}

var (
	n = flag.String("n", "", "app name")
	e = flag.String("e", "", "env")
	c = flag.String("c", "", "config file")
)

func New() Application {
	flag.Parse()

	var appName string
	if n == nil {
		appName = strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0]))
	} else {
		appName = *n
	}

	var env string
	if e == nil {
		env = ""
	} else {
		env = *e
	}

	var conf string
	if c == nil {
		if env == "" {
			conf = ".bingo.js"
		} else {
			conf = ".bingo." + env + ".js"
		}
	} else {
		conf = *c
	}

	a := &application{
		name:           appName,
		modules:        make(map[string]module.Module),
		defaultRouter:  NewRouter(),
		endRunning:     make(chan bool, 1),
		env:            env,
		configFilePath: conf,
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

func (a *application) Run() {
	// 解析配置文件
	p := config.JsParser{}
	a.bingoConfig = p.Parse(a.configFilePath)
	// 找到app配置
	appConfig := a.bingoConfig.FindApp(a.name)
	if appConfig == nil {
		panic(fmt.Sprintf("run app failed! not found app config by fileName %s ", a.name))
	}
	// log
	loggers := module.Loggers{}
	for name, c := range appConfig.Logs {
		loggers[name] = log.NewLogger(c)
	}
	a.AddModule(module.NewLogModule(loggers))
	// default logger
	if defLogger, ok := loggers["default"]; ok && a.bingoConfig.Config.EnableBingoLog {
		log.DefaultLogger = defLogger
	}
	log.DefaultLogger.SetLevel(a.bingoConfig.Config.BingoLogLevel)
	// config
	a.config = utils.NewValueMap()
	for k, v := range appConfig.Config {
		a.config.Put(k, v)
	}
	// db
	mongoes := make(map[string]module.MongoDB)
	mysqls := make(map[string]module.MySqlDB)
	for k, c := range appConfig.Db {
		switch strings.ToLower(c.Type) {
		case "mongo":
			mongoes[k] = module.NewMongoDB(c.Addr, c.User, c.Pwd, c.Db)
		case "mysql":
			mysqls[k] = module.NewMysqlDB(c.Addr, c.User, c.Pwd, c.Db, c.TbPrefix)
		}
	}
	if len(mongoes) > 0 {
		a.AddModule(module.NewMongoModule(mongoes))
	}
	if len(mysqls) > 0 {
		a.AddModule(module.NewMysqlModule(mysqls))
	}
	// init mvc objects
	for _, obj := range a.controllers {
		if IsController(obj) {
			builder := newRouterBuild(a.defaultRouter, obj)
			obj.(Controller).Route(builder)
			builder.Build()
		}
	}
	// rpc
	var rpcServer *rpc.Server
	if appConfig.RpcPort > 0 {
		rpcServer = &rpc.Server{}
		rpcServer.Listen(appConfig.Name, appConfig.Package, appConfig.RpcPort, func(server *rpc.Server) {
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
	for _, serverPkg := range appConfig.RpcTo {
		c := &rpc.Client{}
		servConfigs := a.bingoConfig.FindAppsByPackage(serverPkg)
		var address []string
		for i := range servConfigs {
			address = append(address, fmt.Sprintf("%s:%d", servConfigs[i].Domain, servConfigs[i].RpcPort))
		}
		c.Connect(appConfig.Name, appConfig.Package, serverPkg, address)
		rpcClients = append(rpcClients, c)
	}
	a.AddModule(module.NewRPCModule(appConfig.Name, rpcClients, rpcServer))
	// service
	services := module.Services{}
	for k, s := range appConfig.Service {
		var c codec.Codec
		if strings.ToLower(s.Codec) == "json" {
			c = codec.NewCodec(codec.Json)
		} else if strings.ToLower(s.Codec) == "protobuf" {
			c = codec.NewCodec(codec.Protobuf)
		} else {
			c = codec.NewCodec(codec.Json)
		}

		var n net.Protocol
		if strings.ToLower(s.Net) == "tcp" {
			n = net.Tcp
		} else if strings.ToLower(s.Net) == "ws" {
			n = net.WebSocket
		} else if strings.ToLower(s.Net) == "kcp" {
			n = net.Kcp
		} else if strings.ToLower(s.Net) == "http" {
			go func() {
				log.D("http service start on %s", strconv.Itoa(s.Port))
				if err := fasthttp.ListenAndServe(":"+strconv.Itoa(s.Port), func(req *fasthttp.RequestCtx) {
					log.D("====> %s %s", string(req.Path()), string(req.Request.Body()))
					ctx := a.WebApiCtxPool().Acquire().(*WebApiContext)
					defer a.WebApiCtxPool().Release(ctx)
					ctx.RequestCtx = req
					ctx.Codec = c
					a.defaultRouter.OnHandleRequest(ctx)
				}); err != nil {
					log.E("-- startup http service failed! %s --", err.Error())
				}
			}()
			return
		} else {
			log.E("Wrong net protocol, must be one of http, tcp, kcp, ws.")
			return
		}
		serv, err := net.GoListen(n, s.Port, func(conn net.Conn, msgId net.MessageId, body net.MessageBody) {
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
		if err != nil {
			log.E("Startup %s server failed, %s, %s", s.Net, k, err.Error())
			return
		}
		services[k] = serv
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

func (a *application) RegisterController(objs ...interface{}) {
	a.controllers = append(a.controllers, objs)
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

func (a *application) BingoConfig() *config.BingoConfig {
	return a.bingoConfig
}

func (a *application) Env() string {
	return a.env
}

func (a *application) Destroy() {
	for _, m := range a.modules {
		if m != nil {
			m.Close()
		}
	}
}
