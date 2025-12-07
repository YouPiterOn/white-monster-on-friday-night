package ast

import "youpiteron.dev/white-monster-on-friday-night/internal/lexer"

type Node interface {
	Pos() *lexer.SourcePos
	Visit(v Visitor[any]) any
}
