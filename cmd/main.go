package main

import (
	"compiler_project/lexer"
	"fmt"
)

func main() {
	lexerName := *lexer.NewLexer("->")
	lexerName.LexerAnalysis()

	fmt.Print("End!")
}
