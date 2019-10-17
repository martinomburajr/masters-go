package evolution

import (
	"fmt"
	"github.com/martinomburajr/masters-go/eval"
	"github.com/martinomburajr/masters-go/utils"
)

const DeletionTypeMalicious = 1
const DeletionTypeSafe = 0

// TODO generate AST treeNode from polynomial expression
type Program struct {
	ID string
	T  *DualTree
}

func GenerateProgramID(count int) string {
	randString := RandString(2)
	return fmt.Sprintf("%s-%s-%d", "PROG", randString, count)
}

// ApplyStrategy takes a given Strategy and applies a transformation to the given program.
// depth defines the exact depth the treeNode can evolve to given the transformation.
// Depth of a treeNode increases exponentially. So keep depths small e.g. 1,2,3
// Ensure to place the independent variabel e.g X at the start of the SymbolicExpression terminals array.
// Otherwise there is less of a chance of having the independent variable propagate.
// The system is designed such that the first element of the terminals array will be the most prominent with regards
// to appearance.
func (p *Program) ApplyStrategy(strategy Strategy, terminals []SymbolicExpression,
	nonTerminals []SymbolicExpression, mutationProbability float64, nonTerminalMutationProbability float64,
	depth int, deletionStrategy int) (err error) {

	switch strategy {
	case StrategyDeleteNonTerminal: // CHANGE TO DeleteNonTerminal
		err = p.T.DeleteSubTree(deletionStrategy)
		break

	case StrategyDeleteMalicious:
		err = p.T.DeleteMalicious()
		break

	case StrategyFellTree:
		err = p.T.FellTree()
		break

	case StrategyDeleteTerminal:
		err = p.T.DeleteTerminal()
		break

	case StrategyMutateNonTerminal:
		err = p.T.MutateNonTerminal(nonTerminals)
		break

	case StrategyMutateTerminal:
		err = p.T.MutateTerminal(terminals)
		break

	case StrategyReplaceBranch:
		var tree *DualTree
		tree, err = GenerateRandomTreeEnforceIndependentVariable(depth, terminals[0], terminals, nonTerminals)
		err = p.T.ReplaceBranch(*tree)
		break

	case StrategyAddSubTree:
		var tree *DualTree
		tree, err = GenerateRandomTreeEnforceIndependentVariable(depth, terminals[0], terminals, nonTerminals)
		err = p.T.AddSubTree(tree)
		break

	case StrategyAddToLeaf:
		var tree *DualTree
		tree, err = GenerateRandomTreeEnforceIndependentVariable(depth, terminals[0], terminals, nonTerminals)
		err = p.T.AddToLeaf(*tree)
		break

	case StrategyAddMult:
		var tree *DualTree
		tree, err = GenerateRandomTreeEnforceIndependentVariable(depth, terminals[0], terminals,
			[]SymbolicExpression{{arity:2, value:"*", kind:1}})
		err = p.T.AddToLeaf(*tree)
		break

	case StrategyAddSub:
		var tree *DualTree
		tree, err = GenerateRandomTreeEnforceIndependentVariable(depth, terminals[0], terminals,
			[]SymbolicExpression{{arity:2, value:"-", kind:1}})
		err = p.T.AddToLeaf(*tree)
		break

	case StrategyAddAdd:
		var tree *DualTree
		tree, err = GenerateRandomTreeEnforceIndependentVariable(depth, terminals[0], terminals,
			[]SymbolicExpression{{arity:2, value:"+", kind:1}})
		err = p.T.AddToLeaf(*tree)
		break
	default:
		break
	}
	return err
}

// Eval is a simple helper function that takes in an independent variable,
// uses the programs treeNode to compute the resultant value
func (p *Program) EvalMulti(independentVariables IndependentVariableMap, expressionString string) (float64, error) {
	if p.T == nil {
		return -1, fmt.Errorf("program: %v -> treeNode is nil", p.ID)
	}

	return EvaluateMathematicalExpression(expressionString, independentVariables)
}

// EvaluateMathematicalExpression evaluates a valid expression using the given independentVar
func EvaluateMathematicalExpression(expressionString string, independentVariables IndependentVariableMap) (float64, error) {
	if expressionString == "" {
		return -1, fmt.Errorf("EvaluateMathematicalExpression | expressionString cannot be empty")
	}
	//
	//expression, err := gval.Evaluate(expressionString, independentVariables)
	//if err != nil {
	//	return -1, err
	//}


	expression, err :=  eval.CalculateWithVar(expressionString, independentVariables)
	if err != nil {
		return -1, err
	}

	ans, err := utils.ConvertToFloat64(expression)
	if err != nil {
		return -1, err
	}

	return ans, nil
}

func (p Program) Clone() (Program, error) {
	if p.T != nil {
		dualTree, err := p.T.Clone()
		if err != nil {
			return Program{}, err
		}
		p.T = &dualTree
	}
	p.ID = GenerateProgramID(0)
	return p, nil
}

func (p Program) CloneWithTree(tree DualTree) Program {
	p.T = &tree
	p.ID = GenerateProgramID(0)
	return p
}

type Bug *Program
type Test *Program
