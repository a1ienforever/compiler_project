package tac

import (
	"compiler_project/parser/ast" // замени на реальный путь к твоему ast пакету
	"fmt"
	"strconv"
	"strings"
)

type TACInstruction struct {
	Op   string
	Arg1 string
	Arg2 string
	Res  string
}

type TACBuilder struct {
	instructions []TACInstruction
	tempCount    int
	labelCount   int
}

func NewTACBuilder() *TACBuilder {
	return &TACBuilder{}
}

func (b *TACBuilder) newTemp() string {
	b.tempCount++
	return fmt.Sprintf("t%d", b.tempCount)
}

func (b *TACBuilder) Instructions() []TACInstruction {
	return b.instructions
}

func (b *TACBuilder) newLabel() string {
	b.labelCount++
	return fmt.Sprintf("L%d", b.labelCount)
}

func (b *TACBuilder) Generate(node ast.ExpressionNode) string {
	switch n := node.(type) {

	case *ast.NumberNode:
		return n.Number.Text

	case *ast.FloatNode:
		return n.Float.Text

	case *ast.StringNode:
		return fmt.Sprintf(`"%s"`, n.String.Text)

	case *ast.BooleanNode:
		return n.Boolean.Text

	case *ast.VariableNode:
		return n.Variable.Text

	case *ast.TypedAssignNode:
		val := b.Generate(n.Value)
		b.instructions = append(b.instructions, TACInstruction{
			Op:   "=",
			Arg1: val,
			Res:  n.Variable.Text,
		})
		return n.Variable.Text

	case *ast.BinOperationNode:
		leftConst, leftIsConst := extractConstant(n.LeftNode)
		rightConst, rightIsConst := extractConstant(n.RightNode)

		if leftIsConst && rightIsConst {
			result := evalConstantBinary(n.Operator.Text, leftConst, rightConst)
			return result
		}

		left := b.Generate(n.LeftNode)
		right := b.Generate(n.RightNode)
		temp := b.newTemp()
		b.instructions = append(b.instructions, TACInstruction{
			Op:   n.Operator.Text,
			Arg1: left,
			Arg2: right,
			Res:  temp,
		})
		return temp

	case *ast.ShowNode:
		val := b.Generate(n.Variable)
		b.instructions = append(b.instructions, TACInstruction{
			Op:   "show",
			Arg1: val,
		})
		return ""

	case *ast.StatementsNode:
		for _, stmt := range n.CodeStrings {
			b.Generate(stmt)
		}
		return ""

	case *ast.IfNode:
		cond := b.Generate(n.Condition)
		elseLabel := b.newLabel()
		endLabel := b.newLabel()

		b.instructions = append(b.instructions, TACInstruction{
			Op:   "iffalse",
			Arg1: cond,
			Res:  elseLabel,
		})

		b.Generate(n.TrueBranch)
		b.instructions = append(b.instructions, TACInstruction{
			Op:  "goto",
			Res: endLabel,
		})

		b.instructions = append(b.instructions, TACInstruction{
			Op:  "label",
			Res: elseLabel,
		})

		if n.FalseBranch != nil {
			b.Generate(n.FalseBranch)
		}

		b.instructions = append(b.instructions, TACInstruction{
			Op:  "label",
			Res: endLabel,
		})
		return ""

	case *ast.WhileNode:
		startLabel := b.newLabel()
		endLabel := b.newLabel()

		b.instructions = append(b.instructions, TACInstruction{
			Op:  "label",
			Res: startLabel,
		})

		cond := b.Generate(n.Condition)
		b.instructions = append(b.instructions, TACInstruction{
			Op:   "iffalse",
			Arg1: cond,
			Res:  endLabel,
		})

		b.Generate(n.Body)

		b.instructions = append(b.instructions, TACInstruction{
			Op:  "goto",
			Res: startLabel,
		})

		b.instructions = append(b.instructions, TACInstruction{
			Op:  "label",
			Res: endLabel,
		})
		return ""

	case *ast.FunctionDeclarationNode:
		b.instructions = append(b.instructions, TACInstruction{
			Op:  "func",
			Res: n.Name.Text,
		})

		b.Generate(n.Body)

		b.instructions = append(b.instructions, TACInstruction{
			Op:  "endfunc",
			Res: n.Name.Text,
		})

		return ""

	case *ast.FunctionCallNode:
		var argTemps []string
		for _, arg := range n.Arguments {
			tempVar := b.Generate(arg)
			argTemps = append(argTemps, tempVar)
		}

		resultTemp := b.newTemp()
		b.instructions = append(b.instructions, TACInstruction{
			Op:   "call",
			Arg1: n.Name.Text,
			Arg2: fmt.Sprintf("%d", len(argTemps)),
			Res:  resultTemp,
		})

		for i, arg := range argTemps {
			b.instructions = append(b.instructions, TACInstruction{
				Op:   fmt.Sprintf("arg%d", i),
				Arg1: arg,
			})
		}

		return resultTemp

	default:
		panic(fmt.Sprintf("неподдерживаемый тип узла: %T", node))
	}
}

