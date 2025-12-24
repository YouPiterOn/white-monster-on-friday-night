package compiler

func RegisterStdGlobals(gt *GlobalTable) {
	gt.DefineFunctionVariable(
		"println",
		false,
		TypeNativeFunction(),
		&FuncSignature{
			CallArgs:   []Type{TypeInt()},
			ReturnType: TypeNull(),
			Vararg:     false,
		},
	)
}
