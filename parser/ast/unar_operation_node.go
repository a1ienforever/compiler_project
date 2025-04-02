package ast

import "compiler_project/lexer"

type UnarOperationNode struct {
	operator lexer.Token
	operand  ExpressionNode
}
