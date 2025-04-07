package ast

import "compiler_project/lexer"

type NumberNode struct {
	ExpressionNode
	Number lexer.Token
}

func NewNumberNode(number lexer.Token) *NumberNode {
	return &NumberNode{Number: number}
}
