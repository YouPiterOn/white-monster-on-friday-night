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
	LOAD_GLOBAL
	LOAD_UPVAR
	STORE_VAR
	ASSIGN_GLOBAL
	ASSIGN_UPVAR
	ADD_INT
	SUB_INT
	MUL_INT
	DIV_INT
	EQ_INT
	EQ_BOOL
	NE_INT
	NE_BOOL
	GT_INT
	GTE_INT
	LT_INT
	LTE_INT
	AND_BOOL
	OR_BOOL
	CLOSURE
	CALL
	RETURN
	JUMP_IF_FALSE
	JUMP
)

func (o OpCode) String() string {
	return [...]string{
		"LOAD_CONST",
		"LOAD_VAR",
		"LOAD_GLOBAL",
		"LOAD_UPVAR",
		"STORE_VAR",
		"ASSIGN_GLOBAL",
		"ASSIGN_UPVAR",
		"ADD_INT",
		"SUB_INT",
		"MUL_INT",
		"DIV_INT",
		"EQ_INT",
		"EQ_BOOL",
		"NE_INT",
		"NE_BOOL",
		"GT_INT",
		"GTE_INT",
		"LT_INT",
		"LTE_INT",
		"AND_BOOL",
		"OR_BOOL",
		"CLOSURE",
		"CALL",
		"RETURN",
		"JUMP_IF_FALSE",
		"JUMP",
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

func InstrLoadGlobal(reg int, slot int) Instruction {
	return Instruction{
		OpCode: LOAD_GLOBAL,
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

func InstrAssignGlobal(reg int, slot int) Instruction {
	return Instruction{
		OpCode: ASSIGN_GLOBAL,
		Args:   []int{reg, slot},
	}
}

func InstrAssignUpvar(reg int, slot int) Instruction {
	return Instruction{
		OpCode: ASSIGN_UPVAR,
		Args:   []int{reg, slot},
	}
}

func InstrBinary(op OpCode, regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: op,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func InstrAddInt(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: ADD_INT,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func InstrSubInt(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: SUB_INT,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func InstrMulInt(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: MUL_INT,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func InstrDivInt(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: DIV_INT,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func InstrEqualInt(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: EQ_INT,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func InstrEqualBool(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: EQ_BOOL,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func InstrNotEqualInt(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: NE_INT,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func InstrNotEqualBool(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: NE_BOOL,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func InstrGreaterInt(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: GT_INT,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func InstrGreaterEqualInt(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: GTE_INT,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func InstrLessInt(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: LT_INT,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func InstrLessEqualInt(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: LTE_INT,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func InstrAndBool(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: AND_BOOL,
		Args:   []int{regResult, regLeft, regRight},
	}
}

func InstrOrBool(regResult int, regLeft int, regRight int) Instruction {
	return Instruction{
		OpCode: OR_BOOL,
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

func InstrJumpIfFalse(reg int, target int) Instruction {
	return Instruction{
		OpCode: JUMP_IF_FALSE,
		Args:   []int{reg, target},
	}
}

func InstrJump(target int) Instruction {
	return Instruction{
		OpCode: JUMP,
		Args:   []int{target},
	}
}
