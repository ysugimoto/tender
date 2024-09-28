package ast

import "github.com/ysugimoto/tiny-template/token"

type For struct {
	*Trim
	Token     token.Token
	Interator Expression
	Arg1      Expression
	Arg2      Expression
	Block     []Node
	End       *EndFor
}

func (n *For) GetToken() token.Token { return n.Token }
func (n *For) control()              {}

type EndFor struct {
	*Trim
	Token token.Token
}

func (n *EndFor) GetToken() token.Token { return n.Token }
func (n *EndFor) control()              {}

type If struct {
	*Trim
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
	*Trim
	Token       token.Token
	Condition   Expression
	Consequence []Node
}

func (n *ElseIf) GetToken() token.Token { return n.Token }
func (n *ElseIf) control()              {}

type Else struct {
	*Trim
	Token       token.Token
	Consequence []Node
}

func (n *Else) GetToken() token.Token { return n.Token }
func (n *Else) control()              {}

type EndIf struct {
	*Trim
	Token token.Token
}

func (n *EndIf) GetToken() token.Token { return n.Token }
func (n *EndIf) control()              {}
