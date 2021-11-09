package main

import (
	_ "embed"
	"fmt"
	"github.com/robertkrimen/otto"
	"gjs/modules/core"
)

//go:embed example/index.js
var entry string

func main() {
	defer func() {
		if caught := recover(); caught != nil {
			fmt.Print("Fatal error: ", caught)
			return
		}
	}()

	vm := otto.New()
	vm.Interrupt = make(chan func(), 1)

	_ = vm.Set("setTimeout", core.SetTimeout)
	_ = vm.Set("setInterval", core.SetInterval)
	_ = vm.Set("clearTimeout", core.ClearTimeout)
	_ = vm.Set("clearInterval", core.ClearTimeout)

	_ = vm.Set("require", core.Require)

	core.RegisterBuffer(vm)

	_, err := vm.Run(entry)
	if err != nil {
		fmt.Print("Runtime error: ", err)
		return
	}

	core.Loop(vm)
}
