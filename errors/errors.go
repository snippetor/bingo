package errors

import (
	"github.com/snippetor/bingo/log"
	"runtime/debug"
)

func Check(err error) bool {
	if err != nil {
		panic(err)
		return false
	}
	return true
}

func CatchError(callback func(error)) bool {
	if err := recover(); err != nil {
		callback(err.(error))
		return false
	} else {
		return true
	}
}
