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

// keywords is the list of keywords of the programming language
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
	IDENT = "IDENTIFIER" // add, foobar, x, y, ...
	// INT for integers
	INT = "INTEGER" // 1343456
	// DOUBLE for integers
	DOUBLE = "DOUBLE" // 13434.56
	// BOOL for booleans
	BOOL = "BOOLEAN" // let x = false
	// STRING for strings
	STRING = "STRING" // let x = "hello, world"

	// ASSIGN for assignment operator
	ASSIGN = "="

	// PLUS for plus operator
	PLUS = "+"
	// MINUS for minus operator
	MINUS = "-"
	// ASTERISK for multiplication operator
	ASTERISK = "*"
	// SLASH for division operator
	SLASH = "/"
	// POWER for power operator
	POWER = "^"
	// MODULUS for modulus operator
	MODULUS = "%"

	// BANG for bang (not) operator
	BANG = "!"
	// INCREMENT for increment operator
	INCREMENT = "++"
	// DECREMENT for decrement operator
	DECREMENT = "--"

	// EQ for greater equal to operator
	EQ = "=="
	// NOTEQ for greater not equal to operator
	NOTEQ = "!="

	// LT for less than operator
	LT = "<"
	// GT for greater than operator
	GT = ">"
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

	// PERIOD for period symbol
	PERIOD = "."
	// COMMA for comma symbol
	COMMA = ","
	// SEMICOLON for semicolon symbol
	SEMICOLON = ";"
	// COLON for colon symbol
	COLON = ":"
	// LPAREN for left parenthesis symbol
	LPAREN = "("
	// RPAREN for right parsnthesis symbol
	RPAREN = ")"
	// LBRACE for left bracket symbol
	LBRACE = "{"
	// RBRACE for right bracket symbol
	RBRACE = "}"
	// LBRACKET for right square bracket symbol
	LBRACKET = "["
	// RBRACKET for right square bracket symbol
	RBRACKET = "]"

	// COMMENT for comments
	COMMENT = "//"

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
