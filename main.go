package main

import (
	"fmt"
	"github.com/martinomburajr/masters-go/evolution"
	"log"
)

func main() {
	strategies :=  []evolution.Strategy{evolution.StrategyAddSubTree, evolution.StrategyDeleteSubTree,
		evolution.StrategyMutateSubTree, evolution.StrategyMutateNode}

	spec := evolution.Spec2X

	params := evolution.EvolutionParams{
		EachPopulationSize: 100,
		AntagonistMaxStrategies: 4,
		AntagonistStrategyLength: 3,
		ProtagonistMaxStrategies: 4,
		ProtagonistStrategyLength: 3,
		MaxDepth: 3,
		DepthPenaltyStrategyPenalization: 10,
	}

	engine := evolution.EvolutionEngine{
		StartIndividual: evolution.ProgTreeT_NT_T_0,
		Spec: spec,
		AvailableStrategies: strategies,

		GenerationCount: 50,

		AvailableNonTerminalSet: []evolution.SymbolicExpression{evolution.Add, evolution.Mult, evolution.Sub},
		AvailableTerminalSet: []evolution.SymbolicExpression{
			evolution.Const1, evolution.Const2, evolution.Const3, evolution.Const4, evolution.Const5, evolution.X1,
		},


		FitnessStrategy:     evolution.FitnessProtagonistThresholdTally,
		EvaluationThreshold: 0.1,
		EvaluationMinThreshold: 0.01,

		ParentSelection:  evolution.ParentSelectionTournament,
		TournamentSize:   3,
		Parameters: params,
	}

	result, err := engine.Start()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(result)
}

