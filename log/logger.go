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

package log

import (
	"log"
	"os"
	"time"
	"fmt"
	"path/filepath"
	"strconv"
	"bytes"
	"sync"
	"github.com/snippetor/bingo/utils"
)

type Config struct {
	Level              Level
	OutputType         OutputType
	LogFileRollingType RollingType
	LogFileOutputDir   string
	LogFileName        string
	LogFileMaxSize     int64 // 字节
}

type Level int
type OutputType int
type RollingType int

const (
	Info    Level = iota
	Debug
	Warning
	Error
)

const (
	Console OutputType = 1 << iota
	File
)

const (
	RollingNone              = 0
	RollingDaily RollingType = 1 << iota
	RollingSize
)

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
	TB
)

type Logger interface {
	SetLevel(level Level)
	SetPrefixes(prefix ...string)
	I(format string, v ...interface{})
	D(format string, v ...interface{})
	W(format string, v ...interface{})
	E(format string, v ...interface{})
	Close()
}

var _ Logger = (*logger)(nil)

type logger struct {
	config *Config
	// 内置logger
	lg *log.Logger
	// 日志队列
	c chan *OutputLog
	// 当前日志文件
	f *os.File
	// 检查文件monitor是否在运行
	isMonitorRunning bool
	// 日志前缀，将写在日期和等级后面，日志内容前面
	prefixes []string
	pool     *sync.Pool
	// 日志日期格式
	dateFormat       string
	scanFileInterval time.Duration
}

type OutputLog struct {
	level   Level
	content string
}

func NewLogger(config *Config) Logger {
	l := &logger{}
	l.setConfig(config)
	l.init()
	return l
}

func (l *logger) SetLevel(level Level) {
	l.config.Level = level
}

func (l *logger) init() {
	l.c = make(chan *OutputLog, 5000)
	l.pool = &sync.Pool{}
	l.pool.New = func() interface{} {
		return &OutputLog{}
	}
	l.dateFormat = "20060102"
	l.scanFileInterval = 5 * time.Minute
	// log write
	go func() {
		switch l.config.LogFileRollingType {
		case RollingNone:
			for {
				select {
				case s := <-l.c:
					l.printLog(s)
				}
			}
		case RollingDaily:
			now := time.Now()
			leftSecond := time.Duration(86400 - now.Hour()*3600 - now.Minute()*60 - now.Second())
			date := now.Format(l.dateFormat)
			dailyCheckChan := time.After(time.Millisecond*1000*leftSecond + 500)
			for {
				select {
				case s := <-l.c:
					l.printLog(s)
				case <-dailyCheckChan:
					l.rollingFile(date)
					leftSecond = 86400000
					now = time.Now()
					date = now.Format(l.dateFormat)
					dailyCheckChan = time.After(time.Millisecond*1000*leftSecond + 500)
				}
			}
		case RollingSize:
			sizeCheckChan := time.NewTicker(l.scanFileInterval).C
			for {
				select {
				case s := <-l.c:
					l.printLog(s)
				case <-sizeCheckChan:
					l.rollingFile("")
				}
			}
		case RollingSize | RollingDaily:
			now := time.Now()
			leftSecond := time.Duration(86400 - now.Hour()*3600 - now.Minute()*60 - now.Second())
			date := now.Format(l.dateFormat)
			dailyCheckChan := time.After(time.Millisecond*1000*leftSecond + 500)
			sizeCheckChan := time.NewTicker(l.scanFileInterval).C
			for {
				select {
				case s := <-l.c:
					l.printLog(s)
				case <-dailyCheckChan:
					l.rollingFile(date)
					leftSecond = 86400
					now = time.Now()
					date = now.Format(l.dateFormat)
					dailyCheckChan = time.After(time.Millisecond*1000*leftSecond + 500)
				case <-sizeCheckChan:
					l.rollingFile(date)
				}
			}
		}
	}()
}

func (l *logger) setConfig(c *Config) {
	l.config = c
	if c.OutputType&File == File {
		l.makeFile()
	}
}

func (l *logger) SetPrefixes(prefix ...string) {
	l.prefixes = prefix
}

func (l *logger) formatPrefixes() string {
	var buf bytes.Buffer
	if len(l.prefixes) > 0 {
		for _, p := range l.prefixes {
			buf.WriteString(p)
			buf.WriteString(" ")
		}
	}
	return buf.String()
}

