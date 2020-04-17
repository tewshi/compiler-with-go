package lexer

import (
	"bytes"
	"monkey/token"
	"strconv"
	"strings"
)

// Lexer the lexer type
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

// NewLexer creates and returns a Lexer
func NewLexer(input string) *Lexer {
	l := &Lexer{input: string(normalizeNewLine([]byte(input)))}
	l.readChar()
	return l
}

// normalizeNewLine normalizes \r\n (windows) and \r (mac) into \n (unix)
func normalizeNewLine(lines []byte) []byte {
	// replace CR LF: \r\n (windows) with LF: \n (unix)
	lines = bytes.ReplaceAll(lines, []byte{13, 10}, []byte{10})
	// replace CR: \r (mac) with LF: \n (unix)
	lines = bytes.ReplaceAll(lines, []byte{13}, []byte{10})
	return lines
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition

	l.readPosition++
}

// peekChar returns the char after the current char if there's one
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken returns the next token
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: "=="}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '.':
		tok = newToken(token.PERIOD, l.ch)
	case '+':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.PLUSEQ, Literal: "+="}
		} else if l.peekChar() == '+' {
			l.readChar()
			tok = token.Token{Type: token.INCREMENT, Literal: "++"}
		} else {
			tok = newToken(token.PLUS, l.ch)
		}
	case '-':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.MINUSEQ, Literal: "-="}
		} else if l.peekChar() == '-' {
			l.readChar()
			tok = token.Token{Type: token.DECREMENT, Literal: "--"}
		} else {
			tok = newToken(token.MINUS, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.NOTEQ, Literal: "!="}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.SLASHEQ, Literal: "/="}
		} else if l.peekChar() == '/' {
			tok.Type = token.COMMENT
			tok.Literal = l.readComment()
		} else {
			tok = newToken(token.SLASH, l.ch)
		}
	case '*':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.ASTERISKEQ, Literal: "*="}
		} else {
			tok = newToken(token.ASTERISK, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.LTEQ, Literal: "<="}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.GTEQ, Literal: ">="}
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case '^':
		tok = token.Token{Type: token.POWER, Literal: "^"}
	case '%':
		tok = token.Token{Type: token.MODULUS, Literal: "%"}
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.LookupIdent()
			return tok
		} else if isDigit(l.ch) {
			position := l.position
			tok.Literal = l.readNumber()
			// for doubles
			if l.ch == '.' {
				l.readChar()
				if _, err := strconv.ParseInt(l.readNumber(), 10, 64); err == nil {
					tok.Literal = l.input[position:l.position]
					tok.Type = token.DOUBLE
				} else {
					tok.Literal += "." + l.NextToken().Literal
					tok.Type = token.ILLEGAL
				}
			} else {
				tok.Type = token.INT
			}
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch)
	}
	l.readChar()
	return tok
}

// newToken creates a token given its type and the character
func newToken(tokenType token.Type, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readString() string {
	position := l.position + 1

	for {
		l.readChar()
		if (l.ch == '"' && l.input[l.position-1] != '\\') || l.ch == 0 {
			break
		}
	}
	replacements := []struct {
		Find    string
		Replace string
	}{
		{Find: "\\\"", Replace: "\""},
		{Find: "\\n", Replace: "\n"},
		{Find: "\\t", Replace: "\t"},
		{Find: "\\r", Replace: "\r"},
	}
	var value string = l.input[position:l.position]

	for _, replacement := range replacements {
		value = strings.ReplaceAll(value, replacement.Find, replacement.Replace)
	}

	return value
}

// readIdentifier reads and returns an identifier from the input
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// isLetter returns true if the char is an alphabet or underscore
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// readNumber reads and returns a number from the input
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// isDigit returns true if the char is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// readComment reads and returns a line comment from the input
func (l *Lexer) readComment() string {
	position := l.position
	for !isEndOfLine(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// isDigit returns true if the char is a digit
func isEndOfLine(ch byte) bool {
	// ch == '\n' || ch == '\r' || ch == EOF
	return ch == 0x0a || ch == 0x0d || ch == 0
}

// skipWhitespace skips whitespace from the input
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
