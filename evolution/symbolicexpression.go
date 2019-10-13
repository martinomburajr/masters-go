package evolution

import (
	"fmt"
	"math"
	"strings"
)

// SymbolicExpressionSet represents a mathematical expression broken into symbolic expressions.
// For Example x+1 will be broken into a SymbolicExpressionSet of size 3,
// containing both terminal and non terminal information
type SymbolicExpressionSet []SymbolicExpression

type SymbolicExpression struct {
	arity int
	value string
	kind  int //0 terminal >0 non-terminal
}

func (n *SymbolicExpression) CreateNonTerminal(arity int, value string) {
	n.arity = arity
	n.value = value
	n.kind = 1
}

func (n *SymbolicExpression) CreateTerminal(arity int, value string) {
	n.arity = arity
	n.value = value
	n.kind = 0
}

func (n *SymbolicExpression) ToDualTreeNode(key string) *DualTreeNode {
	return &DualTreeNode{
		value: n.value,
		arity: n.arity,
		left:  nil,
		right: nil,
		key:   key,
	}
}

func GenerateTerminals(count int, symbolList []string) ([]SymbolicExpression, error) {
	if count > len(symbolList) {
		count = len(symbolList)
	}
	if count < 0 {
		count = 0
	}
	if symbolList == nil {
		return nil, fmt.Errorf("symbol list cannot be nil")
	}
	if len(symbolList) < 1 {
		return nil, fmt.Errorf("symbol list cannot be empty")
	}
	se := make([]SymbolicExpression, count)

	for i := 0; i < count; i++ {
		str := symbolList[i]
		sExp := SymbolicExpression{
			value: str,
			arity: 0,
			kind:  0,
		}
		se[i] = sExp
	}

	return se, nil
}

func GenerateNonTerminals(count int, symbolList []string) ([]SymbolicExpression, error) {
	if count > len(symbolList) {
		count = len(symbolList)
	}
	if count < 0 {
		count = 0
	}
	if symbolList == nil {
		return nil, fmt.Errorf("symbol list cannot be nil")
	}
	if len(symbolList) < 1 {
		return nil, fmt.Errorf("symbol list cannot be empty")
	}
	se := make([]SymbolicExpression, count)

	for i := 0; i < count; i++ {
		str := symbolList[i]
		sExp := SymbolicExpression{
			value: str,
			arity: 1,
			kind:  1,
		}
		se[i] = sExp
	}

	return se, nil
}

// ParseString parses a given mathematical expression into a set of terminals and nonTerminals within the string.
// It assumes mathematical expressions in particular non-terminals have an arity of two and take in two arguments e.
// g. * / - + are some examples
func ParseString(expression string, validNonTerminals []string) (terminals, nonTerminals,
	mathematicalExpression []SymbolicExpression,
	err error) {
	if expression == "" {
		return nil, nil, nil, fmt.Errorf("ParseKind | empty")
	}

	expression = strings.TrimSpace(expression)
	expression = strings.ReplaceAll(expression, " ", "")

	if len(expression)%2 == 0 {
		return nil, nil, nil, fmt.Errorf("ParseKind | expression length must be odd as any valid mathematical" +
			" expression where the operators have dual arity must take in two arguments at a minimum. " +
			"This means N (even) number of operators will have N+1 number of terminals resulting in an odd expression" +
			" " +
			"| nil")
	}
	if len(expression) > 0 && validNonTerminals == nil {
		return nil, nil, nil, fmt.Errorf("ParseKind | at least one kind of validNonTerminal is required | nil")
	}
	if len(expression) > 0 && len(validNonTerminals) < 1 {
		return nil, nil, nil, fmt.Errorf("ParseKind | at least one kind of validNonTerminal is required | empty")
	}
	if len(expression) == 1 {
		for i := range validNonTerminals {
			if validNonTerminals[i] == string(expression[0]) {
				return nil, nil, nil, fmt.Errorf("ParseKind | single element in expression must be a terminal and not part" +
					" of the validNonTerminals list")
			}
		}
		terminal := SymbolicExpression{
			arity: 0,
			kind:  0,
			value: string(expression[0]),
		}
		return []SymbolicExpression{terminal}, []SymbolicExpression{}, []SymbolicExpression{terminal}, nil
	}
	containsCount := 0
	for i := range validNonTerminals {
		if validNonTerminals[i] == "" {
			return nil, nil, nil, fmt.Errorf("ParseKind | validNonTerminal cannot contain empty string")
		}
		if strings.Contains(expression, validNonTerminals[i]) {
			containsCount++
		}
	}
	if containsCount < 1 && len(expression) > 1 {
		return nil, nil, nil, fmt.Errorf("ParseKind | validNonTerminal does not contain an unknown item in the" +
			" expression")
	}

	chars := strings.Split(expression, "")

	terminals = make([]SymbolicExpression, int(math.Ceil(float64(len(expression))/2)))
	nonTerminals = make([]SymbolicExpression, int(math.Floor(float64(len(expression))/2)))
	mathematicalExpression = make([]SymbolicExpression, len(expression))

	countx := 0
	county := 0
	for i := 0; i < len(chars); i++ {
		var item SymbolicExpression
		if i%2 == 0 {
			item = SymbolicExpression{
				arity: 0,
				kind:  0,
				value: chars[i],
			}
			terminals[countx] = item
			countx++
		}
		if i%2 == 1 {

			item = SymbolicExpression{
				arity: 2,
				kind:  1,
				value: chars[i],
			}
			nonTerminals[county] = item
			county++
		}
		mathematicalExpression[i] = item
	}

	return terminals, nonTerminals, mathematicalExpression, nil
}

