package fs

import (
	"encoding/base64"
	"github.com/robertkrimen/otto"
	"io/ioutil"
)

func validateEncoding(name string) bool {
	switch name {
	case "utf8", "utf-8", "base64":
		return true
	}

	return false
}

func ReadFileSync(call otto.FunctionCall) otto.Value {
	path := call.Argument(0)
	if !path.IsString() {
		panic("path must be a string")
		return otto.UndefinedValue()
	}

	options := call.Argument(1)
	var encoding = "utf8"
	if !options.IsUndefined() {
		if !options.IsObject() {
			panic("options must be an object")
			return otto.UndefinedValue()
		}

		oo := options.Object()
		encodingValue, _ := oo.Get("encoding")

		if !encodingValue.IsString() {
			panic("invalid encoding value")
			return otto.UndefinedValue()
		}

		t, _ := encodingValue.ToString()
		encoding = t

		if !validateEncoding(encoding) {
			panic("invalid encoding value")
			return otto.UndefinedValue()
		}
	}

	pathString, _ := path.ToString()
	data, err := ioutil.ReadFile(pathString)
	if err != nil {
		panic("error reading file")
		return otto.UndefinedValue()
	}

	var ret otto.Value

	switch encoding {
	case "utf8", "utf-8":
		ret, _ = otto.ToValue(string(data))
	case "base64":
		ret, _ = otto.ToValue(base64.StdEncoding.EncodeToString(data))
	}

	return ret
}
