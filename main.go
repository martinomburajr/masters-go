package main

import (
	"github.com/martinomburajr/masters-go/evolution"
	"log"
	"math/rand"
	"time"
)

func main() {
	Evolution1()
}

func Evolution1() {
	rand.Seed(time.Now().UTC().UnixNano()) //Set seed

	strategies := []evolution.Strategy{
		evolution.StrategyAddSubTree,
		evolution.StrategyDeleteSubTree,
		//evolution.StrategyMutateSubTree,
		evolution.StrategyMutateNode}

	spec := evolution.SpecXBy5

	// TODO only perform parent selection on loser
	params := evolution.EvolutionParams{
		Generations:        50,
		EachPopulationSize: 20, // Must be an even number to prevent awkward ordering of children.
		// This will fail if odd.
		AntagonistMaxStrategies: 4,
		//AntagonistStrategyLength:              3,
		ProtagonistMaxStrategies: 4,
		//ProtagonistStrategyLength:             3,
		MaxDepth:                         10,
		DepthPenaltyStrategyPenalization: 10,
		ProbabilityOfMutation:            0.1,
		ProbabilityOfNonTerminalMutation: 0.1,
		DepthOfRandomNewTrees:            1,
		DeletionType:                     evolution.DeletionTypeSafe,
		EnforceIndependentVariable:       true,
		//Strategies:                            strategies,
		ProtagonistAvailableStrategies:        strategies,
		AntagonistAvailableStrategies:         strategies,
		SetEqualStrategyLength:                true,
		CrossoverPercentage:                   0.2,
		MaintainCrossoverGeneTransferEquality: true,
		NonTerminalSet:                        []evolution.SymbolicExpression{evolution.Add, evolution.Mult, evolution.Sub},
		TerminalSet: []evolution.SymbolicExpression{
			evolution.X1, evolution.Const1, evolution.Const2, evolution.Const3, evolution.Const4, evolution.Const5,
			evolution.Const6, evolution.Const7, evolution.Const8, evolution.Const9, evolution.Const0,
		},
		FitnessStrategy:        evolution.FitnessProtagonistThresholdTally,
		EvaluationThreshold:    10,
		EvaluationMinThreshold: 9,
		TournamentSize:         2,
		StrategyLengthLimit:    10,
		SurvivorPercentage:     0.5,
		StartIndividual:        evolution.ProgTreeXby5,
		Spec:                   spec,
		ParentSelection:        evolution.ParentSelectionTournament,
		EqualStrategiesLength:  3,
	}

	engine := evolution.EvolutionEngine{
		Parameters:  params,
		Generations: []*evolution.Generation{},
	}

	_, err := engine.Start()
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Print(result)
}
