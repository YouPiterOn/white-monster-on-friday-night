package ast

import "youpiteron.dev/white-monster-on-friday-night/internal/lexer"

type Parameter struct {
	Name  string
	PosAt *lexer.SourcePos
}

func (p *Parameter) Pos() *lexer.SourcePos { return p.PosAt }
func (p *Parameter) Visit(v Visitor[any]) any {
	return v.VisitParameter(p)
}

type Function struct {
	Name   string
	Params []Parameter
	Body   Block
	PosAt  *lexer.SourcePos
}

func (f *Function) Pos() *lexer.SourcePos { return f.PosAt }
func (f *Function) statementNode()        {}
func (f *Function) Visit(v Visitor[any]) any {
	return v.VisitFunction(f)
}
