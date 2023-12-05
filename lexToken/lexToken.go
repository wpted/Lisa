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

// reservedWordsTable is the table for reserved words.
var reservedWordsTable = map[string]LexicalType{
	"fn":     FUNCTION,
	"var":    VAR,
	"return": RETURN,
}

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

// LookUpReservedWords checks the reservedWordsTable to see whether the given identifier is a reserved word.
// If ident is a reserved word, return the Lexical Type of the token.
// If it isn't, return IDENT.
func LookUpReservedWords(ident string) LexicalType {
	lt, ok := reservedWordsTable[ident]
	if !ok {
		return IDENT
	}
	return lt
}
