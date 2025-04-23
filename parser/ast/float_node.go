package ast

import "compiler_project/lexer"

// FloatNode представляет вещественное число, например 3.14
type FloatNode struct {
	Float lexer.Token
}

func NewFloatNode(token lexer.Token) *FloatNode {
	return &FloatNode{Float: token}
}

func (*FloatNode) isExpression() {}
