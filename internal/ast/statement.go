package ast

import "youpiteron.dev/white-monster-on-friday-night/internal/lexer"

type Statement interface {
	Node
	statementNode()
}

type Declaration struct {
	Specifier  lexer.KeywordSubkind
	Identifier *Identifier
	Value      Expression
	PosAt      *lexer.SourcePos
}

func (d *Declaration) Pos() *lexer.SourcePos { return d.PosAt }
func (d *Declaration) statementNode()        {}
func (d *Declaration) Visit(v Visitor[any]) any {
	return v.VisitDeclaration(d)
}

type Assignment struct {
	Identifier *Identifier
	Value      Expression
	PosAt      *lexer.SourcePos
}

func (a *Assignment) Pos() *lexer.SourcePos { return a.PosAt }
func (a *Assignment) statementNode()        {}
func (a *Assignment) Visit(v Visitor[any]) any {
	return v.VisitAssignment(a)
}

type Return struct {
	Value Expression
	PosAt *lexer.SourcePos
}

func (r *Return) Pos() *lexer.SourcePos { return r.PosAt }
func (r *Return) statementNode()        {}
func (r *Return) Visit(v Visitor[any]) any {
	return v.VisitReturn(r)
}
