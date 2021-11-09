package path

import "github.com/robertkrimen/otto"

func join(_ otto.FunctionCall) otto.Value {
	return otto.UndefinedValue()
}

func CreateModule(vm *otto.Otto) otto.Value {
	object, err := vm.Object("new Object")
	if err != nil {
		panic("cant create object")
		return otto.UndefinedValue()
	}

	_ = object.Set("join", join)

	ret, _ := otto.ToValue(object)
	return ret
}
