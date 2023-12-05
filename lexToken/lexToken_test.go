package lexToken

import "testing"

func Test_LookUpReservedWord(t *testing.T) {
	testCases := []struct {
		word    string
		lexType LexicalType
	}{
		{"fn", FUNCTION},
		{"var", VAR},
		{"hello", IDENT},
	}

	for i, tc := range testCases {
		lt := LookUpReservedWord(tc.word)
		if lt != tc.lexType {
			t.Errorf("tests[%d] - error looking up reserved word: expected %s, got %s.\n", i, tc.lexType, lt)
		}
	}
}
