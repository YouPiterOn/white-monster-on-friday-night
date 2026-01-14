package compiler

import (
	"youpiteron.dev/white-monster-on-friday-night/internal/ast"
	"youpiteron.dev/white-monster-on-friday-night/internal/lexer"
)

type BinaryOpInfo struct {
	ResultType *ast.Type
	OpCode     OpCode
}

var BinaryOpTable = map[lexer.OperatorSubkind]map[ast.TypeEnum]map[ast.TypeEnum]BinaryOpInfo{
	lexer.OperatorPlus: {
		ast.TYPE_INT: {
			ast.TYPE_INT: {ResultType: ast.TypeInt(), OpCode: ADD_INT},
		},
		ast.TYPE_FLOAT: {
			ast.TYPE_FLOAT: {ResultType: ast.TypeFloat(), OpCode: ADD_FLOAT},
		},
		ast.TYPE_STRING: {
			ast.TYPE_STRING: {ResultType: ast.TypeString(), OpCode: ADD_STRING},
		},
	},
	lexer.OperatorMinus: {
		ast.TYPE_INT: {
			ast.TYPE_INT: {ResultType: ast.TypeInt(), OpCode: SUB_INT},
		},
		ast.TYPE_FLOAT: {
			ast.TYPE_FLOAT: {ResultType: ast.TypeFloat(), OpCode: SUB_FLOAT},
		},
	},
	lexer.OperatorStar: {
		ast.TYPE_INT: {
			ast.TYPE_INT: {ResultType: ast.TypeInt(), OpCode: MUL_INT},
		},
		ast.TYPE_FLOAT: {
			ast.TYPE_FLOAT: {ResultType: ast.TypeFloat(), OpCode: MUL_FLOAT},
		},
	},
	lexer.OperatorSlash: {
		ast.TYPE_INT: {
			ast.TYPE_INT: {ResultType: ast.TypeInt(), OpCode: DIV_INT},
		},
		ast.TYPE_FLOAT: {
			ast.TYPE_FLOAT: {ResultType: ast.TypeFloat(), OpCode: DIV_FLOAT},
		},
	},
	lexer.OperatorEqual: {
		ast.TYPE_INT: {
			ast.TYPE_INT: {ResultType: ast.TypeBool(), OpCode: EQ_INT},
		},
		ast.TYPE_BOOL: {
			ast.TYPE_BOOL: {ResultType: ast.TypeBool(), OpCode: EQ_BOOL},
		},
		ast.TYPE_FLOAT: {
			ast.TYPE_FLOAT: {ResultType: ast.TypeBool(), OpCode: EQ_FLOAT},
		},
		ast.TYPE_STRING: {
			ast.TYPE_STRING: {ResultType: ast.TypeBool(), OpCode: EQ_STRING},
		},
	},
	lexer.OperatorNotEqual: {
		ast.TYPE_INT: {
			ast.TYPE_INT: {ResultType: ast.TypeBool(), OpCode: NE_INT},
		},
		ast.TYPE_BOOL: {
			ast.TYPE_BOOL: {ResultType: ast.TypeBool(), OpCode: NE_BOOL},
		},
		ast.TYPE_FLOAT: {
			ast.TYPE_FLOAT: {ResultType: ast.TypeBool(), OpCode: NE_FLOAT},
		},
		ast.TYPE_STRING: {
			ast.TYPE_STRING: {ResultType: ast.TypeBool(), OpCode: NE_STRING},
		},
	},
	lexer.OperatorGreater: {
		ast.TYPE_INT: {
			ast.TYPE_INT: {ResultType: ast.TypeBool(), OpCode: GT_INT},
		},
		ast.TYPE_FLOAT: {
			ast.TYPE_FLOAT: {ResultType: ast.TypeBool(), OpCode: GT_FLOAT},
		},
	},
	lexer.OperatorGreaterEqual: {
		ast.TYPE_INT: {
			ast.TYPE_INT: {ResultType: ast.TypeBool(), OpCode: GTE_INT},
		},
		ast.TYPE_FLOAT: {
			ast.TYPE_FLOAT: {ResultType: ast.TypeBool(), OpCode: GTE_FLOAT},
		},
	},
	lexer.OperatorLess: {
		ast.TYPE_INT: {
			ast.TYPE_INT: {ResultType: ast.TypeBool(), OpCode: LT_INT},
		},
		ast.TYPE_FLOAT: {
			ast.TYPE_FLOAT: {ResultType: ast.TypeBool(), OpCode: LT_FLOAT},
		},
	},
	lexer.OperatorLessEqual: {
		ast.TYPE_INT: {
			ast.TYPE_INT: {ResultType: ast.TypeBool(), OpCode: LTE_INT},
		},
		ast.TYPE_FLOAT: {
			ast.TYPE_FLOAT: {ResultType: ast.TypeBool(), OpCode: LTE_FLOAT},
		},
	},
	lexer.OperatorAnd: {
		ast.TYPE_BOOL: {
			ast.TYPE_BOOL: {ResultType: ast.TypeBool(), OpCode: AND_BOOL},
		},
	},
	lexer.OperatorOr: {
		ast.TYPE_BOOL: {
			ast.TYPE_BOOL: {ResultType: ast.TypeBool(), OpCode: OR_BOOL},
		},
	},
}

func ResolveBinaryOp(
	op lexer.OperatorSubkind,
	left *ast.Type,
	right *ast.Type,
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
