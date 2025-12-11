package compiler

import "youpiteron.dev/white-monster-on-friday-night/internal/ast"

type InstructionsVisitor struct {
	instructions []Instruction
}

func (v *InstructionsVisitor) VisitProgram(n *ast.Program) any {
	for _, statement := range n.Statements {
		statement.Visit(v)
	}
	return nil
}

func (v *InstructionsVisitor) VisitAssignment(n *ast.Assignment) any {
	return nil
}

func (v *InstructionsVisitor) VisitReturn(n *ast.Return) any {
	return nil
}

func (v *InstructionsVisitor) VisitNumberLiteral(n *ast.NumberLiteral) any {
	return nil
}

func (v *InstructionsVisitor) VisitIdentifier(n *ast.Identifier) any {
	return nil
}

func (v *InstructionsVisitor) VisitBinaryExpr(n *ast.BinaryExpr) any {
	return nil
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
	return nil
}
