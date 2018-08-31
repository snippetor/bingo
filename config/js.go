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
	config, err := o.Get("config")
	errors.Check(err)
	if !config.IsObject() {
		panic("Parse config file failed, 'config' must be object.")
	}
	bingoConfig.Config = p.parseGlobalConfig(config.Object())
	apps, err := o.Get("apps")
	errors.Check(err)
	for _, name := range apps.Object().Keys() {
		app, err := apps.Object().Get(name)
		errors.Check(err)
		bingoConfig.Apps = append(bingoConfig.Apps, p.parseApp(name, app.Object()))
	}
	return bingoConfig
}

func (p *JsParser) parseGlobalConfig(v *otto.Object) *GlobalConfig {
	return &GlobalConfig{
		EnableBingoLog: p.parseBoolMust("enableBingoLog", v, false),
	}
}

func (p *JsParser) parseApp(name string, v *otto.Object) *AppConfig {
	app := &AppConfig{Name: name}
	// Package
	app.Package = p.parseString("package", v)
	// etds
	app.Etds = p.parseStrings("etcds", v)
	// Service map[string]*Service
	service := p.parseObjects("service", v)
	app.Service = make(map[string]*Service)
	for k, v := range service {
		if m, ok := v.(map[string]interface{}); ok {
			app.Service[k] = &Service{
				Net:   p.parseIString("service.net", m["net"]),
				Port:  int(p.parseIInt("service.port", m["port"])),
				Codec: p.parseIStringMust("service.codec", m["codec"], "protobuf"),
			}
		}
	}
	// RpcPort int
	app.RpcPort = int(p.parseInt("rpcPort", v))
	// RpcTo   []string
	app.RpcTo = p.parseStrings("rpcTo", v)
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
	// Db      map[string]*DBConfig
	dbs := p.parseObjects("db", v)
	app.Db = make(map[string]*DBConfig)
	for k, v := range dbs {
		if m, ok := v.(map[string]interface{}); ok {
			app.Db[k] = &DBConfig{
				Type:     p.parseIString("dbs.type", m["type"]),
				Addr:     p.parseIString("dbs.addr", m["addr"]),
				User:     p.parseIStringMust("dbs.user", m["user"], ""),
				Pwd:      p.parseIStringMust("dbs.pwd", m["pwd"], ""),
				Db:       p.parseIStringMust("dbs.db", m["db"], ""),
				TbPrefix: p.parseIStringMust("dbs.tbPrefix", m["tbPrefix"], ""),
			}
		}
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

func (p *JsParser) parseObjects(tag string, v *otto.Object) map[string]interface{} {
	t, err := v.Get(tag)
	errors.Check(err)
	if !t.IsObject() {
		panic("Parse config file failed, 'app." + tag + "' must be object array.")
	}
	if i, err := t.Export(); err != nil {
		panic("Parse config file failed, 'app." + tag + "' must be object array.")
	} else {
		if array, ok := i.(map[string]interface{}); !ok {
			if _, ok := i.([]interface{}); !ok {
				panic("Parse config file failed, 'app." + tag + "' must be object array.")
			} else {
				return map[string]interface{}{}
			}
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
			if a, ok := i.([]interface{}); !ok {
				panic("Parse config file failed, 'app." + tag + "' must be string array.")
			} else {
				var strArray []string
				for _, s := range a {
					if str, ok := s.(string); !ok {
						panic("Parse config file failed, 'app." + tag + "' must be string array.")
					} else {
						strArray = append(strArray, str)
					}
				}
				return strArray
			}
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

func (p *JsParser) parseIString(tag string, v interface{}) string {
	if str, ok := v.(string); ok {
		return str
	} else {
		panic("Parse config file failed, 'app." + tag + "' must be string.")
	}
}

func (p *JsParser) parseStringMust(tag string, v *otto.Object, defValue string) string {
	t, err := v.Get(tag)
	if err != nil {
		return defValue
	}
	return p.mustString(t, defValue)
}

func (p *JsParser) parseIStringMust(tag string, v interface{}, defValue string) string {
	if str, ok := v.(string); ok {
		return str
	} else {
		return defValue
	}
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

func (p *JsParser) parseIInt(tag string, v interface{}) int64 {
	if str, ok := v.(int64); ok {
		return str
	} else {
		panic("Parse config file failed, 'app." + tag + "' must be int.")
	}
}

func (p *JsParser) parseIntMust(tag string, v *otto.Object, defValue int64) int64 {
	t, err := v.Get(tag)
	if err != nil {
		return defValue
	}
	return p.mustInt(t, defValue)
}

func (p *JsParser) parseBoolMust(tag string, v *otto.Object, defValue bool) bool {
	t, err := v.Get(tag)
	if err != nil {
		return defValue
	}
	return p.mustBool(t, defValue)
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

func (p *JsParser) mustBool(v otto.Value, defValue bool) bool {
	if !v.IsBoolean() {
		return defValue
	}
	i, err := v.ToBoolean()
	if err != nil {
		return defValue
	}
	return i
}
