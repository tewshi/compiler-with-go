package ast

import "monkey/token"

// Node a tree node
type Node interface {
	TokenLiteral() string
}

// Expression a single program expression
type Expression interface {
	Node
	expressionNode() string
}

// Statement a single program statement
type Statement interface {
	Node
	statementNode() string
}

// Nodes list of nodes
type Nodes []Node

// Expressions list of expressions
type Expressions []Expression

// Statements list of statements
type Statements []Statement

// Program root of the AST
type Program struct {
	Statements Statements
}

// TokenLiteral the literal value of the token
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// LetStatement represents a let statement in the AST
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral the literal value of the let statement token
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// Identifier represents an identifier in a statement
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral the literal value of the identifier token
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
