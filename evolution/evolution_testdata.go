package evolution

var EvolutionEngineTestNil = EvolutionEngine{}

var EvolutionEngineTest0 = EvolutionEngine{
	StatisticsOutput:  "stats.json",

	Parameters: EvolutionParams{
		MaxDepth:               10,
		FitnessStrategy:        FitnessProtagonistThresholdTally,
		EvaluationMinThreshold: 0.001,
		EvaluationThreshold:    0.01,
	},
}

var EvolutionEngineTest1 = EvolutionEngine{
	StatisticsOutput:  "stats.json",

	Parameters: EvolutionParams{
		MaxDepth:               10,
		FitnessStrategy:        FitnessProtagonistThresholdTally,
		EvaluationMinThreshold: 0.001,
		EvaluationThreshold:    0.01,
	},
}
