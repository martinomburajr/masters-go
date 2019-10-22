package main

import (
	"fmt"
	"github.com/martinomburajr/masters-go/eval"
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
		evolution.StrategyFellTree,
		evolution.StrategyDeleteNonTerminal,
		evolution.StrategyDeleteMalicious,
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

	// TODO Include terminals and non terminals as part of strategy?
	params := evolution.EvolutionParams{
		GenerationsCount:                      50,
		EachPopulationSize:                    2, // Must be an even number to prevent awkward ordering of children.
		AntagonistMaxStrategies:               20,
		ProtagonistMaxStrategies:              20,
		ProbabilityOfMutation:                 0.1,
		DepthOfRandomNewTrees:                 1,
		EnforceIndependentVariable:            true,
		ProtagonistAvailableStrategies:        strategies,
		AntagonistAvailableStrategies:         strategies,
		SetEqualStrategyLength:                true,
		CrossoverPercentage:                   0.2,
		MaintainCrossoverGeneTransferEquality: true,
		FitnessStrategy:                       evolution.FitnessDualThresholdedRatioFitness,
		EvaluationThreshold:                   12,
		TournamentSize:                        3,
		StrategyLengthLimit:                   10,
		SurvivorPercentage:                    0.2,
		ParentSelection:                       evolution.FitnessDualThresholdedRatioFitness,
		EqualStrategiesLength:                 20,
		ThresholdMultiplier:                   1.5,
		AntagonistThresholdMultiplier:         30,
		ProtagonistThresholdMultiplier:        1.3,


		SpecParam: evolution.SpecParam{
			Range: 100,
			Expression: "x*x",
			Seed: -500,
			AvailableVariablesAndOperators: evolution.AvailableVariablesAndOperators{
				Constants: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
				Variables: []string{"x"},
				Operators: []string{"*", "+", "-"},
			},
		},

		ShouldRunInteractiveTerminal: false,
	}

	params.SpecParam.Expression = eval.MartinsReplace(params.SpecParam.Expression, " ", "")
	specCount := params.SpecParam.Range

	constantTerminals, err := evolution.GenerateTerminals(10, params.SpecParam.AvailableVariablesAndOperators.Constants)
	if err != nil {
		log.Fatal(err)
	}
	variableTerminals, err := evolution.GenerateTerminals(10, params.SpecParam.AvailableVariablesAndOperators.Variables)
	if err != nil {
		log.Fatal(err)
	}
	nonTerminals, err := evolution.GenerateNonTerminals(3, params.SpecParam.AvailableVariablesAndOperators.Operators)
	if err != nil {
		log.Fatal(err)
	}
	terminals := append(variableTerminals, constantTerminals...)

	_, _, mathematicalExpression, err := evolution.ParseString(params.SpecParam.Expression, params.SpecParam.AvailableVariablesAndOperators.Operators, params.SpecParam.AvailableVariablesAndOperators.Variables)
	if err != nil {
		log.Fatal(err)
	}

	starterTree := evolution.DualTree{}
	err = starterTree.FromSymbolicExpressionSet2(mathematicalExpression)
	if err != nil {
		log.Fatal("main | cannot parse symbolic expression tree to convert starter tree to a mathematical expression")
	}
	starterTreeAsMathematicalExpression, err := starterTree.ToMathematicalString()
	if err != nil {
		log.Fatal("main | failed to convert starter tree to a mathematical expression")
	}

	startProgram := evolution.Program{
		T: &starterTree,
	}
	var spec evolution.SpecMulti

	switch params.FitnessStrategy {
	case evolution.FitnessMonoThresholdedRatioFitness:
		spec, err = evolution.GenerateSpecSimple(params.SpecParam.Expression, specCount, params.SpecParam.Seed,
			params.ProtagonistThresholdMultiplier, params.ProtagonistThresholdMultiplier)
		if err != nil {
			log.Fatalf("MAIN | failed to create a valid spec | %s", err.Error())
		}
	case evolution.FitnessDualThresholdedRatioFitness:
		spec, err = evolution.GenerateSpecSimple(params.SpecParam.Expression, specCount, params.SpecParam.Seed,
			params.AntagonistThresholdMultiplier, params.ProtagonistThresholdMultiplier)
		if err != nil {
			log.Fatalf("MAIN | failed to create a valid spec | %s", err.Error())
		}
	default:
		spec, err = evolution.GenerateSpecSimple(params.SpecParam.Expression, specCount, params.SpecParam.Seed,
			params.ProtagonistThresholdMultiplier, params.ProtagonistThresholdMultiplier)
		if err != nil {
			log.Fatalf("MAIN | failed to create a valid spec | %s", err.Error())
		}
	}

	fmt.Printf("Protagonist vs Antagonist Competitive Coevolution:\nMathematical Expression: %s\nSpec:%s\n",
		starterTreeAsMathematicalExpression,
		spec.ToString(),
	)

	// Set extra params
	params.Spec = spec
	params.TerminalSet = terminals
	params.NonTerminalSet = nonTerminals
	params.StartIndividual = startProgram
	params.VariableTerminals = variableTerminals

	engine := evolution.EvolutionEngine{
		Parameters:  params,
		Generations: make([]*evolution.Generation, params.GenerationsCount),
	}

	// ########################### OUTPUT STATISTICS  #######################################################3
	fmt.Printf("Generation Count: %d\n", engine.Parameters.GenerationsCount)
	fmt.Printf("Each Individual Count: %d\n", engine.Parameters.EachPopulationSize)

	switch engine.Parameters.FitnessStrategy {
	case evolution.FitnessAbsolute:
		engine.IsMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessAbsolute")
		break
	case evolution.FitnessRatio:
		engine.IsMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessRatio")
		break
	case evolution.FitnessRatioThresholder:
		engine.IsMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessRatioThresholder")
		break
	case evolution.FitnessProtagonistThresholdTally:
		fmt.Printf("Fitness Strategy: %s\n", "FitnessProtagonistThresholdTally")
		engine.IsMoreFitnessBetter = false
		break
	case evolution.FitnessImproverTally:
		engine.IsMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessImproverTally")
		break
	case evolution.FitnessMonoThresholdedRatioFitness:
		engine.IsMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessMonoThresholdedRatioFitness")
		break
	case evolution.FitnessDualThresholdedRatioFitness:
		engine.IsMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessDualThresholdedRatioFitness")
		break
	default:
		engine.IsMoreFitnessBetter = true
		log.Printf("Fitness Strategy: %s\n", "Unknown")
	}
	fmt.Printf("Fitness Strategy: %d\n", engine.Parameters.FitnessStrategy)
	fmt.Printf("Is More Fitness Better: %t\n", engine.IsMoreFitnessBetter)
	fmt.Println()

	// ########################### START THE EVOLUTION PROCESS ##################################################3
	evolutionResult, err := engine.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = evolutionResult.Analyze(engine.Generations, engine.IsMoreFitnessBetter)
	if err != nil {
		log.Fatal(err)
	}
	antagonistSummary, err := evolutionResult.PrintTopIndividualSummary(evolution.IndividualAntagonist)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(antagonistSummary.String())

	protagonistSummary, err := evolutionResult.PrintTopIndividualSummary(evolution.IndividualProtagonist)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(protagonistSummary.String())

	averageGenerationSummary, err := evolutionResult.PrintAverageGenerationSummary()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(averageGenerationSummary.String())

	if params.ShouldRunInteractiveTerminal {
		err = evolutionResult.StartInteractiveTerminal()
		if err != nil {
			log.Fatal(err)
		}
	}
}
