package ast

import (
	"youpiteron.dev/white-monster-on-friday-night/internal/common"
)

type Statement interface {
	Node
	statementNode()
}

type Declaration struct {
	IsMutable  bool
	IsTyped    bool
	TypeOf     *Type
	Identifier *Identifier
	Value      Expression
	PosAt      *common.SourcePos
}

func (d *Declaration) Pos() *common.SourcePos { return d.PosAt }
func (d *Declaration) statementNode()         {}
func (d *Declaration) Visit(v Visitor[any]) any {
	return v.VisitDeclaration(d)
}

type Assignment struct {
	Identifier *Identifier
	Value      Expression
	PosAt      *common.SourcePos
}

func (a *Assignment) Pos() *common.SourcePos { return a.PosAt }
func (a *Assignment) statementNode()         {}
func (a *Assignment) Visit(v Visitor[any]) any {
	return v.VisitAssignment(a)
}

type Block struct {
	Statements []Statement
	PosAt      *common.SourcePos
}

func (b *Block) Pos() *common.SourcePos { return b.PosAt }
func (b *Block) statementNode()         {}
func (b *Block) Visit(v Visitor[any]) any {
	return v.VisitBlock(b)
}

type Return struct {
	Value Expression
	PosAt *common.SourcePos
}

func (r *Return) Pos() *common.SourcePos { return r.PosAt }
func (r *Return) statementNode()         {}
func (r *Return) Visit(v Visitor[any]) any {
	return v.VisitReturn(r)
}

type If struct {
	Condition Expression
	Body      []Statement
	ElseBody  []Statement
	PosAt     *common.SourcePos
}

func (i *If) Pos() *common.SourcePos { return i.PosAt }
func (i *If) statementNode()         {}
func (i *If) Visit(v Visitor[any]) any {
	return v.VisitIf(i)
}
