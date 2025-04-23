package ast

import "compiler_project/lexer"

type BinOperationNode struct {
	Operator  lexer.Token
	LeftNode  ExpressionNode
	RightNode ExpressionNode
}

func NewBinOperationNode(op lexer.Token, left, right ExpressionNode) *BinOperationNode {
	return &BinOperationNode{
		Operator:  op,
		LeftNode:  left,
		RightNode: right,
	}
}
func (*BinOperationNode) isExpression() {}
