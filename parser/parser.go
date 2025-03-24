package parser

import "compiler_project/lexer"

type Parser struct {
	Tokens []lexer.Token
	Pos    int
	Scope  map[string]any
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{Tokens: tokens}
}

func (p *Parser) match() *lexer.Token {
	return nil
}

func (p *Parser) require() *lexer.Token {
	return nil
}
