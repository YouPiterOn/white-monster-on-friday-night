package ast

import (
	"youpiteron.dev/white-monster-on-friday-night/internal/common"
	"youpiteron.dev/white-monster-on-friday-night/internal/lexer"
)

type Expression interface {
	Statement
	expressionNode()
}

type Factor interface {
	Expression
	factorNode()
}

type IntLiteral struct {
	Value int
	PosAt *common.SourcePos
}

func (n *IntLiteral) Pos() *common.SourcePos { return n.PosAt }
func (n *IntLiteral) statementNode()         {}
func (n *IntLiteral) expressionNode()        {}
func (n *IntLiteral) factorNode()            {}
func (n *IntLiteral) Visit(v Visitor[any]) any {
	return v.VisitIntLiteral(n)
}

type BoolLiteral struct {
	Value bool
	PosAt *common.SourcePos
}

func (b *BoolLiteral) Pos() *common.SourcePos { return b.PosAt }
func (b *BoolLiteral) statementNode()         {}
func (b *BoolLiteral) expressionNode()        {}
func (b *BoolLiteral) factorNode()            {}
func (b *BoolLiteral) Visit(v Visitor[any]) any {
	return v.VisitBoolLiteral(b)
}

type NullLiteral struct {
	PosAt *common.SourcePos
}

func (n *NullLiteral) Pos() *common.SourcePos { return n.PosAt }
func (n *NullLiteral) statementNode()         {}
func (n *NullLiteral) expressionNode()        {}
func (n *NullLiteral) factorNode()            {}
func (n *NullLiteral) Visit(v Visitor[any]) any {
	return v.VisitNullLiteral(n)
}

type Identifier struct {
	Name  string
	PosAt *common.SourcePos
}

func (i *Identifier) Pos() *common.SourcePos { return i.PosAt }
func (i *Identifier) statementNode()         {}
func (i *Identifier) expressionNode()        {}
func (i *Identifier) factorNode()            {}
func (i *Identifier) Visit(v Visitor[any]) any {
	return v.VisitIdentifier(i)
}

type BinaryExpr struct {
	Left     Expression
	Operator lexer.OperatorSubkind
	Right    Expression
	PosAt    *common.SourcePos
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

type Block struct {
	Statements []Statement
	PosAt      *common.SourcePos
}

func (b *Block) Pos() *common.SourcePos { return b.PosAt }
func (b *Block) statementNode()         {}
func (b *Block) expressionNode()        {}
func (b *Block) Visit(v Visitor[any]) any {
	return v.VisitBlock(b)
}
