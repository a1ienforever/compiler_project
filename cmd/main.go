package main

import (
	"compiler_project/lexer"
	"compiler_project/parser"
)

func main() {
	code := "var x = 5 + 3;"

	lexer := lexer.NewLexer(code)
	lexer.LexerAnalysis()
	parser := parser.NewParser(lexer.Tokens)
	rootNode := parser.ParseCode()
	parser.Run(rootNode)
}
