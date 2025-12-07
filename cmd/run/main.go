package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"youpiteron.dev/white-monster-on-friday-night/internal/ast"
	"youpiteron.dev/white-monster-on-friday-night/internal/debug"
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

	scanner := bufio.NewScanner(bytes.NewReader(buffer))
	lexer := lexer.NewLexer()

	for scanner.Scan() {
		line := scanner.Text()

		result := lexer.Lex(line)
		parser := ast.NewParser(result.Tokens)
		statement := parser.ParseStatement()

		debugVisitor := debug.DebugVisitor{}

		statement.Visit(&debugVisitor)
	}
}
