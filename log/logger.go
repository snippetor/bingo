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
	"strings"
	"path/filepath"
	"strconv"
	"github.com/Unknwon/goconfig"
	"io/ioutil"
	"bytes"
	"regexp"
)

type Config struct {
	Level                  Level
	OutputType             OutputType
	LogFileRollingType     RollingType
	LogFileOutputDir       string
	LogFileName            string
	LogFileNameDatePattern string
	LogFileNameExt         string
	LogFileMaxSize         int64 // 字节
	LogFileScanInterval    int   // 秒
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

var DEFAULT_CONFIG = &Config{
	Level:                  Info,
	OutputType:             Console | File,
	LogFileRollingType:     RollingDaily,
	LogFileOutputDir:       ".",
	LogFileName:            "bingo",
	LogFileNameDatePattern: "20060102",
	LogFileNameExt:         ".log",
	LogFileMaxSize:         500 * MB,
	LogFileScanInterval:    600,
}

type Logger struct {
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
}

type OutputLog struct {
	level   Level
	content string
}

func NewLogger(configFile string) *Logger {
	// 默认配置
	l := &Logger{}
	l.setConfigFile(configFile)
	l.init()
	return l
}

func NewLoggerWithConfig(config *Config) *Logger {
	// 默认配置
	l := &Logger{}
	l.setConfig(config)
	l.init()
	return l
}

func (l *Logger) SetLevel(level Level) {
	l.config.Level = level
}

func (l *Logger) init() {
	l.c = make(chan *OutputLog, 5000)
	// log write
	go func() {
		for {
			s := <-l.c
			if l.config.OutputType&Console == Console {
				if s.level == Info {
					fmt.Println("\x1B[0;32m" + time.Now().Format("15:04:05.9999999") + " " + s.content + "\x1B[0m")
				} else if s.level == Debug {
					fmt.Println("\x1B[0;34m" + time.Now().Format("15:04:05.9999999") + " " + s.content + "\x1B[0m")
				} else if s.level == Warning {
					fmt.Println("\x1B[0;33m" + time.Now().Format("15:04:05.9999999") + " " + s.content + "\x1B[0m")
				} else if s.level == Error {
					fmt.Println("\x1B[0;31m" + time.Now().Format("15:04:05.9999999") + " " + s.content + "\x1B[0m")
				}
			}
			if l.config.OutputType&File == File {
				if l.f == nil || l.lg == nil {
					l.makeFile()
				}
				l.lg.Output(2, s.content)
			}
		}
	}()
}

func (l *Logger) setConfigFile(configFile string) {
	ini, err := goconfig.LoadConfigFile(configFile)
	if err != nil {
		log.Println("=========== parse config file failed!!! ==========", err)
		return
	}
	mode := ini.MustValue("", "workMode", "prod")
	if _, err := ini.GetSection(mode); err != nil {
		log.Println("=========== no section ["+mode+"] found in config file!!! ==========", err)
		return
	}
	c := &Config{}
	c.Level = Level(ini.MustInt(mode, "level", int(DEFAULT_CONFIG.Level)))
	c.OutputType = OutputType(ini.MustInt(mode, "outputType", int(DEFAULT_CONFIG.OutputType)))
	c.LogFileOutputDir = strings.TrimSpace(ini.MustValue(mode, "logFileOutputDir", DEFAULT_CONFIG.LogFileOutputDir))
	c.LogFileRollingType = RollingType(ini.MustInt(mode, "logFileRollingType", int(DEFAULT_CONFIG.LogFileRollingType)))
	c.LogFileName = strings.TrimSpace(ini.MustValue(mode, "logFileName", DEFAULT_CONFIG.LogFileName))
	c.LogFileNameDatePattern = strings.TrimSpace(ini.MustValue(mode, "logFileNameDatePattern", DEFAULT_CONFIG.LogFileNameDatePattern))
	c.LogFileNameExt = strings.TrimSpace(ini.MustValue(mode, "logFileNameExt", DEFAULT_CONFIG.LogFileNameExt))
	size := strings.TrimSpace(ini.MustValue(mode, "logFileMaxSize", "500MB"))
	i, err := strconv.ParseInt(size, 10, 64)
	if err == nil {
		c.LogFileMaxSize = i
	} else {
		i, err = strconv.ParseInt(size[:len(size)-2], 10, 64)
		if err == nil {
			unit := strings.ToUpper(size[len(size)-2:])
			if unit == "KB" {
				c.LogFileMaxSize = i * KB
			} else if unit == "MB" {
				c.LogFileMaxSize = i * MB
			} else if unit == "GB" {
				c.LogFileMaxSize = i * GB
			} else if unit == "TB" {
				c.LogFileMaxSize = i * TB
			}
		}
	}
	c.LogFileScanInterval = ini.MustInt(mode, "logFileScanInterval", 1)
	l.setConfig(c)
}

func (l *Logger) setConfig(c *Config) {
	l.config = c
	//l.makeFile()
	if c.OutputType&File == File {
		l.startFileCheckMonitor()
	}
}

func (l *Logger) SetPrefixes(prefix ...string) {
	l.prefixes = prefix
}

func (l *Logger) formatPrefixes() string {
	var buf bytes.Buffer
	if len(l.prefixes) > 0 {
		for _, p := range l.prefixes {
			buf.WriteString(p)
			buf.WriteString(" ")
		}
	}
	return buf.String()
}

func (l *Logger) I(format string, v ...interface{}) {
	if Info >= l.config.Level {
		if len(v) == 0 {
			l.c <- &OutputLog{Info, "[I] " + l.formatPrefixes() + format}
		} else {
			l.c <- &OutputLog{Info, "[I] " + l.formatPrefixes() + fmt.Sprintf(format, v...)}
		}
	}
}

func (l *Logger) D(format string, v ...interface{}) {
	if Debug >= l.config.Level {
		if len(v) == 0 {
			l.c <- &OutputLog{Debug, "[D] " + l.formatPrefixes() + format}
		} else {
			l.c <- &OutputLog{Debug, "[D] " + l.formatPrefixes() + fmt.Sprintf(format, v...)}
		}
	}
}

func (l *Logger) W(format string, v ...interface{}) {
	if Warning >= l.config.Level {
		if len(v) == 0 {
			l.c <- &OutputLog{Warning, "[W] " + l.formatPrefixes() + format}
		} else {
			l.c <- &OutputLog{Warning, "[W] " + l.formatPrefixes() + fmt.Sprintf(format, v...)}
		}
	}
}

func (l *Logger) E(format string, v ...interface{}) {
	if Error >= l.config.Level {
		if len(v) == 0 {
			l.c <- &OutputLog{Error, "[E] " + l.formatPrefixes() + format}
		} else {
			l.c <- &OutputLog{Error, "[E] " + l.formatPrefixes() + fmt.Sprintf(format, v...)}
		}
	}
}

func (l *Logger) startFileCheckMonitor() {
	if l.isMonitorRunning {
		return
	}
	l.isMonitorRunning = true
	// file check monitor
	go func() {
		monitorTimer := time.NewTicker(time.Duration(l.config.LogFileScanInterval) * time.Second)
		for {
			select {
			case <-monitorTimer.C:
				l.checkFile()
			}
		}
	}()
}

// 初始化日志文件
func (l *Logger) makeFile() {
	if l.config.OutputType&File != File {
		return
	}
	if l.f == nil {
		var err error
		var fileName = l.config.LogFileName
		if l.config.LogFileRollingType&RollingDaily == RollingDaily {
			t := time.Now().Format(l.config.LogFileNameDatePattern)
			fileName += "-" + t
		}
		if l.config.LogFileRollingType&RollingSize == RollingSize {
			fileName += "-" + l.getNextFileSeq(fileName)
		}
		os.Mkdir(l.config.LogFileOutputDir, os.ModePerm)
		l.f, err = os.OpenFile(filepath.Join(l.config.LogFileOutputDir, fileName+l.config.LogFileNameExt), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Println("=========== create log file failed!!! ========", err)
			return
		}
	}
	if l.lg == nil {
		if l.config.LogFileRollingType&RollingDaily == RollingDaily {
			l.lg = log.New(l.f, "", log.Lmicroseconds)
		} else {
			l.lg = log.New(l.f, "", log.Ldate|log.Lmicroseconds)
		}
	} else {
		if l.config.LogFileRollingType&RollingDaily == RollingDaily {
			l.lg.SetFlags(log.Lmicroseconds)
		} else {
			l.lg.SetFlags(log.Ldate | log.Lmicroseconds)
		}
		l.lg.SetOutput(l.f)
	}
}

// 检查文件是否需要重新创建
func (l *Logger) checkFile() {
	if l.config.OutputType&File != File || l.f == nil {
		return
	}
	needRecreate, newFileName := false, l.config.LogFileName
	if l.config.LogFileRollingType&RollingDaily == RollingDaily {
		dateString := time.Now().Format(l.config.LogFileNameDatePattern)
		t, _ := time.Parse(l.config.LogFileNameDatePattern, dateString)
		if len(l.f.Name()) >= len(l.config.LogFileName)+len(l.config.LogFileNameExt)+len(l.config.LogFileNameDatePattern)+1 {
			d, err := time.Parse(l.config.LogFileNameDatePattern, l.f.Name()[len(l.config.LogFileName)+1:len(l.config.LogFileName)+len(l.config.LogFileNameDatePattern)+1])
			if err != nil {
				log.Println("============== parse date failed!!! ===============")
			}
			if t.After(d) {
				needRecreate = true
				newFileName += "-" + dateString
				newFileName += "-1"
			}
		} else {
			needRecreate = true
			newFileName += "-" + dateString
			newFileName += "-1"
		}
	}

	if l.config.LogFileRollingType&RollingSize == RollingSize && !needRecreate {
		info, err := os.Stat(filepath.Join(l.config.LogFileOutputDir, l.f.Name()))
		if err != nil {
			log.Println("============= check file size failed!!! ==========", err)
			return
		}
		if info.Size() >= l.config.LogFileMaxSize {
			if needRecreate {
				newFileName += "-" + l.getNextFileSeq(newFileName)
			} else {
				needRecreate = true
				dateString := time.Now().Format(l.config.LogFileNameDatePattern)
				newFileName += "-" + dateString
				newFileName += "-" + l.getNextFileSeq(newFileName)
			}
		}
	}

	if needRecreate {
		l.f.Close()
		os.Mkdir(l.config.LogFileOutputDir, os.ModePerm)
		var err error
		l.f, err = os.OpenFile(filepath.Join(l.config.LogFileOutputDir, newFileName+l.config.LogFileNameExt), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Println("=========== open log file failed!!! ========", err)
			return
		}
		l.lg.SetOutput(l.f)
		if l.config.LogFileRollingType&RollingDaily == RollingDaily {
			l.lg.SetFlags(log.Lmicroseconds)
		} else {
			l.lg.SetFlags(log.Ldate | log.Lmicroseconds)
		}
	}
}

// 获取下一个文件序列号
func (l *Logger) getNextFileSeq(fileNamePrefix string) string {
	if files, err := ioutil.ReadDir(l.config.LogFileOutputDir); err != nil {
		return "1"
	} else {
		var maxFileSeq int64 = 1
		reg := regexp.MustCompile(fileNamePrefix + `-([0-9]+)\..*`)
		for _, info := range files {
			if !info.IsDir() {
				res := reg.FindStringSubmatch(info.Name())
				if len(res) >= 2 {
					if seq, err := strconv.ParseInt(res[1], 10, 64); err == nil && seq > maxFileSeq {
						maxFileSeq = seq
					}
				}
			}
		}
		return strconv.FormatInt(maxFileSeq+1, 10)
	}
}
