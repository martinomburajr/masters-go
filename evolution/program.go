package evolution

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/martinomburajr/masters-go/utils"
	"math/rand"
	"strings"
	"time"
)

const DeletionTypeMalicious = 1
const DeletionTypeSafe = 0

// TODO generate AST treeNode from polynomial expression
type Program struct {
	ID string
	T  DualTree
}

func GenerateProgramID(count int) string {
	randString := RandString(2)
	return fmt.Sprintf("%s-%s-%d", "PROG", randString, count)
}

// ApplyStrategy takes a given strategy and applies a transformation to the given program.
// depth defines the exact depth the treeNode can evolve to given the transformation.
// Depth of a treeNode increases exponentially. So keep depths small e.g. 1,2,3
// Ensure to place the independent variabel e.g X at the start of the SymbolicExpression terminals array.
// Otherwise there is less of a chance of having the independent variable propagate.
// The system is designed such that the first element of the terminals array will be the most prominent with regards
// to appearance.
func (p *Program) ApplyStrategy(strategy Strategy, terminals []SymbolicExpression,
	nonTerminals []SymbolicExpression, mutationProbability float32, nonTerminalMutationProbability float32,
	depth int, deletionStrategy int) (err error) {

	switch strategy {
	case StrategyAddSubTree:
		var tree DualTree
		tree, err = GenerateRandomTreeEnforceIndependentVariable(depth, terminals[0], terminals, nonTerminals)
		err = p.T.AddSubTree(tree)
		break
	case StrategyDeleteSubTree:
		err = p.T.DeleteSubTree(deletionStrategy)
		break
	case StrategyMutateNode:
		rand.Seed(time.Now().UnixNano())
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

// Mutation is an evolutionary technique used to randomly change parts of a Program.
func Mutation(prog Program) (Program, error) {
	return Program{}, nil
}

// Eval is a simple helper function that takes in an independent variable,
// uses the programs treeNode to compute the resultant value
func (p *Program) Eval(independentVar float32) (float32, error) {
	if p.T.root == nil {
		return -1, fmt.Errorf("program: %v -> treeNode.root is nil", p.ID)
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

func (p Program) Clone() (Program, error) {
	dualTree, err := p.T.Clone()
	if err != nil {
		return Program{}, err
	}
	p.T = dualTree
	p.ID = GenerateProgramID(0)
	return p, nil
}

type Bug *Program
type Test *Program
