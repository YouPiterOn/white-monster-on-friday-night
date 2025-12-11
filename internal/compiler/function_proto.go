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
	NumLocals    int
	NumParams    int
	Instructions []Instruction
	Upvars       []UpvarDesc
}

func NewFunctionBuilder(name string) *FunctionBuilder {
	return &FunctionBuilder{Name: name, NumLocals: 0, NumParams: 0, Instructions: []Instruction{}, Upvars: []UpvarDesc{}}
}

func (f *FunctionBuilder) AddInstruction(instruction Instruction) {
	f.Instructions = append(f.Instructions, instruction)
}

func (f *FunctionBuilder) AddUpvar(upvar UpvarDesc) {
	f.Upvars = append(f.Upvars, upvar)
}

func (f *FunctionBuilder) Build() FunctionProto {
	return FunctionProto{Name: f.Name, NumLocals: f.NumLocals, NumParams: f.NumParams, Instructions: f.Instructions, Upvars: f.Upvars}
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
