package main

import (
	_ "embed"
	"fmt"
	"github.com/dop251/goja"
	"gjs/modules/core"
	"gjs/modules/core/globals"
	"os"
)

//go:embed example/index.js
var entry string

func main() {
	go func() {
		if caught := recover(); caught != nil {
			panic(caught)
			return
		}
	}()

	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

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

	_, err := vm.RunString(entry)
	if jse, ok := err.(*goja.Exception); ok {
		var e = jse.Value().Export()
		if e != nil {
			_, _ = os.Stderr.WriteString(e.(string))
		}
	} else {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("Error: %s", err.Error()))
	}

	core.Loop()
}
