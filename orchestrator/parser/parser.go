package parser

import (
	"fmt"
	"strconv"
	"strings"

	"DistributedArithmeticExpressionCalculator/orchestrator/models"
)

// Tokenize разбивает входное выражение на токены.
// Поддерживаются числа, операторы и скобки; пробелы игнорируются.
func Tokenize(expr string) []string {
	var tokens []string
	var number strings.Builder

	for _, ch := range expr {
		if ch == ' ' {
			continue
		}
		if (ch >= '0' && ch <= '9') || ch == '.' {
			number.WriteRune(ch)
		} else {
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
			}
			if strings.ContainsRune("+-*/()", ch) {
				tokens = append(tokens, string(ch))
			} else {
				return nil // обнаружен недопустимый символ
			}
		}
	}
	if number.Len() > 0 {
		tokens = append(tokens, number.String())
	}
	return tokens
}

func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}

func precedence(op string) int {
	if op == "+" || op == "-" {
		return 1
	}
	if op == "*" || op == "/" {
		return 2
	}
	return 0
}

// ParseExpression преобразует строковое выражение в AST.
func ParseExpression(expr string) (*models.Node, error) {
	tokens := Tokenize(expr)
	if tokens == nil || len(tokens) == 0 {
		return nil, fmt.Errorf("некорректное выражение")
	}

	var outputQueue []*models.Node
	var opStack []string

	for _, token := range tokens {
		if isOperator(token) {
			// Если оператор, проверяем приоритеты
			for len(opStack) > 0 {
				top := opStack[len(opStack)-1]
				if isOperator(top) && precedence(top) >= precedence(token) {
					opStack = opStack[:len(opStack)-1]
					if len(outputQueue) < 2 {
						return nil, fmt.Errorf("некорректное выражение")
					}
					right := outputQueue[len(outputQueue)-1]
					left := outputQueue[len(outputQueue)-2]
					outputQueue = outputQueue[:len(outputQueue)-2]
					node := &models.Node{
						Operator: top,
						Left:     left,
						Right:    right,
						Computed: false,
					}
					left.Parent = node
					right.Parent = node
					outputQueue = append(outputQueue, node)
				} else {
					break
				}
			}
			opStack = append(opStack, token)
		} else if token == "(" {
			opStack = append(opStack, token)
		} else if token == ")" {
			foundParen := false
			for len(opStack) > 0 {
				top := opStack[len(opStack)-1]
				opStack = opStack[:len(opStack)-1]
				if top == "(" {
					foundParen = true
					break
				}
				// top — оператор, создаём узел
				if len(outputQueue) < 2 {
					return nil, fmt.Errorf("некорректное выражение")
				}
				right := outputQueue[len(outputQueue)-1]
				left := outputQueue[len(outputQueue)-2]
				outputQueue = outputQueue[:len(outputQueue)-2]
				node := &models.Node{
					Operator: top,
					Left:     left,
					Right:    right,
					Computed: false,
				}
				left.Parent = node
				right.Parent = node
				outputQueue = append(outputQueue, node)
			}
			if !foundParen {
				return nil, fmt.Errorf("скобки не сбалансированы")
			}
		} else {
			// Токен должен быть числом
			val, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return nil, fmt.Errorf("некорректное число: %s", token)
			}
			node := &models.Node{
				Value:    val,
				Computed: true,
			}
			outputQueue = append(outputQueue, node)
		}
	}

	// Обработка оставшихся операторов
	for len(opStack) > 0 {
		top := opStack[len(opStack)-1]
		opStack = opStack[:len(opStack)-1]
		if top == "(" || top == ")" {
			return nil, fmt.Errorf("скобки не сбалансированы")
		}
		if len(outputQueue) < 2 {
			return nil, fmt.Errorf("некорректное выражение")
		}
		right := outputQueue[len(outputQueue)-1]
		left := outputQueue[len(outputQueue)-2]
		outputQueue = outputQueue[:len(outputQueue)-2]
		node := &models.Node{
			Operator: top,
			Left:     left,
			Right:    right,
			Computed: false,
		}
		left.Parent = node
		right.Parent = node
		outputQueue = append(outputQueue, node)
	}
	if len(outputQueue) != 1 {
		return nil, fmt.Errorf("некорректное выражение")
	}
	return outputQueue[0], nil
}
