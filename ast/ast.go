package ast

import "github.com/ysugimoto/tender/token"

type Node interface {
	GetToken() token.Token
}

type Control interface {
	Node
	control()
}

type Expression interface {
	Node
	expression()
}

type Literal struct {
	Token token.Token
	Value string
}

func (n *Literal) GetToken() token.Token { return n.Token }

type Interporation struct {
	Token token.Token
	Value *Ident
}

func (n *Interporation) GetToken() token.Token { return n.Token }
