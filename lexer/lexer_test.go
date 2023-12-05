package lexer

import (
	"Lisa/token"
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	allInput := `=+(){},;`

	tokenTests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(allInput)
	for i, tc := range tokenTests {
		tok := l.NextToken()
		if tok.Type != tc.expectedType {
			// %q is a single-quoted character literal safely escaped with Go syntax.
			t.Errorf("tests[%d] - error token type: expected %q, got %q.\n", i, tc.expectedType, tok.Type)
		}

		if tok.Literal != tc.expectedLiteral {
			t.Errorf("tests[%d] - error token literal value: expected %q, got %q.\n", i, tc.expectedLiteral, tok.Literal)
		}
	}
}
