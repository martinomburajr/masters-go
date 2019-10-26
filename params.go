package main

import "github.com/martinomburajr/masters-go/evolution"

var AllExpressions = []string{
	"1", "2", "8",
	"x+1", "x+20", "x-20", "",
	"x", "x*x", "(0-1*x*x)", "0-1*x*x*x", "x*x*x*x",
	"x*x+2*x+1", "x*x+2*x-1", "x*x+2*x+10",
	"x*x*x+2*x+1", "x*x+2*x+10", "x*x*x*x*3x+1",
}

var AllRanges = []int{5, 10, 25, 50, 100, 500, 1000}
var AllSeed = []int{0, 5, 50, 100, 250, 500, 1000}
var AllGenerationsCount = []int{50, 75, 100, 150, 250, 500}
var AllEachPopulationSize = []int{100, 250, 500, 1000, 2000}
var AllReproduction = []evolution.Reproduction{
	evolution.Reproduction{ProbabilityOfMutation: 0.1, CrossoverPercentage: 0.1},
	evolution.Reproduction{ProbabilityOfMutation: 0.1, CrossoverPercentage: 0.2},
	evolution.Reproduction{ProbabilityOfMutation: 0.1, CrossoverPercentage: 0.3},
	evolution.Reproduction{ProbabilityOfMutation: 0.1, CrossoverPercentage: 0.4},
	evolution.Reproduction{ProbabilityOfMutation: 0.1, CrossoverPercentage: 0.5},
	evolution.Reproduction{ProbabilityOfMutation: 0.1, CrossoverPercentage: 0.6},
	evolution.Reproduction{ProbabilityOfMutation: 0.1, CrossoverPercentage: 0.7},
	evolution.Reproduction{ProbabilityOfMutation: 0.1, CrossoverPercentage: 0.8},
	evolution.Reproduction{ProbabilityOfMutation: 0.1, CrossoverPercentage: 0.9},
	evolution.Reproduction{ProbabilityOfMutation: 0.1, CrossoverPercentage: 1.0},
}

var AllDepthOfRandomNewTree = []int{1, 2, 5, 10}
var AllAntagonistStrategyCount = []int{2, 5, 10, 15, 20, 30, 50, 75, 100, 250}
var AllProtagonistStrategyCount = []int{2, 5, 10, 15, 20, 30, 50, 75, 100, 250}

var AllFitnessStrategyType = []string{evolution.FitnessThresholdedAntagonistRatio,
	evolution.FitnessDualThresholdedRatio, evolution.FitnessRatio}
var AllFitStratAntThreshMult = []float64{5, 10, 25, 50, 100, 250}
var AllFitStratProThreshMult = []float64{1, 1.1, 1.25, 1.5, 1.8, 2.5, 5}

var AllSelectionParentType = []int{evolution.ParentSelectionTournament, evolution.ParentSelectionElitism}
var AllSelectionSurvivorPercentage = []float64{0.1, 0.2, 0.4, 0.5, 0.8, 0.9, 1.0}

var AllPossibleStrategies = [][]evolution.Strategy{AllStrategies, AllStrategiesNoDelete, AllStrategiesNoFell,
	AllStrategiesNoSkip, AllStrategiesNoX, AllStrategiesX, AllStrategiesNoAddRandom, AllStrategiesNoMutate}
