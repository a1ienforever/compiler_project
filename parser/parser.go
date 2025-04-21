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

// Match tries to match the current token with the expected types and advances the position.
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

// Require ensures the next token matches one of the expected types, otherwise it panics.
func (p *Parser) Require(expected ...lexer.TokenType) *lexer.Token {
	token := p.Match(expected...)
	if token == nil {
		log.Printf("Ошибка: на позиции %d ожидается %s", p.Pos, expected[0].Name)
		panic(fmt.Sprintf("На позиции %d ожидается %s", p.Pos, expected[0].Name))
	}
	return token
}

// parseVariableOrNumber parses either a number or a variable.
func (p *Parser) parseVariableOrNumber() ast.ExpressionNode {
	tokenTypes := *lexer.TokenTypeList
	number := p.Match(tokenTypes["INTEGER"])
	if number != nil {
		return ast.NewNumberNode(*number)
	}

	variable := p.Match(tokenTypes["VARIABLE"])
	if variable != nil {
		return ast.NewVariableNode(*variable)
	}
	panic(fmt.Sprintf("Ожидается переменная или число на %d", p.Pos))
}

// parseParentheses parses expressions within parentheses.
func (p *Parser) parseParentheses() ast.ExpressionNode {
	tokenTypes := *lexer.TokenTypeList
	if p.Match(tokenTypes["LPAREN"]) != nil {
		node := p.parseFormula()
		p.Require(tokenTypes["RPAREN"])
		return node
	}
	return p.parseVariableOrNumber()
}

// parseFormula parses a formula, supporting binary operations like +, -, *, and /.
func (p *Parser) parseFormula() ast.ExpressionNode {
	tokenTypes := *lexer.TokenTypeList
	leftNode := p.parseParentheses()

	operator := p.Match(tokenTypes["MINUS"], tokenTypes["PLUS"], tokenTypes["MULTIPLY"], tokenTypes["DIVIDE"])

	for operator != nil {
		rightNode := p.parseParentheses()
		leftNode = ast.NewBinOperationNode(*operator, leftNode, rightNode)
		operator = p.Match(tokenTypes["MINUS"], tokenTypes["PLUS"], tokenTypes["MULTIPLY"], tokenTypes["DIVIDE"])
	}

	return leftNode
}

// ParseExpression handles the parsing of expressions, including assignments.
func (p *Parser) ParseExpression() ast.ExpressionNode {
	// Handle print operations (e.g. `show`)
	tokenTypes := *lexer.TokenTypeList
	if p.Match(tokenTypes["SHOW"]) != nil {
		return p.parsePrint()
	}

	variableNode := p.parseVariableOrNumber()
	assignOperator := p.Match(tokenTypes["ASSIGN"])

	if assignOperator != nil {
		rightFormulNode := p.parseFormula()
		binaryNode := ast.NewBinOperationNode(*assignOperator, variableNode, rightFormulNode)
		return binaryNode
	}

	return variableNode
}

// ParseCode parses the entire program code consisting of multiple statements.
func (p *Parser) ParseCode() *ast.StatementsNode {
	root := new(ast.StatementsNode)
	tokenTypes := *lexer.TokenTypeList
	for p.Pos < len(p.Tokens) {
		codeStringNode := p.ParseExpression()
		p.Require(tokenTypes["SEMICOLON"])
		root.AddNode(codeStringNode)
	}
	return root
}

// parsePrint parses print operations (e.g. `show <expression>`).
func (p *Parser) parsePrint() ast.ExpressionNode {
	tokenTypes := *lexer.TokenTypeList
	p.Pos--
	token := p.Match(tokenTypes["SHOW"])
	if token != nil {
		return ast.NewUnarOperationNode(*token, p.parseFormula())
	}
	panic(fmt.Sprintf("Ожидается унарный оператор LOG на позиции %d", p.Pos))
}

// Run executes the parsed AST, calculating values for the expressions.
func (p *Parser) Run(node ast.ExpressionNode) any {
	tokenTypes := *lexer.TokenTypeList
	switch n := node.(type) {
	case *ast.NumberNode:
		num, err := strconv.Atoi(n.Number.Text)
		if err != nil {
			panic(fmt.Sprintf("Невозможно преобразовать %s в число", n.Number.Text))
		}
		return num

	case *ast.UnarOperationNode:
		switch n.Operator.TypeToken.Name {
		case tokenTypes["SHOW"].Name:
			fmt.Println(p.Run(n.Operand))
			return nil
		}

	case *ast.BinOperationNode:
		// Сначала обрабатываем присваивание отдельно
		if n.Operator.TypeToken.Name == tokenTypes["ASSIGN"].Name {
			if variable, ok := n.LeftNode.(*ast.VariableNode); ok {
				rightValue := p.Run(n.RightNode)
				p.Scope[variable.Variable.Text] = rightValue
				fmt.Printf("Добавлена переменная %s в Scope с значением %v\n", variable.Variable.Text, rightValue)
				return rightValue
			}
			panic("Левая часть ASSIGN — не VariableNode")
		}
		leftValue := p.Run(n.LeftNode)
		rightValue := p.Run(n.RightNode)

		switch n.Operator.TypeToken.Name {
		case tokenTypes["PLUS"].Name:
			return leftValue.(int) + rightValue.(int)
		case tokenTypes["MINUS"].Name:
			return leftValue.(int) - rightValue.(int)
		case tokenTypes["MULTIPLY"].Name:
			return leftValue.(int) * rightValue.(int)
		case tokenTypes["DIVIDE"].Name:
			return leftValue.(int) / rightValue.(int)
		case tokenTypes["ASSIGN"].Name:
			if variable, ok := n.LeftNode.(*ast.VariableNode); ok {
				p.Scope[variable.Variable.Text] = rightValue
				fmt.Printf("Добавлена переменная %s в Scope с значением %v\n", variable.Variable.Text, rightValue)
				return rightValue
			}
			panic("Левая часть ASSIGN — не VariableNode")
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
