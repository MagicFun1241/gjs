package fs

import (
	"github.com/dop251/goja"
)

func CreateModule(vm *goja.Runtime) *goja.Object {
	fs := &Module{runtime: vm}

	object := vm.NewObject()
	_ = object.Set("readFileSync", fs.readFileSync)
	_ = object.Set("writeFileSync", fs.writeFileSync)
	_ = object.Set("existsSync", fs.existsSync)
	return object
}
