package ast

import (
	token "Lisa/lexToken"
)

// Node is the base element of an AST tree.
type Node interface {
	// TokenLiteral is the literal of the parsed node.
	TokenLiteral() string
}

// Statement node is an instruction ( action ) node. It performs actions and doesn't produce value.
type Statement interface {
	// Node is embedded to make sure all nodes in the AST fits our rule for nodes.
	Node
	// statementNode is a dummy method to distinguish itself from Expression nodes.
	statementNode()
}

// Expression node indicates that a value is being resolved to it.
// 'var a = 5' and 'var a = add(2, 3)' is in fact the same thing, resolving value 5 to 'a'.
type Expression interface {
	// Node is embedded to make sure all nodes in the AST fits our rule for nodes.
	Node
	// expressionNode is a dummy method to distinguish itself from Statement nodes.
	expressionNode()
}

// ProgramRoot is the root node of every AST our parser produces.
type ProgramRoot struct {
	// Statements stores a series of statements (which is an interface, any node that fits Statement counts.) that is contained in our program.
	Statements []Statement
}

func (p *ProgramRoot) TokenLiteral() string {
	if len(p.Statements) > 0 {
		// Return token literal for the first statement (whichever node that fits Statement interface).
		return p.Statements[0].TokenLiteral()
	}
	// No statements stored in root node.
	return ""
}

// VarStatement node. Should hold the Name of the identifier, the value for the expression, and its own Token.
type VarStatement struct {
	// Token is the token.Var token.
	Token *token.Token
	// Name is the name of the variable.
	Name *Identifier
	// Value is the field that points to the expression on the right side of the equal sign.
	Value Expression
}

func (v *VarStatement) TokenLiteral() string { return v.Token.Literal }

// statementNode categorizes VarStatement node as a statement node.
func (v *VarStatement) statementNode() {}

// Identifier node. In the context of a programming language's abstract syntax tree (AST),
// an identifier is typically associated with a declaration or a statement node.
// This is because an identifier is often used to name and reference variables, functions, or other program entities,
// and these entities are typically introduced or declared through statements.
type Identifier struct {
	Token *token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string {
	// Return the parsed literal from the token.
	return i.Token.Literal
}

// expressionNode categorizes Identifier as an expression node.
func (i *Identifier) expressionNode() {}
