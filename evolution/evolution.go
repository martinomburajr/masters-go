package evolution

import (
	"fmt"
)

type EvolutionParams struct {
	Generations          int
	EnableParallelism    bool
	survivorSelection    int
	parentSelection      int
	ElitismPercentage    float32
	ProgramEval          func() float32
	StatisticsOutput     string
	MaxDepth             int
	DepthPenaltyStrategy int
	DepthPenaltyStrategyPenalization float32
	Threshold            float64
	MinThreshold         float64
	FitnessStrategy      int
	TournamentSize       int
	EachPopulationSize int
	ProbabilityOfRecombination       float32
	ProbabilityOfMutation            float32
	ProbabilityOfNonTerminalMutation float32
	AntagonistMaxStrategies int
	AntagonistStrategyLength int
	ProtagonistMaxStrategies int
	ProtagonistStrategyLength int
	SurvivorPercentage float32
}

const (
	DepthPenaltyStrategyIgnore = 0
	DepthPenaltyStrategyPenalize = 2
	DepthPenaltyStrategyTrim = 1
)

type EvolutionEngine struct {
	StartIndividual                  Program
	Spec                             Spec
	GenerationCount                  int
	Parallelize                      bool
	Generations                      []*Generation
	AvailableStrategies              []Strategy
	AvailableTerminalSet             SymbolicExpressionSet
	AvailableNonTerminalSet          SymbolicExpressionSet
	SurvivorSelection                int
	ParentSelection                  int
	ElitismPercentage                float32
	StatisticsOutput                 string
	MaxDepth                         int
	DepthPenalty                     int
	EvaluationThreshold              float64
	EvaluationMinThreshold           float64
	FitnessStrategy                  int
	TournamentSize                   int
	Parameters EvolutionParams
}



func (e *EvolutionEngine) Start() (*EvolutionResult, error) {
	err := e.validate()
	if err != nil {
		return nil, err
	}

	// Init Population
	e.Generations = make([]*Generation, e.GenerationCount)
	// Set First Generation - TODO Parallelize Individual Creation
	antagonists, err := GenerateRandomIndividuals(e.Parameters.EachPopulationSize, "ANT", IndividualAntagonist,
		e.Parameters.AntagonistStrategyLength, e.Parameters.AntagonistMaxStrategies,
		e.AvailableStrategies, 1, e.AvailableTerminalSet, e.AvailableNonTerminalSet)
	if err != nil {
		return nil, err
	}
	protagonists, err := GenerateRandomIndividuals(e.Parameters.EachPopulationSize, "PRO", IndividualProtagonist,
		e.Parameters.ProtagonistStrategyLength, e.Parameters.ProtagonistMaxStrategies,
		e.AvailableStrategies, 1, e.AvailableTerminalSet, e.AvailableNonTerminalSet)
	if err != nil {
		return nil, err
	}

	//// create the 1st gen0, and begin
	genID := GenerateGenerationID(0)
	gen0 := Generation{
		count:          0,
		GenerationID:   genID,
		Protagonists:   protagonists,
		Antagonists:    antagonists,
		engine: e,
	}
	e.Generations[0] = &gen0

	// cycle through generationCount
	for i := 0; i < e.GenerationCount; i++ {
		e.Generations[i], err = e.Generations[i].Start()
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// Todo Implement EvolutionProcess validate
func (e *EvolutionEngine) validate() error {
	if e.GenerationCount < 1 {
		return fmt.Errorf("set number of generationCount by calling e.Generations(x)")
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


//
//// InitialIndividual returns the input individual
//func (e *EvolutionEngine) GetInitialIndividual() *InitialProgram {
//	return e.StartIndividual
//}
//
//// Options Sets Options to be used in the Evolutionary Process
//func (e *EvolutionEngine) Options(params EvolutionParams) *EvolutionEngine {
//	return e
//}
//
//
//// FitnessEval is a function provided that gives the engine and individuals a means to calculate fitness.
//func (e *EvolutionEngine) FitnessEval(fitnessFunc func() float32) *EvolutionEngine {
//	return e
//}
//
//// ProgramEval
//func (e *EvolutionEngine) ProgramEval(programFunc func() float32) *EvolutionEngine {
//	return e
//}
//
//// Protagonist sets the protagonists count as well as defines a fitness function that is used to calculate its
//// fitness. If you are using sharedFitness,
//// set fitnessFunc to nil. The protagonist is also initialized with a set of strategies it can use.
//// If nil it will pull from a list of available strategies
//func (e *EvolutionEngine) Protagonist(count int, fitnessFunc func() float32, strategies []Strategable) *EvolutionEngine {
//	return e
//}
//

//
//// AvailableStrategies represents a list of strategies available to the population
//func (e *EvolutionEngine) AvailableStrategies(strategies []Strategable) *EvolutionEngine {
//	return e
//}
//
//// Generations indicates the maximum number of generationCount before the simulation ends.
//func (e *EvolutionEngine) Generations(i int) *EvolutionEngine {
//	return e
//}
//
//// GenerationsByError uses maxError to determine how much to minimize the solutions error by before terminating the
//// evolutionary process
//func (e *EvolutionEngine) GenerationsByError(maxError float32) *EvolutionEngine {
//	return e
//}
//
//func (e *EvolutionEngine) ParentSelection(b bool) *EvolutionEngine {
//	return e
//}
//
//func (e *EvolutionEngine) SurvivorSelection(b bool) *EvolutionEngine {
//	return e
//}
//
//func (e *EvolutionEngine) OptimizationStrategy(b bool) *EvolutionEngine {
//	return e
//}
//
//func (e *EvolutionEngine) Parallelize(b bool) *EvolutionEngine {
//	return e
//}
//
//// GenerateStatistics will output statistics to a given file
//func (e *EvolutionEngine) GenerateStatistics(s string) *EvolutionEngine {
//	return e
//}
//
//// Start begines the evolutionary engine, and starts the evolutionary process returning an EvolutionaryProcess
//func (e *EvolutionEngine) Start() *EvolutionProcess {
//	e.validate()
//	return nil
//}
//
//func (e *EvolutionEngine) SetProbabilityOfMutation(probabilityOfMutation float32) {
//	e.probabilityOfMutation = probabilityOfMutation
//}
//
//func (e *EvolutionEngine) SetProbabilityOfNonTerminalMutation(probabilityOfNonTerminalMutation float32) {
//	e.probabilityOfNonTerminalMutation = probabilityOfNonTerminalMutation
//}
//
//// ZeroSumFitness is a measure where both protagonist and antagonist compete from a shared fitness pool.
//// The more one side gets the less the other gets. If this strategy is chosen,
//// you cannot set different fitness strategies for the protagonist and antagonists
//func (e *EvolutionEngine) ZeroSumFitness(i func() float32) *EvolutionEngine {
//	return nil
//}
//
//type InitialProgram struct {
//	ID   string
//	T    *DualTree
//	spec Spec
//}
