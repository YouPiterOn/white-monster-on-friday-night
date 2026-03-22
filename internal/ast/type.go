package ast

import (
	"youpiteron.dev/white-monster-on-friday-night/internal/common"
)

type TypeDeclaration struct {
	Name       string
	Properties []TypeProperty
	PosAt      *common.SourcePos
}

func (t *TypeDeclaration) Pos() *common.SourcePos { return t.PosAt }
func (t *TypeDeclaration) statementNode()         {}
func (t *TypeDeclaration) Visit(v Visitor[any]) any {
	return v.VisitTypeDeclaration(t)
}

type TypeProperty struct {
	Name   string
	TypeOf TypeExpression
	PosAt  *common.SourcePos
}
