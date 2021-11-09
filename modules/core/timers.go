package core

import (
	"github.com/robertkrimen/otto"
	"time"
)

type _timer struct {
	timer    *time.Timer
	duration time.Duration
	interval bool
	call     otto.FunctionCall
}

var registry = map[*_timer]*_timer{}
var ready = make(chan *_timer)

func newTimer(call otto.FunctionCall, interval bool) (*_timer, otto.Value) {
	delay, _ := call.Argument(1).ToInteger()
	if 0 >= delay {
		delay = 1
	}

	timer := &_timer{
		duration: time.Duration(delay) * time.Millisecond,
		call:     call,
		interval: interval,
	}
	registry[timer] = timer

	timer.timer = time.AfterFunc(timer.duration, func() {
		ready <- timer
	})

	value, err := call.Otto.ToValue(timer)
	if err != nil {
		panic(err)
	}

	return timer, value
}

func SetTimeout(call otto.FunctionCall) otto.Value {
	_, value := newTimer(call, false)
	return value
}

func SetInterval(call otto.FunctionCall) otto.Value {
	_, value := newTimer(call, true)
	return value
}

func ClearTimeout(call otto.FunctionCall) otto.Value {
	timer, _ := call.Argument(0).Export()
	if timer, ok := timer.(*_timer); ok {
		timer.timer.Stop()
		delete(registry, timer)
	}
	return otto.UndefinedValue()
}
