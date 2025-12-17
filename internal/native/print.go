package native

import (
	"fmt"

	"youpiteron.dev/white-monster-on-friday-night/internal/api"
	"youpiteron.dev/white-monster-on-friday-night/internal/compiler"
)

func Println(vm api.VM, args ...compiler.Value) (compiler.Value, error) {
	values := make([]any, 0, len(args))
	for _, arg := range args {
		switch arg.TypeOf {
		case compiler.VAL_INT:
			values = append(values, arg.Int)
		case compiler.VAL_BOOL:
			values = append(values, arg.Bool)
		case compiler.VAL_NULL:
			values = append(values, nil)
		case compiler.VAL_CLOSURE:
			values = append(values, arg.Closure.Proto.String())
		case compiler.VAL_NATIVE_FUNCTION:
			values = append(values, arg.Native)
		}
	}

	fmt.Println(values...)
	return compiler.NewNullValue(), nil
}
