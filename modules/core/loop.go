package core

import "github.com/robertkrimen/otto"

func Loop(vm *otto.Otto) {
	for {
		select {
		case timer := <-ready:
			var arguments []interface{}
			if len(timer.call.ArgumentList) > 2 {
				tmp := timer.call.ArgumentList[2:]
				arguments = make([]interface{}, 2+len(tmp))
				for i, value := range tmp {
					arguments[i+2] = value
				}
			} else {
				arguments = make([]interface{}, 1)
			}

			arguments[0] = timer.call.ArgumentList[0]
			_, err := vm.Call(`Function.call.call`, nil, arguments...)

			if err != nil {
				for _, timer := range registry {
					timer.timer.Stop()
					delete(registry, timer)
					return
				}
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
