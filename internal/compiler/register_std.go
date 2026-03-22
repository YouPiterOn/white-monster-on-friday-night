package compiler

import "youpiteron.dev/white-monster-on-friday-night/internal/ast"

func RegisterStdGlobals(gt *GlobalTable) {
	gt.DefineVariable(
		"println",
		false,
		ast.TypeFunction(&ast.FuncSignature{
			CallArgs:   []*ast.Type{ast.TypeArrayOf(ast.TypeAny())},
			ReturnType: ast.TypeVoid(),
			Vararg:     true,
		}),
	)
	gt.DefineVariable(
		"append",
		false,
		ast.TypeFunction(&ast.FuncSignature{
			CallArgs:   []*ast.Type{ast.TypeArrayOf(ast.TypeInt()), ast.TypeInt()},
			ReturnType: ast.TypeArrayOf(ast.TypeInt()),
			Vararg:     false,
		}),
	)

	gt.DefineType("Int", TypeInt())
	gt.DefineType("Float", TypeFloat())
	gt.DefineType("String", TypeString())
	gt.DefineType("Bool", TypeBool())
	gt.DefineType("null", TypeNull())
	gt.DefineType("void", TypeVoid())
}
