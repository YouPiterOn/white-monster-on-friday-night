package compiler

func RegisterStdGlobals(gt *GlobalTable) {
	gt.DefineFunctionVariable(
		"println",
		false,
		VAL_NATIVE_FUNCTION,
		&FuncSignature{
			CallArgs:   []ValueType{VAL_INT},
			ReturnType: VAL_NULL,
			Vararg:     false,
		},
	)
}