var AllStrategies = []evolution.Strategy{
	evolution.StrategyDeleteMalicious,
	evolution.StrategyDeleteNonTerminal,
	evolution.StrategyDeleteTerminal,
	evolution.StrategyMutateNonTerminal,
	evolution.StrategyMutateTerminal,
	evolution.StrategyReplaceBranch,
	evolution.StrategyReplaceBranchX,
	evolution.StrategyAddRandomSubTree,
	evolution.StrategyAddRandomSubTreeX,
	evolution.StrategyAddToLeaf,
	evolution.StrategyAddMult,
	evolution.StrategyAddMultX,
	evolution.StrategyAddSub,
	evolution.StrategyAddSubX,
	evolution.StrategyAddAdd,
	evolution.StrategyAddAddX,
	evolution.StrategySkip,
	evolution.StrategyFellTree,
}
var AllStrategiesNoDelete = []evolution.Strategy{
	evolution.StrategyMutateNonTerminal,
	evolution.StrategyMutateTerminal,
	evolution.StrategyReplaceBranch,
	evolution.StrategyReplaceBranchX,
	evolution.StrategyAddRandomSubTree,
	evolution.StrategyAddRandomSubTreeX,
	evolution.StrategyAddToLeaf,
	evolution.StrategyAddMult,
	evolution.StrategyAddMultX,
	evolution.StrategyAddSub,
	evolution.StrategyAddSubX,
	evolution.StrategyAddAdd,
	evolution.StrategyAddAddX,
	evolution.StrategySkip,
}
var AllStrategiesNoX = []evolution.Strategy{
	evolution.StrategyDeleteMalicious,
	evolution.StrategyDeleteNonTerminal,
	evolution.StrategyDeleteTerminal,
	evolution.StrategyMutateNonTerminal,
	evolution.StrategyMutateTerminal,
	evolution.StrategyReplaceBranch,
	evolution.StrategyAddRandomSubTree,
	evolution.StrategyAddToLeaf,
	evolution.StrategyAddMult,
	evolution.StrategyAddSub,
	evolution.StrategyAddAdd,
	evolution.StrategySkip,
	evolution.StrategyFellTree,
}
var AllStrategiesX = []evolution.Strategy{
	evolution.StrategyDeleteMalicious,
	evolution.StrategyDeleteNonTerminal,
	evolution.StrategyDeleteTerminal,
	evolution.StrategyMutateNonTerminal,
	evolution.StrategyMutateTerminal,
	evolution.StrategyReplaceBranchX,
	evolution.StrategyAddRandomSubTreeX,
	evolution.StrategyAddMultX,
	evolution.StrategyAddSubX,
	evolution.StrategyAddAddX,
	evolution.StrategySkip,
	evolution.StrategyFellTree,
}
var AllStrategiesNoSkip = []evolution.Strategy{
	evolution.StrategyDeleteMalicious,
	evolution.StrategyDeleteNonTerminal,
	evolution.StrategyDeleteTerminal,
	evolution.StrategyMutateNonTerminal,
	evolution.StrategyMutateTerminal,
	evolution.StrategyReplaceBranch,
	evolution.StrategyReplaceBranchX,
	evolution.StrategyAddRandomSubTree,
	evolution.StrategyAddRandomSubTreeX,
	evolution.StrategyAddToLeaf,
	evolution.StrategyAddMult,
	evolution.StrategyAddMultX,
	evolution.StrategyAddSub,
	evolution.StrategyAddSubX,
	evolution.StrategyAddAdd,
	evolution.StrategyAddAddX,
	evolution.StrategyFellTree,
}
var AllStrategiesNoFell = []evolution.Strategy{
	evolution.StrategyDeleteMalicious,
	evolution.StrategyDeleteNonTerminal,
	evolution.StrategyDeleteTerminal,
	evolution.StrategyMutateNonTerminal,
	evolution.StrategyMutateTerminal,
	evolution.StrategyReplaceBranch,
	evolution.StrategyReplaceBranchX,
	evolution.StrategyAddRandomSubTree,
	evolution.StrategyAddRandomSubTreeX,
	evolution.StrategyAddToLeaf,
	evolution.StrategyAddMult,
	evolution.StrategyAddMultX,
	evolution.StrategyAddSub,
	evolution.StrategyAddSubX,
	evolution.StrategyAddAdd,
	evolution.StrategyAddAddX,
	evolution.StrategySkip,
}
var AllStrategiesNoMutate = []evolution.Strategy{
	evolution.StrategyDeleteMalicious,
	evolution.StrategyDeleteNonTerminal,
	evolution.StrategyDeleteTerminal,
	evolution.StrategyReplaceBranch,
	evolution.StrategyReplaceBranchX,
	evolution.StrategyAddRandomSubTree,
	evolution.StrategyAddRandomSubTreeX,
	evolution.StrategyAddToLeaf,
	evolution.StrategyAddMult,
	evolution.StrategyAddMultX,
	evolution.StrategyAddSub,
	evolution.StrategyAddSubX,
	evolution.StrategyAddAdd,
	evolution.StrategyAddAddX,
	evolution.StrategySkip,
	evolution.StrategyFellTree,
}
var AllStrategiesNoAddRandom = []evolution.Strategy{
	evolution.StrategyDeleteMalicious,
	evolution.StrategyDeleteNonTerminal,
	evolution.StrategyDeleteTerminal,
	evolution.StrategyMutateNonTerminal,
	evolution.StrategyMutateTerminal,
	evolution.StrategyReplaceBranch,
	evolution.StrategyReplaceBranchX,
	evolution.StrategyAddMult,
	evolution.StrategyAddMultX,
	evolution.StrategyAddSub,
	evolution.StrategyAddSubX,
	evolution.StrategyAddAdd,
	evolution.StrategyAddAddX,
	evolution.StrategySkip,
	evolution.StrategyFellTree,
}

