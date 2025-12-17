package compiler

import (
	"fmt"
	"os"

	"youpiteron.dev/white-monster-on-friday-night/internal/ast"
)

type CompileResult struct {
	FunctionProtos []FunctionProto
	GlobalTable    *GlobalTable
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
	for _, functionProto := range c.instructionsVisitor.functionProtos {
		fmt.Printf("function proto: %s\n", functionProto.String())
	}
	return &CompileResult{FunctionProtos: c.instructionsVisitor.functionProtos, GlobalTable: c.instructionsVisitor.globalTable}
}
