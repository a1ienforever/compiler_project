package ast

import "compiler_project/lexer"

type VariableNode struct {
	ExpressionNode
	variable lexer.Token
}

func NewVariableNode(variable lexer.Token) *VariableNode {
	return &VariableNode{variable: variable}
}
