package compiler

type FuncSignature struct {
	CallArgs   []ValueType
	ReturnType ValueType
	Vararg     bool
}

type Variable struct {
	Name          string
	Slot          int
	Mutable       bool
	TypeOf        ValueType
	FuncSignature *FuncSignature
}

type Upvar struct {
	Name          string
	Mutable       bool
	LocalSlot     int
	SlotInParent  int
	IsFromParent  bool
	TypeOf        ValueType
	FuncSignature *FuncSignature
}

type Context interface {
	ImplementContext() Context
	DefineVariable(name string, mutable bool, typeOf ValueType) int
	DefineFunctionVariable(name string, mutable bool, typeOf ValueType, funcSignature *FuncSignature) int
	FindLocalVariable(name string) (*Variable, bool)
	FindUpvar(name string) (*Upvar, bool)
	FindVariable(name string) (*Variable, *Upvar, bool)
	AddInstruction(instruction Instruction) int
	SetInstruction(index int, instruction Instruction)
	AddConstant(value Value) int
	AddParam(param ValueType)

	// Getters
	VarSlot() int
	InstructionsLength() int
	Parent() Context
	ReturnType() ValueType
	Params() []ValueType
}
