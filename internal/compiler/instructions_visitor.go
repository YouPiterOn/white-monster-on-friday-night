package compiler

import (
	"fmt"

	"youpiteron.dev/white-monster-on-friday-night/internal/ast"
	"youpiteron.dev/white-monster-on-friday-night/internal/common"
)

type VisitExprResult struct {
	Reg           int
	TypeOf        Type
	FuncSignature *FuncSignature
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
	context        Context
	globalTable    *GlobalTable
	errors         []common.Error
	warnings       []common.Error
	reg            int
	functionProtos []FunctionProto
	moduleProtos   []ModuleProto
}

// ---------- Constructor ----------

func NewInstructionsVisitor() *InstructionsVisitor {
	globalTable := NewGlobalTable()
	RegisterStdGlobals(globalTable)
	return &InstructionsVisitor{context: nil, globalTable: globalTable, errors: []common.Error{}, warnings: []common.Error{}, reg: 0, functionProtos: []FunctionProto{}}
}

// ---------- Helpers ----------

func (v *InstructionsVisitor) EnterModuleContext() {
	v.context = NewModuleContext()
}

func (v *InstructionsVisitor) ExitModuleContext() {
	moduleContext := CastModuleContext(v.context)
	moduleProto := BuildModuleProto(*moduleContext, v.functionProtos)
	v.moduleProtos = append(v.moduleProtos, *moduleProto)
	v.context = nil
	v.functionProtos = []FunctionProto{}
}

func (v *InstructionsVisitor) EmitModuleProto() *ModuleProto {
	moduleContext := CastModuleContext(v.context)
	moduleProto := BuildModuleProto(*moduleContext, v.functionProtos)
	moduleContext.ClearInstructions()
	return moduleProto
}

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

func (v *InstructionsVisitor) enterFunctionContext(returnType Type) {
	v.context = NewFunctionContext(v.context, returnType)
}

func (v *InstructionsVisitor) exitFunctionContext() int {
	functionProto := BuildFunctionProto(v.context)
	v.functionProtos = append(v.functionProtos, *functionProto)
	v.context = v.context.Parent()
	return len(v.functionProtos) - 1
}

func (v *InstructionsVisitor) enterBlockContext() {
	v.context = NewBlockContext(v.context)
}

func (v *InstructionsVisitor) exitBlockContext() {
	v.context = v.context.Parent()
}

// ---------- Visitor Implementations ----------

func (v *InstructionsVisitor) VisitProgram(n *ast.Program) any {
	for _, statement := range n.Statements {
		statement.Visit(v)
		v.resetReg()
	}
	return nil
}

func (v *InstructionsVisitor) VisitDeclaration(n *ast.Declaration) any {
	_, ok := v.context.FindLocalVariable(n.Identifier.Name)
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
			if !resultVisitExpr.TypeOf.IsEqual(TypeFromTypeSubkind(n.TypeOf)) {
				v.addError(fmt.Sprintf("variable %s is of type %s, but declaration is of type %s", n.Identifier.Name, resultVisitExpr.TypeOf, TypeFromTypeSubkind(n.TypeOf)), n.Identifier.Pos())
				return nil
			}
		}

		slot := v.context.DefineVariable(n.Identifier.Name, n.IsMutable, resultVisitExpr.TypeOf)

		v.context.AddInstruction(InstrStoreVar(resultVisitExpr.Reg, slot))
	} else {
		if !n.IsTyped {
			v.addError(fmt.Sprintf("type is required for declaration of variable %s with default value", n.Identifier.Name), n.Identifier.Pos())
			return nil
		}
		if !n.IsMutable {
			v.addError(fmt.Sprintf("constant %s must have a value", n.Identifier.Name), n.Identifier.Pos())
			return nil
		}
		defaultValue := TypeFromTypeSubkind(n.TypeOf).DefaultValue()
		constIndex := v.context.AddConstant(defaultValue)
		slot := v.context.DefineVariable(n.Identifier.Name, n.IsMutable, TypeFromTypeSubkind(n.TypeOf))

		reg := v.nextReg()
		v.context.AddInstruction(InstrLoadConst(reg, constIndex))
		v.context.AddInstruction(InstrStoreVar(reg, slot))
	}

	return nil
}

