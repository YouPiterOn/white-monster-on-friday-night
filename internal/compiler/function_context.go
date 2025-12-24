package compiler

import "youpiteron.dev/white-monster-on-friday-night/internal/ast"

type FunctionContext struct {
	parent           Context
	currentVarSlot   int
	currentUpvarSlot int
	variables        map[string]Variable
	upvarsMap        map[string]Upvar

	params       []*ast.Type
	returnType   *ast.Type
	instructions []Instruction
	constants    []Value
}

func NewFunctionContext(parent Context, returnType *ast.Type) *FunctionContext {
	return &FunctionContext{parent: parent, variables: make(map[string]Variable), upvarsMap: make(map[string]Upvar), currentVarSlot: 0, currentUpvarSlot: 0, returnType: returnType}
}

func CastFunctionContext(context Context) *FunctionContext {
	functionContext, ok := context.(*FunctionContext)
	if !ok {
		panic("COMPILER ERROR: context is not a function context")
	}
	return functionContext
}

func (c *FunctionContext) ImplementContextInterface() Context {
	return c
}

func (c *FunctionContext) DefineVariable(name string, mutable bool, typeOf *ast.Type) int {
	slot := c.currentVarSlot
	c.variables[name] = Variable{Name: name, Slot: slot, Mutable: mutable, TypeOf: typeOf, FuncSignature: nil}
	c.currentVarSlot++
	return slot
}

func (c *FunctionContext) DefineFunctionVariable(name string, mutable bool, typeOf *ast.Type, funcSignature *FuncSignature) int {
	slot := c.currentVarSlot
	c.variables[name] = Variable{Name: name, Slot: slot, Mutable: mutable, TypeOf: typeOf, FuncSignature: funcSignature}
	c.currentVarSlot++
	return slot
}

func (c *FunctionContext) FindLocalVariable(name string) (*Variable, bool) {
	variable, ok := c.variables[name]
	if ok {
		return &variable, true
	}

	return nil, false
}

func (c *FunctionContext) FindUpvar(name string) (*Upvar, bool) {
	upvar, ok := c.upvarsMap[name]
	if ok {
		return &upvar, true
	}
	if c.parent == nil {
		return nil, false
	}

	parentLocal, ok := c.parent.FindLocalVariable(name)
	if ok {
		upvar := Upvar{
			Name:         name,
			Mutable:      parentLocal.Mutable,
			LocalSlot:    c.currentUpvarSlot,
			SlotInParent: parentLocal.Slot,
			IsFromParent: true,
			TypeOf:       parentLocal.TypeOf,
		}
		c.upvarsMap[name] = upvar
		c.currentUpvarSlot++
		return &upvar, true
	}
	parentUpvar, ok := c.parent.FindUpvar(name)
	if ok {
		upvar := Upvar{
			Name:         name,
			Mutable:      parentUpvar.Mutable,
			LocalSlot:    c.currentUpvarSlot,
			SlotInParent: parentUpvar.SlotInParent,
			IsFromParent: false,
			TypeOf:       parentUpvar.TypeOf,
		}
		c.upvarsMap[name] = upvar
		c.currentUpvarSlot++
		return &upvar, true
	}
	return nil, false
}

func (c *FunctionContext) FindVariable(name string) (*Variable, *Upvar, bool) {
	localVar, ok := c.FindLocalVariable(name)
	if ok {
		return localVar, nil, true
	}
	upvar, ok := c.FindUpvar(name)
	if ok {
		return nil, upvar, true
	}
	return nil, nil, false
}

func (c *FunctionContext) AddInstruction(instruction Instruction) int {
	c.instructions = append(c.instructions, instruction)
	return len(c.instructions) - 1
}

func (c *FunctionContext) SetInstruction(index int, instruction Instruction) {
	c.instructions[index] = instruction
}

func (c *FunctionContext) AddConstant(value Value) int {
	c.constants = append(c.constants, value)
	return len(c.constants) - 1
}

func (c *FunctionContext) AddParam(param *ast.Type) {
	c.params = append(c.params, param)
}

// ---------- Getters ----------

func (c *FunctionContext) VarSlot() int {
	return c.currentVarSlot
}

func (c *FunctionContext) InstructionsLength() int {
	return len(c.instructions)
}

func (c *FunctionContext) Parent() Context {
	return c.parent
}

func (c *FunctionContext) ReturnType() *ast.Type {
	return c.returnType
}

func (c *FunctionContext) Params() []*ast.Type {
	return c.params
}
