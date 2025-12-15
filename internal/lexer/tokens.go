package lexer

import (
	"fmt"

	"youpiteron.dev/white-monster-on-friday-night/internal/common"
)

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
	OperatorEqual
	OperatorNotEqual
	OperatorGreater
	OperatorGreaterEqual
	OperatorLess
	OperatorLessEqual
	OperatorAnd
	OperatorOr
)

func (k OperatorSubkind) String() string {
	return [...]string{
		"+",
		"-",
		"*",
		"/",
		"==",
		"!=",
		">",
		">=",
		"<",
		"<=",
		"&&",
		"||",
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

// ---- Token ----

type Token struct {
	Lexeme  string
	Kind    TokenKind
	Subkind any
	Pos     *common.SourcePos
}

func (t Token) String() string {
	return fmt.Sprintf("Token{Lexeme: %s, Kind: %s, Subkind: %v, Pos: %v}", t.Lexeme, t.Kind, t.Subkind, t.Pos)
}
