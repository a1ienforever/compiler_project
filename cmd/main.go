package main

import (
	"compiler_project/lexer"
	"compiler_project/llvmgen"
	"compiler_project/parser"
	"compiler_project/semantics"
	"compiler_project/tac"
	"fmt"
	"os"
)

func main() {
	scope := map[string]interface{}{}
	checker := semantics.NewTypeChecker()

	text := `int a=1;
while a less 10{
	show a;
	int a = a + 1;
	};`

	l := lexer.NewLexer(text)
	l.LexerAnalysis()
	fmt.Println(l.Tokens)

	p := parser.NewParser(l.Tokens)
	p.Scope = scope
	rootNode := p.ParseCode()

	_, err := checker.Check(rootNode)
	if err != nil {
		fmt.Printf("Ошибка семантики: %v\n", err)
		return
	}
	p.Run(rootNode)
	scope = p.Scope

	fmt.Println("=== Трёхадресный код ===")
	builder := tac.NewTACBuilder()
	builder.Generate(rootNode)
	builder.Optimize()
	builder.Print()

	llvm := llvmgen.NewLLVMBuilder()
	llvm.GenerateFromTAC(builder.Instructions())
	irop := llvm.IR()

	outFile, _ := os.Create("output.ll")
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {
			panic(err)
		}
	}(outFile)
	_, err = outFile.WriteString(irop.String())
	if err != nil {
		panic(err)
	}

}
