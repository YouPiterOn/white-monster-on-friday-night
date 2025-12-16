package compiler

import (
	"fmt"

	"youpiteron.dev/white-monster-on-friday-night/internal/ast"
	"youpiteron.dev/white-monster-on-friday-night/internal/common"
)

type VisitExprResult struct {
	Reg      int
	TypeOf   ValueType
	Callable *Callable
}

func CastVisitExprResult(result any) (*VisitExprResult, bool) {
	if result == nil {
		return nil, false
	}
	resultVisitExpr, ok := result.(*VisitExprResult)
	if !ok {
		panic(fmt.Sprintf("COMPILER ERROR: result %v is not a VisitExprResult", resultVisitExpr))
	}
	return resultVisitExpr, true
}

type InstructionsVisitor struct {
	scope            *Scope
	errors           []common.Error
	warnings         []common.Error
	reg              int
	functionBuilders []FunctionBuilder
	functionProtos   []FunctionProto
}

// ---------- Constructor ----------

func NewInstructionsVisitor() *InstructionsVisitor {
	return &InstructionsVisitor{scope: nil, errors: []common.Error{}, warnings: []common.Error{}, reg: 0, functionBuilders: []FunctionBuilder{}, functionProtos: []FunctionProto{}}
}

// ---------- Helpers ----------

func (v *InstructionsVisitor) addError(message string, pos *common.SourcePos) {
	v.errors = append(v.errors, common.Error{Message: message, Pos: pos})
}

func (v *InstructionsVisitor) addWarning(message string, pos *common.SourcePos) {
	v.warnings = append(v.warnings, common.Error{Message: message, Pos: pos})
}

func (v *InstructionsVisitor) nextReg() int {
	reg := v.reg
	v.reg++
	return reg
}

func (v *InstructionsVisitor) resetReg() int {
	reg := v.reg
	v.reg = 0
	return reg
}

func (v *InstructionsVisitor) addInstruction(instruction Instruction) int {
	functionBuilder := v.currentFunctionBuilder()
	functionBuilder.AddInstruction(instruction)
	return len(functionBuilder.Instructions) - 1
}

func (v *InstructionsVisitor) setInstruction(index int, instruction Instruction) {
	functionBuilder := v.currentFunctionBuilder()
	functionBuilder.Instructions[index] = instruction
}

func (v *InstructionsVisitor) instructionsLength() int {
	functionBuilder := v.currentFunctionBuilder()
	return len(functionBuilder.Instructions)
}

func (v *InstructionsVisitor) addConstant(value Value) int {
	if len(v.functionBuilders) == 0 {
		panic("COMPILER ERROR: no function builder found")
	}
	return v.functionBuilders[len(v.functionBuilders)-1].AddConstant(value)
}

func (v *InstructionsVisitor) addParam(param ValueType) {
	if len(v.functionBuilders) == 0 {
		panic("COMPILER ERROR: no function builder found")
	}
	v.functionBuilders[len(v.functionBuilders)-1].AddParam(param)
}

func (v *InstructionsVisitor) addFunctionBuilder(functionBuilder FunctionBuilder) {
	v.functionBuilders = append(v.functionBuilders, functionBuilder)
}

func (v *InstructionsVisitor) buildFunctionProto(scope *Scope) int {
	if len(v.functionBuilders) == 0 {
		panic("COMPILER ERROR: no function builder found")
	}
	functionBuilder := v.functionBuilders[len(v.functionBuilders)-1]
	v.functionProtos = append(v.functionProtos, functionBuilder.Build(scope))
	v.functionBuilders = v.functionBuilders[:len(v.functionBuilders)-1]
	return len(v.functionProtos) - 1
}

func (v *InstructionsVisitor) currentFunctionBuilder() *FunctionBuilder {
	if len(v.functionBuilders) == 0 {
		panic("COMPILER ERROR: no function builder found")
	}
	return &v.functionBuilders[len(v.functionBuilders)-1]
}

func (v *InstructionsVisitor) enterFunctionScope(returnType ValueType) {
	if v.scope == nil {
		v.scope = NewScope(true)
	} else {
		v.scope = v.scope.NewChildScope(true)
	}
	v.addFunctionBuilder(*NewFunctionBuilder(returnType))
}

