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

	p.Require(tokenTypes["SHOW"])

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
		op := p.Match(types["EQUAL"], types["NONEQUAL"], types["MORE"], types["LESS"])
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

	if p.Match(types["LPAREN"]) != nil {
		expr := p.parseFormula()
		p.Require(types["RPAREN"])
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

	p.Require(types["WHILE"])

	condition := p.ParseExpression()

	p.Require(types["LBRACE"])

	body := &ast.StatementsNode{}
	for p.Match(types["RBRACE"]) == nil { // }
		stmt := p.ParseStatement()
		body.AddNode(stmt)
		p.Require(types["SEMICOLON"])
	}

	return ast.NewWhileNode(condition, body)
}

func (p *Parser) Run(node ast.ExpressionNode) interface{} {
	types := *lexer.TokenTypeList

	switch n := node.(type) {
	case *ast.NumberNode:
		val, _ := strconv.Atoi(n.Number.Text)
		return val
	case *ast.FloatNode:
		val, _ := strconv.ParseFloat(n.Float.Text, 64)
		return val
	case *ast.StringNode:
		return n.String.Text
	case *ast.BooleanNode:
		return n.Boolean.TypeToken.Name == "TRUE"
	case *ast.VariableNode:
		return p.Scope[n.Variable.Text]
	case *ast.TypedAssignNode:
		val := p.Run(n.Value)

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
		condVal, _ := cond.(bool)
		if condVal {
			return p.Run(n.TrueBranch)
		}
		if n.FalseBranch != nil {
			return p.Run(n.FalseBranch)
		}
		return nil
	case *ast.WhileNode:
		cond := p.Run(n.Condition)
		condVal, _ := cond.(bool)

		for condVal {
			p.Run(n.Body)
			cond = p.Run(n.Condition)
			condVal, _ = cond.(bool)
		}
		return nil

	case *ast.BinOperationNode:
		left := p.Run(n.LeftNode)
		right := p.Run(n.RightNode)

		switch l := left.(type) {
		case int:
			r, _ := right.(int)
			switch n.Operator.TypeToken {
			case types["PLUS"]:
				return l + r
			case types["MINUS"]:
				return l - r
			case types["MULTIPLY"]:
				return l * r
			case types["DIVIDE"]:
				return l / r
			case types["EQUAL"]:
				return l == r
			case types["NONEQUAL"]:
				return l != r
			case types["MORE"]:
				return l > r
			case types["LESS"]:
				return l < r
			}
		case float64:
			r, _ := right.(float64)
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
		case string:
			r, _ := right.(string)
			if n.Operator.TypeToken == types["PLUS"] {
				return l + r
			}
		case bool:
			r, _ := right.(bool)
			switch n.Operator.TypeToken {
			case types["EQUAL"]:
				return l == r
			case types["NONEQUAL"]:
				return l != r
			case types["AND"]:
				return l && r
			case types["OR"]:
				return l || r
			}
		default:
			panic("Неподдерживаемые типы в бинарной операции")
		}
	default:
		panic("Неизвестная нода")
	}
	return nil
}
