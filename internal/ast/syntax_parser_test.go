package ast

import (
	"testing"

	"youpiteron.dev/white-monster-on-friday-night/internal/common"
	"youpiteron.dev/white-monster-on-friday-night/internal/lexer"
)

func makeToken(lexeme string, kind lexer.TokenKind, subkind any, offset, line, col int) lexer.Token {
	return lexer.Token{
		Lexeme:  lexeme,
		Kind:    kind,
		Subkind: subkind,
		Pos: &common.SourcePos{
			BasePos: common.BasePos{
				Offset: offset,
				Line:   line,
				Column: col,
			},
			Length: len(lexeme),
		},
	}
}

// ---------- ParseDeclaration Tests ----------

func TestParseDeclaration_VarWithTypeAndValue(t *testing.T) {
	tokens := []lexer.Token{
		makeToken("var", lexer.Keyword, lexer.KeywordVar, 0, 1, 1),
		makeToken("x", lexer.Identifier, lexer.IdentifierName, 4, 1, 5),
		makeToken(":", lexer.Punctuator, lexer.Colon, 5, 1, 6),
		makeToken("int", lexer.Type, lexer.TypeInt, 6, 1, 7),
		makeToken("=", lexer.Punctuator, lexer.Assign, 10, 1, 11),
		makeToken("42", lexer.Constant, lexer.Integer, 12, 1, 13),
		makeToken(";", lexer.Punctuator, lexer.StatementEnd, 14, 1, 15),
	}

	parser := NewParser(tokens)
	decl := parser.ParseDeclaration()

	if decl == nil {
		t.Fatal("expected declaration but got nil")
	}
	if len(parser.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", parser.Errors)
	}
	if !decl.IsMutable {
		t.Error("expected IsMutable to be true")
	}
	if !decl.IsTyped {
		t.Error("expected IsTyped to be true")
	}
	if decl.TypeOf != lexer.TypeInt {
		t.Errorf("expected TypeOf to be TypeInt, got %v", decl.TypeOf)
	}
	if decl.Identifier.Name != "x" {
		t.Errorf("expected identifier name 'x', got '%s'", decl.Identifier.Name)
	}
	if decl.Value == nil {
		t.Fatal("expected value but got nil")
	}
}

func TestParseDeclaration_ConstWithoutType(t *testing.T) {
	tokens := []lexer.Token{
		makeToken("const", lexer.Keyword, lexer.KeywordConst, 0, 1, 1),
		makeToken("y", lexer.Identifier, lexer.IdentifierName, 6, 1, 7),
		makeToken("=", lexer.Punctuator, lexer.Assign, 8, 1, 9),
		makeToken("true", lexer.Constant, lexer.Boolean, 10, 1, 11),
		makeToken(";", lexer.Punctuator, lexer.StatementEnd, 14, 1, 15),
	}

	parser := NewParser(tokens)
	decl := parser.ParseDeclaration()

	if decl == nil {
		t.Fatal("expected declaration but got nil")
	}
	if len(parser.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", parser.Errors)
	}
	if decl.IsMutable {
		t.Error("expected IsMutable to be false")
	}
	if decl.IsTyped {
		t.Error("expected IsTyped to be false")
	}
	if decl.Identifier.Name != "y" {
		t.Errorf("expected identifier name 'y', got '%s'", decl.Identifier.Name)
	}
}

func TestParseDeclaration_VarWithoutValue(t *testing.T) {
	tokens := []lexer.Token{
		makeToken("var", lexer.Keyword, lexer.KeywordVar, 0, 1, 1),
		makeToken("z", lexer.Identifier, lexer.IdentifierName, 4, 1, 5),
		makeToken(":", lexer.Punctuator, lexer.Colon, 5, 1, 6),
		makeToken("bool", lexer.Type, lexer.TypeBool, 6, 1, 7),
		makeToken(";", lexer.Punctuator, lexer.StatementEnd, 10, 1, 11),
	}

	parser := NewParser(tokens)
	decl := parser.ParseDeclaration()

	if decl == nil {
		t.Fatal("expected declaration but got nil")
	}
	if len(parser.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", parser.Errors)
	}
	if decl.Value != nil {
		t.Error("expected Value to be nil")
	}
	if decl.TypeOf != lexer.TypeBool {
		t.Errorf("expected TypeOf to be TypeBool, got %v", decl.TypeOf)
	}
}

