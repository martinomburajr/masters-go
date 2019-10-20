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
		evolution.StrategyDeleteNonTerminal,
		evolution.StrategyDeleteMalicious,
		evolution.StrategyDeleteTerminal,
		evolution.StrategyMutateNonTerminal,
		evolution.StrategyMutateTerminal,
		evolution.StrategyReplaceBranch,
		evolution.StrategyAddSubTree,
		evolution.StrategyAddToLeaf,
		evolution.StrategyAddMult,
		evolution.StrategyAddSub,
		evolution.StrategyAddAdd,
		evolution.StrategyFellTree,
	}

	// TODO Include terminals and non terminals as part of strategy?
	params := evolution.EvolutionParams{
		Generations:                           50,
		EachPopulationSize:                    200, // Must be an even number to prevent awkward ordering of children.
		AntagonistMaxStrategies:               20,
		ProtagonistMaxStrategies:              20,
		DepthPenaltyStrategyPenalization:      10,
		ProbabilityOfMutation:                 0.1,
		ProbabilityOfNonTerminalMutation:      0.1,
		DepthOfRandomNewTrees:                 1,
		DeletionType:                          evolution.DeletionTypeSafe,
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
		AntagonistThresholdMultiplier:         16,
		ProtagonistThresholdMultiplier:        1.2,
	}

	expression := "x*x*x"
	expression = eval.MartinsReplace(expression, " ", "")
	specCount := 20

	constants := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	variables := []string{"x"}
	operators := []string{"*", "+", "-"}

	constantTerminals, err := evolution.GenerateTerminals(10, constants)
	if err != nil {
		log.Fatal(err)
	}
	variableTerminals, err := evolution.GenerateTerminals(10, variables)
	if err != nil {
		log.Fatal(err)
	}
	nonTerminals, err := evolution.GenerateNonTerminals(3, operators)
	if err != nil {
		log.Fatal(err)
	}
	terminals := append(variableTerminals, constantTerminals...)

	_, _, mathematicalExpression, err := evolution.ParseString(expression, operators, variables)
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
		spec, err = evolution.GenerateSpecSimple(expression, specCount, -1*(specCount/2),
			params.ProtagonistThresholdMultiplier, params.ProtagonistThresholdMultiplier)
		if err != nil {
			log.Fatalf("MAIN | failed to create a valid spec | %s", err.Error())
		}
	case evolution.FitnessDualThresholdedRatioFitness:
		spec, err = evolution.GenerateSpecSimple(expression, specCount, -1*(specCount/2),
			params.AntagonistThresholdMultiplier, params.ProtagonistThresholdMultiplier)
		if err != nil {
			log.Fatalf("MAIN | failed to create a valid spec | %s", err.Error())
		}
	default:
		spec, err = evolution.GenerateSpecSimple(expression, specCount, -1*(specCount/2),
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
		Generations: make([]*evolution.Generation, params.Generations),
	}

	// ########################### START THE EVOLUTION PROCESS ##################################################3
	evolutionResult, err := engine.Start()
	if err != nil {
		log.Fatal(err)
	}

	// ########################### OUTPUT STATISTICS  #######################################################3
	fmt.Printf("Generation Count: %d\n", engine.Parameters.Generations)
	fmt.Printf("Each Individual Count: %d\n", engine.Parameters.EachPopulationSize)
	isMoreFitnessBetter := true

	switch engine.Parameters.FitnessStrategy {
	case evolution.FitnessAbsolute:
		isMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessAbsolute")
		break
	case evolution.FitnessRatio:
		isMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessRatio")
		break
	case evolution.FitnessRatioThresholder:
		isMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessRatioThresholder")
		break
	case evolution.FitnessProtagonistThresholdTally:
		fmt.Printf("Fitness Strategy: %s\n", "FitnessProtagonistThresholdTally")
		isMoreFitnessBetter = false
		break
	case evolution.FitnessImproverTally:
		isMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessImproverTally")
		break
	case evolution.FitnessMonoThresholdedRatioFitness:
		isMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessMonoThresholdedRatioFitness")
		break
	case evolution.FitnessDualThresholdedRatioFitness:
		isMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessDualThresholdedRatioFitness")
		break
	default:
		isMoreFitnessBetter = true
		log.Printf("Fitness Strategy: %s\n", "Unknown")
	}
	fmt.Printf("Fitness Straegy: %d\n", engine.Parameters.FitnessStrategy)
	fmt.Println()

	err = evolutionResult.Analyze(engine.Generations, isMoreFitnessBetter, 3)
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

	fmt.Println()
}

func GenerateMathExpression() {}
