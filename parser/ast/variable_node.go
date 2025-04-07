package ast

import "compiler_project/lexer"

type VariableNode struct {
	ExpressionNode
	Variable lexer.Token
}

func NewVariableNode(variable lexer.Token) *VariableNode {
	return &VariableNode{Variable: variable}
}
