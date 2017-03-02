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
)

type Config struct {
	Level               Level
	OutType             OutputType
	OutDir              string
	LogFileName         string
	LogFileMaxSize      int64         // 字节
	LogFileScanInterval time.Duration // 秒
}

type Level int
type OutputType int
type SizeUnit int64

const (
	Info    Level = iota
	Debug
	Warning
	Error
)

const (
	Console          OutputType = iota
	FileRollingDaily
	FileRollingSize
)

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
	TB
)

const DATE_FORMAT = "2006-01-02"

type Logger struct {
	config *Config
	// 内置logger
	lg *log.Logger
	// 日志队列
	c chan string
	// 当前日志文件
	f *os.File
	// 检查文件monitor是否在运行
	isMonitorRunning bool
}

func NewLogger(configFile string) *Logger {
	// 默认配置
	l := &Logger{}
	l.c = make(chan string, 5000)
	l.SetConfigFile(configFile)
	// log write
	go func() {
		for {
			s := <-l.c
			if l.config.OutType == Console {
				fmt.Println(s)
			} else {
				if l.f == nil || l.lg == nil {
					l.makeFile()
				}
				l.lg.Output(2, s)
			}
		}
	}()
	return l
}

func (l *Logger) SetConfigFile(configFile string) {
	ini, err := goconfig.LoadConfigFile(configFile)
	if err != nil {
		log.Println("=========== parse config file failed!!! ==========", err)
		return
	}
	d, _ := os.Getwd()
	c := &Config{}
	c.Level = Level(ini.MustInt("", "level", int(Info)))
	c.OutType = OutputType(ini.MustInt("", "outputType", int(Console)))
	c.OutDir = strings.TrimSpace(ini.MustValue("", "outputDir", d))
	c.LogFileName = strings.TrimSpace(ini.MustValue("", "logFileName", "bingo"))
	size := strings.TrimSpace(ini.MustValue("", "logFileMaxSize", "500MB"))
	i, err := strconv.ParseInt(size, 10, 64)
	if err == nil {
		c.LogFileMaxSize = i
	} else {
		i, err = strconv.ParseInt(size[:len(size) - 2], 10, 64)
		if err == nil {
			unit := strings.ToUpper(size[len(size) - 2:])
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
	c.LogFileScanInterval = time.Duration(ini.MustInt("", "logFileScanInterval", 1)) * time.Second
	l.SetConfig(c)
}

func (l *Logger) SetConfig(c *Config) {
	if l.config != nil && l.config.OutType == Console && c.OutType != Console {
		l.config = c
		l.makeFile()
		l.startFileCheckMonitor()
	} else {
		l.config = c
	}
}

func (l *Logger) I(format string, v ...interface{}) {
	if Info >= l.config.Level {
		if len(v) == 0 {
			l.c <- "[I] " + format
		} else {
			l.c <- "[I] " + fmt.Sprintf(format, v...)
		}
	}
}

func (l *Logger) D(format string, v ...interface{}) {
	if Debug >= l.config.Level {
		if len(v) == 0 {
			l.c <- "[D] " + format
		} else {
			l.c <- "[D] " + fmt.Sprintf(format, v...)
		}
	}
}

func (l *Logger) W(format string, v ...interface{}) {
	if Warning >= l.config.Level {
		if len(v) == 0 {
			l.c <- "[W] " + format
		} else {
			l.c <- "[W] " + fmt.Sprintf(format, v...)
		}
	}
}

func (l *Logger) E(format string, v ...interface{}) {
	if Error >= l.config.Level {
		if len(v) == 0 {
			l.c <- "[E] " + format
		} else {
			l.c <- "[E] " + fmt.Sprintf(format, v...)
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
		monitorTimer := time.NewTicker(l.config.LogFileScanInterval)
		for {
			select {
			case <-monitorTimer.C:
				l.checkFile()
			}
		}
	}()
}

func (l *Logger) makeFile() {
	if l.config.OutType == Console {
		return
	}
	if l.f == nil {
		var err error = nil
		if l.config.OutType == FileRollingDaily {
			t := time.Now().Format(DATE_FORMAT)
			l.f, err = os.OpenFile(filepath.Join(l.config.OutDir, l.config.LogFileName+"_"+t), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		} else if l.config.OutType == FileRollingSize {
			l.f, err = os.OpenFile(filepath.Join(l.config.OutDir, l.config.LogFileName+"_1"), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		}
		if err != nil {
			log.Println("=========== create log file failed!!! ========", err)
			return
		}
	}
	if l.f == nil {
		log.Println("=========== check log file failed, not found log file!!! ========")
		return
	}
	if l.lg == nil {
		if l.config.OutType == FileRollingDaily {
			l.lg = log.New(l.f, "", log.Lmicroseconds)
		} else if l.config.OutType == FileRollingSize {
			l.lg = log.New(l.f, "", log.Ldate|log.Lmicroseconds)
		}
	} else {
		l.lg.SetOutput(l.f)
		if l.config.OutType == FileRollingDaily {
			l.lg.SetFlags(log.Lmicroseconds)
		} else if l.config.OutType == FileRollingSize {
			l.lg.SetFlags(log.Ldate | log.Lmicroseconds)
		}
	}
}

func (l *Logger) checkFile() {
	if l.config.OutType == Console {
		return
	}
	if l.config.OutType == FileRollingDaily {
		dateString := time.Now().Format(DATE_FORMAT)
		t, _ := time.Parse(DATE_FORMAT, dateString)
		d, _ := time.Parse(DATE_FORMAT, strings.Replace(l.f.Name(), l.config.LogFileName+"_", "", 1))
		if t.After(d) {
			l.f.Close()
			var err error
			l.f, err = os.OpenFile(filepath.Join(l.config.OutDir, l.config.LogFileName+"_"+dateString), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				log.Println("=========== create log file failed!!! ========", err)
				return
			}
			l.lg.SetOutput(l.f)
			if l.config.OutType == FileRollingDaily {
				l.lg.SetFlags(log.Lmicroseconds)
			} else if l.config.OutType == FileRollingSize {
				l.lg.SetFlags(log.Ldate | log.Lmicroseconds)
			}
		}
	} else {
		l.f.Name()
		info, err := os.Stat(filepath.Join(l.config.OutDir, l.f.Name()))
		if err != nil {
			log.Println("============= check file size failed!!! ==========", err)
			return
		}
		if info.Size() >= l.config.LogFileMaxSize {
			seq, e := strconv.Atoi(strings.Replace(l.f.Name(), l.config.LogFileName+"_", "", 1))
			if e != nil {
				log.Println("============= check file sequence number failed!!! ==========", err)
				return
			}
			l.f.Close()
			var err error
			l.f, err = os.OpenFile(filepath.Join(l.config.OutDir, l.config.LogFileName+"_"+strconv.Itoa(seq + 1)), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				log.Println("=========== create log file failed!!! ========", err)
				return
			}
			l.lg.SetOutput(l.f)
			if l.config.OutType == FileRollingDaily {
				l.lg.SetFlags(log.Lmicroseconds)
			} else if l.config.OutType == FileRollingSize {
				l.lg.SetFlags(log.Ldate | log.Lmicroseconds)
			}
		}
	}
}
