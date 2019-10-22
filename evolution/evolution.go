package evolution

import (
	"fmt"
)

type EvolutionParams struct {
	GenerationsCount                 int `json:"generationCount"`
	EachPopulationSize               int `json:"eachPopulationSize"`
	EnableParallelism                bool `json:"enableParallelism"`
	ElitismPercentage                float64

	FitnessStrategy                  int
	TournamentSize                   int
	// EachPopulationSize represents the size of each protagonist or antagonist population.
	// This value must be even otherwise pairwise operations such as crossover will fail

	//ProbabilityOfRecombination       float64
	ProbabilityOfMutation            float64
	//ProbabilityOfNonTerminalMutation float64

	// ThresholdMultiplier is used when the FitnessRatioThreshold option is selected.
	// It creates a threshold value based on the cumulative value of the dependent variable in the spec.
	// It cannot be less than 1 as a value less than 1 would mean that the approximater functions would need to be
	// better than the spec and that is not possible
	ThresholdMultiplier float64

	AntagonistMaxStrategies int
	//AntagonistStrategyLength         int
	ProtagonistMaxStrategies int
	//ProtagonistStrategyLength        int


	// Strategies is a list of available strategies for each individual.
	// These can be randomly allocated to individuals and duplicates are expected.
	//Strategies            []Strategy
	DepthOfRandomNewTrees int
	// StrategyLengthPenalty is the penalty given to strategies that exceed a given length
	StrategyLengthPenalty float64
	// StrategyLengthLimit is the maximum length a Strategy can reach before being penalized
	StrategyLengthLimit int

	// DeletionType pertains to the different kinds of deletion operations possible for a given Tree.
	//DeletionType int

	// EnforceIndependentVariable ensures that during individual Generation at the start of the evolution,
	// independent variables are injected to the program meaning every program will at least have one independent
	// variable e.g. X
	EnforceIndependentVariable bool

	// CrossoverPercentage pertains to the amount of genetic material crossed-over.
	// This is a percentage represented as a float64. A value of 1 means all material is swapped.
	// A value of 0 means no material is swapped (which in effect are the same thing).
	// Avoid 0 or 1 use values in between
	CrossoverPercentage float64

	//MaintainGeneTransferEquality pertains to a scenario where in the event of differing Strategy lengths betweween
	// two individuals, if N objects are to be swapped during crossover,
	// then N objects will be swapped from both individuals.
	// This prevents one individual receiving more genetic material from another individual and vice-versa
	MaintainCrossoverGeneTransferEquality bool

	TerminalSet    SymbolicExpressionSet
	NonTerminalSet SymbolicExpressionSet

	EvaluationThreshold float64
	ParentSelection     int
	StartIndividual     Program

	SurvivorSelection   int
	SurvivorPercentage float64

	ProtagonistAvailableStrategies []Strategy
	AntagonistAvailableStrategies  []Strategy
	SetEqualStrategyLength         bool
	EqualStrategiesLength          int

	// VariableTerminals represent all the potential variables that may appear.
	// An effort is made to differentiate them from constants so that constants do not get overwritten as variables
	// would when calculating the spec.
	VariableTerminals []SymbolicExpression
	// AntagonistThresholdMultiplier is the multiplier applied to the antagonist delta when calculating fitness.
	// A large value means that antagonists have to attain a greater delta from the spec in order to gain adequate
	// fitness, conversely a smaller value gives the antagonists more slack to not manipulate the program excessively.
	// For good results set it to a value greater than that of the protagonist delta.
	// This value is only used when using DualThresholdedRatioFitness.
	AntagonistThresholdMultiplier float64

	// ProtagonistThresholdMultiplier is the multiplier applied to the protagonist delta when calculating fitness.
	// A large value means that protagonist can be less precise and gain adequate fitness,
	// conversely a smaller value gives the protagonist little room for mistake between its delta and that of the spec.
	// this value is used in both DualThresholdedRatioFitness and ThresholdedRatioFitness as a fitness value for
	// both antagonist and protagonists thresholds.
	ProtagonistThresholdMultiplier float64

	//AvailableOperators           AvailableVariablesAndOperators
	// ShouldRunInteractiveTerminal ensures the interactive terminal is run at the end of the evolution that allows
	// users to query all individuals in all generations.
	ShouldRunInteractiveTerminal bool

	SpecParam SpecParam `json:"spec"`

	Selection Selection `json:"selection"`
	//FitnessStrategy FitnessStrategy `json:"fitnessStrategy"`
	Spec                SpecMulti
}


type AvailableVariablesAndOperators struct {
	Constants []string
	Variables []string
	Operators []string
}

type Strategies struct {
	AvailableStrategies []string
	AntagonistAvailableStrategies []string
	ProtagonistAvailableStrategies []string
	StrategySize int
}

type FitnessStrategy struct {
	Type int
	AntagonistThresholdMultiplier float64
	ProtagonistThresholdMultiplier float64
}

