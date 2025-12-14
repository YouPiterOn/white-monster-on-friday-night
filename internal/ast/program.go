package ast

import "youpiteron.dev/white-monster-on-friday-night/internal/common"

type Program struct {
	Statements []Statement
	PosAt      *common.SourcePos
}

func (p *Program) Pos() *common.SourcePos { return p.PosAt }
func (p *Program) Visit(v Visitor[any]) any {
	return v.VisitProgram(p)
}
