package parser

import (
	"Lisa/ast"
	token "Lisa/lexToken"
	"Lisa/lexer"
	"fmt"
	"strconv"
)

const (
	// The following is the precedence for operator orders.
	_ int = iota
	// LOWEST is the lowest precedence for operator orders.
	LOWEST
	// EQUALS is the order of a '==' operator
	EQUALS
	// LESSGREATER is the order of '>' or '<' operator
	LESSGREATER
	// SUM is the order of a '+' operator
	SUM
	PRODUCT
	PREFIX
	CALL
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

	prefixParseFns map[token.LexicalType]prefixParseFn
	infixParseFns  map[token.LexicalType]infixParseFn
}

// New initializes a Parser instance.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         make([]string, 0),
		curToken:       nil,
		nextToken:      nil,
		prefixParseFns: make(map[token.LexicalType]prefixParseFn),
		infixParseFns:  make(map[token.LexicalType]infixParseFn),
	}

	// Read two tokens, so curToken and nextToken are both set.
	// curToken: nil, nextToken: nil.
	// curToken: nil, nextToken: token(0).
	// curToken: token(0), nextToken: token(1).

	// When input is empty in the given lexer, curToken and nextToken are both type Token.EOF (readPosition >= len(l.input)).
	p.readNextToken()
	p.readNextToken()

	// Register parser functions for parsing expressions.
	p.registerParserFunctionForPrefix(token.IDENT, p.parseIdentifier)
	p.registerParserFunctionForPrefix(token.INT, p.parseIntegerLiteral)
	p.registerParserFunctionForPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerParserFunctionForPrefix(token.EXCLAMATION, p.parsePrefixExpression)

	return p
}

