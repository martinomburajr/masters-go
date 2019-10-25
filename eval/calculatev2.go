package eval

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

func CalculateV2(mathematicalExpression string) (float64, error) {
	if mathematicalExpression == "" {
		return 0, fmt.Errorf("CalculateV2 | mathematicalExpression cannot be empty")
	}
	expr, err := parser.ParseExpr(mathematicalExpression)
	if err != nil {
		return 0, err
	}
	return Eval(expr)
}

func CalculateV2Var(mathematicalExpression string, variables map[string]float64) (float64, error) {
	if mathematicalExpression == "" {
		return 0, fmt.Errorf("CalculateV2 | mathematicalExpression cannot be empty")
	}
	for key := range variables {
		floatStr := strconv.FormatFloat(variables[key], 'g', 10, 64)
		mathematicalExpression = strings.ReplaceAll(mathematicalExpression, key, floatStr)
	}
	expr, err := parser.ParseExpr(mathematicalExpression)
	if err != nil {
		return 0, err
	}
	return Eval(expr)
}

func Eval(exp ast.Expr) (float64, error) {
	switch exp := exp.(type) {
	case *ast.UnaryExpr:
		return EvalUnaryExpr(exp)
	case *ast.BinaryExpr:
		return EvalBinaryExpr(exp)
	case *ast.BasicLit:
		switch exp.Kind {
		case token.INT:
			i, err := strconv.Atoi(exp.Value)
			return float64(i), err
		case token.FLOAT:
			return strconv.ParseFloat(exp.Value, 64)
		}
	}
	return 0, nil
}

func EvalUnaryExpr(exp *ast.UnaryExpr) (float64, error) {
	left, err := Eval(exp.X)
	if err != nil {
		return 0, err
	}
	switch exp.Op {
	case token.SUB:
		return 0 - left, nil
	default:
		return 0, nil
	}
}

func EvalBinaryExpr(exp *ast.BinaryExpr) (float64, error) {
	left, err := Eval(exp.X)
	right, err := Eval(exp.Y)
	if err != nil {
		return 0, err
	}

	switch exp.Op {
	case token.ADD:
		return left + right, nil
	case token.SUB:
		return left - right, nil
	case token.MUL:
		return left * right, nil
	case token.QUO:
		return left / right, nil
	}

	return 0, nil
}
