package native

import (
	"youpiteron.dev/white-monster-on-friday-night/internal/api"
	"youpiteron.dev/white-monster-on-friday-night/internal/compiler"
)

func Append(vm api.VM, args ...compiler.Value) (compiler.Value, error) {
	array := args[0].Array
	element := args[1]
	array = append(array, element)
	return compiler.Value{TypeOf: compiler.VAL_ARRAY, Array: array}, nil
}
