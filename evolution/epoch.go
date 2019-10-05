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
	generation                       *Generation
	program                          Program
	protagonistBegins                bool
	isComplete                       bool
	probabilityOfMutation            float32
	probabilityOfNonTerminalMutation float32
	terminalSet                      []SymbolicExpression
	nonTerminalSet                   []SymbolicExpression
	hasAntagonistApplied             bool
	hasProtagonistApplied            bool
}

// CreateEpochID generates a given epoch id with some useful information
func CreateEpochID(count int, generationId, antagonistId, protagonistId string) string {
	return fmt.Sprintf("EPOCH-%d-GEN-%s-ANTAGON-%s-PROTAGON-%s", count, generationId, antagonistId, protagonistId)
}

// Start creates the Epoch process. This process applies the antagonist strategy first,
// and then the protagonist strategy second.
// It then appends the fitness values to each individual in the epoch.
// Protagonists will by default have nil Trees as their trees depend on those of the Antagonists
func (e *Epoch) Start() error {
	if e.antagonist.Program.T.root == nil {
		return fmt.Errorf("epoch cannot have nil antagonist tree root")
	}

	err := e.applyAntagonistStrategy()
	if err != nil {
		return err
	}
	e.antagonist.hasAppliedStrategy = true

	err = e.applyProtagonistStrategy(e.antagonist.Program.T)
	if err != nil {
		return err
	}
	e.protagonist.hasAppliedStrategy = true

	if !e.hasProtagonistApplied && !e.hasAntagonistApplied {
		return fmt.Errorf("antagonist and protagonist havent applied strategy to program")
	}

	antagonistFitness, protagonistFitness := 0, 0
	switch e.generation.engine.Parameters.FitnessStrategy {
	case FitnessProtagonistThresholdTally:
		antagonistFitness, protagonistFitness, err = ProtagonistThresholdTally(e.generation.engine.Parameters.Spec,
			&e.program, e.generation.engine.Parameters.EvaluationThreshold,
			e.generation.engine.Parameters.EvaluationMinThreshold)
		if err != nil {
			return err
		}
	}
	e.antagonist.age++
	e.protagonist.age++

	e.antagonist.fitness = append(e.antagonist.fitness, antagonistFitness)
	e.protagonist.fitness = append(e.protagonist.fitness, protagonistFitness)

	return nil
}

// applyAntagonistStrategy applies the Antagonist strategies to program.
func (e *Epoch) applyAntagonistStrategy() error {
	for _, strategy := range e.antagonist.strategy {
		err := e.antagonist.Program.ApplyStrategy(strategy,
			e.terminalSet,
			e.nonTerminalSet,
			e.probabilityOfMutation,
			e.probabilityOfNonTerminalMutation,
			e.generation.engine.Parameters.DepthOfRandomNewTrees,
			e.generation.engine.Parameters.DeletionType)
		if err != nil {
			return err
		}
	}
	e.hasAntagonistApplied = true
	return nil
}

// applyProtagonistStrategy Apply Protagonist strategies to program.
func (e *Epoch) applyProtagonistStrategy(antagonistTree DualTree) error {
	if e.protagonist.strategy == nil {
		return fmt.Errorf("protagonist stategy cannot be nil")
	}
	if len(e.protagonist.strategy) < 1 {
		return fmt.Errorf("protagonist strategy cannot be empty")
	}
	if antagonistTree.root == nil {
		return fmt.Errorf("applyProtagonistStrategy | antagonist supplied to protagonist has a nill root tree")
	}
	e.protagonist.Program.T = antagonistTree

	for _, strategy := range e.protagonist.strategy {
		err := e.protagonist.Program.ApplyStrategy(strategy,
			e.terminalSet,
			e.nonTerminalSet,
			e.probabilityOfMutation,
			e.probabilityOfNonTerminalMutation,
			e.generation.engine.Parameters.DepthOfRandomNewTrees,
			e.generation.engine.Parameters.DeletionType)
		if err != nil {
			return err
		}
	}
	e.hasProtagonistApplied = true
	return nil
}
