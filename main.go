package main

import (
	"github.com/martinomburajr/masters-go/evolution"
	"github.com/martinomburajr/masters-go/simulation"
	"log"
)

func main() {
	simulation := simulation.Simulation{
		NumberOfRunsPerState: 5,
		Name:                 "simulation-1",
		OutputDir:            "",
	}

	params := evolution.EvolutionParams{ //f, err := os.Create("test.json")
		StatisticsOutput: evolution.StatisticsOutput{ //if err != nil {
			OutputPath: "", //	log.Fatal(err)
		}, //}
		SpecParam: evolution.SpecParam{ //json.NewEncoder(f).Encode(params)
			Range:      5,
			Expression: "5*x*x*x",
			Seed:       -1000,
			AvailableVariablesAndOperators: evolution.AvailableVariablesAndOperators{
				Constants: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
				Variables: []string{"x"},
				Operators: []string{"*", "+", "-"},
			},
		},
		GenerationsCount:   50,
		EachPopulationSize: 100, // Must be an even number to prevent awkward ordering of children.

		FitnessStrategy: evolution.FitnessStrategy{
			Type:                           evolution.FitnessDualThresholdedRatio,
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
				evolution.StrategyAddTreeWithMult,
				evolution.StrategyAddMultX,
				evolution.StrategyAddTreeWithSub,
				evolution.StrategyAddSubX,
				evolution.StrategyAddTreeWithAdd,
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
				evolution.StrategyAddTreeWithMult,
				evolution.StrategyAddMultX,
				evolution.StrategyAddTreeWithSub,
				evolution.StrategyAddSubX,
				evolution.StrategyAddTreeWithAdd,
				evolution.StrategyAddAddX,
				evolution.StrategySkip,
			},
			AntagonistStrategyCount:  5,
			ProtagonistStrategyCount: 5,
			DepthOfRandomNewTrees:    1,
		},
		//FitnessCalculatorType: 0,
		//ShouldRunInteractiveTerminal: shouldRunInteractive,
	}

	err := simulation.Begin(params)
	if err != nil {
		log.Fatal(err)
	}

	_, err = simulation.CoalesceFiles()
	if err != nil {
		log.Fatal(err)
	}

	//cmd := exec.Command("Rscript", "launchCoalesced.R", coalescedFilesPath)
	//log.Fatal(cmd.Run())
}
