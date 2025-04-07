package parser

import (
	lexer "compiler_project/lexer"
	"compiler_project/parser/ast"
	"fmt"
	"log"
	"strconv"
)

type Parser struct {
	Tokens []lexer.Token
	Pos    int
	Scope  map[string]any
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{Tokens: tokens, Scope: make(map[string]any)}
}
func (p *Parser) Match(expected ...lexer.TokenType) *lexer.Token {
	if p.Pos < len(p.Tokens) {
		currentToken := &p.Tokens[p.Pos]
		for _, typ := range expected {
			if typ.Name == currentToken.TypeToken.Name {
				p.Pos++
				return currentToken
			}
		}
	}
	return nil
}

func (p *Parser) Require(expected ...lexer.TokenType) *lexer.Token {
	token := p.Match(expected...)
	if token == nil {
		log.Printf("Ошибка: на позиции %d ожидается %s", p.Pos, expected[0].Name)
		panic(fmt.Sprintf("На позиции %d ожидается %s", p.Pos, expected[0].Name))
	}
	return token
}

func (p *Parser) parseVariableOrNumber() ast.ExpressionNode {
	number := p.Match((*lexer.TokenTypeList)["NUMBER"])
	if number != nil {
		return ast.NewNumberNode(*number)
	}

	variable := p.Match((*lexer.TokenTypeList)["VARIABLE"])
	if variable != nil {
		return ast.NewVariableNode(*variable)
	}
	panic(fmt.Sprintf("Ожидается переменная или число на %d", p.Pos))
	return nil
}

func (p *Parser) parseParentheses() ast.ExpressionNode {
	if p.Match((*lexer.TokenTypeList)["LPAREN"]) != nil {
		node := p.parseFormula()
		p.Require((*lexer.TokenTypeList)["RPAREN"])
		return node
	}
	return p.parseVariableOrNumber()
}

func (p *Parser) parseFormula() ast.ExpressionNode {
	leftNode := p.parseParentheses()
	operator := p.Match((*lexer.TokenTypeList)["MINUS"], (*lexer.TokenTypeList)["PLUS"])

	for operator != nil {
		rightNode := p.parseParentheses()
		leftNode = ast.NewBinOperationNode(*operator, leftNode, rightNode)
		operator = p.Match((*lexer.TokenTypeList)["MINUS"], (*lexer.TokenTypeList)["PLUS"])
	}

	return leftNode
}

func (p *Parser) ParseExpression() ast.ExpressionNode {
	if p.Match((*lexer.TokenTypeList)["VARIABLE"]) == nil {
		printNode := p.parsePrint()
		return printNode
	}
	p.Pos--
	variableNode := p.parseVariableOrNumber()
	assignOperator := p.Match((*lexer.TokenTypeList)["ASSIGN"])
	if assignOperator != nil {
		rightFormulNode := p.parseFormula()
		binaryNode := ast.NewBinOperationNode(*assignOperator, variableNode, rightFormulNode)
		return binaryNode
	}
	panic(fmt.Sprintf("После переменной ожидается оператор присвоения на позиции %d", p.Pos))
	return nil
}

func (p *Parser) ParseCode() *ast.StatementsNode {
	root := new(ast.StatementsNode)
	for p.Pos < len(p.Tokens) {
		codeStringNode := p.ParseExpression()
		token, _ := (*lexer.TokenTypeList)["SEMICOLON"]
		p.Require(token)
		root.AddNode(codeStringNode)
	}
	return root
}

func (p *Parser) parsePrint() ast.ExpressionNode {
	token := p.Match((*lexer.TokenTypeList)["LOG"])
	if token != nil {
		return ast.NewUnarOperationNode(*token, p.parseFormula())
	}
	panic(fmt.Sprintf("Ожидается унарный оператор LOG на позиции %d", p.Pos))
}

func (p *Parser) Run(node ast.ExpressionNode) any {
	switch n := node.(type) {

	case *ast.NumberNode:
		num, err := strconv.Atoi(n.Number.Text)
		if err != nil {
			panic(fmt.Sprintf("Невозможно преобразовать %s в число", n.Number.Text))
		}
		return num

	case *ast.UnarOperationNode:
		switch n.Operator.TypeToken.Name {
		case (*lexer.TokenTypeList)["LOG"].Name:
			fmt.Println(p.Run(n.Operand))
			return nil
		}

	case *ast.BinOperationNode:
		switch n.Operator.TypeToken.Name {
		case (*lexer.TokenTypeList)["PLUS"].Name:
			return p.Run(n.LeftNode).(int) + p.Run(n.RightNode).(int)

		case (*lexer.TokenTypeList)["MINUS"].Name:
			return p.Run(n.LeftNode).(int) - p.Run(n.RightNode).(int)

		case (*lexer.TokenTypeList)["ASSIGN"].Name:
			result := p.Run(n.RightNode)
			if variable, ok := n.LeftNode.(*ast.VariableNode); ok {
				p.Scope[variable.Variable.Text] = result
				return result
			} else {
				panic("Левая часть ASSIGN — не VariableNode")
			}
		}

	case *ast.VariableNode:
		val, ok := p.Scope[n.Variable.Text]
		if !ok {
			panic(fmt.Sprintf("Переменная '%s' не найдена", n.Variable.Text))
		}
		return val

	case *ast.StatementsNode:
		for _, stmt := range n.CodeStrings {
			p.Run(stmt)
		}
		return nil
	}
	panic("Ошибка! Неизвестный тип узла")
}
