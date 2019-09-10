package evolution

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/martinomburajr/masters-go/utils"
	"math/rand"
	"strings"
)

// TODO generate AST tree from polynomial expression
type Program struct {
	ID string
	T  *DualTree
}

func GenerateProgramID(count int) string {
	return fmt.Sprintf("PROG-%d", count)
}

// ApplyStrategy takes a given strategy and applies a transformation to the given program.
// depth defines the exact depth the tree can evolve to given the transformation.
// Depth of a tree increases exponentially. So keep depths small e.g. 1,2,3
func (p *Program) ApplyStrategy(strategy Strategy, terminals []SymbolicExpression,
	nonTerminals []SymbolicExpression, mutationProbability float32, nonTerminalMutationProbability float32, depth int) (err error) {

	switch strategy {
	case StrategyAddSubTree:
		var tree *DualTree
		tree, err = GenerateRandomTree(depth, terminals, nonTerminals)
		err = p.T.AddSubTree(tree)
		break
	case StrategyDeleteSubTree:
		err = p.T.DeleteSubTree()
		break
	case StrategyMutateNode:
		chanceOfMutation := rand.Float32()
		if mutationProbability > chanceOfMutation {
			if nonTerminalMutationProbability > chanceOfMutation {
				err = p.T.MutateNonTerminal(nonTerminals)
			}
			err = p.T.MutateTerminal(terminals)
		}
		break
	default:
		break
	}
	return err
}

func (p *Program) Fitness() (float32, error) {
	return -1, fmt.Errorf("")
}

// Crossover is a evolutionary technique used to take two parents swap their genetic material and form two new children.
func Crossover(prog1 Program, prog2 Program) (child1 Program, child2 Program, err error) {
	return Program{}, Program{}, nil
}

// Mutation is an evolutionary technique used to randomly change parts of a Program.
func Mutation(prog Program) (Program, error) {
	return Program{}, nil
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

type Bug *Program
type Test *Program
