package compiler

import (
	"fmt"

	"youpiteron.dev/white-monster-on-friday-night/internal/api"
	"youpiteron.dev/white-monster-on-friday-night/internal/lexer"
)

type ValueType int

const (
	VAL_INT ValueType = iota
	VAL_BOOL
	VAL_CLOSURE
	VAL_NULL
	VAL_NATIVE_FUNCTION
)

func (t ValueType) String() string {
	return [...]string{
		"INT",
		"BOOL",
		"CLOSURE",
		"NULL",
		"NATIVE_FUNCTION",
	}[t]
}

func ValueTypeFromTypeSubkind(typeSubkind lexer.TypeSubkind) ValueType {
	switch typeSubkind {
	case lexer.TypeInt:
		return VAL_INT
	case lexer.TypeBool:
		return VAL_BOOL
	case lexer.TypeNull:
		return VAL_NULL
	}
	panic(fmt.Sprintf("COMPILER ERROR: invalid type subkind %v", typeSubkind))
}

func (t ValueType) DefaultValue() Value {
	switch t {
	case VAL_INT:
		return NewIntValue(0)
	case VAL_BOOL:
		return NewBoolValue(false)
	case VAL_NULL:
		return NewNullValue()
	}
	panic(fmt.Sprintf("COMPILER ERROR: invalid type %v", t))
}

type Value struct {
	TypeOf  ValueType
	Int     int
	Bool    bool
	Closure Closure
	Native  NativeFunction
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

type UpvalueCell struct {
	Ptr *Value
}

type Closure struct {
	Proto    *FunctionProto
	Upvalues []*UpvalueCell
}

func NewClosure(proto *FunctionProto) *Closure {
	return &Closure{Proto: proto, Upvalues: make([]*UpvalueCell, len(proto.Upvars))}
}

type NativeFunction func(vm api.VM, args ...Value) (Value, error)
