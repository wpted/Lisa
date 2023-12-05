package lexToken

const (
	// ILLEGAL signifies a token/character we don't know illegal.
	ILLEGAL = "ILLEGAL"
	// EOF means 'End Of File', which tells our parser later on that it can stop.
	EOF = "EOF"

	// IDENT represents the lexical type of variable or function names.
	IDENT    = "IDENT"
	VAR      = "VAR"
	FUNCTION = "FUNCTION"
	RETURN   = "RETURN"

	// INT is the integer type.
	INT = "INT"

	ASSIGN = "="
	PLUS   = "+"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
)

type LexicalType string

// Token is the transformation result of lexing source code.
type Token struct {
	Type LexicalType
	// Literal is the parsed value of a token.
	Literal string
}

// New creates a new *token.
func New(tt LexicalType, literal string) *Token {
	return &Token{
		Type:    tt,
		Literal: literal,
	}
}