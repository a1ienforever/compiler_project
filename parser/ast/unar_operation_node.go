package ast

import "compiler_project/lexer"

type UnarOperationNode struct {
	ExpressionNode
	Operator lexer.Token
	Operand  ExpressionNode
}

func NewUnarOperationNode(operator lexer.Token, operand ExpressionNode) *UnarOperationNode {
	return &UnarOperationNode{Operator: operator, Operand: operand}
}
