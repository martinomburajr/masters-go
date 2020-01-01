package simulation

import "github.com/martinomburajr/masters-go/evolution"

var AllExpressions = []string{
	"x",
	"x*x*x+2*x/3*x+9",
}
var AllRanges = []int{30}
var AllSeed = []int{-10}
var AllGenerationsCount = []int{50}
var AllEachPopulationSize = []int{250}

var AllReproduction = []evolution.Reproduction{
	evolution.Reproduction{ProbabilityOfMutation: 0.1, CrossoverPercentage: 0.3},
}
var AllDepthOfRandomNewTree = []int{0, 3}
var AllAntagonistStrategyCount = []int{15}
var AllProtagonistStrategyCount = []int{15}

var AllFitnessStrategyType = []string{evolution.FitnessDualThresholdedRatio}
var AllFitStratAntThreshMult = []float64{50}
var AllFitStratProThreshMult = []float64{1}

var AllSelectionParentType = []string{evolution.ParentSelectionTournament}
var AllTournamentSizesType = []int{3}
var AllSelectionSurvivorPercentage = []float64{0.2, 0.7}

var AllDivByZeroStrategy = []string{
	evolution.DivByZeroIgnore,
	evolution.DivByZeroPenalize,
}
var AllDivByZeroPenalty = []float64{-2}

var AllPossibleStrategies = [][]evolution.Strategy{
	AllStrategies,
	AllStrategiesDeterministic,
	AllStrategiesRandom,
	//AllStrategiesNoDelete,
	// AllStrategiesNoFell,
	//AllStrategiesNoSkip,
	// AllStrategiesNoX,
	// AllStrategiesX,
	// AllStrategiesNoAddRandom,
	// AllStrategiesNoMutate
}
var AllStrategies = []evolution.Strategy{
	evolution.StrategyDeleteMalicious,
	evolution.StrategyDeleteNonTerminal,
	evolution.StrategyDeleteTerminal,
	evolution.StrategyMutateNonTerminal,
	evolution.StrategyMutateTerminal,
	evolution.StrategyReplaceBranch,
	evolution.StrategyReplaceBranchX,
	evolution.StrategyAddRandomSubTree,
	evolution.StrategyAddToLeaf,
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
var AllStrategiesDeterministic = []evolution.Strategy{
	evolution.StrategySkip,
	evolution.StrategyFellTree,
	evolution.StrategyMultXD,
	evolution.StrategyAddXD,
	evolution.StrategySubXD,
	evolution.StrategyDivXD,
	evolution.StrategyAddTreeWithDiv,
}

var AllStrategiesRandom = []evolution.Strategy{
	evolution.StrategyDeleteMalicious,
	evolution.StrategyDeleteNonTerminal,
	evolution.StrategyDeleteTerminal,
	evolution.StrategyMutateNonTerminal,
	evolution.StrategyMutateTerminal,
	evolution.StrategyReplaceBranch,
	evolution.StrategyReplaceBranchX,
	evolution.StrategyAddRandomSubTree,
	evolution.StrategyAddToLeaf,
	evolution.StrategyAddTreeWithMult,
	evolution.StrategyAddTreeWithSub,
	evolution.StrategyAddTreeWithAdd,
}
var AllStrategiesNoDelete = []evolution.Strategy{
	evolution.StrategyMutateNonTerminal,
	evolution.StrategyMutateTerminal,
	evolution.StrategyReplaceBranch,
	evolution.StrategyReplaceBranchX,
	evolution.StrategyAddRandomSubTree,
	evolution.StrategyAddToLeaf,
	evolution.StrategyAddTreeWithMult,
	evolution.StrategyAddTreeWithSub,
	evolution.StrategyAddTreeWithAdd,
	evolution.StrategySkip,
	evolution.StrategyMultXD,
	evolution.StrategyAddXD,
	evolution.StrategySubXD,
	evolution.StrategyDivXD,
	evolution.StrategyAddTreeWithDiv,
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
	evolution.StrategyAddTreeWithMult,
	evolution.StrategyAddTreeWithSub,
	evolution.StrategyAddTreeWithAdd,
	evolution.StrategySkip,
	evolution.StrategyFellTree,
	evolution.StrategyAddTreeWithDiv,
}
var AllStrategiesX = []evolution.Strategy{
	evolution.StrategyDeleteMalicious,
	evolution.StrategyDeleteNonTerminal,
	evolution.StrategyDeleteTerminal,
	evolution.StrategyMutateNonTerminal,
	evolution.StrategyMutateTerminal,
	evolution.StrategyReplaceBranchX,
	evolution.StrategySkip,
	evolution.StrategyFellTree,
	evolution.StrategyMultXD,
	evolution.StrategyAddXD,
	evolution.StrategySubXD,
	evolution.StrategyDivXD,
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
	evolution.StrategyAddToLeaf,
	evolution.StrategyAddTreeWithMult,
	evolution.StrategyAddTreeWithSub,
	evolution.StrategyAddTreeWithAdd,
	evolution.StrategyFellTree,
	evolution.StrategyMultXD,
	evolution.StrategyAddXD,
	evolution.StrategySubXD,
	evolution.StrategyDivXD,
	evolution.StrategyAddTreeWithDiv,
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
	evolution.StrategyAddToLeaf,
	evolution.StrategyAddTreeWithMult,
	evolution.StrategyAddTreeWithSub,
	evolution.StrategyAddTreeWithAdd,
	evolution.StrategySkip,
	evolution.StrategyMultXD,
	evolution.StrategyAddXD,
	evolution.StrategySubXD,
	evolution.StrategyDivXD,
	evolution.StrategyAddTreeWithDiv,
}
var AllStrategiesNoMutate = []evolution.Strategy{
	evolution.StrategyDeleteMalicious,
	evolution.StrategyDeleteNonTerminal,
	evolution.StrategyDeleteTerminal,
	evolution.StrategyReplaceBranch,
	evolution.StrategyReplaceBranchX,
	evolution.StrategyAddRandomSubTree,
	evolution.StrategyAddToLeaf,
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
var AllStrategiesNoAddRandom = []evolution.Strategy{
	evolution.StrategyDeleteMalicious,
	evolution.StrategyDeleteNonTerminal,
	evolution.StrategyDeleteTerminal,
	evolution.StrategyMutateNonTerminal,
	evolution.StrategyMutateTerminal,
	evolution.StrategyReplaceBranch,
	evolution.StrategyReplaceBranchX,
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

