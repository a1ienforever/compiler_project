package tac

import (
	"compiler_project/parser/ast" // замени на реальный путь к твоему ast пакету
	"fmt"
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
		default:
			fmt.Printf("// неизвестная инструкция: %+v\n", instr)
		}
	}
}
