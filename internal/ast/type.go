package ast

import (
	"fmt"

	"youpiteron.dev/white-monster-on-friday-night/internal/lexer"
)

type TypeEnum int

const (
	TYPE_ANY TypeEnum = iota
	TYPE_INT
	TYPE_FLOAT
	TYPE_STRING
	TYPE_BOOL
	TYPE_NULL
	TYPE_VOID
	TYPE_CLOSURE
	TYPE_NATIVE_FUNCTION
	TYPE_ARRAY
)

func (t TypeEnum) String() string {
	return [...]string{
		"any",
		"int",
		"float",
		"string",
		"bool",
		"null",
		"void",
		"closure",
		"native_function",
		"array",
	}[t]
}

type FuncSignature struct {
	CallArgs   []*Type
	ReturnType *Type
	Vararg     bool
}
type Type struct {
	Type          TypeEnum
	Proto         *Type
	ElementType   *Type
	Properties    map[string]*Type
	FuncSignature *FuncSignature
}

func TypeAny() *Type {
	return &Type{Type: TYPE_ANY}
}

func TypeInt() *Type {
	return &Type{Type: TYPE_INT}
}

func TypeFloat() *Type {
	return &Type{Type: TYPE_FLOAT}
}

func TypeString() *Type {
	return &Type{Type: TYPE_STRING}
}

func TypeBool() *Type {
	return &Type{Type: TYPE_BOOL}
}

func TypeNull() *Type {
	return &Type{Type: TYPE_NULL}
}

func TypeVoid() *Type {
	return &Type{Type: TYPE_VOID}
}

func TypeClosure(funcSignature *FuncSignature) *Type {
	return &Type{Type: TYPE_CLOSURE, FuncSignature: funcSignature}
}

func TypeNativeFunction(funcSignature *FuncSignature) *Type {
	return &Type{Type: TYPE_NATIVE_FUNCTION, FuncSignature: funcSignature}
}

func TypeArrayOf(elementType *Type) *Type {
	return &Type{Type: TYPE_ARRAY, ElementType: elementType}
}

func TypeFromTypeSubkind(typeSubkind lexer.TypeSubkind) *Type {
	switch typeSubkind {
	case lexer.TypeInt:
		return TypeInt()
	case lexer.TypeFloat:
		return TypeFloat()
	case lexer.TypeString:
		return TypeString()
	case lexer.TypeBool:
		return TypeBool()
	case lexer.TypeNull:
		return TypeNull()
	case lexer.TypeVoid:
		return TypeVoid()
	default:
		panic(fmt.Sprintf("invalid type subkind %s", typeSubkind.String()))
	}
}

func (t *Type) IsEqual(other *Type) bool {
	if t.Type != other.Type {
		return false
	}
	if t.ElementType == nil {
		return true
	}
	return t.ElementType.IsEqual(other.ElementType)
}

func (t *Type) String() string {
	if t.ElementType == nil {
		return t.Type.String()
	}
	return fmt.Sprintf("[]%s", t.ElementType.String())
}
