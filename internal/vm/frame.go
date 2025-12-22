package vm

import "youpiteron.dev/white-monster-on-friday-night/internal/compiler"

type Frame struct {
	constants []compiler.Value
	upvalues  []*compiler.UpvalueCell
	locals    []compiler.Value
	registers []compiler.Value
	ip        int
	retval    *compiler.Value
}

func NewFrame(proto compiler.Proto, upvalues []*compiler.UpvalueCell) *Frame {
	return &Frame{constants: proto.Constants(), upvalues: upvalues, locals: make([]compiler.Value, proto.NumLocals()), registers: make([]compiler.Value, 0), ip: 0, retval: nil}
}

func (f *Frame) GetLocal(slot int) *compiler.Value {
	return &f.locals[slot]
}

func (f *Frame) GetRegister(slot int) *compiler.Value {
	return &f.registers[slot]
}

func (f *Frame) GetUpvar(slot int) *compiler.Value {
	if slot >= len(f.upvalues) {
		panic("VM ERROR: upvar slot out of bounds")
	}
	return f.upvalues[slot].Ptr
}

func (f *Frame) SetLocal(slot int, value compiler.Value) {
	if slot >= len(f.locals) {
		newLocals := make([]compiler.Value, (slot+1)*2)
		copy(newLocals, f.locals)
		f.locals = newLocals
	}
	f.locals[slot] = value
}

func (f *Frame) SetRegister(slot int, value compiler.Value) {
	if slot >= len(f.registers) {
		newRegisters := make([]compiler.Value, (slot+1)*2)
		copy(newRegisters, f.registers)
		f.registers = newRegisters
	}
	f.registers[slot] = value
}

func (f *Frame) SetUpvar(slot int, value compiler.Value) {
	*f.upvalues[slot].Ptr = value
}

func (f *Frame) SetConstants(constants []compiler.Value) {
	f.constants = constants
}

func (f *Frame) GetConstant(index int) compiler.Value {
	return f.constants[index]
}

func (f *Frame) AdvanceIp() {
	f.ip++
}

func (f *Frame) SetIp(ip int) {
	f.ip = ip
}

func (f *Frame) SetRetval(retval *compiler.Value) {
	f.retval = retval
}
