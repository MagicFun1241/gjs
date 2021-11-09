package fs

import "github.com/robertkrimen/otto"

func CreateModule(vm *otto.Otto) otto.Value {
	object, err := vm.Object("new Object")
	if err != nil {
		panic("cant create object")
		return otto.UndefinedValue()
	}

	_ = object.Set("readFileSync", ReadFileSync)

	ret, _ := otto.ToValue(object)
	return ret
}