// ParseStringLiberal parses a given mathematical expression into a set of terminals and nonTerminals within the string.
// It assumes mathematical expressions in particular non-terminals have an arity of two and take in two arguments e.
// g. * / - + are some examples
func ParseStringLiberal(expression string) (terminals, nonTerminals,
	mathematicalExpression []SymbolicExpression,
	err error) {
	if expression == "" {
		return nil, nil, nil, fmt.Errorf("ParseKind | empty")
	}

	expression = strings.TrimSpace(expression)
	expression = strings.ReplaceAll(expression, " ", "")

	if len(expression)%2 == 0 {
		return nil, nil, nil, fmt.Errorf("ParseKind | expression length must be odd as any valid mathematical" +
			" expression where the operators have dual arity must take in two arguments at a minimum. " +
			"This means N (even) number of operators will have N+1 number of terminals resulting in an odd expression" +
			" " +
			"| nil")
	}
	if len(expression) == 1 {
		terminal := SymbolicExpression{
			arity: 0,
			kind:  0,
			value: string(expression[0]),
		}
		return []SymbolicExpression{terminal}, []SymbolicExpression{}, []SymbolicExpression{terminal}, nil
	}

	chars := strings.Split(expression, "")

	terminals = make([]SymbolicExpression, int(math.Ceil(float64(len(expression))/2)))
	nonTerminals = make([]SymbolicExpression, int(math.Floor(float64(len(expression))/2)))
	mathematicalExpression = make([]SymbolicExpression, len(expression))

	countx := 0
	county := 0
	for i := 0; i < len(chars); i++ {
		var item SymbolicExpression
		if i%2 == 0 {
			item = SymbolicExpression{
				arity: 0,
				kind:  0,
				value: chars[i],
			}
			terminals[countx] = item
			countx++
		}
		if i%2 == 1 {

			item = SymbolicExpression{
				arity: 2,
				kind:  1,
				value: chars[i],
			}
			nonTerminals[county] = item
			county++
		}
		mathematicalExpression[i] = item
	}

	return terminals, nonTerminals, mathematicalExpression, nil
}

func GenerateExpression(expression string) ([]SymbolicExpression, error) {
	//if count > len(symbolList) {
	//	count = len(symbolList)
	//}
	//if count < 0 {
	//	count = 0
	//}
	//if symbolList == nil {
	//	return nil, fmt.Errorf("symbol list cannot be nil")
	//}
	//if len(symbolList) < 1{
	//	return nil, fmt.Errorf("symbol list cannot be empty")
	//}
	//se := make([]SymbolicExpression, count)
	//
	//for i := 0; i < count; i++ {
	//	str := symbolList[i]
	//	sExp := SymbolicExpression{
	//		value: str,
	//		arity: 1,
	//		kind: 1,
	//	}
	//	se[i] = sExp
	//}
	//
	//return se, nil
	return nil, nil
}

// GenerateN generates a random SymbolicExpressionSet representing a valid mathematical expression.
// If size is less than 0, it reverts it to 0
func GenerateRandomSymbolicExpressionSet(size int) []SymbolicExpression {
	if size < 0 {
		size = 0
	}
	symbolicExpressions := make([]SymbolicExpression, 1)
	symbolicExpressions[0] = X1

	if size < 3 {
		return symbolicExpressions
	}
	for i := 1; i < size; i += 2 {
		if i%2 == 1 && i < (size-1) {
			symbolicExpressions = append(symbolicExpressions, Add)
			symbolicExpressions = append(symbolicExpressions, X1)
		}
	}

	return symbolicExpressions
}
