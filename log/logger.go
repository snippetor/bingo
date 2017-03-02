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

func SetConfigFile(inifile string) {
	ini, err := goconfig.LoadConfigFile(inifile)
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
	SetConfig(c)
}

func SetConfig(c *Config) {
	if config != nil && config.OutType == Console && c.OutType != Console {
		config = c
		makeFile()
		startFileCheckMonitor()
	} else {
		config = c
	}
}

func I(format string, v ...interface{}) {
	if Info >= config.Level {
		if len(v) == 0 {
			c <- "[I] " + format
		} else {
			c <- "[I] " + fmt.Sprintf(format, v...)
		}
	}
}

func D(format string, v ...interface{}) {
	if Debug >= config.Level {
		if len(v) == 0 {
			c <- "[D] " + format
		} else {
			c <- "[D] " + fmt.Sprintf(format, v...)
		}
	}
}

func W(format string, v ...interface{}) {
	if Warning >= config.Level {
		if len(v) == 0 {
			c <- "[W] " + format
		} else {
			c <- "[W] " + fmt.Sprintf(format, v...)
		}
	}
}

func E(format string, v ...interface{}) {
	if Error >= config.Level {
		if len(v) == 0 {
			c <- "[E] " + format
		} else {
			c <- "[E] " + fmt.Sprintf(format, v...)
		}
	}
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

// 输出类型
var (
	config *Config
	// 内置logger
	lg *log.Logger
	// 日志队列
	c chan string = make(chan string, 5000)
	// 当前日志文件
	f *os.File
	// 检查文件monitor是否在运行
	isMonitorRunning bool = false
)

func init() {
	// 默认配置
	d, _ := os.Getwd()
	config = &Config{
		Level:              Info,
		OutType:            Console,
		OutDir:             d,
		LogFileName:        "bingo",
		LogFileMaxSize:     500 * MB,
		LogFileScanInterval:1,
	}
	// log write
	go func() {
		for {
			s := <-c
			if config.OutType == Console {
				fmt.Println(s)
			} else {
				if f == nil || lg == nil {
					makeFile()
				}
				lg.Output(2, s)
			}
		}
	}()
}

func startFileCheckMonitor() {
	if isMonitorRunning {
		return
	}
	isMonitorRunning = true
	// file check monitor
	go func() {
		monitorTimer := time.NewTicker(config.LogFileScanInterval)
		for {
			select {
			case <-monitorTimer.C:
				checkFile()
			}
		}
	}()
}

func makeFile() {
	if config.OutType == Console {
		return
	}
	if f == nil {
		var err error = nil
		if config.OutType == FileRollingDaily {
			t := time.Now().Format(DATE_FORMAT)
			f, err = os.OpenFile(filepath.Join(config.OutDir, config.LogFileName+"_"+t), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		} else if config.OutType == FileRollingSize {
			f, err = os.OpenFile(filepath.Join(config.OutDir, config.LogFileName+"_1"), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		}
		if err != nil {
			log.Println("=========== create log file failed!!! ========", err)
			return
		}
	}
	if f == nil {
		log.Println("=========== check log file failed, not found log file!!! ========")
		return
	}
	if lg == nil {
		lg = log.New(f, "", log.Ldate|log.Lmicroseconds)
	} else {
		lg.SetOutput(f)
	}
}

func checkFile() {
	if config.OutType == Console {
		return
	}
	if config.OutType == FileRollingDaily {
		dateString := time.Now().Format(DATE_FORMAT)
		t, _ := time.Parse(DATE_FORMAT, dateString)
		d, _ := time.Parse(DATE_FORMAT, strings.Replace(f.Name(), config.LogFileName+"_", "", 1))
		if t.After(d) {
			f.Close()
			var err error
			f, err = os.OpenFile(filepath.Join(config.OutDir, config.LogFileName+"_"+dateString), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				log.Println("=========== create log file failed!!! ========", err)
				return
			}
			lg.SetOutput(f)
		}
	} else {
		f.Name()
		info, err := os.Stat(filepath.Join(config.OutDir, f.Name()))
		if err != nil {
			log.Println("============= check file size failed!!! ==========", err)
			return
		}
		if info.Size() >= config.LogFileMaxSize {
			seq, e := strconv.Atoi(strings.Replace(f.Name(), config.LogFileName+"_", "", 1))
			if e != nil {
				log.Println("============= check file sequence number failed!!! ==========", err)
				return
			}
			f.Close()
			var err error
			f, err = os.OpenFile(filepath.Join(config.OutDir, config.LogFileName+"_"+strconv.Itoa(seq + 1)), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				log.Println("=========== create log file failed!!! ========", err)
				return
			}
			lg.SetOutput(f)
		}
	}
}
