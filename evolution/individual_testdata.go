package evolution

// IndividualNilProgNil
var IndividualNilProgNil = Individual{
	Program:              &ProgNil,
	Kind:                 0,
	Id:                   "ANTAGONIST-",
	Age:                  0,
	HasCalculatedFitness: false,
}

// IndividualProg0-Kind-0
var IndividualProg0Kind0 = Individual{
	Program:                  &ProgX,
	Kind:                     0,
	Id:                       "ANTAGONIST-",
	Age:                      0,
	HasCalculatedFitness:     false,
	Strategy:                 []Strategy{StrategyMutateTerminal, StrategyDeleteNonTerminal, StrategyAddSubTree},
	Fitness:                  []float64{},
	FitnessCalculationMethod: FitnessProtagonistThresholdTally,
	HasAppliedStrategy:       false,
}

// IndividualProg0-Kind-1
var IndividualProg0Kind1 = Individual{
	Program:                  &ProgX,
	Kind:                     1,
	Id:                       "PROTAGONIST-",
	Age:                      0,
	HasCalculatedFitness:     false,
	Strategy:                 []Strategy{StrategyMutateTerminal, StrategyDeleteNonTerminal, StrategyAddSubTree},
	Fitness:                  []float64{},
	FitnessCalculationMethod: FitnessProtagonistThresholdTally,
	HasAppliedStrategy:       false,
}

// IndividualProg1Kind1
var IndividualProg1Kind1 = Individual{
	Program:                  &Prog1,
	Kind:                     1,
	Id:                       "PROTAGONIST-",
	Age:                      0,
	HasCalculatedFitness:     false,
	Strategy:                 []Strategy{StrategyMutateTerminal, StrategyDeleteNonTerminal, StrategyAddSubTree},
	Fitness:                  []float64{},
	FitnessCalculationMethod: FitnessProtagonistThresholdTally,
	HasAppliedStrategy:       false,
}

// IndividualProgTreeT_NT_T_0
var IndividualProgTreeT_NT_T_0 = Individual{
	Program:                  &ProgTreeT_NT_T_0,
	Kind:                     1,
	Id:                       "PROTAGONIST-",
	Age:                      0,
	HasCalculatedFitness:     false,
	Strategy:                 []Strategy{StrategyMutateTerminal, StrategyDeleteNonTerminal, StrategyAddSubTree},
	Fitness:                  []float64{},
	FitnessCalculationMethod: FitnessProtagonistThresholdTally,
	HasAppliedStrategy:       false,
}

// IndividualProgTreeT_NT_T_4
var IndividualProgTreeT_NT_T_4 = Individual{
	Program:                  &ProgTreeT_NT_T_4,
	Kind:                     1,
	Id:                       "PROTAGONIST-",
	Age:                      0,
	HasCalculatedFitness:     false,
	Strategy:                 []Strategy{StrategyMutateTerminal, StrategyDeleteNonTerminal, StrategyAddSubTree},
	Fitness:                  []float64{},
	FitnessCalculationMethod: FitnessProtagonistThresholdTally,
	HasAppliedStrategy:       false,
}

// IndividualProgTreeT_NT_T_4
var IndividualProgTreeT_NT_T_1 = Individual{
	Program:                  &ProgTreeT_NT_T_1,
	Kind:                     1,
	Id:                       "PROTAGONIST-",
	Age:                      0,
	HasCalculatedFitness:     false,
	Strategy:                 []Strategy{StrategyMutateTerminal, StrategyDeleteNonTerminal, StrategyAddSubTree},
	Fitness:                  []float64{},
	FitnessCalculationMethod: FitnessProtagonistThresholdTally,
	HasAppliedStrategy:       false,
}

// IndividualProgTreeT_NT_T_NT_T_0
var IndividualProgTreeT_NT_T_NT_T_0 = Individual{
	Program:                  &ProgTreeT_NT_T_NT_T_0,
	Kind:                     1,
	Id:                       "PROTAGONIST-",
	Age:                      0,
	HasCalculatedFitness:     false,
	Strategy:                 []Strategy{StrategyMutateTerminal, StrategyDeleteNonTerminal, StrategyAddSubTree},
	Fitness:                  []float64{},
	FitnessCalculationMethod: FitnessProtagonistThresholdTally,
	HasAppliedStrategy:       false,
}

// IndividualProgTreeT_NT_T_NT_T_NT_T_0
var IndividualProgTreeT_NT_T_NT_T_NT_T_0 = Individual{
	Program:                  &ProgTreeT_NT_T_NT_T_NT_T_0,
	Kind:                     1,
	Id:                       "PROTAGONIST-",
	Age:                      0,
	HasCalculatedFitness:     false,
	Strategy:                 []Strategy{StrategyMutateTerminal, StrategyDeleteNonTerminal, StrategyAddSubTree},
	Fitness:                  []float64{},
	FitnessCalculationMethod: FitnessProtagonistThresholdTally,
	HasAppliedStrategy:       false,
}

// IndividualProgTreeT_NT_T_NT_T_NT_T_NT_T_1
var IndividualProgTreeT_NT_T_NT_T_NT_T_NT_T_1 = Individual{
	Program:                  &ProgTreeXXXX4,
	Kind:                     1,
	Id:                       "PROTAGONIST-",
	Age:                      0,
	HasCalculatedFitness:     false,
	Strategy:                 []Strategy{StrategyMutateTerminal, StrategyDeleteNonTerminal, StrategyAddSubTree},
	Fitness:                  []float64{},
	FitnessCalculationMethod: FitnessProtagonistThresholdTally,
	HasAppliedStrategy:       false,
}

// IndividualProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0Kind0-Kind-1
var IndividualProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0Kind0 = Individual{
	Program:                  &ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0,
	Kind:                     0,
	Id:                       "ANTAGONIST-",
	Age:                      0,
	HasCalculatedFitness:     false,
	Strategy:                 []Strategy{StrategyMutateTerminal, StrategyDeleteNonTerminal, StrategyAddSubTree},
	Fitness:                  []float64{},
	FitnessCalculationMethod: FitnessProtagonistThresholdTally,
	HasAppliedStrategy:       false,
}

// IndividualProg0-Kind-1
var IndividualProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0Kind1 = Individual{
	Program:                  &ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0,
	Kind:                     1,
	Id:                       "PROTAGONIST-",
	Age:                      0,
	HasCalculatedFitness:     false,
	Strategy:                 []Strategy{StrategyMutateTerminal, StrategyDeleteNonTerminal, StrategyAddSubTree},
	Fitness:                  []float64{},
	FitnessCalculationMethod: FitnessProtagonistThresholdTally,
	HasAppliedStrategy:       false,
}
