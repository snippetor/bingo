package log

import (
)

var (
	defaultLogger *Logger
)

func init() {
	defaultLogger = &Logger{}
	defaultLogger.setConfig(DEFAULT_CONFIG)
	defaultLogger.init()
}

func SetDefaultLogConfigFile(configFile string) {
	defaultLogger.setConfigFile(configFile)
}

func SetDefaultLogConfig(config *Config) {
	defaultLogger.setConfig(config)
}

func I(format string, v ...interface{}) {
	defaultLogger.I(format, v...)
}

func D(format string, v ...interface{}) {
	defaultLogger.D(format, v...)
}

func W(format string, v ...interface{}) {
	defaultLogger.W(format, v...)
}

func E(format string, v ...interface{}) {
	defaultLogger.E(format, v...)
}
