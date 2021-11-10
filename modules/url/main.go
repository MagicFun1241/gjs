package url

import (
	"github.com/dop251/goja"
	"gjs/modules/core/converters"
	"net/url"
	"reflect"
	"strconv"
)

type Module struct {
	runtime *goja.Runtime
}

type URL struct {
	Host     string
	Port     uint16
	Href     string
	Protocol string
	Pathname string
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

	return converters.InterfaceToObject(u.runtime, URL{
		Host: nativeUrl.Host,
		Port: uint16(port),
		// TODO: Имплементировать нормальное поведение
		Href:     urlStr,
		Protocol: nativeUrl.Scheme,
		Pathname: nativeUrl.Path,
	})
}

func CreateModule(vm *goja.Runtime) *goja.Object {
	u := &Module{runtime: vm}

	object := vm.NewObject()
	_ = object.Set("parse", u.parse)
	return object
}
