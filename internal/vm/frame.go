package vm

type Frame struct {
	closure   *Closure
	locals    []Value
	registers []Value
}

func NewFrame(closure *Closure) *Frame {
	return &Frame{closure: closure, locals: make([]Value, closure.proto.NumLocals), registers: make([]Value, 0)}
}

func (f *Frame) GetLocal(slot int) *Value {
	return &f.locals[slot]
}

func (f *Frame) GetRegister(slot int) *Value {
	return &f.registers[slot]
}

func (f *Frame) GetUpvar(slot int) *Value {
	if slot >= len(f.closure.upvalues) {
		panic("VM ERROR: upvar slot out of bounds")
	}
	return f.closure.upvalues[slot].Ptr
}

func (f *Frame) SetLocal(slot int, value Value) {
	if slot >= len(f.locals) {
		newLocals := make([]Value, (slot+1)*2)
		copy(newLocals, f.locals)
		f.locals = newLocals
	}
	f.locals[slot] = value
}

func (f *Frame) SetRegister(slot int, value Value) {
	if slot >= len(f.registers) {
		newRegisters := make([]Value, (slot+1)*2)
		copy(newRegisters, f.registers)
		f.registers = newRegisters
	}
	f.registers[slot] = value
}

func (f *Frame) SetUpvar(slot int, value Value) {
	*f.closure.upvalues[slot].Ptr = value
}
