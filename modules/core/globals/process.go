package globals

import (
	"github.com/dop251/goja"
	"os"
)

type ProcessModule struct {
	runtime *goja.Runtime
}

func (p *ProcessModule) cwd() goja.Value {
	wd, err := os.Getwd()
	if err != nil {
		panic(p.runtime.NewGoError(err))
		return goja.Undefined()
	}

	return p.runtime.ToValue(wd)
}

func RegisterProcess(vm *goja.Runtime) {
	p := &ProcessModule{runtime: vm}

	o := vm.NewObject()
	_ = o.Set("cwd", p.cwd)

	_ = vm.GlobalObject().Set("process", o)
}
