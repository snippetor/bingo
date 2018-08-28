package module

import "github.com/snippetor/bingo/log"

type Loggers map[string]log.Logger

// log module
type LogModule interface {
	Module
	DefaultLogger() log.Logger
	Logger(name string) log.Logger
}

type logModule struct {
	loggers Loggers
}

func NewLogModule(loggers Loggers) LogModule {
	return &logModule{loggers}
}

func (m *logModule) DefaultLogger() log.Logger {
	return m.Logger("default")
}

func (m *logModule) Logger(name string) log.Logger {
	if m.loggers != nil {
		if logger, ok := m.loggers[name]; ok {
			return logger
		}
	}
	logger := log.NewLogger(&log.Config{
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
