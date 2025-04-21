package main

import (
	"compiler_project/lexer"
	"compiler_project/parser"
	"fmt"
)

func main() {
	code := "x = 10;" +
		"y = x + 5;" +
		"show y;"

	lexer := lexer.NewLexer(code)
	tokens := lexer.LexerAnalysis()

	for _, token := range *tokens {
		fmt.Printf("Token: Type=%s, Text= '%s', Pos=%d\n", token.TypeToken.Name, token.Text, token.Pos)
	}
	parser := parser.NewParser(lexer.Tokens)
	rootNode := parser.ParseCode()
	parser.Run(rootNode)
}
