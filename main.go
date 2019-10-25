package main

import (
	"fmt"
	"github.com/martinomburajr/masters-go/eval"
	"github.com/martinomburajr/masters-go/evolution"
	"log"
	"math/rand"
	"os/exec"
	"time"
)

func main() {
	name := "run.json"

	// TODO Include terminals and non terminals as part of strategy?
	params := evolution.EvolutionParams{
		StatisticsOutput: evolution.StatisticsOutput{
			OutputPath: "",
		},
		SpecParam: evolution.SpecParam{
			Range:      5,
			Expression: "x*x",
			Seed:       0,
			AvailableVariablesAndOperators: evolution.AvailableVariablesAndOperators{
				Constants: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
				Variables: []string{"x"},
				Operators: []string{"*", "+", "-"},
			},
		},
		GenerationsCount:   50,
		EachPopulationSize: 20, // Must be an even number to prevent awkward ordering of children.

		FitnessStrategy: evolution.FitnessStrategy{
			Type:                           evolution.FitnessThresholdedAntagonistRatio,
			IsMoreFitnessBetter:            false,
			AntagonistThresholdMultiplier:  40,
			ProtagonistThresholdMultiplier: 1.2,
		},

		Selection: evolution.Selection{
			Parent: evolution.ParentSelection{
				Type:           evolution.ParentSelectionTournament,
				TournamentSize: 3,
			},
			Survivor: evolution.SurvivorSelection{
				Type:               1,
				SurvivorPercentage: 0.2,
			},
		},

		Reproduction: evolution.Reproduction{
			ProbabilityOfMutation: 0.01,
			CrossoverPercentage:   0.5,
		},

		Strategies: evolution.Strategies{
			ProtagonistAvailableStrategies: []evolution.Strategy{
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
			},
			AntagonistAvailableStrategies: []evolution.Strategy{
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
			},
			AntagonistStrategyCount:  5,
			ProtagonistStrategyCount: 5,
			DepthOfRandomNewTrees:    1,
		},
		FitnessCalculatorType: 0,
		//ShouldRunInteractiveTerminal: shouldRunInteractive,
	}

	simulation := Simulation{
		EvolutionStates:      []evolution.EvolutionParams{params},
		NumberOfRunsPerState: 10,
		Name:                 "simulation-1",
	}

	simulation.Begin()
	cmd := exec.Command("Rscript", "launch.R")
	log.Fatal(cmd.Run())
}
