package main

func main() {

	var gen Generation
		gen.
			Cycle()

	var numberGenerator NumberGenerator
		numberGenerator.
			GenerateRandomInt()

	var mutateStrategy MutateStrategy
		mutateStrategy.
			TargetNonTerminals()
			TargetTerminals()

	//var addSubTreeStrategy AddSubTreeStrategy
	//	addSubTreeStrategy.

	var evolutionProcess *EvolutionProcess
	evolutionProcess = EvolutionEngine.
		Options(EvolutionParams{}).
		SetStartIndividual(Tree{}, Spec{}). // Todo Implement EvolutionProcess SetStartIndividual
		ZeroSumFitness(func() float32 { return 0}).

		FitnessEval(func() float32 { return 0} ). // Todo Implement EvolutionProcess FitnessEval
		ProgramEval( func() float32 { return 0} ). // Todo Implement EvolutionProcess ProgramEval
		Protagonist(100, func() float32 { return 0}, []Strategy{}). // Todo Implement EvolutionProcess Protagonist
		Antagonist(100, func() float32 { return 0}, []Strategy{}). // Todo Implement EvolutionProcess Antagonist
		AvailableStrategies([]Strategy{}).
		Generations(300). // Todo Implement EvolutionProcess Generations
		ParentSelection(EvolutionaryStrategy.Tournament).
		SurvivorSelection(EvolutionaryStrategy.Rank).
		OptimizationStrategy(EvolutionaryStrategy.Minimization).
		Parallelize(true). // Todo Implement EvolutionProcess Parallelize()
		GenerateStatistics("./stats.json"). // Todo Implement EvolutionProcess GenerateStatistics
		Start() // Todo Implement EvolutionProcess Start

	//evolutionResult.
	//	TopAntagonist()

	var epochSimulator *EpochSimulator
	epochSimulator.
		Antagonist().
		Protagonist()

	var epochResult *EpochResult
	epochResult = epochSimulator.Start()

	epochResult.
		Record()

	var prog Program
		prog.

}

func FitnessEval(i func() float32) {

}

type Spec []*EquationPairing

type EquationPairing struct {
	Dependent float32
	Independent float32
}

type MutateStrategy Strategy
type AddSubTreeStrategy Strategy


type Evolution struct {
	RootProgram *Program
}

func (e *Evolution) Evolve() {

}

type ProgramGenerator struct {

}


type Strategy interface { Apply(t *Tree) }
type Mutable interface { Mutate() *Tree }
type Evaluable interface { Eval() float32 }
type Fitnessable interface { Fitness() float32 }
type ApplyStrategeable interface { ApplyStrategy([]Strategy) }

type TreeGenerator interface { Generate() *Tree }
type ProgramGenerator interface { Generate() *Program }
type EvolutionGenerator interface {}
type Starter interface { Start()}