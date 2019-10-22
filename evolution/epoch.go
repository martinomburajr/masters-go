package evolution

import (
	"fmt"
)

// Epoch is defined as a coevolutionary step where protagonist and antagonist compete.
// For example an epoch could represent a distinct interaction between two parties.
// For instance a bug mutated program (antagonist) can be challenged a variety of times (
// specified by {iterations}) by the tests (protagonist).
// The test will use up the strategies it contains and attempt to chew away at the antagonists Fitness,
// to maximize its own
type Epoch struct {
	id                               string
	protagonist                      Individual
	antagonist                       Individual
	generation                       *Generation
	program                          Program
	isComplete                       bool
	terminalSet                      []SymbolicExpression
	nonTerminalSet                   []SymbolicExpression
	hasAntagonistApplied             bool
	hasProtagonistApplied            bool
}

// CreateEpochID generates a given epoch Id with some useful information
func CreateEpochID(count int, generationId, antagonistId, protagonistId string) string {
	return fmt.Sprintf("EPOCH-%d-GEN-%s-ANTAGON-%s-PROTAGON-%s", count, generationId, antagonistId, protagonistId)
}

// Start creates the Epoch process. This process applies the antagonist Strategy first,
// and then the protagonist Strategy second.
// It then appends the Fitness values to each individual in the epoch.
func (e *Epoch) Start() error {
	if e.protagonist.Program == nil {
		return fmt.Errorf("epoch cannot have nil protagonist")
	}
	if e.antagonist.Program == nil {
		return fmt.Errorf("epoch cannot have nil antagonist")
	}

	err := e.applyAntagonistStrategy()
	if err != nil {
		return err
	}
	e.antagonist.HasAppliedStrategy = true

	err = e.applyProtagonistStrategy(*e.antagonist.Program.T)
	if err != nil {
		return err
	}
	e.protagonist.HasAppliedStrategy = true

	if !e.hasProtagonistApplied && !e.hasAntagonistApplied {
		return fmt.Errorf("antagonist and protagonist havent applied Strategy to program")
	}

	antagonistFitness, protagonistFitness := 0.0, 0.0
	switch e.generation.engine.Parameters.FitnessStrategy {
	case FitnessProtagonistThresholdTally:
		antagonistFitness, protagonistFitness, err = ProtagonistThresholdTally(e.generation.engine.Parameters.Spec,
			e.protagonist.Program, e.generation.engine.Parameters.EvaluationThreshold)
		if err != nil {
			return err
		}
		break
	case FitnessRatio:
		antagonistFitness, protagonistFitness, err = RatioFitness(e.generation.engine.Parameters.Spec, e.antagonist.Program,
			e.protagonist.Program)
		if err != nil {
			return err
		}
		break

	case FitnessRatioThresholder:
		antagonistFitness, protagonistFitness, err = RatioFitnessThresholded(e.generation.engine.Parameters.Spec,
			e.generation.engine.Parameters.ThresholdMultiplier, e.antagonist.Program,
			e.protagonist.Program)
		if err != nil {
			return err
		}
		break
	case FitnessDualThresholdedRatioFitness:
		antagonistFitness, protagonistFitness, err = ThresholdedRatioFitness(e.generation.engine.Parameters.Spec, e.antagonist.Program,
			e.protagonist.Program)
		if err != nil {
			return err
		}
		break
	case FitnessMonoThresholdedRatioFitness:
		antagonistFitness, protagonistFitness, err = ThresholdedRatioFitness(e.generation.engine.Parameters.Spec, e.antagonist.Program,
			e.protagonist.Program)
		if err != nil {
			return err
		}
		break
	default:
		err = fmt.Errorf("Unknown Fitness Strategy selected")
	}

	e.antagonist.Fitness = append(e.antagonist.Fitness, antagonistFitness)
	e.protagonist.Fitness = append(e.protagonist.Fitness, protagonistFitness)

	//antString := e.antagonist.ToString()
	//fmt.Println(antString.String())
	//proString := e.protagonist.ToString()
	//fmt.Println(proString.String())

	//program, err := e.generation.engine.Parameters.StartIndividual.Clone()
	//if err != nil {
	//	return err
	//}
	//e.antagonist.Program = &program
	//e.protagonist.Program = &Program{}
	return nil
}

// applyAntagonistStrategy applies the Antagonist strategies to program.
func (e *Epoch) applyAntagonistStrategy() error {
	program, err := e.generation.engine.Parameters.StartIndividual.Clone()
	if err != nil {
		return err
	}
	e.antagonist.Program = &program
	for _, strategy := range e.antagonist.Strategy {
		err := e.antagonist.Program.ApplyStrategy(strategy,
			e.terminalSet,
			e.nonTerminalSet,
			e.generation.engine.Parameters.DepthOfRandomNewTrees)
		if err != nil {
			return err
		}
	}
	e.hasAntagonistApplied = true
	return nil
}

// applyProtagonistStrategy Apply Protagonist strategies to program.
func (e *Epoch) applyProtagonistStrategy(antagonistTree DualTree) error {
	//if e.protagonist == nil {
	//	return fmt.Errorf("protagonist cannot be nil")
	//}
	if e.protagonist.Strategy == nil {
		return fmt.Errorf("protagonist stategy cannot be nil")
	}
	if len(e.protagonist.Strategy) < 1 {
		return fmt.Errorf("protagonist Strategy cannot be empty")
	}
	if antagonistTree.root == nil {
		return fmt.Errorf("applyProtagonistStrategy | antagonist supplied to protagonist has a nill root Tree")
	}
	tree, err := antagonistTree.Clone()
	if err != nil {
		return err
	}
	e.protagonist.Program.T = &tree

	for _, strategy := range e.protagonist.Strategy {
		err := e.protagonist.Program.ApplyStrategy(strategy,
			e.terminalSet,
			e.nonTerminalSet,
			e.generation.engine.Parameters.DepthOfRandomNewTrees)
		if err != nil {
			return err
		}
	}
	e.hasProtagonistApplied = true
	return nil
}
