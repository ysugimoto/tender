package parser

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/ysugimoto/tiny-template/ast"
	"github.com/ysugimoto/tiny-template/token"
)

func (p *Parser) parseIdent() *ast.Ident {
	return &ast.Ident{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseString() *ast.String {
	return &ast.String{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseInt() (*ast.Int, error) {
	v, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		return nil, errors.WithStack(TypeConversionError(p.curToken, "INTEGER"))
	}

	return &ast.Int{
		Token: p.curToken,
		Value: v,
	}, nil
}

func (p *Parser) parseFloat() (*ast.Float, error) {
	v, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		return nil, errors.WithStack(TypeConversionError(p.curToken, "FLOAT"))
	}

	return &ast.Float{
		Token: p.curToken,
		Value: v,
	}, nil
}

func (p *Parser) parseBool() *ast.Bool {
	return &ast.Bool{
		Token: p.curToken,
		Value: p.curToken.Type == token.TRUE,
	}
}
