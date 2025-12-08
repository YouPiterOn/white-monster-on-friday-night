package debug

import (
	"fmt"

	"youpiteron.dev/white-monster-on-friday-night/internal/ast"
)

type DebugVisitor struct {
}

func (v *DebugVisitor) VisitProgram(n *ast.Program) any {
	fmt.Printf("program: %v\n", n.Pos())
	for _, statement := range n.Statements {
		statement.Visit(v)
	}
	return nil
}

func (v *DebugVisitor) VisitAssignment(n *ast.Assignment) any {
	fmt.Printf("assignment: %v\n", n.Pos())
	fmt.Printf("  specifier: %v\n", n.Specifier)
	n.Identifier.Visit(v)
	n.Value.Visit(v)
	return nil
}

func (v *DebugVisitor) VisitReturn(n *ast.Return) any {
	fmt.Printf("return: %v\n", n.Pos())
	n.Value.Visit(v)
	return nil
}

func (v *DebugVisitor) VisitNumberLiteral(n *ast.NumberLiteral) any {
	fmt.Printf("number literal: %v\n", n.Pos())
	fmt.Printf("  value: %d\n", n.Value)
	return nil
}

func (v *DebugVisitor) VisitIdentifier(n *ast.Identifier) any {
	fmt.Printf("identifier: %v\n", n.Pos())
	fmt.Printf("  name: %s\n", n.Name)
	return nil
}

func (v *DebugVisitor) VisitBinaryExpr(n *ast.BinaryExpr) any {
	fmt.Printf("binary expression: %v\n", n.Pos())
	fmt.Printf("  operator: %v\n", n.Operator)
	n.Left.Visit(v)
	n.Right.Visit(v)
	return nil
}

func (v *DebugVisitor) VisitParameter(n *ast.Parameter) any {
	fmt.Printf("parameter: %v\n", n.Pos())
	fmt.Printf("  name: %s\n", n.Name)
	return nil
}

func (v *DebugVisitor) VisitFunction(n *ast.Function) any {
	fmt.Printf("function: %v\n", n.Pos())
	fmt.Printf("  name: %s\n", n.Name)
	for _, param := range n.Params {
		param.Visit(v)
	}
	n.Body.Visit(v)
	return nil
}

func (v *DebugVisitor) VisitBlock(n *ast.Block) any {
	fmt.Printf("block: %v\n", n.Pos())
	for _, statement := range n.Statements {
		statement.Visit(v)
	}
	return nil
}

func (v *DebugVisitor) VisitCallExpr(n *ast.CallExpr) any {
	fmt.Printf("call expression: %v\n", n.Pos())
	n.Function.Visit(v)
	for _, argument := range n.Arguments {
		argument.Visit(v)
	}
	return nil
}
