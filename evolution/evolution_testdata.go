package evolution

var EvolutionEngineTestNil = EvolutionEngine{}

var EvolutionEngineTest0 = EvolutionEngine{
	EvaluationMinThreshold: 0.001,
	EvaluationThreshold:    0.01,
	Spec:                   SpecX,
	FitnessStrategy:        FitnessProtagonistThresholdTally,
	ParentSelection:        ParentSelectionElitism,
	ElitismPercentage:      1,
	AvailableStrategies:    []Strategy{StrategyAddSubTree, StrategyDeleteSubTree, StrategyMutateNode},
	GenerationCount:        10,
	MaxDepth:               10,
	StartIndividual:        Prog1,
	DepthPenalty:           5,
	StatisticsOutput:       "stats.json",
	SurvivorSelection:      SurvivorSelectionGenerational,
}

var EvolutionEngineTest1 = EvolutionEngine{
	EvaluationMinThreshold: 0.001,
	EvaluationThreshold:    0.01,
	Spec:                   SpecX,
	FitnessStrategy:        FitnessProtagonistThresholdTally,
	ParentSelection:        ParentSelectionElitism,
	ElitismPercentage:      1,
	AvailableStrategies:    []Strategy{StrategyAddSubTree, StrategyDeleteSubTree, StrategyMutateNode},
	GenerationCount:        100,
	MaxDepth:               10,
	StartIndividual:        Prog1,
	DepthPenalty:           5,
	StatisticsOutput:       "stats.json",
	SurvivorSelection:      SurvivorSelectionGenerational,
}
