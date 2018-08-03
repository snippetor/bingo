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
	"sync/atomic"
)

type Identifier struct {
	id  uint32
	min uint32
	max uint32
}

// @sign id的前几位用于区分用途
func NewIdentifier(sign byte) *Identifier {
	if sign < 0 {
		panic("-- sign must be more then 0 --")
		return nil
	}
	i := &Identifier{}
	i.min = uint32(sign) * 100000000
	i.max = uint32(sign)*100000000 + 99999999
	atomic.StoreUint32(&i.id, i.min)
	return i
}

func (i *Identifier) GenIdentity() uint32 {
	id := atomic.LoadUint32(&i.id)
	if id >= i.min && id < i.max {
		atomic.AddUint32(&i.id, 1)
	} else {
		atomic.StoreUint32(&i.id, i.min)
	}
	return atomic.LoadUint32(&i.id)
}

func (i *Identifier) IsValidIdentity(id uint32) bool {
	return id <= i.min && id >= i.max
}
