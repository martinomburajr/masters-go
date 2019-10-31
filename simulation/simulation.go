package simulation

import (
	"encoding/json"
	"fmt"
	"github.com/martinomburajr/masters-go/evolution"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Simulation struct {
	CurrentEvolutionState evolution.EvolutionParams
	NumberOfRunsPerState  int
	Name                  string
	// Output-Only
	OutputDir string
}

func (s *Simulation) Begin(params evolution.EvolutionParams) error {
	os.Mkdir("data", 0755)

	var folder string
	for i := 0; i < s.NumberOfRunsPerState; i++ {
		engine := PrepareSimulation(params, i)
		folder = engine.Parameters.StatisticsOutput.OutputDir
		err := StartEngine(engine)
		if err != nil {
			return err
		}
	}

	s.OutputDir = folder
	return nil
}

// CoalesceFiles will coalesce all files into a single document that can be analyzed
func (s *Simulation) CoalesceFiles() (string, error) {
	files := make([]string, 0)
	err := filepath.Walk(s.OutputDir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return "", err
	}

	newFiles := make([]string, 0)
	for i := range files {
		if strings.Contains(files[i], "json") {
			newFiles = append(newFiles, files[i])
		}
	}
	if newFiles == nil {
		return "",fmt.Errorf("CoalesceFiles | could not coalesce files - nil")
	}
	if len(newFiles) == 0 {
		return "",fmt.Errorf("CoalesceFiles | no files to coalesce")
	}

	jsonOutputs := make([]evolution.JSONOutput, len(newFiles))
	for i := 1; i < len(newFiles); i++ {
		filePath := fmt.Sprintf("%s", newFiles[i])
		file, err := os.Open(filePath)
		if err != nil {
			return "", err
		}

		split := strings.Split(filePath, "/")
		topLevelDir := split[0]
		subInfoDir := split[1]
		subSubNameDir := split[2]

		absolutePath, err := filepath.Abs(filePath)
		if err != nil {
			return "", err
		}
		workingDir := strings.ReplaceAll(absolutePath, filePath, "")
		statsPath := fmt.Sprintf("%s%s/%s/%s/%s", workingDir, topLevelDir, subInfoDir, subSubNameDir, "stats")
		RLaunchPath := fmt.Sprintf("%s%s", workingDir, "R/launch.R")
		go func(){
			cmd := exec.Command("Rscript",
				RLaunchPath,
				absolutePath,
				statsPath)
		err = cmd.Start()
		if err != nil {
				log.Println(err.Error())
			}
		}()



		var jsonOutput evolution.JSONOutput
		err = json.NewDecoder(file).Decode(&jsonOutput)
		if err != nil {
			return "",err
		}
		jsonOutputs[i-1] = jsonOutput
	}

	return coalesce(jsonOutputs, s.OutputDir)
}


func coalesce(files []evolution.JSONOutput, outputDir string) (string, error) {
	if files == nil {
		return "", fmt.Errorf("coalesce | json files cannot be nil")
	}
	if len(files) < 1 {
		return "", fmt.Errorf("coalesce | json files cannot be empty")
	}
	if outputDir == "" {
		return "", fmt.Errorf("outputDir empty")
	}

	coalesced := evolution.JSONCoalescedOutput{
		Name: outputDir[:len(outputDir)-1],
		CoalescedOutput: files,
	}

	path := fmt.Sprintf("%s%s", outputDir, "coalesced.json")
	err := os.Mkdir(outputDir, 0755)
	file, err := os.Create(path)
	if err != nil {
		return path, err
	}
	fmt.Printf("\nWrote to file: %s", path)
	err = json.NewEncoder(file).Encode(coalesced)
	if err != nil {
		return path, err
	}

	return path, nil
}


func StartEngine(engine *evolution.EvolutionEngine) error {
	evolutionResult, err := engine.Start()
	if err != nil {
		return err
	}

	err = evolutionResult.Analyze(engine.Generations, engine.Parameters.FitnessStrategy.IsMoreFitnessBetter, engine.Parameters)
	if err != nil {
		return err
	}

	antagonistSummary, err := evolutionResult.PrintTopIndividualSummary(evolution.IndividualAntagonist)
	if err != nil {
		return err
	}
	fmt.Println(antagonistSummary.String())

	protagonistSummary, err := evolutionResult.PrintTopIndividualSummary(evolution.IndividualProtagonist)
	if err != nil {
		return err
	}
	fmt.Println(protagonistSummary.String())

	averageGenerationSummary, err := evolutionResult.PrintAverageGenerationSummary()
	if err != nil {
		return err
	}
	fmt.Println(averageGenerationSummary.String())

	if engine.Parameters.ShouldRunInteractiveTerminal {
		err = evolutionResult.StartInteractiveTerminal(engine.Parameters)
		if err != nil {
			return err
		}
	}

	return nil
}

// PrepareSimulation takes in the given evolution parameters and a count variable and returns the engine that can be
// started run the simulation. The evolution engine will run count times.
func PrepareSimulation(params evolution.EvolutionParams, count int) *evolution.EvolutionEngine {
	rand.Seed(time.Now().UTC().UnixNano()) //Set seed

	if params.SpecParam.Seed < 0 {
		params.FitnessCalculatorType = 1
	}
	params.SpecParam.Expression = strings.ReplaceAll(params.SpecParam.Expression, " ", "")

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
	case evolution.FitnessMonoThresholdedRatio:
		spec, err = evolution.GenerateSpecSimple(params.SpecParam, params.FitnessStrategy, params.FitnessCalculatorType)
		if err != nil {
			log.Fatalf("MAIN | failed to create a valid spec | %s", err.Error())
		}
	case evolution.FitnessDualThresholdedRatio:
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
		engine.Parameters.FitnessStrategy.IsMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessAbsolute")
		break
	case evolution.FitnessRatio:
		engine.Parameters.FitnessStrategy.IsMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessRatio")
		break
	case evolution.FitnessProtagonistThresholdTally:
		fmt.Printf("Fitness Strategy: %s\n", "FitnessProtagonistThresholdTally")
		engine.Parameters.FitnessStrategy.IsMoreFitnessBetter = false
		break
	case evolution.FitnessThresholdedAntagonistRatio:
		engine.Parameters.FitnessStrategy.IsMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessThresholdedAntagonistRatio")
		break
	case evolution.FitnessMonoThresholdedRatio:
		engine.Parameters.FitnessStrategy.IsMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessMonoThresholdedRatio")
		break
	case evolution.FitnessDualThresholdedRatio:
		engine.Parameters.FitnessStrategy.IsMoreFitnessBetter = true
		fmt.Printf("Fitness Strategy: %s\n", "FitnessDualThresholdedRatio")
		break
	default:
		engine.Parameters.FitnessStrategy.IsMoreFitnessBetter = true
		log.Printf("Fitness Strategy: %s\n", "Unknown")
	}
	fmt.Printf("Fitness Strategy: %s\n", engine.Parameters.FitnessStrategy)
	fmt.Printf("Is More Fitness Better: %t\n", engine.Parameters.FitnessStrategy.IsMoreFitnessBetter)
	fmt.Println()

	var outputPath string
	var folder string
	timeStr := strings.ReplaceAll(strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", ""), "+", "")
	if engine.Parameters.StatisticsOutput.Name == "" {
		folder = fmt.Sprintf("data/%s/", engine.Parameters.ToString())
		outputPath = fmt.Sprintf("%s%s-%d.json", folder,
			timeStr, count)
	} else {
		folder = fmt.Sprintf("data/%s-%s/", engine.Parameters.StatisticsOutput.Name, engine.Parameters.ToString())
		outputPath = fmt.Sprintf("%s%s-%d.json", folder, timeStr, count)
	}
	engine.Parameters.StatisticsOutput.OutputDir = folder
	engine.Parameters.StatisticsOutput.OutputPath = outputPath
	return engine
	// ########################### START THE EVOLUTION PROCESS ##################################################3
}


// BeginToil will work through a multidimensional set of data to try all possible combination of parameters for ideal
// parameter tuning
func (s *Simulation) BeginToil(indexFile string) error {
	var loadedIndex IndexProgress
	file, err := os.Open(indexFile)
	if err != nil {
		return err
	}
	err = json.NewDecoder(file).Decode(&loadedIndex)
	if err != nil {
		return err
	}

	for expressionIndex := loadedIndex.ExpressionIndex; expressionIndex < len(AllExpressions); expressionIndex++ {
		for rangesIndex := loadedIndex.RangesIndex; rangesIndex < len(AllRanges); rangesIndex++ {
			for seedIndex := loadedIndex.SeedIndex; seedIndex < len(AllSeed); seedIndex++ {
				for generationsCountIndex := loadedIndex.GenerationsCountIndex; generationsCountIndex < len(
					AllGenerationsCount); generationsCountIndex++ {
					for eachPopulationIndex := loadedIndex.EachPopulationIndex; eachPopulationIndex < len(
						AllEachPopulationSize); eachPopulationIndex++ {
						for reproductionIndex := loadedIndex.ReproductionIndex; reproductionIndex < len(
							AllReproduction); reproductionIndex++ {
							for allDepthOfRandomNewTreeIndex := loadedIndex.
								AllDepthOfRandomNewTreeIndex; allDepthOfRandomNewTreeIndex < len(
								AllDepthOfRandomNewTree); allDepthOfRandomNewTreeIndex++ {
								for antagonistStrategyCountIndex := loadedIndex.
									AntagonistStrategyCountIndex; antagonistStrategyCountIndex < len(
									AllAntagonistStrategyCount); antagonistStrategyCountIndex++ {
									for protagonisttStrategyCountIndex := loadedIndex.
										ProtagonisttStrategyCountIndex; protagonisttStrategyCountIndex < len(
										AllProtagonistStrategyCount); protagonisttStrategyCountIndex++ {
										//for availableStrategyIndex = 0; availableStrategyIndex < len(AllAvailableStrategy); availableStrategyIndex++ {
										for fitnessStrategyTypeIndex := loadedIndex.
											FitnessStrategyTypeIndex; fitnessStrategyTypeIndex < len(
											AllFitnessStrategyType); fitnessStrategyTypeIndex++ {
											for fitStratAntThresMultIndex := loadedIndex.
												FitStratAntThresMultIndex; fitStratAntThresMultIndex < len(
												AllFitStratAntThreshMult); fitStratAntThresMultIndex++ {
												for fitStratProtThreshMultIndex := loadedIndex.
													FitStratProtThreshMultIndex; fitStratProtThreshMultIndex < len(
													AllFitStratProThreshMult); fitStratProtThreshMultIndex++ {
													for selectParentTypeIndex := loadedIndex.
														SelectParentTypeIndex; selectParentTypeIndex < len(
														AllSelectionSurvivorPercentage); selectParentTypeIndex++ {
														for strategiesAntagonistIndex := loadedIndex.
															StrategiesAntagonistIndex; strategiesAntagonistIndex < len(
															AllStrategies); strategiesAntagonistIndex++ {
															for selectSurvivorPercentIndex := loadedIndex.SelectSurvivorPercentIndex; selectSurvivorPercentIndex < len(AllSelectionSurvivorPercentage); selectSurvivorPercentIndex++ {
																for strategiesProtagonistIndex := loadedIndex.
																	StrategiesProtagonistIndex; strategiesProtagonistIndex < len(
																	AllStrategies); strategiesProtagonistIndex++ {
																	for numberOfRunsPerState := loadedIndex.
																		NumberOfRunsPerState; numberOfRunsPerState < s.
																		NumberOfRunsPerState; numberOfRunsPerState++ {

																		// TODO Add Parallelism
																		params := evolution.EvolutionParams{
																			GenerationsCount:   AllGenerationsCount[generationsCountIndex],
																			EachPopulationSize: AllEachPopulationSize[eachPopulationIndex],
																			SpecParam: evolution.SpecParam{
																				Seed:       AllSeed[seedIndex],
																				Expression: AllExpressions[expressionIndex],
																				Range:      AllRanges[rangesIndex],
																				AvailableVariablesAndOperators: evolution.AvailableVariablesAndOperators{
																					Constants: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
																					Variables: []string{"x"},
																					Operators: []string{"*", "+", "-"},
																				},
																			},
																			Reproduction: AllReproduction[reproductionIndex],
																			Strategies: evolution.Strategies{
																				AntagonistAvailableStrategies:  AllPossibleStrategies[strategiesAntagonistIndex],
																				ProtagonistAvailableStrategies: AllPossibleStrategies[strategiesProtagonistIndex],
																				AntagonistStrategyCount:        AllAntagonistStrategyCount[antagonistStrategyCountIndex],
																				ProtagonistStrategyCount:       AllProtagonistStrategyCount[protagonisttStrategyCountIndex],
																			},
																			Selection: evolution.Selection{
																				Survivor: evolution.SurvivorSelection{
																					Type:               1,
																					SurvivorPercentage: AllSelectionSurvivorPercentage[selectSurvivorPercentIndex],
																				},
																			},
																			FitnessStrategy: evolution.FitnessStrategy{
																				Type:                           AllFitnessStrategyType[fitnessStrategyTypeIndex],
																				AntagonistThresholdMultiplier:  AllFitStratAntThreshMult[fitStratAntThresMultIndex],
																				ProtagonistThresholdMultiplier: AllFitStratProThreshMult[fitStratProtThreshMultIndex],
																			},
																			StatisticsOutput: evolution.StatisticsOutput{
																				OutputPath: "",
																			},
																		}

																		engine := PrepareSimulation(params, numberOfRunsPerState)
																		err := StartEngine(engine)
																		if err != nil {
																			return err
																		}

																		// Write index progress to file
																		indexProgress := IndexProgress{
																			NumberOfRunsPerState:           numberOfRunsPerState,
																			AllDepthOfRandomNewTreeIndex:   allDepthOfRandomNewTreeIndex,
																			AntagonistStrategyCountIndex:   antagonistStrategyCountIndex,
																			EachPopulationIndex:            eachPopulationIndex,
																			ExpressionIndex:                expressionIndex,
																			FitnessStrategyTypeIndex:       fitnessStrategyTypeIndex,
																			FitStratAntThresMultIndex:      fitStratAntThresMultIndex,
																			FitStratProtThreshMultIndex:    fitStratProtThreshMultIndex,
																			GenerationsCountIndex:          generationsCountIndex,
																			ProtagonisttStrategyCountIndex: protagonisttStrategyCountIndex,
																			RangesIndex:                    rangesIndex,
																			ReproductionIndex:              reproductionIndex,
																			SeedIndex:                      seedIndex,
																			SelectParentTypeIndex:          selectParentTypeIndex,
																			SelectSurvivorPercentIndex:     selectSurvivorPercentIndex,
																			StrategiesAntagonistIndex:      strategiesAntagonistIndex,
																			StrategiesProtagonistIndex:     strategiesProtagonistIndex,
																		}
																		err = WriteIndexProgressToFile(indexProgress, indexFile)
																		if err != nil {
																			return err
																		}

																		if numberOfRunsPerState == loadedIndex.
																			NumberOfRunsPerState-1 {

																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}

					}
				}
			}
		}
	}

	return nil

}



func WriteIndexProgressToFile(indexProgress IndexProgress, indexFile string) error {
	file, err := os.Create(indexFile)
	if err != nil {
		return err
	}
	return json.NewEncoder(file).Encode(indexProgress)
}

type IndexProgress struct {
	ExpressionIndex                int `expressionIndex`
	RangesIndex                    int `json:"rangesIndex"`
	SeedIndex                      int `json:"seedIndex"`
	GenerationsCountIndex          int `json:"generationsCountIndex"`
	EachPopulationIndex            int `json:"eachPopulationIndex"`
	ReproductionIndex              int `json:"reproductionIndex"`
	AllDepthOfRandomNewTreeIndex   int `json:"allDepthOfRandomNewTreeIndex"`
	AntagonistStrategyCountIndex   int `json:"antagonistStrategyCountIndex"`
	ProtagonisttStrategyCountIndex int `json:"protagonisttStrategyCountIndex"`
	FitnessStrategyTypeIndex       int `json:"fitnessStrategyTypeIndex"`
	FitStratAntThresMultIndex      int `json:"fitStratAntThresMultIndex"`
	FitStratProtThreshMultIndex    int `json:"fitStratProtThreshMultIndex"`
	SelectParentTypeIndex          int `json:"selectParentTypeIndex"`
	SelectSurvivorPercentIndex     int `json:"selectSurvivorPercentIndex"`
	StrategiesAntagonistIndex      int `json:"strategiesAntagonistIndex"`
	StrategiesProtagonistIndex     int `json:"strategiesProtagonistIndex"`
	NumberOfRunsPerState           int `json:"numberOfRunsPerState"`
}

