package ast

import "compiler_project/lexer"

type StringNode struct {
	ExpressionNode
	String lexer.Token
}

func NewStringNode(str lexer.Token) *StringNode {
	return &StringNode{String: str}
}
