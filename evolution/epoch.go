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
	protagonist       *Individual
	antagonist        *Individual
	program           Program
	protagonistBegins bool
	isComplete        bool
	probabilityOfMutation float32
	probabilityOfNonTerminalMutation float32
	terminalSet []SymbolicExpression
	nonTerminalSet []SymbolicExpression
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
func (e *Epoch) SetProbabilityOfMutation(probability float32) *Epoch {
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

// InitSimulator creates the Epoch Simulator. You must call Start to begin this process
func (e *Epoch) InitSimulator() *EpochSimulator {
	return &EpochSimulator {
		e, false, false,
	}
}

// EpochSimulator is responsible for simulating actions in a given Epoch
type EpochSimulator struct {
	epoch                 *Epoch
	hasAntagonistApplied  bool
	hasProtagonistApplied bool
}

// applyAntagonistStrategy applies the Antagonist strategies to program.
func (e *EpochSimulator) applyAntagonistStrategy() (*EpochSimulator, error) {
	for _, strategy := range e.epoch.antagonist.strategy {
		err := e.epoch.program.ApplyStrategy(strategy,
			e.epoch.terminalSet,
			e.epoch.nonTerminalSet,
			e.epoch.probabilityOfMutation,
			e.epoch.probabilityOfNonTerminalMutation)
		if err != nil {
			return nil, err
		}
	}
	e.hasAntagonistApplied = true
	return e, nil
}

// applyProtagonistStrategy Apply Protagonist strategies to program.
func (e *EpochSimulator) applyProtagonistStrategy() (*EpochSimulator, error) {
	for _, strategy := range e.epoch.protagonist.strategy {
		err := e.epoch.program.ApplyStrategy(strategy,
			e.epoch.terminalSet,
			e.epoch.nonTerminalSet,
			e.epoch.probabilityOfMutation,
			e.epoch.probabilityOfNonTerminalMutation)
		if err != nil {
			return nil, err
		}
	}
	e.hasProtagonistApplied = true
	return e, nil
}

// Start begins the epoch simulation by allowing the competing individuals to compete
func (e *EpochSimulator) Start() (*EpochResult, error) {
	if e.epoch.protagonist == nil {
		return nil, fmt.Errorf("epoch cannot have nil protagonist")
	}
	if e.epoch.antagonist == nil {
		return nil, fmt.Errorf("epoch cannot have nil antagonist")
	}

	antagonistEpoch, err := e.applyAntagonistStrategy()
	if err != nil {
		return nil, err
	}

	combinedEpoch, err := antagonistEpoch.applyProtagonistStrategy()
	if err != nil {
		return nil, err
	}

	if combinedEpoch.hasProtagonistApplied && combinedEpoch.hasAntagonistApplied {
		return nil, nil
	}

	return &EpochResult{
		combinedEpoch.epoch,
	}, nil

}

type EpochResult struct {
	epoch *Epoch
}
