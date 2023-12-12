package parser

import (
	"Lisa/ast"
	"Lisa/lexer"
	"testing"
)

func TestParser_ParseProgram(t *testing.T) {
	t.Run("Correct 'Var' statements", func(t *testing.T) {
		input := `var x = 5;
			      var y = 10;
                  var foo_bar = 838383;`

		l := lexer.New(input)
		p := New(l)

		astRoot := p.ParseProgram()
		if astRoot == nil {
			t.Errorf("Error parsing program - nil program root.\n")
		}

		if len(astRoot.Statements) != 3 {
			t.Errorf("Error statement length for program root: expected %d, got %d.", 3, len(astRoot.Statements))
		}

		expectedIdentifier := []string{"x", "y", "foo_bar"}

		for i, stmt := range astRoot.Statements {
			if stmt.TokenLiteral() != "var" {
				t.Errorf("Error s.TokenLiteral: expected %q, got %q.\n", "var", stmt.TokenLiteral())
			}

			varStmt, ok := stmt.(*ast.VarStatement)
			if !ok {
				t.Errorf("Error statement type: expected *ast.VarStatement, got %T.\n", stmt)
			}

			// 1. Check the Identifier.

			// Name is the identifier associated with varStmt.
			// The value of the identifier should be the same as the parsed literal.
			if varStmt.Name.Value != expectedIdentifier[i] {
				t.Errorf("Error VarStatement.Name.Value: expected '%s', got '%s'", expectedIdentifier[i], varStmt.Name.Value)
			}

			// The token.TokenLiteral should be the same as the parsed literal.
			if varStmt.Name.TokenLiteral() != expectedIdentifier[i] {
				t.Errorf("Error token literal: expected '%s', got '%s'.\n", expectedIdentifier[i], varStmt.Name.Value)
			}

			// TODO: 2. Check the Value expression.
		}
	})

	t.Run("Incorrect 'Var' statements", func(t *testing.T) {
		input := `var = 5;
			      var y 10;
                  var foo_bar = ;`

		l := lexer.New(input)
		p := New(l)

		astRoot := p.ParseProgram()
		if astRoot == nil {
			t.Errorf("Error parsing program - nil program root.\n")
		}

		if len(astRoot.Statements) != 3 {
			t.Errorf("Error statement length for program root: expected %d, got %d.", 3, len(astRoot.Statements))
		}

		if len(p.errors) != 3 {
			t.Errorf("Error length of expected errors, expected %d, got %d.", 3, len(p.errors))
		}
	})
}