func (v *InstructionsVisitor) exitFunctionScope() int {
	functionSlot := v.buildFunctionProto(v.scope)
	v.scope = v.scope.parent
	return functionSlot
}

func (v *InstructionsVisitor) enterBlockScope() {
	if v.scope == nil {
		v.scope = NewScope(false)
	} else {
		v.scope = v.scope.NewChildScope(false)
	}
}

func (v *InstructionsVisitor) exitBlockScope() {
	v.scope = v.scope.parent
}

// ---------- Getters ----------

func (v *InstructionsVisitor) Errors() []common.Error {
	return v.errors
}

func (v *InstructionsVisitor) Warnings() []common.Error {
	return v.warnings
}

func (v *InstructionsVisitor) FunctionProtos() []FunctionProto {
	return v.functionProtos
}

// ---------- Visitor Implementations ----------

func (v *InstructionsVisitor) VisitProgram(n *ast.Program) any {
	v.enterFunctionScope(VAL_INT)
	for _, statement := range n.Statements {
		statement.Visit(v)
		v.resetReg()
	}
	v.exitFunctionScope()
	return nil
}

func (v *InstructionsVisitor) VisitDeclaration(n *ast.Declaration) any {
	_, ok := v.scope.FindLocalVariable(n.Identifier.Name)
	if ok {
		v.addError(fmt.Sprintf("variable %s already defined", n.Identifier.Name), n.Identifier.Pos())
		return nil
	}
	if n.Value != nil {
		result := n.Value.Visit(v)
		resultVisitExpr, ok := CastVisitExprResult(result)
		if !ok {
			return nil
		}

		if n.IsTyped {
			if resultVisitExpr.TypeOf != ValueTypeFromTypeSubkind(n.TypeOf) {
				v.addError(fmt.Sprintf("variable %s is of type %s, but declaration is of type %s", n.Identifier.Name, resultVisitExpr.TypeOf, ValueTypeFromTypeSubkind(n.TypeOf)), n.Identifier.Pos())
				return nil
			}
		}

		slot := v.scope.DefineVariable(n.Identifier.Name, n.IsMutable, resultVisitExpr.TypeOf)

		v.addInstruction(InstrStoreVar(resultVisitExpr.Reg, slot))
	} else {
		if !n.IsTyped {
			v.addError(fmt.Sprintf("type is required for declaration of variable %s with default value", n.Identifier.Name), n.Identifier.Pos())
			return nil
		}
		if !n.IsMutable {
			v.addError(fmt.Sprintf("constant %s must have a value", n.Identifier.Name), n.Identifier.Pos())
			return nil
		}
		defaultValue := ValueTypeFromTypeSubkind(n.TypeOf).DefaultValue()
		constIndex := v.addConstant(defaultValue)
		slot := v.scope.DefineVariable(n.Identifier.Name, n.IsMutable, ValueTypeFromTypeSubkind(n.TypeOf))

		reg := v.nextReg()
		v.addInstruction(InstrLoadConst(reg, constIndex))
		v.addInstruction(InstrStoreVar(reg, slot))
	}

	return nil
}

func (v *InstructionsVisitor) VisitAssignment(n *ast.Assignment) any {
	localVar, upvar, ok := v.scope.FindVariable(n.Identifier.Name)
	if !ok {
		v.addError(fmt.Sprintf("variable %s not found", n.Identifier.Name), n.Identifier.Pos())
		return nil
	}
	result := n.Value.Visit(v)
	resultVisitExpr, ok := CastVisitExprResult(result)
	if !ok {
		return nil
	}

	if localVar != nil {
		if !localVar.Mutable {
			v.addError(fmt.Sprintf("variable %s is not mutable", n.Identifier.Name), n.Identifier.Pos())
			return nil
		}
		if resultVisitExpr.TypeOf != localVar.TypeOf {
			v.addError(fmt.Sprintf("variable %s is of type %s, but assignment is of type %s", n.Identifier.Name, localVar.TypeOf, resultVisitExpr.TypeOf), n.Identifier.Pos())
			return nil
		}
		v.addInstruction(InstrStoreVar(resultVisitExpr.Reg, localVar.Slot))
	} else if upvar != nil {
		if !upvar.Mutable {
			v.addError(fmt.Sprintf("variable %s is not mutable", n.Identifier.Name), n.Identifier.Pos())
			return nil
		}
		if resultVisitExpr.TypeOf != upvar.TypeOf {
			v.addError(fmt.Sprintf("variable %s is of type %s, but assignment is of type %s", n.Identifier.Name, upvar.TypeOf, resultVisitExpr.TypeOf), n.Identifier.Pos())
			return nil
		}
		v.addInstruction(InstrAssignUpvar(resultVisitExpr.Reg, upvar.LocalSlot))
	}

	return nil
}

