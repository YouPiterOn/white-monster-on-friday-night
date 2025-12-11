package ast

import "youpiteron.dev/white-monster-on-friday-night/internal/lexer"

type Function struct {
	Name   string
	Params []Expression
	Body   Block
	PosAt  *lexer.SourcePos
}

func (f *Function) Pos() *lexer.SourcePos { return f.PosAt }
func (f *Function) statementNode()        {}
func (f *Function) Visit(v Visitor[any]) any {
	return v.VisitFunction(f)
}
