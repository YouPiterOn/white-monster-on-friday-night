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

func InstrLoadConst(reg int, constIndex int) Instruction {
	return Instruction{
		OpCode: LOAD_CONST,
		Args:   []int{reg, constIndex},
	}
}

func InstrLoadVar(reg int, slot int) Instruction {
	return Instruction{
		OpCode: LOAD_VAR,
		Args:   []int{reg, slot},
	}
}

func InstrLoadUpvar(reg int, slot int) Instruction {
	return Instruction{
		OpCode: LOAD_UPVAR,
		Args:   []int{reg, slot},
	}
}

func InstrStoreVar(reg int, slot int) Instruction {
	return Instruction{
		OpCode: STORE_VAR,
		Args:   []int{reg, slot},
	}
}

func InstrAssignUpvar(reg int, slot int) Instruction {
	return Instruction{
		OpCode: ASSIGN_UPVAR,
		Args:   []int{reg, slot},
	}
}

func InstrAdd(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: ADD,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func InstrSub(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: SUB,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func InstrMul(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: MUL,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func InstrDiv(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: DIV,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func InstrClosure(resultReg int, functionSlot int) Instruction {
	return Instruction{
		OpCode: CLOSURE,
		Args:   []int{resultReg, functionSlot},
	}
}

func InstrReturn(reg int) Instruction {
	return Instruction{
		OpCode: RETURN,
		Args:   []int{reg},
	}
}

func InstrCall(resultReg int, functionReg int, args []int) Instruction {
	return Instruction{
		OpCode: CALL,
		Args:   append([]int{resultReg, functionReg}, args...),
	}
}
