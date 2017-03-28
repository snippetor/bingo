package utils

import (
	"sync"
)

// ID生成
type Identity int32

type Identifier struct {
	_identify_      Identity
	l               *sync.Mutex
	MINMUM_IDENTIFY Identity
	MAXMUM_IDENTIFY Identity
}

// @sign id的前几位用于区分用途
func NewIdentifier(sign byte) *Identifier {
	if sign < 0 {
		panic("-- sign must be more then 0 --")
		return nil
	}
	i := &Identifier{}
	i.l = &sync.Mutex{}
	i.MINMUM_IDENTIFY = Identity(int32(sign) * 10000000)
	i.MAXMUM_IDENTIFY = Identity(int32(sign)*10000000 + 9999999)
	return i
}

func (i *Identifier) GenIdentity() Identity {
	i.l.Lock()
	if i._identify_ >= i.MINMUM_IDENTIFY && i._identify_ < i.MAXMUM_IDENTIFY {
		i._identify_++
	} else {
		i._identify_ = i.MINMUM_IDENTIFY
	}
	i.l.Unlock()
	return i._identify_
}

func (i *Identifier) IsValidIdentity(id Identity) bool {
	return id <= i.MAXMUM_IDENTIFY && id >= i.MINMUM_IDENTIFY
}
