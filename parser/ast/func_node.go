package ast

import "compiler_project/lexer"

type FunctionDeclarationNode struct {
	Name   *lexer.Token
	Params []*lexer.Token
	Body   *StatementsNode
}

func NewFunctionDeclarationNode(name *lexer.Token, params []*lexer.Token, body *StatementsNode) *FunctionDeclarationNode {
	return &FunctionDeclarationNode{Name: name, Params: params, Body: body}
}

func (n *FunctionDeclarationNode) isExpression() {}

// Вызов функции
type FunctionCallNode struct {
	Name      *lexer.Token
	Arguments []ExpressionNode
}

func NewFunctionCallNode(name *lexer.Token, arguments []ExpressionNode) *FunctionCallNode {
	return &FunctionCallNode{Name: name, Arguments: arguments}
}

func (n *FunctionCallNode) isExpression() {}
