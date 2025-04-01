package parser

import (
	"compiler_project/lexer"
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

func (p *Parser) match(expected *map[string]lexer.TokenType) *lexer.Token {
	if p.Pos < len(p.Tokens) {
		currentToken := p.Tokens[p.Pos]
		fmt.Println(currentToken)
		_, ok := (*expected)[currentToken.TypeToken.Name]
		if ok {
			p.Pos++
			return &currentToken
		}
	}
	return nil
}

func (p *Parser) require(expected *map[string]lexer.TokenType) *lexer.Token {
	if token := p.match(expected); token == nil {
		log.Print(token)
		panic(fmt.Sprintf("На позиции %d ожидается "))
	}
	return nil
}
