package vm

import "youpiteron.dev/white-monster-on-friday-night/internal/compiler"

type FunctionInstance struct {
	numLocals int
	constants []compiler.Value
	upvars    []compiler.UpvarDesc
}

func NewFunctionInstance(functionProto *compiler.FunctionProto) *FunctionInstance {
	return &FunctionInstance{numLocals: functionProto.NumLocals(), constants: functionProto.Constants(), upvars: functionProto.Upvars()}
}
