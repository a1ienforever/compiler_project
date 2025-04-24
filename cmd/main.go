package main

import (
	"bufio"
	"compiler_project/lexer"
	"compiler_project/parser"
	"fmt"
	"os"
)

func main() {
	//code := "double a = 10.5; show a;"
	reader := bufio.NewReader(os.Stdin)
	scope := map[string]interface{}{}

	for {
		fmt.Print(">>> ")
		text, _ := reader.ReadString('\n')
		if text == "exit\n" {
			break
		}
		l := lexer.NewLexer(text)
		l.LexerAnalysis()
		fmt.Println(l.Tokens)
		p := parser.NewParser(l.Tokens)
		p.Scope = scope
		rootNode := p.ParseCode()
		p.Run(rootNode)
		scope = p.Scope
	}
	//lexer := lexer.NewLexer(code)
	//lexer.LexerAnalysis()
	//
	//rootNode := parser.ParseCode()
	//parser.Run(rootNode)
}
