package simulation

import "github.com/martinomburajr/masters-go/evolution"

var AllMaxGenerations = []int{500}
var AllExpressions = []string{
	"x",
	"1/1*x*x",
	"x*x*x*x*x*x*x*x*x",
	"x*x*x+2*x/3*x*x+5",
}
var AllReproduction = []evolution.Reproduction{
	{
		CrossoverStrategy:     evolution.CrossoverSinglePoint,
		ProbabilityOfMutation: 0.6,
	},
	{
		CrossoverStrategy:     evolution.CrossoverUniform,
		ProbabilityOfMutation: 0.6,
	},
	{
		CrossoverStrategy:     evolution.CrossoverSinglePoint,
		ProbabilityOfMutation: 0.1,
	},
	{
		CrossoverStrategy:     evolution.CrossoverUniform,
		ProbabilityOfMutation: 0.1,
	},
}
var AllSurvivorSelection = []evolution.SurvivorSelection{
	{
		Type:               evolution.SurvivorSelectionFitnessBased,
		SurvivorPercentage: 0.3,
	},
	{
		Type:               evolution.SurvivorSelectionFitnessBased,
		SurvivorPercentage: 0.7,
	},
}
var AllRanges = []int{20}
var AllSeed = []int{-10}
var AllGenerationsCount = []int{50}
var AllEachPopulationSize = []int{50,100}

var AllDepthOfRandomNewTree = []int{1}
var AllAntagonistStrategyCount = []int{15}
var AllProtagonistStrategyCount = []int{15}

var AllFitnessStrategyType = []string{evolution.FitnessDualThresholdedRatio}
var AllFitStratAntThreshMult = []float64{10}
var AllFitStratProThreshMult = []float64{1}

var AllSelectionParentType = []string{evolution.ParentSelectionTournament}
var AllTournamentSizesType = []int{3}


var AllDivByZeroStrategy = []string{
	evolution.DivByZeroSteadyPenalize,
}
var AllDivByZeroPenalty = []float64{-1}

var AllPossibleStrategies = [][]evolution.Strategy{
	AllStrategies,
}

