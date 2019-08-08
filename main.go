package main

func main() {

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
		SetStartIndividual(Tree{}, Spec{}).
		FitnessEval(func() float32 { return 0} ).
		ProgramEval( func() float32 { return 0} ).
		ProtagonistCount(100).
		AntagonistCount(100).
		AvailableStrategies([]Strategdfy{}).
		Generations(300).
		ParentSelection(EvolutionaryStrategy.Tournament).
		SurvivorSelection(EvolutionaryStrategy.Rank).
		OptimizationStrategy(EvolutionaryStrategy.Minimization).
		Parallelize(true).
		GenerateStatistics("./stats.json").
		Start()

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