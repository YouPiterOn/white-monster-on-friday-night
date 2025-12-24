package ast

import (
	"fmt"

	"youpiteron.dev/white-monster-on-friday-night/internal/common"
	"youpiteron.dev/white-monster-on-friday-night/internal/lexer"
)

type Parser struct {
	tokens []lexer.Token
	idx    int
	Errors []common.Error
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

	got := fmt.Sprintf("%v \"%v\"", t.Kind, t.Subkind)
	want := fmt.Sprintf("%v \"%v\"", kind, subkind)
	p.addError(fmt.Sprintf("expected %s but got %s", want, got), t.Pos)
	return nil
}

func (p *Parser) addError(msg string, pos *common.SourcePos) {
	p.Errors = append(p.Errors, common.Error{
		Message: msg,
		Pos:     pos,
	})
}

func (p *Parser) ParseProgram() *Program {
	statements := []Statement{}
	t := p.peek(0)
	if t == nil {
		return &Program{Statements: statements, PosAt: nil}
	}
	for {
		statement := p.ParseStatement()
		if statement == nil {
			break
		}
		statements = append(statements, statement)
	}
	return &Program{Statements: statements, PosAt: t.Pos}
}

func (p *Parser) ParseStatement() Statement {
	t := p.peek(0)
	if t == nil {
		return nil
	}

	next := p.peek(1)

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
		return p.ParseDeclaration()
	}

	if t.Kind == lexer.Identifier && t.Subkind == lexer.IdentifierName && next != nil && next.Kind == lexer.Punctuator && next.Subkind == lexer.Assign {
		return p.ParseAssignment()
	}

	if t.Kind == lexer.Keyword && t.Subkind == lexer.KeywordIf {
		return p.ParseIf()
	}

	expression := p.ParseExpression(true)
	if expression == nil {
		p.addError(fmt.Sprintf("expected statement but got %v(%v)", t.Kind, t.Subkind), t.Pos)
		return nil
	}
	semicolon := p.eatExpected(lexer.Punctuator, lexer.StatementEnd, "expected ';'")
	if semicolon == nil {
		return nil
	}

	return expression
}

