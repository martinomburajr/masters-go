package main

import (
	"fmt"
	"io"
)

type EvolutionParams struct {
	Generations       int
	EnableParallelism bool
}

type EvolutionEngine struct {
	startIndividual     *Program
	spec                Spec
	generations         int
	parallelize         bool
	availableStrategies []*Strategy
	Fitnessable
	programEval      func() float32
	statisticsOutput string
}

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
	if e.Fitnessable
	return nil
}

func (e *EvolutionEngine) Options(params EvolutionParams) *EvolutionEngine {
	return nil
}

func (e *EvolutionEngine) SetStartIndividual(individual Tree, spec Spec) *EvolutionEngine {
	e.spec = &spec
	return nil
}

func (e *EvolutionEngine) FitnessEval(fitnessFunc func() float32) *EvolutionEngine {
	return nil
}

func (e *EvolutionEngine) ProgramEval(programFunc func() float32) *EvolutionEngine {
	return nil
}

func (e *EvolutionEngine) ProtagonistCount(count int) *EvolutionEngine {
	return nil
}

func (e *EvolutionEngine) AntagonistCount(count int) *EvolutionEngine {
	return nil
}

func (e *EvolutionEngine) AvailableStrategies(strategies []Strategy) *EvolutionEngine {
	return nil
}

func (e *EvolutionEngine) Generations(i int) *EvolutionEngine {
	return nil
}

func (e *EvolutionEngine) ParentSelection(b bool) *EvolutionEngine {
	return nil
}

func (e *EvolutionEngine) SurvivorSelection(b bool) *EvolutionEngine {
	return nil
}

func (e *EvolutionEngine) OptimizationStrategy(b bool) *EvolutionEngine {
	return nil
}

func (e *EvolutionEngine) Parallelize(b bool) *EvolutionEngine {
	return nil
}

func (e *EvolutionEngine) GenerateStatistics(s string) *EvolutionEngine {
	return nil
}

func (e *EvolutionEngine) Start() *EvolutionProcess{
	e.validate()
	return nil
}

/**
 EvolutionProcess represents the state of an evolutionary process once the evolution engine starts
 */
type EvolutionProcess struct {
	currentGeneration int
	protagonists []*Program
	antagonists []*Program
	engine *EvolutionEngine
	strategiesApplied []Strategy
}

type IEvolutionProcess interface {

}


type Epoch struct {
	protagonist *Program
	antagonist *Program
	engine *EvolutionEngine
	isComplete bool
	generation int
}


type EpochSimulator struct {

}

func (e *EpochSimulator) Antagonist() *EpochSimulator{
	return nil
}

func (e *EpochSimulator) Protagonist() *EpochSimulator{
	return nil
}

func (e *EpochSimulator) Start() *EpochResult{
	return nil
}

type EpochResult struct {
	engine *EvolutionEngine //Reference to underlying engine
}