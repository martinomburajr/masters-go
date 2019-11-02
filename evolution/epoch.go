package evolution

import (
	"fmt"
	"math"
)

// Epoch is defined as a coevolutionary step where protagonist and antagonist compete.
// For example an epoch could represent a distinct interaction between two parties.
// For instance a bug mutated program (antagonist) can be challenged a variety of times (
// specified by {iterations}) by the tests (protagonist).
// The test will use up the strategies it contains and attempt to chew away at the antagonists Fitness,
// to maximize its own
type Epoch struct {
	id                    string
	protagonist           *Individual
	antagonist            *Individual
	generation            *Generation
	program               Program
	isComplete            bool
	terminalSet           []SymbolicExpression
	nonTerminalSet        []SymbolicExpression
	hasAntagonistApplied  bool
	hasProtagonistApplied bool
}

// CreateEpochID generates a given epoch Id with some useful information
func CreateEpochID(count int, generationId, antagonistId, protagonistId string) string {
	return fmt.Sprintf("EPOCH-%d-GEN-%s-ANTAGON-%s-PROTAGON-%s", count, generationId, antagonistId, protagonistId)
}

// Start creates the Epoch process. This process applies the antagonist Strategy first,
// and then the protagonist Strategy second.
// It then appends the Fitness values to each individual in the epoch.
func (e *Epoch) Start(perfectTreeMap map[string]PerfectTree) error {
	if e.protagonist == nil {
		return fmt.Errorf("epoch cannot have nil protagonist")
	}
	if e.antagonist == nil {
		return fmt.Errorf("epoch cannot have nil antagonist")
	}
	if perfectTreeMap == nil {
		return fmt.Errorf("perfectTreeMap cannot be nil")
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
	switch e.generation.engine.Parameters.FitnessStrategy.Type {
	case FitnessProtagonistThresholdTally:
		antagonistFitness, protagonistFitness, err = ProtagonistThresholdTally(e.generation.engine.Parameters.Spec,
			e.protagonist.Program, e.generation.engine.Parameters.FitnessStrategy.AntagonistThresholdMultiplier)
		if err != nil {
			return err
		}
		break
	case FitnessThresholdedAntagonistRatio:
		antagonistFitness, protagonistFitness, err = ThresholdedAntagonistRatioFitness(e.generation.engine.Parameters.Spec, e.antagonist.Program,
			e.protagonist.Program, e.generation.engine.Parameters.FitnessCalculatorType)
		if err != nil {
			return err
		}
		break
	case FitnessRatio:
		antagonistFitness, protagonistFitness, err = RatioFitness(e.generation.engine.Parameters.Spec, e.antagonist.Program,
			e.protagonist.Program, e.generation.engine.Parameters.FitnessCalculatorType)
		if err != nil {
			return err
		}
		break

	case FitnessDualThresholdedRatio:
		antagonistFitness, protagonistFitness, err = ThresholdedRatioFitness(e.generation.engine.Parameters.Spec, e.antagonist.Program,
			e.protagonist.Program, e.generation.engine.Parameters.FitnessCalculatorType)
		if err != nil {
			return err
		}
		break
	case FitnessMonoThresholdedRatio:
		antagonistFitness, protagonistFitness, err = ThresholdedRatioFitness(e.generation.engine.Parameters.Spec,
			e.protagonist.Program,
			e.protagonist.Program, e.generation.engine.Parameters.FitnessCalculatorType)
		if err != nil {
			return err
		}
		break
	default:
		err = fmt.Errorf("unknown Fitness Strategy selected")
	}

	if perfectTreeMap[e.antagonist.Id].Program == nil {
		perfectTreeMap[e.antagonist.Id] = PerfectTree{FitnessValue: math.MinInt64}
	}
	perfectTreeAntagonist := perfectTreeMap[e.antagonist.Id]
	if perfectTreeAntagonist.FitnessValue < antagonistFitness {
		perfectTreeAntagonist.Program = e.antagonist.Program
		perfectTreeAntagonist.FitnessValue = antagonistFitness
		perfectTreeAntagonist.FitnessDetla = antagonistFitness
		perfectTreeMap[e.antagonist.Id] = perfectTreeAntagonist
	}

	if perfectTreeMap[e.protagonist.Id].Program == nil {
		perfectTreeMap[e.protagonist.Id] = PerfectTree{FitnessValue: math.MinInt64}
	}
	perfectTreeProtagonist := perfectTreeMap[e.protagonist.Id]
	if perfectTreeProtagonist.FitnessValue < protagonistFitness {
		perfectTreeProtagonist.Program = e.protagonist.Program
		perfectTreeProtagonist.FitnessValue = protagonistFitness
		perfectTreeProtagonist.FitnessDetla = protagonistFitness
		perfectTreeMap[e.protagonist.Id] = perfectTreeProtagonist
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

// AggregateFitness simply adds all the Fitness values of a given individual to come up with a total number.
// If the Fitness array is nil or empty return MaxInt8 as values such as -1 or 0 have a differnt meaning
func AggregateFitness(individual Individual) (float64, error) {
	if individual.Fitness == nil {
		return math.MaxInt8, fmt.Errorf("individuals Fitness arr cannot be nil")
	}
	if len(individual.Fitness) == 0 {
		return math.MaxInt8, fmt.Errorf("individuals Fitness arr cannot be empty")
	}

	sum := 0.0
	for i := 0; i < len(individual.Fitness); i++ {
		sum += individual.Fitness[i]
	}
	return sum, nil
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
			e.generation.engine.Parameters.Strategies.DepthOfRandomNewTrees)
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
			e.generation.engine.Parameters.Strategies.DepthOfRandomNewTrees)
		if err != nil {
			return err
		}
	}
	e.hasProtagonistApplied = true
	return nil
}
