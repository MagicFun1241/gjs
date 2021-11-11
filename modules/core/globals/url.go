package globals

import (
	"fmt"
	"github.com/dop251/goja"
	"net/url"
	"reflect"
	"strconv"
)

type UrlModule struct {
	runtime *goja.Runtime
}

func (u *UrlModule) urlConstructor(call goja.ConstructorCall) *goja.Object {
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

	_ = call.This.DefineAccessorProperty("hostname", u.runtime.ToValue(func() goja.Value {
		return u.runtime.ToValue(nativeUrl.Hostname())
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
		return u.runtime.ToValue(strconv.Itoa(port))
	}), u.runtime.ToValue(func(call goja.FunctionCall) goja.Value {
		val := call.Argument(0)
		port = int(val.ToInteger())
		return goja.Undefined()
	}), goja.FLAG_FALSE, goja.FLAG_TRUE)

	_ = call.This.DefineAccessorProperty("pathname", u.runtime.ToValue(func() goja.Value {
		return u.runtime.ToValue(nativeUrl.Path)
	}), u.runtime.ToValue(func(call goja.FunctionCall) goja.Value {
		val := call.Argument(0)
		nativeUrl.Path = val.String()
		return goja.Undefined()
	}), goja.FLAG_FALSE, goja.FLAG_TRUE)

	_ = call.This.DefineAccessorProperty("search", u.runtime.ToValue(func() goja.Value {
		return u.runtime.ToValue(fmt.Sprintf("?%s", nativeUrl.RawQuery))
	}), u.runtime.ToValue(func(call goja.FunctionCall) goja.Value {
		val := call.Argument(0)
		nativeUrl.Path = val.String()
		return goja.Undefined()
	}), goja.FLAG_FALSE, goja.FLAG_TRUE)

	_ = call.This.DefineAccessorProperty("origin", u.runtime.ToValue(func() goja.Value {
		if port == 0 {
			return u.runtime.ToValue(fmt.Sprintf("%s://%s", nativeUrl.Scheme, nativeUrl.Host))
		} else {
			return u.runtime.ToValue(fmt.Sprintf("%s://%s:%d", nativeUrl.Scheme, nativeUrl.Hostname(), port))
		}
	}), u.runtime.ToValue(func(call goja.FunctionCall) goja.Value {
		val := call.Argument(0)
		port = int(val.ToInteger())
		return goja.Undefined()
	}), goja.FLAG_FALSE, goja.FLAG_TRUE)

	_ = call.This.DefineAccessorProperty("href", u.runtime.ToValue(func() goja.Value {
		if port == 0 {
			return u.runtime.ToValue(fmt.Sprintf("%s://%s%s?%s", nativeUrl.Scheme, nativeUrl.Host, nativeUrl.Path, nativeUrl.RawQuery))
		} else {
			return u.runtime.ToValue(fmt.Sprintf("%s://%s:%d%s?%s", nativeUrl.Scheme, nativeUrl.Hostname(), port, nativeUrl.Path, nativeUrl.RawQuery))
		}
	}), u.runtime.ToValue(func(call goja.FunctionCall) goja.Value {
		val := call.Argument(0)
		port = int(val.ToInteger())
		return goja.Undefined()
	}), goja.FLAG_FALSE, goja.FLAG_TRUE)

	return nil
}

func (u *UrlModule) urlSearchParamsConstructor() *goja.Object {
	return nil
}

func RegisterUrl(vm *goja.Runtime) {
	u := &UrlModule{runtime: vm}

	_ = vm.GlobalObject().Set("URL", u.urlConstructor)
	_ = vm.GlobalObject().Set("URLSearchParams", u.urlSearchParamsConstructor)
}
