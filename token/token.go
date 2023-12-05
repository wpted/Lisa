package token

const (
	// ILLEGAL signifies a token/character we don't know illegal.
	ILLEGAL = "ILLEGAL"
	// EOF means 'End Of File', which tells our parser later on that it can stop.
	EOF = "EOF"

	IDENT = "IDENT"
	INT   = "INT"

	ASSIGN = "="
	PLUS   = "+"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	VAR      = "VAR"
)

type TokenType string

// Token is the transformation result of lexing source code.
type Token struct {
	Type TokenType
	// Literal is the parsed value of a token.
	Literal string
}

// New creates a new *token.
func New(tt TokenType, literal string) *Token {
	return &Token{
		Type:    tt,
		Literal: literal,
	}
}
