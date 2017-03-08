package log

var (
	defaultLogger *Logger
)

func SetDefaultLogConfigFile(configFile string) {
	if defaultLogger == nil {
		defaultLogger = &Logger{}
		defaultLogger.setConfigFile(configFile)
		defaultLogger.init()
	} else {
		defaultLogger.setConfigFile(configFile)
	}
}

func SetDefaultLogConfig(config *Config) {
	defaultLogger.setConfig(config)
	if defaultLogger == nil {
		defaultLogger = &Logger{}
		defaultLogger.setConfig(config)
		defaultLogger.init()
	} else {
		defaultLogger.setConfig(config)
	}
}

func initDefaultLogger() {
	if defaultLogger == nil {
		defaultLogger = &Logger{}
		defaultLogger.setConfig(DEFAULT_CONFIG)
		defaultLogger.init()
	}
}

func I(format string, v ...interface{}) {
	initDefaultLogger()
	defaultLogger.I(format, v...)
}

func D(format string, v ...interface{}) {
	initDefaultLogger()
	defaultLogger.D(format, v...)
}

func W(format string, v ...interface{}) {
	initDefaultLogger()
	defaultLogger.W(format, v...)
}

func E(format string, v ...interface{}) {
	initDefaultLogger()
	defaultLogger.E(format, v...)
}
