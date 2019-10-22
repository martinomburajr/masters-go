package evolution

var EvolutionEngineTestNil = EvolutionEngine{}

var EvolutionEngineTest0 = EvolutionEngine{
	StatisticsOutput: "stats.json",

	Parameters: EvolutionParams{
		FitnessStrategy:     FitnessProtagonistThresholdTally,
		EvaluationThreshold: 0.01,
	},
}

var EvolutionEngineTest1 = EvolutionEngine{
	StatisticsOutput: "stats.json",

	Parameters: EvolutionParams{
		FitnessStrategy:     FitnessProtagonistThresholdTally,
		EvaluationThreshold: 0.01,
	},
}
