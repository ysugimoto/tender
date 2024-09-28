package ast

import "github.com/ysugimoto/tiny-template/token"

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (n *PrefixExpression) GetToken() token.Token { return n.Token }
func (n *PrefixExpression) expression()           {}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (n *InfixExpression) GetToken() token.Token { return n.Token }
func (n *InfixExpression) expression()           {}

type GroupedExpression struct {
	Token token.Token
	Right Expression
}

func (n *GroupedExpression) GetToken() token.Token { return n.Token }
func (n *GroupedExpression) expression()           {}
