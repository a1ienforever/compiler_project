package ast

type StatementsNode struct {
	ExpressionNode
	CodeStrings []ExpressionNode
}

func (s *StatementsNode) AddNode(node ExpressionNode) {
	s.CodeStrings = append(s.CodeStrings, node)
}