func (b *TACBuilder) Print() {
	for _, instr := range b.instructions {
		switch instr.Op {
		case "+", "-", "*", "/", "equal", "non-equal", "less", "more":
			fmt.Printf("%s = %s %s %s\n", instr.Res, instr.Arg1, instr.Op, instr.Arg2)
		case "", "=":
			fmt.Printf("%s = %s\n", instr.Res, instr.Arg1)
		case "show":
			fmt.Printf("show %s\n", instr.Arg1)
		case "goto":
			fmt.Printf("goto %s\n", instr.Res)
		case "iffalse":
			fmt.Printf("iffalse %s goto %s\n", instr.Arg1, instr.Res)
		case "label":
			fmt.Printf("%s:\n", instr.Res)
		case "func":
			fmt.Printf("func %s\n", instr.Res)
		case "endfunc":
			fmt.Printf("endfunc %s\n", instr.Res)
		case "call":
			fmt.Printf("%s = call %s with %s args\n", instr.Res, instr.Arg1, instr.Arg2)
		default:
			if len(instr.Op) > 3 && instr.Op[:3] == "arg" {
				fmt.Printf("param %s\n", instr.Arg1)
			} else {
				fmt.Printf("// неизвестная инструкция: %+v\n", instr)
			}
		}
	}
}

func (b *TACBuilder) Optimize() {
	var optimized []TACInstruction
	used := make(map[string]bool)

	// Проход 1: Определим, какие переменные используются
	for _, instr := range b.instructions {
		if instr.Arg1 != "" {
			used[instr.Arg1] = true
		}
		if instr.Arg2 != "" {
			used[instr.Arg2] = true
		}
	}

	// Проход 2: Удалим инструкции, результат которых не используется
	for _, instr := range b.instructions {
		isTemp := strings.HasPrefix(instr.Res, "t")
		if instr.Res != "" &&
			!used[instr.Res] &&
			isTemp &&
			instr.Op != "show" && instr.Op != "call" && instr.Op != "goto" && instr.Op != "iffalse" && instr.Op != "label" && instr.Op != "func" && instr.Op != "endfunc" {
			continue // Удаляем только временные переменные
		}

		optimized = append(optimized, instr)
	}

	b.instructions = optimized
}

func extractConstant(node ast.ExpressionNode) (string, bool) {
	switch n := node.(type) {
	case *ast.NumberNode:
		return n.Number.Text, true
	case *ast.FloatNode:
		return n.Float.Text, true
	case *ast.BooleanNode:
		return n.Boolean.Text, true
	default:
		return "", false
	}
}

func evalConstantBinary(op, a, b string) string {
	switch op {
	case "+", "-", "*", "/":
		// Попробуем как целые
		ai, err1 := strconv.Atoi(a)
		bi, err2 := strconv.Atoi(b)
		if err1 == nil && err2 == nil {
			switch op {
			case "+":
				return strconv.Itoa(ai + bi)
			case "-":
				return strconv.Itoa(ai - bi)
			case "*":
				return strconv.Itoa(ai * bi)
			case "/":
				if bi != 0 {
					return strconv.Itoa(ai / bi)
				}
				return "0" // защита от деления на ноль
			}
		}
	case "equal":
		if a == b {
			return "true"
		}
		return "false"
	case "non-equal":
		if a != b {
			return "true"
		}
		return "false"
	}
	return "0"
}
