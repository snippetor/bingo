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
