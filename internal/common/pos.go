package common

import "fmt"

type BasePos struct {
	Offset int
	Line   int
	Column int
}

type SourcePos struct {
	BasePos
	Length int
}

func (p SourcePos) String() string {
	return fmt.Sprintf("Offset: %d, Line: %d, Column: %d, Length: %d", p.Offset, p.Line, p.Column, p.Length)
}
