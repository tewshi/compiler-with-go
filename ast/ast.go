package ast

import (
	"bytes"
	"monkey/token"
)

// Node a tree node
type Node interface {
	TokenLiteral() string
	String() string
}

// Expression a single program expression
type Expression interface {
	Node
	expressionNode()
}

// Statement a single program statement
type Statement interface {
	Node
	statementNode()
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

// String string representation of the program
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
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

// String string representation of a let statement
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

// ReturnStatement represents a return statement in the AST
type ReturnStatement struct {
	Token token.Token // the token.RETURN token
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral the literal value of the return statement token
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// String string representation of a return statement
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement represents a expression statement in the AST
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral the literal value of the expression statement token
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// String string representation of an expression statement
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// Identifier represents an identifier in a statement
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral the literal value of the identifier token
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// String string representation of an identifier
func (i *Identifier) String() string { return i.Value }

// IntegerLiteral represents an integer in a statement
type IntegerLiteral struct {
	Token token.Token // the token.INT token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral the literal value of the integer token
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// String string representation of an integer
func (il *IntegerLiteral) String() string { return il.Token.Literal }
