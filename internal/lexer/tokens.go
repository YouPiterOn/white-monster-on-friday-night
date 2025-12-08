package lexer

import "fmt"

type TokenKind int

const (
	Keyword TokenKind = iota
	Identifier
	Punctuator
	Constant
	Operator
)

func (k TokenKind) String() string {
	return [...]string{
		"keyword",
		"identifier",
		"punctuator",
		"constant",
		"operator",
	}[k]
}

type KeywordSubkind int

const (
	KeywordConst KeywordSubkind = iota
	KeywordVar
	KeywordReturn
	KeywordFunction
)

func (k KeywordSubkind) String() string {
	return [...]string{
		"const",
		"var",
		"return",
		"function",
	}[k]
}

type IdentifierSubkind int

const (
	IdentifierName IdentifierSubkind = iota
)

func (k IdentifierSubkind) String() string {
	return [...]string{
		"name",
	}[k]
}

type PunctuatorSubkind int

const (
	Assign PunctuatorSubkind = iota
	BlockStart
	BlockEnd
	ParenOpen
	ParenClose
	StatementEnd
	Comma
)

func (k PunctuatorSubkind) String() string {
	return [...]string{
		"=",
		"{",
		"}",
		"(",
		")",
		";",
		",",
	}[k]
}

type ConstantSubkind int

const (
	Numeric ConstantSubkind = iota
)

func (k ConstantSubkind) String() string {
	return [...]string{
		"numeric",
	}[k]
}

type OperatorSubkind int

const (
	Plus OperatorSubkind = iota
	Minus
	Star
)

func (k OperatorSubkind) String() string {
	return [...]string{
		"+",
		"-",
		"*",
	}[k]
}

// ---- Source positions ----

type BasePos struct {
	Offset int
	Line   int
	Column int
}

type SourcePos struct {
	BasePos
	Length int
}

func (p SourcePos) String() string {
	return fmt.Sprintf("SourcePos{Offset: %d, Line: %d, Column: %d, Length: %d}", p.Offset, p.Line, p.Column, p.Length)
}

// ---- Token ----

type Token struct {
	Lexeme  string
	Kind    TokenKind
	Subkind any
	Pos     *SourcePos
}

func (t Token) String() string {
	return fmt.Sprintf("Token{Lexeme: %s, Kind: %s, Subkind: %v, Pos: %v}", t.Lexeme, t.Kind, t.Subkind, t.Pos)
}
