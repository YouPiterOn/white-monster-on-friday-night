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
	LOAD_UPVAR
	STORE_VAR
	ASSIGN_UPVAR
	ADD
	SUB
	MUL
	DIV
	CLOSURE
	CALL
	RETURN
)

func (o OpCode) String() string {
	return [...]string{
		"LOAD_CONST",
		"LOAD_VAR",
		"LOAD_UPVAR",
		"STORE_VAR",
		"ASSIGN_UPVAR",
		"ADD",
		"SUB",
		"MUL",
		"DIV",
		"CLOSURE",
		"CALL",
		"RETURN",
	}[o]
}

func LoadConst(reg int, value int) Instruction {
	return Instruction{
		OpCode: LOAD_CONST,
		Args:   []int{reg, value},
	}
}

func LoadVar(reg int, slot int) Instruction {
	return Instruction{
		OpCode: LOAD_VAR,
		Args:   []int{reg, slot},
	}
}

func LoadUpvar(reg int, slot int) Instruction {
	return Instruction{
		OpCode: LOAD_UPVAR,
		Args:   []int{reg, slot},
	}
}

func StoreVar(reg int, slot int) Instruction {
	return Instruction{
		OpCode: STORE_VAR,
		Args:   []int{reg, slot},
	}
}

func AssignUpvar(reg int, slot int) Instruction {
	return Instruction{
		OpCode: ASSIGN_UPVAR,
		Args:   []int{reg, slot},
	}
}

func Add(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: ADD,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func Sub(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: SUB,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func Mul(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: MUL,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func Div(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: DIV,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func Return(reg int) Instruction {
	return Instruction{
		OpCode: RETURN,
		Args:   []int{reg},
	}
}

func Call(resultReg int, functionReg int, args []int) Instruction {
	return Instruction{
		OpCode: CALL,
		Args:   append([]int{resultReg, functionReg}, args...),
	}
}
