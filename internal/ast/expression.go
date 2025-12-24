package ast

import (
	"youpiteron.dev/white-monster-on-friday-night/internal/common"
	"youpiteron.dev/white-monster-on-friday-night/internal/lexer"
)

type Expression interface {
	Statement
	expressionNode()
}

type IntLiteral struct {
	Value       int
	PosAt       *common.SourcePos
	IsStatement bool
}

func (n *IntLiteral) Pos() *common.SourcePos { return n.PosAt }
func (n *IntLiteral) statementNode()         {}
func (n *IntLiteral) expressionNode()        {}
func (n *IntLiteral) Visit(v Visitor[any]) any {
	return v.VisitIntLiteral(n)
}

type BoolLiteral struct {
	Value       bool
	PosAt       *common.SourcePos
	IsStatement bool
}

func (b *BoolLiteral) Pos() *common.SourcePos { return b.PosAt }
func (b *BoolLiteral) statementNode()         {}
func (b *BoolLiteral) expressionNode()        {}
func (b *BoolLiteral) Visit(v Visitor[any]) any {
	return v.VisitBoolLiteral(b)
}

type NullLiteral struct {
	PosAt       *common.SourcePos
	IsStatement bool
}

func (n *NullLiteral) Pos() *common.SourcePos { return n.PosAt }
func (n *NullLiteral) statementNode()         {}
func (n *NullLiteral) expressionNode()        {}
func (n *NullLiteral) Visit(v Visitor[any]) any {
	return v.VisitNullLiteral(n)
}

type ArrayLiteral struct {
	Elements    []Expression
	PosAt       *common.SourcePos
	IsStatement bool
}

func (a *ArrayLiteral) Pos() *common.SourcePos { return a.PosAt }
func (a *ArrayLiteral) statementNode()         {}
func (a *ArrayLiteral) expressionNode()        {}
func (a *ArrayLiteral) Visit(v Visitor[any]) any {
	return v.VisitArrayLiteral(a)
}

type Identifier struct {
	Name        string
	PosAt       *common.SourcePos
	IsStatement bool
}

func (i *Identifier) Pos() *common.SourcePos { return i.PosAt }
func (i *Identifier) statementNode()         {}
func (i *Identifier) expressionNode()        {}
func (i *Identifier) Visit(v Visitor[any]) any {
	return v.VisitIdentifier(i)
}

type BinaryExpr struct {
	Left        Expression
	Operator    lexer.OperatorSubkind
	Right       Expression
	PosAt       *common.SourcePos
	IsStatement bool
}

func (b *BinaryExpr) Pos() *common.SourcePos { return b.PosAt }
func (b *BinaryExpr) statementNode()         {}
func (b *BinaryExpr) expressionNode()        {}
func (b *BinaryExpr) Visit(v Visitor[any]) any {
	return v.VisitBinaryExpr(b)
}

type CallExpr struct {
	Identifier Identifier
	Arguments  []Expression
	PosAt      *common.SourcePos
}

func (c *CallExpr) Pos() *common.SourcePos { return c.PosAt }
func (c *CallExpr) statementNode()         {}
func (c *CallExpr) expressionNode()        {}
func (c *CallExpr) Visit(v Visitor[any]) any {
	return v.VisitCallExpr(c)
}

type IndexExpr struct {
	Array       Expression
	Index       Expression
	PosAt       *common.SourcePos
	IsStatement bool
}

func (i *IndexExpr) Pos() *common.SourcePos { return i.PosAt }
func (i *IndexExpr) statementNode()         {}
func (i *IndexExpr) expressionNode()        {}
func (i *IndexExpr) Visit(v Visitor[any]) any {
	return v.VisitIndexExpr(i)
}
