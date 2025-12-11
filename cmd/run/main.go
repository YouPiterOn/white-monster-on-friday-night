package main

import (
	"fmt"
	"os"

	"youpiteron.dev/white-monster-on-friday-night/internal/ast"
	"youpiteron.dev/white-monster-on-friday-night/internal/compiler"
	"youpiteron.dev/white-monster-on-friday-night/internal/lexer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: run <file>")
		os.Exit(1)
	}

	path := os.Args[1]

	buffer, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("failed to read file %s: %v\n", path, err)
		os.Exit(1)
	}

	fmt.Printf("read %d bytes from %s\n", len(buffer), path)
	fmt.Printf("buffer: %s\n", string(buffer))

	lexer := lexer.NewLexer()

	lexerResult := lexer.Lex(string(buffer))
	if len(lexerResult.Errors) > 0 {
		fmt.Printf("failed to lex file %s: %v\n", path, lexerResult.Errors)
		os.Exit(1)
	}

	for _, token := range lexerResult.Tokens {
		fmt.Printf("token: %s\n", token.String())
	}

	parser := ast.NewParser(lexerResult.Tokens)
	program := parser.ParseProgram()
	if len(parser.Errors) > 0 {
		fmt.Printf("failed to parse tokens from file %s\n", path)
		for _, error := range parser.Errors {
			fmt.Printf("  %s at %v\n", error.Message, error.Pos)
		}
		os.Exit(1)
	}

	instructionsVisitor := compiler.NewInstructionsVisitor()
	program.Visit(instructionsVisitor)
	if len(instructionsVisitor.Errors()) > 0 {
		fmt.Printf("failed to generate instructions from file %s\n", path)
		for _, error := range instructionsVisitor.Errors() {
			fmt.Printf("  %s at %v\n", error.Message, error.Pos)
		}
		os.Exit(1)
	}
	for _, functionProto := range instructionsVisitor.FunctionProtos() {
		fmt.Printf("function proto: %s\n", functionProto.String())
	}
}
