package log


var (
	defaultLogger *Logger
)

func init() {
	defaultLogger = &NewLogger()
}

func SetDefaultLogConfig()  {
	
}
