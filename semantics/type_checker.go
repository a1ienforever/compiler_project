package semantics

import (
	"compiler_project/lexer"
	"compiler_project/parser/ast"
	"fmt"
)

type TypeChecker struct {
	Scope map[string]string // имя переменной → тип (например: "x" → "int")
}

func NewTypeChecker() *TypeChecker {
	return &TypeChecker{Scope: map[string]string{}}
}

func (tc *TypeChecker) Check(node ast.ExpressionNode) (string, error) {
	switch n := node.(type) {

	case *ast.NumberNode:
		return "int", nil

	case *ast.FloatNode:
		return "double", nil

	case *ast.StringNode:
		return "string", nil

	case *ast.BooleanNode:
		return "boolean", nil

	case *ast.VariableNode:
		t, ok := tc.Scope[n.Variable.Text]
		if !ok {
			return "", fmt.Errorf("переменная %s не определена", n.Variable.Text)
		}
		return t, nil

	case *ast.TypedAssignNode:
		declaredType := normalizeTypeName(n.Type.Type)
		// Временно запоминаем тип переменной, чтобы она была видна внутри Check(n.Value)
		tc.Scope[n.Variable.Text] = declaredType

		valType, err := tc.Check(n.Value)
		if err != nil {
			return "", err
		}

		if valType != declaredType {
			return "", fmt.Errorf("тип переменной %s задан как %s, но присваивается %s", n.Variable.Text, declaredType, valType)
		}

		return declaredType, nil

	case *ast.StatementsNode:
		for _, stmt := range n.CodeStrings {
			_, err := tc.Check(stmt)
			if err != nil {
				return "", err
			}
		}
		return "void", nil
	case *ast.BinOperationNode:
		types := *lexer.TokenTypeList

		leftType, err := tc.Check(n.LeftNode)
		if err != nil {
			return "", err
		}
		rightType, err := tc.Check(n.RightNode)
		if err != nil {
			return "", err
		}

		// Проводим проверку типов для операций EQUAL и NONEQUAL
		switch n.Operator.TypeToken {
		case types["EQUAL"], types["NONEQUAL"]:
			// Проверка на типы, которые поддерживают операцию сравнения
			if leftType != rightType {
				return "", fmt.Errorf("недопустимое сравнение типов: %s и %s", leftType, rightType)
			}
			// Можно добавить дополнительные проверки для типов, если они должны быть ограничены
			switch leftType {
			case "int", "double", "string", "boolean":
				// Поддерживаем сравнение для этих типов
				return "boolean", nil
			default:
				return "", fmt.Errorf("операция %s не поддерживается для типа %s", n.Operator.TypeToken, leftType)
			}

		// Поддержка других типов бинарных операций, например, для чисел
		case types["GREATER"], types["LESS"]:
			if leftType != rightType {
				return "", fmt.Errorf("недопустимое сравнение типов: %s и %s", leftType, rightType)
			}
			if leftType == "boolean" {
				return "", fmt.Errorf("операция %s не поддерживается для типа boolean", n.Operator.TypeToken)
			}
			return "boolean", nil

		// Остальные бинарные операции, например, AND, OR, которые могут быть логическими
		case types["AND"], types["OR"]:
			if leftType != "boolean" || rightType != "boolean" {
				return "", fmt.Errorf("логическая операция %s требует типов boolean", n.Operator.TypeToken)
			}
			return "boolean", nil
		case types["PLUS"], types["MINUS"], types["MULTIPLY"], types["DIVIDE"]:
			if (leftType == "int" || leftType == "double") && leftType == rightType {
				return leftType, nil
			}
			return "", fmt.Errorf("арифметическая операция %s требует совпадающих числовых типов, получено: %s и %s", n.Operator.TypeToken, leftType, rightType)

		default:
			return "", fmt.Errorf("неподдерживаемая операция %s для типов %s и %s", n.Operator.TypeToken, leftType, rightType)
		}
	case *ast.IfNode:
		condType, err := tc.Check(n.Condition)
		if err != nil {
			return "", err
		}
		if condType != "boolean" {
			return "", fmt.Errorf("условие в if должно быть boolean, получено: %s", condType)
		}
		_, err = tc.Check(n.TrueBranch)
		if err != nil {
			return "", err
		}
		if n.FalseBranch != nil {
			_, err = tc.Check(n.FalseBranch)
			if err != nil {
				return "", err
			}
		}
		return "void", nil
	case *ast.WhileNode:
		// Проверка типа условия в цикле while
		condType, err := tc.Check(n.Condition)
		if err != nil {
			return "", err
		}
		if condType != "boolean" {
			return "", fmt.Errorf("условие в while должно быть boolean, получено: %s", condType)
		}

		// Проверка тела цикла
		_, err = tc.Check(n.Body)
		if err != nil {
			return "", err
		}
		return "void", nil

	case *ast.ShowNode:
		_, err := tc.Check(n.Variable)
		if err != nil {
			return "", err
		}
		return "void", nil

	default:
		return "", fmt.Errorf("неизвестный тип AST узла: %T", node)
	}
}

func normalizeTypeName(t string) string {
	switch t {
	case "int", "double", "string", "boolean":
		return t
	case "INT":
		return "int"
	case "DOUB":
		return "double"
	case "STR":
		return "string"
	case "BOOLEAN":
		return "boolean"
	default:
		return t
	}
}
