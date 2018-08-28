package errors

import (
	"fmt"
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

type ErrCode uint32

const (
	ErrCodeNo       ErrCode = 0x0
	ErrCodeProtocol ErrCode = 0x1
	ErrCodeInternal ErrCode = 0x2
	ErrCodeConnect  ErrCode = 0x3
)

var errCodeName = map[ErrCode]string{
	ErrCodeNo:       "NO_ERROR",
	ErrCodeProtocol: "PROTOCOL_ERROR",
	ErrCodeInternal: "INTERNAL_ERROR",
	ErrCodeConnect:  "CONNECT_ERROR",
}

func (e ErrCode) String() string {
	if s, ok := errCodeName[e]; ok {
		return s
	}
	return fmt.Sprintf("unknown error code 0x%x", uint32(e))
}

type ConnectionError ErrCode

func (e ConnectionError) Error() string { return fmt.Sprintf("connection error: %s", ErrCode(e)) }

type UnknownNetTypeError struct {
	ErrCode
	UnknownType int
}

func (e UnknownNetTypeError) Error() string { return fmt.Sprintf("unknown net type error: %d", e.UnknownType) }
