package ast

type StatementsNode struct {
	ExpressionNode
	codeStrings []ExpressionNode
}

func (s *StatementsNode) AddNode(node ExpressionNode) {
	s.codeStrings = append(s.codeStrings, node)
}
