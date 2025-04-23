package ast

import "compiler_project/lexer"

type VariableNode struct {
	Variable lexer.Token
}

func NewVariableNode(token lexer.Token) *VariableNode {
	return &VariableNode{Variable: token}
}
func (*VariableNode) isExpression() {}
