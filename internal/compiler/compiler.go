package compiler

import "youpiteron.dev/white-monster-on-friday-night/internal/ast"

type Compiler struct {
	program *ast.Program
	scope   *Scope
}

func NewCompiler(program *ast.Program) *Compiler {
	return &Compiler{program: program, scope: NewScope()}
}
