package lexer

import (
	token "Lisa/lexToken"
	"testing"
)

func Test_New(t *testing.T) {
	testInput := "test"
	l := New("test")
	if l.input != testInput {
		t.Errorf("error initializing lexer - wrong input: expected %s, have %s.\n", testInput, l.input)
	}
	if l.position != 0 {
		t.Errorf("error initializing lexer - wrong position: expected %d, have %d.\n", 0, l.position)
	}
	if l.readPosition != 1 {
		t.Errorf("error initializing lexer - wrong read position: expected %d, have %d.\n", 1, l.readPosition)
	}
	if l.ch != testInput[0] {
		t.Errorf("error initializing lexer - wrong character: expected %q, have %q.\n", testInput[0], l.ch)
	}
}

func TestLexer_readChar(t *testing.T) {
	testCases := []struct {
		input            string
		jumps            int
		expectedPosition int
		expectedCurrChar byte
	}{
		{"test 32", 3, 3, 't'},
		{"32", 3, 3, 0},
	}
	var lex *Lexer
	for i, tc := range testCases {
		lex = New(tc.input)
		for i := 0; i < tc.jumps; i++ {
			lex.readChar()
		}
		if tc.expectedPosition != lex.position {
			t.Errorf("tests[%d] - error reading char [wrong position]: expected %d, got %d.\n", i, tc.expectedPosition, lex.position)
		}
		if tc.expectedCurrChar != lex.ch {
			t.Errorf("tests[%d] - error reading char [wrong character]: expected %q, got %q.\n", i, tc.expectedCurrChar, lex.ch)
		}
	}
}

func TestLexer_readIdentifier(t *testing.T) {
	testCases := []struct {
		input         string
		expectedIdent string
	}{
		{"test 32", "test"},
		{"32", ""},
		{"test_abc 32", "test_abc"},
	}
	var lex *Lexer
	for i, tc := range testCases {
		lex = New(tc.input)
		re := lex.readIdentifier()
		if re != tc.expectedIdent {
			t.Errorf("tests[%d] - error reading number: expected %s, got %s.\n", i, tc.expectedIdent, re)
		}
	}
}

func Test_isLetter(t *testing.T) {
	testCases := []struct {
		ch       byte
		isLetter bool
	}{
		{'0', false},
		{'a', true},
		{' ', false},
		{'A', true},
		{'_', true},
	}

	for i, tc := range testCases {
		re := isLetter(tc.ch)
		if re != tc.isLetter {
			t.Errorf("tests[%d] - error isDigit: expected %v, got %v.\n", i, tc.isLetter, re)
		}
	}
}

func TestLexer_readNumber(t *testing.T) {
	testCases := []struct {
		input          string
		expectedNumber string
	}{
		{"32 fds", "32"},
		{"32", "32"},
	}
	var lex *Lexer
	for i, tc := range testCases {
		lex = New(tc.input)
		re := lex.readNumber()
		if re != tc.expectedNumber {
			t.Errorf("tests[%d] - error reading number: expected %s, got %s.\n", i, tc.expectedNumber, re)
		}
	}
}

func Test_isDigit(t *testing.T) {
	testCases := []struct {
		ch      byte
		isDigit bool
	}{
		{'0', true},
		{'a', false},
		{' ', false},
	}

	for i, tc := range testCases {
		re := isDigit(tc.ch)
		if re != tc.isDigit {
			t.Errorf("tests[%d] - error isDigit: expected %v, got %v.\n", i, tc.isDigit, re)
		}
	}
}

func TestLexer_skipWhiteSpace(t *testing.T) {
	testCases := []struct {
		input            string
		expectedPosition int
	}{
		{" fds", 1},
		{"   fds", 3},
		{"\nfds", 1},
		{"\tfds", 1},
		{"\rfds", 1},
	}
	var lex *Lexer
	for i, tc := range testCases {
		lex = New(tc.input)
		lex.skipWhiteSpace()
		if lex.position != tc.expectedPosition {
			t.Errorf("tests[%d] - error reading number: expected %d, got %d.\n", i, tc.expectedPosition, lex.position)
		}
	}
}