func (v *InstructionsVisitor) VisitReturn(n *ast.Return) any {
	result := n.Value.Visit(v)
	resultVisitExpr, ok := CastVisitExprResult(result)
	if !ok {
		return nil
	}
	builder := v.currentFunctionBuilder()
	if resultVisitExpr.TypeOf != builder.ReturnType {
		v.addError(fmt.Sprintf("return value must be of type %s, but got %s", builder.ReturnType, resultVisitExpr.TypeOf), n.Value.Pos())
		return nil
	}
	v.addInstruction(InstrReturn(resultVisitExpr.Reg))
	return nil
}

func (v *InstructionsVisitor) VisitIntLiteral(n *ast.IntLiteral) any {
	reg := v.nextReg()
	constIndex := v.addConstant(NewIntValue(n.Value))
	v.addInstruction(InstrLoadConst(reg, constIndex))
	return &VisitExprResult{Reg: reg, TypeOf: VAL_INT}
}

func (v *InstructionsVisitor) VisitBoolLiteral(n *ast.BoolLiteral) any {
	reg := v.nextReg()
	constIndex := v.addConstant(NewBoolValue(n.Value))
	v.addInstruction(InstrLoadConst(reg, constIndex))
	return &VisitExprResult{Reg: reg, TypeOf: VAL_BOOL}
}

func (v *InstructionsVisitor) VisitNullLiteral(n *ast.NullLiteral) any {
	reg := v.nextReg()
	constIndex := v.addConstant(NewNullValue())
	v.addInstruction(InstrLoadConst(reg, constIndex))
	return &VisitExprResult{Reg: reg, TypeOf: VAL_NULL}
}

func (v *InstructionsVisitor) VisitIdentifier(n *ast.Identifier) any {
	localVar, upvar, ok := v.scope.FindVariable(n.Name)
	if !ok {
		v.addError(fmt.Sprintf("variable %s not found", n.Name), n.Pos())
		return nil
	}
	reg := v.nextReg()
	var typeOf ValueType
	var callable *Callable
	if localVar != nil {
		v.addInstruction(InstrLoadVar(reg, localVar.Slot))
		typeOf = localVar.TypeOf
		callable = localVar.Callable
	} else if upvar != nil {
		v.addInstruction(InstrLoadUpvar(reg, upvar.LocalSlot))
		typeOf = upvar.TypeOf
		callable = upvar.Callable
	}
	return &VisitExprResult{Reg: reg, TypeOf: typeOf, Callable: callable}
}

func (v *InstructionsVisitor) VisitBinaryExpr(n *ast.BinaryExpr) any {
	leftResult := n.Left.Visit(v)
	leftVisitExpr, ok := CastVisitExprResult(leftResult)
	if !ok {
		return nil
	}
	rightResult := n.Right.Visit(v)
	rightVisitExpr, ok := CastVisitExprResult(rightResult)
	if !ok {
		return nil
	}

	opInfo, ok := ResolveBinaryOp(n.Operator, leftVisitExpr.TypeOf, rightVisitExpr.TypeOf)
	if !ok {
		v.addError(fmt.Sprintf("binary operator %s is not supported for types %s and %s", n.Operator, leftVisitExpr.TypeOf, rightVisitExpr.TypeOf), n.Pos())
		return nil
	}
	reg := v.nextReg()
	v.addInstruction(InstrBinary(opInfo.OpCode, reg, leftVisitExpr.Reg, rightVisitExpr.Reg))
	return &VisitExprResult{Reg: reg, TypeOf: opInfo.ResultType}
}

func (v *InstructionsVisitor) VisitParam(n *ast.Param) any {
	v.addParam(ValueTypeFromTypeSubkind(n.TypeOf))
	v.scope.DefineVariable(n.Name, false, ValueTypeFromTypeSubkind(n.TypeOf))
	return nil
}

