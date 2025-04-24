package ast

type WhileNode struct {
	Condition ExpressionNode
	Body      *StatementsNode
}

func NewWhileNode(condition ExpressionNode, body *StatementsNode) *WhileNode {
	return &WhileNode{
		Condition: condition,
		Body:      body,
	}
}

func (node *WhileNode) isExpression() {}
