package ast

import "youpiteron.dev/white-monster-on-friday-night/internal/common"

type TypeExpression interface {
	Node
	typeExpressionNode()
}

type TypeIdentifier struct {
	Name    string
	IsArray bool
	PosAt   *common.SourcePos
}

func (t *TypeIdentifier) Pos() *common.SourcePos { return t.PosAt }
func (t *TypeIdentifier) typeExpressionNode()    {}
func (t *TypeIdentifier) Visit(v Visitor[any]) any {
	return v.VisitTypeIdentifier(t)
}
