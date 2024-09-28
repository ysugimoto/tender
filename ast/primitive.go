package ast

import "github.com/ysugimoto/tiny-template/token"

type Ident struct {
	Token token.Token
	Value string
}

func (n *Ident) GetToken() token.Token { return n.Token }
func (n *Ident) expression()           {}

type String struct {
	Token token.Token
	Value string
}

func (n *String) GetToken() token.Token { return n.Token }
func (n *String) expression()           {}

type Int struct {
	Token token.Token
	Value int64
}

func (n *Int) GetToken() token.Token { return n.Token }
func (n *Int) expression()           {}

type Float struct {
	Token token.Token
	Value float64
}

func (n *Float) GetToken() token.Token { return n.Token }
func (n *Float) expression()           {}

type Bool struct {
	Token token.Token
	Value bool
}

func (n *Bool) GetToken() token.Token { return n.Token }
func (n *Bool) expression()           {}
