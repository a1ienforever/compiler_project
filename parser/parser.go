package parser

import (
	"compiler_project/lexer"
	"compiler_project/parser/ast"
	"fmt"
	"log"
)

type Parser struct {
	Tokens []lexer.Token
	Pos    int
	Scope  map[string]any
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{Tokens: tokens}
}

func (p *Parser) Match(expected ...lexer.TokenType) *lexer.Token {
	if p.Pos < len(p.Tokens) {
		currentToken := &p.Tokens[p.Pos]
		for _, typ := range expected {
			if typ.Name == currentToken.TypeToken.Name {
				p.Pos++
				return currentToken
			}
		}
	}
	return nil
}

func (p *Parser) Require(expected ...lexer.TokenType) *lexer.Token {
	token := p.Match(expected...)
	if token == nil {
		log.Printf("Ошибка: на позиции %d ожидается %s", p.Pos, expected[0].Name)
		panic(fmt.Sprintf("На позиции %d ожидается %s", p.Pos, expected[0].Name))
	}
	return token
}

func (p *Parser) ParseExpression() *ast.ExpressionNode {
	return nil
}

func (p *Parser) ParseCode() *ast.StatementsNode {
	root := new(ast.StatementsNode)
	for p.Pos < len(p.Tokens) {
		//codeStringNode := p.ParseExpression()
		p.Require()
		//root.AddNode(codeStringNode)
	}
	return root
}
