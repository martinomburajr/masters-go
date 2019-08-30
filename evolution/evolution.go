package evolution

import (
	"fmt"
)

type EvolutionParams struct {
	Generations       int
	EnableParallelism bool
}

type EvolutionEngine struct {
	startIndividual                  Program
	spec                             Spec
	generationResults                []*GenerationResult
	generations                      int
	eachPopulation                   int
	parallelize                      bool
	probabilityOfRecombination       float32
	probabilityOfMutation            float32
	probabilityOfNonTerminalMutation float32
	availableStrategies              []Strategy
	survivorSelection                int
	parentSelection                  int
	elitismPercentage                float32
	programEval                      func() float32
	statisticsOutput                 string
	maxDepth int
	depthPenalty int
	threshold float64
	minThreshold float64
	fitnessStrategy int
}

func (engine *EvolutionEngine) StatisticsOutput() string {
	return engine.statisticsOutput
}

func (engine *EvolutionEngine) SetStatisticsOutput(statisticsOutput string) {
	engine.statisticsOutput = statisticsOutput
}

// ProgramEval is a function provided that gives the engine and individuals a means to calculate fitness.
func (engine *EvolutionEngine) ProgramEval() func() float32 {
	return engine.programEval
}

// ProgramEval is a function provided that gives the engine and individuals a means to calculate fitness.
func (engine *EvolutionEngine) SetProgramEval(programEval func() float32) {
	engine.programEval = programEval
}

func (engine *EvolutionEngine) ProbabilityOfNonTerminalMutation() float32 {
	return engine.probabilityOfNonTerminalMutation
}

func (engine *EvolutionEngine) SetProbabilityOfNonTerminalMutation(probabilityOfNonTerminalMutation float32) {
	engine.probabilityOfNonTerminalMutation = probabilityOfNonTerminalMutation
}

func (engine *EvolutionEngine) ProbabilityOfMutation() float32 {
	return engine.probabilityOfMutation
}

func (engine *EvolutionEngine) SetProbabilityOfMutation(probabilityOfMutation float32) {
	engine.probabilityOfMutation = probabilityOfMutation
}

func (engine *EvolutionEngine) Generations() int {
	return engine.generations
}

func (engine *EvolutionEngine) SetGenerations(generations int) {
	engine.generations = generations
}

func (engine *EvolutionEngine) Spec() Spec {
	return engine.spec
}

func (engine *EvolutionEngine) SetSpec(spec Spec) {
	engine.spec = spec
}

func (engine *EvolutionEngine) StartIndividual() Program {
	return engine.startIndividual
}

// SetStartIndividual sets the starting individual along with the spec. Both must be provided
func (engine *EvolutionEngine) SetStartIndividual(startIndividual Program) {
	engine.startIndividual = startIndividual
}

func (engine *EvolutionEngine) AvailableStrategies() []Strategy {
	return engine.availableStrategies
}

func (engine *EvolutionEngine) SetAvailableStrategies(availableStrategies []Strategy) {
	engine.availableStrategies = availableStrategies
}

func (engine *EvolutionEngine) Parallelize() bool {
	return engine.parallelize
}

func (engine *EvolutionEngine) SetParallelize(parallelize bool) {
	engine.parallelize = parallelize
}

func (e *EvolutionEngine) Start() (*EvolutionResult, error) {
	// Init Population
	//err := e.validate()
	//if err != nil {
	//	return nil, err
	//}
	//// Set First Generation - TODO Parallelize Individual Creation
	//antagonists, err := GenerateRandomIndividuals(e.eachPopulation, "ANT", "BUG", 3, 4, e.availableStrategies)
	//if err != nil {
	//	return nil, err
	//}
	//protagonists, err := GenerateRandomIndividuals(e.eachPopulation, "PRO", "TEST", 3, 4, e.availableStrategies)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// create the 1st gen0, and begin
	//gen0 := Generation{}
	//generationResult, err := gen0.Start()

	// cycle through generations
	//for i := 0; i < e.generations; i++ {
	//	generationId := fmt.Sprintf("GEN#-%d", i)
	//	Generation{}.Start()
	//}
	return nil, nil
}

// Todo Implement EvolutionProcess validate
func (e *EvolutionEngine) validate() error {
	if e.generations < 1 {
		return fmt.Errorf("set number of generations by calling e.Generations(x)")
	}
	//if e.startIndividual == Program{} {
	//	return fmt.Errorf("set a start individuals")
	//}
	//err := e.startIndividual.Validate()
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
//	return e.startIndividual
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
//// Generations indicates the maximum number of generations before the simulation ends.
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
