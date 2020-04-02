package token

// Type the token type
type Type string

// Token the token
type Token struct {
	Type    Type
	Literal string
}

// Tokens list of tokens
type Tokens []Token

// LookupIdent set the identifier type based on the literal
func (t *Token) LookupIdent() {
	if tokenType, ok := keywords[t.Literal]; ok {
		t.Type = tokenType
	} else {
		t.Type = IDENT
	}
}

var keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

const (
	// ILLEGAL for unknown characters
	ILLEGAL = "ILLEGAL"
	// EOF for end of file
	EOF = "EOF"

	// IDENT for identifiers
	IDENT = "IDENT" // add, foobar, x, y, ...
	// INT for integers
	INT = "INT" // 1343456
	// BOOL for booleans
	BOOL = "BOOL" // let x: bool = false

	// ASSIGN for assignment operator
	ASSIGN = "="
	// PLUS for plus operator
	PLUS = "+"
	// MINUS for minus operator
	MINUS = "-"
	// BANG for bang operator
	BANG = "!"
	// ASTERISK for multiplication operator
	ASTERISK = "*"
	// SLASH for division operator
	SLASH = "/"
	// LT for less than operator
	LT = "<"
	// GT for greater than operator
	GT = ">"
	// EQ for greater equal to operator
	EQ = "=="
	// NOTEQ for greater not equal to operator
	NOTEQ = "!="
	// LTEQ for greater less than or equal to operator
	LTEQ = "<="
	// GTEQ for greater than or equal to operator
	GTEQ = ">="
	// PLUSEQ for plus equal to operator
	PLUSEQ = "+="
	// MINUSEQ for minus equal to operator
	MINUSEQ = "-="
	// SLASHEQ for slash equal to operator
	SLASHEQ = "/="
	// ASTERISKEQ for asterisk equal to operator
	ASTERISKEQ = "*="

	// COMMA for comma symbol
	COMMA = ","
	// SEMICOLON for semicolon symbol
	SEMICOLON = ";"
	// LPAREN for left parenthesis symbol
	LPAREN = "("
	// RPAREN for right parsnthesis symbol
	RPAREN = ")"
	// LBRACE for left bracket symbol
	LBRACE = "{"
	// RBRACE for right bracket symbol
	RBRACE = "}"

	// FUNCTION for function keyword
	FUNCTION = "FUNCTION" // func add() {}
	// LET for let keyword
	LET = "LET" // let x...
	// TRUE for boolean true
	TRUE = "TRUE" // let x = true
	// FALSE for boolean false
	FALSE = "FALSE" // let x = false
	// IF for conditional if
	IF = "IF"
	// ELSE for conditional if
	ELSE = "ELSE"
	// RETURN for return in functions
	RETURN = "RETURN"
)
