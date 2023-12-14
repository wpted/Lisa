package parser

import (
    "Lisa/ast"
    "Lisa/lexer"
    "testing"
)

func TestParser_ParseStatement(t *testing.T) {
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

            // 2. TODO: Check the Value of the expression.
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

    t.Run("Correct 'Return' statements", func(t *testing.T) {
        input := `return 5;
			      return 10;`

        l := lexer.New(input)
        p := New(l)
        astRoot := p.ParseProgram()
        if astRoot == nil {
            t.Errorf("Error parsing program - nil program root.\n")
        }

        if len(astRoot.Statements) != 2 {
            t.Errorf("Error statement length for program root: expected %d, got %d.", 3, len(astRoot.Statements))
        }

        for _, stmt := range astRoot.Statements {
            if stmt.TokenLiteral() != "return" {
                t.Errorf("Error s.TokenLiteral: expected %q, got %q.\n", "return", stmt.TokenLiteral())
            }

            _, ok := stmt.(*ast.ReturnStatement)
            if !ok {
                t.Errorf("Error statement type: expected *ast.ReturnStatement, got %T.\n", stmt)
            }
        }
    })

    t.Run("Incorrect 'Return' statements", func(t *testing.T) {
        input := `return ;
			      return 10`

        l := lexer.New(input)
        p := New(l)
        astRoot := p.ParseProgram()
        if astRoot == nil {
            t.Errorf("Error parsing program - nil program root.\n")
        }

        if len(astRoot.Statements) != 2 {
            t.Errorf("Error statement length for program root: expected %d, got %d.", 2, len(astRoot.Statements))
        }
        if len(p.errors) != 2 {
            t.Errorf("Error length of expected errors, expected %d, got %d.", 2, len(p.errors))
        }
    })
}

func TestParser_ParseExpression(t *testing.T) {
    t.Run("Test Expression - Identifiers", func(t *testing.T) {
        input := "foobar;"
        l := lexer.New(input)
        p := New(l)
        astRoot := p.ParseProgram()
        if len(astRoot.Statements) != 1 {
            t.Errorf("Error statement length for program root: expected %d, got %d.", 1, len(astRoot.Statements))
        }

        if len(astRoot.Statements) > 0 {
            // Check if a statement in the parsed statements is an *ast.ExpressionStatement.
            stmt, ok := astRoot.Statements[0].(*ast.ExpressionStatement)
            if !ok {
                t.Errorf("Error statement type: expected *ast.ExpressionStatement, got %T", stmt)
            } else {
                // If a statement is not an *ast.ExpressionStatement type, calling Expression field will lead to runtime panic.
                // Check if expression in the Expression field is an *ast.IdentifierExpression.
                ident, ok := stmt.Expression.(*ast.IdentifierExpression)
                if !ok {
                    t.Errorf("Error expression type: expected *ast.IdentifierExpression, got %T", ident)
                } else {
                    if ident.Value != "foobar" {
                        t.Errorf("Error expression value in Identifier: expected %s, got %s.", "foobar", ident.Value)
                    }

                    if ident.TokenLiteral() != "foobar" {
                        t.Errorf("Error expression TokenLiteral() in Identifier: expected %s, got %s.", "foobar", ident.Value)
                    }
                }
            }
        }
    })

    t.Run("Test Expression - Integer Literals", func(t *testing.T) {
        input := "5;"
        l := lexer.New(input)
        p := New(l)
        astRoot := p.ParseProgram()
        if len(astRoot.Statements) != 1 {
            t.Errorf("Error statement length for program root: expected %d, got %d.", 1, len(astRoot.Statements))
        }

        if len(astRoot.Statements) > 0 {
            // Check if a statement in the parsed statements is an *ast.ExpressionStatement.
            stmt, ok := astRoot.Statements[0].(*ast.ExpressionStatement)
            if !ok {
                t.Errorf("Error statement type: expected *ast.ExpressionStatement, got %T", stmt)
            } else {
                // If a statement is not an *ast.ExpressionStatement type, calling Expression field will lead to runtime panic.
                // Check if expression in the Expression field is an *ast.IdentifierExpression.
                intLitExp, ok := stmt.Expression.(*ast.IntegerLiteralExpression)
                if !ok {
                    t.Errorf("Error expression type: expected *ast.IdentifierExpression, got %T", intLitExp)
                } else {
                    if intLitExp.Value != int64(5) {
                        t.Errorf("Error expression value in Identifier: expected %d, got %d.", int64(5), intLitExp.Value)
                    }

                    if intLitExp.TokenLiteral() != "5" {
                        t.Errorf("Error expression TokenLiteral() in Identifier: expected %s, got %s.", "5", "5")
                    }
                }
            }
        }
    })

    t.Run("Test Expression - Prefix Operators", func(t *testing.T) {

    })

    t.Run("Test Expression - Infix Operators", func(t *testing.T) {

    })

}
