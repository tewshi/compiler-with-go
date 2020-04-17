package ast

import (
	"bytes"
	"monkey/token"
	"strconv"
	"strings"
)

// TAB the tab character
const TAB = "  "

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

// BlockStatement represents a block statement in the AST
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// TokenLiteral the literal value of the block statement token
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

// String string representation of aa block statement
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Identifier represents an identifier in a statement
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

// Identifiers list of identifier struct
type Identifiers []*Identifier

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
func (il *IntegerLiteral) String() string { return strconv.FormatInt(il.Value, 10) }

// DoubleLiteral represents a double in a statement
type DoubleLiteral struct {
	Token     token.Token // the token.INT token
	Value     float64
	Precision int // the double's precision
}

func (dl *DoubleLiteral) expressionNode() {}

// TokenLiteral the literal value of the double token
func (dl *DoubleLiteral) TokenLiteral() string { return dl.Token.Literal }

// String string representation of an double
func (dl *DoubleLiteral) String() string { return strconv.FormatFloat(dl.Value, 'f', dl.Precision, 64) }

// PrefixExpression represents a prefix expression
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !x or -x or ++x or --x
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral the literal value of the prefix expression token
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

// String string representation of a prefix expression
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// SuffixExpression represents a suffix expression
type SuffixExpression struct {
	Token    token.Token // The suffix token, e.g. x++ or x--
	Operator string
	Left     Expression
}

func (se *SuffixExpression) expressionNode() {}

// TokenLiteral the literal value of the suffix expression token
func (se *SuffixExpression) TokenLiteral() string { return se.Token.Literal }

// String string representation of a suffix expression
func (se *SuffixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(se.Left.String())
	out.WriteString(se.Operator)
	out.WriteString(")")
	return out.String()
}

// InfixExpression represents an infix expression
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. 5 + 9;
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode() {}

// TokenLiteral the literal value of the infix expression token
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }

// String string representation of a infix expression
func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")
	return out.String()
}

// IfExpression represents an if expression
type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

// TokenLiteral the literal value of the if expression token
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

// String string representation of an if expression
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if (")
	out.WriteString(ie.Condition.String())

	out.WriteString(") {")
	out.WriteString(TAB + ie.Consequence.String())
	out.WriteString("; }")
	if ie.Alternative != nil {
		out.WriteString(" else {")
		out.WriteString(ie.Alternative.String())
		out.WriteString(";}")
	}
	return out.String()
}

// Boolean the boolean struct
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

// TokenLiteral the literal value of the boolean token
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

// String string representation of a boolean
func (b *Boolean) String() string { return b.Token.Literal }

// StringLiteral represents a string in a statement
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

// TokenLiteral the literal value of the string token
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

// String string representation of a string
func (sl *StringLiteral) String() string { return sl.Token.Literal }

// ArrayLiteral represents a array in a statement
type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements Expressions
}

func (al *ArrayLiteral) expressionNode() {}

// TokenLiteral the literal value of the string token
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }

// String string representation of a string
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// HashLiteral represents a hash in a statement
type HashLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}

// TokenLiteral the literal value of the hash token
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }

// String string representation of a hash
func (hl *HashLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+": "+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// IndexExpression represents an index in a statement: arr[1]
type IndexExpression struct {
	Token token.Token // The [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

// TokenLiteral the literal value of the index expression token
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }

// String string representation of a index expression
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

// FunctionLiteral represents a function in a statement
type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters Identifiers
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

// TokenLiteral the literal value of the function token
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

// String string representation of a function literal
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString(" (")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") { ")
	for _, s := range fl.Body.Statements {
		out.WriteString(s.String() + "; ")
	}
	out.WriteString("};")
	return out.String()
}

// CallExpression represents a function call expression
type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments Expressions
}

func (ce *CallExpression) expressionNode() {}

// TokenLiteral the literal value of the function call token
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

// String string representation of a function call literal
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// CommentLiteral represents a line comment in a statement
type CommentLiteral struct {
	Token token.Token
	Value string
}

func (cl *CommentLiteral) expressionNode() {}

func (cl *CommentLiteral) statementNode() {}

// TokenLiteral the literal value of the line comment token
func (cl *CommentLiteral) TokenLiteral() string { return cl.Token.Literal }

// String string representation of a line comment
func (cl *CommentLiteral) String() string { return cl.Token.Literal }
