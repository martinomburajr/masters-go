package simulation

import "github.com/martinomburajr/masters-go/evolution"

var AllMaxGenerations = []int{500}
var AllExpressions = []string{
	"x",
	"x*x*x*x*x*x*x*x",
	"x*x*x+2*x/3*x*x+5",
}
var AllTopologies = []evolution.Topology{
	{
		Type:                  evolution.TopologyHallOfFame,
		HoFGenerationInterval: 0.1,
	},
	{
		Type: evolution.TopologySingleEliminationTournament,
		SETNoOfTournaments: 0.2,
	},
	{
		Type: evolution.TopologyRoundRobin,
	},
	{
		Type:     evolution.TopologyKRandom,
		KRandomK: 3,
	},
}

var AllReproduction = []evolution.Reproduction{
	{
		CrossoverStrategy:     evolution.CrossoverSinglePoint,
		ProbabilityOfMutation: 0.3,
	},
	{
		CrossoverStrategy:     evolution.CrossoverUniform,
		ProbabilityOfMutation: 0.3,
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
var AllMinimumGenerationMeanBeforeTerminate = []float64{0.05}
var AllMinimumTopProtagonistMeanBeforeTerminate = []float64{0.1}
var AllProtagonistMinGenAvgFit = []float64{0.7}
var AllRanges = []int{20}
var AllSeed = []int{-10}
var AllGenerationsCount = []int{50}
var AllEachPopulationSize = []int{64}

var AllDepthOfRandomNewTree = []int{1}
var AllAntagonistStrategyCount = []int{16}
var AllProtagonistStrategyCount = []int{16}

var AllFitnessStrategyType = []string{evolution.FitnessDualThresholdedRatio}
var AllFitStratAntThreshMult = []float64{16}
var AllFitStratProThreshMult = []float64{1}

var AllSelectionParentType = []evolution.ParentSelection{
	{
		Type:           evolution.ParentSelectionTournament,
		TournamentSize: 3,
	},
}

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
