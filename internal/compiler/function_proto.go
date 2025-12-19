package compiler

import "fmt"

type UpvarDesc struct {
	SlotInParent int
	IsFromParent bool
}

func (u *UpvarDesc) String() string {
	return fmt.Sprintf("SlotInParent: %d, IsFromParent: %t", u.SlotInParent, u.IsFromParent)
}

type FunctionProto struct {
	numLocals    int
	instructions []Instruction
	upvars       []UpvarDesc
	constants    []Value
}

func (f *FunctionProto) ImplementProtoInterface() Proto {
	return f
}

func (f *FunctionProto) String() string {
	instructionsString := ""
	for _, instruction := range f.instructions {
		instructionsString += instruction.String() + "\n"
	}
	upvarsString := ""
	for _, upvar := range f.upvars {
		upvarsString += upvar.String() + "\n"
	}
	return fmt.Sprintf("Instructions:\n%s\nUpvars:\n%s", instructionsString, upvarsString)
}

// ---------- Getters ----------

func (f *FunctionProto) NumLocals() int {
	return f.numLocals
}

func (f *FunctionProto) Instructions() []Instruction {
	return f.instructions
}

func (f *FunctionProto) Constants() []Value {
	return f.constants
}

func (f *FunctionProto) Upvars() []UpvarDesc {
	return f.upvars
}

func BuildFunctionProto(context Context) *FunctionProto {
	functionContext := CastFunctionContext(context)
	numLocals := functionContext.currentVarSlot
	upvars := []UpvarDesc{}
	for _, upvar := range functionContext.upvarsMap {
		upvars = append(upvars, UpvarDesc{SlotInParent: upvar.SlotInParent, IsFromParent: upvar.IsFromParent})
	}
	return &FunctionProto{
		numLocals:    numLocals,
		instructions: functionContext.instructions,
		upvars:       upvars,
		constants:    functionContext.constants,
	}
}
