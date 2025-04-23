package ast

type TypeNode struct {
	Type string // "int", "string", "double"
}

func NewTypeNode(typ string) *TypeNode {
	return &TypeNode{Type: typ}
}

func (*TypeNode) isExpression() {}
