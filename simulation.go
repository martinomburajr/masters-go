package main

import (
	"fmt"
	"github.com/martinomburajr/masters-go/eval"
	"github.com/martinomburajr/masters-go/evolution"
	"log"
	"math/rand"
	"time"
)

type Simulation struct {
	EvolutionStates      []evolution.EvolutionParams
	NumberOfRunsPerState int
	Name                 string
}

func (s *Simulation) Begin() {
	for i := 0; i < len(s.EvolutionStates); i++ {
		for j := 0; j < s.NumberOfRunsPerState; j++ {

			outputPath := fmt.Sprintf("%s/%s-%d.json", s.EvolutionStates[i].StatisticsOutput.Name,
				time.Now().Format(time.RFC3339), i)
			s.EvolutionStates[i].StatisticsOutput.OutputPath = outputPath
			engine, evolutionParams := PrepareSimulation(s.EvolutionStates[i])
			StartEngine(engine, evolutionParams)
		}
	}
}

func StartEngine(engine *evolution.EvolutionEngine, params evolution.EvolutionParams) {
	evolutionResult, err := engine.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = evolutionResult.Analyze(engine.Generations, params.FitnessStrategy.IsMoreFitnessBetter, engine.Parameters)
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
		err = evolutionResult.StartInteractiveTerminal(params)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func PrepareSimulation(params evolution.EvolutionParams) (*evolution.EvolutionEngine, evolution.EvolutionParams) {
	rand.Seed(time.Now().UTC().UnixNano()) //Set seed

	if params.SpecParam.Seed < 0 {
		params.FitnessCalculatorType = 1
	}
	params.SpecParam.Expression = eval.MartinsReplace(params.SpecParam.Expression, " ", "")

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

	switch params.FitnessStrategy.Type {
	case evolution.FitnessMonoThresholdedRatioFitness:
		spec, err = evolution.GenerateSpecSimple(params.SpecParam, params.FitnessStrategy, params.FitnessCalculatorType)
		if err != nil {
			log.Fatalf("MAIN | failed to create a valid spec | %s", err.Error())
		}
	case evolution.FitnessDualThresholdedRatioFitness:
		spec, err = evolution.GenerateSpecSimple(params.SpecParam, params.FitnessStrategy, params.FitnessCalculatorType)
		if err != nil {
			log.Fatalf("MAIN | failed to create a valid spec | %s", err.Error())
		}
	default:
		spec, err = evolution.GenerateSpecSimple(params.SpecParam, params.FitnessStrategy, params.FitnessCalculatorType)
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
	params.StartIndividual = startProgram
	params.SpecParam.AvailableSymbolicExpressions.Terminals = append(variableTerminals, constantTerminals...)
	params.SpecParam.AvailableSymbolicExpressions.NonTerminals = nonTerminals

	engine := &evolution.EvolutionEngine{
		Parameters:  params,
		Generations: make([]*evolution.Generation, params.GenerationsCount),
	}

	// ########################### OUTPUT STATISTICS  #######################################################3
	fmt.Printf("Generation Count: %d\n", engine.Parameters.GenerationsCount)
	fmt.Printf("Each Individual Count: %d\n", engine.Parameters.EachPopulationSize)

	switch engine.Parameters.FitnessCalculatorType {
	case 0:
		fmt.Printf("Fitness Calculation Method: %s\n", "Speedy")
	case 1:
		fmt.Printf("Fitness Calculation Method: %s\n", "Heavy")
	default:
		fmt.Printf("Fitness Calculation Method: %s\n", "Unknown")
	}

	switch engine.Parameters.FitnessStrategy.Type {
	case evolution.FitnessAbsolute:
		params.FitnessStrategy.IsMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessAbsolute")
		break
	case evolution.FitnessRatio:
		params.FitnessStrategy.IsMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessRatio")
		break
	case evolution.FitnessProtagonistThresholdTally:
		fmt.Printf("Fitness Strategy: %s\n", "FitnessProtagonistThresholdTally")
		params.FitnessStrategy.IsMoreFitnessBetter = false
		break
	case evolution.FitnessThresholdedAntagonistRatio:
		params.FitnessStrategy.IsMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessThresholdedAntagonistRatio")
		break
	case evolution.FitnessMonoThresholdedRatioFitness:
		params.FitnessStrategy.IsMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessMonoThresholdedRatioFitness")
		break
	case evolution.FitnessDualThresholdedRatioFitness:
		params.FitnessStrategy.IsMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessDualThresholdedRatioFitness")
		break
	default:
		params.FitnessStrategy.IsMoreFitnessBetter = true
		log.Printf("Fitness Strategy: %s\n", "Unknown")
	}
	fmt.Printf("Fitness Strategy: %d\n", engine.Parameters.FitnessStrategy)
	fmt.Printf("Is More Fitness Better: %t\n", params.FitnessStrategy.IsMoreFitnessBetter)
	fmt.Println()

	return engine, params
	// ########################### START THE EVOLUTION PROCESS ##################################################3
}
