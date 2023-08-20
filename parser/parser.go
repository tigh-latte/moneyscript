package parser

import (
	"errors"
	"fmt"

	"git.tigh.dev/tigh-latte/monkeyscript/ast"
	"git.tigh.dev/tigh-latte/monkeyscript/lexer"
	"git.tigh.dev/tigh-latte/monkeyscript/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // iX or !X
	CALL        // myFunction(X)
)

type (
	prefixParseFunc func() ast.Expression
	infixParseFunc  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	errors []error

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFunc
	infixParseFns  map[token.TokenType]infixParseFunc
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: make([]error, 0),
	}

	p.prefixParseFns = map[token.TokenType]prefixParseFunc{
		token.IDENT:   p.parseIdentifier,
		token.INT:     p.parseIntegerLiteral,
		token.MINUS:   p.parsePrefixExpression,
		token.EXCLAIM: p.parsePrefixExpression,
	}

	// Call twice to set both curToken and nextToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

//func (p *Parser) registerPrefix(tt token.TokenType, fn prefixParseFunc) {
//	p.prefixParseFns[tt] = fn
//}

func (p *Parser) registerInfix(tt token.TokenType, fn infixParseFunc) {
	p.infixParseFns[tt] = fn
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		// if stmt != nil {
		program.Statements = append(program.Statements, stmt)
		//}
		p.nextToken()
	}

	return program
}

func (p *Parser) expectPeek(token token.TokenType) bool {
	if p.peekToken.Type != token {
		err := fmt.Errorf("%w: expected %q got %q", ErrUnexpectedToken, token, p.peekToken.Type)
		p.errors = append(p.errors, err)
		return false
	}

	return true
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatemeant()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	if ok := p.expectPeek(token.IDENT); !ok {
		return nil
	}
	stmt := &ast.LetStatement{Token: p.curToken}

	p.nextToken() // Advance to identifier

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if ok := p.expectPeek(token.ASSIGN); !ok {
		return nil
	}

	p.nextToken() // advance to assigment

	// TODO: skipping expressions until we get a semicolon
	for p.curToken.Type != token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatemeant() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	// TODO: skip expressions until we hit semi colon
	for p.curToken.Type != token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	fn := p.prefixParseFns[p.curToken.Type]
	if fn == nil {
		p.errors = append(p.errors, fmt.Errorf("no prefix parse function for %s found", p.curToken.Type))
		return nil
	}

	return fn()
}

func (p *Parser) Errors() error {
	return errors.Join(p.errors...)
}