func (p *Parser) ParseType() *Type {
	tok := p.peek(0)
	if tok == nil {
		return nil
	}
	if tok.Kind == lexer.Punctuator && tok.Subkind == lexer.BracketOpen {
		p.eat()
		arrClose := p.eatExpected(lexer.Punctuator, lexer.BracketClose, "expected ']'")
		if arrClose == nil {
			return nil
		}
		elementType := p.eatExpected(lexer.Type, nil, "expected type")
		if elementType == nil {
			return nil
		}
		return TypeArrayOf(TypeFromTypeSubkind(elementType.Subkind.(lexer.TypeSubkind)))
	} else if tok.Kind == lexer.Type {
		p.eat()
		return TypeFromTypeSubkind(tok.Subkind.(lexer.TypeSubkind))
	}
	p.addError(fmt.Sprintf("expected type but got %v(%v)", tok.Kind, tok.Subkind), tok.Pos)
	return nil
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

	params := []Param{}
	for {
		param := p.ParseParam()
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
	colon := p.eatExpected(lexer.Punctuator, lexer.Colon, "expected ':'")
	if colon == nil {
		return nil
	}
	returnType := p.ParseType()
	if returnType == nil {
		return nil
	}
	body := p.ParseBody()

	return &Function{Name: idTok.Lexeme, Params: params, Body: body, ReturnType: returnType, PosAt: kw.Pos}
}

func (p *Parser) ParseParam() *Param {
	idTok := p.eatExpected(lexer.Identifier, lexer.IdentifierName, "expected identifier")
	if idTok == nil {
		return nil
	}
	colon := p.eatExpected(lexer.Punctuator, lexer.Colon, "expected ':'")
	if colon == nil {
		return nil
	}
	typeOf := p.ParseType()
	if typeOf == nil {
		return nil
	}
	return &Param{Name: idTok.Lexeme, TypeOf: typeOf, PosAt: idTok.Pos}
}

func (p *Parser) ParseBody() []Statement {
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
	return statements
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

	value := p.ParseExpression(false)

	semicolon := p.eatExpected(lexer.Punctuator, lexer.StatementEnd, "expected ';'")
	if semicolon == nil {
		return nil
	}

	return &Return{
		Value: value,
		PosAt: kw.Pos,
	}
}

func (p *Parser) ParseDeclaration() *Declaration {
	kw := p.eatExpected(lexer.Keyword, nil, "expected 'var' or 'const'")
	if kw == nil {
		return nil
	}

	specifier := kw.Subkind.(lexer.KeywordSubkind)
	if specifier != lexer.KeywordVar && specifier != lexer.KeywordConst {
		p.addError(fmt.Sprintf("expected 'var' or 'const' but got %v(%v)", kw.Kind, kw.Subkind), kw.Pos)
		return nil
	}
	isMutable := specifier == lexer.KeywordVar

	idTok := p.eatExpected(lexer.Identifier, lexer.IdentifierName, "expected identifier")
	if idTok == nil {
		return nil
	}

	isTyped := false
	var typeOf *Type
	colonTok := p.peek(0)
	if colonTok != nil && colonTok.Kind == lexer.Punctuator && colonTok.Subkind == lexer.Colon {
		isTyped = true
		p.eat()
		typeOf = p.ParseType()
		if typeOf == nil {
			return nil
		}
	}

	eqTok := p.peek(0)
	var value Expression = nil
	if eqTok != nil && eqTok.Kind == lexer.Punctuator && eqTok.Subkind == lexer.Assign {
		p.eat()
		value = p.ParseExpression(false)
		if value == nil {
			return nil
		}
	}

	semicolon := p.eatExpected(lexer.Punctuator, lexer.StatementEnd, "expected ';'")
	if semicolon == nil {
		return nil
	}

	return &Declaration{
		IsMutable:  isMutable,
		IsTyped:    isTyped,
		TypeOf:     typeOf,
		Identifier: &Identifier{Name: idTok.Lexeme, PosAt: idTok.Pos},
		Value:      value,
		PosAt:      kw.Pos,
	}
}

func (p *Parser) ParseAssignment() *Assignment {
	idTok := p.eatExpected(lexer.Identifier, lexer.IdentifierName, "expected identifier")
	if idTok == nil {
		return nil
	}

	eq := p.eatExpected(lexer.Punctuator, lexer.Assign, "expected '='")
	if eq == nil {
		return nil
	}

	value := p.ParseExpression(false)
	if value == nil {
		return nil
	}

	semicolon := p.eatExpected(lexer.Punctuator, lexer.StatementEnd, "expected ';'")
	if semicolon == nil {
		return nil
	}

	return &Assignment{
		Identifier: &Identifier{Name: idTok.Lexeme, PosAt: idTok.Pos},
		Value:      value,
		PosAt:      idTok.Pos,
	}
}

func (p *Parser) ParseIf() *If {
	kw := p.eatExpected(lexer.Keyword, lexer.KeywordIf, "expected 'if'")
	if kw == nil {
		return nil
	}
	lparen := p.eatExpected(lexer.Punctuator, lexer.ParenOpen, "expected '('")
	if lparen == nil {
		return nil
	}
	condition := p.ParseExpression(false)
	if condition == nil {
		return nil
	}
	rparen := p.eatExpected(lexer.Punctuator, lexer.ParenClose, "expected ')'")
	if rparen == nil {
		return nil
	}
	body := p.ParseBody()
	if body == nil {
		return nil
	}
	elseBody := []Statement{}
	elseKw := p.peek(0)
	if elseKw != nil && elseKw.Kind == lexer.Keyword && elseKw.Subkind == lexer.KeywordElse {
		p.eat()
		elseBody = p.ParseBody()
	}
	return &If{Condition: condition, Body: body, ElseBody: elseBody, PosAt: kw.Pos}
}

func (p *Parser) ParseExpression(isStatement bool) Expression {
	return p.ParseLogicalOrExpr(isStatement)
}

func (p *Parser) ParseLogicalOrExpr(isStatement bool) Expression {
	left := p.ParseLogicalAndExpr(isStatement)
	if left == nil {
		return nil
	}

	for {
		op := p.peek(0)
		if op == nil || op.Kind != lexer.Operator || op.Subkind != lexer.OperatorOr {
			break
		}

		p.eat()
		right := p.ParseLogicalAndExpr(isStatement)
		if right == nil {
			return nil
		}

		left = &BinaryExpr{Left: left, Operator: lexer.OperatorOr, Right: right, IsStatement: isStatement}
	}

	return left
}

func (p *Parser) ParseLogicalAndExpr(isStatement bool) Expression {
	left := p.ParseEqualityExpr(isStatement)
	if left == nil {
		return nil
	}

	for {
		op := p.peek(0)
		if op == nil || op.Kind != lexer.Operator || op.Subkind != lexer.OperatorAnd {
			break
		}

		p.eat()
		right := p.ParseEqualityExpr(isStatement)
		if right == nil {
			return nil
		}

		left = &BinaryExpr{
			Left:        left,
			Operator:    lexer.OperatorAnd,
			Right:       right,
			PosAt:       op.Pos,
			IsStatement: isStatement,
		}
	}

	return left
}

func (p *Parser) ParseEqualityExpr(isStatement bool) Expression {
	left := p.ParseComparisonExpr(isStatement)
	if left == nil {
		return nil
	}

	for {
		op := p.peek(0)
		if op == nil || op.Kind != lexer.Operator ||
			(op.Subkind != lexer.OperatorEqual &&
				op.Subkind != lexer.OperatorNotEqual) {
			break
		}

		p.eat()
		right := p.ParseComparisonExpr(isStatement)
		if right == nil {
			return nil
		}

		left = &BinaryExpr{
			Left:        left,
			Operator:    op.Subkind.(lexer.OperatorSubkind),
			Right:       right,
			PosAt:       op.Pos,
			IsStatement: isStatement,
		}
	}

	return left
}

func (p *Parser) ParseComparisonExpr(isStatement bool) Expression {
	left := p.ParseAdditiveExpr(isStatement)
	if left == nil {
		return nil
	}

	for {
		op := p.peek(0)
		if op == nil || op.Kind != lexer.Operator ||
			(op.Subkind != lexer.OperatorLess &&
				op.Subkind != lexer.OperatorLessEqual &&
				op.Subkind != lexer.OperatorGreater &&
				op.Subkind != lexer.OperatorGreaterEqual) {
			break
		}

		p.eat()
		right := p.ParseAdditiveExpr(isStatement)
		if right == nil {
			return nil
		}

		left = &BinaryExpr{
			Left:        left,
			Operator:    op.Subkind.(lexer.OperatorSubkind),
			Right:       right,
			PosAt:       op.Pos,
			IsStatement: isStatement,
		}
	}

	return left
}

func (p *Parser) ParseAdditiveExpr(isStatement bool) Expression {
	left := p.ParseMultiplicativeExpr(isStatement)
	if left == nil {
		return nil
	}

	for {
		op := p.peek(0)
		if op == nil || op.Kind != lexer.Operator ||
			(op.Subkind != lexer.OperatorPlus && op.Subkind != lexer.OperatorMinus) {
			break
		}

		p.eat()
		right := p.ParseMultiplicativeExpr(isStatement)
		if right == nil {
			return nil
		}

		left = &BinaryExpr{
			Left:        left,
			Operator:    op.Subkind.(lexer.OperatorSubkind),
			Right:       right,
			PosAt:       op.Pos,
			IsStatement: isStatement,
		}
	}

	return left
}

func (p *Parser) ParseMultiplicativeExpr(isStatement bool) Expression {
	left := p.ParsePrimaryExpr(isStatement)
	if left == nil {
		return nil
	}

	for {
		op := p.peek(0)
		if op == nil || op.Kind != lexer.Operator ||
			(op.Subkind != lexer.OperatorStar && op.Subkind != lexer.OperatorSlash) {
			break
		}

		p.eat()
		right := p.ParsePrimaryExpr(isStatement)
		if right == nil {
			return nil
		}

		left = &BinaryExpr{
			Left:        left,
			Operator:    op.Subkind.(lexer.OperatorSubkind),
			Right:       right,
			PosAt:       op.Pos,
			IsStatement: isStatement,
		}
	}

	return left
}

func (p *Parser) ParsePrimaryExpr(isStatement bool) Expression {
	tok := p.peek(0)

	if tok.Kind == lexer.Punctuator && tok.Subkind == lexer.ParenOpen {
		p.eat()

		expr := p.ParseExpression(isStatement)
		if expr == nil {
			return nil
		}

		rparen := p.eatExpected(lexer.Punctuator, lexer.ParenClose, "expected ')'")
		if rparen == nil {
			return nil
		}

		return expr
	}

	return p.ParseAtomExpr(isStatement)
}

func (p *Parser) ParseAtomExpr(isStatement bool) Expression {
	t := p.peek(0)
	if t == nil {
		return nil
	}

	if t.Kind == lexer.Constant && t.Subkind == lexer.Integer {
		return p.ParseIntLiteral(isStatement)
	}

	if t.Kind == lexer.Constant && t.Subkind == lexer.Boolean {
		return p.ParseBoolLiteral(isStatement)
	}

	if t.Kind == lexer.Constant && t.Subkind == lexer.Null {
		return p.ParseNullLiteral(isStatement)
	}

	if t.Kind == lexer.Punctuator && t.Subkind == lexer.BracketOpen {
		return p.ParseArrayLiteral(isStatement)
	}

	if t.Kind == lexer.Identifier {
		t := p.peek(1)
		if t != nil && t.Kind == lexer.Punctuator && t.Subkind == lexer.ParenOpen {
			return p.ParseCallExpr()
		}

		if t != nil && t.Kind == lexer.Punctuator && t.Subkind == lexer.BracketOpen {
			return p.ParseIndexExpr(isStatement)
		}

		return p.ParseIdentifier(isStatement)
	}

	return nil
}

// ---------- Atoms ----------

func (p *Parser) ParseIntLiteral(isStatement bool) *IntLiteral {
	t := p.eat()
	if t == nil {
		return nil
	}

	return &IntLiteral{
		Value:       atoi(t.Lexeme),
		PosAt:       t.Pos,
		IsStatement: isStatement,
	}
}

func (p *Parser) ParseBoolLiteral(isStatement bool) *BoolLiteral {
	t := p.eat()
	if t == nil {
		return nil
	}

	if t.Lexeme != "true" && t.Lexeme != "false" {
		p.addError(fmt.Sprintf("expected 'true' or 'false' but got %s", t.Lexeme), t.Pos)
		return nil
	}

	return &BoolLiteral{Value: t.Lexeme == "true", PosAt: t.Pos, IsStatement: isStatement}
}

func (p *Parser) ParseNullLiteral(isStatement bool) *NullLiteral {
	t := p.eat()
	if t == nil {
		return nil
	}

	if t.Lexeme != "null" {
		p.addError(fmt.Sprintf("expected 'null' but got %s", t.Lexeme), t.Pos)
		return nil
	}

	return &NullLiteral{PosAt: t.Pos, IsStatement: isStatement}
}

func (p *Parser) ParseArrayLiteral(isStatement bool) *ArrayLiteral {
	lparen := p.eatExpected(lexer.Punctuator, lexer.BracketOpen, "expected '['")
	if lparen == nil {
		return nil
	}
	elements := []Expression{}
	for {
		t := p.peek(0)
		if t == nil {
			return nil
		}
		if t.Kind == lexer.Punctuator && t.Subkind == lexer.BracketClose {
			break
		}
		element := p.ParseExpression(false)
		if element == nil {
			return nil
		}
		elements = append(elements, element)
		commaTok := p.peek(0)
		if commaTok == nil || !(commaTok.Kind == lexer.Punctuator && commaTok.Subkind == lexer.Comma) {
			break
		}
		p.eat()
	}
	rparen := p.eatExpected(lexer.Punctuator, lexer.BracketClose, "expected ']'")
	if rparen == nil {
		return nil
	}
	return &ArrayLiteral{Elements: elements, PosAt: lparen.Pos, IsStatement: isStatement}
}

func (p *Parser) ParseIdentifier(isStatement bool) *Identifier {
	idTok := p.eatExpected(lexer.Identifier, lexer.IdentifierName, "expected identifier")
	if idTok == nil {
		return nil
	}

	return &Identifier{Name: idTok.Lexeme, PosAt: idTok.Pos, IsStatement: isStatement}
}

func (p *Parser) ParseCallExpr() *CallExpr {
	identifier := p.ParseIdentifier(false)
	if identifier == nil {
		return nil
	}
	lparen := p.eatExpected(lexer.Punctuator, lexer.ParenOpen, "expected '('")
	if lparen == nil {
		return nil
	}
	arguments := []Expression{}
	t := p.peek(0)
	if t == nil {
		return nil
	}
	if t.Kind == lexer.Punctuator && t.Subkind == lexer.ParenClose {
		return &CallExpr{Identifier: *identifier, Arguments: arguments, PosAt: lparen.Pos}
	}
	for {
		argument := p.ParseExpression(false)
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
	return &CallExpr{Identifier: *identifier, Arguments: arguments, PosAt: lparen.Pos}
}

func (p *Parser) ParseIndexExpr(isStatement bool) *IndexExpr {
	array := p.ParseIdentifier(false)
	if array == nil {
		return nil
	}
	bracketOpen := p.eatExpected(lexer.Punctuator, lexer.BracketOpen, "expected '['")
	if bracketOpen == nil {
		return nil
	}
	index := p.ParseExpression(false)
	if index == nil {
		return nil
	}
	bracketClose := p.eatExpected(lexer.Punctuator, lexer.BracketClose, "expected ']'")
	if bracketClose == nil {
		return nil
	}
	return &IndexExpr{Array: array, Index: index, PosAt: array.PosAt, IsStatement: isStatement}
}

func atoi(s string) int {
	var n int
	for i := 0; i < len(s); i++ {
		n = n*10 + int(s[i]-'0')
	}
	return n
}
