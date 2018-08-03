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

type App struct {
	Name       string
	ModelName  string                 `json:"app"`
	Domain     int
	Services   []*Service             `json:"service"`
	Rpc        []*RPC                 `json:"rpc"`
	Config     map[string]interface{} `json:"config"`
	LogConfigs []*LogConfig           `json:"log"`
}

type Config struct {
	Domains []string `json:"domains"`
	Apps    []*App   `json:"apps"`
}

var (
	config *Config
	apps   = make(map[string]Application)
)

func init() {
	config = &Config{}
}

// 解析配置文件
func Parse(configPath string) {
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

func findApp(name string) *App {
	for _, n := range config.Apps {
		if n.Name == name {
			return n
		}
	}
	return nil
}

func runApp(n *App) {
	// create app
	m, ok := getappModel(n.ModelName)
	if !ok {
		fwlogger.E("-- not found model by name %s --", n.ModelName)
		return
	}
	m.setappName(n.Name)
	m.initModules()

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
	if n.LogConfigs != nil {
		loggers := make(map[string]*log.Logger)
		for _, c := range n.LogConfigs {
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
		m.setLoggers(loggers)
	}

	// config
	vm := &utils.ValueMap{}
	for k, v := range n.Config {
		vm.Put(k, v)
	}
	m.setConfig(vm)

	// init
	m.OnInit()

	// rpc
	if n.RpcPort > 0 {
		s := &rpc.Server{}
		s.Listen(n.Name, n.ModelName, n.RpcPort)
		m.setRPCServer(s)
	}
	for _, rpcServerapp := range n.RpcTo {
		c := &rpc.Client{}
		serverapp := findApp(rpcServerapp)
		if serverapp != nil {
			c.Connect(n.Name, n.ModelName, config.Domains[serverapp.Domain]+":"+strconv.Itoa(serverapp.RpcPort))
			m.putRPCClient(serverapp.Name, c)
		}
	}
	// service
	for _, s := range n.Services {
		switch strings.ToLower(s.Type) {
		case "tcp":
			serv := net.GoListen(net.Tcp, s.Port, func(conn net.IConn, msgId net.MessageId, body net.MessageBody) {
				switch msgId {
				case net.MSGID_CONNECT_CONNECTED:
					m.OnServiceClientConnected(s.Name, conn)
				case net.MSGID_CONNECT_DISCONNECT:
					m.OnServiceClientDisconnected(s.Name, conn)
				default:
					m.OnReceiveServiceMessage(conn, msgId, body)
				}
			})
			m.putService(s.Name, serv)
		case "ws":
			serv := net.GoListen(net.WebSocket, s.Port, func(conn net.IConn, msgId net.MessageId, body net.MessageBody) {
				switch msgId {
				case net.MSGID_CONNECT_CONNECTED:
					m.OnServiceClientConnected(s.Name, conn)
				case net.MSGID_CONNECT_DISCONNECT:
					m.OnServiceClientDisconnected(s.Name, conn)
				default:
					m.OnReceiveServiceMessage(conn, msgId, body)
				}
			})
			m.putService(s.Name, serv)
		case "http":
			go func() {
				fwlogger.D("-- http service start on %s --", strconv.Itoa(s.Port))
				if err := fasthttp.ListenAndServe(":"+strconv.Itoa(s.Port), func(ctx *fasthttp.RequestCtx) {
					fwlogger.D("====> %s %s", string(ctx.Path()), string(ctx.Request.Body()))
					m.OnReceiveHttpServiceRequest(ctx)
				}); err != nil {
					fwlogger.E("-- startup http service failed! %s --", err.Error())
				}
			}()
		}
	}

	apps[n.Name] = &m
}

func Run(appName string) {
	n := findApp(appName)
	if n == nil {
		fwlogger.E("-- run app failed! not found app by name %s --", appName)
		return
	}
	runApp(n)
}

func Stop(appName string) {
	n := findApp(appName)
	if n == nil {
		fwlogger.E("-- stop app failed! not found app by name %s --", appName)
		return
	}
	if m, ok := apps[appName]; ok {
		m.destroy()
	}
}

func RunAll() {
	for _, n := range config.Apps {
		runApp(n)
	}
}

func StopAll() {
	for _, m := range apps {
		if m != nil {
			m.destroy()
		}
	}
}
