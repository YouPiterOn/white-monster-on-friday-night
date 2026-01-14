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

type Type struct {
	Type        TypeEnum
	ElementType *Type
}

var primitiveTypes = map[TypeEnum]*Type{
	TYPE_ANY:             {Type: TYPE_ANY},
	TYPE_INT:             {Type: TYPE_INT},
	TYPE_FLOAT:           {Type: TYPE_FLOAT},
	TYPE_STRING:          {Type: TYPE_STRING},
	TYPE_BOOL:            {Type: TYPE_BOOL},
	TYPE_NULL:            {Type: TYPE_NULL},
	TYPE_VOID:            {Type: TYPE_VOID},
	TYPE_CLOSURE:         {Type: TYPE_CLOSURE},
	TYPE_NATIVE_FUNCTION: {Type: TYPE_NATIVE_FUNCTION},
}

func TypeAny() *Type {
	return primitiveTypes[TYPE_ANY]
}

func TypeInt() *Type {
	return primitiveTypes[TYPE_INT]
}

func TypeFloat() *Type {
	return primitiveTypes[TYPE_FLOAT]
}

func TypeString() *Type {
	return primitiveTypes[TYPE_STRING]
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
