package compiler

import "fmt"

type ModuleProto struct {
	numLocals    int
	instructions []Instruction
	constants    []Value
	functions    []FunctionProto
}

func (m *ModuleProto) ImplementProtoInterface() Proto {
	return m
}

func (m *ModuleProto) String() string {
	instructionsString := ""
	for _, instruction := range m.instructions {
		instructionsString += instruction.String() + "\n"
	}
	functionsString := ""
	for _, function := range m.functions {
		functionsString += function.String() + "\n"
	}
	return fmt.Sprintf("Instructions:\n%s\nConstants:\n%v\nFunctions:\n%s", instructionsString, m.constants, functionsString)
}

// ---------- Getters ----------

func (m *ModuleProto) NumLocals() int {
	return m.numLocals
}

func (m *ModuleProto) Instructions() []Instruction {
	return m.instructions
}

func (m *ModuleProto) Constants() []Value {
	return m.constants
}

func (m *ModuleProto) Functions() []FunctionProto {
	return m.functions
}

func BuildModuleProto(context Context, functions []FunctionProto) *ModuleProto {
	moduleContext := CastModuleContext(context)
	return &ModuleProto{
		numLocals:    moduleContext.currentVarSlot,
		instructions: moduleContext.instructions,
		constants:    moduleContext.constants,
		functions:    functions,
	}
}
