package evolution

var GenerationNil = Generation{}

var GenerationTest0 = Generation{
	engine:                       &EvolutionEngineTest0,
	hasParentSelectionHappened:   false,
	GenerationID:                 "gen0",
	isComplete:                   false,
	Protagonists:                 []Individual{IndividualProg0Kind0, IndividualProg0Kind1},
	Antagonists:                  []Individual{IndividualProg1Kind1, IndividualNilProgNil},
	hasSurvivorSelectionHappened: false,
}

// GenerationTest1 1 less protagonist
var GenerationTest1 = Generation{
	engine:                       &EvolutionEngineTest0,
	hasParentSelectionHappened:   false,
	GenerationID:                 "gen0",
	isComplete:                   false,
	Protagonists:                 []Individual{IndividualProg0Kind0},
	Antagonists:                  []Individual{IndividualProg1Kind1, IndividualNilProgNil},
	hasSurvivorSelectionHappened: false,
}
