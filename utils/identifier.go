// Copyright 2017 bingo Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
