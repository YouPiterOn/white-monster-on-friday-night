package compiler

import "youpiteron.dev/white-monster-on-friday-night/internal/ast"

type GlobalTable struct {
	ids       map[string]int
	variables []Variable
	types     map[string]*Type
}

func NewGlobalTable() *GlobalTable {
	return &GlobalTable{ids: make(map[string]int), variables: make([]Variable, 0)}
}

func (g *GlobalTable) DefineVariable(name string, mutable bool, typeOf *ast.Type) int {
	slot := len(g.variables)
	g.variables = append(g.variables, Variable{Name: name, Slot: slot, Mutable: mutable, TypeOf: typeOf})
	g.ids[name] = slot
	return slot
}

func (g *GlobalTable) FindVariable(name string) (*Variable, bool) {
	slot, ok := g.ids[name]
	if ok {
		return &g.variables[slot], true
	}
	return nil, false
}

func (g *GlobalTable) VariablesLength() int {
	return len(g.variables)
}

func (g *GlobalTable) DefineType(name string, typeOf *Type) {
	g.types[name] = typeOf
}

func (g *GlobalTable) FindType(name string) (*Type, bool) {
	typeOf, ok := g.types[name]
	if ok {
		return typeOf, true
	}
	return nil, false
}
