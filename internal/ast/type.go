package ast

import (
	"fmt"

	"youpiteron.dev/white-monster-on-friday-night/internal/lexer"
)

type TypeEnum int

const (
	TYPE_INT TypeEnum = iota
	TYPE_BOOL
	TYPE_NULL
	TYPE_VOID
	TYPE_CLOSURE
	TYPE_NATIVE_FUNCTION
	TYPE_ARRAY
)

func (t TypeEnum) String() string {
	return [...]string{
		"int",
		"bool",
		"null",
		"void",
		"closure",
		"native_function",
		"array",
	}[t]
}

type Type struct {
	Type        TypeEnum
	ElementType *Type
}

var primitiveTypes = map[TypeEnum]*Type{
	TYPE_INT:             {Type: TYPE_INT},
	TYPE_BOOL:            {Type: TYPE_BOOL},
	TYPE_NULL:            {Type: TYPE_NULL},
	TYPE_VOID:            {Type: TYPE_VOID},
	TYPE_CLOSURE:         {Type: TYPE_CLOSURE},
	TYPE_NATIVE_FUNCTION: {Type: TYPE_NATIVE_FUNCTION},
}

func TypeInt() *Type {
	return primitiveTypes[TYPE_INT]
}

func TypeBool() *Type {
	return primitiveTypes[TYPE_BOOL]
}

func TypeNull() *Type {
	return primitiveTypes[TYPE_NULL]
}

func TypeVoid() *Type {
	return primitiveTypes[TYPE_VOID]
}

func TypeClosure() *Type {
	return primitiveTypes[TYPE_CLOSURE]
}

func TypeNativeFunction() *Type {
	return primitiveTypes[TYPE_NATIVE_FUNCTION]
}

func TypeArrayOf(elementType *Type) *Type {
	return &Type{Type: TYPE_ARRAY, ElementType: elementType}
}

func TypeFromTypeSubkind(typeSubkind lexer.TypeSubkind) *Type {
	switch typeSubkind {
	case lexer.TypeInt:
		return TypeInt()
	case lexer.TypeBool:
		return TypeBool()
	case lexer.TypeNull:
		return TypeNull()
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
