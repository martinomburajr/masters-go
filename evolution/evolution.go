package evolution

import (
	"fmt"
	"github.com/martinomburajr/masters-go/program"

	//"github.com/martinomburajr/masters-go/program"
	"github.com/martinomburajr/masters-go/program/tree/dualtree"
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

// Antagonist sets the antagonists count as well as defines a fitness function that is used to calculate its fitness.
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
	return e
}

// ZeroSumFitness is a measure where both protagonist and antagonist compete from a shared fitness pool.
// The more one side gets the less the other gets. If this strategy is chosen,
// you cannot set different fitness strategies for the protagonist and antagonists
func (e *EvolutionEngine) ZeroSumFitness(i func() float32) *EvolutionEngine {

}

/**
EvolutionProcess represents the state of an evolutionary process once the evolution engine starts
*/
type EvolutionProcess struct {
	currentGeneration *Generation
	engine            *EvolutionEngine
}

type Generation struct {
	GenerationID       int
	PreviousGeneration *Generation
	NextGeneration     *Generation
	Protagonists       []*program.Program //Protagonists in a given generation
	Antagonists        []*program.Program //Antagonists in a given generation
	FittestProtagonist *program.Program
	FittestAntagonist  *program.Program
	engine             *EvolutionEngine // Reference to Engine
}

// Next returns the next generation
func (g *Generation) Next() *Generation {
	return g.NextGeneration
}

// Engine returns a reference to the Evolution Engine in use
func (g *Generation) Engine() *EvolutionEngine {
	return g.Engine()
}

// Previous returns the previous generation
func (g *Generation) Previous() *Generation {
	return g.PreviousGeneration
}

// Start begins the generational evolutionary cycle.
// It creates a new generation that it links the {NextGeneration} field to. Similar to the way a LinkedList works
func (g *Generation) Start() *Generation {
	return g.PreviousGeneration
}

// Restart is similar to StartHOG but it restarts the evolutionary process from the selected Generation.
// All future generations are deleted to make way for this evolutionary process
func (g *Generation) Restart() *Generation {
	return g.PreviousGeneration
}

// StartHOG is a unique version of start. It clears future history and jumps to a given generation,
// inserts generational material into the generation, and creates a new evolutionary propagation from it.
func (g *Generation) StartHOG(gen Generation) *Generation {
	return g.PreviousGeneration
}

// Epoch is defined as a coevolutionary step where protagonist and antagonist compete.
// For example an epoch could represent a distinct interaction between two parties.
// For instance a bug mutated program (antagonist) can be challenged a variety of times (
// specified by {iterations}) by the tests (protagonist).
// The test will use up the strategies it contains and attempt to chew away at the antagonists fitness,
// to maximize its own
type Epoch struct {
	Protagonist *program.Program
	Antagonist  *program.Program
	engine      *EvolutionEngine
	iterations  int
	isComplete  bool
	generation  *Generation
}

func (e *Epoch) Iterations(iterations int) *Epoch {
	e.iterations = iterations
	return e
}

// EpochSimulator is responsible for simulating actions in a given Epoch
type EpochSimulator struct {
	epoch *Epoch
}

// Start begins the epoch simulation by allowing the competing individuals to do their thing
func (e *Epoch) Start() *EpochResult {
	return nil
}

type EpochResult struct {
	engine *EvolutionEngine //Reference to underlying engine
}

type InitialProgram struct {
	ID   string
	T    *dualtree.DualTree
	spec Spec
}

func (p *InitialProgram) Spec(spec Spec) *InitialProgram {
	p.spec = spec
	return p
}

func (p *InitialProgram) GetSpec() Spec {
	return p.spec
}

func (p *InitialProgram) Validate() error {

}