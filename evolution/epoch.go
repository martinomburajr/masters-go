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
	program           *Program
	protagonistBegins bool
	isComplete        bool
}

func (e *Epoch) GetProtagonistBegins() bool {
	return e.protagonistBegins
}

// Program sets the program for the epoch
func (e *Epoch) Program(program *Program) *Epoch {
	e.program = program
	return e
}

// ProtagonistBegins states whether the protagonist should start the epoch
func (e *Epoch) ProtagonistBegins(protagonistBegins bool) *Epoch {
	e.protagonistBegins = protagonistBegins
	return e
}

// Protagonist sets the protagonist for the epoch
func (e *Epoch) Protagonist(protagonist *Individual) *Epoch {
	e.protagonist = protagonist
	return e
}

// Antagonist sets the antagonist for the epoch
func (e *Epoch) Antagonist(antagonist *Individual) *Epoch {
	e.antagonist = antagonist
	return e
}

// EpochSimulator is responsible for simulating actions in a given Epoch
type EpochSimulator struct {
	epoch                 *Epoch
	hasAntagonistApplied  bool
	hasProtagonistApplied bool
}

// Start begins the epoch simulation by allowing the competing individuals to do their thing
func (e *EpochSimulator) Start() (*EpochResult, error) {
	if e.epoch.program == nil {
		return nil, fmt.Errorf("epoch cannot have nil program")
	}
	if e.epoch.protagonist == nil {
		return nil, fmt.Errorf("epoch cannot have nil protagonist")
	}
	if e.epoch.antagonist == nil {
		return nil, fmt.Errorf("epoch cannot have nil antagonist")
	}

	return nil, nil
}

type EpochResult struct {
	engine *EvolutionEngine //Reference to underlying engine
}
