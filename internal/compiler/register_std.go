package compiler

import "youpiteron.dev/white-monster-on-friday-night/internal/ast"

func RegisterStdGlobals(gt *GlobalTable) {
	gt.DefineFunctionVariable(
		"println",
		false,
		ast.TypeNativeFunction(),
		&FuncSignature{
			CallArgs:   []*ast.Type{ast.TypeInt()},
			ReturnType: ast.TypeNull(),
			Vararg:     false,
		},
	)
}
