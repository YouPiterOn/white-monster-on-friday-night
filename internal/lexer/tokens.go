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
		"endline",
		"operator",
	}[k]
}

type KeywordSubkind int

const (
	KeywordConst KeywordSubkind = iota
	KeywordVar
	KeywordReturn
)

type IdentifierSubkind int

const (
	IdentifierName IdentifierSubkind = iota
)

type PunctuatorSubkind int

const (
	Assign PunctuatorSubkind = iota
)

type ConstantSubkind int

const (
	Numeric ConstantSubkind = iota
)

type EndlineSubkind int

const (
	EndlineNone EndlineSubkind = iota
)

type OperatorSubkind int

const (
	Plus OperatorSubkind = iota
	Minus
	Star
)

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
