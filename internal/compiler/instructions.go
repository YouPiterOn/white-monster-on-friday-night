package compiler

import (
	"fmt"
)

type Instruction struct {
	OpCode OpCode
	Args   []int
}

func (i *Instruction) String() string {
	return fmt.Sprintf("%s %v", i.OpCode, i.Args)
}

type OpCode int

const (
	LOAD_CONST OpCode = iota
	LOAD_VAR
	STORE_VAR
	ADD
	SUB
	MUL
	DIV
	CALL
	RETURN
)

func (o OpCode) String() string {
	return [...]string{
		"LOAD_CONST",
		"LOAD_VAR",
		"STORE_VAR",
		"ADD",
		"SUB",
		"MUL",
		"DIV",
		"CALL",
		"RETURN",
	}[o]
}
