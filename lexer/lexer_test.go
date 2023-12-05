package lexer

import (
	"Lisa/token"
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	tokenTests := []struct {
		input        string
		parsedResult []struct {
			expectedType    token.LexicalType
			expectedLiteral string
		}
	}{
		{
			input: `=+(){},;`,
			parsedResult: []struct {
				expectedType    token.LexicalType
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
			},
		},
		{
			input: `var five = 5;
					var ten = 10;
					var add = fn(x, y) {
     					return x + y;
					};
   					var result = add(five, ten);
   					`,
			parsedResult: []struct {
				expectedType    token.LexicalType
				expectedLiteral string
			}{
				// var five = 5;
				{expectedType: token.VAR, expectedLiteral: "var"},
				{expectedType: token.IDENT, expectedLiteral: "five"},
				{expectedType: token.ASSIGN, expectedLiteral: "="},
				{expectedType: token.INT, expectedLiteral: "5"},
				{expectedType: token.SEMICOLON, expectedLiteral: ";"},

				// var ten = 10;
				{expectedType: token.VAR, expectedLiteral: "var"},
				{expectedType: token.IDENT, expectedLiteral: "ten"},
				{expectedType: token.ASSIGN, expectedLiteral: "="},
				{expectedType: token.INT, expectedLiteral: "10"},
				{expectedType: token.SEMICOLON, expectedLiteral: ";"},

				// var add = fn(x, y){return x + y;};
				{expectedType: token.VAR, expectedLiteral: "var"},
				{expectedType: token.IDENT, expectedLiteral: "add"},
				{expectedType: token.ASSIGN, expectedLiteral: "="},
				{expectedType: token.FUNCTION, expectedLiteral: "fn"},
				{expectedType: token.LPAREN, expectedLiteral: "("},
				{expectedType: token.IDENT, expectedLiteral: "x"},
				{expectedType: token.COMMA, expectedLiteral: ","},
				{expectedType: token.IDENT, expectedLiteral: "y"},
				{expectedType: token.LPAREN, expectedLiteral: ")"},
				{expectedType: token.LBRACE, expectedLiteral: "{"},
				{expectedType: token.RETURN, expectedLiteral: "return"},
				{expectedType: token.IDENT, expectedLiteral: "x"},
				{expectedType: token.PLUS, expectedLiteral: "+"},
				{expectedType: token.IDENT, expectedLiteral: "y"},
				{expectedType: token.SEMICOLON, expectedLiteral: ";"},
				{expectedType: token.RBRACE, expectedLiteral: "}"},
				{expectedType: token.SEMICOLON, expectedLiteral: ";"},

				{expectedType: token.VAR, expectedLiteral: "var"},
				{expectedType: token.IDENT, expectedLiteral: "result"},
				{expectedType: token.ASSIGN, expectedLiteral: "="},
				{expectedType: token.IDENT, expectedLiteral: "add"},
				{expectedType: token.LPAREN, expectedLiteral: "("},
				{expectedType: token.IDENT, expectedLiteral: "five"},
				{expectedType: token.COMMA, expectedLiteral: ","},
				{expectedType: token.IDENT, expectedLiteral: "ten"},
				{expectedType: token.LPAREN, expectedLiteral: ")"},
				{expectedType: token.SEMICOLON, expectedLiteral: ";"},
			},
		},
	}

	l := new(Lexer)
	for i, tt := range tokenTests {
		l = New(tt.input)
		for tokenIdx, tc := range tt.parsedResult {
			tok := l.NextToken()
			if tok.Type != tc.expectedType {
				// %q is a single-quoted character literal safely escaped with Go syntax.
				t.Errorf("tests[%d], token[%d] - error token type: expected %q, got %q.\n", i, tokenIdx, tc.expectedType, tok.Type)
			}

			if tok.Literal != tc.expectedLiteral {
				t.Errorf("tests[%d], token[%d] - error token literal value: expected %q, got %q.\n", i, tokenIdx, tc.expectedLiteral, tok.Literal)
			}
		}
	}
}
