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
		},
	}

	if program.String() != "let myVar = anotherVar;// this is a comment" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
