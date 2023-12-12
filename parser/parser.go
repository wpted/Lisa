package parser

import (
	"Lisa/ast"
	token "Lisa/lexToken"
	"Lisa/lexer"
	"fmt"
)

// Parser is a component that takes the input data, and builds a data structure, checking for correct syntax in the process.
// Out parser utilizes a *lexer.Lexer and is responsible for building an *ast.ProgramRoot (which is a tree) from the input.
type Parser struct {
	// l is a pointer to an instance of the lexer.
	l      *lexer.Lexer
	errors []string

	// curToken points to the current token.
	curToken *token.Token
	// nextToken points to the next token. We need to examine the next token to decide what to do next with the current token.
	// For example, when encounter a curToken token.INT, we need to decide to end the line(var x = 5),
	// or it actually follows with an arithmetic expression (var y = 5 + 6;).
	nextToken *token.Token
}

// New initializes a Parser instance.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:         l,
		errors:    make([]string, 0),
		curToken:  nil,
		nextToken: nil,
	}

	// Read two tokens, so curToken and nextToken are both set.
	// curToken: nil, nextToken: nil.
	// curToken: nil, nextToken: token(0).
	// curToken: token(0), nextToken: token(1).

	// When input is empty in the given lexer, curToken and nextToken are both type Token.EOF (readPosition >= len(l.input)).
	p.readNextToken()
	p.readNextToken()

	return p
}

// Errors returns all the errors occur when parsing an input with a Parser.
func (p *Parser) Errors() []string {
	return p.errors
}

// readNextToken helps advances the token in the lexer by 1, simultaneously setting the curToken and nextToken.
func (p *Parser) readNextToken() {
	p.curToken = p.nextToken
	// Get the lexical transformation of character to token.
	p.nextToken = p.l.ReadNextToken()
}

// ParseProgram parses the input and return a converted AST tree.
func (p *Parser) ParseProgram() *ast.ProgramRoot {
	// Create an empty ast.ProgramRoot.
	astRoot := &ast.ProgramRoot{
		// ast.Statement is an interface.
		Statements: make([]ast.Statement, 0),
	}

	// Determine when to stop the recursive parsing -> When current token type is token.EOF.
	for !p.curTokenTypeIs(token.EOF) {
		// Parse the statement one by one from the input and append it to astRoot.
		stmt := p.parseStatement()
		// Nil stmt means that there's something wrong in the given input.
		if stmt != nil {
			astRoot.Statements = append(astRoot.Statements, stmt)
		}
		p.readNextToken()
	}

	return astRoot
}

// parseStatement is a helper function that checks the type of the curToken and determine what statement to return to ParseProgram.
// parseStatement stops when reaches a ';' (type token.SEMICOLON).
func (p *Parser) parseStatement() ast.Statement {

	// Check what statement we're currently parsing.
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	//case token.RETURN:
	default:
		return nil
	}
}

// parseVarStatement parses a statement that starts with 'var' and ends with ';'.
// If there's any elements missing from a standard 'var x = 5;' statement, the parser stores all the error in errors and returns a nil ast.VarStatement.
func (p *Parser) parseVarStatement() *ast.VarStatement {
	invalidStmt := false

	// When parsing 'var x = 5;'

	// 1. Assign 'var' token to the 'Token' field for the VarStatement.
	stmt := &ast.VarStatement{Token: p.curToken}

	// 2. Peek the next token, which we expect it is the variable name 'x', a type token.IDENT.
	// If the token is an expected type, advance the token pointer.
	if !p.expectNext(token.IDENT) {
		// Not an expected Identifier, store the error.
		p.storeNextTokenTypeError(token.IDENT)
		invalidStmt = true
		// Since we're missing the identifier, we expect the next token type to be a type token.ASSIGN.
	}

	// Set the name of stmt to the parsed variable name.
	stmt.Name = &ast.IdentifierExpression{
		// Since we already advanced the token when checking whether the token is a type token.IDENT,
		// p.curToken is now the name of the variable. ('x' in 'var x = 5;')
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	// 3. Peek the next token, which we expect it is a '=', a type token.ASSIGN.
	// If the token is an expected type, advance the token pointer.
	if !p.expectNext(token.ASSIGN) {
		p.storeNextTokenTypeError(token.ASSIGN)
		invalidStmt = true
	}

	// 4. Check and assign the expression to stmt.Value.
	p.readNextToken()

	// 5. Check the semicolon at the end of a var statement.
	if !p.expectNext(token.SEMICOLON) {
		// Not an expected Identifier, store the error.
		p.storeNextTokenTypeError(token.SEMICOLON)
		invalidStmt = true
	}

	// After reached a supposed ';' position return the parsed statement, check whether the statement is valid.
	if invalidStmt {
		return nil
	}
	return stmt
}

// curTokenTypeIs checks whether the type of the current token is identical as the given type.
func (p *Parser) curTokenTypeIs(expectType token.LexicalType) bool {
	return p.curToken.Type == expectType
}

// nextTokenTypeIs checks whether the type of the next token is identical as the given type.
func (p *Parser) nextTokenTypeIs(expectType token.LexicalType) bool {
	return p.nextToken.Type == expectType
}

// expectNext checks the next token type, and if the nextToken is identical as the given type, advance the pointer and return true.
func (p *Parser) expectNext(expectType token.LexicalType) bool {
	if p.nextTokenTypeIs(expectType) {
		// If the token is the expected type, advance to the token.
		p.readNextToken()
		return true
	}
	return false
}

func (p *Parser) storeNextTokenTypeError(expectType token.LexicalType) {
	errMsg := fmt.Sprintf("error next token type: expected TYPE(%s), got TYPE(%s).", expectType, p.nextToken.Type)
	p.errors = append(p.errors, errMsg)
}

//func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
//	return nil
//}
//
//func (p *Parser) parseExpression() ast.Expression {
//	return nil
//}
//