func TestParseDeclaration_MissingIdentifier(t *testing.T) {
	tokens := []lexer.Token{
		makeToken("var", lexer.Keyword, lexer.KeywordVar, 0, 1, 1),
		makeToken(";", lexer.Punctuator, lexer.StatementEnd, 4, 1, 5),
	}

	parser := NewParser(tokens)
	decl := parser.ParseDeclaration()

	if decl != nil {
		t.Error("expected nil declaration due to missing identifier")
	}
	if len(parser.Errors) == 0 {
		t.Error("expected errors but got none")
	}
}

// ---------- ParseFunction Tests ----------

func TestParseFunction_Basic(t *testing.T) {
	tokens := []lexer.Token{
		makeToken("function", lexer.Keyword, lexer.KeywordFunction, 0, 1, 1),
		makeToken("add", lexer.Identifier, lexer.IdentifierName, 9, 1, 10),
		makeToken("(", lexer.Punctuator, lexer.ParenOpen, 12, 1, 13),
		makeToken("x", lexer.Identifier, lexer.IdentifierName, 13, 1, 14),
		makeToken(":", lexer.Punctuator, lexer.Colon, 14, 1, 15),
		makeToken("int", lexer.Type, lexer.TypeInt, 15, 1, 16),
		makeToken(")", lexer.Punctuator, lexer.ParenClose, 18, 1, 19),
		makeToken(":", lexer.Punctuator, lexer.Colon, 19, 1, 20),
		makeToken("int", lexer.Type, lexer.TypeInt, 20, 1, 21),
		makeToken("{", lexer.Punctuator, lexer.BlockStart, 23, 1, 24),
		makeToken("}", lexer.Punctuator, lexer.BlockEnd, 24, 1, 25),
	}

	parser := NewParser(tokens)
	fnStmt := parser.ParseFunction()

	if fnStmt == nil {
		t.Fatal("expected function but got nil")
	}
	if len(parser.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", parser.Errors)
	}
	fn, ok := fnStmt.(*Function)
	if !ok {
		t.Fatalf("expected *Function, got %T", fnStmt)
	}
	if fn.Name != "add" {
		t.Errorf("expected function name 'add', got '%s'", fn.Name)
	}
	if len(fn.Params) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(fn.Params))
	}
	if fn.ReturnType != lexer.TypeInt {
		t.Errorf("expected return type TypeInt, got %v", fn.ReturnType)
	}
}

// ---------- ParseReturn Tests ----------

func TestParseReturn_WithValue(t *testing.T) {
	tokens := []lexer.Token{
		makeToken("return", lexer.Keyword, lexer.KeywordReturn, 0, 1, 1),
		makeToken("42", lexer.Constant, lexer.Integer, 7, 1, 8),
		makeToken(";", lexer.Punctuator, lexer.StatementEnd, 9, 1, 10),
	}

	parser := NewParser(tokens)
	retStmt := parser.ParseReturn()

	if retStmt == nil {
		t.Fatal("expected return statement but got nil")
	}
	if len(parser.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", parser.Errors)
	}
	ret, ok := retStmt.(*Return)
	if !ok {
		t.Fatalf("expected *Return, got %T", retStmt)
	}
	if ret.Value == nil {
		t.Fatal("expected return value but got nil")
	}
}

// ---------- ParseAssignment Tests ----------

func TestParseAssignment_Basic(t *testing.T) {
	tokens := []lexer.Token{
		makeToken("x", lexer.Identifier, lexer.IdentifierName, 0, 1, 1),
		makeToken("=", lexer.Punctuator, lexer.Assign, 2, 1, 3),
		makeToken("10", lexer.Constant, lexer.Integer, 4, 1, 5),
		makeToken(";", lexer.Punctuator, lexer.StatementEnd, 6, 1, 7),
	}

	parser := NewParser(tokens)
	assign := parser.ParseAssignment()

	if assign == nil {
		t.Fatal("expected assignment but got nil")
	}
	if len(parser.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", parser.Errors)
	}
	if assign.Identifier.Name != "x" {
		t.Errorf("expected identifier name 'x', got '%s'", assign.Identifier.Name)
	}
	if assign.Value == nil {
		t.Fatal("expected value but got nil")
	}
}

