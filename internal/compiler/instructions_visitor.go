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
	instructions []Instruction
	scope        *Scope
	errors       []SyntaxError
	reg          int
}

// ---------- Constructor ----------

func NewInstructionsVisitor() *InstructionsVisitor {
	return &InstructionsVisitor{scope: NewScope(), instructions: []Instruction{}, errors: []SyntaxError{}, reg: 0}
}

// ---------- Helpers ----------

func (v *InstructionsVisitor) addError(message string, pos *lexer.SourcePos) {
	v.errors = append(v.errors, SyntaxError{Message: message, Pos: pos})
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

// ---------- Getters ----------

func (v *InstructionsVisitor) Errors() []SyntaxError {
	return v.errors
}

func (v *InstructionsVisitor) Instructions() []Instruction {
	return v.instructions
}

func (v *InstructionsVisitor) VisitProgram(n *ast.Program) any {
	for _, statement := range n.Statements {
		statement.Visit(v)
		v.resetReg()
	}
	return nil
}

// ---------- Visitor Implementations ----------

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
	v.instructions = append(v.instructions, StoreVar(valueRxInt, slot))

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
			v.addError(fmt.Sprintf("variable %s is not mutable", n.Identifier.Name), n.Identifier.Pos())
			return nil
		}
		v.instructions = append(v.instructions, StoreVar(valueRxInt, localVar.slot))
	} else if upvar != nil {
		if !upvar.mutable {
			v.addError(fmt.Sprintf("variable %s is not mutable", n.Identifier.Name), n.Identifier.Pos())
			return nil
		}
		v.instructions = append(v.instructions, AssignUpvar(valueRxInt, upvar.localSlot))
	}
	return nil
}

func (v *InstructionsVisitor) VisitReturn(n *ast.Return) any {
	valueRx := n.Value.Visit(v)
	valueRxInt, ok := valueRx.(int)
	if !ok {
		panic(fmt.Sprintf("COMPILER ERROR: value %v is not an integer", valueRxInt))
	}
	v.instructions = append(v.instructions, Return(valueRxInt))
	return nil
}

func (v *InstructionsVisitor) VisitNumberLiteral(n *ast.NumberLiteral) any {
	reg := v.reg
	v.instructions = append(v.instructions, LoadConst(reg, n.Value))
	v.reg++
	return reg
}

func (v *InstructionsVisitor) VisitIdentifier(n *ast.Identifier) any {
	localVar, upvar, ok := v.scope.FindVariable(n.Name)
	if !ok {
		v.addError(fmt.Sprintf("variable %s not found", n.Name), n.Pos())
		return nil
	}
	if localVar != nil {
		v.instructions = append(v.instructions, LoadVar(v.resetReg(), localVar.slot))
	} else if upvar != nil {
		v.instructions = append(v.instructions, LoadUpvar(v.resetReg(), upvar.localSlot))
	}
	return nil
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
		v.instructions = append(v.instructions, Add(reg, leftRxInt, rightRxInt))
	case lexer.OperatorMinus:
		v.instructions = append(v.instructions, Sub(reg, leftRxInt, rightRxInt))
	case lexer.OperatorStar:
		v.instructions = append(v.instructions, Mul(reg, leftRxInt, rightRxInt))
	case lexer.OperatorSlash:
		v.instructions = append(v.instructions, Div(reg, leftRxInt, rightRxInt))
	}
	return reg
}

func (v *InstructionsVisitor) VisitParameter(n *ast.Parameter) any {
	return nil
}

func (v *InstructionsVisitor) VisitFunction(n *ast.Function) any {
	return nil
}

func (v *InstructionsVisitor) VisitBlock(n *ast.Block) any {
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
	v.instructions = append(v.instructions, Call(identifierRxInt, args))
	return nil
}
