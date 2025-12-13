package compiler

type ValueKind int

const (
	VAL_INT ValueKind = iota
	VAL_BOOL
	VAL_CLOSURE
	VAL_NULL
)

type Value struct {
	Kind    ValueKind
	Int     int
	Bool    bool
	Closure Closure
}

func NewIntValue(value int) Value {
	return Value{Kind: VAL_INT, Int: value}
}

func NewBoolValue(value bool) Value {
	return Value{Kind: VAL_BOOL, Bool: value}
}

func NewNullValue() Value {
	return Value{Kind: VAL_NULL}
}

func NewClosureValue(proto *FunctionProto) Value {
	return Value{Kind: VAL_CLOSURE, Closure: *NewClosure(proto)}
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
