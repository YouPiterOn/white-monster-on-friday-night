package compiler

type Type int

const (
	Int Type = iota
	Float
	String
	Bool
	Function
)

type Variable struct {
	Name  string
	Type  Type
	Value any
}
