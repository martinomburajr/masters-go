package program

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/martinomburajr/masters-go/program/tree/dualtree"
	"github.com/martinomburajr/masters-go/utils"
	"strings"
)

// TODO generate AST tree from polynomial expression
type Program struct {
	ID                   string
	T                    *dualtree.DualTree
	hasAppliedStrategies bool
}

func (p *Program) ApplyStrategy() {
	return
}

func (p *Program) Fitness() float32 {
	return 0
}

// Eval is a simple helper function that takes in an independent variable,
// uses the programs tree to compute the resultant value
func (p *Program) Eval(independentVar float32) (float32, error) {
	if p.T == nil {
		return -1, fmt.Errorf("program: %v -> tree is nil", p.ID)
	}

	err := p.T.Validate()
	if err != nil {
		return -1, err
	}

	expressionString, err := p.T.ToMathematicalString()
	if err != nil {
		return -1, err
	}

	indepStr := fmt.Sprintf("%f", independentVar)
	mathematicalExpression := strings.ReplaceAll(expressionString, "x", indepStr)

	expression, err := govaluate.NewEvaluableExpression(mathematicalExpression)
	if err != nil {
		return -1, err
	}

	result, err := expression.Evaluate(nil)
	if err != nil {
		return -1, err
	}

	ans, err := utils.ConvertToFloat(result)
	if err != nil {
		return -1, err
	}

	return ans, nil
}



//func (p *Program) Terminals() []*dualtree.Terminal {
//	return nil
//}
//
//func (p *Program) NonTerminals() []*dualtree.NodeType {
//	return nil
//}

func (p *Program) Mutate() {

}

func (p *Program) Recombine() {

}

func (p *Program) Validate() error {
	return nil
}

type Bug *Program
type Test *Program

