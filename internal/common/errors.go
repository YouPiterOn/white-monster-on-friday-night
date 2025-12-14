package common

import "fmt"

type Error struct {
	Message string
	Pos     *SourcePos
}

func (e *Error) String() string {
	return fmt.Sprintf("Error: %s at %v", e.Message, e.Pos)
}
