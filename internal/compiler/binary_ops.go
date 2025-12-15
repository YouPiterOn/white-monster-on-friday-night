package compiler

import "youpiteron.dev/white-monster-on-friday-night/internal/lexer"

type BinaryOpInfo struct {
	ResultType ValueType
	OpCode     OpCode
}

var BinaryOpTable = map[lexer.OperatorSubkind]map[ValueType]map[ValueType]BinaryOpInfo{
	lexer.OperatorPlus: {
		VAL_INT: {
			VAL_INT: {ResultType: VAL_INT, OpCode: ADD_INT},
		},
	},
	lexer.OperatorMinus: {
		VAL_INT: {
			VAL_INT: {ResultType: VAL_INT, OpCode: SUB_INT},
		},
	},
	lexer.OperatorStar: {
		VAL_INT: {
			VAL_INT: {ResultType: VAL_INT, OpCode: MUL_INT},
		},
	},
	lexer.OperatorSlash: {
		VAL_INT: {
			VAL_INT: {ResultType: VAL_INT, OpCode: DIV_INT},
		},
	},
	lexer.OperatorEqual: {
		VAL_INT: {
			VAL_INT: {ResultType: VAL_BOOL, OpCode: EQ_INT},
		},
		VAL_BOOL: {
			VAL_BOOL: {ResultType: VAL_BOOL, OpCode: EQ_BOOL},
		},
	},
	lexer.OperatorNotEqual: {
		VAL_INT: {
			VAL_INT: {ResultType: VAL_BOOL, OpCode: NE_INT},
		},
		VAL_BOOL: {
			VAL_BOOL: {ResultType: VAL_BOOL, OpCode: NE_BOOL},
		},
	},
	lexer.OperatorGreater: {
		VAL_INT: {
			VAL_INT: {ResultType: VAL_BOOL, OpCode: GT_INT},
		},
	},
	lexer.OperatorGreaterEqual: {
		VAL_INT: {
			VAL_INT: {ResultType: VAL_BOOL, OpCode: GTE_INT},
		},
	},
	lexer.OperatorLess: {
		VAL_INT: {
			VAL_INT: {ResultType: VAL_BOOL, OpCode: LT_INT},
		},
	},
	lexer.OperatorLessEqual: {
		VAL_INT: {
			VAL_INT: {ResultType: VAL_BOOL, OpCode: LTE_INT},
		},
	},
	lexer.OperatorAnd: {
		VAL_BOOL: {
			VAL_BOOL: {ResultType: VAL_BOOL, OpCode: AND_BOOL},
		},
	},
	lexer.OperatorOr: {
		VAL_BOOL: {
			VAL_BOOL: {ResultType: VAL_BOOL, OpCode: OR_BOOL},
		},
	},
}

func ResolveBinaryOp(
	op lexer.OperatorSubkind,
	left ValueType,
	right ValueType,
) (BinaryOpInfo, bool) {
	opMap, ok := BinaryOpTable[op]
	if !ok {
		return BinaryOpInfo{}, false
	}
	leftMap, ok := opMap[left]
	if !ok {
		return BinaryOpInfo{}, false
	}
	info, ok := leftMap[right]
	return info, ok
}
