package compiler

import "fmt"

type UpvarDesc struct {
	SlotInParent int
	IsFromParent bool
}

func (u *UpvarDesc) String() string {
	return fmt.Sprintf("SlotInParent: %d, IsFromParent: %t", u.SlotInParent, u.IsFromParent)
}

type FunctionBuilder struct {
	Name         string
	Params       []ValueType
	ReturnType   ValueType
	Instructions []Instruction
	Constants    []Value
}

func NewFunctionBuilder(name string, returnType ValueType) *FunctionBuilder {
	return &FunctionBuilder{Name: name, Params: []ValueType{}, ReturnType: returnType, Instructions: []Instruction{}, Constants: []Value{}}
}

func (f *FunctionBuilder) AddInstruction(instruction Instruction) {
	f.Instructions = append(f.Instructions, instruction)
}

func (f *FunctionBuilder) AddConstant(value Value) int {
	f.Constants = append(f.Constants, value)
	return len(f.Constants) - 1
}

func (f *FunctionBuilder) AddParam(param ValueType) {
	f.Params = append(f.Params, param)
}

func (f *FunctionBuilder) Build(scope *Scope) FunctionProto {
	numLocals := scope.currentVarSlot
	upvars := []UpvarDesc{}
	for _, upvar := range scope.upvarsMap {
		upvars = append(upvars, UpvarDesc{SlotInParent: upvar.SlotInParent, IsFromParent: upvar.IsFromParent})
	}
	return FunctionProto{Name: f.Name, NumLocals: numLocals, Params: f.Params, ReturnType: f.ReturnType, Instructions: f.Instructions, Upvars: upvars, Constants: f.Constants}
}

type FunctionProto struct {
	Name         string
	NumLocals    int
	Params       []ValueType
	ReturnType   ValueType
	Instructions []Instruction
	Upvars       []UpvarDesc
	Constants    []Value
}

func (f *FunctionProto) String() string {
	instructionsString := ""
	for _, instruction := range f.Instructions {
		instructionsString += instruction.String() + "\n"
	}
	upvarsString := ""
	for _, upvar := range f.Upvars {
		upvarsString += upvar.String() + "\n"
	}
	return fmt.Sprintf("%s:\nInstructions:\n%s\nUpvars:\n%s", f.Name, instructionsString, upvarsString)
}
