package ast

import "compiler_project/lexer"

type BinOperationNode struct {
	ExpressionNode
	operator  lexer.Token
	leftNode  ExpressionNode
	rightNode ExpressionNode
}

func NewBinOperationNode(operator lexer.Token, leftNode ExpressionNode, rightNode ExpressionNode) *BinOperationNode {
	return &BinOperationNode{operator: operator, leftNode: leftNode, rightNode: rightNode}
}
