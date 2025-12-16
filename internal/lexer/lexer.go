package lexer

import (
	"fmt"

	"youpiteron.dev/white-monster-on-friday-night/internal/common"
)

// ---------- Lexer ----------

type LexerState int

const (
	StateInitial LexerState = iota
	StateIdentifier
	StateNumber
	StateOperator
)

type LexResult struct {
	Tokens []Token
	Errors []common.Error
}

type Lexer struct {
	input string
	idx   int
	line  int
	col   int

	state    LexerState
	buf      string
	startPos *common.BasePos
}

func NewLexer() *Lexer {
	return &Lexer{}
}

func (l *Lexer) Lex(input string) LexResult {
	l.reset(input)

	var tokens []Token
	var errors []common.Error

	for !l.eof() {
		ch := l.peek()

		switch l.state {

		case StateInitial:
			// whitespace
			if isWs(ch) {
				l.next()
				continue
			}

			// punctuator ';'
			if ch == ';' {
				pos := l.posSpan(1)
				l.next()
				tokens = append(tokens, Token{
					Lexeme:  ";",
					Kind:    Punctuator,
					Subkind: StatementEnd,
					Pos:     &pos,
				})
				continue
			}

			// punctuator ','
			if ch == ',' {
				pos := l.posSpan(1)
				l.next()
				tokens = append(tokens, Token{
					Lexeme:  ",",
					Kind:    Punctuator,
					Subkind: Comma,
					Pos:     &pos,
				})
				continue
			}

			// punctuator '{'
			if ch == '{' {
				pos := l.posSpan(1)
				l.next()
				tokens = append(tokens, Token{
					Lexeme:  "{",
					Kind:    Punctuator,
					Subkind: BlockStart,
					Pos:     &pos,
				})
				continue
			}

			// punctuator '}'
			if ch == '}' {
				pos := l.posSpan(1)
				l.next()
				tokens = append(tokens, Token{
					Lexeme:  "}",
					Kind:    Punctuator,
					Subkind: BlockEnd,
					Pos:     &pos,
				})
				continue
			}

			// punctuator '('
			if ch == '(' {
				pos := l.posSpan(1)
				l.next()
				tokens = append(tokens, Token{
					Lexeme:  "(",
					Kind:    Punctuator,
					Subkind: ParenOpen,
					Pos:     &pos,
				})
				continue
			}

			// punctuator ')'
			if ch == ')' {
				pos := l.posSpan(1)
				l.next()
				tokens = append(tokens, Token{
					Lexeme:  ")",
					Kind:    Punctuator,
					Subkind: ParenClose,
					Pos:     &pos,
				})
				continue
			}

			// punctuator ':'
			if ch == ':' {
				pos := l.posSpan(1)
				l.next()
				tokens = append(tokens, Token{
					Lexeme:  ":",
					Kind:    Punctuator,
					Subkind: Colon,
					Pos:     &pos,
				})
				continue
			}

			// operator
			if isOperatorStart(ch) {
				l.state = StateOperator
				l.startPos = l.capturePos()
				l.buf = string(ch)
				l.next()
				continue
			}

			// number
			if isDigit(ch) {
				l.state = StateNumber
				l.startPos = l.capturePos()
				continue
			}

			// identifier
			if isIdentStart(ch) {
				l.state = StateIdentifier
				l.startPos = l.capturePos()
				l.buf = string(ch)
				l.next()
				continue
			}

			// unexpected character
			pos := l.posSpan(1)
			errors = append(errors, common.Error{
				Message: fmt.Sprintf("unexpected symbol '%c'", ch),
				Pos:     &pos,
			})
			l.next()
			continue

		case StateIdentifier:
			if !l.eof() && isIdentContinue(l.peek()) {
				l.buf += string(l.next())
				continue
			}

			// finalize identifier / keyword
			if tok, err := l.flushIdentifier(); tok != nil {
				tokens = append(tokens, *tok)
			} else if err != nil {
				errors = append(errors, *err)
			}

			l.state = StateInitial
			l.buf = ""
			l.startPos = nil
			continue

		case StateNumber:
			if !l.eof() && isDigit(l.peek()) {
				l.buf += string(l.next())
				continue
			}

			if tok, err := l.flushNumber(); tok != nil {
				tokens = append(tokens, *tok)
			} else if err != nil {
				errors = append(errors, *err)
			}

			l.state = StateInitial
			l.buf = ""
			l.startPos = nil
			continue

		case StateOperator:
			if !l.eof() && isOperatorContinue(l.peek()) {
				l.buf += string(l.next())
				continue
			}
			if tok, err := l.flushOperator(); tok != nil {
				tokens = append(tokens, *tok)
			} else if err != nil {
				errors = append(errors, *err)
			}
			l.state = StateInitial
			l.buf = ""
			l.startPos = nil
			continue
		}
	}

	// final buffer flush
	if l.buf != "" && l.startPos != nil {
		switch l.state {
		case StateIdentifier:
			if tok, err := l.flushIdentifier(); tok != nil {
				tokens = append(tokens, *tok)
			} else if err != nil {
				errors = append(errors, *err)
			}
		case StateNumber:
			if tok, err := l.flushNumber(); tok != nil {
				tokens = append(tokens, *tok)
			} else if err != nil {
				errors = append(errors, *err)
			}
		}
	}

	return LexResult{Tokens: tokens, Errors: errors}
}

// ---------- Flushers ----------

