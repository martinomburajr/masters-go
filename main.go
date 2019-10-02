package main

import (
	"fmt"
	"github.com/martinomburajr/masters-go/evolution"
	"log"
)

func main() {
	strategies := []evolution.Strategy{
		evolution.StrategyAddSubTree,
		evolution.StrategyDeleteSubTree,
		//evolution.StrategyMutateSubTree,
		evolution.StrategyMutateNode}

	spec := evolution.Spec2X

	params := evolution.EvolutionParams{
		Generations: 50,
		EachPopulationSize:                    6, // Must be an even number to prevent awkward ordering of children.
		// This will fail if odd.
		AntagonistMaxStrategies:               4,
		AntagonistStrategyLength:              3,
		ProtagonistMaxStrategies:              4,
		ProtagonistStrategyLength:             3,
		MaxDepth:                              4,
		DepthPenaltyStrategyPenalization:      10,
		ProbabilityOfMutation:                 0.1,
		ProbabilityOfNonTerminalMutation:      0.1,
		DepthOfRandomNewTrees:                 1,
		DeletionType:                          evolution.DeletionTypeSafe,
		EnforceIndependentVariable:            true,
		Strategies:                            strategies,
		CrossoverPercentage:                   0.2,
		MaintainCrossoverGeneTransferEquality: true,
		NonTerminalSet: []evolution.SymbolicExpression{evolution.Add, evolution.Mult, evolution.Sub},
		TerminalSet: []evolution.SymbolicExpression{
			evolution.X1, evolution.Const1, evolution.Const2, evolution.Const3, evolution.Const4, evolution.Const5,
			evolution.Const6, evolution.Const7, evolution.Const8, evolution.Const9, evolution.Const0,
		},
		FitnessStrategy:        evolution.FitnessProtagonistThresholdTally,
		EvaluationThreshold:    0.1,
		EvaluationMinThreshold: 0.01,
		TournamentSize:  3,
		ParentSelection: evolution.ParentSelectionTournament,
		StrategyLengthLimit: 10,
		SurvivorPercentage: 0.5,
	}

	engine := evolution.EvolutionEngine{
		StartIndividual: evolution.ProgTreeT_NT_T_0,
		Spec:            spec,
		ParentSelection: evolution.ParentSelectionTournament,
		Parameters:      params,
	}

	result, err := engine.Start()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(result)
}
