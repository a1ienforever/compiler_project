package ast

import "compiler_project/lexer"

type BinOperationNode struct {
	ExpressionNode
	Operator  lexer.Token
	LeftNode  ExpressionNode
	RightNode ExpressionNode
}

func NewBinOperationNode(operator lexer.Token, leftNode ExpressionNode, rightNode ExpressionNode) *BinOperationNode {
	return &BinOperationNode{Operator: operator, LeftNode: leftNode, RightNode: rightNode}
}
