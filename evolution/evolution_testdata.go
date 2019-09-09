package evolution

var EvolutionEngineTestNil = EvolutionEngine{}

var EvolutionEngineTest0 = EvolutionEngine{
	EvaluationMinThreshold:           0.001,
	EvaluationThreshold:              0.01,
	Spec:                             SpecX,
	FitnessStrategy:                  FitnessProtagonistThresholdTally,
	ParentSelection:                  ParentSelectionElitism,
	ElitismPercentage:                1,
	AvailableStrategies:              []Strategy{{Kind: StrategyAddSubTree}, {Kind: StrategyDeleteSubTree}, {Kind: StrategyMutateNode}},
	GenerationCount:                  10,
	ProbabilityOfNonTerminalMutation: 0.05,
	ProbabilityOfMutation:            0.03,
	MaxDepth:                         10,
	StartIndividual:                  Prog1,
	DepthPenalty:                     5,
	EachPopulation:                   100,
	StatisticsOutput:                 "stats.json",
	SurvivorSelection:                SurvivorSelectionGenerational,
}

var EvolutionEngineTest1 = EvolutionEngine{
	EvaluationMinThreshold:           0.001,
	EvaluationThreshold:              0.01,
	Spec:                             SpecX,
	FitnessStrategy:                  FitnessProtagonistThresholdTally,
	ParentSelection:                  ParentSelectionElitism,
	ElitismPercentage:                1,
	AvailableStrategies:              []Strategy{{Kind: StrategyAddSubTree}, {Kind: StrategyDeleteSubTree}, {Kind: StrategyMutateNode}},
	GenerationCount:                  100,
	ProbabilityOfNonTerminalMutation: 0.05,
	ProbabilityOfMutation:            0.03,
	MaxDepth:                         10,
	StartIndividual:                  Prog1,
	DepthPenalty:                     5,
	EachPopulation:                   100,
	StatisticsOutput:                 "stats.json",
	SurvivorSelection:                SurvivorSelectionGenerational,
}
