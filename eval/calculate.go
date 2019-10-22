package eval

import (
	"fmt"
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
			return 0, false
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

const (
	ERR  = iota // error
	NUM         // number
	LPAR        // left parenthesis
	RPAR        // right parenthesis
	OP          // operator
)

func (lexer *Lexer) Init(data string) *Lexer {
	lexer.data = data
	lexer.pos = 0
	return lexer
}

func (l *Lexer) Next() int {
	n := len(l.data)
	l.Kind = ERR
	if l.pos < n {
		switch char := l.data[l.pos]; char {
		case '+', '-', '*', '/':
			l.pos++
			l.Kind = OP
			l.Oper = char
		case '(':
			l.pos++
			l.Kind = LPAR
			l.Oper = char
		case ')':
			l.pos++
			l.Kind = RPAR
			l.Oper = char
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
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
		node, ok := p.Parse()
		if !ok {
			return nil, false
		}
		if p.lexer.Kind == RPAR {
			p.lexer.Next()
		}
		return node, true
	}
	return nil, false
}

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
	}
	return lhs, true
}

func CalcEvaler() {

}

//Calculate is a streamlined evalutor that takes in a mathematical string and performs either +,
// - / * operations on the given string. It performs about 4 times better than existing methods as per the benchmarks..
func Calculate(substitutedExpression string) (float64, error) {
	var node Node
	var result Number
	var p *Parser
	var parseOk, evalOk bool

	//substitutedExpression = MartinsReplace(substitutedExpression, " ", "")
	//substitutedExpression = NegativeNumberParser(substitutedExpression)
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
			return 0.0, fmt.Errorf("%s = Evaluation error\n", substitutedExpression)
		}
	} else {
		return 0.0, fmt.Errorf("%s = Syntax error\n", substitutedExpression)
	}
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

// MartinsReplace can only replace a single item string
func MartinsReplace(str string, old, new string) string {
	sb := strings.Builder{}
	for i := 0; i < len(str); i++ {
		if str[i] == []byte(old)[0] {
			sb.WriteString(new)
			continue
		}
		sb.WriteByte(str[i])
	}
	return sb.String()
}

// Calcu

// NegativeNumberParser is used to generate a representation that is aware of negative numbers in all possible
// mathematical variations e.g. -1 1--1 etc. It acts as a sanitization method before the actual evaluation happens
func NegativeNumberParser(str string) string {
	builder := strings.Builder{}
	if len(str) == 0 {
		return ""
	}
	if str[0] == '-' {
		builder.WriteString("0-") //Change to write byte
		if str[1] == '(' && str[2] == '-' {
			builder.WriteString("(0") //Change to write byte
			fixerAttempts := 0
			for i := 2; i < len(str); i++ {
				if str[i] == '-' && str[i+1] == '-' {
					builder.WriteString("+")
					i++
					fixerAttempts++
				} else {
					builder.WriteByte(str[i])
				}
			}
			for i := len(str) - fixerAttempts; i < len(str)-1; i++ {
				if str[i] != '-' {
					builder.WriteByte(str[i])
				}
			}
			return builder.String()
		}
	} else if str[0] == '(' && str[1] == '-' {
		builder.WriteString("(0") //Change to write byte
	} else {
		builder.WriteByte(str[0])
	}
	fixerAttempts := 0
	addEndParenth := 0
	for i := 1; i < len(str); i++ {
		if str[i] == '-' && str[i+1] == '-' {
			builder.WriteString("+")
			i++
			fixerAttempts++
		} else if str[i] == '(' && str[i+1] == '-' {
			builder.WriteString("(0-")
			i++
			fixerAttempts++
		} else if str[i] == '*' && str[i+1] == '-' {
			builder.WriteString("*(0-")
			i += 1
			fixerAttempts += 1
			addEndParenth += 1
		} else if str[i] == '+' && str[i+1] == '-' {
			builder.WriteString("+(0-")
			i += 1
			fixerAttempts += 1
			addEndParenth += 1
		} else if str[i] == '/' && str[i+1] == '-' {
			builder.WriteString("/(0-")
			i += 1
			fixerAttempts += 1
			addEndParenth += 1
		} else {
			builder.WriteByte(str[i])
		}
	}
	//g := builder.String()
	//log.Println(g)
	for i := len(str) - (fixerAttempts); i < len(str); i++ {
		//u := string(str[i])
		//log.Print(u)
		if str[i] != '-' {
			builder.WriteByte(str[i])
		}
	}
	for i := 0; i < addEndParenth; i++ {
		builder.WriteByte(')')
	}

	return builder.String()
}

// NegativeNumberParser is used to generate a representation that is aware of negative numbers in all possible
// mathematical variations e.g. -1 1--1 etc. It acts as a sanitization method before the actual evaluation happens
//func NegativeNumberParser2(str string) string {
//	str1 := strings.ReplaceAll(str, "(-", "(0-")
//	str2 := strings.ReplaceAll(str1, "(-", "(0-")
//	str3 := strings.ReplaceAll(str2, "(-", "(0-")
//}
