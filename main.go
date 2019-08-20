package main

import (
	"github.com/martinomburajr/masters-go/evolution"
	"log"
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

	var initialProgram evolution.InitialProgram
	initialProgram.Spec(nil)

	var evolutionEngine = evolution.EvolutionEngine{} //Create the Evolution Engine

	var evolutionProcess *evolution.EvolutionProcess
	evolutionProcess = evolutionEngine.
		SetStartIndividual(initialProgram). // Todo Implement EvolutionProcess SetStartIndividual
		ZeroSumFitness(func() float32 { return 0 }).
		FitnessEval(func() float32 { return 0 }).                                 // Todo Implement EvolutionProcess FitnessEval
		ProgramEval(func() float32 { return 0 }).                                 // Todo Implement EvolutionProcess ProgramEval
		Protagonist(100, func() float32 { return 0 }, []evolution.Strategable{}). // Todo Implement EvolutionProcess Protagonist
		Antagonist(100, func() float32 { return 0 }, []evolution.Strategable{}).  // Todo Implement EvolutionProcess antagonist
		AvailableStrategies([]evolution.Strategable{}).
		Generations(300). // Todo Implement EvolutionProcess Generations
		ParentSelection(evolution.EvolutionaryStrategy.Tournament).
		SurvivorSelection(evolution.EvolutionaryStrategy.Rank).
		OptimizationStrategy(evolution.EvolutionaryStrategy.Minimization).
		Parallelize(true).                  // Todo Implement EvolutionProcess Parallelize()
		GenerateStatistics("./stats.json"). // Todo Implement EvolutionProcess GenerateStatistics
		Options(evolution.EvolutionParams{}).
		Start() // Todo Implement EvolutionProcess Start

	log.Print(evolutionProcess)

	//evolutionResult.
	//	TopAntagonist()

	//var epochSimulator *EpochSimulator
	//epochSimulator.
	//	antagonist().
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

type MutateStrategy evolution.Strategable
type AddSubTreeStrategy evolution.Strategable

type Evolution struct {
	RootProgram *evolution.Program
}

func (e *Evolution) Evolve() {

}

// TODO Create MutateTerminal strategy
// TODO Create AddSubTree strategy
// TODO Create DeleteSubTree strategy
// TODO Create SoftDeleteSubTree strategy
// TODO Create SwapSubTree strategy

//type Mutable interface { MutateTerminal() *program.DualTree }
//type Evaluable interface { Eval() float32 }
//type Fitnessable interface { Fitness() float32 }
//type ApplyStrategeable interface { ApplyStrategy([]program.Strategable) }
//
//type TreeGenerator interface { Generate() *program.DualTree }
//type ProgramGenerator interface { Generate() *program.Program }
//type EvolutionGenerator interface {}
//type Starter interface { Start()}

// TODO Create Error Function
