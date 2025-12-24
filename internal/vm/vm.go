package vm

import (
	"fmt"

	"youpiteron.dev/white-monster-on-friday-night/internal/compiler"
	"youpiteron.dev/white-monster-on-friday-night/internal/native"
)

type VM struct {
	frames         []Frame
	moduleInstance *ModuleInstance
	globals        []compiler.Value
}

func (v *VM) ImplementVMInterface() {}

func NewVM(gt *compiler.GlobalTable) *VM {
	vm := &VM{frames: make([]Frame, 0), moduleInstance: nil, globals: make([]compiler.Value, gt.Length())}
	vm.initStdlibValues(gt)
	return vm
}

func (v *VM) RunModuleProto(moduleProto *compiler.ModuleProto) int {
	v.moduleInstance = NewModuleInstance(moduleProto)
	if len(v.frames) == 0 {
		frame := NewFrame(moduleProto, make([]*compiler.UpvalueCell, 0))
		v.frames = append(v.frames, *frame)
	} else {
		v.currentFrame().SetConstants(moduleProto.Constants())
	}
	retval := v.runInstructions(moduleProto.Instructions())
	if retval == nil {
		return 0
	}
	return retval.Int
}

func (v *VM) initStdlibValues(gt *compiler.GlobalTable) {
	// println
	if variable, ok := gt.FindVariable("println"); ok {
		v.globals[variable.Slot] = compiler.Value{
			TypeOf: compiler.VAL_NATIVE_FUNCTION,
			Native: native.Println,
		}
	}
	// append
	if variable, ok := gt.FindVariable("append"); ok {
		v.globals[variable.Slot] = compiler.Value{
			TypeOf: compiler.VAL_NATIVE_FUNCTION,
			Native: native.Append,
		}
	}
}

func (v *VM) currentFrame() *Frame {
	return &v.frames[len(v.frames)-1]
}

