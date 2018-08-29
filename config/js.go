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

package config

import (
	"github.com/snippetor/bingo/errors"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"github.com/snippetor/bingo/log"
)

type JsParser struct {
}

// 解析配置文件
func (p *JsParser) Parse(configPath string) *BingoConfig {
	content, err := ioutil.ReadFile(configPath)
	errors.Check(err)
	o := otto.New()
	_, err = o.Run(content)
	errors.Check(err)

	bingoConfig := &BingoConfig{}
	apps, err := o.Get("apps")
	errors.Check(err)
	for _, name := range apps.Object().Keys() {
		app, err := apps.Object().Get(name)
		errors.Check(err)
		bingoConfig.Apps = append(bingoConfig.Apps, p.parseApp(name, app.Object()))
	}
	return bingoConfig
}

func (p *JsParser) parseApp(name string, v *otto.Object) *AppConfig {
	app := &AppConfig{Name: name}
	// Package
	app.Package = p.parseString("package", v)
	// etds
	app.Etds = p.parseStrings("etds", v)
	// Service []*Service
	service := p.parseObjects("service", v)
	for _, v := range service {
		app.Service = append(app.Service, &Service{
			Net:   p.parseString("net", v),
			Port:  int(p.parseInt("net", v)),
			Codec: p.parseStringMust("codec", v, "protobuf"),
		})
	}
	// RpcPort int
	app.RpcPort = int(p.parseInt("rpcPort", v))
	// RpcTo   []string
	app.RpcTo = p.parseStrings("rpcPort", v)
	// Logs  map[string]*LogConfig
	m := p.parseStringMap("logs", v)
	app.Logs = make(map[string]*log.Config)
	for k, v := range m {
		app.Logs[k] = &log.Config{
			Level:              log.Level(p.mustInt(v, int64(log.Info))),
			OutputType:         log.OutputType(p.mustInt(v, int64(log.Console))),
			LogFileRollingType: log.RollingType(p.mustInt(v, int64(log.RollingNone))),
			LogFileOutputDir:   p.mustString(v, "."),
			LogFileMaxSize:     p.mustInt(v, 1*log.GB),
		}
	}
	// Db      []*DBConfig
	dbs := p.parseObjects("db", v)
	for _, v := range dbs {
		app.Db = append(app.Db, &DBConfig{
			Type:     p.parseString("type", v),
			Addr:     p.parseString("addr", v),
			User:     p.parseStringMust("user", v, ""),
			Pwd:      p.parseStringMust("pwd", v, ""),
			Db:       p.parseStringMust("db", v, ""),
			TbPrefix: p.parseStringMust("tbPrefix", v, ""),
		})
	}
	// config
	m = p.parseStringMap("config", v)
	app.Config = make(map[string]interface{})
	for k, v := range m {
		i, err := v.Export()
		errors.Check(err)
		app.Config[k] = i
	}
	return app
}

func (p *JsParser) parseStringMap(tag string, v *otto.Object) map[string]otto.Value {
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

func (p *JsParser) parseObjects(tag string, v *otto.Object) []*otto.Object {
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

func (p *JsParser) parseStrings(tag string, v *otto.Object) []string {
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

func (p *JsParser) parseString(tag string, v *otto.Object) string {
	t, err := v.Get(tag)
	errors.Check(err)
	if !t.IsString() {
		panic("Parse config file failed, 'app." + tag + "' must be string.")
	}
	str, err := t.ToString()
	errors.Check(err)
	return str
}

func (p *JsParser) parseStringMust(tag string, v *otto.Object, defValue string) string {
	t, err := v.Get(tag)
	if err != nil {
		return defValue
	}
	return p.mustString(t, defValue)
}

func (p *JsParser) parseInt(tag string, v *otto.Object) int64 {
	t, err := v.Get(tag)
	errors.Check(err)
	if !t.IsNumber() {
		panic("Parse config file failed, 'app." + tag + "' must be int.")
	}
	i, err := t.ToInteger()
	errors.Check(err)
	return i
}

func (p *JsParser) parseIntMust(tag string, v *otto.Object, defValue int64) int64 {
	t, err := v.Get(tag)
	if err != nil {
		return defValue
	}
	return p.mustInt(t, defValue)
}

func (p *JsParser) mustString(v otto.Value, defValue string) string {
	if !v.IsString() {
		return defValue
	}
	i, err := v.ToString()
	if err != nil {
		return defValue
	}
	return i
}

func (p *JsParser) mustInt(v otto.Value, defValue int64) int64 {
	if !v.IsNumber() {
		return defValue
	}
	i, err := v.ToInteger()
	if err != nil {
		return defValue
	}
	return i
}
