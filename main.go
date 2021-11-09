package main

import (
	_ "embed"
	"fmt"
	"github.com/dop251/goja"
	"gjs/modules/core"
	"gjs/modules/core/globals"
)

//go:embed example/index.js
var entry string

func main() {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	m := &core.Module{Runtime: vm}

	_ = vm.Set("setTimeout", m.SetTimeout)
	_ = vm.Set("setInterval", m.SetInterval)
	_ = vm.Set("clearTimeout", m.ClearTimeout)
	_ = vm.Set("clearInterval", m.ClearTimeout)

	_ = vm.Set("require", m.Require)

	globals.RegisterConsole(vm)
	globals.RegisterBuffer(vm)

	_, err := vm.RunString(entry)
	if err != nil {
		fmt.Print(err)
		return
	}

	core.Loop(vm)
}
