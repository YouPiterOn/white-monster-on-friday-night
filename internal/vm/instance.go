package vm

import "youpiteron.dev/white-monster-on-friday-night/internal/compiler"

type Instance interface {
	ImplementInstanceInterface() Instance
	NumLocals() int
	Constants() []compiler.Value
}
