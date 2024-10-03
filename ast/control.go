package ast

import "github.com/ysugimoto/tender/token"

type For struct {
	Token    token.Token
	Iterator *Ident
	Arg1     *Ident
	Arg2     *Ident
	Block    []Node
	End      *EndFor
}

func (n *For) GetToken() token.Token { return n.Token }
func (n *For) control()              {}

type EndFor struct {
	Token token.Token
}

func (n *EndFor) GetToken() token.Token { return n.Token }
func (n *EndFor) control()              {}

type If struct {
	Token       token.Token
	Condition   Expression
	Another     []*ElseIf
	Consequence []Node
	Alternative *Else
	End         *EndIf
}

func (n *If) GetToken() token.Token { return n.Token }
func (n *If) control()              {}

type ElseIf struct {
	Token       token.Token
	Condition   Expression
	Consequence []Node
}

func (n *ElseIf) GetToken() token.Token { return n.Token }
func (n *ElseIf) control()              {}

type Else struct {
	Token       token.Token
	Consequence []Node
}

func (n *Else) GetToken() token.Token { return n.Token }
func (n *Else) control()              {}

type EndIf struct {
	Token token.Token
}

func (n *EndIf) GetToken() token.Token { return n.Token }
func (n *EndIf) control()              {}
