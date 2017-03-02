package log

import (
	"testing"
	"time"
	"os"
)

func TestConsole(t *testing.T) {

	SetConfig(&Config{
		Level:  Info,
		OutType:Console,
	})

	I("test info....%s %s", "k", "x")
	D("test debug....")
	W("test warning....")
	E("test error....")
}

func TestFile(t *testing.T) {

	SetConfig(&Config{
		Level:              Info,
		OutType:            FileRollingDaily,
		OutDir:             "",
		LogFileName:        "test",
		LogFileMaxSize:     5 * MB,
		LogFileScanInterval:1,
	})

	I("test info....%s %s", "k", "x")
	I("test info....%s %s", "k", "x")
	I("test info....%s %s", "k", "x")
	D("test debug....")
	W("test warning....")
	E("test error....")

	time.Sleep(10*time.Second)
	os.Exit(1)
}

func BenchmarkFile(b *testing.B) {
	l := NewLogger()
	l.SetConfig(&Config{
		Level:              Info,
		OutType:            FileRollingDaily,
		OutDir:             "",
		LogFileName:        "test",
		LogFileMaxSize:     5 * MB,
		LogFileScanInterval:1,
	})
	for i := 0; i < b.N; i++ {
		l.I("test info....%s %s", "k", "x")
	}
}