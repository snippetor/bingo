package module

import "github.com/snippetor/bingo/log"

// log module
type LogModule interface {
	Module
	GetLogger(name string) log.Logger
}

type Loggers map[string]log.Logger

type logModule struct {
	loggers Loggers
}

func NewLogModule(loggers Loggers) LogModule {
	return &logModule{loggers}
}

func (m *logModule) GetLogger(name string) log.Logger {
	if m.loggers != nil {
		if logger, ok := m.loggers[name]; ok {
			return logger
		}
	}
	logger := log.NewLoggerWithConfig(&log.Config{
		Level:      log.Info,
		OutputType: log.Console,
	})
	return logger
}

func (m *logModule) Close() {
	for _, logger := range m.loggers {
		if logger != nil {
			logger.Close()
		}
	}

}
