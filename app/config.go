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
	"github.com/snippetor/bingo/errors"
	"github.com/robertkrimen/otto"
	"github.com/snippetor/bingo/log"
)

type Service struct {
	Net   string
	Port  int
	Codec string
}

type DBConfig struct {
	Type     string
	Addr     string
	User     string
	Pwd      string
	Db       string
	TbPrefix string
}

type AppConfig struct {
	Name    string
	Package string
	Etds    []string
	Service []*Service
	RpcPort int
	RpcTo   []string
	Logs    map[string]*log.Config
	Db      []*DBConfig
	Config  map[string]interface{}
}

type BingoConfig struct {
	Apps []*AppConfig
}

var (
	bingoConfig *BingoConfig
)

func BingoConfiguration() *BingoConfig {
	return bingoConfig
}

// 解析配置文件
func parse(configPath string) {
	content, err := ioutil.ReadFile(configPath)
	errors.Check(err)
	o := otto.New()
	_, err = o.Run(content)
	errors.Check(err)

	bingoConfig = &BingoConfig{}
	apps, err := o.Get("apps")
	errors.Check(err)
	for _, name := range apps.Object().Keys() {
		app, err := apps.Object().Get(name)
		errors.Check(err)
		bingoConfig.Apps = append(bingoConfig.Apps, parseApp(name, app.Object()))
	}
}

func parseApp(name string, v *otto.Object) *AppConfig {
	app := &AppConfig{Name: name}
	// Package
	app.Package = parseString("package", v)
	// etds
	app.Etds = parseStrings("etds", v)
	// Service []*Service
	service := parseObjects("service", v)
	for _, v := range service {
		app.Service = append(app.Service, &Service{
			Net:   parseString("net", v),
			Port:  int(parseInt("net", v)),
			Codec: parseStringMust("codec", v, "protobuf"),
		})
	}
	// RpcPort int
	app.RpcPort = int(parseInt("rpcPort", v))
	// RpcTo   []string
	app.RpcTo = parseStrings("rpcPort", v)
	// Logs  map[string]*LogConfig
	m := parseStringMap("logs", v)
	app.Logs = make(map[string]*log.Config)
	for k, v := range m {
		app.Logs[k] = &log.Config{
			Level:              log.Level(mustInt(v, int64(log.Info))),
			OutputType:         log.OutputType(mustInt(v, int64(log.Console))),
			LogFileRollingType: log.RollingType(mustInt(v, int64(log.RollingNone))),
			LogFileOutputDir:   mustString(v, "."),
			LogFileMaxSize:     mustInt(v, 1*log.GB),
		}
	}
	// Db      []*DBConfig
	dbs := parseObjects("db", v)
	for _, v := range dbs {
		app.Db = append(app.Db, &DBConfig{
			Type:     parseString("type", v),
			Addr:     parseString("addr", v),
			User:     parseStringMust("user", v, ""),
			Pwd:      parseStringMust("pwd", v, ""),
			Db:       parseStringMust("db", v, ""),
			TbPrefix: parseStringMust("tbPrefix", v, ""),
		})
	}
	// config
	m = parseStringMap("config", v)
	app.Config = make(map[string]interface{})
	for k, v := range m {
		i, err := v.Export()
		errors.Check(err)
		app.Config[k] = i
	}
	return app
}

func parseStringMap(tag string, v *otto.Object) map[string]otto.Value {
	t, err := v.Get(tag)
	errors.Check(err)
	if !t.IsObject() {
		panic("Parse config file failed, 'app." + tag + "' must be object map.")
	}
	m := make(map[string]otto.Value)
	for _, name := range t.Object().Keys() {
		obj, err := t.Object().Get(name)
		errors.Check(err)
		m[name] = obj
	}
	return m
}

func parseObjects(tag string, v *otto.Object) []*otto.Object {
	t, err := v.Get(tag)
	errors.Check(err)
	if !t.IsObject() {
		panic("Parse config file failed, 'app." + tag + "' must be object array.")
	}
	if i, err := t.Export(); err != nil {
		panic("Parse config file failed, 'app." + tag + "' must be object array.")
	} else {
		if array, ok := i.([]*otto.Object); !ok {
			panic("Parse config file failed, 'app." + tag + "' must be object array.")
		} else {
			return array
		}
	}
}

func parseStrings(tag string, v *otto.Object) []string {
	t, err := v.Get(tag)
	errors.Check(err)
	if !t.IsObject() {
		panic("Parse config file failed, 'app." + tag + "' must be string array.")
	}
	if i, err := t.Export(); err != nil {
		panic("Parse config file failed, 'app." + tag + "' must be string array.")
	} else {
		if array, ok := i.([]string); !ok {
			panic("Parse config file failed, 'app." + tag + "' must be string array.")
		} else {
			return array
		}
	}
}

func parseString(tag string, v *otto.Object) string {
	t, err := v.Get(tag)
	errors.Check(err)
	if !t.IsString() {
		panic("Parse config file failed, 'app." + tag + "' must be string.")
	}
	str, err := t.ToString()
	errors.Check(err)
	return str
}

func parseStringMust(tag string, v *otto.Object, defValue string) string {
	t, err := v.Get(tag)
	if err != nil {
		return defValue
	}
	return mustString(t, defValue)
}

func parseInt(tag string, v *otto.Object) int64 {
	t, err := v.Get(tag)
	errors.Check(err)
	if !t.IsNumber() {
		panic("Parse config file failed, 'app." + tag + "' must be int.")
	}
	i, err := t.ToInteger()
	errors.Check(err)
	return i
}

func parseIntMust(tag string, v *otto.Object, defValue int64) int64 {
	t, err := v.Get(tag)
	if err != nil {
		return defValue
	}
	return mustInt(t, defValue)
}

func mustString(v otto.Value, defValue string) string {
	if !v.IsString() {
		return defValue
	}
	i, err := v.ToString()
	if err != nil {
		return defValue
	}
	return i
}

func mustInt(v otto.Value, defValue int64) int64 {
	if !v.IsNumber() {
		return defValue
	}
	i, err := v.ToInteger()
	if err != nil {
		return defValue
	}
	return i
}

func findApp(name string) *AppConfig {
	for _, n := range bingoConfig.Apps {
		if n.Name == name {
			return n
		}
	}
	return nil
}
