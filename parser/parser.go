package parser

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/ysugimoto/tiny-template/ast"
	"github.com/ysugimoto/tiny-template/lexer"
	"github.com/ysugimoto/tiny-template/token"
)

const (
	LOWEST int = iota + 1
	OR
	AND
	EQUALS
	LESS_GREATER
	PREFIX
	GROUP
	END
)

var precedences = map[token.TokenType]int{
	token.EQUAL:              EQUALS,
	token.NOT_EQUAL:          EQUALS,
	token.GREATER_THAN:       LESS_GREATER,
	token.GREATER_THAN_EQUAL: LESS_GREATER,
	token.LESS_THAN:          LESS_GREATER,
	token.LESS_THAN_EQUAL:    LESS_GREATER,
	token.STRING:             PREFIX,
	token.IDENT:              PREFIX,
	token.IF:                 PREFIX,
	token.LEFT_PAREN:         GROUP,
	token.AND:                AND,
	token.OR:                 OR,
}

type (
	prefixParser  func() (ast.Expression, error)
	infixParser   func(ast.Expression) (ast.Expression, error)
	postfixParser func(ast.Expression) (ast.Expression, error)
)

type Parser struct {
	l *lexer.Lexer

	prevToken token.Token
	curToken  token.Token
	peekToken token.Token

	prefixParsers map[token.TokenType]prefixParser
	infixParsers  map[token.TokenType]infixParser
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}
	p.prefixParsers = map[token.TokenType]prefixParser{
		token.IDENT:      func() (ast.Expression, error) { return p.parseIdent(), nil },
		token.STRING:     func() (ast.Expression, error) { return p.parseString(), nil },
		token.INT:        func() (ast.Expression, error) { return p.parseInt() },
		token.FLOAT:      func() (ast.Expression, error) { return p.parseFloat() },
		token.NOT:        func() (ast.Expression, error) { return p.parsePrefixExpression() },
		token.TRUE:       func() (ast.Expression, error) { return p.parseBool(), nil },
		token.FALSE:      func() (ast.Expression, error) { return p.parseBool(), nil },
		token.LEFT_PAREN: func() (ast.Expression, error) { return p.parseGroupedExpression() },
	}
	p.infixParsers = map[token.TokenType]infixParser{
		// token.STRING:             p.parseInfixStringConcatExpression,
		// token.IDENT:              p.parseInfixStringConcatExpression,
		token.EQUAL:              p.parseInfixExpression,
		token.NOT_EQUAL:          p.parseInfixExpression,
		token.GREATER_THAN:       p.parseInfixExpression,
		token.GREATER_THAN_EQUAL: p.parseInfixExpression,
		token.LESS_THAN:          p.parseInfixExpression,
		token.LESS_THAN_EQUAL:    p.parseInfixExpression,
		token.AND:                p.parseInfixExpression,
		token.OR:                 p.parseInfixExpression,
	}

	p.NextToken()
	p.NextToken()

	return p
}

func (p *Parser) prevTokenIs(t token.TokenType) bool {
	return p.prevToken.Type == t
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) curPrecedence() int {
	if v, ok := precedences[p.curToken.Type]; ok {
		return v
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if v, ok := precedences[p.peekToken.Type]; ok {
		return v
	}
	return LOWEST
}

func (p *Parser) NextToken() {
	p.prevToken = p.curToken
	p.curToken = p.peekToken

	p.peekToken = p.l.NextToken()
}

func (p *Parser) Parse() ([]ast.Node, error) {
	var parsed []ast.Node

	for !p.curTokenIs(token.EOF) {
		node, err := p.parse()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		parsed = append(parsed, node)
	}
	return parsed, nil
}

func (p *Parser) parse() (ast.Node, error) {
	var leftTrim bool
	if p.curTokenIs(token.TILDA) {
		leftTrim = true
		p.NextToken()
	}

	switch p.curToken.Type {
	case token.LITERAL:
		return &ast.Literal{Token: p.curToken}, nil
	case token.FOR:
		return p.parseForControl(leftTrim)
	case token.ENDFOR:
		return p.parseEndForControl(leftTrim)
	case token.IF:
		return p.parseIfControl(leftTrim)
	case token.ELSEIF:
		return p.parseElseIfControl(leftTrim)
	case token.ELSE:
		return p.parseElseControl(leftTrim)
	case token.ENDIF:
		return p.parseEndIfControl(leftTrim)
	case token.INTERPORATION:
		node, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return &ast.Interporation{
			Value: node,
		}, nil
	default:
		return nil, errors.WithStack(fmt.Errorf(
			`Unexpted token "%s" found on line %d, position %d`,
			p.curToken.Type,
			p.curToken.Line,
			p.curToken.Position,
		))
	}
}
