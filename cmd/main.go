package main

import (
	"compiler_project/lexer"
	"fmt"
)

func main() {
	lexerName := *lexer.NewLexer("var")
	lexerName.LexerAnalysis()

	fmt.Print()
}
