package fs

import (
	"encoding/base64"
	"errors"
	"github.com/dop251/goja"
	ioFilesystem "io/fs"
	"io/ioutil"
	"os"
	"reflect"
)

type Module struct {
	runtime *goja.Runtime
}

var (
	writeFileError = errors.New("write error")
)

func validateEncoding(name string) bool {
	switch name {
	case "utf8", "utf-8", "base64":
		return true
	}

	return false
}

func dynamicArrayToBytes(a goja.DynamicArray) []byte {
	r := make([]byte, a.Len())

	for i := 0; i <= a.Len(); i++ {
		item := a.Get(i)
		r[i] = byte(item.ToInteger())
	}

	return r
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

func (fs *Module) mkdirSync(call goja.FunctionCall) goja.Value {
	p := call.Argument(0)
	optionsValue := call.Argument(1)

	var recursive = false
	var mode ioFilesystem.FileMode = 0777

	if !goja.IsUndefined(optionsValue) {
		if options, ok := optionsValue.(*goja.Object); ok {
			modeValue := options.Get("mode")
			if !goja.IsUndefined(modeValue) {
				if modeValue.ExportType().Kind() != reflect.Int64 {
					panic(fs.runtime.NewTypeError("mode option must be a int"))
					return goja.Undefined()
				}

				mode = ioFilesystem.FileMode(modeValue.ToInteger())
			}

			recursiveValue := options.Get("recursive")
			if !goja.IsUndefined(recursiveValue) {
				if recursiveValue.ExportType().Kind() != reflect.Bool {
					panic(fs.runtime.NewTypeError("recursive option must be a bool"))
					return goja.Undefined()
				}

				recursive = recursiveValue.ToBoolean()
			}
		} else {
			panic(fs.runtime.NewTypeError("options must be a object"))
			return goja.Undefined()
		}
	}

	if p.ExportType().Kind() != reflect.String {
		panic(fs.runtime.NewTypeError("dir path must be a string"))
		return goja.Undefined()
	}

	if recursive {
		_ = os.MkdirAll(p.String(), mode)
		// TODO: Реализовать рекурсивное создание подпапок тк по стандарту необходимо возвращать путь к первой созданной
		return goja.Undefined()
	} else {
		_ = os.Mkdir(p.String(), mode)
		return goja.Undefined()
	}
}

func (fs *Module) writeFileSync(call goja.FunctionCall) goja.Value {
	fileValue := call.Argument(0)
	dataValue := call.Argument(1)
	optionsValue := call.Argument(2)

	var encoding = "utf8"
	var mode ioFilesystem.FileMode = 0666

	if !goja.IsUndefined(optionsValue) {
		if options, ok := optionsValue.(*goja.Object); ok {
			encodingValue := options.Get("encoding")

			if encodingValue.ExportType().Kind() != reflect.String {
				panic(fs.runtime.NewTypeError("invalid encoding passed"))
				return goja.Undefined()
			}

			encoding = encodingValue.String()
			if !validateEncoding(encoding) {
				panic(fs.runtime.NewTypeError("invalid encoding value"))
				return goja.Undefined()
			}
		}
	}

	switch fileValue.ExportType().Kind() {
	case reflect.String:
		file := fileValue.String()

		var d []byte

		if dataArray, ok := dataValue.(goja.DynamicArray); ok {
			d = dynamicArrayToBytes(dataArray)
		} else {
			switch dataValue.ExportType().Kind() {
			case reflect.String:
				d = []byte(dataValue.String())
			}
		}

		switch encoding {
		case "base64":
			d = []byte(base64.StdEncoding.EncodeToString(d))
		}

		err := ioutil.WriteFile(file, d, mode)
		if err != nil {
			panic(fs.runtime.NewGoError(writeFileError))
		}
	}

	return goja.Undefined()
}

func (fs *Module) existsSync(call goja.FunctionCall) goja.Value {
	pathValue := call.Argument(0)

	if pathValue.ExportType().Kind() != reflect.String {
		panic(fs.runtime.NewTypeError("path must be a string"))
		return goja.Undefined()
	}

	path := pathValue.String()
	_, err := os.Stat(path)

	return fs.runtime.ToValue(err == nil)
}

func CreateModule(vm *goja.Runtime) *goja.Object {
	fs := &Module{runtime: vm}

	object := vm.NewObject()
	_ = object.Set("readFileSync", fs.readFileSync)
	_ = object.Set("writeFileSync", fs.writeFileSync)
	_ = object.Set("existsSync", fs.existsSync)
	_ = object.Set("mkdirSync", fs.mkdirSync)
	return object
}
