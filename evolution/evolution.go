package evolution

import (
	"fmt"
)

type EvolutionParams struct {
	Generations                      int
	EnableParallelism                bool
	survivorSelection                int
	parentSelection                  int
	ElitismPercentage                float32
	ProgramEval                      func() float32
	MaxDepth                         int
	DepthPenaltyStrategy             int
	DepthPenaltyStrategyPenalization float32
	Threshold                        float64
	MinThreshold                     float64
	FitnessStrategy                  int
	TournamentSize                   int
	// EachPopulationSize represents the size of each protagonist or antagonist population.
	// This value must be even otherwise pairwise operations such as crossover will fail
	EachPopulationSize               int
	ProbabilityOfRecombination       float32
	ProbabilityOfMutation            float32
	ProbabilityOfNonTerminalMutation float32
	AntagonistMaxStrategies          int
	AntagonistStrategyLength         int
	ProtagonistMaxStrategies         int
	ProtagonistStrategyLength        int
	SurvivorPercentage               float32
	// Strategies is a list of available strategies for each individual.
	// These can be randomly allocated to individuals and duplicates are expected.
	Strategies            []Strategy
	DepthOfRandomNewTrees int
	// StrategyLengthPenalty is the penalty given to strategies that exceed a given length
	StrategyLengthPenalty float32
	// StrategyLengthLimit is the maximum length a strategy can reach before being penalized
	StrategyLengthLimit int

	// DeletionType pertains to the different kinds of deletion operations possible for a given tree.
	DeletionType int

	// EnforceIndependentVariable ensures that during individual generation at the start of the evolution,
	// independent variables are injected to the program meaning every program will at least have one independent
	// variable e.g. X
	EnforceIndependentVariable bool

	// CrossoverPercentage pertains to the amount of genetic material crossed-over.
	// This is a percentage represented as a float32. A value of 1 means all material is swapped.
	// A value of 0 means no material is swapped (which in effect are the same thing).
	// Avoid 0 or 1 use values in between
	CrossoverPercentage float32

	//MaintainGeneTransferEquality pertains to a scenario where in the event of differing strategy lengths betweween
	// two individuals, if N objects are to be swapped during crossover,
	// then N objects will be swapped from both individuals.
	// This prevents one individual receiving more genetic material from another individual and vice-versa
	MaintainCrossoverGeneTransferEquality bool

	TerminalSet    SymbolicExpressionSet
	NonTerminalSet SymbolicExpressionSet

	EvaluationThreshold    float64
	EvaluationMinThreshold float64
	ParentSelection int
}

const (
	DepthPenaltyStrategyIgnore   = 0
	DepthPenaltyStrategyPenalize = 2
	DepthPenaltyStrategyTrim     = 1
)

type EvolutionEngine struct {
	StartIndividual Program
	Spec            Spec

	Parallelize bool
	Generations []*Generation

	SurvivorSelection int
	ParentSelection   int
	ElitismPercentage float32
	StatisticsOutput  string

	Parameters EvolutionParams
}

func (e *EvolutionEngine) Start() (*EvolutionResult, error) {
	err := e.validate()
	if err != nil {
		return nil, err
	}

	// Init Population
	e.Generations = make([]*Generation, e.Parameters.Generations)
	// Set First Generation - TODO Parallelize Individual Creation
	antagonists, err := GenerateRandomIndividuals(e.Parameters.EachPopulationSize, "ANT", IndividualAntagonist,
		e.Parameters.AntagonistStrategyLength, e.Parameters.AntagonistMaxStrategies,
		e.Parameters.Strategies, 1, e.Parameters.TerminalSet, e.Parameters.NonTerminalSet, e.Parameters.EnforceIndependentVariable)
	if err != nil {
		return nil, err
	}
	protagonists, err := GenerateRandomIndividuals(e.Parameters.EachPopulationSize, "PRO", IndividualProtagonist,
		e.Parameters.ProtagonistStrategyLength, e.Parameters.ProtagonistMaxStrategies,
		e.Parameters.Strategies, 1, e.Parameters.TerminalSet, e.Parameters.NonTerminalSet, e.Parameters.EnforceIndependentVariable)
	if err != nil {
		return nil, err
	}

	//// create the 1st gen0, and begin
	genID := GenerateGenerationID(0)
	gen0 := Generation{
		count:        0,
		GenerationID: genID,
		Protagonists: protagonists,
		Antagonists:  antagonists,
		engine:       e,
	}
	e.Generations[0] = &gen0

	// cycle through generationCount
	for i := 0; i < e.Parameters.Generations; i++ {
		e.Generations[i], err = e.Generations[i].Start()
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// Todo Implement EvolutionProcess validate
func (e *EvolutionEngine) validate() error {
	if e.Parameters.Generations < 1 {
		return fmt.Errorf("set number of generationCount by calling e.Generations(x)")
	}
	if e.Parameters.EachPopulationSize % 2 != 0 {
		return fmt.Errorf("set number of EachPopulationSize to an Even number")
	}
	//if e.StartIndividual == Program{} {
	//	return fmt.Errorf("set a start individuals")
	//}
	//err := e.StartIndividual.Validate()
	//if err != nil {
	//	return err
	//}
	//if e.spec == nil {
	//	return fmt.Errorf("set a valid spec")
	//}
	//if len(e.spec) < 3 {
	//	return fmt.Errorf("a small spec will hamper evolutionary accuracy")
	//}
	return nil
}