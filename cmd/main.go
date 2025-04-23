package main

import (
	"compiler_project/lexer"
	"compiler_project/parser"
	"fmt"
)

func main() {
	code := "int a = 10; int b = 1; show a; show b;"

	lexer := lexer.NewLexer(code)
	lexer.LexerAnalysis()

	parser := parser.NewParser(lexer.Tokens)
	rootNode := parser.ParseCode()
	fmt.Println(*rootNode)
	parser.Run(rootNode)
}
