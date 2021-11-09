package path

import (
	"github.com/dop251/goja"
)

type Module struct {
	runtime *goja.Runtime
}

func (p *Module) join(_ goja.FunctionCall) goja.Value {
	return goja.Undefined()
}

func CreateModule(vm *goja.Runtime) *goja.Object {
	path := &Module{runtime: vm}

	object := vm.NewObject()
	_ = object.Set("join", path.join)
	return object
}
