package compiler

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

func TypeInt() Type {
	return Type{Type: TYPE_INT}
}

func TypeBool() Type {
	return Type{Type: TYPE_BOOL}
}

func TypeNull() Type {
	return Type{Type: TYPE_NULL}
}

func TypeVoid() Type {
	return Type{Type: TYPE_VOID}
}

func TypeClosure() Type {
	return Type{Type: TYPE_CLOSURE}
}

func TypeNativeFunction() Type {
	return Type{Type: TYPE_NATIVE_FUNCTION}
}

func TypeArrayOf(elementType Type) Type {
	return Type{Type: TYPE_ARRAY, ElementType: &elementType}
}

func TypeFromTypeSubkind(typeSubkind lexer.TypeSubkind) Type {
	switch typeSubkind {
	case lexer.TypeInt:
		return TypeInt()
	case lexer.TypeBool:
		return TypeBool()
	case lexer.TypeNull:
		return TypeNull()
	default:
		panic(fmt.Sprintf("COMPILER ERROR: invalid type subkind %s", typeSubkind.String()))
	}
}

func (t Type) DefaultValue() Value {
	switch t.Type {
	case TYPE_INT:
		return NewIntValue(0)
	case TYPE_BOOL:
		return NewBoolValue(false)
	case TYPE_NULL:
		return NewNullValue()
	}
	panic(fmt.Sprintf("COMPILER ERROR: invalid type %v", t))
}

func (t Type) IsEqual(other Type) bool {
	if t.Type != other.Type {
		return false
	}
	if t.ElementType == nil {
		return true
	}
	return t.ElementType.IsEqual(*other.ElementType)
}

func (t Type) String() string {
	if t.ElementType == nil {
		return t.Type.String()
	}
	return fmt.Sprintf("%s<%s>", t.Type.String(), t.ElementType.String())
}
