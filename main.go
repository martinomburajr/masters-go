package main

import (
	"fmt"
	"github.com/martinomburajr/masters-go/evolution"
	"github.com/martinomburajr/masters-go/simulation"
	"log"
	"path/filepath"
)

func main() {

	absolutePath, err := filepath.Abs(".")
	if err != nil {
		log.Println(err)
		return
	}
	simulation := simulation.Simulation{
		NumberOfRunsPerState: 5,
		Name:                 "simulation-1",
		OutputDir:            "",
		RPath:                fmt.Sprintf("%s%s", absolutePath, "/R/runScript.R"),
	}
	params := evolution.EvolutionParams{
		StatisticsOutput: evolution.StatisticsOutput{
			OutputPath: "",
		},
		SpecParam: evolution.SpecParam{
			Range:      20,
			Expression: "x*x",
			Seed:       -10,
			AvailableVariablesAndOperators: evolution.AvailableVariablesAndOperators{
				Constants: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
				Variables: []string{"x"},
				Operators: []string{"*", "+", "-", "/"},
			},
			DivideByZeroStrategy: evolution.DivByZeroPenalize,
			DivideByZeroPenalty:  -2,
		},
		GenerationsCount:   50,
		EachPopulationSize: 4, // Must be an even number to prevent awkward ordering of children.

		FitnessStrategy: evolution.FitnessStrategy{
			Type:                           evolution.FitnessDualThresholdedRatio,
			AntagonistThresholdMultiplier:  30,
			ProtagonistThresholdMultiplier: 1,
		},

		Selection: evolution.Selection{
			Parent: evolution.ParentSelection{
				Type:           evolution.ParentSelectionTournament,
				TournamentSize: 3,
			},
			Survivor: evolution.SurvivorSelection{
				Type:               "SteadyState",
				SurvivorPercentage: 0.5,
			},
		},
		Reproduction: evolution.Reproduction{
			ProbabilityOfMutation: 0.01,
			CrossoverPercentage:   0.2,
		},
		Strategies: evolution.Strategies{
			ProtagonistAvailableStrategies: []evolution.Strategy{
				evolution.StrategyMutateNonTerminal,
				evolution.StrategyMutateTerminal,
				evolution.StrategyReplaceBranch,
				evolution.StrategyReplaceBranchX,
				evolution.StrategyAddRandomSubTree,
				evolution.StrategyAddToLeaf,
				evolution.StrategyAddTreeWithMult,
				evolution.StrategyAddTreeWithSub,
				evolution.StrategyAddTreeWithAdd,
				evolution.StrategyAddTreeWithDiv,
				evolution.StrategySkip,
				evolution.StrategyMultXD,
				evolution.StrategyAddXD,
				evolution.StrategySubXD,
				evolution.StrategyDivXD,
			},
			AntagonistAvailableStrategies: []evolution.Strategy{
				evolution.StrategyMutateNonTerminal,
				evolution.StrategyMutateTerminal,
				evolution.StrategyReplaceBranch,
				evolution.StrategyReplaceBranchX,
				evolution.StrategyAddRandomSubTree,
				evolution.StrategyAddToLeaf,
				evolution.StrategyAddTreeWithMult,
				evolution.StrategyAddTreeWithSub,
				evolution.StrategyAddTreeWithAdd,
				evolution.StrategyAddTreeWithDiv,
				evolution.StrategySkip,
				evolution.StrategySkip,
				evolution.StrategyMultXD,
				evolution.StrategyAddXD,
				evolution.StrategySubXD,
				evolution.StrategyDivXD,
			},
			AntagonistStrategyCount:  15,
			ProtagonistStrategyCount: 15,
			DepthOfRandomNewTrees:    1,
		},
		//FitnessCalculatorType: 0,
		//ShouldRunInteractiveTerminal: shouldRunInteractive,
	}

	_, err = simulation.Begin(params)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Print(finalParams)

	//err = simulation.CoalesceFiles(finalParams)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//cmd := exec.Command("Rscript", "launchCoalesced.R", coalescedFilesPath)
	//log.Fatal(cmd.Run())
}
