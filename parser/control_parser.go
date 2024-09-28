package parser

import (
	"github.com/pkg/errors"
	"github.com/ysugimoto/tiny-template/ast"
	"github.com/ysugimoto/tiny-template/token"
)

func (p *Parser) parseControl(isRoot bool) (ast.Control, error) {
	leftTrim := p.curToken.LeftTrim

	// point to control keyword and copy trimming flag
	p.NextToken()
	p.curToken.LeftTrim = leftTrim

	switch p.curToken.Type {
	case token.FOR:
		return p.parseForControl()
	case token.IF:
		return p.parseIfControl()

	// Following control is forbidden to present on root
	case token.ENDFOR:
		if isRoot {
			return nil, errors.WithStack(UnexpectedToken(p.curToken))
		}
		return p.parseEndForControl()
	case token.ELSEIF:
		if isRoot {
			return nil, errors.WithStack(UnexpectedToken(p.curToken))
		}
		return p.parseElseIfControl()
	case token.ELSE:
		if isRoot {
			return nil, errors.WithStack(UnexpectedToken(p.curToken))
		}
		return p.parseElseControl()
	case token.ENDIF:
		if isRoot {
			return nil, errors.WithStack(UnexpectedToken(p.curToken))
		}
		return p.parseEndIfControl()
	default:
		return nil, errors.WithStack(UnexpectedToken(p.curToken))
	}
}

func (p *Parser) parseForControl() (*ast.For, error) {
	node := &ast.For{
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
	node.Interator = p.parseIdent()

	if !p.peekTokenIs(token.CONTROL_END) {
		return nil, errors.WithStack(UnexpectedToken(p.curToken, token.CONTROL_END))
	}
	p.NextToken()
	node.Token.RightTrim = p.curToken.RightTrim

	p.NextToken() // point to inside of control

	for {
		switch p.curToken.Type {
		case token.LITERAL:
			node.Block = append(node.Block, &ast.Literal{
				Token: p.curToken,
				Value: p.curToken.Literal,
			})
		case token.CONTROL_START:
			control, err := p.parseControl(false)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if end, ok := control.(*ast.EndFor); ok {
				node.End = end
				goto OUT
			}
			node.Block = append(node.Block, control)
		case token.INTERPORATION:
			node.Block = append(node.Block, &ast.Interporation{
				Token: p.curToken,
				Value: p.parseIdent(),
			})
		default:
			return nil, errors.WithStack(UnexpectedToken(p.curToken))
		}
		p.NextToken()
	}
OUT:

	return node, nil
}

func (p *Parser) parseEndForControl() (*ast.EndFor, error) {
	node := &ast.EndFor{
		Token: p.curToken,
	}

	if !p.peekTokenIs(token.CONTROL_END) {
		return nil, errors.WithStack(UnexpectedToken(p.curToken, token.CONTROL_END))
	}
	p.NextToken() // point to CONTROL_END
	node.Token.RightTrim = p.curToken.RightTrim

	return node, nil
}

func (p *Parser) parseIfControl() (*ast.If, error) {
	node := &ast.If{
		Token: p.curToken,
	}

	p.NextToken() // point to first condition token

	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	node.Condition = exp

	if !p.peekTokenIs(token.CONTROL_END) {
		return nil, errors.WithStack(UnexpectedToken(p.curToken, token.CONTROL_END))
	}
	p.NextToken() // point to CONTROL_END
	node.Token.RightTrim = p.curToken.RightTrim

	p.NextToken() // point to inside of control

	appendTarget := func(n ast.Node) {
		switch {
		case node.Alternative != nil:
			node.Alternative.Consequence = append(node.Alternative.Consequence, n)
		case len(node.Another) > 0:
			node.Another[len(node.Another)-1].Consequence = append(node.Another[len(node.Another)-1].Consequence, n)
		default:
			node.Consequence = append(node.Consequence, n)
		}
	}

	for {
		switch p.curToken.Type {
		case token.LITERAL:
			appendTarget(&ast.Literal{
				Token: p.curToken,
				Value: p.curToken.Literal,
			})
		case token.CONTROL_START:
			control, err := p.parseControl(false)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			switch t := control.(type) {
			case *ast.ElseIf:
				node.Another = append(node.Another, t)
			case *ast.Else:
				node.Alternative = t
			case *ast.EndIf:
				node.End = t
				goto OUT
			default:
				appendTarget(&ast.Literal{Token: p.curToken})
			}
		case token.INTERPORATION:
			appendTarget(&ast.Interporation{
				Token: p.curToken,
				Value: p.parseIdent(),
			})
		default:
			return nil, errors.WithStack(UnexpectedToken(p.curToken))
		}
	}
OUT:

	return node, nil
}

func (p *Parser) parseElseIfControl() (*ast.ElseIf, error) {
	node := &ast.ElseIf{
		Token: p.curToken,
	}

	p.NextToken() // point to first condition token

	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	node.Condition = exp

	if !p.peekTokenIs(token.CONTROL_END) {
		return nil, errors.WithStack(UnexpectedToken(p.curToken, token.CONTROL_END))
	}
	p.NextToken() // point to CONTROL_END
	node.Token.RightTrim = p.curToken.RightTrim

	return node, nil
}

func (p *Parser) parseElseControl() (*ast.Else, error) {
	node := &ast.Else{
		Token: p.curToken,
	}

	p.NextToken() // point to first condition token

	if !p.peekTokenIs(token.CONTROL_END) {
		return nil, errors.WithStack(UnexpectedToken(p.curToken, token.CONTROL_END))
	}
	p.NextToken() // point to CONTROL_END
	node.Token.RightTrim = p.curToken.RightTrim

	return node, nil
}

func (p *Parser) parseEndIfControl() (*ast.EndIf, error) {
	node := &ast.EndIf{
		Token: p.curToken,
	}

	p.NextToken() // point to first condition token

	if !p.peekTokenIs(token.CONTROL_END) {
		return nil, errors.WithStack(UnexpectedToken(p.curToken, token.CONTROL_END))
	}
	p.NextToken() // point to CONTROL_END
	node.Token.RightTrim = p.curToken.RightTrim

	return node, nil
}
