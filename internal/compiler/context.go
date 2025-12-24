package compiler

type FuncSignature struct {
	CallArgs   []Type
	ReturnType Type
	Vararg     bool
}

type Variable struct {
	Name          string
	Slot          int
	Mutable       bool
	TypeOf        Type
	FuncSignature *FuncSignature
}

type Upvar struct {
	Name          string
	Mutable       bool
	LocalSlot     int
	SlotInParent  int
	IsFromParent  bool
	TypeOf        Type
	FuncSignature *FuncSignature
}

type Context interface {
	ImplementContextInterface() Context
	DefineVariable(name string, mutable bool, typeOf Type) int
	DefineFunctionVariable(name string, mutable bool, typeOf Type, funcSignature *FuncSignature) int
	FindLocalVariable(name string) (*Variable, bool)
	FindUpvar(name string) (*Upvar, bool)
	FindVariable(name string) (*Variable, *Upvar, bool)
	AddInstruction(instruction Instruction) int
	SetInstruction(index int, instruction Instruction)
	AddConstant(value Value) int
	AddParam(param Type)

	// Getters
	VarSlot() int
	InstructionsLength() int
	Parent() Context
	ReturnType() Type
	Params() []Type
}
