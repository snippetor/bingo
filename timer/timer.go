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

package timer

import "time"

type CallFunc func(v ...interface{})

type DelayCall struct {
	f    CallFunc
	args []interface{}
}

func (t *DelayCall) call() {
	t.f(t.args...)
}

type Timer struct {
	delay     time.Duration
	delayCall *DelayCall
}

func NewTimer(delay time.Duration, f CallFunc, args []interface{}) *Timer {
	return &Timer{
		delay: delay,
		delayCall: &DelayCall{
			f:    f,
			args: args,
		},
	}
}

func (t *Timer) Run() {
	go func() {
		time.Sleep(t.delay)
		t.delayCall.call()
	}()
}

func (t *Timer) GetDurations() time.Duration {
	return t.delay
}

func (t *Timer) GetFunc() *DelayCall {
	return t.delayCall
}
