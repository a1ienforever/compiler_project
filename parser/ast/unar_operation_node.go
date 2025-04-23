package ast

import "compiler_project/lexer"

type UnarOperationNode struct {
	Operator lexer.Token
	Operand  ExpressionNode
}

func NewUnarOperationNode(op lexer.Token, operand ExpressionNode) *UnarOperationNode {
	return &UnarOperationNode{
		Operator: op,
		Operand:  operand,
	}
}
func (*UnarOperationNode) isExpression() {}
