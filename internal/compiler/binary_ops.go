package compiler

import "youpiteron.dev/white-monster-on-friday-night/internal/lexer"

type BinaryOpInfo struct {
	ResultType *Type
	OpCode     OpCode
}

var BinaryOpTable = map[lexer.OperatorSubkind]map[TypeEnum]map[TypeEnum]BinaryOpInfo{
	lexer.OperatorPlus: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeInt(), OpCode: ADD_INT},
		},
	},
	lexer.OperatorMinus: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeInt(), OpCode: SUB_INT},
		},
	},
	lexer.OperatorStar: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeInt(), OpCode: MUL_INT},
		},
	},
	lexer.OperatorSlash: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeInt(), OpCode: DIV_INT},
		},
	},
	lexer.OperatorEqual: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeBool(), OpCode: EQ_INT},
		},
		TYPE_BOOL: {
			TYPE_BOOL: {ResultType: TypeBool(), OpCode: EQ_BOOL},
		},
	},
	lexer.OperatorNotEqual: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeBool(), OpCode: NE_INT},
		},
		TYPE_BOOL: {
			TYPE_BOOL: {ResultType: TypeBool(), OpCode: NE_BOOL},
		},
	},
	lexer.OperatorGreater: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeBool(), OpCode: GT_INT},
		},
	},
	lexer.OperatorGreaterEqual: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeBool(), OpCode: GTE_INT},
		},
	},
	lexer.OperatorLess: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeBool(), OpCode: LT_INT},
		},
	},
	lexer.OperatorLessEqual: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeBool(), OpCode: LTE_INT},
		},
	},
	lexer.OperatorAnd: {
		TYPE_BOOL: {
			TYPE_BOOL: {ResultType: TypeBool(), OpCode: AND_BOOL},
		},
	},
	lexer.OperatorOr: {
		TYPE_BOOL: {
			TYPE_BOOL: {ResultType: TypeBool(), OpCode: OR_BOOL},
		},
	},
}

func ResolveBinaryOp(
	op lexer.OperatorSubkind,
	left *Type,
	right *Type,
) (BinaryOpInfo, bool) {
	opMap, ok := BinaryOpTable[op]
	if !ok {
		return BinaryOpInfo{}, false
	}
	leftMap, ok := opMap[left.Type]
	if !ok {
		return BinaryOpInfo{}, false
	}
	info, ok := leftMap[right.Type]
	return info, ok
}
