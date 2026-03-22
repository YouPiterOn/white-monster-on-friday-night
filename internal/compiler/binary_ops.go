package compiler

import (
	"youpiteron.dev/white-monster-on-friday-night/internal/lexer"
)

type BinaryOpInfo struct {
	ResultType *Type
	OpCode     OpCode
}

var BinaryOpTable = map[lexer.OperatorSubkind]map[PrimitiveType]map[PrimitiveType]BinaryOpInfo{
	lexer.OperatorPlus: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeInt(), OpCode: ADD_INT},
		},
		TYPE_FLOAT: {
			TYPE_FLOAT: {ResultType: TypeFloat(), OpCode: ADD_FLOAT},
		},
		TYPE_STRING: {
			TYPE_STRING: {ResultType: TypeString(), OpCode: ADD_STRING},
		},
	},
	lexer.OperatorMinus: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeInt(), OpCode: SUB_INT},
		},
		TYPE_FLOAT: {
			TYPE_FLOAT: {ResultType: TypeFloat(), OpCode: SUB_FLOAT},
		},
	},
	lexer.OperatorStar: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeInt(), OpCode: MUL_INT},
		},
		TYPE_FLOAT: {
			TYPE_FLOAT: {ResultType: TypeFloat(), OpCode: MUL_FLOAT},
		},
	},
	lexer.OperatorSlash: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeInt(), OpCode: DIV_INT},
		},
		TYPE_FLOAT: {
			TYPE_FLOAT: {ResultType: TypeFloat(), OpCode: DIV_FLOAT},
		},
	},
	lexer.OperatorEqual: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeBool(), OpCode: EQ_INT},
		},
		TYPE_BOOL: {
			TYPE_BOOL: {ResultType: TypeBool(), OpCode: EQ_BOOL},
		},
		TYPE_FLOAT: {
			TYPE_FLOAT: {ResultType: TypeBool(), OpCode: EQ_FLOAT},
		},
		TYPE_STRING: {
			TYPE_STRING: {ResultType: TypeBool(), OpCode: EQ_STRING},
		},
	},
	lexer.OperatorNotEqual: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeBool(), OpCode: NE_INT},
		},
		TYPE_BOOL: {
			TYPE_BOOL: {ResultType: TypeBool(), OpCode: NE_BOOL},
		},
		TYPE_FLOAT: {
			TYPE_FLOAT: {ResultType: TypeBool(), OpCode: NE_FLOAT},
		},
		TYPE_STRING: {
			TYPE_STRING: {ResultType: TypeBool(), OpCode: NE_STRING},
		},
	},
	lexer.OperatorGreater: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeBool(), OpCode: GT_INT},
		},
		TYPE_FLOAT: {
			TYPE_FLOAT: {ResultType: TypeBool(), OpCode: GT_FLOAT},
		},
	},
	lexer.OperatorGreaterEqual: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeBool(), OpCode: GTE_INT},
		},
		TYPE_FLOAT: {
			TYPE_FLOAT: {ResultType: TypeBool(), OpCode: GTE_FLOAT},
		},
	},
	lexer.OperatorLess: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeBool(), OpCode: LT_INT},
		},
		TYPE_FLOAT: {
			TYPE_FLOAT: {ResultType: TypeBool(), OpCode: LT_FLOAT},
		},
	},
	lexer.OperatorLessEqual: {
		TYPE_INT: {
			TYPE_INT: {ResultType: TypeBool(), OpCode: LTE_INT},
		},
		TYPE_FLOAT: {
			TYPE_FLOAT: {ResultType: TypeBool(), OpCode: LTE_FLOAT},
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
