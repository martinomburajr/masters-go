package evolution

import (
	"fmt"
)

type EvolutionParams struct {
	Generations       int
	EnableParallelism bool
}

type EvolutionEngine struct {
	startIndividual     *InitialProgram
	spec                Spec
	generations         int
	parallelize         bool
	availableStrategies []*Strategable
	programEval         func() float32
	statisticsOutput    string
}

// Todo Implement EvolutionProcess validate
func (e *EvolutionEngine) validate() error {
	if e.generations == 0 {
		return fmt.Errorf("set number of generations by calling e.Generations(x)")
	}
	if e.startIndividual == nil {
		return fmt.Errorf("set a start generation")
	}
	err := e.startIndividual.Validate()
	if err != nil {
		return err
	}
	if e.spec == nil {
		return fmt.Errorf("set a valid spec")
	}
	if len(e.spec) < 3 {
		return fmt.Errorf("a small spec will hamper evolutionary accuracy")
	}

	return nil
}

// InitialIndividual returns the input individual
func (e *EvolutionEngine) GetInitialIndividual() *InitialProgram {
	return e.startIndividual
}

// Options Sets Options to be used in the Evolutionary Process
func (e *EvolutionEngine) Options(params EvolutionParams) *EvolutionEngine {
	return e
}

// SetStartIndividual sets the starting individual along with the spec. Both must be provided
func (e *EvolutionEngine) SetStartIndividual(program InitialProgram) *EvolutionEngine {
	e.spec = program.GetSpec()
	return e
}

// FitnessEval is a function provided that gives the engine and individuals a means to calculate fitness.
func (e *EvolutionEngine) FitnessEval(fitnessFunc func() float32) *EvolutionEngine {
	return e
}

// ProgramEval
func (e *EvolutionEngine) ProgramEval(programFunc func() float32) *EvolutionEngine {
	return e
}

// Protagonist sets the protagonists count as well as defines a fitness function that is used to calculate its
// fitness. If you are using sharedFitness,
// set fitnessFunc to nil. The protagonist is also initialized with a set of strategies it can use.
// If nil it will pull from a list of available strategies
func (e *EvolutionEngine) Protagonist(count int, fitnessFunc func() float32, strategies []Strategable) *EvolutionEngine {
	return e
}

// antagonist sets the antagonists count as well as defines a fitness function that is used to calculate its fitness.
// If you are using sharedFitness, set fitnessFunc to nil.
// The antagonist is also initialized with a set of strategies it can use.
// If nil it will pull from a list of available strategies
func (e *EvolutionEngine) Antagonist(count int, fitnessFunc func() float32, strategies []Strategable) *EvolutionEngine {
	return e
}

// AvailableStrategies represents a list of strategies available to the population
func (e *EvolutionEngine) AvailableStrategies(strategies []Strategable) *EvolutionEngine {
	return e
}

// Generations indicates the maximum number of generations before the simulation ends.
func (e *EvolutionEngine) Generations(i int) *EvolutionEngine {
	return e
}

// GenerationsByError uses maxError to determine how much to minimize the solutions error by before terminating the
// evolutionary process
func (e *EvolutionEngine) GenerationsByError(maxError float32) *EvolutionEngine {
	return e
}

func (e *EvolutionEngine) ParentSelection(b bool) *EvolutionEngine {
	return e
}

func (e *EvolutionEngine) SurvivorSelection(b bool) *EvolutionEngine {
	return e
}

func (e *EvolutionEngine) OptimizationStrategy(b bool) *EvolutionEngine {
	return e
}

func (e *EvolutionEngine) Parallelize(b bool) *EvolutionEngine {
	return e
}

// GenerateStatistics will output statistics to a given file
func (e *EvolutionEngine) GenerateStatistics(s string) *EvolutionEngine {
	return e
}

// Start begines the evolutionary engine, and starts the evolutionary process returning an EvolutionaryProcess
func (e *EvolutionEngine) Start() *EvolutionProcess {
	e.validate()
	return nil
}

// ZeroSumFitness is a measure where both protagonist and antagonist compete from a shared fitness pool.
// The more one side gets the less the other gets. If this strategy is chosen,
// you cannot set different fitness strategies for the protagonist and antagonists
func (e *EvolutionEngine) ZeroSumFitness(i func() float32) *EvolutionEngine {
	return nil
}

type InitialProgram struct {
	ID   string
	T    *DualTree
	spec Spec
}
