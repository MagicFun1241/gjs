package main

import (
	_ "embed"
	"fmt"
	"github.com/MagicFun1241/gjs/modules/core"
	"github.com/MagicFun1241/gjs/modules/core/globals"
	"github.com/dop251/goja"
	"os"
)

//go:embed example/index.js
var entry string

func main() {
	vm := goja.New()

	defer func() {
		if err := recover(); err != nil {
			_, _ = os.Stderr.WriteString(err.(string))
		}
	}()

	m := &core.Module{Runtime: vm}

	_ = vm.Set("setTimeout", m.SetTimeout)
	_ = vm.Set("setInterval", m.SetInterval)
	_ = vm.Set("clearTimeout", m.ClearTimeout)
	_ = vm.Set("clearInterval", m.ClearTimeout)

	_ = vm.Set("require", m.Require)

	core.RegisterCompatibility(vm)

	globals.RegisterConsole(vm)
	globals.RegisterBuffer(vm)
	globals.RegisterProcess(vm)
	globals.RegisterUrl(vm)

	_, err := vm.RunString(entry)
	if err != nil {
		if jse, ok := err.(*goja.Exception); ok {
			_, _ = os.Stderr.WriteString(jse.String())
		} else {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("Error: %s", err.Error()))
		}
	}

	core.Loop()
}
