package ast

import "compiler_project/lexer"

type TypedAssignNode struct {
	Type     *TypeNode
	Variable lexer.Token
	Value    ExpressionNode
}

func NewTypedAssignNode(typeToken lexer.Token, variable lexer.Token, value ExpressionNode) *TypedAssignNode {
	return &TypedAssignNode{
		Type:     NewTypeNode(typeToken.TypeToken.Name),
		Variable: variable,
		Value:    value,
	}
}
func (*TypedAssignNode) isExpression() {}
