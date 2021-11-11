package url

import (
	"fmt"
	"github.com/dop251/goja"
	"net/url"
	"reflect"
	"strconv"
)

type Module struct {
	runtime *goja.Runtime
}

func (u *Module) parse(call goja.FunctionCall) goja.Value {
	urlValue := call.Argument(0)

	if urlValue.ExportType().Kind() != reflect.String {
		panic(u.runtime.NewTypeError("url must be a string"))
		return goja.Undefined()
	}

	urlStr := urlValue.String()
	nativeUrl, _ := url.Parse(urlStr)

	var port = 0

	portStr := nativeUrl.Port()
	if portStr != "" {
		port, _ = strconv.Atoi(portStr)
	}

	object := u.runtime.NewObject()

	_ = object.DefineAccessorProperty("host", u.runtime.ToValue(func() goja.Value {
		return u.runtime.ToValue(nativeUrl.Host)
	}), u.runtime.ToValue(func(call goja.FunctionCall) goja.Value {
		val := call.Argument(0)
		nativeUrl.Host = val.String()
		return goja.Undefined()
	}), goja.FLAG_FALSE, goja.FLAG_TRUE)

	_ = object.DefineAccessorProperty("protocol", u.runtime.ToValue(func() goja.Value {
		return u.runtime.ToValue(nativeUrl.Scheme)
	}), u.runtime.ToValue(func(call goja.FunctionCall) goja.Value {
		val := call.Argument(0)
		nativeUrl.Scheme = val.String()
		return goja.Undefined()
	}), goja.FLAG_FALSE, goja.FLAG_TRUE)

	_ = object.DefineAccessorProperty("port", u.runtime.ToValue(func() goja.Value {
		return u.runtime.ToValue(port)
	}), u.runtime.ToValue(func(call goja.FunctionCall) goja.Value {
		val := call.Argument(0)
		port = int(val.ToInteger())
		return goja.Undefined()
	}), goja.FLAG_FALSE, goja.FLAG_TRUE)

	_ = object.DefineAccessorProperty("search", u.runtime.ToValue(func() goja.Value {
		if nativeUrl.RawQuery != "" {
			return u.runtime.ToValue(fmt.Sprintf("?%s", nativeUrl.RawQuery))
		} else {
			return goja.Null()
		}
	}), u.runtime.ToValue(func(call goja.FunctionCall) goja.Value {
		val := call.Argument(0)
		nativeUrl.Path = val.String()
		return goja.Undefined()
	}), goja.FLAG_FALSE, goja.FLAG_TRUE)

	_ = object.DefineAccessorProperty("query", u.runtime.ToValue(func() goja.Value {
		if nativeUrl.RawQuery != "" {
			return u.runtime.ToValue(nativeUrl.RawQuery)
		} else {
			return goja.Null()
		}
	}), u.runtime.ToValue(func(call goja.FunctionCall) goja.Value {
		val := call.Argument(0)
		nativeUrl.Path = val.String()
		return goja.Undefined()
	}), goja.FLAG_FALSE, goja.FLAG_TRUE)

	_ = object.DefineAccessorProperty("path", u.runtime.ToValue(func() goja.Value {
		return u.runtime.ToValue(fmt.Sprintf("%s?%s", nativeUrl.Path, nativeUrl.RawQuery))
	}), u.runtime.ToValue(func(call goja.FunctionCall) goja.Value {
		val := call.Argument(0)
		nativeUrl.Path = val.String()
		return goja.Undefined()
	}), goja.FLAG_FALSE, goja.FLAG_TRUE)

	_ = object.DefineAccessorProperty("pathname", u.runtime.ToValue(func() goja.Value {
		return u.runtime.ToValue(nativeUrl.Path)
	}), u.runtime.ToValue(func(call goja.FunctionCall) goja.Value {
		val := call.Argument(0)
		nativeUrl.Path = val.String()
		return goja.Undefined()
	}), goja.FLAG_FALSE, goja.FLAG_TRUE)

	_ = object.DefineAccessorProperty("href", u.runtime.ToValue(func() goja.Value {
		if port == 0 {
			return u.runtime.ToValue(fmt.Sprintf("%s://%s%s", nativeUrl.Scheme, nativeUrl.Host, nativeUrl.Path))
		} else {
			return u.runtime.ToValue(fmt.Sprintf("%s://%s:%d%s", nativeUrl.Scheme, nativeUrl.Host, port, nativeUrl.Path))
		}
	}), u.runtime.ToValue(func(call goja.FunctionCall) goja.Value {
		val := call.Argument(0)
		port = int(val.ToInteger())
		return goja.Undefined()
	}), goja.FLAG_FALSE, goja.FLAG_TRUE)

	return object
}

func CreateModule(vm *goja.Runtime) *goja.Object {
	u := &Module{runtime: vm}

	object := vm.NewObject()
	_ = object.Set("parse", u.parse)
	return object
}
