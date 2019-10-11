package main

import (
	"fmt"
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
		evolution.StrategyMutateNode,
	}

	//starterTreeExpression := "( x ) * ( 5 )"
	//expressionSet := []evolution.SymbolicExpression{evolution.X1, evolution.Mult, evolution.Const5}
	//expressionSet := []evolution.SymbolicExpression{evolution.X1}
	constants := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	variables := []string{"x", "y", "a", "b", "c", "d"}
	operators := []string{"*", "+", "-"}

	constantTerminals, err := evolution.GenerateTerminals(3, constants)
	if err != nil {
		log.Fatal(err)
	}
	variableTerminals, err := evolution.GenerateTerminals(3, variables)
	if err != nil {
		log.Fatal(err)
	}
	nonTerminals, err := evolution.GenerateNonTerminals(3, operators)
	if err != nil {
		log.Fatal(err)
	}
	terminals := append(variableTerminals, constantTerminals...)

	_, _, mathematicalExpression, err := evolution.ParseString("x * 5", operators)
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
	spec, err := evolution.GenerateSpec(starterTreeAsMathematicalExpression, 10, 0)
	if err != nil {
		log.Fatalf("MAIN | failed to create a valid spec | %s", err.Error())
	}

	fmt.Printf("Protagonist vs Antagonist Competitive Coevolution:\nMathematical Expression: %s\nSpec: %s\n",
		starterTreeAsMathematicalExpression,
		spec.ToString())

	// TODO only perform parent selection on loser
	// TODO Do children undergo tournament selection
	// TODO Include terminals and non terminals as part of strategy?
	// TODO Should threshold increase given spec
	// TODO Should we pick most recent individual even if fitness is the same?
	params := evolution.EvolutionParams{
		Generations:                           50,
		EachPopulationSize:                    4, // Must be an even number to prevent awkward ordering of children.
		AntagonistMaxStrategies:               4,
		ProtagonistMaxStrategies:              4,
		DepthPenaltyStrategyPenalization:      10,
		ProbabilityOfMutation:                 0.1,
		ProbabilityOfNonTerminalMutation:      0.1,
		DepthOfRandomNewTrees:                 1,
		DeletionType:                          evolution.DeletionTypeSafe,
		EnforceIndependentVariable:            true,
		ProtagonistAvailableStrategies:        strategies,
		AntagonistAvailableStrategies:         strategies,
		SetEqualStrategyLength:                true,
		CrossoverPercentage:                   0.3,
		MaintainCrossoverGeneTransferEquality: true,
		NonTerminalSet:                        nonTerminals,
		TerminalSet:                           terminals,
		FitnessStrategy:                       evolution.FitnessRatio,
		EvaluationThreshold:                   12,
		TournamentSize:                        3,
		StrategyLengthLimit:                   10,
		SurvivorPercentage:                    0.5,
		StartIndividual:                       startProgram,
		Spec:                                  spec,
		ParentSelection:                       evolution.ParentSelectionTournament,
		EqualStrategiesLength:                 2,
		ThresholdMultiplier:                   1.5,
	}

	engine := evolution.EvolutionEngine{
		Parameters:  params,
		Generations: make([]*evolution.Generation, params.Generations),
	}

	// ########################### START THE EVOLUTION PROCESS ##################################################3
	evolutionResult, err := engine.Start()
	if err != nil {
		log.Fatal(err)
	}

	// ########################### OUPUT STATISTICS  #######################################################3
	fmt.Printf("Generation Count: %d\n", engine.Parameters.Generations)
	fmt.Printf("Each Individual Count: %d\n", engine.Parameters.EachPopulationSize)
	fmt.Println()

	_, _ = evolutionResult.Analyze(engine.Generations, 3)
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

	cumGenerationSummary, err := evolutionResult.PrintCumGenerationSummary()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cumGenerationSummary.String())

	fmt.Println()
	//fmt.Print(result)
}

func GenerateMathExpression() {}
