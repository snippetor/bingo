package net

import ()

var (
	globalPacker iMessagePacker
)

func init() {
	globalPacker = iMessagePacker(&defaultMessagePacker{})
}

func GetMessagePacker() iMessagePacker {
	return globalPacker
}
