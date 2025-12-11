package compiler

type Variable struct {
	name    string
	slot    int
	mutable bool
}

type Upvar struct {
	name         string
	mutable      bool
	localSlot    int
	slotInParent int
	isFromParent bool
}

type Scope struct {
	parent           *Scope
	isFunctionScope  bool
	currentVarSlot   int
	currentUpvarSlot int
	variables        map[string]Variable
	upvarsMap        map[string]Upvar
}

func NewScope(isFunctionScope bool) *Scope {
	return &Scope{isFunctionScope: isFunctionScope, variables: make(map[string]Variable), upvarsMap: make(map[string]Upvar), currentVarSlot: 0, currentUpvarSlot: 0}
}

func (s *Scope) NewChildScope(isFunctionScope bool) *Scope {
	currentVarSlot := 0
	currentUpvarSlot := 0
	if isFunctionScope {
		currentVarSlot = s.currentVarSlot
		currentUpvarSlot = s.currentUpvarSlot
	}
	return &Scope{
		parent:           s,
		isFunctionScope:  isFunctionScope,
		variables:        make(map[string]Variable),
		upvarsMap:        make(map[string]Upvar),
		currentVarSlot:   currentVarSlot,
		currentUpvarSlot: currentUpvarSlot,
	}
}

func (s *Scope) DefineVariable(name string, mutable bool) int {
	slot := s.currentVarSlot
	s.variables[name] = Variable{name: name, slot: slot, mutable: mutable}
	s.currentVarSlot++
	return slot
}

func (s *Scope) FindLocalVariable(name string) (*Variable, bool) {
	variable, ok := s.variables[name]
	if ok {
		return &variable, true
	}
	if !s.isFunctionScope {
		variable, ok := s.parent.FindLocalVariable(name)
		if ok {
			return variable, true
		}
	}

	return nil, false
}

func (s *Scope) FindUpvar(name string) (*Upvar, bool) {
	if !s.isFunctionScope {
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
		upvar := Upvar{name: name, mutable: parentLocal.mutable, localSlot: s.currentUpvarSlot, slotInParent: parentLocal.slot, isFromParent: true}
		s.upvarsMap[name] = upvar
		s.currentUpvarSlot++
		return &upvar, true
	}
	parentUpvar, ok := s.parent.FindUpvar(name)
	if ok {
		upvar := Upvar{name: name, mutable: parentUpvar.mutable, localSlot: s.currentUpvarSlot, slotInParent: parentUpvar.slotInParent, isFromParent: false}
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
