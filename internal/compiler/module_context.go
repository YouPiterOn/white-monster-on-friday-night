package compiler

import "youpiteron.dev/white-monster-on-friday-night/internal/ast"

type ModuleContext struct {
	currentVarSlot int
	variables      map[string]Variable

	returnType   *ast.Type
	instructions []Instruction
	constants    []Value
}

func NewModuleContext() *ModuleContext {
	return &ModuleContext{currentVarSlot: 0, variables: make(map[string]Variable), returnType: ast.TypeInt(), instructions: make([]Instruction, 0), constants: make([]Value, 0)}
}

func CastModuleContext(context Context) *ModuleContext {
	moduleContext, ok := context.(*ModuleContext)
	if !ok {
		panic("COMPILER ERROR: context is not a module context")
	}
	return moduleContext
}

func (c *ModuleContext) ImplementContextInterface() Context {
	return c
}

func (c *ModuleContext) DefineVariable(name string, mutable bool, typeOf *ast.Type) int {
	slot := c.currentVarSlot
	c.variables[name] = Variable{Name: name, Slot: slot, Mutable: mutable, TypeOf: typeOf, FuncSignature: nil}
	c.currentVarSlot++
	return slot
}

func (c *ModuleContext) DefineFunctionVariable(name string, mutable bool, typeOf *ast.Type, funcSignature *FuncSignature) int {
	slot := c.currentVarSlot
	c.variables[name] = Variable{Name: name, Slot: slot, Mutable: mutable, TypeOf: typeOf, FuncSignature: funcSignature}
	c.currentVarSlot++
	return slot
}

func (c *ModuleContext) FindLocalVariable(name string) (*Variable, bool) {
	variable, ok := c.variables[name]
	if ok {
		return &variable, true
	}
	return nil, false
}

func (c *ModuleContext) FindUpvar(name string) (*Upvar, bool) {
	return nil, false
}

func (c *ModuleContext) FindVariable(name string) (*Variable, *Upvar, bool) {
	localVar, ok := c.FindLocalVariable(name)
	if ok {
		return localVar, nil, true
	}
	return nil, nil, false
}

func (c *ModuleContext) AddInstruction(instruction Instruction) int {
	c.instructions = append(c.instructions, instruction)
	return len(c.instructions) - 1
}

func (c *ModuleContext) ClearInstructions() {
	c.instructions = make([]Instruction, 0)
}

func (c *ModuleContext) SetInstruction(index int, instruction Instruction) {
	c.instructions[index] = instruction
}

func (c *ModuleContext) AddConstant(value Value) int {
	c.constants = append(c.constants, value)
	return len(c.constants) - 1
}

func (c *ModuleContext) AddParam(param *ast.Type) {
	panic("COMPILER ERROR: cannot add param to module context")
}

// ---------- Getters ----------

func (c *ModuleContext) VarSlot() int {
	return c.currentVarSlot
}

func (c *ModuleContext) InstructionsLength() int {
	return len(c.instructions)
}

func (c *ModuleContext) Parent() Context {
	return nil
}

func (c *ModuleContext) ReturnType() *ast.Type {
	return c.returnType
}

func (c *ModuleContext) Params() []*ast.Type {
	return []*ast.Type{}
}
