package converters

import (
	"github.com/dop251/goja"
	"reflect"
	"strings"
)

func replaceFirst(str string, replacement byte) string {
	out := []byte(str)
	out[0] = replacement
	return string(out)
}

func InterfaceToObject(vm *goja.Runtime, v interface{}) *goja.Object {
	var r = vm.NewObject()

	t := reflect.Indirect(reflect.ValueOf(v))
	ty := reflect.TypeOf(v)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fn := ty.Field(i).Name
		fn = replaceFirst(fn, []byte(strings.ToLower(string(fn[0])))[0])

		switch f.Kind() {
		case reflect.Int64:
			_ = r.Set(fn, f.Int())
		case reflect.String:
			_ = r.Set(fn, f.String())
		case reflect.Bool:
			_ = r.Set(fn, f.Bool())
		case reflect.Func:
			if f.IsNil() {
				break
			}

			_ = r.Set(fn, f.Interface())
		}
	}

	return r
}
