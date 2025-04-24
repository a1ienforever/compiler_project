package semantics

import (
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
		valType, err := tc.Check(n.Value)
		if err != nil {
			return "", err
		}
		declaredType := normalizeTypeName(n.Type.Type)
		if valType != declaredType {
			return "", fmt.Errorf("тип переменной %s задан как %s, но присваивается %s", n.Variable.Text, declaredType, valType)
		}
		tc.Scope[n.Variable.Text] = declaredType
		return declaredType, nil

	case *ast.StatementsNode:
		for _, stmt := range n.CodeStrings {
			_, err := tc.Check(stmt)
			if err != nil {
				return "", err
			}
		}
		return "void", nil

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
