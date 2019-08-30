package evolution

var EvolutionNil = EvolutionEngine{}

var Evolution0 = EvolutionEngine{
	minThreshold: 0.001,
	threshold: 0.01,
	spec:SpecX,
	fitnessStrategy:FitnessProtagonistThresholdTally,
	parentSelection:ParentSelectionElitism,
	elitismPercentage: 1,
	availableStrategies: []Strategy{{Kind:StrategyAddSubTree}, {Kind:StrategyDeleteSubTree}, {Kind:StrategyMutateNode}},
	generations: 10,
	probabilityOfNonTerminalMutation: 0.05,
	probabilityOfMutation: 0.03,
	maxDepth: 10,
	startIndividual:Prog1,
	depthPenalty: 5,
	eachPopulation: 100,
	statisticsOutput: "stats.json",
	survivorSelection: SurvivorSelectionGenerational,
}

var Evolution1 = EvolutionEngine{
	minThreshold: 0.001,
	threshold: 0.01,
	spec:SpecX,
	fitnessStrategy:FitnessProtagonistThresholdTally,
	parentSelection:ParentSelectionElitism,
	elitismPercentage: 1,
	availableStrategies: []Strategy{{Kind:StrategyAddSubTree}, {Kind:StrategyDeleteSubTree}, {Kind:StrategyMutateNode}},
	generations: 100,
	probabilityOfNonTerminalMutation: 0.05,
	probabilityOfMutation: 0.03,
	maxDepth: 10,
	startIndividual:Prog1,
	depthPenalty: 5,
	eachPopulation: 100,
	statisticsOutput: "stats.json",
	survivorSelection: SurvivorSelectionGenerational,
}
