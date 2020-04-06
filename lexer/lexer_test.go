package lexer

import (
	"fmt"
	"testing"

	"monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	let ten = 10;

	let add = fn(x, y) {
		x + y;
	};

	let result = add(five, ten);
	!-/*5;
	5 < 10 > 5;

	if (5 < 10) {
		return true;
	} else {
		return false;
	}

	10 == 10;
	10 != 9;
	x -= 9;
	x += 9;
	x /= 9;
	x *= 9;
	`
	tests := token.Tokens{
		token.Token{Type: token.LET, Literal: "let"},
		token.Token{Type: token.IDENT, Literal: "five"},
		token.Token{Type: token.ASSIGN, Literal: "="},
		token.Token{Type: token.INT, Literal: "5"},
		token.Token{Type: token.SEMICOLON, Literal: ";"},

		token.Token{Type: token.LET, Literal: "let"},
		token.Token{Type: token.IDENT, Literal: "ten"},
		token.Token{Type: token.ASSIGN, Literal: "="},
		token.Token{Type: token.INT, Literal: "10"},
		token.Token{Type: token.SEMICOLON, Literal: ";"},

		token.Token{Type: token.LET, Literal: "let"},
		token.Token{Type: token.IDENT, Literal: "add"},
		token.Token{Type: token.ASSIGN, Literal: "="},
		token.Token{Type: token.FUNCTION, Literal: "fn"},
		token.Token{Type: token.LPAREN, Literal: "("},
		token.Token{Type: token.IDENT, Literal: "x"},
		token.Token{Type: token.COMMA, Literal: ","},
		token.Token{Type: token.IDENT, Literal: "y"},
		token.Token{Type: token.RPAREN, Literal: ")"},
		token.Token{Type: token.LBRACE, Literal: "{"},
		token.Token{Type: token.IDENT, Literal: "x"},
		token.Token{Type: token.PLUS, Literal: "+"},
		token.Token{Type: token.IDENT, Literal: "y"},
		token.Token{Type: token.SEMICOLON, Literal: ";"},
		token.Token{Type: token.RBRACE, Literal: "}"},
		token.Token{Type: token.SEMICOLON, Literal: ";"},

		token.Token{Type: token.LET, Literal: "let"},
		token.Token{Type: token.IDENT, Literal: "result"},
		token.Token{Type: token.ASSIGN, Literal: "="},
		token.Token{Type: token.IDENT, Literal: "add"},
		token.Token{Type: token.LPAREN, Literal: "("},
		token.Token{Type: token.IDENT, Literal: "five"},
		token.Token{Type: token.COMMA, Literal: ","},
		token.Token{Type: token.IDENT, Literal: "ten"},
		token.Token{Type: token.RPAREN, Literal: ")"},
		token.Token{Type: token.SEMICOLON, Literal: ";"},

		token.Token{Type: token.BANG, Literal: "!"},
		token.Token{Type: token.MINUS, Literal: "-"},
		token.Token{Type: token.SLASH, Literal: "/"},
		token.Token{Type: token.ASTERISK, Literal: "*"},
		token.Token{Type: token.INT, Literal: "5"},
		token.Token{Type: token.SEMICOLON, Literal: ";"},

		token.Token{Type: token.INT, Literal: "5"},
		token.Token{Type: token.LT, Literal: "<"},
		token.Token{Type: token.INT, Literal: "10"},
		token.Token{Type: token.GT, Literal: ">"},
		token.Token{Type: token.INT, Literal: "5"},
		token.Token{Type: token.SEMICOLON, Literal: ";"},

		token.Token{Type: token.IF, Literal: "if"},
		token.Token{Type: token.LPAREN, Literal: "("},
		token.Token{Type: token.INT, Literal: "5"},
		token.Token{Type: token.LT, Literal: "<"},
		token.Token{Type: token.INT, Literal: "10"},
		token.Token{Type: token.RPAREN, Literal: ")"},
		token.Token{Type: token.LBRACE, Literal: "{"},
		token.Token{Type: token.RETURN, Literal: "return"},
		token.Token{Type: token.TRUE, Literal: "true"},
		token.Token{Type: token.SEMICOLON, Literal: ";"},
		token.Token{Type: token.RBRACE, Literal: "}"},
		token.Token{Type: token.ELSE, Literal: "else"},
		token.Token{Type: token.LBRACE, Literal: "{"},
		token.Token{Type: token.RETURN, Literal: "return"},
		token.Token{Type: token.FALSE, Literal: "false"},
		token.Token{Type: token.SEMICOLON, Literal: ";"},
		token.Token{Type: token.RBRACE, Literal: "}"},

		token.Token{Type: token.INT, Literal: "10"},
		token.Token{Type: token.EQ, Literal: "=="},
		token.Token{Type: token.INT, Literal: "10"},
		token.Token{Type: token.SEMICOLON, Literal: ";"},

		token.Token{Type: token.INT, Literal: "10"},
		token.Token{Type: token.NOTEQ, Literal: "!="},
		token.Token{Type: token.INT, Literal: "9"},
		token.Token{Type: token.SEMICOLON, Literal: ";"},

		token.Token{Type: token.IDENT, Literal: "x"},
		token.Token{Type: token.MINUSEQ, Literal: "-="},
		token.Token{Type: token.INT, Literal: "9"},
		token.Token{Type: token.SEMICOLON, Literal: ";"},

		token.Token{Type: token.IDENT, Literal: "x"},
		token.Token{Type: token.PLUSEQ, Literal: "+="},
		token.Token{Type: token.INT, Literal: "9"},
		token.Token{Type: token.SEMICOLON, Literal: ";"},

		token.Token{Type: token.IDENT, Literal: "x"},
		token.Token{Type: token.SLASHEQ, Literal: "/="},
		token.Token{Type: token.INT, Literal: "9"},
		token.Token{Type: token.SEMICOLON, Literal: ";"},
		token.Token{Type: token.IDENT, Literal: "x"},

		token.Token{Type: token.ASTERISKEQ, Literal: "*="},
		token.Token{Type: token.INT, Literal: "9"},
		token.Token{Type: token.SEMICOLON, Literal: ";"},
		token.Token{Type: token.EOF, Literal: ""},
	}

	l := NewLexer(input)
	for i, tt := range tests {
		tok := l.NextToken()
		fmt.Println(tok)
		fmt.Println(tt)
		if tok.Type != tt.Type {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.Type, tok.Type)
		}
		if tok.Literal != tt.Literal {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.Literal, tok.Literal)
		}
	}
}
