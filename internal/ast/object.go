package ast

import "youpiteron.dev/white-monster-on-friday-night/internal/common"

type ObjectDeclaration struct {
	Name         string
	Properties   []*ObjectProperty
	IsTypeObject bool
	Extends      []*Identifier
	Implements   []TypeExpression
	PosAt        *common.SourcePos
}

func (o *ObjectDeclaration) Pos() *common.SourcePos { return o.PosAt }
func (o *ObjectDeclaration) statementNode()         {}
func (o *ObjectDeclaration) Visit(v Visitor[any]) any {
	return v.VisitObjectDeclaration(o)
}

type ObjectProperty struct {
	Name  string
	Value Expression
	PosAt *common.SourcePos
}
