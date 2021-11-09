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

	if d.ExportType().Kind() != reflect.String {
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

	if d.ExportType().Kind() != reflect.Array {
		return b.runtime.ToValue(false)
	}

	// TODO: Проверять, является ли каждый элемент массива числом
	return b.runtime.ToValue(true)
}

func RegisterBuffer(vm *goja.Runtime) {
	b := &BufferModule{runtime: vm}

	o := vm.NewObject()
	_ = o.Set("from", b.bufferFrom)
	_ = o.Set("isBuffer", b.isBuffer)

	_ = vm.Set("Buffer", o)
}
