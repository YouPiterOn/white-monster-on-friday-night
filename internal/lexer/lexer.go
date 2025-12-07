package lexer

import "fmt"

// ---------- Errors ----------

type LexError struct {
	Message string
	Pos     SourcePos
}

// ---------- Lexer ----------

type LexerState int

const (
	StateInitial LexerState = iota
	StateIdentifier
	StateNumber
)

type LexResult struct {
	Tokens []Token
	Errors []LexError
}

type Lexer struct {
	input string
	idx   int
	line  int
	col   int

	state    LexerState
	buf      string
	startPos *BasePos
}

func NewLexer() *Lexer {
	return &Lexer{}
}

func (l *Lexer) Lex(input string) LexResult {
	l.reset(input)

	var tokens []Token
	var errors []LexError

	for !l.eof() {
		ch := l.peek()

		switch l.state {

		case StateInitial:
			// whitespace
			if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
				l.next()
				continue
			}

			// punctuator '='
			if ch == '=' {
				pos := l.posSpan(1)
				l.next()
				tokens = append(tokens, Token{
					Lexeme:  "=",
					Kind:    Punctuator,
					Subkind: Assign,
					Pos:     &pos,
				})
				continue
			}

			// operator
			if op, ok := operatorSubkind(ch); ok {
				pos := l.posSpan(1)
				l.next()
				tokens = append(tokens, Token{
					Lexeme:  string(ch),
					Kind:    Operator,
					Subkind: op,
					Pos:     &pos,
				})
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
				continue
			}

			// unexpected character
			pos := l.posSpan(1)
			errors = append(errors, LexError{
				Message: fmt.Sprintf("unexpected symbol '%c'", ch),
				Pos:     pos,
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

func (l *Lexer) flushIdentifier() (*Token, *LexError) {
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

	// identifier
	return &Token{
		Lexeme:  lex,
		Kind:    Identifier,
		Subkind: IdentifierName,
		Pos:     &pos,
	}, nil
}

func (l *Lexer) flushNumber() (*Token, *LexError) {
	if l.startPos == nil {
		return nil, nil
	}

	lex := l.buf
	pos := l.finishPos(*l.startPos, len(lex))

	return &Token{
		Lexeme:  lex,
		Kind:    Constant,
		Subkind: Numeric,
		Pos:     &pos,
	}, nil
}

// ---------- Helpers ----------

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

func operatorSubkind(ch byte) (OperatorSubkind, bool) {
	switch ch {
	case '+':
		return Plus, true
	case '-':
		return Minus, true
	case '*':
		return Star, true
	default:
		return 0, false
	}
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

func (l *Lexer) posSpan(length int) SourcePos {
	return SourcePos{
		BasePos: BasePos{
			Offset: l.idx,
			Line:   l.line,
			Column: l.col,
		},
		Length: length,
	}
}

func (l *Lexer) capturePos() *BasePos {
	return &BasePos{
		Offset: l.idx,
		Line:   l.line,
		Column: l.col,
	}
}

func (l *Lexer) finishPos(start BasePos, length int) SourcePos {
	return SourcePos{
		BasePos: BasePos{
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
