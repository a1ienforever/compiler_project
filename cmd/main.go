package main

import (
	"compiler_project/lexer"
	"compiler_project/parser"
	"fmt"
)

func main() {
	l := *lexer.NewLexer("123 abc")
	l.LexerAnalysis()

	p := *parser.NewParser(l.Tokens)
	fmt.Println(p)
	fmt.Println(l.Tokens)
	fmt.Print("End!")
}
