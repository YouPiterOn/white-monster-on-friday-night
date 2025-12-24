package compiler

import (
	"fmt"

	"youpiteron.dev/white-monster-on-friday-night/internal/api"
	"youpiteron.dev/white-monster-on-friday-night/internal/ast"
)

type ValueType int

const (
	VAL_INT ValueType = iota
	VAL_BOOL
	VAL_CLOSURE
	VAL_NULL
	VAL_NATIVE_FUNCTION
	VAL_ARRAY
)

func (t ValueType) String() string {
	return [...]string{
		"INT",
		"BOOL",
		"CLOSURE",
		"NULL",
		"NATIVE_FUNCTION",
		"ARRAY",
	}[t]
}

type Value struct {
	TypeOf  ValueType
	Int     int
	Bool    bool
	Closure Closure
	Native  NativeFunction
	Array   []Value
}

func NewIntValue(value int) Value {
	return Value{TypeOf: VAL_INT, Int: value}
}

func NewBoolValue(value bool) Value {
	return Value{TypeOf: VAL_BOOL, Bool: value}
}

func NewNullValue() Value {
	return Value{TypeOf: VAL_NULL}
}

func NewClosureValue(proto *FunctionProto) Value {
	return Value{TypeOf: VAL_CLOSURE, Closure: *NewClosure(proto)}
}

func NewNativeFunctionValue(function NativeFunction) Value {
	return Value{TypeOf: VAL_NATIVE_FUNCTION, Native: function}
}

func NewArrayValue(elements []Value) Value {
	return Value{TypeOf: VAL_ARRAY, Array: elements}
}

func DefaultValue(typeOf *ast.Type) Value {
	switch typeOf.Type {
	case ast.TYPE_INT:
		return NewIntValue(0)
	case ast.TYPE_BOOL:
		return NewBoolValue(false)
	case ast.TYPE_NULL:
		return NewNullValue()
	case ast.TYPE_ARRAY:
		return NewArrayValue([]Value{})
	}
	panic(fmt.Sprintf("invalid type %v", typeOf))
}

type UpvalueCell struct {
	Ptr *Value
}

type Closure struct {
	Proto    Proto
	Upvalues []*UpvalueCell
}

func NewClosure(proto Proto) *Closure {
	return &Closure{Proto: proto, Upvalues: make([]*UpvalueCell, proto.NumLocals())}
}

type NativeFunction func(vm api.VM, args ...Value) (Value, error)
