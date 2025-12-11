package vm

import "youpiteron.dev/white-monster-on-friday-night/internal/compiler"

type ValueKind int

const (
	VAL_INT ValueKind = iota
	VAL_CLOSURE
)

type Value struct {
	Kind    ValueKind
	Int     int
	Closure Closure
}

type UpvalueCell struct {
	Ptr *Value
}

type Closure struct {
	proto    *compiler.FunctionProto
	upvalues []*UpvalueCell
}

func NewClosure(proto *compiler.FunctionProto) *Closure {
	return &Closure{proto: proto, upvalues: make([]*UpvalueCell, len(proto.Upvars))}
}
