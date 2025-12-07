package ast

import "youpiteron.dev/white-monster-on-friday-night/internal/lexer"

type Expression interface {
	Node
	expressionNode()
}

type Factor interface {
	Expression
	factorNode()
}

type NumberLiteral struct {
	Value int
	PosAt *lexer.SourcePos
}

func (n *NumberLiteral) Pos() *lexer.SourcePos { return n.PosAt }
func (n *NumberLiteral) expressionNode()       {}
func (n *NumberLiteral) factorNode()           {}
func (n *NumberLiteral) Visit(v Visitor[any]) any {
	return v.VisitNumberLiteral(n)
}

type Identifier struct {
	Name  string
	PosAt *lexer.SourcePos
}

func (i *Identifier) Pos() *lexer.SourcePos { return i.PosAt }
func (i *Identifier) expressionNode()       {}
func (i *Identifier) factorNode()           {}
func (i *Identifier) Visit(v Visitor[any]) any {
	return v.VisitIdentifier(i)
}

type BinaryExpr struct {
	Left     Expression
	Operator lexer.OperatorSubkind
	Right    Expression
	PosAt    *lexer.SourcePos
}

func (b *BinaryExpr) Pos() *lexer.SourcePos { return b.PosAt }
func (b *BinaryExpr) expressionNode()       {}
func (b *BinaryExpr) Visit(v Visitor[any]) any {
	return v.VisitBinaryExpr(b)
}
