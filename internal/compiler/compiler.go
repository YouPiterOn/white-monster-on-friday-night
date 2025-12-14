package compiler

import (
	"fmt"
	"os"

	"youpiteron.dev/white-monster-on-friday-night/internal/ast"
)

type Compiler struct {
	instructionsVisitor *InstructionsVisitor
}

func NewCompiler() *Compiler {
	return &Compiler{instructionsVisitor: NewInstructionsVisitor()}
}

func (c *Compiler) Compile(program *ast.Program) []FunctionProto {
	program.Visit(c.instructionsVisitor)
	if len(c.instructionsVisitor.Errors()) > 0 {
		fmt.Printf("COMPILATION ERROR: failed to generate instructions from program\n")
		for _, error := range c.instructionsVisitor.Errors() {
			fmt.Printf("  %s at %v\n", error.Message, error.Pos)
		}
		os.Exit(1)
	}
	return c.instructionsVisitor.FunctionProtos()
}
