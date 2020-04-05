package ast

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

// TokenLiteral the literal vale of the token
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
