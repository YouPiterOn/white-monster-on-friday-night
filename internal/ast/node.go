package ast

import "youpiteron.dev/white-monster-on-friday-night/internal/common"

type Node interface {
	Pos() *common.SourcePos
	Visit(v Visitor[any]) any
}
