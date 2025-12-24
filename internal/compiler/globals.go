package compiler

type GlobalTable struct {
	ids       map[string]int
	variables []Variable
}

func NewGlobalTable() *GlobalTable {
	return &GlobalTable{ids: make(map[string]int), variables: make([]Variable, 0)}
}

func (g *GlobalTable) DefineVariable(name string, mutable bool, typeOf Type) int {
	slot := len(g.variables)
	g.variables = append(g.variables, Variable{Name: name, Slot: slot, Mutable: mutable, TypeOf: typeOf, FuncSignature: nil})
	g.ids[name] = slot
	return slot
}

func (g *GlobalTable) DefineFunctionVariable(name string, mutable bool, typeOf Type, funcSignature *FuncSignature) int {
	slot := len(g.variables)
	g.variables = append(g.variables, Variable{Name: name, Slot: slot, Mutable: mutable, TypeOf: typeOf, FuncSignature: funcSignature})
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

func (g *GlobalTable) Length() int {
	return len(g.variables)
}
