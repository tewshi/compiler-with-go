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

// tokenLiteral the literal vale of the token
func (p *Program) tokenLiteral() string {
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

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) tokenLiteral() string { return ls.Token.Literal }

// Identifier represents an identifier in a statement
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) tokenLiteral() string { return i.Token.Literal }