type SpecParam struct {
	// SpecRange defines a range of variables on either side of the X axis. A range of 4 will include -2, -1,
	// 0 and 1.
	Range int
	//Expression is the actual expression being tested.
	// It is the initial function that is converted to the startIndividual
	Expression string
	Seed int
	AvailableVariablesAndOperators AvailableVariablesAndOperators
}

type Reproduction struct {
	// CrossoverPercentage pertains to the amount of genetic material crossed-over.
	// This is a percentage represented as a float64. A value of 1 means all material is swapped.
	// A value of 0 means no material is swapped (which in effect are the same thing).
	// Avoid 0 or 1 use values in between
	CrossoverPercentage float64
	ProbabilityOfMutation            float64
}
type Selection struct {
	Parent ParentSelection
	Survivor SurvivorSelection
}

type ParentSelection struct {
	Type int
	TournamentSize int
}

type SurvivorSelection struct {
	Type int
	// SurvivorPercentage represents how many individulas in the parent vs child population should continue.
	// 1 means all parents move on. 0 means only children move on. Any number in betwee is a percentage value.
	// It cannot be greater than 1 or less than 0.
	SurvivorPercentage float64
}


type EvolutionEngine struct {
	Parallelize         bool `json:"parallelize"`
	Generations         []*Generation
	StatisticsOutput    string          `json:"statisticsOutput"`
	Parameters          EvolutionParams `json:"parameters"`
	IsMoreFitnessBetter bool            `json:"isMoreFitnessBetter"`
}

func (e *EvolutionEngine) Start() (EvolutionResult, error) {
	err := e.validate()
	if err != nil {
		return EvolutionResult{}, err
	}

	// Set First Generation - TODO Parallelize Individual Creation
	genID := GenerateGenerationID(0)
	gen0 := Generation{
		count:        0,
		GenerationID: genID,
		Protagonists: nil,
		Antagonists:  nil,
		engine:       e,
	}
	e.Generations[0] = &gen0

	antagonists, err := e.Generations[0].GenerateRandomIndividual(IndividualAntagonist,
		e.Parameters.StartIndividual)
	if err != nil {
		return EvolutionResult{}, err
	}

	protagonists, err := e.Generations[0].GenerateRandomIndividual(IndividualProtagonist,
		Program{})
	if err != nil {
		return EvolutionResult{}, err
	}

	gen0.Protagonists = protagonists
	gen0.Antagonists = antagonists

	// cycle through generationCount
	e.Generations[0] = &gen0
	for i := 0; i < e.Parameters.GenerationsCount-1; i++ {
		protagonistsCleanse, err := CleansePopulation(e.Generations[i].Protagonists, *e.Parameters.StartIndividual.T)
		if err != nil {
			return EvolutionResult{}, err
		}
		antagonistsCleanse, err := CleansePopulation(e.Generations[i].Antagonists, *e.Parameters.StartIndividual.T)
		if err != nil {
			return EvolutionResult{}, err
		}

		e.Generations[i].Protagonists = protagonistsCleanse
		e.Generations[i].Antagonists = antagonistsCleanse
		nextGeneration, err := e.Generations[i].Start(i)
		if err != nil {
			return EvolutionResult{}, err
		}
		e.Generations[i+1] = nextGeneration
	}

	evolutionResult := EvolutionResult{}
	err = evolutionResult.Analyze(e.Generations, e.IsMoreFitnessBetter)
	if err != nil {
		return EvolutionResult{}, err
	}

	return evolutionResult, nil
}

// Todo Implement EvolutionProcess validate
func (e *EvolutionEngine) validate() error {
	if e.Parameters.GenerationsCount < 1 {
		return fmt.Errorf("set number of generationCount by calling e.GenerationsCount(x)")
	}
	if e.Parameters.EachPopulationSize%2 != 0 {
		return fmt.Errorf("set number of EachPopulationSize to an Even number")
	}
	if e.Parameters.SetEqualStrategyLength == true && e.Parameters.EqualStrategiesLength < 1 {
		return fmt.Errorf("cannot SetEqualStrategyLength to true and EqualStrategiesLength less than 1")
	}
	if e.Parameters.StartIndividual.T == nil {
		return fmt.Errorf("start individual cannot have a nil Tree")
	}
	if e.Parameters.Spec == nil {
		return fmt.Errorf("spec cannot be nil")
	}
	if len(e.Parameters.Spec) < 1 {
		return fmt.Errorf("spec cannot be empty")
	}
	if e.Parameters.FitnessStrategy == FitnessRatioThresholder && e.Parameters.ThresholdMultiplier < 1 {
		return fmt.Errorf("ThresholdMultiplier cannot be less than 1")
	}
	if e.Parameters.SurvivorPercentage > 1 || e.Parameters.SurvivorPercentage < 0 {
		return fmt.Errorf("SurvivorPercentage cannot be less than 0 or greater than 1. It is a percent value")
	}
	//err := e.StartIndividual.Validate()
	//if err != nil {
	//	return err
	//}

	if len(e.Parameters.Spec) < 3 {
		return fmt.Errorf("a small spec will hamper evolutionary accuracy")
	}
	return nil
}
