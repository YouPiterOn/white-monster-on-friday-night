package compiler

import (
	"fmt"
	"os"

	"youpiteron.dev/white-monster-on-friday-night/internal/ast"
)

type CompileResult struct {
	ModuleProto ModuleProto
	GlobalTable *GlobalTable
}

type Compiler struct {
	instructionsVisitor *InstructionsVisitor
}

func NewCompiler() *Compiler {
	return &Compiler{instructionsVisitor: NewInstructionsVisitor()}
}

func (c *Compiler) Compile(program *ast.Program) *CompileResult {
	program.Visit(c.instructionsVisitor)
	if len(c.instructionsVisitor.errors) > 0 {
		fmt.Printf("COMPILATION ERROR: failed to generate instructions from program\n")
		for _, error := range c.instructionsVisitor.errors {
			fmt.Printf("  %s at %v\n", error.Message, error.Pos)
		}
		os.Exit(1)
	}
	moduleProto := c.instructionsVisitor.moduleProtos[len(c.instructionsVisitor.moduleProtos)-1]
	fmt.Printf("module proto: %s\n", moduleProto.String())
	return &CompileResult{ModuleProto: moduleProto, GlobalTable: c.instructionsVisitor.globalTable}
}
