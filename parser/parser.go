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
	case "if":
		return p.parseIfStatement()
	case "while":
		return p.parseWhileStatement()
	default:
		expr := p.ParseExpression()
		return expr
	}
}

func (p *Parser) parseIfStatement() ast.ExpressionNode {
	types := *lexer.TokenTypeList

	p.Require(types["IF"])

	condition := p.ParseExpression()
	p.Require(types["LBRACE"]) // {

	trueBranch := &ast.StatementsNode{}
	for p.Match(types["RBRACE"]) == nil { // }
		stmt := p.ParseStatement()
		trueBranch.AddNode(stmt)
		p.Require(types["SEMICOLON"])
	}

	var falseBranch *ast.StatementsNode
	if p.Match(types["ELSE"]) != nil {
		p.Require(types["LBRACE"]) // {

		falseBranch = &ast.StatementsNode{}
		for p.Match(types["RBRACE"]) == nil { // }
			stmt := p.ParseStatement()
			falseBranch.AddNode(stmt)
			p.Require(types["SEMICOLON"])
		}
	}

	return ast.NewIfNode(condition, trueBranch, falseBranch)
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
	typeToken := p.Match(types["INT"], types["DOUB"], types["STR"], types["BOOLEAN"])
	if typeToken == nil {
		return nil
	}
	variable := p.Require(types["VARIABLE"])
	p.Require(types["ASSIGN"])
	value := p.parseFormula()
	return ast.NewTypedAssignNode(*typeToken, *variable, value)
}

func (p *Parser) parseFormula() ast.ExpressionNode {
	return p.parseLogicalOr()
}

func (p *Parser) parseLogicalOr() ast.ExpressionNode {
	node := p.parseLogicalAnd()
	types := *lexer.TokenTypeList
	for {
		op := p.Match(types["OR"])
		if op == nil {
			break
		}
		right := p.parseLogicalAnd()
		node = ast.NewBinOperationNode(*op, node, right)
	}
	return node
}

func (p *Parser) parseLogicalAnd() ast.ExpressionNode {
	node := p.parseEquality()
	types := *lexer.TokenTypeList
	for {
		op := p.Match(types["AND"])
		if op == nil {
			break
		}
		right := p.parseEquality()
		node = ast.NewBinOperationNode(*op, node, right)
	}
	return node
}

func (p *Parser) parseEquality() ast.ExpressionNode {
	node := p.parseTerm()
	types := *lexer.TokenTypeList
	for {
		op := p.Match(types["EQUAL"], types["NONEQUAL"])
		if op == nil {
			break
		}
		right := p.parseTerm()
		node = ast.NewBinOperationNode(*op, node, right)
	}
	return node
}

func (p *Parser) parseTerm() ast.ExpressionNode {
	node := p.parseFactor()
	types := *lexer.TokenTypeList
	for {
		op := p.Match(types["PLUS"], types["MINUS"])
		if op == nil {
			break
		}
		right := p.parseFactor()
		node = ast.NewBinOperationNode(*op, node, right)
	}
	return node
}

func (p *Parser) parseFactor() ast.ExpressionNode {
	node := p.parsePrimary()
	types := *lexer.TokenTypeList
	for {
		op := p.Match(types["MULTIPLY"], types["DIVIDE"])
		if op == nil {
			break
		}
		right := p.parsePrimary()
		node = ast.NewBinOperationNode(*op, node, right)
	}
	return node
}

func (p *Parser) parsePrimary() ast.ExpressionNode {
	types := *lexer.TokenTypeList

	if p.Match(types["LPAREN"]) != nil { // обрабатываем (
		expr := p.parseFormula()   // рекурсивно разбираем вложенное выражение
		p.Require(types["RPAREN"]) // обрабатываем )
		return expr
	}

	if b := p.Match(types["TRUE"], types["FALSE"]); b != nil {
		return ast.NewBooleanNode(*b)
	}
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

func (p *Parser) parseWhileStatement() ast.ExpressionNode {
	types := *lexer.TokenTypeList

	// Ожидаем while
	p.Require(types["WHILE"])

	// Условие
	condition := p.ParseExpression()

	// Ожидаем {
	p.Require(types["LBRACE"])

	// Блок команд внутри цикла
	body := &ast.StatementsNode{}
	for p.Match(types["RBRACE"]) == nil { // }
		stmt := p.ParseStatement()
		body.AddNode(stmt)
		p.Require(types["SEMICOLON"])
	}

	// Создаём и возвращаем узел while
	return ast.NewWhileNode(condition, body)
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
	case *ast.BooleanNode:
		return n.Boolean.TypeToken.Name == "TRUE"
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
		case "boolean":
			if _, ok := val.(bool); !ok {
				panic(fmt.Sprintf("Ожидался тип boolean, но получено %T", val))
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
	case *ast.IfNode:
		cond := p.Run(n.Condition)
		condVal, ok := cond.(bool)
		if !ok {
			panic("Условие в if должно быть boolean")
		}
		if condVal {
			return p.Run(n.TrueBranch)
		}
		if n.FalseBranch != nil {
			return p.Run(n.FalseBranch)
		}
		return nil
	case *ast.WhileNode:
		cond := p.Run(n.Condition)
		condVal, ok := cond.(bool)
		if !ok {
			panic("Условие в while должно быть boolean")
		}
		for condVal {
			p.Run(n.Body)
			cond = p.Run(n.Condition)
			condVal, ok = cond.(bool)
			if !ok {
				panic("Условие в while должно быть boolean")
			}
		}
		return nil

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
			case types["DIVIDE"]:
				return l / r
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
			case types["DIVIDE"]:
				return l / r
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
		case bool:
			r, ok := right.(bool)
			if !ok {
				panic(fmt.Sprintf("Ожидался boolean справа, но получено %T", right))
			}
			switch n.Operator.TypeToken {
			case types["EQUAL"]:
				return l == r
			case types["NONEQUAL"]:
				return l != r
			case types["AND"]:
				return l && r
			case types["OR"]:
				return l || r
			default:
				panic("Неподдерживаемая логическая операция")
			}
		default:
			panic("Неподдерживаемые типы в бинарной операции")
		}
	default:
		panic("Неизвестная нода")
	}
}
