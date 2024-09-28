package parser

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/ysugimoto/tiny-template/ast"
	"github.com/ysugimoto/tiny-template/token"
)

func (p *Parser) parseForControl(leftTrim bool) (*ast.For, error) {
	node := &ast.For{
		Trim: &ast.Trim{
			Left: leftTrim,
		},
		Token: p.curToken,
	}

	p.NextToken() // point to first loop argument

	if !p.curTokenIs(token.IDENT) {
		return nil, errors.WithStack(UnexpectedToken(p.curToken, token.IDENT))
	}
	node.Arg1 = p.parseIdent()

	// If next token is COMMA, for-loop has two arguments
	if p.peekTokenIs(token.COMMA) {
		p.NextToken() // point to COMMA
		p.NextToken() // point to second argument
		if !p.curTokenIs(token.IDENT) {
			return nil, errors.WithStack(UnexpectedToken(p.curToken, token.IDENT))
		}
		node.Arg2 = p.parseIdent()
	}

	// Expect "in" keyword
	p.NextToken()
	if !p.curTokenIs(token.IN) {
		return nil, errors.WithStack(UnexpectedToken(p.curToken, token.IN))
	}

	// Expect iterator
	p.NextToken()
	if !p.curTokenIs(token.IDENT) {
		return nil, errors.WithStack(UnexpectedToken(p.curToken, token.IDENT))
	}

	if p.peekTokenIs(token.TILDA) {
		node.Trim.Right = true
		p.NextToken()
	}

	p.NextToken() // point to inside for block

	for {
		block, err := p.parse()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if end, ok := block.(*ast.EndFor); ok {
			node.End = end
			break
		}
		node.Block = append(node.Block, block)
		p.NextToken()
	}

	return node, nil
}

func (p *Parser) parseEndForControl(leftTrim bool) (*ast.EndFor, error) {
	node := &ast.EndFor{
		Trim: &ast.Trim{
			Left: leftTrim,
		},
		Token: p.curToken,
	}

	if p.peekTokenIs(token.TILDA) {
		node.Trim.Right = true
		p.NextToken()
	}

	return node, nil
}

func (p *Parser) parseIfControl(leftTrim bool) (*ast.If, error) {
	return nil, fmt.Errorf("[IF] Not Implemented")
}

func (p *Parser) parseElseIfControl(leftTrim bool) (*ast.ElseIf, error) {
	return nil, fmt.Errorf("[ELSEIF] Not Implemented")
}

func (p *Parser) parseElseControl(leftTrim bool) (*ast.Else, error) {
	return nil, fmt.Errorf("[ELSE] Not Implemented")
}

func (p *Parser) parseEndIfControl(leftTrim bool) (*ast.EndIf, error) {
	node := &ast.EndIf{
		Trim: &ast.Trim{
			Left: leftTrim,
		},
		Token: p.curToken,
	}

	if p.peekTokenIs(token.TILDA) {
		node.Trim.Right = true
		p.NextToken()
	}

	return node, nil
}
