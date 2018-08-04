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

package bingo

import (
	"io/ioutil"
	"encoding/json"
	"github.com/snippetor/bingo/rpc"
	"github.com/snippetor/bingo/net"
	"strings"
	"github.com/valyala/fasthttp"
	"strconv"
	"github.com/snippetor/bingo/log/fwlogger"
	"github.com/snippetor/bingo/utils"
	"github.com/snippetor/bingo/log"
	"github.com/snippetor/bingo/module"
	"github.com/snippetor/bingo/route"
	"github.com/snippetor/bingo/app"
	"github.com/snippetor/bingo/mvc"
)

type Service struct {
	Name  string
	Type  string
	Port  int
	Codec string
}

type RPC struct {
	Port int      `json:"port"`
	To   []string `json:"to"`
}

type LogConfig struct {
	Name                string
	Level               int
	OutputType          int
	OutputDir           string
	RollingType         int
	FileName            string
	FileNameDatePattern string
	FileNameExt         string
	FileMaxSize         string
	FileScanInterval    int
}

type AppConfig struct {
	Name       string
	ModelName  string                 `json:"application"`
	Domain     int
	Services   []*Service             `json:"service"`
	Rpc        *RPC                   `json:"rpc"`
	Config     map[string]interface{} `json:"config"`
	LogConfigs []*LogConfig           `json:"log"`
}

type Config struct {
	Domains []string     `json:"domains"`
	Apps    []*AppConfig `json:"apps"`
}

type StartUpFunc func(app.Application) []interface{}

var (
	config     *Config
	router     route.Router
	runningApp = make(map[string]app.Application)
	startUpFunc = make(map[string]StartUpFunc)
)

func init() {
	config = &Config{}
	router = route.NewRouter()
}

func RegisterApp(appName string, startupFunc StartUpFunc) {
	startUpFunc[appName] = startupFunc
}

// 解析配置文件
func parse(configPath string) {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		fwlogger.E("-- parse config file failed! %s --", err)
		return
	}
	if err = json.Unmarshal(content, config); err != nil {
		fwlogger.E("-- parse config file failed! %s --", err)
		return
	}
	//TODO check
}

func findApp(name string) *AppConfig {
	for _, n := range config.Apps {
		if n.Name == name {
			return n
		}
	}
	return nil
}

