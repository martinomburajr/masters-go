package evolution

import (
	"fmt"
)

type EvolutionParams struct {
	Generations                      int
	EnableParallelism                bool
	survivorSelection                int
	parentSelection                  int
	ElitismPercentage                float64
	ProgramEval                      func() float64
	MaxDepth                         int
	DepthPenaltyStrategy             int
	DepthPenaltyStrategyPenalization float64
	Threshold                        float64
	MinThreshold                     float64
	FitnessStrategy                  int
	TournamentSize                   int
	// EachPopulationSize represents the size of each protagonist or antagonist population.
	// This value must be even otherwise pairwise operations such as crossover will fail
	EachPopulationSize               int
	ProbabilityOfRecombination       float64
	ProbabilityOfMutation            float64
	ProbabilityOfNonTerminalMutation float64

	// ThresholdMultiplier is used when the FitnessRatioThreshold option is selected.
	// It creates a threshold value based on the cumulative value of the dependent variable in the spec.
	// It cannot be less than 1 as a value less than 1 would mean that the approximater functions would need to be
	// better than the spec and that is not possible
	ThresholdMultiplier float64

	AntagonistMaxStrategies int
	//AntagonistStrategyLength         int
	ProtagonistMaxStrategies int
	//ProtagonistStrategyLength        int
	SurvivorPercentage float64
	// Strategies is a list of available strategies for each individual.
	// These can be randomly allocated to individuals and duplicates are expected.
	//Strategies            []Strategy
	DepthOfRandomNewTrees int
	// StrategyLengthPenalty is the penalty given to strategies that exceed a given length
	StrategyLengthPenalty float64
	// StrategyLengthLimit is the maximum length a Strategy can reach before being penalized
	StrategyLengthLimit int

	// DeletionType pertains to the different kinds of deletion operations possible for a given Tree.
	DeletionType int

	// EnforceIndependentVariable ensures that during individual Generation at the start of the evolution,
	// independent variables are injected to the program meaning every program will at least have one independent
	// variable e.g. X
	EnforceIndependentVariable bool

	// CrossoverPercentage pertains to the amount of genetic material crossed-over.
	// This is a percentage represented as a float64. A value of 1 means all material is swapped.
	// A value of 0 means no material is swapped (which in effect are the same thing).
	// Avoid 0 or 1 use values in between
	CrossoverPercentage float64

	//MaintainGeneTransferEquality pertains to a scenario where in the event of differing Strategy lengths betweween
	// two individuals, if N objects are to be swapped during crossover,
	// then N objects will be swapped from both individuals.
	// This prevents one individual receiving more genetic material from another individual and vice-versa
	MaintainCrossoverGeneTransferEquality bool

	TerminalSet    SymbolicExpressionSet
	NonTerminalSet SymbolicExpressionSet

	EvaluationThreshold float64
	ParentSelection     int
	StartIndividual     Program
	Spec                SpecMulti
	SurvivorSelection   int

	ProtagonistAvailableStrategies []Strategy
	AntagonistAvailableStrategies  []Strategy
	SetEqualStrategyLength         bool
	EqualStrategiesLength          int

	// VariableTerminals represent all the potential variables that may appear.
	// An effort is made to differentiate them from constants so that constants do not get overwritten as variables
	// would when calculating the spec.
	VariableTerminals              []SymbolicExpression
}

const (
	DepthPenaltyStrategyIgnore   = 0
	DepthPenaltyStrategyPenalize = 2
	DepthPenaltyStrategyTrim     = 1
)

type EvolutionEngine struct {
	Parallelize      bool
	Generations      []*Generation
	StatisticsOutput string
	Parameters       EvolutionParams
}

func (e *EvolutionEngine) Start() (EvolutionResult, error) {
	err := e.validate()
	if err != nil {
		return EvolutionResult{}, err
	}

	// Set First Generation - TODO Parallelize Individual Creation
	genID := GenerateGenerationID(0)
	gen0 := Generation{
		count:        0,
		GenerationID: genID,
		Protagonists: nil,
		Antagonists:  nil,
		engine:       e,
	}
	e.Generations[0] = &gen0

	antagonists, err := e.Generations[0].GenerateRandomIndividual("ANT", e.Parameters.StartIndividual)
	if err != nil {
		return EvolutionResult{}, err
	}

	protagonists, err := e.Generations[0].GenerateRandomIndividual("PRO", e.Parameters.StartIndividual)
	if err != nil {
		return EvolutionResult{}, err
	}

	gen0.Protagonists = protagonists
	gen0.Antagonists = antagonists

	// cycle through generationCount
	e.Generations[0] = &gen0
	for i := 0; i < e.Parameters.Generations-1; i++ {
		//if i != e.Parameters.Generations-2 {
		protagonistsCleanse, err := CleansePopulation(e.Generations[i].Protagonists, *e.Parameters.StartIndividual.T)
		if err != nil {
			return EvolutionResult{}, err
		}
		antagonistsCleanse, err := CleansePopulation(e.Generations[i].Antagonists, *e.Parameters.StartIndividual.T)
		if err != nil {
			return EvolutionResult{}, err
		}

		e.Generations[i].Protagonists = protagonistsCleanse
		e.Generations[i].Antagonists = antagonistsCleanse
		//}
		nextGeneration, err := e.Generations[i].Start()
		if err != nil {
			return EvolutionResult{}, err
		}
		e.Generations[i+1] = nextGeneration
	}

	// Sort individuals in all generations
	for i := range e.Generations {
		e.Generations[i].Protagonists = SortIndividuals(e.Generations[i].Protagonists)
		e.Generations[i].Antagonists = SortIndividuals(e.Generations[i].Antagonists)
	}

	evolutionResult := EvolutionResult{}
	_, err = evolutionResult.Analyze(e.Generations, 3)
	if err != nil {
		return EvolutionResult{}, err
	}

	return evolutionResult, nil
}

// Todo Implement EvolutionProcess validate
func (e *EvolutionEngine) validate() error {
	if e.Parameters.Generations < 1 {
		return fmt.Errorf("set number of generationCount by calling e.Generations(x)")
	}
	if e.Parameters.EachPopulationSize%2 != 0 {
		return fmt.Errorf("set number of EachPopulationSize to an Even number")
	}
	if e.Parameters.SetEqualStrategyLength == true && e.Parameters.EqualStrategiesLength < 1 {
		return fmt.Errorf("cannot SetEqualStrategyLength to true and EqualStrategiesLength less than 1")
	}
	if e.Parameters.StartIndividual.T == nil {
		return fmt.Errorf("start individual cannot have a nil Tree")
	}
	if e.Parameters.Spec == nil {
		return fmt.Errorf("spec cannot be nil")
	}
	if len(e.Parameters.Spec) < 1 {
		return fmt.Errorf("spec cannot be empty")
	}
	if e.Parameters.FitnessStrategy == FitnessRatioThresholder && e.Parameters.ThresholdMultiplier < 1 {
		return fmt.Errorf("ThresholdMultiplier cannot be less than 1")
	}
	//err := e.StartIndividual.Validate()
	//if err != nil {
	//	return err
	//}

	if len(e.Parameters.Spec) < 3 {
		return fmt.Errorf("a small spec will hamper evolutionary accuracy")
	}
	return nil
}
