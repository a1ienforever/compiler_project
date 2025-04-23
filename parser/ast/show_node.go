package ast

import (
	"fmt"
)

type ShowNode struct {
	Variable ExpressionNode
}

func NewShowNode(variable ExpressionNode) *ShowNode {
	return &ShowNode{Variable: variable}
}

func (n *ShowNode) String() string {
	return fmt.Sprintf("Show(%s)", n.Variable)
}
func (n *ShowNode) isExpression() {}