func (l *logger) I(format string, v ...interface{}) {
	if Info >= l.config.Level {
		output := l.pool.Get().(*OutputLog)
		output.level = Info
		if len(v) == 0 {
			output.content = "[I] " + l.formatPrefixes() + format
		} else {
			output.content = "[I] " + l.formatPrefixes() + fmt.Sprintf(format, v...)
		}
		l.c <- output
	}
}

func (l *logger) D(format string, v ...interface{}) {
	if Debug >= l.config.Level {
		output := l.pool.Get().(*OutputLog)
		output.level = Debug
		if len(v) == 0 {
			output.content = "[D] " + l.formatPrefixes() + format
		} else {
			output.content = "[D] " + l.formatPrefixes() + fmt.Sprintf(format, v...)
		}
		l.c <- output
	}
}

func (l *logger) W(format string, v ...interface{}) {
	if Warning >= l.config.Level {
		output := l.pool.Get().(*OutputLog)
		output.level = Warning
		if len(v) == 0 {
			output.content = "[W] " + l.formatPrefixes() + format
		} else {
			output.content = "[W] " + l.formatPrefixes() + fmt.Sprintf(format, v...)
		}
		l.c <- output
	}
}

func (l *logger) E(format string, v ...interface{}) {
	if Error >= l.config.Level {
		output := l.pool.Get().(*OutputLog)
		output.level = Error
		if len(v) == 0 {
			output.content = "[E] " + l.formatPrefixes() + format
		} else {
			output.content = "[E] " + l.formatPrefixes() + fmt.Sprintf(format, v...)
		}
		l.c <- output
	}
}

func (l *logger) printLog(output *OutputLog) {
	if l.config.OutputType&Console == Console {
		if output.level == Info {
			fmt.Println("\x1B[0;32m" + time.Now().Format("15:04:05.9999999") + " " + output.content + "\x1B[0m")
		} else if output.level == Debug {
			fmt.Println("\x1B[0;34m" + time.Now().Format("15:04:05.9999999") + " " + output.content + "\x1B[0m")
		} else if output.level == Warning {
			fmt.Println("\x1B[0;33m" + time.Now().Format("15:04:05.9999999") + " " + output.content + "\x1B[0m")
		} else if output.level == Error {
			fmt.Println("\x1B[0;31m" + time.Now().Format("15:04:05.9999999") + " " + output.content + "\x1B[0m")
		}
	}
	if l.config.OutputType&File == File {
		if l.f == nil || l.lg == nil {
			l.makeFile()
		}
		l.lg.Output(2, output.content)
	}
	l.pool.Put(output)
}

// 初始化日志文件
func (l *logger) makeFile() {
	if l.config.OutputType&File != File {
		return
	}
	if l.f == nil {
		var err error
		os.Mkdir(l.config.LogFileOutputDir, os.ModePerm)
		l.f, err = os.OpenFile(filepath.Join(l.config.LogFileOutputDir, l.config.LogFileName), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
	}
	if l.lg == nil {
		l.lg = log.New(l.f, "", log.Ldate|log.Lmicroseconds)
	} else {
		l.lg.SetOutput(l.f)
	}
}

func (l *logger) rollingFile(date string) {
	if l.config.OutputType&File != File || l.f == nil {
		return
	}
	os.Mkdir(l.config.LogFileOutputDir, os.ModePerm)
	newFileName := l.config.LogFileName
	if l.config.LogFileRollingType&RollingDaily == RollingDaily {
		newFileName += "." + date
	}
	if l.config.LogFileRollingType&RollingSize == RollingSize {
		i := 1
		tmpFileName := newFileName + "." + strconv.Itoa(i)
		for utils.IsFileExists(filepath.Join(l.config.LogFileOutputDir, tmpFileName)) {
			tmpFileName = newFileName + "." + strconv.Itoa(i)
			i += 1
		}
		newFileName = tmpFileName
	}
	l.f.Close()
	if err := os.Rename(filepath.Join(l.config.LogFileOutputDir, l.config.LogFileName), filepath.Join(l.config.LogFileOutputDir, newFileName)); err != nil {
		panic(err)
	}
	newFile, err := os.OpenFile(filepath.Join(l.config.LogFileOutputDir, l.config.LogFileName), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	l.f = newFile
	l.lg.SetOutput(l.f)
}

func (l *logger) Close() {
	if l.f != nil {
		l.f.Close()
	}
}
