package core

import (
	"github.com/dop251/goja"
)

func Loop(vm *goja.Runtime) {
	for {
		select {
		case timer := <-ready:
			var arguments []interface{}
			if len(timer.call.Arguments) > 2 {
				tmp := timer.call.Arguments[2:]
				arguments = make([]interface{}, 2+len(tmp))
				for i, value := range tmp {
					arguments[i+2] = value
				}
			} else {
				arguments = make([]interface{}, 1)
			}

			arguments[0] = timer.call.Arguments[0]

			var fn func(...interface{}) string
			_ = vm.ExportTo(vm.Get("Function.call.call"), &fn)
			fn(arguments...)

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