func (v *VM) runInstructions(instructions []compiler.Instruction) *compiler.Value {
	for {
		frame := v.currentFrame()
		if frame.ip >= len(instructions) || frame.retval != nil {
			break
		}
		instruction := instructions[frame.ip]
		switch instruction.OpCode {
		case compiler.LOAD_CONST:
			v.opLoadConst(instruction.Args)
		case compiler.LOAD_VAR:
			v.opLoadVar(instruction.Args)
		case compiler.LOAD_GLOBAL:
			v.opLoadGlobal(instruction.Args)
		case compiler.LOAD_UPVAR:
			v.opLoadUpvar(instruction.Args)
		case compiler.STORE_VAR:
			v.opStoreVar(instruction.Args)
		case compiler.ASSIGN_GLOBAL:
			v.opAssignGlobal(instruction.Args)
		case compiler.ASSIGN_UPVAR:
			v.opAssignUpvar(instruction.Args)
		case compiler.ADD_INT:
			v.opAddInt(instruction.Args)
		case compiler.SUB_INT:
			v.opSubInt(instruction.Args)
		case compiler.MUL_INT:
			v.opMulInt(instruction.Args)
		case compiler.DIV_INT:
			v.opDivInt(instruction.Args)
		case compiler.EQ_INT:
			v.opEqInt(instruction.Args)
		case compiler.EQ_BOOL:
			v.opEqBool(instruction.Args)
		case compiler.NE_INT:
			v.opNeInt(instruction.Args)
		case compiler.NE_BOOL:
			v.opNeBool(instruction.Args)
		case compiler.GT_INT:
			v.opGtInt(instruction.Args)
		case compiler.GTE_INT:
			v.opGteInt(instruction.Args)
		case compiler.LT_INT:
			v.opLtInt(instruction.Args)
		case compiler.LTE_INT:
			v.opLteInt(instruction.Args)
		case compiler.AND_BOOL:
			v.opAndBool(instruction.Args)
		case compiler.OR_BOOL:
			v.opOrBool(instruction.Args)
		case compiler.CLOSURE:
			v.opClosure(instruction.Args)
		case compiler.CALL:
			v.opCall(instruction.Args)
		case compiler.RETURN:
			v.opReturn(instruction.Args)
		case compiler.JUMP_IF_FALSE:
			v.opJumpIfFalse(instruction.Args)
		case compiler.JUMP:
			v.opJump(instruction.Args)
		case compiler.MAKE_ARRAY:
			v.opMakeArray(instruction.Args)
		case compiler.INDEX_ARRAY:
			v.opIndexArray(instruction.Args)
		}
		frame.AdvanceIp()
	}
	var retval *compiler.Value = &compiler.Value{TypeOf: compiler.VAL_NULL}
	if v.currentFrame().retval != nil {
		*retval = *v.currentFrame().retval
		v.currentFrame().retval = nil
	}
	v.currentFrame().SetIp(0)
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

func (v *VM) opLoadGlobal(args []int) {
	value := v.globals[args[1]]
	v.currentFrame().SetRegister(args[0], value)
}

func (v *VM) opLoadUpvar(args []int) {
	value := v.currentFrame().GetUpvar(args[1])
	v.currentFrame().SetRegister(args[0], *value)
}

func (v *VM) opStoreVar(args []int) {
	value := v.currentFrame().GetRegister(args[0])
	v.currentFrame().SetLocal(args[1], *value)
}

func (v *VM) opAssignGlobal(args []int) {
	value := v.currentFrame().GetRegister(args[0])
	v.globals[args[1]] = *value
}

func (v *VM) opAssignUpvar(args []int) {
	value := v.currentFrame().GetRegister(args[0])
	v.currentFrame().SetUpvar(args[1], *value)
}

func (v *VM) opAddInt(args []int) {
	left := v.currentFrame().GetRegister(args[1])
	right := v.currentFrame().GetRegister(args[2])
	result := compiler.Value{TypeOf: compiler.VAL_INT, Int: left.Int + right.Int}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opSubInt(args []int) {
	left := v.currentFrame().GetRegister(args[1])
	right := v.currentFrame().GetRegister(args[2])
	result := compiler.Value{TypeOf: compiler.VAL_INT, Int: left.Int - right.Int}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opMulInt(args []int) {
	left := v.currentFrame().GetRegister(args[1])
	right := v.currentFrame().GetRegister(args[2])
	result := compiler.Value{TypeOf: compiler.VAL_INT, Int: left.Int * right.Int}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opDivInt(args []int) {
	left := v.currentFrame().GetRegister(args[1])
	right := v.currentFrame().GetRegister(args[2])
	result := compiler.Value{TypeOf: compiler.VAL_INT, Int: left.Int / right.Int}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opEqInt(args []int) {
	left := v.currentFrame().GetRegister(args[1])
	right := v.currentFrame().GetRegister(args[2])
	result := compiler.Value{TypeOf: compiler.VAL_BOOL, Bool: left.Int == right.Int}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opEqBool(args []int) {
	left := v.currentFrame().GetRegister(args[1])
	right := v.currentFrame().GetRegister(args[2])
	result := compiler.Value{TypeOf: compiler.VAL_BOOL, Bool: left.Bool == right.Bool}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opNeInt(args []int) {
	left := v.currentFrame().GetRegister(args[1])
	right := v.currentFrame().GetRegister(args[2])
	result := compiler.Value{TypeOf: compiler.VAL_BOOL, Bool: left.Int != right.Int}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opNeBool(args []int) {
	left := v.currentFrame().GetRegister(args[1])
	right := v.currentFrame().GetRegister(args[2])
	result := compiler.Value{TypeOf: compiler.VAL_BOOL, Bool: left.Bool != right.Bool}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opGtInt(args []int) {
	left := v.currentFrame().GetRegister(args[1])
	right := v.currentFrame().GetRegister(args[2])
	result := compiler.Value{TypeOf: compiler.VAL_BOOL, Bool: left.Int > right.Int}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opGteInt(args []int) {
	left := v.currentFrame().GetRegister(args[1])
	right := v.currentFrame().GetRegister(args[2])
	result := compiler.Value{TypeOf: compiler.VAL_BOOL, Bool: left.Int >= right.Int}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opLtInt(args []int) {
	left := v.currentFrame().GetRegister(args[1])
	right := v.currentFrame().GetRegister(args[2])
	result := compiler.Value{TypeOf: compiler.VAL_BOOL, Bool: left.Int < right.Int}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opLteInt(args []int) {
	left := v.currentFrame().GetRegister(args[1])
	right := v.currentFrame().GetRegister(args[2])
	result := compiler.Value{TypeOf: compiler.VAL_BOOL, Bool: left.Int <= right.Int}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opAndBool(args []int) {
	left := v.currentFrame().GetRegister(args[1])
	right := v.currentFrame().GetRegister(args[2])
	result := compiler.Value{TypeOf: compiler.VAL_BOOL, Bool: left.Bool && right.Bool}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opOrBool(args []int) {
	left := v.currentFrame().GetRegister(args[1])
	right := v.currentFrame().GetRegister(args[2])
	result := compiler.Value{TypeOf: compiler.VAL_BOOL, Bool: left.Bool || right.Bool}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opClosure(args []int) {
	proto := v.moduleInstance.functions[args[1]]
	closure := &compiler.Closure{Proto: &proto, Upvalues: make([]*compiler.UpvalueCell, len(proto.Upvars()))}
	for i, upvar := range proto.Upvars() {
		if upvar.IsFromParent {
			closure.Upvalues[i] = &compiler.UpvalueCell{Ptr: v.currentFrame().GetLocal(upvar.SlotInParent)}
		} else {
			closure.Upvalues[i] = &compiler.UpvalueCell{Ptr: v.currentFrame().GetUpvar(upvar.SlotInParent)}
		}
	}
	value := compiler.Value{TypeOf: compiler.VAL_CLOSURE, Closure: *closure}
	v.currentFrame().SetRegister(args[0], value)
}

func (v *VM) opCall(args []int) {
	function := v.currentFrame().GetRegister(args[1])
	funcArgs := args[2:]
	switch function.TypeOf {
	case compiler.VAL_CLOSURE:
		frame := NewFrame(function.Closure.Proto, function.Closure.Upvalues)
		for i, argument := range funcArgs {
			value := v.currentFrame().GetRegister(argument)
			frame.SetLocal(i, *value)
		}
		v.frames = append(v.frames, *frame)
		functionProto := function.Closure.Proto
		retval := v.runInstructions(functionProto.Instructions())
		v.frames = v.frames[:len(v.frames)-1]
		v.currentFrame().SetRegister(args[0], *retval)
		return
	case compiler.VAL_NATIVE_FUNCTION:
		values := make([]compiler.Value, len(funcArgs))
		for i, argument := range funcArgs {
			values[i] = *v.currentFrame().GetRegister(argument)
		}
		retval, err := function.Native(v, values...)
		if err != nil {
			panic(fmt.Sprintf("VM ERROR: native function %v returned error: %v", function.Native, err))
		}
		v.currentFrame().SetRegister(args[0], retval)
		return
	}
}

func (v *VM) opReturn(args []int) {
	val := v.currentFrame().GetRegister(args[0])
	v.currentFrame().SetRetval(val)
}

func (v *VM) opJumpIfFalse(args []int) {
	condition := v.currentFrame().GetRegister(args[0])
	if !condition.Bool {
		v.currentFrame().SetIp(args[1])
	}
}

func (v *VM) opJump(args []int) {
	v.currentFrame().SetIp(args[0])
}

func (v *VM) opMakeArray(args []int) {
	elements := args[1:]
	values := make([]compiler.Value, len(elements))
	for i, element := range elements {
		values[i] = *v.currentFrame().GetRegister(element)
	}
	result := compiler.Value{TypeOf: compiler.VAL_ARRAY, Array: values}
	v.currentFrame().SetRegister(args[0], result)
}

func (v *VM) opIndexArray(args []int) {
	array := v.currentFrame().GetRegister(args[1])
	index := v.currentFrame().GetRegister(args[2])
	result := array.Array[index.Int]
	v.currentFrame().SetRegister(args[0], result)
}
