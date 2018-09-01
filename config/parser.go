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
	Service map[string]*Service
	RpcPort int
	RpcTo   []string
	Logs    map[string]*log.Config
	Db      map[string]*DBConfig
	Config  map[string]interface{}
}

type GlobalConfig struct {
	EnableBingoLog bool
	BingoLogLevel  log.Level
}

type BingoConfig struct {
	Apps   []*AppConfig
	Config *GlobalConfig
}

type Parser interface {
	Parse(filePath string) *BingoConfig
}

func (c *BingoConfig) FindApp(name string) *AppConfig {
	for _, c := range c.Apps {
		if c.Name == name {
			return c
		}
	}
	return nil
}
