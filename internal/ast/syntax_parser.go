package ast

import (
	"fmt"

	"youpiteron.dev/white-monster-on-friday-night/internal/lexer"
)

type SyntaxError struct {
	Message string
	Pos     *lexer.SourcePos
}

type Parser struct {
	tokens []lexer.Token
	idx    int
	Errors []SyntaxError
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) peek(offset int) *lexer.Token {
	i := p.idx + offset
	if i >= len(p.tokens) {
		return nil
	}
	return &p.tokens[i]
}

func (p *Parser) eat() *lexer.Token {
	if p.idx >= len(p.tokens) {
		return nil
	}
	t := &p.tokens[p.idx]
	p.idx++
	return t
}

func (p *Parser) eatExpected(kind lexer.TokenKind, subkind any, msg string) *lexer.Token {
	t := p.peek(0)
	if t == nil {
		p.addError(msg, nil)
		return nil
	}
	if t.Kind == kind && (subkind == nil || t.Subkind == subkind) {
		return p.eat()
	}

	got := fmt.Sprintf("%v(%v)", t.Kind, t.Subkind)
	want := fmt.Sprintf("%v(%v)", kind, subkind)
	p.addError(fmt.Sprintf("expected %s but got %s", want, got), t.Pos)
	return nil
}

func (p *Parser) addError(msg string, pos *lexer.SourcePos) {
	p.Errors = append(p.Errors, SyntaxError{
		Message: msg,
		Pos:     pos,
	})
}

func (p *Parser) ParseProgram() *Program {
	statements := []Statement{}
	pos := p.peek(0).Pos
	for {
		statement := p.ParseStatement()
		if statement == nil {
			break
		}
		statements = append(statements, statement)
	}
	return &Program{Statements: statements, PosAt: pos}
}

func (p *Parser) ParseStatement() Statement {
	t := p.peek(0)
	if t == nil {
		return nil
	}

	if t.Kind == lexer.Punctuator && t.Subkind == lexer.BlockStart {
		return p.ParseBlock()
	}

	if t.Kind == lexer.Keyword && t.Subkind == lexer.KeywordFunction {
		return p.ParseFunction()
	}

	if t.Kind == lexer.Keyword && t.Subkind == lexer.KeywordReturn {
		return p.ParseReturn()
	}

	if t.Kind == lexer.Keyword && (t.Subkind == lexer.KeywordVar || t.Subkind == lexer.KeywordConst) {
		return p.ParseAssignment()
	}
	expression := p.ParseExpression()
	if expression == nil {
		return nil
	}
	semicolon := p.eatExpected(lexer.Punctuator, lexer.StatementEnd, "expected ';'")
	if semicolon == nil {
		return nil
	}
	return &Return{Value: expression, PosAt: semicolon.Pos}
}

func (p *Parser) ParseFunction() Statement {
	kw := p.eatExpected(lexer.Keyword, lexer.KeywordFunction, "expected 'function'")
	if kw == nil {
		return nil
	}

	idTok := p.eatExpected(lexer.Identifier, lexer.IdentifierName, "expected identifier")
	if idTok == nil {
		return nil
	}

	lparen := p.eatExpected(lexer.Punctuator, lexer.ParenOpen, "expected '('")
	if lparen == nil {
		return nil
	}

	params := []Parameter{}
	for {
		param := p.ParseParameter()
		if param == nil {
			break
		}
		params = append(params, *param)
		commaTok := p.peek(0)
		if commaTok == nil || !(commaTok.Kind == lexer.Punctuator && commaTok.Subkind == lexer.Comma) {
			break
		}
		p.eat()
	}
	rparen := p.eatExpected(lexer.Punctuator, lexer.ParenClose, "expected ')'")
	if rparen == nil {
		return nil
	}

	body := p.ParseBlock()
	if body == nil {
		return nil
	}

	return &Function{Name: idTok.Lexeme, Params: params, Body: *body, PosAt: kw.Pos}
}

func (p *Parser) ParseParameter() *Parameter {

	idTok := p.eatExpected(lexer.Identifier, lexer.IdentifierName, "expected identifier")
	if idTok == nil {
		return nil
	}
	return &Parameter{Name: idTok.Lexeme, PosAt: idTok.Pos}
}

func (p *Parser) ParseBlock() *Block {
	lbrace := p.eatExpected(lexer.Punctuator, lexer.BlockStart, "expected '{'")
	if lbrace == nil {
		return nil
	}
	statements := []Statement{}
	for {
		t := p.peek(0)
		if t == nil {
			break
		}
		if t.Kind == lexer.Punctuator && t.Subkind == lexer.BlockEnd {
			break
		}
		statement := p.ParseStatement()
		if statement == nil {
			break
		}
		statements = append(statements, statement)
	}
	rbrace := p.eatExpected(lexer.Punctuator, lexer.BlockEnd, "expected '}'")
	if rbrace == nil {
		return nil
	}
	return &Block{Statements: statements, PosAt: lbrace.Pos}
}

