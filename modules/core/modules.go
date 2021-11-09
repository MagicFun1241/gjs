package core

import (
	"fmt"
	"github.com/robertkrimen/otto"
	coreFs "gjs/modules/fs"
	corePath "gjs/modules/path"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const ModuleLocationNative = 1
const ModuleLocationPackage = 2
const ModuleLocationRelative = 3

func moduleExists(name string) (exist, native bool, moduleLocation uint8) {
	switch name {
	case "fs", "path":
		return true, true, ModuleLocationNative
	}

	if strings.HasPrefix(name, "./") || strings.HasPrefix(name, "../") {
		return true, false, ModuleLocationRelative
	} else {
		modulePath := path.Join("node_modules", name)

		_, err := os.Stat(modulePath)
		if os.IsNotExist(err) {
			return false, false, 0
		}

		_, err = os.Stat(path.Join(modulePath, "index.js"))
		if os.IsNotExist(err) {
			return false, false, 0
		}

		return true, false, ModuleLocationPackage
	}
}

func importModule(vm *otto.Otto, file string) otto.Value {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(fmt.Sprintf("error reading file %s", file))
		return otto.UndefinedValue()
	}

	script, err := vm.Compile("", string(content))
	if err != nil {
		panic(err.Error())
		return otto.UndefinedValue()
	}

	v, err := vm.Run(script)
	if err != nil {
		panic(err.Error())
		return otto.UndefinedValue()
	}

	return v
}

func Require(call otto.FunctionCall) otto.Value {
	if !call.Argument(0).IsString() {
		panic("module must be a string")
		return otto.UndefinedValue()
	}

	moduleName, _ := call.Argument(0).ToString()

	exist, native, location := moduleExists(moduleName)
	if !exist {
		panic("module is not exists")
		return otto.UndefinedValue()
	}

	var o otto.Value

	if native {
		switch moduleName {
		case "fs":
			o = coreFs.CreateModule(call.Otto)
		case "path":
			o = corePath.CreateModule(call.Otto)
		}
	} else {
		switch location {
		case ModuleLocationPackage:
			return importModule(call.Otto, moduleName)
		case ModuleLocationRelative:
			return importModule(call.Otto, moduleName)
		}
	}

	return o
}
