package ast

import "youpiteron.dev/white-monster-on-friday-night/internal/common"

type ObjectStatement struct {
	Properties   []*ObjectProperty
	IsTypeObject bool
	Extends      []*Identifier
	Implements   []*TypeIdentifier
	PosAt        *common.SourcePos
}

func (o *ObjectStatement) Pos() *common.SourcePos { return o.PosAt }
func (o *ObjectStatement) statementNode()         {}
func (o *ObjectStatement) Visit(v Visitor[any]) any {
	return v.VisitObjectStatement(o)
}

type ObjectProperty struct {
	Name   string
	TypeOf *Type
	PosAt  *common.SourcePos
}
