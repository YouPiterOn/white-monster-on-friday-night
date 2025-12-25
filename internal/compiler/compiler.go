package compiler

import (
	"fmt"
	"os"

	"youpiteron.dev/white-monster-on-friday-night/internal/ast"
	"youpiteron.dev/white-monster-on-friday-night/internal/common"
)

type CompileResult struct {
	ModuleProto ModuleProto
	GlobalTable *GlobalTable
}

type Compiler struct {
	replMode            bool
	instructionsVisitor *InstructionsVisitor
}

func NewCompiler() *Compiler {
	return &Compiler{replMode: false, instructionsVisitor: NewInstructionsVisitor()}
}

func (c *Compiler) CompileToModuleProto(program *ast.Program) *CompileResult {
	c.instructionsVisitor.EnterModuleContext()
	program.Visit(c.instructionsVisitor)
	c.instructionsVisitor.ExitModuleContext()
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

func (c *Compiler) StartREPL() *GlobalTable {
	c.replMode = true
	c.instructionsVisitor.EnterModuleContext()
	return c.instructionsVisitor.globalTable
}

func (c *Compiler) EndREPL() {
	if !c.replMode {
		panic("COMPILER ERROR: cannot end REPL mode without starting it")
	}
	c.replMode = false
	c.instructionsVisitor.ExitModuleContext()
}

func (c *Compiler) CompileREPLChunk(program *ast.Program) (*CompileResult, []common.Error) {
	if !c.replMode {
		panic("COMPILER ERROR: cannot compile REPL chunk without starting REPL mode")
	}
	program.Visit(c.instructionsVisitor)
	if len(c.instructionsVisitor.errors) > 0 {
		fmt.Printf("COMPILATION ERROR: failed to generate instructions from program\n")
		for _, error := range c.instructionsVisitor.errors {
			fmt.Printf("  %s at %v\n", error.Message, error.Pos)
		}
		return nil, c.instructionsVisitor.errors
	}
	moduleProto := c.instructionsVisitor.EmitModuleProto()
	return &CompileResult{ModuleProto: *moduleProto, GlobalTable: c.instructionsVisitor.globalTable}, nil
}
