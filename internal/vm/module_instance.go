package vm

import "youpiteron.dev/white-monster-on-friday-night/internal/compiler"

type ModuleInstance struct {
	numLocals int
	constants []compiler.Value
	functions []compiler.FunctionProto
}

func NewModuleInstance(moduleProto *compiler.ModuleProto) *ModuleInstance {
	return &ModuleInstance{numLocals: moduleProto.NumLocals(), constants: moduleProto.Constants(), functions: moduleProto.Functions()}
}
