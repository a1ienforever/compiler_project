package lexer

import "fmt"

type Lexer struct {
	code   string
	pos    int
	tokens []Token
}

func NewLexer(code string) *Lexer {
	return &Lexer{code: code}
}

func (l *Lexer) LexerAnalysis() *[]Token {
	fmt.Print(*l)
	return &l.tokens
}
