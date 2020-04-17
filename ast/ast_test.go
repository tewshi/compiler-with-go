package ast

import (
	"monkey/token"
	"testing"
)

// TODO: add more tests here
func TestString(t *testing.T) {
	program := &Program{
		Statements: Statements{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			}, &CommentLiteral{
				Token: token.Token{Type: token.COMMENT, Literal: "// this is a comment"},
				Value: "this is a comment",
			},
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "x"},
					Value: "x",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
			},
			&ExpressionStatement{
				Expression: &InfixExpression{
					Token:    token.Token{Type: token.PLUSEQ, Literal: token.PLUSEQ},
					Left:     &Identifier{Token: token.Token{Type: token.IDENT, Literal: "x"}, Value: "x"},
					Operator: token.PLUSEQ,
					Right:    &Identifier{Token: token.Token{Type: token.IDENT, Literal: "x"}, Value: "x"},
				},
			},
			&ExpressionStatement{
				Expression: &SuffixExpression{
					Token:    token.Token{Type: token.INCREMENT, Literal: token.INCREMENT},
					Left:     &Identifier{Token: token.Token{Type: token.IDENT, Literal: "x"}, Value: "x"},
					Operator: token.INCREMENT,
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;let x = myVar;(x += x)(x++)" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
