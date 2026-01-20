package ast

import (
	"fmt"

	"youpiteron.dev/white-monster-on-friday-night/internal/common"
	"youpiteron.dev/white-monster-on-friday-night/internal/lexer"
)

type TypeDeclaration struct {
	Name       string
	Properties []TypeProperty
	PosAt      *common.SourcePos
}

func (t *TypeDeclaration) Pos() *common.SourcePos { return t.PosAt }
func (t *TypeDeclaration) statementNode()         {}
func (t *TypeDeclaration) Visit(v Visitor[any]) any {
	return v.VisitTypeDeclaration(t)
}

type TypeProperty struct {
	Name   string
	TypeOf TypeExpression
	PosAt  *common.SourcePos
}

type TypeEnum int

const (
	TYPE_ANY TypeEnum = iota
	TYPE_INT
	TYPE_FLOAT
	TYPE_STRING
	TYPE_BOOL
	TYPE_NULL
	TYPE_VOID
	TYPE_FUNCTION
	TYPE_ARRAY
)

func (t TypeEnum) String() string {
	return [...]string{
		"any",
		"Int",
		"Float",
		"String",
		"Bool",
		"null",
		"void",
		"Function",
		"Array",
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

func TypeFunction(funcSignature *FuncSignature) *Type {
	return &Type{Type: TYPE_FUNCTION, FuncSignature: funcSignature}
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
