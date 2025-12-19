package compiler

type BlockContext struct {
	parent         Context
	currentVarSlot int
	variables      map[string]Variable
}

func NewBlockContext(parent Context) *BlockContext {
	return &BlockContext{parent: parent, variables: make(map[string]Variable), currentVarSlot: parent.VarSlot()}
}

func (c *BlockContext) ImplementContextInterface() Context {
	return c
}

func (c *BlockContext) DefineVariable(name string, mutable bool, typeOf ValueType) int {
	slot := c.currentVarSlot
	c.variables[name] = Variable{Name: name, Slot: slot, Mutable: mutable, TypeOf: typeOf, FuncSignature: nil}
	c.currentVarSlot++
	return slot
}

func (c *BlockContext) DefineFunctionVariable(name string, mutable bool, typeOf ValueType, funcSignature *FuncSignature) int {
	slot := c.currentVarSlot
	c.variables[name] = Variable{Name: name, Slot: slot, Mutable: mutable, TypeOf: typeOf, FuncSignature: funcSignature}
	c.currentVarSlot++
	return slot
}

func (c *BlockContext) FindLocalVariable(name string) (*Variable, bool) {
	variable, ok := c.variables[name]
	if ok {
		return &variable, true
	}

	if c.parent == nil {
		return nil, false
	}
	parentVar, ok := c.parent.FindLocalVariable(name)
	if ok {
		return parentVar, true
	}

	return nil, false
}

func (c *BlockContext) FindUpvar(name string) (*Upvar, bool) {
	if c.parent == nil {
		return nil, false
	}
	return c.parent.FindUpvar(name)
}

func (c *BlockContext) FindVariable(name string) (*Variable, *Upvar, bool) {
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

func (c *BlockContext) AddInstruction(instruction Instruction) int {
	if c.parent == nil {
		panic("COMPILER ERROR: cannot add instruction to root block context")
	}
	return c.parent.AddInstruction(instruction)
}

func (c *BlockContext) SetInstruction(index int, instruction Instruction) {
	if c.parent == nil {
		panic("COMPILER ERROR: cannot set instruction in root block context")
	}
	c.parent.SetInstruction(index, instruction)
}

func (c *BlockContext) AddConstant(value Value) int {
	if c.parent == nil {
		panic("COMPILER ERROR: cannot add constant to root block context")
	}
	return c.parent.AddConstant(value)
}

func (c *BlockContext) AddParam(param ValueType) {
	if c.parent == nil {
		panic("COMPILER ERROR: cannot add param to root block context")
	}
	c.parent.AddParam(param)
}

// ---------- Getters ----------

func (c *BlockContext) VarSlot() int {
	return c.currentVarSlot
}

func (c *BlockContext) InstructionsLength() int {
	if c.parent == nil {
		panic("COMPILER ERROR: cannot get instructions length in root block context")
	}
	return c.parent.InstructionsLength()
}

func (c *BlockContext) Parent() Context {
	return c.parent
}

func (c *BlockContext) ReturnType() ValueType {
	if c.parent == nil {
		panic("COMPILER ERROR: cannot get return type in root block context")
	}
	return c.parent.ReturnType()
}

func (c *BlockContext) Params() []ValueType {
	if c.parent == nil {
		panic("COMPILER ERROR: cannot get params in root block context")
	}
	return c.parent.Params()
}
