package log

import (
	"log"
	"os"
	"time"
	"fmt"
	"strings"
)

func SetLevel(l Level) {
	level = l
}

func SetOutputType(t OutputType) {
	outType = t
}

func D(format string, v ...interface{}) {

}

func I(format string, v ...interface{}) {

}

func W(format string, v ...interface{}) {

}

func E(format string, v ...interface{}) {

}

type Level byte
type OutputType byte

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

const DATE_FORMAT = "2006-01-02"

// 输出类型
var (
	outType OutputType = Console
	// 输出等级
	level Level = Info
	// 内置logger
	lg *log.Logger = nil
	// 日志队列
	c chan string = make(chan string, 5000)
	// 当前日志文件
	f *os.File = nil
	// 日志文件名
	fn string = "bingo"
	// 日志目录
	dir string = ""
)

func init() {
	dir, _ = os.Getwd()
	go func() {
		for {
			s := <-c
			if outType == Console {
				writeToConsole(s)
			} else {
				writeToFile(s)
			}
		}
	}()
}

func writeToConsole(s string) {
	fmt.Println(s)
}

func writeToFile(s string) {
	if f == nil {

	}
}

func checkFile() {
	if outType == Console {
		return
	}
	if f == nil {
		var err error = nil
		if outType == FileRollingDaily {
			t := time.Now().Format(DATE_FORMAT)
			f, err = os.OpenFile(dir+"/"+fn+"_"+t, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		} else if outType == FileRollingSize {
			f, err = os.OpenFile(dir+"/"+fn+"_1", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
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
	if outType == FileRollingDaily {
		dateString := time.Now().Format(DATE_FORMAT)
		t, _ := time.Parse(DATE_FORMAT, dateString)
		d, _ := time.Parse(DATE_FORMAT, strings.Replace(f.Name(), fn+"_", "", 1))
		if t.After(d) {
			f.Close()
			var err error
			f, err = os.OpenFile(dir+"/"+fn+"_"+dateString, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				log.Println("=========== create log file failed!!! ========", err)
				return
			}
			lg.SetOutput(f)
		}
	} else {
		if maxFileCount > 1 {
			if fileSize(f.dir + "/" + f.filename) >= maxFileSize {
				return true
			}
		}
	}
}
