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
	object  *goja.Object
}

func (u *Module) urlConstructor(call goja.ConstructorCall) *goja.Object {
	urlValue := call.Argument(0)

	if urlValue.ExportType().Kind() != reflect.String {
		panic(u.runtime.NewTypeError("url must be a string"))
		return nil
	}

	urlStr := urlValue.String()
	nativeUrl, _ := url.Parse(urlStr)

	var port = 0

	portStr := nativeUrl.Port()
	if portStr != "" {
		port, _ = strconv.Atoi(portStr)
	}

	_ = call.This.DefineAccessorProperty("host", u.runtime.ToValue(func() goja.Value {
		return u.runtime.ToValue(nativeUrl.Host)
	}), u.runtime.ToValue(func(call goja.FunctionCall) goja.Value {
		val := call.Argument(0)
		nativeUrl.Host = val.String()
		return goja.Undefined()
	}), goja.FLAG_FALSE, goja.FLAG_TRUE)

	_ = call.This.DefineAccessorProperty("protocol", u.runtime.ToValue(func() goja.Value {
		return u.runtime.ToValue(nativeUrl.Scheme)
	}), u.runtime.ToValue(func(call goja.FunctionCall) goja.Value {
		val := call.Argument(0)
		nativeUrl.Scheme = val.String()
		return goja.Undefined()
	}), goja.FLAG_FALSE, goja.FLAG_TRUE)

	_ = call.This.DefineAccessorProperty("port", u.runtime.ToValue(func() goja.Value {
		return u.runtime.ToValue(port)
	}), u.runtime.ToValue(func(call goja.FunctionCall) goja.Value {
		val := call.Argument(0)
		port = int(val.ToInteger())
		return goja.Undefined()
	}), goja.FLAG_FALSE, goja.FLAG_TRUE)

	_ = call.This.DefineAccessorProperty("search", u.runtime.ToValue(func() goja.Value {
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

	_ = call.This.DefineAccessorProperty("query", u.runtime.ToValue(func() goja.Value {
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

	_ = call.This.DefineAccessorProperty("path", u.runtime.ToValue(func() goja.Value {
		return u.runtime.ToValue(fmt.Sprintf("%s?%s", nativeUrl.Path, nativeUrl.RawQuery))
	}), u.runtime.ToValue(func(call goja.FunctionCall) goja.Value {
		val := call.Argument(0)
		nativeUrl.Path = val.String()
		return goja.Undefined()
	}), goja.FLAG_FALSE, goja.FLAG_TRUE)

	_ = call.This.DefineAccessorProperty("pathname", u.runtime.ToValue(func() goja.Value {
		return u.runtime.ToValue(nativeUrl.Path)
	}), u.runtime.ToValue(func(call goja.FunctionCall) goja.Value {
		val := call.Argument(0)
		nativeUrl.Path = val.String()
		return goja.Undefined()
	}), goja.FLAG_FALSE, goja.FLAG_TRUE)

	_ = call.This.DefineAccessorProperty("href", u.runtime.ToValue(func() goja.Value {
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

	return nil
}

func (u *Module) parse(call goja.FunctionCall) goja.Value {
	if c, ok := goja.AssertFunction(u.object.Get("Url")); ok {
		v, _ := c(nil, call.Argument(0))
		return v
	}

	return goja.Undefined()
}

func CreateModule(vm *goja.Runtime) *goja.Object {
	u := &Module{runtime: vm}
	u.object = vm.NewObject()

	_ = u.object.Set("parse", u.parse)
	_ = u.object.Set("Url", u.urlConstructor)
	return u.object
}
