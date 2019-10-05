package evolution

import (
	"fmt"
	"log"
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
	//AntagonistStrategyLength         int
	ProtagonistMaxStrategies int
	//ProtagonistStrategyLength        int
	SurvivorPercentage float32

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
	ParentSelection        int
	StartIndividual        Program
	Spec                   Spec
	SurvivorSelection      int

	ProtagonistAvailableStrategies []Strategy
	AntagonistAvailableStrategies  []Strategy
	SetEqualStrategyLength         bool
	EqualStrategiesLength          int
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

	e.Generations = make([]*Generation, e.Parameters.Generations)

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

	antagonists, err := e.Generations[0].GenerateRandomAntagonists("ANT")
	if err != nil {
		return EvolutionResult{}, err
	}

	protagonists, err := e.Generations[0].GenerateRandomProtagonists("PRO")
	if err != nil {
		return EvolutionResult{}, err
	}

	gen0.Protagonists = protagonists
	gen0.Antagonists = antagonists

	// cycle through generationCount
	e.Generations[0] = &gen0
	for i := 0; i < e.Parameters.Generations-1; i++ {
		nextGeneration, err := e.Generations[i].Start()
		if err != nil {
			return EvolutionResult{}, err
		}
		e.Generations[i+1] = nextGeneration
	}

	evolutionResult := EvolutionResult{}
	_, err = evolutionResult.Analyze(e.Generations, 3)
	if err != nil {
		return EvolutionResult{}, err
	}
	fmt.Println("Top Protagonist Tree")
	topProtagonistTree := evolutionResult.TopProtagonist.tree
	log.Println(topProtagonistTree)
	log.Printf("%#v", evolutionResult.TopAntagonist.result.strategy)

	topAntagonistTree := evolutionResult.TopAntagonist.tree
	fmt.Println("Top Antagonist Tree")
	log.Println(topAntagonistTree)
	log.Printf("%#v", evolutionResult.TopAntagonist.result.strategy)
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
	if e.Parameters.StartIndividual.T.root == nil {
		return fmt.Errorf("start individual cannot have a nil tree root")
	}
	if e.Parameters.Spec == nil {
		return fmt.Errorf("spec cannot be nil")
	}
	if len(e.Parameters.Spec) < 1 {
		return fmt.Errorf("spec cannot be empty")
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
