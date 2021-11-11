package globals

import (
	"github.com/dop251/goja"
	"reflect"
)

type BufferModule struct {
	runtime *goja.Runtime
}

func (b *BufferModule) bufferFrom(call goja.FunctionCall) goja.Value {
	d := call.Argument(0)

	if d.ExportType().Kind() == reflect.String {
		s := d.String()
		v := b.runtime.ToValue([]byte(s))
		return v
	} else {
		panic(b.runtime.NewTypeError("invalid argument type"))
	}

	return goja.Undefined()
}

func (b *BufferModule) isBuffer(call goja.FunctionCall) goja.Value {
	d := call.Argument(0)

	if arr, ok := d.(goja.DynamicArray); ok {
		for i := 0; i <= arr.Len(); i++ {
			arrItem := arr.Get(i)
			if arrItem.ExportType().Kind() != reflect.Int64 {
				return b.runtime.ToValue(false)
			}
		}

		return b.runtime.ToValue(true)
	} else {
		return b.runtime.ToValue(false)
	}
}

func RegisterBuffer(vm *goja.Runtime) {
	b := &BufferModule{runtime: vm}

	f := vm.ToValue(b.bufferFrom).(*goja.Object)
	_ = f.Set("from", b.bufferFrom)
	_ = f.Set("isBuffer", b.isBuffer)

	_ = vm.Set("Buffer", f)
}
