package compiler

import (
	"fmt"

	"youpiteron.dev/white-monster-on-friday-night/internal/ast"
	"youpiteron.dev/white-monster-on-friday-night/internal/lexer"
)

type SyntaxError struct {
	Message string
	Pos     *lexer.SourcePos
}

type InstructionsVisitor struct {
	scope            *Scope
	errors           []SyntaxError
	warnings         []SyntaxError
	reg              int
	functionBuilders []FunctionBuilder
	functionProtos   []FunctionProto
}

// ---------- Constructor ----------

func NewInstructionsVisitor() *InstructionsVisitor {
	return &InstructionsVisitor{scope: nil, errors: []SyntaxError{}, warnings: []SyntaxError{}, reg: 0, functionBuilders: []FunctionBuilder{}, functionProtos: []FunctionProto{}}
}

// ---------- Helpers ----------

func (v *InstructionsVisitor) addError(message string, pos *lexer.SourcePos) {
	v.errors = append(v.errors, SyntaxError{Message: message, Pos: pos})
}

func (v *InstructionsVisitor) addWarning(message string, pos *lexer.SourcePos) {
	v.warnings = append(v.warnings, SyntaxError{Message: message, Pos: pos})
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

func (v *InstructionsVisitor) addInstruction(instruction Instruction) {
	if len(v.functionBuilders) == 0 {
		panic("COMPILER ERROR: no function proto found")
	}
	v.functionBuilders[len(v.functionBuilders)-1].AddInstruction(instruction)
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

func (v *InstructionsVisitor) enterFunctionScope(name string, numParams int) {
	if v.scope == nil {
		v.scope = NewScope(true)
	} else {
		v.scope = v.scope.NewChildScope(true)
	}
	v.addFunctionBuilder(*NewFunctionBuilder(name, numParams))
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

func (v *InstructionsVisitor) Errors() []SyntaxError {
	return v.errors
}

func (v *InstructionsVisitor) Warnings() []SyntaxError {
	return v.warnings
}

func (v *InstructionsVisitor) FunctionProtos() []FunctionProto {
	return v.functionProtos
}

// ---------- Visitor Implementations ----------

func (v *InstructionsVisitor) VisitProgram(n *ast.Program) any {
	v.enterFunctionScope("main", 0)
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
	slot := v.scope.DefineVariable(n.Identifier.Name, n.Specifier == lexer.KeywordVar)
	valueRx := n.Value.Visit(v)

	valueRxInt, ok := valueRx.(int)
	if !ok {
		panic(fmt.Sprintf("COMPILER ERROR: value %v is not an integer", valueRxInt))
	}
	v.addInstruction(StoreVar(valueRxInt, slot))

	return nil
}

func (v *InstructionsVisitor) VisitAssignment(n *ast.Assignment) any {
	localVar, upvar, ok := v.scope.FindVariable(n.Identifier.Name)
	if !ok {
		v.addError(fmt.Sprintf("variable %s not found", n.Identifier.Name), n.Identifier.Pos())
		return nil
	}
	valueRx := n.Value.Visit(v)
	valueRxInt, ok := valueRx.(int)
	if !ok {
		panic(fmt.Sprintf("COMPILER ERROR: value %v is not an integer", valueRxInt))
	}
	if localVar != nil {
		if !localVar.mutable {
			v.addWarning(fmt.Sprintf("variable %s is not mutable", n.Identifier.Name), n.Identifier.Pos())
			return nil
		}
		v.addInstruction(StoreVar(valueRxInt, localVar.slot))
	} else if upvar != nil {
		if !upvar.mutable {
			v.addWarning(fmt.Sprintf("variable %s is not mutable", n.Identifier.Name), n.Identifier.Pos())
			return nil
		}
		v.addInstruction(AssignUpvar(valueRxInt, upvar.localSlot))
	}

	return nil
}

func (v *InstructionsVisitor) VisitReturn(n *ast.Return) any {
	valueRx := n.Value.Visit(v)
	valueRxInt, ok := valueRx.(int)
	if !ok {
		panic(fmt.Sprintf("COMPILER ERROR: value %v is not an integer", valueRxInt))
	}
	v.addInstruction(Return(valueRxInt))
	return nil
}

func (v *InstructionsVisitor) VisitNumberLiteral(n *ast.NumberLiteral) any {
	reg := v.reg
	v.addInstruction(LoadConst(reg, n.Value))
	v.reg++
	return reg
}

func (v *InstructionsVisitor) VisitIdentifier(n *ast.Identifier) any {
	localVar, upvar, ok := v.scope.FindVariable(n.Name)
	if !ok {
		v.addError(fmt.Sprintf("variable %s not found", n.Name), n.Pos())
		return 0
	}
	reg := v.nextReg()
	if localVar != nil {
		v.addInstruction(LoadVar(reg, localVar.slot))
	} else if upvar != nil {
		v.addInstruction(LoadUpvar(reg, upvar.localSlot))
	}
	return reg
}

func (v *InstructionsVisitor) VisitBinaryExpr(n *ast.BinaryExpr) any {
	leftRx := n.Left.Visit(v)
	leftRxInt, ok := leftRx.(int)
	if !ok {
		panic(fmt.Sprintf("COMPILER ERROR: left value %v is not an integer", leftRxInt))
	}
	rightRx := n.Right.Visit(v)
	rightRxInt, ok := rightRx.(int)
	if !ok {
		panic(fmt.Sprintf("COMPILER ERROR: right value %v is not an integer", rightRxInt))
	}
	reg := v.nextReg()
	switch n.Operator {
	case lexer.OperatorPlus:
		v.addInstruction(Add(reg, leftRxInt, rightRxInt))
	case lexer.OperatorMinus:
		v.addInstruction(Sub(reg, leftRxInt, rightRxInt))
	case lexer.OperatorStar:
		v.addInstruction(Mul(reg, leftRxInt, rightRxInt))
	case lexer.OperatorSlash:
		v.addInstruction(Div(reg, leftRxInt, rightRxInt))
	}
	return reg
}

func (v *InstructionsVisitor) VisitParam(n *ast.Param) any {
	v.scope.DefineVariable(n.Name, false)
	return nil
}

func (v *InstructionsVisitor) VisitFunction(n *ast.Function) any {
	_, ok := v.scope.FindLocalVariable(n.Name)
	if ok {
		v.addError(fmt.Sprintf("variable %s already defined", n.Name), n.Pos())
		return nil
	}

	v.enterFunctionScope(n.Name, len(n.Params))

	for _, param := range n.Params {
		param.Visit(v)
	}
	for _, statement := range n.Body {
		statement.Visit(v)
	}

	functionSlot := v.exitFunctionScope()

	slot := v.scope.DefineVariable(n.Name, false)
	reg := v.nextReg()
	v.addInstruction(Closure(reg, functionSlot))
	v.addInstruction(StoreVar(reg, slot))
	return reg
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
	args := []int{}
	for _, argument := range n.Arguments {
		argumentRx := argument.Visit(v)
		argumentRxInt, ok := argumentRx.(int)
		if !ok {
			panic(fmt.Sprintf("COMPILER ERROR: argument registry %v is not an integer", argumentRxInt))
		}
		args = append(args, argumentRxInt)
	}
	identifierRx := n.Identifier.Visit(v)
	identifierRxInt, ok := identifierRx.(int)
	if !ok {
		panic(fmt.Sprintf("COMPILER ERROR: identifier registry %v is not an integer", identifierRxInt))
	}
	resultReg := v.nextReg()
	v.addInstruction(Call(resultReg, identifierRxInt, args))
	return resultReg
}
