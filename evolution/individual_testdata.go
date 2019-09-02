package evolution

// IndividualNilProgNil
var IndividualNilProgNil = Individual{
	Program:              &ProgNil,
	kind:                 0,
	id:                   "ANTAGONIST-",
	age:                  0,
	hasCalculatedFitness: false,
}

// IndividualProg0-Kind-0
var IndividualProg0Kind0 = Individual{
	Program:                  &Prog0,
	kind:                     0,
	id:                       "ANTAGONIST-",
	age:                      0,
	hasCalculatedFitness:     false,
	strategy:                 []Strategy{{Kind: StrategyMutateNode}, {Kind: StrategyDeleteSubTree}, {Kind: StrategyAddSubTree}},
	fitness:                  []int{},
	fitnessCalculationMethod: FitnessProtagonistThresholdTally,
	hasAppliedStrategy:       false,
}

// IndividualProg0-Kind-1
var IndividualProg0Kind1 = Individual{
	Program:                  &Prog0,
	kind:                     1,
	id:                       "PROTAGONIST-",
	age:                      0,
	hasCalculatedFitness:     false,
	strategy:                 []Strategy{{Kind: StrategyMutateNode}, {Kind: StrategyDeleteSubTree}, {Kind: StrategyAddSubTree}},
	fitness:                  []int{},
	fitnessCalculationMethod: FitnessProtagonistThresholdTally,
	hasAppliedStrategy:       false,
}

// IndividualProg1Kind1
var IndividualProg1Kind1 = Individual{
	Program:                  &Prog1,
	kind:                     1,
	id:                       "PROTAGONIST-",
	age:                      0,
	hasCalculatedFitness:     false,
	strategy:                 []Strategy{{Kind: StrategyMutateNode}, {Kind: StrategyDeleteSubTree}, {Kind: StrategyAddSubTree}},
	fitness:                  []int{},
	fitnessCalculationMethod: FitnessProtagonistThresholdTally,
	hasAppliedStrategy:       false,
}

// IndividualProg0-Kind-1
var IndividualProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0Kind0 = Individual{
	Program:                  &ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0,
	kind:                     0,
	id:                       "ANTAGONIST-",
	age:                      0,
	hasCalculatedFitness:     false,
	strategy:                 []Strategy{{Kind: StrategyMutateNode}, {Kind: StrategyDeleteSubTree}, {Kind: StrategyAddSubTree}},
	fitness:                  []int{},
	fitnessCalculationMethod: FitnessProtagonistThresholdTally,
	hasAppliedStrategy:       false,
}

// IndividualProg0-Kind-1
var IndividualProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0Kind1 = Individual{
	Program:                  &ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0,
	kind:                     1,
	id:                       "PROTAGONIST-",
	age:                      0,
	hasCalculatedFitness:     false,
	strategy:                 []Strategy{{Kind: StrategyMutateNode}, {Kind: StrategyDeleteSubTree}, {Kind: StrategyAddSubTree}},
	fitness:                  []int{},
	fitnessCalculationMethod: FitnessProtagonistThresholdTally,
	hasAppliedStrategy:       false,
}
