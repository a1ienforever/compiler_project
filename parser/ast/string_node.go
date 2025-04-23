package ast

import "compiler_project/lexer"

type StringNode struct {
	String lexer.Token
}

func NewStringNode(token lexer.Token) *StringNode {
	return &StringNode{String: token}
}
func (*StringNode) isExpression() {}