func (v *InstructionsVisitor) VisitAssignment(n *ast.Assignment) any {
	var globalVar *Variable
	localVar, upvar, ok := v.context.FindVariable(n.Identifier.Name)
	if !ok {
		globalVar, ok = v.globalTable.FindVariable(n.Identifier.Name)
		if !ok {
			v.addError(fmt.Sprintf("variable %s not found", n.Identifier.Name), n.Identifier.Pos())
			return nil
		}
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
		if !resultVisitExpr.TypeOf.IsEqual(localVar.TypeOf) {
			v.addError(fmt.Sprintf("variable %s is of type %s, but assignment is of type %s", n.Identifier.Name, localVar.TypeOf, resultVisitExpr.TypeOf), n.Identifier.Pos())
			return nil
		}
		v.context.AddInstruction(InstrStoreVar(resultVisitExpr.Reg, localVar.Slot))
	} else if upvar != nil {
		if !upvar.Mutable {
			v.addError(fmt.Sprintf("variable %s is not mutable", n.Identifier.Name), n.Identifier.Pos())
			return nil
		}
		if !resultVisitExpr.TypeOf.IsEqual(upvar.TypeOf) {
			v.addError(fmt.Sprintf("variable %s is of type %s, but assignment is of type %s", n.Identifier.Name, upvar.TypeOf, resultVisitExpr.TypeOf), n.Identifier.Pos())
			return nil
		}
		v.context.AddInstruction(InstrAssignUpvar(resultVisitExpr.Reg, upvar.LocalSlot))
	} else if globalVar != nil {
		if !globalVar.Mutable {
			v.addError(fmt.Sprintf("variable %s is not mutable", n.Identifier.Name), n.Identifier.Pos())
			return nil
		}
		if !resultVisitExpr.TypeOf.IsEqual(globalVar.TypeOf) {
			v.addError(fmt.Sprintf("variable %s is of type %s, but assignment is of type %s", n.Identifier.Name, globalVar.TypeOf, resultVisitExpr.TypeOf), n.Identifier.Pos())
			return nil
		}

		v.context.AddInstruction(InstrAssignGlobal(resultVisitExpr.Reg, globalVar.Slot))
	}

	return nil
}

func (v *InstructionsVisitor) VisitReturn(n *ast.Return) any {
	result := n.Value.Visit(v)
	resultVisitExpr, ok := CastVisitExprResult(result)
	if !ok {
		return nil
	}
	returnType := v.context.ReturnType()
	if !resultVisitExpr.TypeOf.IsEqual(returnType) {
		v.addError(fmt.Sprintf("return value must be of type %s, but got %s", returnType, resultVisitExpr.TypeOf), n.Value.Pos())
		return nil
	}
	v.context.AddInstruction(InstrReturn(resultVisitExpr.Reg))
	return nil
}

func (v *InstructionsVisitor) VisitIntLiteral(n *ast.IntLiteral) any {
	if n.IsStatement {
		return nil
	}
	reg := v.nextReg()
	constIndex := v.context.AddConstant(NewIntValue(n.Value))
	v.context.AddInstruction(InstrLoadConst(reg, constIndex))
	return &VisitExprResult{Reg: reg, TypeOf: TypeInt()}
}

func (v *InstructionsVisitor) VisitBoolLiteral(n *ast.BoolLiteral) any {
	if n.IsStatement {
		return nil
	}
	reg := v.nextReg()
	constIndex := v.context.AddConstant(NewBoolValue(n.Value))
	v.context.AddInstruction(InstrLoadConst(reg, constIndex))
	return &VisitExprResult{Reg: reg, TypeOf: TypeBool()}
}

func (v *InstructionsVisitor) VisitNullLiteral(n *ast.NullLiteral) any {
	if n.IsStatement {
		return nil
	}
	reg := v.nextReg()
	constIndex := v.context.AddConstant(NewNullValue())
	v.context.AddInstruction(InstrLoadConst(reg, constIndex))
	return &VisitExprResult{Reg: reg, TypeOf: TypeNull()}
}

func (v *InstructionsVisitor) VisitIdentifier(n *ast.Identifier) any {
	if n.IsStatement {
		return nil
	}
	var globalVar *Variable
	localVar, upvar, ok := v.context.FindVariable(n.Name)
	if !ok {
		globalVar, ok = v.globalTable.FindVariable(n.Name)
		if !ok {
			v.addError(fmt.Sprintf("variable %s not found", n.Name), n.Pos())
			return nil
		}
	}
	reg := v.nextReg()
	var typeOf Type
	var callable *FuncSignature
	if localVar != nil {
		v.context.AddInstruction(InstrLoadVar(reg, localVar.Slot))
		typeOf = localVar.TypeOf
		callable = localVar.FuncSignature
	} else if upvar != nil {
		v.context.AddInstruction(InstrLoadUpvar(reg, upvar.LocalSlot))
		typeOf = upvar.TypeOf
		callable = upvar.FuncSignature
	} else if globalVar != nil {
		v.context.AddInstruction(InstrLoadGlobal(reg, globalVar.Slot))
		typeOf = globalVar.TypeOf
		callable = globalVar.FuncSignature
	}
	return &VisitExprResult{Reg: reg, TypeOf: typeOf, FuncSignature: callable}
}

func (v *InstructionsVisitor) VisitBinaryExpr(n *ast.BinaryExpr) any {
	leftResult := n.Left.Visit(v)
	leftVisitExpr, leftOk := CastVisitExprResult(leftResult)
	rightResult := n.Right.Visit(v)
	rightVisitExpr, rightOk := CastVisitExprResult(rightResult)

	if !leftOk || !rightOk || n.IsStatement {
		return nil
	}

	opInfo, ok := ResolveBinaryOp(n.Operator, leftVisitExpr.TypeOf, rightVisitExpr.TypeOf)
	if !ok {
		v.addError(fmt.Sprintf("binary operator %s is not supported for types %s and %s", n.Operator, leftVisitExpr.TypeOf, rightVisitExpr.TypeOf), n.Pos())
		return nil
	}
	reg := v.nextReg()
	v.context.AddInstruction(InstrBinary(opInfo.OpCode, reg, leftVisitExpr.Reg, rightVisitExpr.Reg))
	return &VisitExprResult{Reg: reg, TypeOf: opInfo.ResultType}
}

