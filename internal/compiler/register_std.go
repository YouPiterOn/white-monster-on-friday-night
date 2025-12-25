package compiler

import "youpiteron.dev/white-monster-on-friday-night/internal/ast"

func RegisterStdGlobals(gt *GlobalTable) {
	gt.DefineFunctionVariable(
		"println",
		false,
		ast.TypeNativeFunction(),
		&FuncSignature{
			CallArgs:   []*ast.Type{ast.TypeArrayOf(ast.TypeInt())},
			ReturnType: ast.TypeNull(),
			Vararg:     true,
		},
	)
	gt.DefineFunctionVariable(
		"append",
		false,
		ast.TypeNativeFunction(),
		&FuncSignature{
			CallArgs:   []*ast.Type{ast.TypeArrayOf(ast.TypeInt()), ast.TypeInt()},
			ReturnType: ast.TypeArrayOf(ast.TypeInt()),
			Vararg:     false,
		},
	)
}
