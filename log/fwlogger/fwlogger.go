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

package fwlogger

import (
	"github.com/snippetor/bingo/log"
)

var (
	// bingo框架日志
	fwLogger log.Logger
)

func init() {
	fwLogger = log.NewLoggerWithConfig(&log.Config{
		Level:                  log.Warning,
		OutputType:             log.Console | log.File,
		LogFileRollingType:     log.RollingDaily,
		LogFileOutputDir:       ".",
		LogFileName:            "bingo",
		LogFileNameDatePattern: "20060102",
		LogFileNameExt:         ".log",
		LogFileMaxSize:         1 * log.GB,
		LogFileScanInterval:    60,
	})
}

func SetLevel(level log.Level) {
	fwLogger.SetLevel(level)
}

func I(format string, v ...interface{}) {
	fwLogger.I("[FW] "+format, v...)
}

func D(format string, v ...interface{}) {
	fwLogger.D("[FW] "+format, v...)
}

func W(format string, v ...interface{}) {
	fwLogger.W("[FW] "+format, v...)
}

func E(format string, v ...interface{}) {
	fwLogger.E("[FW] "+format, v...)
}
