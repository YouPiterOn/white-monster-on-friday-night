package compiler

import "youpiteron.dev/white-monster-on-friday-night/internal/ast"

func RegisterStdGlobals(gt *GlobalTable) {
	gt.DefineVariable(
		"println",
		false,
		ast.TypeNativeFunction(&ast.FuncSignature{
			CallArgs:   []*ast.Type{ast.TypeArrayOf(ast.TypeAny())},
			ReturnType: ast.TypeVoid(),
			Vararg:     true,
		}),
	)
	gt.DefineVariable(
		"append",
		false,
		ast.TypeNativeFunction(&ast.FuncSignature{
			CallArgs:   []*ast.Type{ast.TypeArrayOf(ast.TypeInt()), ast.TypeInt()},
			ReturnType: ast.TypeArrayOf(ast.TypeInt()),
			Vararg:     false,
		}),
	)
}
