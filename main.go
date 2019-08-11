package main

import (
	"github.com/martinomburajr/masters-go/program"
)

func main() {

	//var gen Generation
	//	gen.
	//		Cycle()
	//
	//var numberGenerator NumberGenerator
	//	numberGenerator.
	//		GenerateRandomInt()
	//
	//var mutateStrategy MutateStrategy
	//	mutateStrategy.
	//		TargetNonTerminals()
	//		TargetTerminals()

	//var addSubTreeStrateg  y AddSubTreeStrategy
	//	addSubTreeStrategy.

	var initialProgram program.InitialProgram
		initialProgram.Spec(Spec{})

	var evolutionEngine = EvolutionEngine{} //Create the Evolution Engine

	var evolutionProcess *EvolutionProcess
	evolutionProcess = evolutionEngine.
		SetStartIndividual(initialProgram). // Todo Implement EvolutionProcess SetStartIndividual
		ZeroSumFitness(func() float32 { return 0}).

		FitnessEval(func() float32 { return 0} ). // Todo Implement EvolutionProcess FitnessEval
		ProgramEval( func() float32 { return 0} ). // Todo Implement EvolutionProcess ProgramEval
		Protagonist(100, func() float32 { return 0}, []Strategable{}). // Todo Implement EvolutionProcess Protagonist
		Antagonist(100, func() float32 { return 0}, []Strategable{}). // Todo Implement EvolutionProcess Antagonist
		AvailableStrategies([]Strategable{}).
		Generations(300). // Todo Implement EvolutionProcess Generations
		ParentSelection(EvolutionaryStrategy.Tournament).
		SurvivorSelection(EvolutionaryStrategy.Rank).
		OptimizationStrategy(EvolutionaryStrategy.Minimization).
		Parallelize(true). // Todo Implement EvolutionProcess Parallelize()
		GenerateStatistics("./stats.json"). // Todo Implement EvolutionProcess GenerateStatistics
		Options(EvolutionParams{}).
		Start() // Todo Implement EvolutionProcess Start



	//evolutionResult.
	//	TopAntagonist()

	//var epochSimulator *EpochSimulator
	//epochSimulator.
	//	Antagonist().
	//	Protagonist()
	//
	//var epochResult *EpochResult
	//epochResult = epochSimulator.Start()
	//
	//epochResult.
	//	Record()
	//
	//var prog program.Program
	//	prog.

}

func FitnessEval(i func() float32) {

}





type MutateStrategy Strategable
type AddSubTreeStrategy Strategable


type Evolution struct {
	RootProgram *program.Program
}

func (e *Evolution) Evolve() {

}

// TODO Create Mutate strategy
// TODO Create AddSubTree strategy
// TODO Create DeleteSubTree strategy
// TODO Create SoftDeleteSubTree strategy
// TODO Create SwapSubTree strategy

type Mutable interface { Mutate() *program.DualTree }
type Evaluable interface { Eval() float32 }
type Fitnessable interface { Fitness() float32 }
type ApplyStrategeable interface { ApplyStrategy([]Strategable) }

type TreeGenerator interface { Generate() *program.DualTree }
type ProgramGenerator interface { Generate() *program.Program }
type EvolutionGenerator interface {}
type Starter interface { Start()}

// TODO Create Error Function
