package evolution

import (
	"fmt"
)

// Epoch is defined as a coevolutionary step where protagonist and antagonist compete.
// For example an epoch could represent a distinct interaction between two parties.
// For instance a bug mutated program (antagonist) can be challenged a variety of times (
// specified by {iterations}) by the tests (protagonist).
// The test will use up the strategies it contains and attempt to chew away at the antagonists fitness,
// to maximize its own
type Epoch struct {
	id                               string
	protagonist                      *Individual
	antagonist                       *Individual
	generation  *Generation
	program                          Program
	protagonistBegins                bool
	isComplete                       bool
	probabilityOfMutation            float32
	probabilityOfNonTerminalMutation float32
	terminalSet                      []SymbolicExpression
	nonTerminalSet                   []SymbolicExpression
	hasAntagonistApplied  bool
	hasProtagonistApplied bool
}

// NewEpoch creates a new epoch. The id string can simply be the index from an iteration that creates multiple epochs
func NewEpoch(id string, protagonist *Individual, antagonist *Individual, program Program, probabilityOfMutation float32, probabilityOfNonTerminalMutation float32, terminalSet []SymbolicExpression, nonTerminalSet []SymbolicExpression) *Epoch {
	id = fmt.Sprintf("Epoch-%s-%s|%s", id, antagonist.id, protagonist.id)
	return &Epoch{id: id, protagonist: protagonist, antagonist: antagonist, program: program, probabilityOfMutation: probabilityOfMutation, probabilityOfNonTerminalMutation: probabilityOfNonTerminalMutation, terminalSet: terminalSet, nonTerminalSet: nonTerminalSet}
}

func (e *Epoch) GetProtagonistBegins() bool {
	return e.protagonistBegins
}

// Program sets the program for the epoch
func (e *Epoch) SetProgram(program Program) *Epoch {
	e.program = program
	return e
}

// ProtagonistBegins states whether the protagonist should start the epoch
func (e *Epoch) SetProtagonistBegins(protagonistBegins bool) *Epoch {
	e.protagonistBegins = protagonistBegins
	return e
}

// Protagonist sets the protagonist for the epoch
func (e *Epoch) SetProtagonist(protagonist *Individual) *Epoch {
	e.protagonist = protagonist
	return e
}

// Antagonist sets the antagonist for the epoch
func (e *Epoch) SetAntagonist(antagonist *Individual) *Epoch {
	e.antagonist = antagonist
	return e
}

// SetProbabilityOfMutation sets the probability that the program will use a mutation strategy.
// Otherwise it will be skipped
func (e Epoch) SetProbabilityOfMutation(probability float32) Epoch {
	e.probabilityOfMutation = probability
	e.probabilityOfNonTerminalMutation = probability
	return e
}

// SetProbabilityOfNonTerminalMutation sets the probability that the program will mutate the non-terminal after
// mutation is deemed as the appropriate strategy. Otherwise it will mutate the terminal instead.
func (e *Epoch) SetProbabilityOfNonTerminalMutation(probability float32) *Epoch {
	e.probabilityOfNonTerminalMutation = probability
	return e
}

// Start creates the Epoch process. This process applies the antagonist strategy first,
// and then the protagonist strategy second.
// It then appends the fitness values to each individual in the epoch.
func (e *Epoch) Start() error {
	if e.protagonist == nil {
		return fmt.Errorf("epoch cannot have nil protagonist")
	}
	if e.antagonist == nil {
		return fmt.Errorf("epoch cannot have nil antagonist")
	}

	err := e.applyAntagonistStrategy()
	if err != nil {
		return err
	}

	err = e.applyProtagonistStrategy()
	if err != nil {
		return err
	}

	if e.hasProtagonistApplied && e.hasAntagonistApplied {
		return fmt.Errorf("antagonist and protagonist havent applied strategy to program")
	}

	antagonistFitness, protagonistFitness := 0, 0
	switch e.generation.engine.parentSelection {
		case FitnessProtagonistThresholdTally:
			antagonistFitness, protagonistFitness, err = ProtagonistThresholdTally(e.generation.engine.spec,
				&e.program, e.generation.engine.threshold,
				e.generation.engine.minThreshold)
				if err != nil {
					return err
				}
	}

	e.antagonist.fitness = append(e.antagonist.fitness, antagonistFitness)
	e.protagonist.fitness = append(e.protagonist.fitness, protagonistFitness)

	return nil
}

// applyAntagonistStrategy applies the Antagonist strategies to program.
func (e *Epoch) applyAntagonistStrategy() (error) {
	for _, strategy := range e.antagonist.strategy {
		err := e.program.ApplyStrategy(strategy,
			e.terminalSet,
			e.nonTerminalSet,
			e.probabilityOfMutation,
			e.probabilityOfNonTerminalMutation,
			e.generation.engine.maxDepth)
		if err != nil {
			return err
		}
	}
	e.hasAntagonistApplied = true
	return nil
}

// applyProtagonistStrategy Apply Protagonist strategies to program.
func (e *Epoch) applyProtagonistStrategy() (error) {
	for _, strategy := range e.protagonist.strategy {
		err := e.program.ApplyStrategy(strategy,
			e.terminalSet,
			e.nonTerminalSet,
			e.probabilityOfMutation,
			e.probabilityOfNonTerminalMutation,
			e.generation.engine.maxDepth)
		if err != nil {
			return err
		}
	}
	e.hasProtagonistApplied = true
	return nil
}