func (l *Lexer) flushIdentifier() (*Token, *common.Error) {
	if l.startPos == nil {
		return nil, nil
	}

	lex := l.buf
	pos := l.finishPos(*l.startPos, len(lex))

	// keywords
	if lex == "return" {
		return &Token{
			Lexeme:  lex,
			Kind:    Keyword,
			Subkind: KeywordReturn,
			Pos:     &pos,
		}, nil
	}
	if lex == "const" {
		return &Token{
			Lexeme:  lex,
			Kind:    Keyword,
			Subkind: KeywordConst,
			Pos:     &pos,
		}, nil
	}
	if lex == "var" {
		return &Token{
			Lexeme:  lex,
			Kind:    Keyword,
			Subkind: KeywordVar,
			Pos:     &pos,
		}, nil
	}
	if lex == "function" {
		return &Token{
			Lexeme:  lex,
			Kind:    Keyword,
			Subkind: KeywordFunction,
			Pos:     &pos,
		}, nil
	}
	if lex == "if" {
		return &Token{
			Lexeme:  lex,
			Kind:    Keyword,
			Subkind: KeywordIf,
			Pos:     &pos,
		}, nil
	}
	if lex == "else" {
		return &Token{
			Lexeme:  lex,
			Kind:    Keyword,
			Subkind: KeywordElse,
			Pos:     &pos,
		}, nil
	}

	// constants
	if lex == "true" {
		return &Token{
			Lexeme:  lex,
			Kind:    Constant,
			Subkind: Boolean,
			Pos:     &pos,
		}, nil
	}
	if lex == "false" {
		return &Token{
			Lexeme:  lex,
			Kind:    Constant,
			Subkind: Boolean,
			Pos:     &pos,
		}, nil
	}
	if lex == "null" {
		return &Token{
			Lexeme:  lex,
			Kind:    Constant,
			Subkind: Null,
			Pos:     &pos,
		}, nil
	}

	if subkind, ok := typeSubkind(lex); ok {
		return &Token{
			Lexeme:  lex,
			Kind:    Type,
			Subkind: subkind,
			Pos:     &pos,
		}, nil
	}

	// identifier
	return &Token{
		Lexeme:  lex,
		Kind:    Identifier,
		Subkind: IdentifierName,
		Pos:     &pos,
	}, nil
}

func (l *Lexer) flushOperator() (*Token, *common.Error) {
	if l.startPos == nil {
		return nil, nil
	}

	lex := l.buf
	pos := l.finishPos(*l.startPos, len(lex))

	// punctuator '='
	if lex == "=" {
		return &Token{
			Lexeme:  lex,
			Kind:    Punctuator,
			Subkind: Assign,
			Pos:     &pos,
		}, nil
	}

	if op, ok := operatorSubkind(lex); ok {
		return &Token{
			Lexeme:  lex,
			Kind:    Operator,
			Subkind: op,
			Pos:     &pos,
		}, nil
	}
	return nil, &common.Error{Message: fmt.Sprintf("invalid operator: %s", lex), Pos: &pos}
}

func (l *Lexer) flushNumber() (*Token, *common.Error) {
	if l.startPos == nil {
		return nil, nil
	}

	lex := l.buf
	pos := l.finishPos(*l.startPos, len(lex))

	return &Token{
		Lexeme:  lex,
		Kind:    Constant,
		Subkind: Integer,
		Pos:     &pos,
	}, nil
}

// ---------- Helpers ----------

func isWs(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		ch == '_'
}

func isIdentStart(ch byte) bool {
	return isLetter(ch)
}

func isIdentContinue(ch byte) bool {
	return isLetter(ch) || isDigit(ch)
}

func isOperatorStart(ch byte) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '=' || ch == '!' || ch == '>' || ch == '<' || ch == '&' || ch == '|'
}

func isOperatorContinue(ch byte) bool {
	return ch == '=' || ch == '&' || ch == '|'
}

func operatorSubkind(lex string) (OperatorSubkind, bool) {
	switch lex {
	case "+":
		return OperatorPlus, true
	case "-":
		return OperatorMinus, true
	case "*":
		return OperatorStar, true
	case "/":
		return OperatorSlash, true
	case "==":
		return OperatorEqual, true
	case "!=":
		return OperatorNotEqual, true
	case ">":
		return OperatorGreater, true
	case ">=":
		return OperatorGreaterEqual, true
	case "<":
		return OperatorLess, true
	case "<=":
		return OperatorLessEqual, true
	case "&&":
		return OperatorAnd, true
	case "||":
		return OperatorOr, true
	default:
		return 0, false
	}
}

func typeSubkind(lex string) (TypeSubkind, bool) {
	switch lex {
	case "int":
		return TypeInt, true
	case "bool":
		return TypeBool, true
	case "null":
		return TypeNull, true
	}
	return 0, false
}

// ---------- Cursor management ----------

func (l *Lexer) eof() bool {
	return l.idx >= len(l.input)
}

func (l *Lexer) peek() byte {
	return l.input[l.idx]
}

func (l *Lexer) next() byte {
	ch := l.input[l.idx]
	l.idx++

	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}

	return ch
}

// ---------- Position helpers ----------

func (l *Lexer) posSpan(length int) common.SourcePos {
	return common.SourcePos{
		BasePos: common.BasePos{
			Offset: l.idx,
			Line:   l.line,
			Column: l.col,
		},
		Length: length,
	}
}

func (l *Lexer) capturePos() *common.BasePos {
	return &common.BasePos{
		Offset: l.idx,
		Line:   l.line,
		Column: l.col,
	}
}

func (l *Lexer) finishPos(start common.BasePos, length int) common.SourcePos {
	return common.SourcePos{
		BasePos: common.BasePos{
			Offset: start.Offset,
			Line:   start.Line,
			Column: start.Column,
		},
		Length: length,
	}
}

// ---------- Reset ----------

func (l *Lexer) reset(input string) {
	l.input = input
	l.idx = 0
	l.line = 1
	l.col = 1
	l.state = StateInitial
	l.buf = ""
	l.startPos = nil
}
