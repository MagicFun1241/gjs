package converters

import "github.com/dop251/goja"

func DynamicArrayToBytes(a goja.DynamicArray) []byte {
	r := make([]byte, a.Len())

	for i := 0; i <= a.Len(); i++ {
		item := a.Get(i)
		r[i] = byte(item.ToInteger())
	}

	return r
}