//// ParamNoDeleteStd contains no Delete Operations
//var ParamNoDeleteStd = evolution.EvolutionParams{
//	StatisticsOutput: evolution.StatisticsOutput{
//		OutputPath: "",
//	},
//	SpecParam: evolution.SpecParam{
//		Range:      5,
//		Expression: "x*x",
//		Seed:       0,
//		AvailableVariablesAndOperators: evolution.AvailableVariablesAndOperators{
//			Constants: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
//			Variables: []string{"x"},
//			Operators: []string{"*", "+", "-"},
//		},
//	},
//	GenerationsCount:   50,
//	EachPopulationSize: 500, // Must be an even number to prevent awkward ordering of children.
//	FitnessStrategy: evolution.FitnessStrategy{
//		Type:                           evolution.FitnessThresholdedAntagonistRatio,
//		IsMoreFitnessBetter:            false,
//		AntagonistThresholdMultiplier:  40,
//		ProtagonistThresholdMultiplier: 1.2,
//	},
//	Selection: evolution.Selection{
//		Parent: evolution.ParentSelection{
//			Type:           evolution.ParentSelectionTournament,
//			TournamentSize: 3,
//		},
//		Survivor: evolution.SurvivorSelection{
//			Type:               1,
//			SurvivorPercentage: 0.2,
//		},
//	},
//
//	Reproduction: evolution.Reproduction{
//		ProbabilityOfMutation: 0.01,
//		CrossoverPercentage:   0.5,
//	},
//
//	Strategies: evolution.Strategies{
//		ProtagonistAvailableStrategies: []evolution.Strategy{
//			evolution.StrategyMutateNonTerminal,
//			evolution.StrategyMutateTerminal,
//			evolution.StrategyReplaceBranch,
//			evolution.StrategyReplaceBranchX,
//			evolution.StrategyAddRandomSubTree,
//			evolution.StrategyAddRandomSubTreeX,
//			evolution.StrategyAddToLeaf,
//			evolution.StrategyAddMult,
//			evolution.StrategyAddMultX,
//			evolution.StrategyAddSub,
//			evolution.StrategyAddSubX,
//			evolution.StrategyAddAdd,
//			evolution.StrategyAddAddX,
//			evolution.StrategySkip,
//		},
//		AntagonistAvailableStrategies: []evolution.Strategy{
//			evolution.StrategyMutateNonTerminal,
//			evolution.StrategyMutateTerminal,
//			evolution.StrategyReplaceBranch,
//			evolution.StrategyReplaceBranchX,
//			evolution.StrategyAddRandomSubTree,
//			evolution.StrategyAddRandomSubTreeX,
//			evolution.StrategyAddToLeaf,
//			evolution.StrategyAddMult,
//			evolution.StrategyAddMultX,
//			evolution.StrategyAddSub,
//			evolution.StrategyAddSubX,
//			evolution.StrategyAddAdd,
//			evolution.StrategyAddAddX,
//			evolution.StrategySkip,
//		},
//		AntagonistStrategyCount:  5,
//		ProtagonistStrategyCount: 5,
//		DepthOfRandomNewTrees:    1,
//	},
//	FitnessCalculatorType: 0,
//	//ShouldRunInteractiveTerminal: shouldRunInteractive,
//}
//
//// ParamNoFellStrat contains no Delete Operations
//var ParamNoFellStrat = evolution.EvolutionParams{
//	StatisticsOutput: evolution.StatisticsOutput{
//		OutputPath: "",
//	},
//	SpecParam: evolution.SpecParam{
//		Range:      10,
//		Expression: "x*x",
//		Seed:       0,
//		AvailableVariablesAndOperators: evolution.AvailableVariablesAndOperators{
//			Constants: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
//			Variables: []string{"x"},
//			Operators: []string{"*", "+", "-"},
//		},
//	},
//	GenerationsCount:   50,
//	EachPopulationSize: 500, // Must be an even number to prevent awkward ordering of children.
//
//	FitnessStrategy: evolution.FitnessStrategy{
//		Type:                           evolution.FitnessThresholdedAntagonistRatio,
//		IsMoreFitnessBetter:            false,
//		AntagonistThresholdMultiplier:  40,
//		ProtagonistThresholdMultiplier: 1.2,
//	},
//
//	Selection: evolution.Selection{
//		Parent: evolution.ParentSelection{
//			Type:           evolution.ParentSelectionTournament,
//			TournamentSize: 3,
//		},
//		Survivor: evolution.SurvivorSelection{
//			Type:               1,
//			SurvivorPercentage: 0.2,
//		},
//	},
//
//	Reproduction: evolution.Reproduction{
//		ProbabilityOfMutation: 0.01,
//		CrossoverPercentage:   0.5,
//	},
//
//	Strategies: evolution.Strategies{
//		ProtagonistAvailableStrategies: []evolution.Strategy{
//			evolution.StrategyDeleteMalicious,
//			evolution.StrategyDeleteNonTerminal,
//			evolution.StrategyDeleteTerminal,
//			evolution.StrategyMutateNonTerminal,
//			evolution.StrategyMutateTerminal,
//			evolution.StrategyReplaceBranch,
//			evolution.StrategyReplaceBranchX,
//			evolution.StrategyAddRandomSubTree,
//			evolution.StrategyAddRandomSubTreeX,
//			evolution.StrategyAddToLeaf,
//			evolution.StrategyAddMult,
//			evolution.StrategyAddMultX,
//			evolution.StrategyAddSub,
//			evolution.StrategyAddSubX,
//			evolution.StrategyAddAdd,
//			evolution.StrategyAddAddX,
//			evolution.StrategySkip,
//		},
//		AntagonistAvailableStrategies: []evolution.Strategy{
//			evolution.StrategyDeleteMalicious,
//			evolution.StrategyDeleteNonTerminal,
//			evolution.StrategyDeleteTerminal,
//			evolution.StrategyMutateNonTerminal,
//			evolution.StrategyMutateTerminal,
//			evolution.StrategyReplaceBranch,
//			evolution.StrategyReplaceBranchX,
//			evolution.StrategyAddRandomSubTree,
//			evolution.StrategyAddRandomSubTreeX,
//			evolution.StrategyAddToLeaf,
//			evolution.StrategyAddMult,
//			evolution.StrategyAddMultX,
//			evolution.StrategyAddSub,
//			evolution.StrategyAddSubX,
//			evolution.StrategyAddAdd,
//			evolution.StrategyAddAddX,
//			evolution.StrategySkip,
//		},
//		AntagonistStrategyCount:  5,
//		ProtagonistStrategyCount: 5,
//		DepthOfRandomNewTrees:    1,
//	},
//	FitnessCalculatorType: 0,
//	//ShouldRunInteractiveTerminal: shouldRunInteractive,
//}
//
//// ParamNoFellStrat contains no Delete Operations
//var ParamAllFellStrat = evolution.EvolutionParams{
//	StatisticsOutput: evolution.StatisticsOutput{
//		OutputPath: "",
//	},
//	SpecParam: evolution.SpecParam{
//		Range:      10,
//		Expression: "x*x",
//		Seed:       0,
//		AvailableVariablesAndOperators: evolution.AvailableVariablesAndOperators{
//			Constants: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
//			Variables: []string{"x"},
//			Operators: []string{"*", "+", "-"},
//		},
//	},
//	GenerationsCount:   50,
//	EachPopulationSize: 500, // Must be an even number to prevent awkward ordering of children.
//
//	FitnessStrategy: evolution.FitnessStrategy{
//		Type:                           evolution.FitnessThresholdedAntagonistRatio,
//		IsMoreFitnessBetter:            false,
//		AntagonistThresholdMultiplier:  40,
//		ProtagonistThresholdMultiplier: 1.2,
//	},
//
//	Selection: evolution.Selection{
//		Parent: evolution.ParentSelection{
//			Type:           evolution.ParentSelectionTournament,
//			TournamentSize: 3,
//		},
//		Survivor: evolution.SurvivorSelection{
//			Type:               1,
//			SurvivorPercentage: 0.2,
//		},
//	},
//
//	Reproduction: evolution.Reproduction{
//		ProbabilityOfMutation: 0.01,
//		CrossoverPercentage:   0.5,
//	},
//
//	Strategies: evolution.Strategies{
//		ProtagonistAvailableStrategies: []evolution.Strategy{
//			evolution.StrategyDeleteMalicious,
//			evolution.StrategyDeleteNonTerminal,
//			evolution.StrategyDeleteTerminal,
//			evolution.StrategyFellTree,
//			evolution.StrategyMutateNonTerminal,
//			evolution.StrategyMutateTerminal,
//			evolution.StrategyReplaceBranch,
//			evolution.StrategyReplaceBranchX,
//			evolution.StrategyAddRandomSubTree,
//			evolution.StrategyAddRandomSubTreeX,
//			evolution.StrategyAddToLeaf,
//			evolution.StrategyAddMult,
//			evolution.StrategyAddMultX,
//			evolution.StrategyAddSub,
//			evolution.StrategyAddSubX,
//			evolution.StrategyAddAdd,
//			evolution.StrategyAddAddX,
//			evolution.StrategySkip,
//		},
//		AntagonistAvailableStrategies: []evolution.Strategy{
//			evolution.StrategyFellTree,
//			evolution.StrategyDeleteMalicious,
//			evolution.StrategyDeleteNonTerminal,
//			evolution.StrategyDeleteTerminal,
//			evolution.StrategyMutateNonTerminal,
//			evolution.StrategyMutateTerminal,
//			evolution.StrategyReplaceBranch,
//			evolution.StrategyReplaceBranchX,
//			evolution.StrategyAddRandomSubTree,
//			evolution.StrategyAddRandomSubTreeX,
//			evolution.StrategyAddToLeaf,
//			evolution.StrategyAddMult,
//			evolution.StrategyAddMultX,
//			evolution.StrategyAddSub,
//			evolution.StrategyAddSubX,
//			evolution.StrategyAddAdd,
//			evolution.StrategyAddAddX,
//			evolution.StrategySkip,
//		},
//		AntagonistStrategyCount:  5,
//		ProtagonistStrategyCount: 5,
//		DepthOfRandomNewTrees:    1,
//	},
//	FitnessCalculatorType: 0,
//	//ShouldRunInteractiveTerminal: shouldRunInteractive,
//}
