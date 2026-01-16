package compiler

import "youpiteron.dev/white-monster-on-friday-night/internal/ast"

type Variable struct {
	Name    string
	Slot    int
	Mutable bool
	TypeOf  *ast.Type
}

type Upvar struct {
	Name         string
	Mutable      bool
	LocalSlot    int
	SlotInParent int
	IsFromParent bool
	TypeOf       *ast.Type
}

type Context interface {
	ImplementContextInterface() Context
	DefineVariable(name string, mutable bool, typeOf *ast.Type) int
	FindLocalVariable(name string) (*Variable, bool)
	FindUpvar(name string) (*Upvar, bool)
	FindVariable(name string) (*Variable, *Upvar, bool)
	AddInstruction(instruction Instruction) int
	SetInstruction(index int, instruction Instruction)
	AddConstant(value Value) int
	AddParam(param *ast.Type)

	// Getters
	VarSlot() int
	InstructionsLength() int
	Parent() Context
	ReturnType() *ast.Type
	Params() []*ast.Type
}
