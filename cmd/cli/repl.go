package main

import (
	"bufio"
	"fmt"
	"os"

	"youpiteron.dev/white-monster-on-friday-night/internal/ast"
	"youpiteron.dev/white-monster-on-friday-night/internal/compiler"
	"youpiteron.dev/white-monster-on-friday-night/internal/lexer"
	"youpiteron.dev/white-monster-on-friday-night/internal/vm"
)

func REPL() {
	reader := bufio.NewReader(os.Stdin)
	lexer := lexer.NewLexer()
	compiler := compiler.NewCompiler()
	globalTable := compiler.StartREPL()
	vm := vm.NewVM(globalTable)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		lexerResult := lexer.Lex(input)
		if len(lexerResult.Errors) > 0 {
			fmt.Println("Error lexing input:", lexerResult.Errors)
			continue
		}
		parser := ast.NewParser(lexerResult.Tokens)
		program := parser.ParseProgram()
		if len(parser.Errors) > 0 {
			fmt.Println("Error parsing input:", parser.Errors)
			continue
		}
		compileResult, _ := compiler.CompileREPLChunk(program)
		retval := vm.RunModuleProto(&compileResult.ModuleProto)
		fmt.Printf("retval: %d\n", retval)
	}
}
