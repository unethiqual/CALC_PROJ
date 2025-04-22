package parser_test

import (
	"fmt"
	"testing"

	"github.com/unethiqual/CALC_PROJ/orchestrator/models"
	"DistributedArithmeticExpressionCalculator/orchestrator/parser"
)

// evaluate recursively computes the value of an AST node.
func evaluate(node *models.Node) (float64, error) {
	if node == nil {
		return 0, fmt.Errorf("nil node")
	}
	if node.Computed {
		return node.Value, nil
	}
	left, err := evaluate(node.Left)
	if err != nil {
		return 0, err
	}
	right, err := evaluate(node.Right)
	if err != nil {
		return 0, err
	}
	switch node.Operator {
	case "+":
		return left + right, nil
	case "-":
		return left - right, nil
	case "*":
		return left * right, nil
	case "/":
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return left / right, nil
	default:
		return 0, fmt.Errorf("unknown operator: %s", node.Operator)
	}
}

func TestParseExpression_WithSpaces(t *testing.T) {
	expr := "2 + 2 * (2 + 5) * 3"
	ast, err := parser.ParseExpression(expr)
	if err != nil {
		t.Fatalf("error parsing expression: %v", err)
	}
	result, err := evaluate(ast)
	if err != nil {
		t.Fatalf("error evaluating AST: %v", err)
	}
	expected := 2 + 2*(2+5)*3 // 2 + 2*7*3 = 2 + 42 = 44
	if result != float64(expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestParseExpression_NoSpaces(t *testing.T) {
	expr := "2+2*(2+5)*3"
	ast, err := parser.ParseExpression(expr)
	if err != nil {
		t.Fatalf("error parsing expression: %v", err)
	}
	result, err := evaluate(ast)
	if err != nil {
		t.Fatalf("error evaluating AST: %v", err)
	}
	expected := 44
	if result != float64(expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestParseExpression_Invalid(t *testing.T) {
	expr := "2+*3"
	_, err := parser.ParseExpression(expr)
	if err == nil {
		t.Fatalf("expected error for invalid expression, but got none")
	}
}
