package net

import "sync"

// ID生成
type Identity int

const (
	MINMUM_IDENTIFY = 1000000
	MAXMUM_IDENTIFY = 9999999
)

var (
	_identify_ Identity = MINMUM_IDENTIFY
	l          *sync.Mutex
)

func init() {
	l = &sync.Mutex{}
}

func genIdentity() Identity {
	l.Lock()
	_identify_++
	if _identify_ > MAXMUM_IDENTIFY {
		_identify_ = MINMUM_IDENTIFY
	}
	l.Unlock()
	return _identify_
}

func isValidIdentity(id Identity) bool {
	return id <= MAXMUM_IDENTIFY && id >= MINMUM_IDENTIFY
}
