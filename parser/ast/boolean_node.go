package ast

import "compiler_project/lexer"

type BooleanNode struct {
	Boolean lexer.Token
}

func NewBooleanNode(tok lexer.Token) *BooleanNode {
	return &BooleanNode{Boolean: tok}
}

func (*BooleanNode) isExpression() {}
