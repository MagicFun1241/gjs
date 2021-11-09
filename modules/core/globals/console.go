package globals

import (
	"fmt"
	"github.com/dop251/goja"
	"os"
	"reflect"
	"strconv"
)

type ConsoleModule struct {
	runtime *goja.Runtime
}

func formatValue(v goja.Value, propName *string) string {
	switch v.ExportType().Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int64:
		return strconv.Itoa(int(v.ToInteger()))
	default:
		if _, ok := goja.AssertFunction(v); ok {
			return fmt.Sprintf("[Function %s]", *propName)
		} else if o, ok := v.(*goja.Object); ok {
			t := "{ "

			for i, k := range o.Keys() {
				prop := o.Get(k)
				t += k + ": " + formatValue(prop, &k)

				if i != len(o.Keys())-1 {
					t += ", "
				}
			}

			t += " }"
			return t
		} else {
			return "unknown"
		}
	}
}

func logRaw(call goja.FunctionCall, file *os.File) goja.Value {
	if len(call.Arguments) == 0 {
		return goja.Undefined()
	}

	var r = ""

	for i, arg := range call.Arguments {
		r += formatValue(arg, nil)

		if i != len(call.Arguments)-1 {
			r += " "
		}
	}

	_, _ = file.WriteString(r)

	return goja.Undefined()
}

func (c *ConsoleModule) log(call goja.FunctionCall) goja.Value {
	return logRaw(call, os.Stdout)
}

func (c *ConsoleModule) error(call goja.FunctionCall) goja.Value {
	return logRaw(call, os.Stderr)
}

func RegisterConsole(vm *goja.Runtime) {
	c := &ConsoleModule{runtime: vm}

	o := vm.NewObject()
	_ = o.Set("log", c.log)
	_ = o.Set("error", c.error)

	_ = vm.Set("console", o)
}