package compiler

type Proto interface {
	ImplementProtoInterface() Proto
	NumLocals() int
	Instructions() []Instruction
	Constants() []Value

	String() string
}
