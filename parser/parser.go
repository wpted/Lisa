package parser

import (
	"Lisa/ast"
	token "Lisa/lexToken"
	"Lisa/lexer"
)

// Parser is a component that takes the input data, and builds a data structure, checking for correct syntax in the process.
// Out parser utilizes a *lexer.Lexer and is responsible for building an *ast.ProgramRoot (which is a tree) from the input.
type Parser struct {
	// l is a pointer to an instance of the lexer.
	l *lexer.Lexer

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
	for p.curToken.Type != token.EOF {
		// Parse the statement one by one from the input and append it to astRoot.
		stmt := p.parseStatement()
		if stmt != nil {
			astRoot.Statements = append(astRoot.Statements, stmt)
		}
		p.readNextToken()
	}

	return astRoot
}

func (p *Parser) parseStatement() ast.Statement {
	return nil
}
