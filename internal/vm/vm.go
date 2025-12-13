package vm

import (
	"youpiteron.dev/white-monster-on-friday-night/internal/compiler"
)

type VM struct {
	frames         []Frame
	functionProtos []compiler.FunctionProto
}

func NewVM(functionProtos []compiler.FunctionProto) *VM {
	return &VM{frames: make([]Frame, 0), functionProtos: functionProtos}
}

func (v *VM) Run() compiler.Value {
	if len(v.functionProtos) == 0 {
		panic("VM ERROR: no functions to run")
	}
	functionProto := &v.functionProtos[len(v.functionProtos)-1]
	closure := compiler.NewClosure(functionProto)
	frame := NewFrame(closure)
	v.frames = append(v.frames, *frame)
	retval := v.runFunction(functionProto)
	return *retval
}

func (v *VM) currentFrame() *Frame {
	return &v.frames[len(v.frames)-1]
}

func (v *VM) runFunction(functionProto *compiler.FunctionProto) *compiler.Value {
	var retval *compiler.Value = nil
	for _, instruction := range functionProto.Instructions {
		if retval != nil {
			break
		}
		switch instruction.OpCode {
		case compiler.LOAD_CONST:
			v.opLoadConst(instruction.Args)
		case compiler.LOAD_VAR:
			v.opLoadVar(instruction.Args)
		case compiler.LOAD_UPVAR:
			v.opLoadUpvar(instruction.Args)
		case compiler.STORE_VAR:
			v.opStoreVar(instruction.Args)
		case compiler.ASSIGN_UPVAR:
			v.opAssignUpvar(instruction.Args)
		case compiler.ADD:
			v.opAdd(instruction.Args)
		case compiler.SUB:
			v.opSub(instruction.Args)
		case compiler.MUL:
			v.opMul(instruction.Args)
		case compiler.DIV:
			v.opDiv(instruction.Args)
		case compiler.CLOSURE:
			v.opClosure(instruction.Args)
		case compiler.CALL:
			v.opCall(instruction.Args)
		case compiler.RETURN:
			retval = v.opReturn(instruction.Args)
		}
	}
	v.frames = v.frames[:len(v.frames)-1]
	return retval
}

func (v *VM) opLoadConst(args []int) {
	value := v.currentFrame().GetConstant(args[1])
	v.currentFrame().SetRegister(args[0], value)
}

func (v *VM) opLoadVar(args []int) {
	value := v.currentFrame().GetLocal(args[1])
	v.currentFrame().SetRegister(args[0], *value)
}

func (v *VM) opLoadUpvar(args []int) {
	value := v.currentFrame().GetUpvar(args[1])
	v.currentFrame().SetRegister(args[0], *value)
}

func (v *VM) opStoreVar(args []int) {
	value := v.currentFrame().GetRegister(args[0])
	v.currentFrame().SetLocal(args[1], *value)
}

func (v *VM) opAssignUpvar(args []int) {
	value := v.currentFrame().GetRegister(args[0])
	v.currentFrame().SetUpvar(args[1], *value)
}

func (v *VM) opAdd(args []int) {
	left := v.currentFrame().GetRegister(args[1])
	right := v.currentFrame().GetRegister(args[2])
	if left.Kind != compiler.VAL_INT || right.Kind != compiler.VAL_INT {
		panic("VM ERROR: invalid operand type for addition")
	}
	result := compiler.Value{Kind: compiler.VAL_INT, Int: left.Int + right.Int}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opSub(args []int) {
	left := v.currentFrame().GetRegister(args[1])
	right := v.currentFrame().GetRegister(args[2])
	if left.Kind != compiler.VAL_INT || right.Kind != compiler.VAL_INT {
		panic("VM ERROR: invalid operand type for subtraction")
	}
	result := compiler.Value{Kind: compiler.VAL_INT, Int: left.Int - right.Int}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opMul(args []int) {
	left := v.currentFrame().GetRegister(args[1])
	right := v.currentFrame().GetRegister(args[2])
	if left.Kind != compiler.VAL_INT || right.Kind != compiler.VAL_INT {
		panic("VM ERROR: invalid operand type for multiplication")
	}
	result := compiler.Value{Kind: compiler.VAL_INT, Int: left.Int * right.Int}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opDiv(args []int) {
	left := v.currentFrame().GetRegister(args[1])
	right := v.currentFrame().GetRegister(args[2])
	if left.Kind != compiler.VAL_INT || right.Kind != compiler.VAL_INT {
		panic("VM ERROR: invalid operand type for division")
	}
	result := compiler.Value{Kind: compiler.VAL_INT, Int: left.Int / right.Int}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opClosure(args []int) {
	proto := v.functionProtos[args[1]]
	closure := &compiler.Closure{Proto: &proto, Upvalues: make([]*compiler.UpvalueCell, len(proto.Upvars))}
	for i, upvar := range proto.Upvars {
		if upvar.IsFromParent {
			closure.Upvalues[i] = &compiler.UpvalueCell{Ptr: v.currentFrame().GetLocal(upvar.SlotInParent)}
		} else {
			closure.Upvalues[i] = &compiler.UpvalueCell{Ptr: v.currentFrame().GetUpvar(upvar.SlotInParent)}
		}
	}
	value := compiler.Value{Kind: compiler.VAL_CLOSURE, Closure: *closure}
	v.currentFrame().SetRegister(args[0], value)
}

func (v *VM) opCall(args []int) {
	function := v.currentFrame().GetRegister(args[1])
	if function.Kind != compiler.VAL_CLOSURE {
		panic("VM ERROR: invalid operand type for call")
	}
	frame := NewFrame(&function.Closure)
	for i, argument := range args[2:] {
		value := v.currentFrame().GetRegister(argument)
		frame.SetLocal(i, *value)
	}
	v.frames = append(v.frames, *frame)
	functionProto := function.Closure.Proto
	retval := v.runFunction(functionProto)
	v.currentFrame().SetRegister(args[0], *retval)
}

func (v *VM) opReturn(args []int) *compiler.Value {
	return v.currentFrame().GetRegister(args[0])
}