func (p *Parser) ParseReturn() Statement {
	kw := p.eatExpected(lexer.Keyword, lexer.KeywordReturn, "expected 'return'")
	if kw == nil {
		return nil
	}

	value := p.ParseExpression()

	semicolon := p.eatExpected(lexer.Punctuator, lexer.StatementEnd, "expected ';'")
	if semicolon == nil {
		return nil
	}

	return &Return{
		Value: value,
		PosAt: kw.Pos,
	}
}

func (p *Parser) ParseAssignment() Statement {
	kw := p.eatExpected(lexer.Keyword, nil, "expected 'var' or 'const'")
	if kw == nil {
		return nil
	}

	specifier := kw.Subkind.(lexer.KeywordSubkind)
	if specifier != lexer.KeywordVar && specifier != lexer.KeywordConst {
		p.addError(fmt.Sprintf("expected 'var' or 'const' but got %v(%v)", kw.Kind, kw.Subkind), kw.Pos)
		return nil
	}

	idTok := p.eatExpected(lexer.Identifier, lexer.IdentifierName, "expected identifier")
	if idTok == nil {
		return nil
	}

	eq := p.eatExpected(lexer.Punctuator, lexer.Assign, "expected '='")
	if eq == nil {
		return nil
	}

	value := p.ParseExpression()
	if value == nil {
		return nil
	}

	semicolon := p.eatExpected(lexer.Punctuator, lexer.StatementEnd, "expected ';'")
	if semicolon == nil {
		return nil
	}

	return &Assignment{
		Specifier:  specifier,
		Identifier: &Identifier{Name: idTok.Lexeme, PosAt: idTok.Pos},
		Value:      value,
		PosAt:      idTok.Pos,
	}
}

func (p *Parser) ParseExpression() Expression {
	next := p.peek(1)
	if next != nil && next.Kind == lexer.Operator {
		return p.ParseBinaryExpr()
	}
	if next != nil && next.Kind == lexer.Punctuator && next.Subkind == lexer.ParenOpen {
		return p.ParseCallExpr()
	}
	return p.ParseFactor()
}

func (p *Parser) ParseCallExpr() Expression {
	identifier := p.ParseIdentifier()
	if identifier == nil {
		return nil
	}
	lparen := p.eatExpected(lexer.Punctuator, lexer.ParenOpen, "expected '('")
	if lparen == nil {
		return nil
	}
	arguments := []Expression{}
	for {
		argument := p.ParseExpression()
		if argument == nil {
			break
		}
		arguments = append(arguments, argument)
		commaTok := p.peek(0)
		if commaTok == nil || !(commaTok.Kind == lexer.Punctuator && commaTok.Subkind == lexer.Comma) {
			break
		}
		p.eat()
	}
	rparen := p.eatExpected(lexer.Punctuator, lexer.ParenClose, "expected ')'")
	if rparen == nil {
		return nil
	}
	return &CallExpr{Function: *identifier, Arguments: arguments, PosAt: lparen.Pos}
}

func (p *Parser) ParseFactor() Expression {
	t := p.peek(0)
	if t == nil {
		return nil
	}

	if t.Kind == lexer.Constant && t.Subkind == lexer.Numeric {
		return p.ParseNumberLiteral()
	}

	if t.Kind == lexer.Identifier {
		return p.ParseIdentifier()
	}

	p.addError(fmt.Sprintf("expected factor but got %v(%v)", t.Kind, t.Subkind), t.Pos)
	return nil
}

func (p *Parser) ParseNumberLiteral() Expression {
	t := p.eat()
	if t == nil {
		return nil
	}

	return &NumberLiteral{
		Value: atoi(t.Lexeme),
		PosAt: t.Pos,
	}
}

func (p *Parser) ParseIdentifier() *Identifier {
	t := p.eat()
	if t == nil {
		return nil
	}

	return &Identifier{Name: t.Lexeme, PosAt: t.Pos}
}

func (p *Parser) ParseBinaryExpr() Expression {
	left := p.ParseFactor()
	if left == nil {
		return nil
	}

	op := p.eat()
	if op == nil || op.Kind != lexer.Operator {
		return nil
	}

	right := p.ParseExpression()
	if right == nil {
		return nil
	}

	return &BinaryExpr{
		Left:     left,
		Operator: op.Subkind.(lexer.OperatorSubkind),
		Right:    right,
		PosAt:    op.Pos,
	}
}

func atoi(s string) int {
	var n int
	for i := 0; i < len(s); i++ {
		n = n*10 + int(s[i]-'0')
	}
	return n
}
