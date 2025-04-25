package main

import (
	"bufio"
	"compiler_project/lexer"
	"compiler_project/parser"
	"compiler_project/semantics"
	"compiler_project/tac"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	scope := map[string]interface{}{}
	checker := semantics.NewTypeChecker()

	for {
		fmt.Print(">>> ")
		text, _ := reader.ReadString('\n')
		if text == "exit\n" {
			break
		}
		l := lexer.NewLexer(text)
		l.LexerAnalysis()
		p := parser.NewParser(l.Tokens)
		p.Scope = scope
		rootNode := p.ParseCode()
		_, err := checker.Check(rootNode)
		if err != nil {
			fmt.Printf("Ошибка семантики: %v\n", err)
			continue
		}
		p.Run(rootNode)
		scope = p.Scope

		builder := tac.NewTACBuilder()
		builder.Generate(rootNode)
		fmt.Println("=== Трёхадресный код ===")
		builder.Print()
	}
}
