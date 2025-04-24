package ast

type IfNode struct {
	Condition   ExpressionNode
	TrueBranch  *StatementsNode
	FalseBranch *StatementsNode
}

func NewIfNode(condition ExpressionNode, trueBranch, falseBranch *StatementsNode) *IfNode {
	return &IfNode{Condition: condition, TrueBranch: trueBranch, FalseBranch: falseBranch}
}

func (n *IfNode) isExpression() {}
