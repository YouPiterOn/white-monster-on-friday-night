package ast

import "youpiteron.dev/white-monster-on-friday-night/internal/lexer"

type Program struct {
	Statements []Statement
	PosAt      *lexer.SourcePos
}

func (p *Program) Pos() *lexer.SourcePos { return p.PosAt }
func (p *Program) Visit(v Visitor[any]) any {
	return v.VisitProgram(p)
}