func (v *InstructionsVisitor) VisitParam(n *ast.Param) any {
	v.context.AddParam(TypeFromTypeSubkind(n.TypeOf))
	v.context.DefineVariable(n.Name, false, TypeFromTypeSubkind(n.TypeOf))
	return nil
}

func (v *InstructionsVisitor) VisitFunction(n *ast.Function) any {
	_, ok := v.context.FindLocalVariable(n.Name)
	if ok {
		v.addError(fmt.Sprintf("variable %s already defined", n.Name), n.Pos())
		return nil
	}

	v.enterFunctionContext(TypeFromTypeSubkind(n.ReturnType))

	for _, param := range n.Params {
		param.Visit(v)
	}
	for _, statement := range n.Body {
		statement.Visit(v)
	}

	params := v.context.Params()
	returnType := v.context.ReturnType()

	functionSlot := v.exitFunctionContext()

	slot := v.context.DefineFunctionVariable(n.Name, false, TypeClosure(), &FuncSignature{CallArgs: params, ReturnType: returnType})
	reg := v.nextReg()
	v.context.AddInstruction(InstrClosure(reg, functionSlot))
	v.context.AddInstruction(InstrStoreVar(reg, slot))
	return &VisitExprResult{Reg: reg, TypeOf: TypeClosure()}
}

func (v *InstructionsVisitor) VisitBlock(n *ast.Block) any {
	v.enterBlockContext()
	for _, statement := range n.Statements {
		statement.Visit(v)
	}
	v.exitBlockContext()
	return nil
}

func (v *InstructionsVisitor) VisitCallExpr(n *ast.CallExpr) any {
	result := n.Identifier.Visit(v)
	resultVisitExpr, ok := CastVisitExprResult(result)
	if !ok {
		return nil
	}
	if !resultVisitExpr.TypeOf.IsEqual(TypeClosure()) && !resultVisitExpr.TypeOf.IsEqual(TypeNativeFunction()) {
		v.addError(fmt.Sprintf("function call must be of type closure, but got %s", resultVisitExpr.TypeOf), n.Identifier.Pos())
		return nil
	}
	if resultVisitExpr.FuncSignature == nil {
		v.addError(fmt.Sprintf("function %s is not callable", n.Identifier.Name), n.Identifier.Pos())
		return nil
	}

	args := []int{}
	isError := false

	if len(n.Arguments) != len(resultVisitExpr.FuncSignature.CallArgs) {
		v.addError(fmt.Sprintf("function %s takes %d arguments, but got %d", n.Identifier.Name, len(resultVisitExpr.FuncSignature.CallArgs), len(n.Arguments)), n.Identifier.Pos())
		isError = true
	}

	for i, argument := range n.Arguments {
		paramType := resultVisitExpr.FuncSignature.CallArgs[i]

		argumentResult := argument.Visit(v)
		argumentVisitExpr, ok := CastVisitExprResult(argumentResult)
		if !ok {
			return nil
		}

		if !argumentVisitExpr.TypeOf.IsEqual(paramType) {
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
	v.context.AddInstruction(InstrCall(resultReg, resultVisitExpr.Reg, args))
	return &VisitExprResult{Reg: resultReg, TypeOf: resultVisitExpr.FuncSignature.ReturnType}
}

func (v *InstructionsVisitor) VisitIf(n *ast.If) any {
	conditionResult := n.Condition.Visit(v)
	conditionVisitExpr, ok := CastVisitExprResult(conditionResult)
	if !ok {
		return nil
	}
	if !conditionVisitExpr.TypeOf.IsEqual(TypeBool()) {
		v.addError(fmt.Sprintf("condition must be of type bool, but got %s", conditionVisitExpr.TypeOf), n.Condition.Pos())
		return nil
	}
	reg := conditionVisitExpr.Reg
	jumpIfFalseIndex := v.context.AddInstruction(InstrJumpIfFalse(reg, -1))

	for _, statement := range n.Body {
		statement.Visit(v)
	}
	elseBodyIndex := -1
	if len(n.ElseBody) > 0 {
		elseBodyIndex = v.context.AddInstruction(InstrJump(-1))
	}
	endIfTarget := v.context.InstructionsLength() - 1
	v.context.SetInstruction(jumpIfFalseIndex, InstrJumpIfFalse(reg, endIfTarget))

	for _, statement := range n.ElseBody {
		statement.Visit(v)
	}

	if elseBodyIndex != -1 {
		endElseTarget := v.context.InstructionsLength() - 1
		v.context.SetInstruction(elseBodyIndex, InstrJump(endElseTarget))
	}
	return nil
}
