package ast

type StatementsNode struct {
	CodeStrings []ExpressionNode
}

func (s *StatementsNode) AddNode(node ExpressionNode) {
	s.CodeStrings = append(s.CodeStrings, node)
}

func (s *StatementsNode) isExpression() {}
