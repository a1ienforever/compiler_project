package ast

import "compiler_project/lexer"

type NumberNode struct {
	Number lexer.Token
}

func NewNumberNode(token lexer.Token) *NumberNode {
	return &NumberNode{Number: token}
}
func (*NumberNode) isExpression() {}
