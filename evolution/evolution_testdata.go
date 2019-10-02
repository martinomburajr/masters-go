package evolution

var EvolutionEngineTestNil = EvolutionEngine{}

var EvolutionEngineTest0 = EvolutionEngine{
	Spec:              SpecX,
	ParentSelection:   ParentSelectionElitism,
	ElitismPercentage: 1,
	StartIndividual:   Prog1,

	StatisticsOutput:  "stats.json",
	SurvivorSelection: SurvivorSelectionGenerational,
	Parameters: EvolutionParams{
		MaxDepth:               10,
		FitnessStrategy:        FitnessProtagonistThresholdTally,
		EvaluationMinThreshold: 0.001,
		EvaluationThreshold:    0.01,
	},
}

var EvolutionEngineTest1 = EvolutionEngine{

	Spec: SpecX,

	ParentSelection:   ParentSelectionElitism,
	ElitismPercentage: 1,

	StartIndividual: Prog1,

	StatisticsOutput:  "stats.json",
	SurvivorSelection: SurvivorSelectionGenerational,
	Parameters: EvolutionParams{
		MaxDepth:               10,
		FitnessStrategy:        FitnessProtagonistThresholdTally,
		EvaluationMinThreshold: 0.001,
		EvaluationThreshold:    0.01,
	},
}
