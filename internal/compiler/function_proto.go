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
	return fmt.Sprintf("Instructions:\n%s\nUpvars:\n%s", instructionsString, upvarsString)
}

func BuildFunctionProto(context Context) *FunctionProto {
	functionContext := CastFunctionContext(context)
	numLocals := functionContext.currentVarSlot
	upvars := []UpvarDesc{}
	for _, upvar := range functionContext.upvarsMap {
		upvars = append(upvars, UpvarDesc{SlotInParent: upvar.SlotInParent, IsFromParent: upvar.IsFromParent})
	}
	return &FunctionProto{
		NumLocals:    numLocals,
		Params:       functionContext.params,
		ReturnType:   functionContext.returnType,
		Instructions: functionContext.instructions,
		Upvars:       upvars,
		Constants:    functionContext.constants,
	}
}
