package parser

import (
	"compiler_project/lexer"
	"compiler_project/parser/ast"
	"fmt"
	"strconv"
)

type Parser struct {
	Tokens   []lexer.Token
	Position int
	Scope    map[string]interface{}
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{Tokens: tokens, Scope: map[string]interface{}{}}
}

func (p *Parser) Match(expected ...lexer.TokenType) *lexer.Token {
	if p.Position >= len(p.Tokens) {
		return nil
	}
	current := p.Tokens[p.Position]
	for _, exp := range expected {
		if current.TypeToken == exp {
			fmt.Printf("Token: Type=%s, Text='%s' - Обработан!\n", current.TypeToken.Name, current.Text)
			p.Position++
			return &current
		}
	}
	return nil
}

func (p *Parser) Require(expected ...lexer.TokenType) *lexer.Token {
	tok := p.Match(expected...)
	if tok == nil {
		panic(fmt.Sprintf("Ожидался токен из множества: %v", expected))
	}
	return tok
}

func (p *Parser) ParseCode() *ast.StatementsNode {
	var statements ast.StatementsNode
	tokenTypes := *lexer.TokenTypeList

	// Парсим выражения и добавляем их в statements
	for p.Position < len(p.Tokens) {
		stmt := p.ParseStatement()
		statements.AddNode(stmt)
		p.Require(tokenTypes["SEMICOLON"]) // Добавляем каждое выражение в список
	}
	fmt.Println(p.Tokens)
	return &statements
}

func (p *Parser) ParseStatement() ast.ExpressionNode {
	current := p.Tokens[p.Position]

	switch current.TypeToken.Name {
	case "int", "double", "string":
		stmt := p.parseTypedAssignment()
		return stmt
	case "show":
		return p.parseShowStatement()
	default:
		expr := p.ParseExpression()
		return expr
	}
}

func (p *Parser) parseShowStatement() ast.ExpressionNode {
	tokenTypes := *lexer.TokenTypeList

	p.Require(tokenTypes["SHOW"]) // "show"

	var_name := p.Require(tokenTypes["VARIABLE"]) // имя переменной
	variable := ast.NewVariableNode(*var_name)
	return ast.NewShowNode(variable)
}

func (p *Parser) ParseExpression() ast.ExpressionNode {

	if node := p.parseTypedAssignment(); node != nil {
		return node
	}
	return p.parseFormula()
}

func (p *Parser) parseTypedAssignment() ast.ExpressionNode {
	types := *lexer.TokenTypeList
	typeToken := p.Match(types["INT"], types["DOUB"], types["STR"])
	if typeToken == nil {
		return nil
	}
	variable := p.Require(types["VARIABLE"])
	p.Require(types["ASSIGN"])
	value := p.parseFormula()
	return ast.NewTypedAssignNode(*typeToken, *variable, value)
}

func (p *Parser) parseFormula() ast.ExpressionNode {
	left := p.parsePrimary()

	for {
		op := p.Match((*lexer.TokenTypeList)["PLUS"], (*lexer.TokenTypeList)["MINUS"])
		if op == nil {
			break
		}
		right := p.parsePrimary()
		left = ast.NewBinOperationNode(*op, left, right)
	}

	return left
}

func (p *Parser) parsePrimary() ast.ExpressionNode {
	types := *lexer.TokenTypeList
	if number := p.Match(types["INTEGER"]); number != nil {
		return ast.NewNumberNode(*number)
	}
	if flt := p.Match(types["DOUBLE"]); flt != nil {
		return ast.NewFloatNode(*flt)
	}
	if str := p.Match(types["STRING"]); str != nil {
		return ast.NewStringNode(*str)
	}
	if variable := p.Match(types["VARIABLE"]); variable != nil {
		return ast.NewVariableNode(*variable)
	}
	panic("Ожидалось выражение")
}

func (p *Parser) Run(node ast.ExpressionNode) interface{} {
	types := *lexer.TokenTypeList

	switch n := node.(type) {
	case *ast.NumberNode:
		val, err := strconv.Atoi(n.Number.Text)
		if err != nil {
			panic("Невозможно преобразовать int")
		}
		return val
	case *ast.FloatNode:
		val, err := strconv.ParseFloat(n.Float.Text, 64)
		if err != nil {
			panic("Невозможно преобразовать float")
		}
		return val
	case *ast.StringNode:
		return n.String.Text
	case *ast.VariableNode:
		return p.Scope[n.Variable.Text]
	case *ast.TypedAssignNode:
		val := p.Run(n.Value)
		switch n.Type.Type {
		case "int":
			if _, ok := val.(int); !ok {
				panic(fmt.Sprintf("Ожидался тип int, но получено %T", val))
			}
		case "double":
			if _, ok := val.(float64); !ok {
				panic(fmt.Sprintf("Ожидался тип double, но получено %T", val))
			}
		case "string":
			if _, ok := val.(string); !ok {
				panic(fmt.Sprintf("Ожидался тип string, но получено %T", val))
			}
		}
		p.Scope[n.Variable.Text] = val
		fmt.Printf("Добавлена типизированная переменная %s типа %s со значением %v\n", n.Variable.Text, n.Type.Type, val)
		return val
	case *ast.StatementsNode:
		var result interface{}
		for _, stmt := range n.CodeStrings {
			result = p.Run(stmt)
		}
		return result

	case *ast.ShowNode:
		val := p.Run(n.Variable)
		fmt.Printf(">> %v\n", val)
		return val
	case *ast.BinOperationNode:
		left := p.Run(n.LeftNode)
		right := p.Run(n.RightNode)

		switch l := left.(type) {
		case int:
			r, ok := right.(int)
			if !ok {
				panic(fmt.Sprintf("Ожидался int справа, но получено %T", right))
			}
			switch n.Operator.TypeToken {
			case types["PLUS"]:
				return l + r
			case types["MINUS"]:
				return l - r
			case types["MULTIPLY"]:
				return l * r
			}
			panic("Незвестный оператор!")
		case float64:
			r, ok := right.(float64)
			if !ok {
				panic(fmt.Sprintf("Ожидался float64 справа, но получено %T", right))
			}
			switch n.Operator.TypeToken {
			case types["PLUS"]:
				return l + r
			case types["MINUS"]:
				return l - r
			case types["MULTIPLY"]:
				return l * r
			}
			panic("Незвестный оператор!")
		case string:
			r, ok := right.(string)
			if !ok {
				panic(fmt.Sprintf("Ожидался string справа, но получено %T", right))
			}
			if n.Operator.TypeToken == types["PLUS"] {
				return l + r
			}
			panic("Операции над строками кроме + не поддерживаются")
		default:
			panic("Неподдерживаемые типы в бинарной операции")
		}
	default:
		panic("Неизвестная нода")
	}
}
