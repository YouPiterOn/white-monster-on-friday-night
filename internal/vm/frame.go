package vm

import "youpiteron.dev/white-monster-on-friday-night/internal/compiler"

type Frame struct {
	closure   *compiler.Closure
	locals    []compiler.Value
	registers []compiler.Value
}

func NewFrame(closure *compiler.Closure) *Frame {
	return &Frame{closure: closure, locals: make([]compiler.Value, closure.Proto.NumLocals), registers: make([]compiler.Value, 0)}
}

func (f *Frame) GetLocal(slot int) *compiler.Value {
	return &f.locals[slot]
}

func (f *Frame) GetRegister(slot int) *compiler.Value {
	return &f.registers[slot]
}

func (f *Frame) GetUpvar(slot int) *compiler.Value {
	if slot >= len(f.closure.Upvalues) {
		panic("VM ERROR: upvar slot out of bounds")
	}
	return f.closure.Upvalues[slot].Ptr
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
	*f.closure.Upvalues[slot].Ptr = value
}

func (f *Frame) GetConstant(index int) compiler.Value {
	return f.closure.Proto.Constants[index]
}
