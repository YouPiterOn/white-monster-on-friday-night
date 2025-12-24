package compiler

import "youpiteron.dev/white-monster-on-friday-night/internal/lexer"

type BinaryOpInfo struct {
	ResultType Type
	OpCode     OpCode
}

var BinaryOpTable = map[lexer.OperatorSubkind]map[Type]map[Type]BinaryOpInfo{
	lexer.OperatorPlus: {
		TypeInt(): {
			TypeInt(): {ResultType: TypeInt(), OpCode: ADD_INT},
		},
	},
	lexer.OperatorMinus: {
		TypeInt(): {
			TypeInt(): {ResultType: TypeInt(), OpCode: SUB_INT},
		},
	},
	lexer.OperatorStar: {
		TypeInt(): {
			TypeInt(): {ResultType: TypeInt(), OpCode: MUL_INT},
		},
	},
	lexer.OperatorSlash: {
		TypeInt(): {
			TypeInt(): {ResultType: TypeInt(), OpCode: DIV_INT},
		},
	},
	lexer.OperatorEqual: {
		TypeInt(): {
			TypeInt(): {ResultType: TypeBool(), OpCode: EQ_INT},
		},
		TypeBool(): {
			TypeBool(): {ResultType: TypeBool(), OpCode: EQ_BOOL},
		},
	},
	lexer.OperatorNotEqual: {
		TypeInt(): {
			TypeInt(): {ResultType: TypeBool(), OpCode: NE_INT},
		},
		TypeBool(): {
			TypeBool(): {ResultType: TypeBool(), OpCode: NE_BOOL},
		},
	},
	lexer.OperatorGreater: {
		TypeInt(): {
			TypeInt(): {ResultType: TypeBool(), OpCode: GT_INT},
		},
	},
	lexer.OperatorGreaterEqual: {
		TypeInt(): {
			TypeInt(): {ResultType: TypeBool(), OpCode: GTE_INT},
		},
	},
	lexer.OperatorLess: {
		TypeInt(): {
			TypeInt(): {ResultType: TypeBool(), OpCode: LT_INT},
		},
	},
	lexer.OperatorLessEqual: {
		TypeInt(): {
			TypeInt(): {ResultType: TypeBool(), OpCode: LTE_INT},
		},
	},
	lexer.OperatorAnd: {
		TypeBool(): {
			TypeBool(): {ResultType: TypeBool(), OpCode: AND_BOOL},
		},
	},
	lexer.OperatorOr: {
		TypeBool(): {
			TypeBool(): {ResultType: TypeBool(), OpCode: OR_BOOL},
		},
	},
}

func ResolveBinaryOp(
	op lexer.OperatorSubkind,
	left Type,
	right Type,
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
