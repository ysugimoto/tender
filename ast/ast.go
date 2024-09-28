package ast

import "github.com/ysugimoto/tiny-template/token"

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

type Trim struct {
	Left  bool
	Right bool
}

type Literal struct {
	Token token.Token
}

func (n *Literal) GetToken() token.Token { return n.Token }

type Interporation struct {
	Token token.Token
	Value Expression
}

func (n *Interporation) GetToken() token.Token { return n.Token }
