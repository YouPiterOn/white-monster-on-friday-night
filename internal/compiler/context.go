package compiler

type Callable struct {
	CallArgs   []ValueType
	ReturnType ValueType
}

type Variable struct {
	Name     string
	Slot     int
	Mutable  bool
	TypeOf   ValueType
	Callable *Callable
}

type Upvar struct {
	Name         string
	Mutable      bool
	LocalSlot    int
	SlotInParent int
	IsFromParent bool
	TypeOf       ValueType
	Callable     *Callable
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

func (s *Context) NewChildContext(isClosureContext bool) *Context {
	currentVarSlot := 0
	currentUpvarSlot := 0
	if !isClosureContext {
		currentVarSlot = s.currentVarSlot
		currentUpvarSlot = s.currentUpvarSlot
	}
	return &Context{
		parent:           s,
		isClosureContext: isClosureContext,
		variables:        make(map[string]Variable),
		upvarsMap:        make(map[string]Upvar),
		currentVarSlot:   currentVarSlot,
		currentUpvarSlot: currentUpvarSlot,
	}
}

func (s *Context) DefineVariable(name string, mutable bool, typeOf ValueType) int {
	slot := s.currentVarSlot
	s.variables[name] = Variable{Name: name, Slot: slot, Mutable: mutable, TypeOf: typeOf, Callable: nil}
	s.currentVarSlot++
	return slot
}

func (s *Context) DefineCallableVariable(name string, mutable bool, typeOf ValueType, callArgs []ValueType, returnType ValueType) int {
	slot := s.currentVarSlot
	s.variables[name] = Variable{Name: name, Slot: slot, Mutable: mutable, TypeOf: typeOf, Callable: &Callable{CallArgs: callArgs, ReturnType: returnType}}
	s.currentVarSlot++
	return slot
}

func (s *Context) FindLocalVariable(name string) (*Variable, bool) {
	variable, ok := s.variables[name]
	if ok {
		return &variable, true
	}
	if !s.isClosureContext {
		variable, ok := s.parent.FindLocalVariable(name)
		if ok {
			return variable, true
		}
	}

	return nil, false
}

func (s *Context) FindUpvar(name string) (*Upvar, bool) {
	if !s.isClosureContext {
		return s.parent.FindUpvar(name)
	}
	upvar, ok := s.upvarsMap[name]
	if ok {
		return &upvar, true
	}
	if s.parent == nil {
		return nil, false
	}

	parentLocal, ok := s.parent.FindLocalVariable(name)
	if ok {
		upvar := Upvar{
			Name:         name,
			Mutable:      parentLocal.Mutable,
			LocalSlot:    s.currentUpvarSlot,
			SlotInParent: parentLocal.Slot,
			IsFromParent: true,
			TypeOf:       parentLocal.TypeOf,
		}
		s.upvarsMap[name] = upvar
		s.currentUpvarSlot++
		return &upvar, true
	}
	parentUpvar, ok := s.parent.FindUpvar(name)
	if ok {
		upvar := Upvar{
			Name:         name,
			Mutable:      parentUpvar.Mutable,
			LocalSlot:    s.currentUpvarSlot,
			SlotInParent: parentUpvar.SlotInParent,
			IsFromParent: false,
			TypeOf:       parentUpvar.TypeOf,
		}
		s.upvarsMap[name] = upvar
		s.currentUpvarSlot++
		return &upvar, true
	}
	return nil, false
}

func (s *Context) FindVariable(name string) (*Variable, *Upvar, bool) {
	localVar, ok := s.FindLocalVariable(name)
	if ok {
		return localVar, nil, true
	}
	upvar, ok := s.FindUpvar(name)
	if ok {
		return nil, upvar, true
	}
	return nil, nil, false
}

func (s *Context) AddInstruction(instruction Instruction) {
	if s.isClosureContext {
		s.instructions = append(s.instructions, instruction)
	} else {
		s.parent.AddInstruction(instruction)
	}
}

func (s *Context) AddConstant(value Value) int {
	if s.isClosureContext {
		s.constants = append(s.constants, value)
		return len(s.constants) - 1
	} else {
		return s.parent.AddConstant(value)
	}
}

func (s *Context) AddParam(param ValueType) {
	if s.isClosureContext {
		s.params = append(s.params, param)
	} else {
		s.parent.AddParam(param)
	}
}

func (s *Context) BuildFunctionProto() (*FunctionProto, bool) {
	if !s.isClosureContext {
		return nil, false
	}
	numLocals := s.currentVarSlot
	upvars := []UpvarDesc{}
	for _, upvar := range s.upvarsMap {
		upvars = append(upvars, UpvarDesc{SlotInParent: upvar.SlotInParent, IsFromParent: upvar.IsFromParent})
	}
	return &FunctionProto{
		NumLocals:    numLocals,
		Params:       s.params,
		ReturnType:   s.returnType,
		Instructions: s.instructions,
		Upvars:       upvars,
		Constants:    s.constants,
	}, true
}
