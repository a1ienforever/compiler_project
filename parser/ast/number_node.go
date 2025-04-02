package ast

import "compiler_project/lexer"

type NumberNode struct {
	number lexer.Token
}

func NewNumberNode(number lexer.Token) *NumberNode {
	return &NumberNode{number: number}
}