func runApp(appName string) {
	a := findApp(appName)
	if n == nil {
		fwlogger.E("-- run application failed! not found application by name %s --", appName)
		return
	}
	// config
	appConfig := utils.NewValueMap()
	for k, v := range a.Config {
		appConfig.Put(k, v)
	}
	// new application
	app := app.New(a.Name, appConfig)
	// init mvc objects
	objs := mvcObjects[appName]
	if objs != nil {
		for _, obj := range objs {
			if mvc.IsController(obj) {
				obj.(mvc.Controller).Route(router)
			} else {
				//TODO model
			}
		}
	}
	// context pool
	serviceCtxPool = route.NewPool(func() route.Context {
		return &route.ServiceContext{
			Context: route.NewContext(app),
		}
	})
	rpcCtxPool = route.NewPool(func() route.Context {
		return &route.RpcContext{
			Context: route.NewContext(app),
		}
	})
	webApiCtxPool = route.NewPool(func() route.Context {
		return &route.WebApiContext{
			Context: route.NewContext(app),
		}
	})
	// log
	/**
	"level": 0,
          "outputType": 3,
          "outputDir": ".",
          "rollingType": 3,
          "fileName": "dev-err",
          "fileNameDatePattern": 20060102,
          "fileNameExt": ".log",
          "fileMaxSize": "1KB",
          "fileScanInterval": 3
	 */
	loggers := module.Loggers{}
	if a.LogConfigs != nil {
		for _, c := range a.LogConfigs {
			config := log.DEFAULT_CONFIG
			if c.Level != 0 {
				config.Level = log.Level(c.Level)
			}
			if c.OutputType != 0 {
				config.OutputType = log.OutputType(c.OutputType)
			}
			if c.RollingType != 0 {
				config.LogFileRollingType = log.RollingType(c.RollingType)
			}
			if c.OutputDir != "" {
				config.LogFileOutputDir = c.OutputDir
			}
			if c.FileName != "" {
				config.LogFileName = c.FileName
			}
			if c.FileNameDatePattern != "" {
				config.LogFileNameDatePattern = c.FileNameDatePattern
			}
			if c.FileNameExt != "" {
				config.LogFileNameExt = c.FileNameExt
			}
			if c.FileMaxSize != "" {
				if i, err := strconv.ParseInt(c.FileMaxSize, 10, 64); err == nil {
					config.LogFileMaxSize = i
				} else {
					if i, err = strconv.ParseInt(c.FileMaxSize[:len(c.FileMaxSize)-2], 10, 64); err == nil {
						unit := strings.ToUpper(c.FileMaxSize[len(c.FileMaxSize)-2:])
						if unit == "KB" {
							config.LogFileMaxSize = i * log.KB
						} else if unit == "MB" {
							config.LogFileMaxSize = i * log.MB
						} else if unit == "GB" {
							config.LogFileMaxSize = i * log.GB
						} else if unit == "TB" {
							config.LogFileMaxSize = i * log.TB
						}
					}
				}
			}
			if c.FileScanInterval != 0 {
				config.LogFileScanInterval = c.FileScanInterval
			}
			loggers[c.Name] = log.NewLoggerWithConfig(config)
		}
	}
	app.AddModule(module.NewLogModule(loggers))
	// rpc
	var rpcServer *rpc.Server
	if a.Rpc.Port > 0 {
		rpcServer = &rpc.Server{}
		rpcServer.Listen(a.Name, a.ModelName, a.Rpc.Port, func(conn net.IConn, caller string, seq uint32, method string, args *rpc.Args) {
			if r, ok := router.(interface {
				OnRPCRequest(method string, ctx route.Context)
			}); ok {
				ctx := rpcCtxPool.Acquire().(*route.RpcContext)
				ctx.Conn = conn
				ctx.CallSeq = seq
				ctx.Method = method
				ctx.Args = args
				ctx.Caller = caller
				defer rpcCtxPool.Release(ctx)
				r.OnRPCRequest(method, ctx)
			}
		})
	}
	var rpcClients []*rpc.Client
	for _, serverName := range a.Rpc.To {
		c := &rpc.Client{}
		serverApp := findApp(serverName)
		if serverApp != nil {
			c.Connect(a.Name, a.ModelName, config.Domains[serverApp.Domain]+":"+strconv.Itoa(serverApp.Rpc.Port), func(conn net.IConn, caller string, seq uint32, method string, args *rpc.Args) {
				if r, ok := router.(interface {
					OnRPCRequest(method string, ctx route.Context)
				}); ok {
					ctx := rpcCtxPool.Acquire().(*route.RpcContext)
					ctx.Conn = conn
					ctx.CallSeq = seq
					ctx.Method = method
					ctx.Args = args
					ctx.Caller = caller
					defer rpcCtxPool.Release(ctx)
					r.OnRPCRequest(method, ctx)
				}
			})
			rpcClients = append(rpcClients, c)
		}
	}
	app.AddModule(module.NewRPCModule(a.Name, rpcClients, rpcServer))
	// service
	services := module.Services{}
	for _, s := range a.Services {
		switch strings.ToLower(s.Type) {
		case "tcp":
			serv := net.GoListen(net.Tcp, s.Port, func(conn net.IConn, msgId net.MessageId, body net.MessageBody) {
				if r, ok := router.(interface {
					OnServiceRequest(net.MessageId, route.Context)
				}); ok {
					ctx := serviceCtxPool.Acquire().(*route.ServiceContext)
					ctx.Conn = conn
					ctx.MessageId = msgId.MsgId()
					ctx.MessageType = msgId.Type()
					ctx.MessageGroup = msgId.Group()
					ctx.MessageExtra = msgId.Extra()
					ctx.MessageBody = body
					defer serviceCtxPool.Release(ctx)
					r.OnServiceRequest(msgId, ctx)
				}
			})
			services[s.Name] = serv
		case "ws":
			serv := net.GoListen(net.WebSocket, s.Port, func(conn net.IConn, msgId net.MessageId, body net.MessageBody) {
				if r, ok := router.(interface {
					OnServiceRequest(net.MessageId, route.Context)
				}); ok {
					ctx := serviceCtxPool.Acquire().(*route.ServiceContext)
					ctx.Conn = conn
					ctx.MessageId = msgId.MsgId()
					ctx.MessageType = msgId.Type()
					ctx.MessageGroup = msgId.Group()
					ctx.MessageExtra = msgId.Extra()
					defer serviceCtxPool.Release(ctx)
					r.OnServiceRequest(msgId, ctx)
				}
			})
			services[s.Name] = serv
		case "http":
			go func() {
				fwlogger.D("-- http service start on %s --", strconv.Itoa(s.Port))
				if err := fasthttp.ListenAndServe(":"+strconv.Itoa(s.Port), func(c *fasthttp.RequestCtx) {
					fwlogger.D("====> %s %s", string(c.Path()), string(c.Request.Body()))
					if r, ok := router.(interface {
						OnWebApiRequest(string, route.Context)
					}); ok {
						ctx := webApiCtxPool.Acquire().(*route.WebApiContext)
						ctx.RequestCtx = c
						defer webApiCtxPool.Release(ctx)
						r.OnWebApiRequest(string(c.Path()), ctx)
					}
				}); err != nil {
					fwlogger.E("-- startup http service failed! %s --", err.Error())
				}
			}()
		}
	}
	app.AddModule(module.NewServiceModule(services))
	runningApp[a.Name] = app
}

func stop(appName string) {
	n := findApp(appName)
	if n == nil {
		fwlogger.E("-- stop application failed! not found application by name %s --", appName)
		return
	}
	if m, ok := runningApp[appName]; ok {
		if r, ok := m.(interface {
			Destroy()
		}); ok {
			r.Destroy()
		}
	}
}

func runAll() {
	for _, n := range config.Apps {
		runApp(n)
	}
}

func stopAll() {
	for _, m := range runningApp {
		if m != nil {
			if r, ok := m.(interface {
				Destroy()
			}); ok {
				r.Destroy()
			}
		}
	}
}
