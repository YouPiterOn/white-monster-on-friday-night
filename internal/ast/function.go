package ast

import (
	"youpiteron.dev/white-monster-on-friday-night/internal/common"
)

type Param struct {
	Name   string
	TypeOf *Type
	PosAt  *common.SourcePos
	Vararg bool
}

func (p *Param) Pos() *common.SourcePos { return p.PosAt }
func (p *Param) Visit(v Visitor[any]) any {
	return v.VisitParam(p)
}

type Function struct {
	Name       string
	Params     []Param
	Vararg     bool
	Body       []Statement
	ReturnType *Type
	PosAt      *common.SourcePos
}

func (f *Function) Pos() *common.SourcePos { return f.PosAt }
func (f *Function) statementNode()         {}
func (f *Function) Visit(v Visitor[any]) any {
	return v.VisitFunction(f)
}
