package fs

import (
	"encoding/base64"
	"github.com/dop251/goja"
	"io/ioutil"
	"reflect"
)

type Module struct {
	runtime *goja.Runtime
}

func validateEncoding(name string) bool {
	switch name {
	case "utf8", "utf-8", "base64":
		return true
	}

	return false
}

func (fs *Module) readFileSync(call goja.FunctionCall) goja.Value {
	path := call.Argument(0)
	if path.ExportType().Kind() != reflect.String {
		panic(fs.runtime.NewTypeError("path must be a string"))
		return goja.Undefined()
	}

	options := call.Argument(1)
	var encoding = "utf8"
	if !goja.IsUndefined(options) {
		if optionsObject, ok := options.(*goja.Object); ok {
			encodingValue := optionsObject.Get("encoding")
			if encodingValue.ExportType().Kind() != reflect.String {
				panic(fs.runtime.NewTypeError("invalid encoding value"))
				return goja.Undefined()
			}

			encoding = encodingValue.String()

			if !validateEncoding(encoding) {
				panic(fs.runtime.NewTypeError("invalid encoding value"))
				return goja.Undefined()
			}
		} else {
			panic(fs.runtime.NewTypeError("options must be an object"))
			return goja.Undefined()
		}
	}

	pathString := path.String()
	data, err := ioutil.ReadFile(pathString)
	if err != nil {
		panic(fs.runtime.NewTypeError("error reading file"))
		return goja.Undefined()
	}

	var ret goja.Value

	switch encoding {
	case "utf8", "utf-8":
		ret = fs.runtime.ToValue(string(data))
	case "base64":
		ret = fs.runtime.ToValue(base64.StdEncoding.EncodeToString(data))
	}

	return ret
}
