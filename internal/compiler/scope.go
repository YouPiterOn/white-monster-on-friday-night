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

type Scope struct {
	parent           *Scope
	isClosureScope   bool
	currentVarSlot   int
	currentUpvarSlot int
	variables        map[string]Variable
	upvarsMap        map[string]Upvar
}

func NewScope(isClosureScope bool) *Scope {
	return &Scope{isClosureScope: isClosureScope, variables: make(map[string]Variable), upvarsMap: make(map[string]Upvar), currentVarSlot: 0, currentUpvarSlot: 0}
}

func (s *Scope) NewChildScope(isClosureScope bool) *Scope {
	currentVarSlot := 0
	currentUpvarSlot := 0
	if !isClosureScope {
		currentVarSlot = s.currentVarSlot
		currentUpvarSlot = s.currentUpvarSlot
	}
	return &Scope{
		parent:           s,
		isClosureScope:   isClosureScope,
		variables:        make(map[string]Variable),
		upvarsMap:        make(map[string]Upvar),
		currentVarSlot:   currentVarSlot,
		currentUpvarSlot: currentUpvarSlot,
	}
}

func (s *Scope) DefineVariable(name string, mutable bool, typeOf ValueType) int {
	slot := s.currentVarSlot
	s.variables[name] = Variable{Name: name, Slot: slot, Mutable: mutable, TypeOf: typeOf, Callable: nil}
	s.currentVarSlot++
	return slot
}

func (s *Scope) DefineCallableVariable(name string, mutable bool, typeOf ValueType, callArgs []ValueType, returnType ValueType) int {
	slot := s.currentVarSlot
	s.variables[name] = Variable{Name: name, Slot: slot, Mutable: mutable, TypeOf: typeOf, Callable: &Callable{CallArgs: callArgs, ReturnType: returnType}}
	s.currentVarSlot++
	return slot
}

func (s *Scope) FindLocalVariable(name string) (*Variable, bool) {
	variable, ok := s.variables[name]
	if ok {
		return &variable, true
	}
	if !s.isClosureScope {
		variable, ok := s.parent.FindLocalVariable(name)
		if ok {
			return variable, true
		}
	}

	return nil, false
}

func (s *Scope) FindUpvar(name string) (*Upvar, bool) {
	if !s.isClosureScope {
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

func (s *Scope) FindVariable(name string) (*Variable, *Upvar, bool) {
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
