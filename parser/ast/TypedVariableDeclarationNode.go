package ast

// TypeNode представляет узел для типа данных (например, int, string, float)
type TypeNode struct {
	Type string // Тип данных, например "int", "string", "float"
}

// NewTypeNode создает новый узел для типа данных
func NewTypeNode(typ string) *TypeNode {
	return &TypeNode{
		Type: typ,
	}
}
