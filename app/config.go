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
	"github.com/snippetor/bingo/log/fwlogger"
	"gopkg.in/yaml.v2"
)

type Service struct {
	Name  string
	Type  string
	Port  int
	Codec string
}

type RPCConfig struct {
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

type DBConfig struct {
	Addr     string
	User     string
	Pwd      string
	Db       string
	TbPrefix string
}

type AppConfig struct {
	Name       string
	ModelName  string       `json:"app"`
	Domain     int
	Services   []*Service   `json:"service"`
	Rpc        *RPCConfig
	Config     map[string]interface{}
	LogConfigs []*LogConfig `json:"log"`
	DB         map[string]*DBConfig
}

type BingoConfig struct {
	Domains []string     `json:"domains"`
	Apps    []*AppConfig `json:"apps"`
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
	if err != nil {
		fwlogger.E("-- parse config file failed! %s --", err)
		return
	}
	if err = yaml.Unmarshal(content, bingoConfig); err != nil {
		fwlogger.E("-- parse config file failed! %s --", err)
		return
	}
}

func findApp(name string) *AppConfig {
	for _, n := range bingoConfig.Apps {
		if n.Name == name {
			return n
		}
	}
	return nil
}
