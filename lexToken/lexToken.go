package lexToken

const (
	// ILLEGAL signifies a token/character we don't know illegal.
	ILLEGAL = "ILLEGAL"
	// EOF means 'End Of File', which tells our parser later on that it can stop.
	EOF = "EOF"

	// IDENT represents the lexical type of variable or function names.
	IDENT = "IDENT"

	VAR      = "VAR"
	FUNCTION = "FUNCTION"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"

	// INT is the integer type.
	INT = "INT"

	ASSIGN      = "="
	PLUS        = "+"
	MINUS       = "-"
	EXCLAMATION = "!"
	CARET       = "^"
	ASTERISK    = "*"
	SLASH       = "/"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN      = "("
	RPAREN      = ")"
	LBRACE      = "{"
	RBRACE      = "}"
	LESSTHAN    = "<"
	GREATERTHAN = ">"
)

type LexicalType string

// reservedWordsTable is the table for reserved words.
var reservedWordsTable = map[string]LexicalType{
	"fn":     FUNCTION,
	"var":    VAR,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
}

// Token is the transformation result of lexing source code.
// Tokens can be classified into three: one-character token (e.g. -), two-character token (e.g. ==), and keyword token (e.g. return)
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

// LookUpReservedWord checks the reservedWordsTable to see whether the given identifier is a reserved word.
// If ident is a reserved word, return the Lexical Type of the token.
// If it isn't, return IDENT.
func LookUpReservedWord(ident string) (lt LexicalType, isReserved bool) {
	lt, isReserved = reservedWordsTable[ident]
	if !isReserved {
		return IDENT, isReserved
	}
	return lt, isReserved
}
