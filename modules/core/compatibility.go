package core

import "github.com/dop251/goja"

func RegisterCompatibility(vm *goja.Runtime) {
	_ = vm.GlobalObject().Set("global", vm.GlobalObject())

	exports := vm.NewObject()
	_ = vm.GlobalObject().Set("exports", exports)

	module := vm.NewObject()
	_ = module.Set("exports", exports)
	_ = vm.GlobalObject().Set("module", module)
}