var AllStrategies = []evolution.Strategy{
	evolution.StrategyDeleteNonTerminal,
	evolution.StrategyDeleteTerminal,
	evolution.StrategyMutateNonTerminal,
	evolution.StrategyMutateTerminal,
	evolution.StrategyReplaceBranch,
	evolution.StrategyReplaceBranchX,
	evolution.StrategyAddRandomSubTree,
	evolution.StrategyAddToLeaf,
	evolution.StrategyAddToLeafX,
	evolution.StrategyAddTreeWithMult,
	evolution.StrategyAddTreeWithSub,
	evolution.StrategyAddTreeWithAdd,
	evolution.StrategySkip,
	evolution.StrategyFellTree,
	evolution.StrategyMultXD,
	evolution.StrategyAddXD,
	evolution.StrategySubXD,
	evolution.StrategyDivXD,
	evolution.StrategyAddTreeWithDiv,
}
//var AllStrategiesDeterministic = []evolution.Strategy{
//	evolution.StrategySkip,
//	evolution.StrategyFellTree,
//	evolution.StrategyMultXD,
//	evolution.StrategyAddXD,
//	evolution.StrategySubXD,
//	evolution.StrategyDivXD,
//	evolution.StrategyAddTreeWithDiv,
//}
//
//var AllStrategiesRandom = []evolution.Strategy{
//	evolution.StrategyDeleteMalicious,
//	evolution.StrategyDeleteNonTerminal,
//	evolution.StrategyDeleteTerminal,
//	evolution.StrategyMutateNonTerminal,
//	evolution.StrategyMutateTerminal,
//	evolution.StrategyReplaceBranch,
//	evolution.StrategyReplaceBranchX,
//	evolution.StrategyAddRandomSubTree,
//	evolution.StrategyAddToLeaf,
//	evolution.StrategyAddTreeWithMult,
//	evolution.StrategyAddTreeWithSub,
//	evolution.StrategyAddTreeWithAdd,
//}
//var AllStrategiesNoDelete = []evolution.Strategy{
//	evolution.StrategyMutateNonTerminal,
//	evolution.StrategyMutateTerminal,
//	evolution.StrategyReplaceBranch,
//	evolution.StrategyReplaceBranchX,
//	evolution.StrategyAddRandomSubTree,
//	evolution.StrategyAddToLeaf,
//	evolution.StrategyAddTreeWithMult,
//	evolution.StrategyAddTreeWithSub,
//	evolution.StrategyAddTreeWithAdd,
//	evolution.StrategySkip,
//	evolution.StrategyMultXD,
//	evolution.StrategyAddXD,
//	evolution.StrategySubXD,
//	evolution.StrategyDivXD,
//	evolution.StrategyAddTreeWithDiv,
//}
//var AllStrategiesNoX = []evolution.Strategy{
//	evolution.StrategyDeleteMalicious,
//	evolution.StrategyDeleteNonTerminal,
//	evolution.StrategyDeleteTerminal,
//	evolution.StrategyMutateNonTerminal,
//	evolution.StrategyMutateTerminal,
//	evolution.StrategyReplaceBranch,
//	evolution.StrategyAddRandomSubTree,
//	evolution.StrategyAddToLeaf,
//	evolution.StrategyAddTreeWithMult,
//	evolution.StrategyAddTreeWithSub,
//	evolution.StrategyAddTreeWithAdd,
//	evolution.StrategySkip,
//	evolution.StrategyFellTree,
//	evolution.StrategyAddTreeWithDiv,
//}
//var AllStrategiesX = []evolution.Strategy{
//	evolution.StrategyDeleteMalicious,
//	evolution.StrategyDeleteNonTerminal,
//	evolution.StrategyDeleteTerminal,
//	evolution.StrategyMutateNonTerminal,
//	evolution.StrategyMutateTerminal,
//	evolution.StrategyReplaceBranchX,
//	evolution.StrategySkip,
//	evolution.StrategyFellTree,
//	evolution.StrategyMultXD,
//	evolution.StrategyAddXD,
//	evolution.StrategySubXD,
//	evolution.StrategyDivXD,
//}
//var AllStrategiesNoSkip = []evolution.Strategy{
//	evolution.StrategyDeleteMalicious,
//	evolution.StrategyDeleteNonTerminal,
//	evolution.StrategyDeleteTerminal,
//	evolution.StrategyMutateNonTerminal,
//	evolution.StrategyMutateTerminal,
//	evolution.StrategyReplaceBranch,
//	evolution.StrategyReplaceBranchX,
//	evolution.StrategyAddRandomSubTree,
//	evolution.StrategyAddToLeaf,
//	evolution.StrategyAddTreeWithMult,
//	evolution.StrategyAddTreeWithSub,
//	evolution.StrategyAddTreeWithAdd,
//	evolution.StrategyFellTree,
//	evolution.StrategyMultXD,
//	evolution.StrategyAddXD,
//	evolution.StrategySubXD,
//	evolution.StrategyDivXD,
//	evolution.StrategyAddTreeWithDiv,
//}
//var AllStrategiesNoFell = []evolution.Strategy{
//	evolution.StrategyDeleteMalicious,
//	evolution.StrategyDeleteNonTerminal,
//	evolution.StrategyDeleteTerminal,
//	evolution.StrategyMutateNonTerminal,
//	evolution.StrategyMutateTerminal,
//	evolution.StrategyReplaceBranch,
//	evolution.StrategyReplaceBranchX,
//	evolution.StrategyAddRandomSubTree,
//	evolution.StrategyAddToLeaf,
//	evolution.StrategyAddTreeWithMult,
//	evolution.StrategyAddTreeWithSub,
//	evolution.StrategyAddTreeWithAdd,
//	evolution.StrategySkip,
//	evolution.StrategyMultXD,
//	evolution.StrategyAddXD,
//	evolution.StrategySubXD,
//	evolution.StrategyDivXD,
//	evolution.StrategyAddTreeWithDiv,
//}
//var AllStrategiesNoMutate = []evolution.Strategy{
//	evolution.StrategyDeleteMalicious,
//	evolution.StrategyDeleteNonTerminal,
//	evolution.StrategyDeleteTerminal,
//	evolution.StrategyReplaceBranch,
//	evolution.StrategyReplaceBranchX,
//	evolution.StrategyAddRandomSubTree,
//	evolution.StrategyAddToLeaf,
//	evolution.StrategyAddTreeWithMult,
//	evolution.StrategyAddTreeWithSub,
//	evolution.StrategyAddTreeWithAdd,
//	evolution.StrategySkip,
//	evolution.StrategyFellTree,
//	evolution.StrategyMultXD,
//	evolution.StrategyAddXD,
//	evolution.StrategySubXD,
//	evolution.StrategyDivXD,
//	evolution.StrategyAddTreeWithDiv,
//}
//var AllStrategiesNoAddRandom = []evolution.Strategy{
//	evolution.StrategyDeleteMalicious,
//	evolution.StrategyDeleteNonTerminal,
//	evolution.StrategyDeleteTerminal,
//	evolution.StrategyMutateNonTerminal,
//	evolution.StrategyMutateTerminal,
//	evolution.StrategyReplaceBranch,
//	evolution.StrategyReplaceBranchX,
//	evolution.StrategyAddTreeWithMult,
//	evolution.StrategyAddTreeWithSub,
//	evolution.StrategyAddTreeWithAdd,
//	evolution.StrategySkip,
//	evolution.StrategyFellTree,
//	evolution.StrategyMultXD,
//	evolution.StrategyAddXD,
//	evolution.StrategySubXD,
//	evolution.StrategyDivXD,
//	evolution.StrategyAddTreeWithDiv,
//}

