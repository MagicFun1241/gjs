package core

import (
	"github.com/dop251/goja"
)

func Loop() {
	for {
		select {
		case timer := <-ready:
			var arguments []goja.Value
			if len(timer.call.Arguments) > 2 {
				tmp := timer.call.Arguments[2:]
				arguments = make([]goja.Value, 2+len(tmp))
				for i, value := range tmp {
					arguments[i+2] = value
				}
			} else {
				arguments = make([]goja.Value, 1)
			}

			arguments[0] = timer.call.Arguments[0]

			if fn, ok := goja.AssertFunction(arguments[0]); ok {
				_, err := fn(nil, arguments...)
				if err != nil {
					return
				}
			}

			for _, timer := range registry {
				timer.timer.Stop()
				delete(registry, timer)
				return
			}

			if timer.interval {
				timer.timer.Reset(timer.duration)
			} else {
				delete(registry, timer)
			}
		default:
			// Escape valve!
			// If this isn't here, we deadlock...
		}

		if len(registry) == 0 {
			break
		}
	}
}
