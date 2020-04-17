package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"monkey/utils"
	"strconv"
)

// Here we use iota to give the following constants incrementing numbers as values,
// with 0 assigned to the blank identifier _.
const (
	_ int = iota
	// LOWEST precedence
	LOWEST
	// INPLACE precedence
	INPLACE
	// LOGICALOR
	LOGICALOR // ||
	// LOGICALAND
	LOGICALAND // &&
	// EQUALS just above lowest in prcecedence
	EQUALS // ==
	// LESSGREATER just above equals in prcecedence
	LESSGREATER // < <= >= >
	// SUM just above less than or greater than in prcecedence
	SUM // +
	// PRODUCT just above sum in prcecedence
	PRODUCT // *
	// POWER just above product in prcecedence
	POWER // ^
	// PREFIX just above power in prcecedence
	PREFIX // -X or !X
	// PERIOD just above prefix in prcecedence
	PERIOD // 1.2
	// CALL just above prefix in prcecedence
	CALL // myFunction(X)
	// INDEX above all others in prcecedence
	INDEX // array[index]
)

// precedence table: it associates token types with their precedence
//  1. () [] -> . ::     (function call, array index, member access)
//  2. ! ~ & ++ --       (most unary ops: not, incr, decr, ...)
//  3. * / %             (multiplication, division, modulo)
//  4. + -               (addition, subtraction)
//  5. << >>             (bitwise shift left, right)
//  6. < <= > >=         (comparisons: lt, lt eq, gt, gt eq)
//  7. == !=             (comparison: eq, not eq)
//  8. &                 (bitwise AND)
//  9. ^                 (bitwise exclusive OR, XOR)
// 10. |                 (bitwise inclusive OR)
// 11. &&                (logical AND)
// 12. ||                (logical OR)
// 13. ?:                (conditional, ternary)
// 14. = += -= *= /= %=  (ltr assignment ops)
// 14. &= |= ^= <<= >>=  (ltr assignment ops)
// 15. ,                 (comma)
var precedences = map[token.Type]int{
	token.PLUSEQ:     INPLACE,
	token.MINUSEQ:    INPLACE,
	token.SLASHEQ:    INPLACE,
	token.ASTERISKEQ: INPLACE,
	token.OR:         LOGICALOR,
	token.AND:        LOGICALAND,
	token.EQ:         EQUALS,
	token.NOTEQ:      EQUALS,
	token.LT:         LESSGREATER,
	token.GT:         LESSGREATER,
	token.LTEQ:       LESSGREATER,
	token.GTEQ:       LESSGREATER,
	token.PLUS:       SUM,
	token.MINUS:      SUM,
	token.SLASH:      PRODUCT,
	token.ASTERISK:   PRODUCT,
	token.MODULUS:    PRODUCT,
	token.POWER:      POWER,
	token.PERIOD:     PERIOD,
	token.LPAREN:     CALL,
	token.LBRACKET:   INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
	suffixParseFn func(ast.Expression) ast.Expression
)

// Parser the parser struct which contains a lexer
type Parser struct {
	l              *lexer.Lexer
	curToken       token.Token
	peekToken      token.Token
	errors         []string
	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
	suffixParseFns map[token.Type]suffixParseFn
}

// NewParser given a lexer, creates and returns a new parser
func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// register prefix parse function for all of our prefix operators
	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.DOUBLE, p.parseDoubleLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)

	// register prefic parse function for comments // ...
	p.registerPrefix(token.COMMENT, p.parseCommentLiteral)

	// register array literal parser
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	// register hash literal parser
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)

	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.INCREMENT, p.parsePrefixExpression)
	p.registerPrefix(token.DECREMENT, p.parsePrefixExpression)

	// register boolean parser
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)

	// register grouped expression parser
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	// register conditional if...else parser
	p.registerPrefix(token.IF, p.parseIfExpression)

	// register function (fn) parser
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	// register infix parse function for all of our infix operators
	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.POWER, p.parseInfixExpression)
	p.registerInfix(token.PERIOD, p.parseInfixExpression)
	p.registerInfix(token.MODULUS, p.parseInfixExpression)

	p.registerInfix(token.PLUSEQ, p.parseInfixExpression)
	p.registerInfix(token.MINUSEQ, p.parseInfixExpression)
	p.registerInfix(token.SLASHEQ, p.parseInfixExpression)
	p.registerInfix(token.ASTERISKEQ, p.parseInfixExpression)

	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOTEQ, p.parseInfixExpression)

	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)

	p.registerInfix(token.LTEQ, p.parseInfixExpression)
	p.registerInfix(token.GTEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	// register suffix parse function for all of our suffix operators
	p.suffixParseFns = make(map[token.Type]suffixParseFn)
	p.registerSuffix(token.INCREMENT, p.parseSuffixExpression)
	p.registerSuffix(token.DECREMENT, p.parseSuffixExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

// nextToken advance to the next token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// curTokenIs check that the current token type matches t
func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

// peekTokenIs check that the next token type matches t
func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

// expectPeek assert that the current token type matches before advancing
func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// curPrecedence returns the precedence associated with the token type of p.curToken
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p

	}
	return LOWEST
}

