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
	NumParams    int
	Instructions []Instruction
}

func NewFunctionBuilder(name string, numParams int) *FunctionBuilder {
	return &FunctionBuilder{Name: name, NumParams: numParams, Instructions: []Instruction{}}
}

func (f *FunctionBuilder) AddInstruction(instruction Instruction) {
	f.Instructions = append(f.Instructions, instruction)
}

func (f *FunctionBuilder) Build(scope *Scope) FunctionProto {
	numLocals := scope.currentVarSlot
	upvars := []UpvarDesc{}
	for _, upvar := range scope.upvarsMap {
		upvars = append(upvars, UpvarDesc{SlotInParent: upvar.slotInParent, IsFromParent: upvar.isFromParent})
	}
	return FunctionProto{Name: f.Name, NumLocals: numLocals, NumParams: f.NumParams, Instructions: f.Instructions, Upvars: upvars}
}

type FunctionProto struct {
	Name         string
	NumLocals    int
	NumParams    int
	Instructions []Instruction
	Upvars       []UpvarDesc
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
