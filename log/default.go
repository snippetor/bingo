package log

var DefaultLogger Logger

func init() {
	DefaultLogger = NewLogger(&Config{
		Level:      Info,
		OutputType: Console,
	})
}

func I(format string, v ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.I(format, v)
	}
}

func D(format string, v ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.D(format, v)
	}
}

func W(format string, v ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.W(format, v)
	}
}

func E(format string, v ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.E(format, v)
	}
}