// Errors returns all the errors occur when parsing an input with a Parser.
func (p *Parser) Errors() []string {
	return p.errors
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

// readNextToken helps advances the token in the lexer by 1, simultaneously setting the curToken and nextToken.
func (p *Parser) readNextToken() {
	p.curToken = p.nextToken
	// Get the lexical transformation of character to token.
	p.nextToken = p.l.ReadNextToken()
}

// parseStatement is a helper function that checks the type of the curToken and determine what statement to return to ParseProgram.
// parseStatement stops when reaches a ';' (type token.SEMICOLON).
func (p *Parser) parseStatement() ast.Statement {

	// Check what statement we're currently parsing.
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseVarStatement parses a statement that starts with 'var' and ends with ';', (e.g 'var x = 5;').
// If there's any elements missing from a standard 'var x = 5;' statement, the parser stores all the error in errors and returns a nil ast.VarStatement.
func (p *Parser) parseVarStatement() *ast.VarStatement {
	var stmtInvalid bool

	// When parsing 'var x = 5;'

	// 1. Assign 'var' token to the 'Token' field for the VarStatement.
	stmt := &ast.VarStatement{Token: p.curToken}

	// 2. Peek the next token, which we expect it is the variable name 'x', a type token.IDENT.
	// If the token is an expected type, advance the token pointer.
	if !p.expectNext(token.IDENT) {
		// Not an expected Identifier, store the error.
		p.storeNextTokenTypeError(token.IDENT)
		stmtInvalid = true
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
		stmtInvalid = true
	}

	// 4. TODO: Check and assign the expression to stmt.Value.
	p.readNextToken()

	// 5. Check the semicolon at the end of a var statement.
	if !p.expectNext(token.SEMICOLON) {
		// Not an expected Identifier, store the error.
		p.storeNextTokenTypeError(token.SEMICOLON)
		stmtInvalid = true
	}

	// After reached a supposed ';' position return the parsed statement, check whether the statement is valid.
	if stmtInvalid {
		return nil
	}
	return stmt
}

// parseReturnStatement parses a statement that starts with a 'return' and ends with a ';', (e.g 'return 5;').
// If there's any elements missing from a standard 'return 5;' statement, the parser stores all the error in errors and returns a nil ast.ReturnStatement.
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	var stmtInvalid bool

	// 1. Assign currentToken 'return' to p.curToken.
	stmt := &ast.ReturnStatement{
		Token:       p.curToken,
		ReturnValue: nil,
	}

	// 2. TODO: Check and assign the expression of the value of return.
	p.readNextToken()

	// 3. Check the semicolon at the end of a return statement.
	if !p.expectNext(token.SEMICOLON) {
		p.storeNextTokenTypeError(token.SEMICOLON)
		stmtInvalid = true
	}

	// After reached a supposed ';' position return the parsed statement, check whether the statement is valid.
	if stmtInvalid {
		return nil
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	var stmtInvalid bool
	stmt := &ast.ExpressionStatement{
		Token: p.curToken,
	}
	stmt.Expression = p.parseExpression(LOWEST)
	if !p.expectNext(token.SEMICOLON) {
		p.storeNextTokenTypeError(token.SEMICOLON)
		stmtInvalid = true
	}

	// After reached a supposed ';' position return the parsed statement, check whether the statement is valid.
	if stmtInvalid {
		return nil
	}
	return stmt
}

func (p *Parser) parseExpression(operatorPrecedence int) ast.Expression {
	// 1. Integer literals (5;)
	// 2. Identifiers (foobar;) -> (PrefixParseFn)

	// Fetch the parser function from the pre-registered functions.
	parseFn := p.prefixParseFns[p.curToken.Type]
	if parseFn != nil {
		expression := parseFn()
		return expression
	}
	// 3. Prefix Operators (-5;)
	// 4. Infix Operators (5 * 2;)
	// - Normal binary operators (+, -, *, /)
	// - Comparison Operators (x >= y; x == y;)
	// - Grouped parentheses to group and reorder evaluation (2 + (5 * 2);)
	// 5. Function call expression (foobar();)
	// 6. Function literals (fn(a, b){return a + b;};)
	return nil
}

// parseIdentifier turns the current token from a parser to an *ast.IdentifierExpression, returned as an ast.Expression interface.
// This function should be registered when starting a new parser, and should be called when parser encounter a token of type token.IDENT.
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.IdentifierExpression{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

// parseIntegerLiteral turns the current token from a parser to an *ast.IntegerLiteralExpression, returned as an ast.Expression interface.
// This function should be registered when starting a new parser, and should be called when parser encounter a token of type token.INT.
func (p *Parser) parseIntegerLiteral() ast.Expression {
	exp := &ast.IntegerLiteralExpression{
		Token: p.curToken,
	}

	// Convert the token of a token literal to a int64 value.
	value, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		errMsg := fmt.Sprintf("error could not parse %q as integer", p.curToken.Literal)
		p.storeParseTokenError(errMsg)
		return nil
	}

	// After successfully reading from the string literal, assign the value to the integer literal expression.
	exp.Value = int64(value)

	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Operator: p.curToken.Literal,
		Token:    p.curToken,
	}

	// Advance the pointer of the lexer, now it's the token after the prefix.
	p.readNextToken()
	exp.RightToken = p.parseExpression(PREFIX)
	return exp
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
		// If the token is the expected type, advance to the expected token.
		p.readNextToken()
		return true
	}
	return false
}

func (p *Parser) storeNextTokenTypeError(expectType token.LexicalType) {
	errMsg := fmt.Sprintf("error next token type: expected TYPE(%s), got TYPE(%s).", expectType, p.nextToken.Type)
	p.errors = append(p.errors, errMsg)
}

func (p *Parser) storeParseTokenError(errMsg string) {
	p.errors = append(p.errors, errMsg)
}

type prefixParseFn func() ast.Expression
type infixParseFn func(expression ast.Expression) ast.Expression

func (p *Parser) registerParserFunctionForPrefix(tokenType token.LexicalType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerParserFunctionForInfix(tokenType token.LexicalType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
