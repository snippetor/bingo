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
	"github.com/snippetor/bingo/utils"
	"io/ioutil"
)

type Config struct {
	Level                  Level
	OutputType             OutputType
	LogFileRollingType     RollingType
	LogFileOutputDir       string
	LogFileContentPattern  string // "DT [L] M"
	LogFileName            string
	LogFileNameDatePattern string
	LogFileExt             string
	LogFileMaxSize         int64         // 字节
	LogFileScanInterval    time.Duration // 秒
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
	LogFileContentPattern:  "T [L] M",
	LogFileName:            "default",
	LogFileNameDatePattern: "20060102",
	LogFileExt:             ".log",
	LogFileMaxSize:         500 * MB,
	LogFileScanInterval:    1 * time.Second,
}

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
	l.setConfigFile(configFile)
	l.init()
	return l
}

func (l *Logger) init() {
	l.c = make(chan string, 5000)
	// log write
	go func() {
		for {
			s := <-l.c
			if l.config.OutputType&Console == Console {
				fmt.Println(time.Now().Format("15:04:05") + " " + s)
			}
			if l.config.OutputType&File == File {
				if l.f == nil || l.lg == nil {
					l.makeFile()
				}
				l.lg.Output(2, s)
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
	d, _ := os.Getwd()
	mode := ini.MustValue("", "workMode", "prod")
	if _, err := ini.GetSection(mode); err != nil {
		log.Println("=========== no section ["+mode+"] found in config file!!! ==========", err)
		return
	}
	c := &Config{}
	c.Level = Level(ini.MustInt(mode, "level", int(Info)))
	c.OutputType = OutputType(ini.MustInt(mode, "outputType", int(Console)))
	c.LogFileOutputDir = strings.TrimSpace(ini.MustValue(mode, "logFileOutputDir", d))
	c.LogFileRollingType = RollingType(ini.MustInt(mode, "logFileRollingType", int(RollingDaily)))
	c.LogFileName = strings.TrimSpace(ini.MustValue(mode, "logFileName", "bingo"))
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
	c.LogFileScanInterval = time.Duration(ini.MustInt(mode, "logFileScanInterval", 1)) * time.Second
	l.setConfig(c)
}

func (l *Logger) setConfig(c *Config) {
	l.config = c
	l.makeFile()
	l.startFileCheckMonitor()
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

// 初始化日志文件
func (l *Logger) makeFile() {
	if l.config.OutputType == Console {
		return
	}
	if l.f == nil {
		var err error
		var fileName string = l.config.LogFileName
		if l.config.LogFileRollingType&RollingDaily == RollingDaily {
			t := time.Now().Format(l.config.LogFileNameDatePattern)
			fileName += "-" + t
		}
		if l.config.LogFileRollingType&RollingSize == RollingSize {
			fileName += "-" + l.genFileSeq()
		}
		l.f, err = os.OpenFile(filepath.Join(l.config.LogFileOutputDir, fileName), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
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
	if l.config.OutputType == Console {
		return
	}
	needRecreate, newFileName := false, l.config.LogFileName
	if l.config.LogFileRollingType&RollingDaily == RollingDaily {
		dateString := time.Now().Format(l.config.LogFileNameDatePattern)
		t, _ := time.Parse(l.config.LogFileNameDatePattern, dateString)
		if len(l.f.Name()) >= len(l.config.LogFileName)+9 {
			d, err := time.Parse(l.config.LogFileNameDatePattern, l.f.Name()[len(l.config.LogFileName)+1:len(l.config.LogFileName)+9])
			if err != nil {
				log.Println("============== parse date failed!!! ===============")
			}
			if t.After(d) {
				needRecreate = true
				newFileName += "-" + dateString
				newFileName += "-1"
				l.resetFileSeq()
			}
		} else {
			needRecreate = true
			newFileName += "-" + dateString
			newFileName += "-1"
			l.resetFileSeq()
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
				newFileName += "-" + l.genFileSeq()
			} else {
				needRecreate = true
				dateString := time.Now().Format(l.config.LogFileNameDatePattern)
				newFileName += "-" + dateString
				newFileName += "-" + l.genFileSeq()
			}
		}
	}

	if needRecreate {
		l.f.Close()
		var err error
		l.f, err = os.OpenFile(filepath.Join(l.config.LogFileOutputDir, newFileName), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Println("=========== create log file failed!!! ========", err)
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

// 生成日志文件序列号，并保存到.seq
func (l *Logger) genFileSeq() string {
	seqFile := filepath.Join(l.config.LogFileOutputDir, ".seq")
	if utils.IsFileExists(seqFile) {
		bytes, err := ioutil.ReadFile(seqFile)
		if err == nil {
			seq, err := strconv.Atoi(string(bytes))
			if err == nil {
				ioutil.WriteFile(seqFile, []byte(strconv.Itoa(seq+1)), 0666)
				return strconv.Itoa(seq + 1)
			}
		}
	}
	ioutil.WriteFile(seqFile, []byte("1"), 0666)
	return "1"
}

// 重置日志文件序列号
func (l *Logger) resetFileSeq() {
	seqFile := filepath.Join(l.config.LogFileOutputDir, ".seq")
	if utils.IsFileExists(seqFile) {
		ioutil.WriteFile(seqFile, []byte("1"), 0666)
	}
}
