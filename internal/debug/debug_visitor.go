package debug

import (
	"fmt"

	"youpiteron.dev/white-monster-on-friday-night/internal/ast"
)

type DebugVisitor struct {
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
