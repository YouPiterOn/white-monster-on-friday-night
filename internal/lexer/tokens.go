package lexer

import "fmt"

type TokenKind int

const (
	Keyword TokenKind = iota
	Identifier
	Punctuator
	Constant
	Operator
	Type
)

func (k TokenKind) String() string {
	return [...]string{
		"keyword",
		"identifier",
		"punctuator",
		"constant",
		"operator",
		"type",
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
	Colon
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
		":",
	}[k]
}

type ConstantSubkind int

const (
	Integer ConstantSubkind = iota
	Boolean
	Null
)

func (k ConstantSubkind) String() string {
	return [...]string{
		"numeric",
		"boolean",
		"null",
	}[k]
}

type OperatorSubkind int

const (
	OperatorPlus OperatorSubkind = iota
	OperatorMinus
	OperatorStar
	OperatorSlash
)

func (k OperatorSubkind) String() string {
	return [...]string{
		"+",
		"-",
		"*",
		"/",
	}[k]
}

type TypeSubkind int

const (
	TypeInt TypeSubkind = iota
	TypeBool
	TypeNull
)

func (k TypeSubkind) String() string {
	return [...]string{
		"int",
		"bool",
		"null",
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
	return fmt.Sprintf("Offset: %d, Line: %d, Column: %d, Length: %d", p.Offset, p.Line, p.Column, p.Length)
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
