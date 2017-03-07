package bingo

import "github.com/snippetor/bingo/log"

var (
	// bingo框架日志
	fwLogger *log.Logger
)
func init() {
	fwLogger = log.NewLoggerWithConfig(log.DEFAULT_CONFIG)
}

func I(format string, v ...interface{}) {
	fwLogger.I(format, v...)
}

func D(format string, v ...interface{}) {
	fwLogger.D(format, v...)
}

func W(format string, v ...interface{}) {
	fwLogger.W(format, v...)
}

func E(format string, v ...interface{}) {
	fwLogger.E(format, v...)
}