// ---------- ParseBlock Tests ----------

func TestParseBlock_Empty(t *testing.T) {
	tokens := []lexer.Token{
		makeToken("{", lexer.Punctuator, lexer.BlockStart, 0, 1, 1),
		makeToken("}", lexer.Punctuator, lexer.BlockEnd, 1, 1, 2),
	}

	parser := NewParser(tokens)
	block := parser.ParseBlock()

	if block == nil {
		t.Fatal("expected block but got nil")
	}
	if len(parser.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", parser.Errors)
	}
	if len(block.Statements) != 0 {
		t.Errorf("expected 0 statements, got %d", len(block.Statements))
	}
}

// ---------- ParseExpression Tests ----------

func TestParseExpression_Identifier(t *testing.T) {
	tokens := []lexer.Token{
		makeToken("x", lexer.Identifier, lexer.IdentifierName, 0, 1, 1),
	}

	parser := NewParser(tokens)
	expr := parser.ParseExpression(false)

	if expr == nil {
		t.Fatal("expected expression but got nil")
	}
	if len(parser.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", parser.Errors)
	}
}

// ---------- ParseCallExpr Tests ----------

func TestParseCallExpr_NoArguments(t *testing.T) {
	tokens := []lexer.Token{
		makeToken("foo", lexer.Identifier, lexer.IdentifierName, 0, 1, 1),
		makeToken("(", lexer.Punctuator, lexer.ParenOpen, 3, 1, 4),
		makeToken(")", lexer.Punctuator, lexer.ParenClose, 4, 1, 5),
	}

	parser := NewParser(tokens)
	callExpr := parser.ParseCallExpr(false)

	if callExpr == nil {
		t.Fatal("expected call expression but got nil")
	}
	if len(parser.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", parser.Errors)
	}
	if callExpr.Identifier.Name != "foo" {
		t.Errorf("expected identifier name 'foo', got '%s'", callExpr.Identifier.Name)
	}
	if len(callExpr.Arguments) != 0 {
		t.Errorf("expected 0 arguments, got %d", len(callExpr.Arguments))
	}
}

// ---------- ParseFactor Tests ----------

func TestParseMultiplicativeExpr_IntLiteral(t *testing.T) {
	tokens := []lexer.Token{
		makeToken("123", lexer.Constant, lexer.Integer, 0, 1, 1),
	}

	parser := NewParser(tokens)
	factor := parser.ParseMultiplicativeExpr(false)

	if factor == nil {
		t.Fatal("expected factor but got nil")
	}
	if len(parser.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", parser.Errors)
	}
	intLit, ok := factor.(*IntLiteral)
	if !ok {
		t.Fatalf("expected IntLiteral, got %T", factor)
	}
	if intLit.Value != 123 {
		t.Errorf("expected value 123, got %d", intLit.Value)
	}
}

// ---------- ParseIntLiteral Tests ----------

func TestParseIntLiteral_Basic(t *testing.T) {
	tokens := []lexer.Token{
		makeToken("456", lexer.Constant, lexer.Integer, 0, 1, 1),
	}

	parser := NewParser(tokens)
	lit := parser.ParseIntLiteral(false)

	if lit == nil {
		t.Fatal("expected int literal but got nil")
	}
	if len(parser.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", parser.Errors)
	}
	if lit.Value != 456 {
		t.Errorf("expected value 456, got %d", lit.Value)
	}
}

// ---------- ParseBoolLiteral Tests ----------

func TestParseBoolLiteral_True(t *testing.T) {
	tokens := []lexer.Token{
		makeToken("true", lexer.Constant, lexer.Boolean, 0, 1, 1),
	}

	parser := NewParser(tokens)
	lit := parser.ParseBoolLiteral(false)

	if lit == nil {
		t.Fatal("expected bool literal but got nil")
	}
	if len(parser.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", parser.Errors)
	}
	if !lit.Value {
		t.Error("expected value to be true")
	}
}

