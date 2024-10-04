package parser

import (
	"github.com/pkg/errors"
	"github.com/ysugimoto/tender/ast"
	"github.com/ysugimoto/tender/token"
)

func (p *Parser) parseControl(cs controlState) (ast.Control, error) {
	leftTrim := p.curToken.LeftTrim

	// point to control keyword and copy trimming flag
	p.NextToken()
	p.curToken.LeftTrim = leftTrim

	parsers, ok := p.controlParsers[cs]
	if !ok {
		return nil, errors.WithStack(UndefinedControlParserState(p.curToken))
	}

	parser, ok := parsers[p.curToken.Type]
	if !ok {
		return nil, errors.WithStack(UnexpectedToken(p.curToken))
	}

	return parser()
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
	node.Iterator = p.parseIdent()

	if !p.peekTokenIs(token.CONTROL_END) {
		return nil, errors.WithStack(UnexpectedToken(p.curToken, token.CONTROL_END))
	}
	p.NextToken()
	node.Token.RightTrim = p.curToken.RightTrim

	p.NextToken() // point to inside of control

	pool := nodePool.Get().(*[]ast.Node) // nolint:errcheck
	blocks := *pool
	defer func() {
		*pool = blocks
		nodePool.Put(pool)
	}()

	blocks = blocks[0:0]
	for {
		switch p.curToken.Type {
		case token.LITERAL:
			blocks = append(blocks, &ast.Literal{
				Token: p.curToken,
			})
		case token.CONTROL_START:
			control, err := p.parseControl(FOR)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if end, ok := control.(*ast.EndFor); ok {
				node.End = end
				goto OUT
			}
			blocks = append(blocks, control)
		case token.INTERPORATION:
			blocks = append(blocks, &ast.Interporation{
				Token: p.curToken,
				Value: p.parseIdent(),
			})
		default:
			return nil, errors.WithStack(UnexpectedToken(p.curToken))
		}
		p.NextToken()
	}
OUT:

	node.Block = make([]ast.Node, len(blocks))
	copy(node.Block, blocks)
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
		Token:       p.curToken,
		Another:     []*ast.ElseIf{},
		Consequence: []ast.Node{},
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

	// Acceptable control statement depends on the parser state
	// So change state on the following for-loop inside.
	//
	// The acceptable controls spec are:
	// IF:     for, if, elseif, else, endif
	// ELSEIF: for, if, elseif, else, endif
	// ELSE:   for, if, endif
	state := IF

	for {
		switch p.curToken.Type {
		case token.LITERAL:
			appendTarget(&ast.Literal{
				Token: p.curToken,
			})
		case token.CONTROL_START:
			control, err := p.parseControl(state)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			switch t := control.(type) {
			case *ast.ElseIf:
				node.Another = append(node.Another, t)
			case *ast.Else:
				node.Alternative = t
				// move state to ELSE
				state = ELSE
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
		p.NextToken()
	}
OUT:

	return node, nil
}

func (p *Parser) parseElseIfControl() (*ast.ElseIf, error) {
	node := &ast.ElseIf{
		Token:       p.curToken,
		Consequence: []ast.Node{},
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

func (p *Parser) parseSeparatedElseIfControl() (*ast.ElseIf, error) {
	node := &ast.ElseIf{
		Token:       p.curToken,
		Consequence: []ast.Node{},
	}

	p.NextToken() // point to separated if token
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

func (p *Parser) parseElseControl() (ast.Control, error) {
	// else if also treats as elseif
	if p.peekTokenIs(token.IF) {
		p.curToken.Literal += " if"
		return p.parseSeparatedElseIfControl()
	}

	node := &ast.Else{
		Token:       p.curToken,
		Consequence: []ast.Node{},
	}

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

	if !p.peekTokenIs(token.CONTROL_END) {
		return nil, errors.WithStack(UnexpectedToken(p.curToken, token.CONTROL_END))
	}
	p.NextToken() // point to CONTROL_END
	node.Token.RightTrim = p.curToken.RightTrim

	return node, nil
}
