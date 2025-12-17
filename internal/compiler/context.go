package compiler

type FuncSignature struct {
	CallArgs   []ValueType
	ReturnType ValueType
	Vararg     bool
}

type Variable struct {
	Name          string
	Slot          int
	Mutable       bool
	TypeOf        ValueType
	FuncSignature *FuncSignature
}

type Upvar struct {
	Name          string
	Mutable       bool
	LocalSlot     int
	SlotInParent  int
	IsFromParent  bool
	TypeOf        ValueType
	FuncSignature *FuncSignature
}

type Context struct {
	parent           *Context
	isClosureContext bool
	currentVarSlot   int
	currentUpvarSlot int
	variables        map[string]Variable
	upvarsMap        map[string]Upvar

	params       []ValueType
	returnType   ValueType
	instructions []Instruction
	constants    []Value
}

func NewContext(isClosureContext bool) *Context {
	return &Context{isClosureContext: isClosureContext, variables: make(map[string]Variable), upvarsMap: make(map[string]Upvar), currentVarSlot: 0, currentUpvarSlot: 0}
}

func (c *Context) NewChildContext(isClosureContext bool) *Context {
	currentVarSlot := 0
	currentUpvarSlot := 0
	if !isClosureContext {
		currentVarSlot = c.currentVarSlot
		currentUpvarSlot = c.currentUpvarSlot
	}
	return &Context{
		parent:           c,
		isClosureContext: isClosureContext,
		variables:        make(map[string]Variable),
		upvarsMap:        make(map[string]Upvar),
		currentVarSlot:   currentVarSlot,
		currentUpvarSlot: currentUpvarSlot,
	}
}

func (c *Context) DefineVariable(name string, mutable bool, typeOf ValueType) int {
	slot := c.currentVarSlot
	c.variables[name] = Variable{Name: name, Slot: slot, Mutable: mutable, TypeOf: typeOf, FuncSignature: nil}
	c.currentVarSlot++
	return slot
}

func (c *Context) DefineFunctionVariable(name string, mutable bool, typeOf ValueType, funcSignature *FuncSignature) int {
	slot := c.currentVarSlot
	c.variables[name] = Variable{Name: name, Slot: slot, Mutable: mutable, TypeOf: typeOf, FuncSignature: funcSignature}
	c.currentVarSlot++
	return slot
}

func (c *Context) FindLocalVariable(name string) (*Variable, bool) {
	variable, ok := c.variables[name]
	if ok {
		return &variable, true
	}
	if !c.isClosureContext {
		variable, ok := c.parent.FindLocalVariable(name)
		if ok {
			return variable, true
		}
	}

	return nil, false
}

func (c *Context) FindUpvar(name string) (*Upvar, bool) {
	if !c.isClosureContext {
		return c.parent.FindUpvar(name)
	}
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

func (c *Context) FindVariable(name string) (*Variable, *Upvar, bool) {
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

func (c *Context) AddInstruction(instruction Instruction) {
	if c.isClosureContext {
		c.instructions = append(c.instructions, instruction)
	} else {
		c.parent.AddInstruction(instruction)
	}
}

func (c *Context) AddConstant(value Value) int {
	if c.isClosureContext {
		c.constants = append(c.constants, value)
		return len(c.constants) - 1
	} else {
		return c.parent.AddConstant(value)
	}
}

func (c *Context) AddParam(param ValueType) {
	if c.isClosureContext {
		c.params = append(c.params, param)
	} else {
		c.parent.AddParam(param)
	}
}

func (c *Context) BuildFunctionProto() (*FunctionProto, bool) {
	if !c.isClosureContext {
		return nil, false
	}
	numLocals := c.currentVarSlot
	upvars := []UpvarDesc{}
	for _, upvar := range c.upvarsMap {
		upvars = append(upvars, UpvarDesc{SlotInParent: upvar.SlotInParent, IsFromParent: upvar.IsFromParent})
	}
	return &FunctionProto{
		NumLocals:    numLocals,
		Params:       c.params,
		ReturnType:   c.returnType,
		Instructions: c.instructions,
		Upvars:       upvars,
		Constants:    c.constants,
	}, true
}
