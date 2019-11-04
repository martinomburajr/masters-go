package eval

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

/* ==== AST ==== */

type Number float64

type Node interface {
	Eval() (Number, bool)
}

// Binary operator AST node
type Binary struct {
	op    byte
	left  Node
	right Node
}

func (n *Binary) Init(op byte, left, right Node) Node {
	n.op = op
	n.left = left
	n.right = right
	return n
}

func (n *Binary) Eval() (Number, bool) {
	left, ok := n.left.Eval()
	if !ok {
		return 0, false
	}
	right, ok := n.right.Eval()
	if !ok {
		return 0, false
	}
	switch n.op {
	case '+':
		return left + right, true
	case '-':
		return left - right, true
	case '*':
		return left * right, true
	case '/':
		if right == 0 {
			return Number(math.NaN()), false
		}
		return left / right, true
	}
	return 0, false
}

//func (n *Binary) String() string {
//	return fmt.Sprintf("(%s %c %s)", n.left, n.op, n.right)
//}

// Leaf value AST node
type Leaf struct {
	value Number
}

func (n *Leaf) Init(value Number) Node {
	n.value = value
	return n
}

func (n *Leaf) Eval() (Number, bool) {
	return n.value, true
}

//func (n *Leaf) String() string {
//	return fmt.Sprintf("%v", n.value)  // %v = default format
//}

/* ==== Lexer ==== */

type Lexer struct {
	data string
	pos  int
	Kind int
	Num  Number
	Oper byte
}

func (lexer *Lexer) Init(data string) *Lexer {
	lexer.data = data
	lexer.pos = 0
	return lexer
}

func (l *Lexer) Next() int {
	n := len(l.data)
	l.Kind = ERR
	if l.pos < n {
		char := l.data[l.pos]
		var prevChar uint8
		if l.pos != 0 {
			prevChar = l.data[l.pos-1]
		}
		switch char {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			var value Number = 0
			var divisor Number = 1
			for ; l.pos < n && '0' <= l.data[l.pos] && l.data[l.pos] <= '9'; l.pos++ {
				value = value*10 + Number(l.data[l.pos]-'0')
			}
			if l.pos < n && l.data[l.pos] == '.' {
				l.pos++
				for ; l.pos < n && '0' <= l.data[l.pos] && l.data[l.pos] <= '9'; l.pos++ {
					value = value*10 + Number(l.data[l.pos]-'0')
					divisor *= 10
				}
			}
			l.Kind = NUM
			l.Num = value / divisor
		case ')':
			l.pos++
			l.Kind = RPAR
			l.Oper = char
		case '(':
			l.pos++
			l.Kind = LPAR
			l.Oper = char
		case '+', '-', '*', '/':
			if char == '-' {
				switch prevChar {
				case '+', '-', '*', '/', '(':
					l.pos++
					l.Kind = NEG
					l.Oper = char
				default:
					l.pos++
					l.Kind = OP
					l.Oper = char
				}
			} else {
				l.pos++
				l.Kind = OP
				l.Oper = char
			}
		}
	}
	return l.Kind
}

/* ==== Parser ==== */

type Parser struct {
	lexer      *Lexer
	precedence map[byte]int
}

func (p *Parser) Init(data string) *Parser {
	p.lexer = new(Lexer).Init(data)
	p.precedence = make(map[byte]int)
	p.lexer.Next()
	return p
}

func (p *Parser) AddOperator(op byte, precedence int) {
	p.precedence[op] = precedence
}

func (p *Parser) Parse() (Node, bool) {
	lhs, ok := p.parsePrimary()
	if !ok {
		return nil, false
	}
	// starting with 1 instead of 0, because
	// map[*]int returns 0 for non-existant items
	node, ok := p.parseOperators(lhs, 1)
	if !ok {
		return nil, false
	}
	return node, true
}

func (p *Parser) parsePrimary() (Node, bool) {
	switch p.lexer.Kind {
	case NUM:
		node := new(Leaf).Init(p.lexer.Num)
		p.lexer.Next()
		return node, true
	case LPAR:
		p.lexer.Next()
		if p.lexer.Kind == NEG {
			p.lexer.Next()
			p.lexer.Num = p.lexer.Num * -1
		}
		node, ok := p.Parse()
		if !ok {
			return nil, false
		}
		if p.lexer.Kind == RPAR {
			p.lexer.Next()
		}
		return node, true
	case NEG:
		p.lexer.Next()
		p.lexer.Num = p.lexer.Num * -1
		node := new(Leaf).Init(p.lexer.Num)
		return node, true
	}
	return nil, false
}

const (
	ERR  = iota // error
	NUM         // number
	LPAR        // left parenthesis
	RPAR        // right parenthesis
	OP          // operator
	NEG         // a negative sign
)

func (p *Parser) parseOperators(lhs Node, min_precedence int) (Node, bool) {
	var ok bool
	var rhs Node
	for p.lexer.Kind == OP && p.precedence[p.lexer.Oper] >= min_precedence {
		op := p.lexer.Oper
		p.lexer.Next()
		rhs, ok = p.parsePrimary()
		if !ok {
			return nil, false
		}
		for p.lexer.Kind == OP && p.precedence[p.lexer.Oper] > p.precedence[op] {
			op2 := p.lexer.Oper
			rhs, ok = p.parseOperators(rhs, p.precedence[op2])
			if !ok {
				return nil, false
			}
		}
		lhs = new(Binary).Init(op, lhs, rhs)
		if p.lexer.pos < len(p.lexer.data) && p.lexer.Num < 0 {
			if p.lexer.Kind == OP {
				continue
			} else {
				p.lexer.Next()
			}
		}
	}
	return lhs, true
}

//Calculate is a streamlined evalutor that takes in a mathematical string and performs either +,
// - / * operations on the given string. It performs about 4 times better than existing methods as per the benchmarks..
// Divide by 0 returns not a number
func Calculate(substitutedExpression string) (float64, error) {
	var node Node
	var result Number
	var p *Parser
	var parseOk, evalOk bool

	substitutedExpression = NegativeNumberParser(substitutedExpression)
	p = new(Parser).Init(substitutedExpression)
	p.AddOperator('+', 1)
	p.AddOperator('-', 1)
	p.AddOperator('*', 2)
	p.AddOperator('/', 2)
	node, parseOk = p.Parse()
	if parseOk {
		result, evalOk = node.Eval()
		if evalOk {
			return float64(result), nil // %v = default format
		} else {
			return float64(result), fmt.Errorf("%s = Invalid Evaluation error\n", substitutedExpression)
		}
	} else {
		return 0.0, fmt.Errorf("%s = Invalid Syntax error\n", substitutedExpression)
	}
}

// MathError Kind = 0, syntax, Kind = 1 NaN
type MathError struct {
	Kind int
}

func CalculateWithVar(substitutedExpression string, variables map[string]float64) (float64, error) {
	if variables == nil {
		return Calculate(substitutedExpression)
	}
	for key := range variables {
		floatStr := strconv.FormatFloat(variables[key], 'g', 10, 64)
		substitutedExpression = strings.ReplaceAll(substitutedExpression, key, floatStr)
	}
	return Calculate(substitutedExpression)
}

// NegativeNumberParser adds a 0 if an odd number is at the start of a string
func NegativeNumberParser(str string) string {
	if str == "" {
		return ""
	}
	builder := strings.Builder{}
	if str[0] == '-' {
		builder.WriteByte('0')
	}
	builder.WriteString(str)
	return builder.String()
}
