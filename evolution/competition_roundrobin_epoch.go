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
	return fmt.Sprintf("EPOCH-%d", count)
}

// TopologyRoundRobin creates the Epoch process. This process applies the antagonist Strategy first,
// and then the protagonist Strategy second.
// It then appends the Fitness values to each individual in the epoch.
func (e *Epoch) Start(perfectTreeMap map[string]PerfectTree, params EvolutionParams) error {
	if e.protagonist == nil {
		return fmt.Errorf("epoch cannot have nil protagonist")
	}
	if e.antagonist == nil {
		return fmt.Errorf("epoch cannot have nil antagonist")
	}
	if perfectTreeMap == nil {
		return fmt.Errorf("perfectTreeMap cannot be nil")
	}

	err := e.antagonist.ApplyAntagonistStrategy(params)
	if err != nil {
		return err
	}
	e.antagonist.HasAppliedStrategy = true

	err = e.protagonist.ApplyProtagonistStrategy(*e.antagonist.Program.T, params)
	if err != nil {
		return err
	}
	e.protagonist.HasAppliedStrategy = true

	if !e.hasProtagonistApplied && !e.hasAntagonistApplied {
		return fmt.Errorf("antagonist and protagonist havent applied Strategy to program")
	}

	antagonistFitness, protagonistFitness, antagonistFitnessDelta, protagonistFitnessDelta := 0.0, 0.0, 0.0, 0.0
	switch e.generation.engine.Parameters.FitnessStrategy.Type {
	case FitnessDualThresholdedRatio:
		antagonistFitness, protagonistFitness, antagonistFitnessDelta,
			protagonistFitnessDelta, err = ThresholdedRatioFitness(e.generation.engine.Parameters.Spec,
			e.antagonist.Program,
			e.protagonist.Program, e.generation.engine.Parameters.SpecParam.DivideByZeroStrategy)
		if err != nil {
			return err
		}
		break
	case FitnessMonoThresholdedRatio:
		antagonistFitness, protagonistFitness, antagonistFitnessDelta, protagonistFitnessDelta,
			err = ThresholdedRatioFitness(e.generation.
			engine.Parameters.Spec,
			e.protagonist.Program,
			e.protagonist.Program, e.generation.engine.Parameters.SpecParam.DivideByZeroStrategy)
		if err != nil {
			return err
		}
		break
	default:
		err = fmt.Errorf("unknown Fitness Strategy selected")
	}

	FitnessResolver(perfectTreeMap, e.antagonist, e.protagonist, antagonistFitness, antagonistFitnessDelta,
		protagonistFitness,
		protagonistFitnessDelta)
	return nil
}

func FitnessResolver(perfectTreeMap map[string]PerfectTree, antagonist, protagonist *Individual,
	antagonistFitness float64, antagonistFitnessDelta float64, protagonistFitness float64, protagonistFitnessDelta float64) {
	if perfectTreeMap[antagonist.Parent.Id].Program == nil {
		perfectTreeMap[antagonist.Parent.Id] = PerfectTree{BestFitnessValue: math.MinInt64}
	}
	perfectTreeAntagonist := perfectTreeMap[antagonist.Parent.Id]
	if perfectTreeAntagonist.BestFitnessValue < antagonistFitness {
		perfectTreeAntagonist.Program = antagonist.Program
		perfectTreeAntagonist.BestFitnessValue = antagonistFitness
		if antagonistFitnessDelta != math.Inf(1) {
			perfectTreeAntagonist.BestFitnessDelta = antagonistFitnessDelta
		}
		perfectTreeMap[antagonist.Parent.Id] = perfectTreeAntagonist
	}
	if perfectTreeMap[protagonist.Parent.Id].Program == nil {
		perfectTreeMap[protagonist.Parent.Id] = PerfectTree{BestFitnessValue: math.MinInt64}
	}
	perfectTreeProtagonist := perfectTreeMap[protagonist.Parent.Id]
	if perfectTreeProtagonist.BestFitnessValue < protagonistFitness {
		perfectTreeProtagonist.Program = protagonist.Program
		perfectTreeProtagonist.BestFitnessValue = protagonistFitness
		if protagonistFitnessDelta != math.Inf(1) {
			perfectTreeProtagonist.BestFitnessDelta = protagonistFitnessDelta
		} else {
			perfectTreeProtagonist.BestFitnessDelta = math.MaxInt16
		}
		perfectTreeMap[protagonist.Parent.Id] = perfectTreeProtagonist
	}
	if antagonistFitnessDelta != math.Inf(1) {
		antagonist.Parent.Deltas = append(antagonist.Parent.Deltas, antagonistFitnessDelta)
	} else {
		antagonist.Parent.Deltas = append(antagonist.Parent.Deltas, 0)
	}
	if protagonistFitnessDelta != math.Inf(1) {
		protagonist.Parent.Deltas = append(protagonist.Parent.Deltas, protagonistFitnessDelta)
	} else {
		protagonist.Parent.Deltas = append(protagonist.Parent.Deltas, math.MaxInt16)
	}
	antagonist.Parent.Fitness = append(antagonist.Parent.Fitness, antagonistFitness)
	protagonist.Parent.Fitness = append(protagonist.Parent.Fitness, protagonistFitness)
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