func TestLexer_NextToken(t *testing.T) {
	tokenTests := []struct {
		input                 string
		expectedParsedResults []struct {
			expectedType    token.LexicalType
			expectedLiteral string
		}
	}{
		{
			input: `=+(){},;`,
			expectedParsedResults: []struct {
				expectedType    token.LexicalType
				expectedLiteral string
			}{
				{expectedType: token.ASSIGN, expectedLiteral: "="},
				{expectedType: token.PLUS, expectedLiteral: "+"},
				{expectedType: token.LPAREN, expectedLiteral: "("},
				{expectedType: token.RPAREN, expectedLiteral: ")"},
				{expectedType: token.LBRACE, expectedLiteral: "{"},
				{expectedType: token.RBRACE, expectedLiteral: "}"},
				{expectedType: token.COMMA, expectedLiteral: ","},
				{expectedType: token.SEMICOLON, expectedLiteral: ";"},
				{expectedType: token.EOF, expectedLiteral: ""},
			},
		},
		{
			input: `var five = 5;
					var ten = 10;
					var add = fn(x, y) {
     					return x + y;
					};
   					var result = add(five, ten);`,
			expectedParsedResults: []struct {
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
				{expectedType: token.RPAREN, expectedLiteral: ")"},
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
				{expectedType: token.RPAREN, expectedLiteral: ")"},
				{expectedType: token.SEMICOLON, expectedLiteral: ";"},
				{expectedType: token.EOF, expectedLiteral: ""},
			},
		},
		{
			// even though !-/*5 doesn't make sense, the lexer still has to parse it into token and let other mechanism check if the logic is correct.
			input: `var five = 5;
					!-/*5^; 
   					5 < 10 > 5;`,
			expectedParsedResults: []struct {
				expectedType    token.LexicalType
				expectedLiteral string
			}{
				// var five = 5;
				{expectedType: token.VAR, expectedLiteral: "var"},
				{expectedType: token.IDENT, expectedLiteral: "five"},
				{expectedType: token.ASSIGN, expectedLiteral: "="},
				{expectedType: token.INT, expectedLiteral: "5"},
				{expectedType: token.SEMICOLON, expectedLiteral: ";"},

				// !-/*5^;
				{expectedType: token.EXCLAMATION, expectedLiteral: "!"},
				{expectedType: token.MINUS, expectedLiteral: "-"},
				{expectedType: token.SLASH, expectedLiteral: "/"},
				{expectedType: token.ASTERISK, expectedLiteral: "*"},
				{expectedType: token.INT, expectedLiteral: "5"},
				{expectedType: token.CARET, expectedLiteral: "^"},
				{expectedType: token.SEMICOLON, expectedLiteral: ";"},

				// 5 < 10 > 5;
				{expectedType: token.INT, expectedLiteral: "5"},
				{expectedType: token.LESSTHAN, expectedLiteral: "<"},
				{expectedType: token.INT, expectedLiteral: "10"},
				{expectedType: token.GREATERTHAN, expectedLiteral: ">"},
				{expectedType: token.INT, expectedLiteral: "5"},
				{expectedType: token.SEMICOLON, expectedLiteral: ";"},
				{expectedType: token.EOF, expectedLiteral: ""},
			},
		},
		{
			input: `fn(a, b){
						if (a > b) {
							return true;
						} else {
							return false;
						}
					}`,
			expectedParsedResults: []struct {
				expectedType    token.LexicalType
				expectedLiteral string
			}{
				// fn(a, b) {
				{expectedType: token.FUNCTION, expectedLiteral: "fn"},
				{expectedType: token.LPAREN, expectedLiteral: "("},
				{expectedType: token.IDENT, expectedLiteral: "a"},
				{expectedType: token.COMMA, expectedLiteral: ","},
				{expectedType: token.IDENT, expectedLiteral: "b"},
				{expectedType: token.RPAREN, expectedLiteral: ")"},
				{expectedType: token.LBRACE, expectedLiteral: "{"},

				// if (a > b) {
				{expectedType: token.IF, expectedLiteral: "if"},
				{expectedType: token.LPAREN, expectedLiteral: "("},
				{expectedType: token.IDENT, expectedLiteral: "a"},
				{expectedType: token.GREATERTHAN, expectedLiteral: ">"},
				{expectedType: token.IDENT, expectedLiteral: "b"},
				{expectedType: token.RPAREN, expectedLiteral: ")"},
				{expectedType: token.LBRACE, expectedLiteral: "{"},

				// return true;
				{expectedType: token.RETURN, expectedLiteral: "return"},
				{expectedType: token.TRUE, expectedLiteral: "true"},
				{expectedType: token.SEMICOLON, expectedLiteral: ";"},

				// } else {
				{expectedType: token.RBRACE, expectedLiteral: "}"},
				{expectedType: token.ELSE, expectedLiteral: "else"},
				{expectedType: token.LBRACE, expectedLiteral: "{"},

				// return false;
				{expectedType: token.RETURN, expectedLiteral: "return"},
				{expectedType: token.FALSE, expectedLiteral: "false"},
				{expectedType: token.SEMICOLON, expectedLiteral: ";"},

				// }
				{expectedType: token.RBRACE, expectedLiteral: "}"},

				// }
				{expectedType: token.RBRACE, expectedLiteral: "}"},
			},
		},
	}

	l := new(Lexer)
	for i, tt := range tokenTests {
		l = New(tt.input)
		for tokenIdx, pr := range tt.expectedParsedResults {
			tok := l.NextToken()
			if tok.Type != pr.expectedType {
				// %q is a single-quoted character literal safely escaped with Go syntax.
				t.Errorf("tests[%d], token[%d] - error token type: expected %q, got %q.\n", i, tokenIdx, pr.expectedType, tok.Type)
			}

			if tok.Literal != pr.expectedLiteral {
				t.Errorf("tests[%d], token[%d] - error token literal value: expected %q, got %q.\n", i, tokenIdx, pr.expectedLiteral, tok.Literal)
			}
		}
	}
}