// peekPrecedence returns the precedence associated with the token type of p.peekToken
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// registerPrefix registers a prefix parser for a token type
func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// noPrefixParseFnError sets the error for a prefix expression that has no registered prefix parser
func (p *Parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// registerInfix registers an infix parser for a token type
func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// noInfixParseFnError sets the error for an infix expression that has no registered infix parser
func (p *Parser) noInfixParseFnError(t token.Type, left string, right string) {
	msg := fmt.Sprintf("no infix parse function for %s %s %s found", left, t, right)
	p.errors = append(p.errors, msg)
}

// registerSuffix registers a suffix parser for a token type
func (p *Parser) registerSuffix(tokenType token.Type, fn suffixParseFn) {
	p.suffixParseFns[tokenType] = fn
}

// noSuffixParseFnError sets the error for a suffix expression that has no registered suffix parser
func (p *Parser) noSuffixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no suffix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// Errors returns the list of errors
func (p *Parser) Errors() []string {
	return p.errors
}

// peekError creates error encountered when the type of
// next token does not match what is expected
func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// parseIdentifier parses the current token as an identifier
func (p *Parser) parseIdentifier() ast.Expression {
	identifier := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	switch p.peekToken.Type {
	case token.INCREMENT, token.DECREMENT:
		suffix := p.suffixParseFns[p.peekToken.Type]
		if suffix == nil {
			p.noSuffixParseFnError(p.peekToken.Type)
			return nil
		}
		p.nextToken()

		return suffix(identifier)
	}

	return identifier
}

// parseIntegerLiteral parses the current token as an integer literal
func (p *Parser) parseIntegerLiteral() ast.Expression {
	// defer untrace(trace("parseIntegerLiteral"))
	value, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit := &ast.IntegerLiteral{Token: p.curToken, Value: value}
	return lit
}

// parseDoubleLiteral parses the current token as a double literal
func (p *Parser) parseDoubleLiteral() ast.Expression {
	// defer untrace(trace("parseDoubleLiteral"))
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as double", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	precision := utils.Precision(p.curToken.Literal)
	lit := &ast.DoubleLiteral{Token: p.curToken, Precision: precision, Value: value}
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	// defer untrace(trace("parseStringLiteral"))
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseCommentLiteral() ast.Expression {
	// defer untrace(trace("parseCommentLiteral"))
	return &ast.CommentLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)
	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.nextToken()
		value := p.parseExpression(LOWEST)
		hash.Pairs[key] = value
		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}
	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	return hash
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return exp
}

// parseExpression parses the current token as an expression based on the registered parsers
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// defer untrace(trace("parseExpression"))
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseExpressionList(end token.Type) ast.Expressions {
	list := ast.Expressions{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}
	return list
}

// parsePrefixExpression parses the current token as a prefix expression
func (p *Parser) parsePrefixExpression() ast.Expression {
	// defer untrace(trace("parsePrefixExpression"))
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// defer untrace(trace("parseInfixExpression"))

	modOps := []interface{}{
		token.PLUSEQ,
		token.MINUSEQ,
		token.SLASHEQ,
		token.ASTERISKEQ,
	}

	if utils.InArray(p.curToken.Literal, modOps) {
		var leftType string

		switch left.(type) {
		case *ast.IntegerLiteral:
			leftType = token.INT
		case *ast.DoubleLiteral:
			leftType = token.DOUBLE
		case *ast.StringLiteral:
			leftType = token.STRING
		case *ast.Boolean:
			leftType = token.BOOL
		default:
			leftType = token.ILLEGAL
		}

		if _, ok := left.(*ast.Identifier); !ok {
			msg := fmt.Sprintf("the infix operator %s requires %s on the left, %s found", p.curToken.Literal, token.IDENT, leftType)
			p.errors = append(p.errors, msg)
			return nil
		}
	}

	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()

	// to make + right associative
	// if expression.Operator == "+" {
	// 	expression.Right = p.parseExpression(precedence - 1)
	// } else {
	// 	expression.Right = p.parseExpression(precedence)
	// }

	expression.Right = p.parseExpression(precedence)

	return expression
}

// parseSuffixExpression parses the current token as a suffix expression
func (p *Parser) parseSuffixExpression(left ast.Expression) ast.Expression {
	// defer untrace(trace("parseSuffixExpression"))
	expression := &ast.SuffixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	return expression
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	// defer untrace(trace("parseCallExpression"))
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	literal.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	literal.Body = p.parseBlockStatement()
	return literal
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = ast.Statements{}
	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)

		}
		p.nextToken()
	}
	return block
}

// parseLetStatement parses a let statement
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement parses a return statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpressionStatement parses an expression statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// parseFunctionParameters parses a function's parameters
func (p *Parser) parseFunctionParameters() ast.Identifiers {
	// defer untrace(trace("parseFunctionParameters"))
	identifiers := ast.Identifiers{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}
	p.nextToken()
	identifiers = append(identifiers, p.parseIdentifier().(*ast.Identifier))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		identifiers = append(identifiers, p.parseIdentifier().(*ast.Identifier))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// parseStatement parses a statement
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.COMMENT: // strip line comments from the program
		return nil
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// ParseProgram parses the input
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = ast.Statements{}
	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}
