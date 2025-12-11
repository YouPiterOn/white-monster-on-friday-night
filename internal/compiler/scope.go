package compiler

type Scope struct {
	parent    *Scope
	constants map[string]bool // set of constants
}

func NewScope() *Scope {
	return &Scope{constants: make(map[string]bool)}
}

func (s *Scope) NewChildScope() *Scope {
	return &Scope{parent: s, constants: make(map[string]bool)}
}

func (s *Scope) DefineConstant(name string) {
	s.constants[name] = true
}

func (s *Scope) IsConstant(name string) bool {
	return s.constants[name]
}
