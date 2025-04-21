package parser

import (
	"compiler_project/lexer"
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
		log.Fatalf("Ошибка: на позиции %d ожидается %s", p.Pos, expected[0].Name)
	}
	return token
}

func (p *Parser) parseVariableOrNumber() ast.ExpressionNode {
	tokenTypes := *lexer.TokenTypeList

	if number := p.Match(tokenTypes["INTEGER"]); number != nil {
		return ast.NewNumberNode(*number)
	}

	if str := p.Match(tokenTypes["STRING"]); str != nil {
		return ast.NewStringNode(*str)
	}

	if variable := p.Match(tokenTypes["VARIABLE"]); variable != nil {
		return ast.NewVariableNode(*variable)
	}

	panic(fmt.Sprintf("Ожидается переменная, строка или число на позиции %d", p.Pos))
}

func (p *Parser) parseParentheses() ast.ExpressionNode {
	tokenTypes := *lexer.TokenTypeList
	if p.Match(tokenTypes["LPAREN"]) != nil {
		node := p.parseFormula()
		p.Require(tokenTypes["RPAREN"])
		return node
	}
	return p.parseVariableOrNumber()
}

func (p *Parser) parseFormula() ast.ExpressionNode {
	tokenTypes := *lexer.TokenTypeList
	leftNode := p.parseParentheses()

	for {
		operator := p.Match(tokenTypes["MINUS"], tokenTypes["PLUS"], tokenTypes["MULTIPLY"], tokenTypes["DIVIDE"])
		if operator == nil {
			break
		}
		rightNode := p.parseParentheses()
		leftNode = ast.NewBinOperationNode(*operator, leftNode, rightNode)
	}

	return leftNode
}

func (p *Parser) ParseExpression() ast.ExpressionNode {
	tokenTypes := *lexer.TokenTypeList
	if p.Match(tokenTypes["SHOW"]) != nil {
		return p.parsePrint()
	}

	variableNode := p.parseVariableOrNumber()
	assignOperator := p.Match(tokenTypes["ASSIGN"])

	if assignOperator != nil {
		rightFormulNode := p.parseFormula()
		return ast.NewBinOperationNode(*assignOperator, variableNode, rightFormulNode)
	}

	return variableNode
}

func (p *Parser) ParseCode() *ast.StatementsNode {
	root := new(ast.StatementsNode)
	tokenTypes := *lexer.TokenTypeList

	for p.Pos < len(p.Tokens) {
		expr := p.ParseExpression()
		p.Require(tokenTypes["SEMICOLON"])
		root.AddNode(expr)
	}

	return root
}

func (p *Parser) parsePrint() ast.ExpressionNode {
	tokenTypes := *lexer.TokenTypeList
	p.Pos--
	token := p.Match(tokenTypes["SHOW"])
	if token != nil {
		return ast.NewUnarOperationNode(*token, p.parseFormula())
	}
	panic(fmt.Sprintf("Ожидается унарный оператор LOG на позиции %d", p.Pos))
}
func (p *Parser) Run(node ast.ExpressionNode) any {
	tokenTypes := *lexer.TokenTypeList

	switch n := node.(type) {
	case *ast.NumberNode:
		num, err := strconv.Atoi(n.Number.Text)
		if err != nil {
			panic(fmt.Sprintf("Невозможно преобразовать %s в число", n.Number.Text))
		}
		return num

	case *ast.StringNode:
		raw := n.String.Text
		if len(raw) >= 2 && raw[0] == '\'' && raw[len(raw)-1] == '\'' {
			return raw[1 : len(raw)-1]
		}
		return raw

	case *ast.VariableNode:
		val, ok := p.Scope[n.Variable.Text]
		if !ok {
			panic(fmt.Sprintf("Переменная '%s' не найдена", n.Variable.Text))
		}
		return val

	case *ast.UnarOperationNode:
		switch n.Operator.TypeToken.Name {
		case tokenTypes["SHOW"].Name:
			val := p.Run(n.Operand)
			fmt.Println(val)
			return nil
		}

	case *ast.BinOperationNode:
		// Обработка присваивания
		if n.Operator.TypeToken.Name == tokenTypes["ASSIGN"].Name {
			if variable, ok := n.LeftNode.(*ast.VariableNode); ok {
				value := p.Run(n.RightNode)
				p.Scope[variable.Variable.Text] = value
				fmt.Printf("Добавлена переменная %s в Scope со значением %v\n", variable.Variable.Text, value)
				return value
			}
			panic("Левая часть ASSIGN — не переменная")
		}

		left := p.Run(n.LeftNode)
		right := p.Run(n.RightNode)

		switch n.Operator.TypeToken.Name {
		case tokenTypes["PLUS"].Name:
			switch lv := left.(type) {
			case int:
				return lv + right.(int)
			case string:
				return lv + right.(string)
			default:
				panic("Операция + доступна только для строк и чисел")
			}
		case tokenTypes["MINUS"].Name:
			return left.(int) - right.(int)
		case tokenTypes["MULTIPLY"].Name:
			return left.(int) * right.(int)
		case tokenTypes["DIVIDE"].Name:
			return left.(int) / right.(int)
		}
	case *ast.StatementsNode:
		for _, stmt := range n.CodeStrings {
			p.Run(stmt)
		}
		return nil
	}

	panic("Ошибка! Неизвестный тип узла")
}
