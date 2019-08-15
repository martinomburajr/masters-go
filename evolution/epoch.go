package evolution

import "github.com/martinomburajr/masters-go/program"

// Epoch is defined as a coevolutionary step where protagonist and antagonist compete.
// For example an epoch could represent a distinct interaction between two parties.
// For instance a bug mutated program (antagonist) can be challenged a variety of times (
// specified by {iterations}) by the tests (protagonist).
// The test will use up the strategies it contains and attempt to chew away at the antagonists fitness,
// to maximize its own
type Epoch struct {
	protagonist *program.Program
	antagonist  *program.Program
	protagonistBegins bool
	iterations  int
	isComplete  bool
	generation  *Generation
}

func (e *Epoch) GetProtagonistBegins() bool {
	return e.protagonistBegins
}

// ProtagonistBegins states whether the protagonist should start the epoch
func (e *Epoch) ProtagonistBegins(protagonistBegins bool) *Epoch {
	e.protagonistBegins = protagonistBegins
	return e
}

// Protagonist sets the protagonist for the epoch
func (e *Epoch) Protagonist(protagonist *program.Program) *Epoch {
	e.protagonist = protagonist
	return e
}

// Antagonist sets the antagonist for the epoch
func (e *Epoch) Antagonist(antagonist *program.Program) *Epoch {
	e.antagonist = antagonist
	return e
}

func (e *Epoch) Iterations(iterations int) *Epoch {
	e.iterations = iterations
	return e
}

// EpochSimulator is responsible for simulating actions in a given Epoch
type EpochSimulator struct {
	epoch *Epoch
	hasAntagonistApplied bool
	hasProtagonistApplied bool
}

// Start begins the epoch simulation by allowing the competing individuals to do their thing
func (e *Epoch) Start() *EpochResult {
	return nil
}

type EpochResult struct {
	engine *EvolutionEngine //Reference to underlying engine
}

