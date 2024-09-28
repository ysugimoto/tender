package parser

import (
	"github.com/pkg/errors"
	"github.com/ysugimoto/tiny-template/ast"
	"github.com/ysugimoto/tiny-template/token"
)

func (p *Parser) parseExpression(precedence int) (ast.Expression, error) {
	prefix, ok := p.prefixParsers[p.curToken.Type]
	if !ok {
		return nil, UndefinedPrefix(p.curToken)
	}

	left, err := prefix()
	if err != nil {
		return nil, err
	}

	for precedence < p.peekPrecedence() {
		infix, ok := p.infixParsers[p.peekToken.Type]
		if !ok {
			return left, nil
		}
		p.NextToken()
		left, err = infix(left)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return left, nil
}

func (p *Parser) parsePrefixExpression() (*ast.PrefixExpression, error) {
	node := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.NextToken() // point to expression start
	right, err := p.parseExpression(PREFIX)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	node.Right = right

	return node, nil
}

func (p *Parser) parseGroupedExpression() (*ast.GroupedExpression, error) {
	node := &ast.GroupedExpression{
		Token: p.curToken,
	}

	p.NextToken() // point to expression start
	right, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	node.Right = right

	if !p.peekTokenIs(token.RIGHT_PAREN) {
		return nil, errors.WithStack(UnexpectedToken(p.peekToken, "RIGHT_PAREN"))
	}
	p.NextToken() // point to RIGHT_PAREN

	return node, nil
}

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, error) {
	node := &ast.InfixExpression{
		Token:    p.curToken, // point to operator token
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.NextToken() // point to right expression start
	right, err := p.parseExpression(precedence)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	node.Right = right

	return node, nil
}
