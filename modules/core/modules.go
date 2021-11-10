package core

import (
	"fmt"
	"github.com/dop251/goja"
	coreFilesystem "gjs/modules/fs"
	corePath "gjs/modules/path"
	coreUrl "gjs/modules/url"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const ModuleLocationNative = 1
const ModuleLocationPackage = 2
const ModuleLocationRelative = 3

type Module struct {
	Runtime *goja.Runtime
}

func moduleExists(name string) (exist, native bool, moduleLocation uint8) {
	switch name {
	case "fs", "url", "http", "https", "path":
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

func (m *Module) importModule(file string) goja.Value {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(fmt.Sprintf("error reading file %s", file))
		return goja.Undefined()
	}

	script, err := goja.Compile("", string(content), false)
	if err != nil {
		panic(err.Error())
		return goja.Undefined()
	}

	v, err := m.Runtime.RunProgram(script)
	if err != nil {
		panic(err.Error())
		return goja.Undefined()
	}

	return v
}

func (m *Module) Require(call goja.FunctionCall) goja.Value {
	moduleValue := call.Argument(0)
	if moduleValue.ExportType().Name() != "string" {
		panic(m.Runtime.NewTypeError("module must be a string"))
		return goja.Undefined()
	}

	moduleName := moduleValue.String()

	exist, native, location := moduleExists(moduleName)
	if !exist {
		m.Runtime.Interrupt(fmt.Sprintf("module '%s' is not exists", moduleName))
		return goja.Undefined()
	}

	var o *goja.Object

	if native {
		switch moduleName {
		case "fs":
			o = coreFilesystem.CreateModule(m.Runtime)
		case "path":
			o = corePath.CreateModule(m.Runtime)
		case "url":
			o = coreUrl.CreateModule(m.Runtime)
		}
	} else {
		switch location {
		case ModuleLocationPackage:
			return m.importModule(moduleName)
		case ModuleLocationRelative:
			return m.importModule(moduleName)
		}
	}

	return o
}
