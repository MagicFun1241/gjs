package core

import (
	"github.com/robertkrimen/otto"
)

func bufferFrom(call otto.FunctionCall) otto.Value {
	d := call.Argument(0)

	if d.IsString() {
		s, _ := d.ToString()
		v, _ := call.Otto.ToValue([]byte(s))
		return v
	} else {
		panic("invalid argument type")
	}

	return otto.UndefinedValue()
}

func isBuffer(call otto.FunctionCall) otto.Value {
	d := call.Argument(0)

	if !d.IsObject() {
		return otto.FalseValue()
	}

	if d.Object().Class() != "GoArray" {
		return otto.FalseValue()
	}

	// TODO: Провертять, является ли каждый элемент массива числом
	return otto.TrueValue()
}

func RegisterBuffer(vm *otto.Otto) {
	o, _ := vm.Object("new Object")
	_ = o.Set("from", bufferFrom)
	_ = o.Set("isBuffer", isBuffer)

	_ = vm.Set("Buffer", o)
}