func (v *InstructionsVisitor) VisitFunction(n *ast.Function) any {
	_, ok := v.scope.FindLocalVariable(n.Name)
	if ok {
		v.addError(fmt.Sprintf("variable %s already defined", n.Name), n.Pos())
		return nil
	}

	v.enterFunctionScope(ValueTypeFromTypeSubkind(n.ReturnType))

	for _, param := range n.Params {
		param.Visit(v)
	}
	for _, statement := range n.Body {
		statement.Visit(v)
	}

	builder := v.currentFunctionBuilder()

	functionSlot := v.exitFunctionScope()

	slot := v.scope.DefineCallableVariable(n.Name, false, VAL_CLOSURE, builder.Params, builder.ReturnType)
	reg := v.nextReg()
	v.addInstruction(InstrClosure(reg, functionSlot))
	v.addInstruction(InstrStoreVar(reg, slot))
	return &VisitExprResult{Reg: reg, TypeOf: VAL_CLOSURE}
}

func (v *InstructionsVisitor) VisitBlock(n *ast.Block) any {
	v.enterBlockScope()
	for _, statement := range n.Statements {
		statement.Visit(v)
	}
	v.exitBlockScope()
	return nil
}

func (v *InstructionsVisitor) VisitCallExpr(n *ast.CallExpr) any {
	result := n.Identifier.Visit(v)
	resultVisitExpr, ok := CastVisitExprResult(result)
	if !ok {
		return nil
	}
	if resultVisitExpr.TypeOf != VAL_CLOSURE {
		v.addError(fmt.Sprintf("function call must be of type closure, but got %s", resultVisitExpr.TypeOf), n.Identifier.Pos())
		return nil
	}
	if resultVisitExpr.Callable == nil {
		v.addError(fmt.Sprintf("function %s is not callable", n.Identifier.Name), n.Identifier.Pos())
		return nil
	}

	args := []int{}
	isError := false

	if len(n.Arguments) != len(resultVisitExpr.Callable.CallArgs) {
		v.addError(fmt.Sprintf("function %s takes %d arguments, but got %d", n.Identifier.Name, len(resultVisitExpr.Callable.CallArgs), len(n.Arguments)), n.Identifier.Pos())
		isError = true
	}

	for i, argument := range n.Arguments {
		paramType := resultVisitExpr.Callable.CallArgs[i]

		argumentResult := argument.Visit(v)
		argumentVisitExpr, ok := CastVisitExprResult(argumentResult)
		if !ok {
			return nil
		}

		if argumentVisitExpr.TypeOf != paramType {
			v.addError(fmt.Sprintf("argument %d must be of type %s, but got %s", i, paramType, argumentVisitExpr.TypeOf), argument.Pos())
			isError = true
			continue
		}
		args = append(args, argumentVisitExpr.Reg)
	}

	if isError {
		return nil
	}

	resultReg := v.nextReg()
	v.addInstruction(InstrCall(resultReg, resultVisitExpr.Reg, args))
	return &VisitExprResult{Reg: resultReg, TypeOf: resultVisitExpr.Callable.ReturnType}
}

func (v *InstructionsVisitor) VisitIf(n *ast.If) any {
	conditionResult := n.Condition.Visit(v)
	conditionVisitExpr, ok := CastVisitExprResult(conditionResult)
	if !ok {
		return nil
	}
	if conditionVisitExpr.TypeOf != VAL_BOOL {
		v.addError(fmt.Sprintf("condition must be of type bool, but got %s", conditionVisitExpr.TypeOf), n.Condition.Pos())
		return nil
	}
	reg := conditionVisitExpr.Reg
	jumpIfFalseIndex := v.addInstruction(InstrJumpIfFalse(reg, -1))

	for _, statement := range n.Body {
		statement.Visit(v)
	}
	elseBodyIndex := -1
	if len(n.ElseBody) > 0 {
		elseBodyIndex = v.addInstruction(InstrJump(-1))
	}
	endIfTarget := v.instructionsLength() - 1
	v.setInstruction(jumpIfFalseIndex, InstrJumpIfFalse(reg, endIfTarget))

	for _, statement := range n.ElseBody {
		statement.Visit(v)
	}

	if elseBodyIndex != -1 {
		endElseTarget := v.instructionsLength() - 1
		v.setInstruction(elseBodyIndex, InstrJump(endElseTarget))
	}
	return nil
}