// ---------- ParseNullLiteral Tests ----------

func TestParseNullLiteral_Basic(t *testing.T) {
	tokens := []lexer.Token{
		makeToken("null", lexer.Constant, lexer.Null, 0, 1, 1),
	}

	parser := NewParser(tokens)
	lit := parser.ParseNullLiteral(false)

	if lit == nil {
		t.Fatal("expected null literal but got nil")
	}
	if len(parser.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", parser.Errors)
	}
}

// ---------- ParseIdentifier Tests ----------

func TestParseIdentifier_Basic(t *testing.T) {
	tokens := []lexer.Token{
		makeToken("myVar", lexer.Identifier, lexer.IdentifierName, 0, 1, 1),
	}

	parser := NewParser(tokens)
	id := parser.ParseIdentifier(false)

	if id == nil {
		t.Fatal("expected identifier but got nil")
	}
	if len(parser.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", parser.Errors)
	}
	if id.Name != "myVar" {
		t.Errorf("expected name 'myVar', got '%s'", id.Name)
	}
}

// ---------- ParseParam Tests ----------

func TestParseParam_Basic(t *testing.T) {
	tokens := []lexer.Token{
		makeToken("x", lexer.Identifier, lexer.IdentifierName, 0, 1, 1),
		makeToken(":", lexer.Punctuator, lexer.Colon, 1, 1, 2),
		makeToken("int", lexer.Type, lexer.TypeInt, 2, 1, 3),
	}

	parser := NewParser(tokens)
	param := parser.ParseParam()

	if param == nil {
		t.Fatal("expected parameter but got nil")
	}
	if len(parser.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", parser.Errors)
	}
	if param.Name != "x" {
		t.Errorf("expected parameter name 'x', got '%s'", param.Name)
	}
	if param.TypeOf != lexer.TypeInt {
		t.Errorf("expected type TypeInt, got %v", param.TypeOf)
	}
}

// ---------- ParseBody Tests ----------

func TestParseBody_WithStatement(t *testing.T) {
	tokens := []lexer.Token{
		makeToken("{", lexer.Punctuator, lexer.BlockStart, 0, 1, 1),
		makeToken("return", lexer.Keyword, lexer.KeywordReturn, 2, 1, 3),
		makeToken("1", lexer.Constant, lexer.Integer, 9, 1, 10),
		makeToken(";", lexer.Punctuator, lexer.StatementEnd, 10, 1, 11),
		makeToken("}", lexer.Punctuator, lexer.BlockEnd, 12, 1, 13),
	}

	parser := NewParser(tokens)
	body := parser.ParseBody()

	if body == nil {
		t.Fatal("expected body but got nil")
	}
	if len(parser.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", parser.Errors)
	}
	if len(body) != 1 {
		t.Errorf("expected 1 statement, got %d", len(body))
	}
}

// ---------- ParseProgram Tests ----------

func TestParseProgram_Empty(t *testing.T) {
	tokens := []lexer.Token{}

	parser := NewParser(tokens)
	program := parser.ParseProgram()

	if program == nil {
		t.Fatal("expected program but got nil")
	}
	if len(program.Statements) != 0 {
		t.Errorf("expected 0 statements, got %d", len(program.Statements))
	}
}

// ---------- ParseStatement Tests ----------

func TestParseStatement_Expression(t *testing.T) {
	tokens := []lexer.Token{
		makeToken("42", lexer.Constant, lexer.Integer, 0, 1, 1),
		makeToken(";", lexer.Punctuator, lexer.StatementEnd, 2, 1, 3),
	}

	parser := NewParser(tokens)
	stmt := parser.ParseStatement()

	if stmt == nil {
		t.Fatal("expected statement but got nil")
	}
	if len(parser.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", parser.Errors)
	}
	ret, ok := stmt.(*IntLiteral)
	if !ok {
		t.Fatalf("expected IntLiteral statement, got %T", stmt)
	}
	if ret.Value != 42 {
		t.Errorf("expected value 42, got %d", ret.Value)
	}
}
