package parser

import (
	"fmt"
	"strconv"

	"git.tigh.dev/tigh-latte/monkeyscript/ast"
)

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	i, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Errorf("could not parse %q as integer: %w", p.curToken.Literal, err))
		return nil
	}
	return &ast.IntegerLiteral{Token: p.curToken, Value: i}
}
