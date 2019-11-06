package simulation

import (
	"encoding/json"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/martinomburajr/masters-go/evolution"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Simulation struct {
	CurrentEvolutionState evolution.EvolutionParams
	NumberOfRunsPerState  int
	Name                  string
	// Output-Only
	OutputDir string
	SimulationStats []SimulationRunStats
}

func (s *Simulation) Begin(params evolution.EvolutionParams) (evolution.EvolutionParams, error) {
	os.Mkdir("data", 0755)
	s.SimulationStats = make([]SimulationRunStats, s.NumberOfRunsPerState)

	//wg := sync.WaitGroup{}
	var folder string
	for i := 0; i < s.NumberOfRunsPerState; i++ {
		//wg.Add(1)
		//go func(params evolution.EvolutionParams, i int){
		//	defer wg.Done()
		params.InternalCount = i
		engine := PrepareSimulation(params, i)
		params = engine.Parameters
		folder = engine.Parameters.StatisticsOutput.OutputDir

		s.StartEngine(engine)
		//}(params, i)
	}
	//wg.Wait()
	log.Println("SYNCHRONIZED!")
	s.OutputDir = folder
	outputPath := fmt.Sprintf("%s/%s", s.OutputDir, "runs.csv")

	statistics, err := s.ToRunStats()
	if err != nil {
		return params, err
	}

	// WRITE TO FILE
	outputFileCSV, err := os.Create(outputPath)
	if err != nil {
		return params,err
	}
	defer outputFileCSV.Close()

	writer := gocsv.DefaultCSVWriter(outputFileCSV)
	if writer.Error() != nil {
		return params,writer.Error()
	}
	err = gocsv.Marshal(statistics, outputFileCSV)
	if err != nil {
		return params,err
	}

	// WRITe STRATEGIES
	strategyStats, err := s.ToStrategyStats(s.OutputDir)
	if err != nil {
		return params,err
	}
	err = WriteRunStrategy(strategyStats, fmt.Sprintf("%s/%s", s.OutputDir, "strategy.csv"))
	if err != nil {
		return params,err
	}

	return params,err
}



func WriteRunStrategy(runStrat []RunStrategyStatistics, outputPath string) error {
	// WRITE TO FILE
	outputFileCSV, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFileCSV.Close()

	writer := gocsv.DefaultCSVWriter(outputFileCSV)
	if writer.Error() != nil {
		return writer.Error()
	}
	err = gocsv.Marshal(runStrat, outputFileCSV)
	if err != nil {
		return err
	}
	return nil
}
type SimulationRunStats struct {
	TopAntagonist evolution.Individual
	TopProtagonist evolution.Individual
	FinalAntagonist evolution.Individual
	FinalProtagonist evolution.Individual
}

// CoalesceFiles will coalesce all files into a single document that can be analyzed
func (s *Simulation) CoalesceFiles(evolutionParams evolution.EvolutionParams) error {

	multiCSVOutput := evolution.MultiOutput{}
	err := multiCSVOutput.Load(s.OutputDir, "", evolutionParams)

	// write params to file
	paramsPath := fmt.Sprintf("%s/%s", s.OutputDir, "params.json")
	paramsFile, err := os.Create(paramsPath)
	if err != nil {
		return err
	}
	err = json.NewEncoder(paramsFile).Encode(evolutionParams)
	if err != nil {
		return err
	}
	fmt.Println("Wrote Params file to: " + paramsPath)

	//multiCSVOutput.Load(s.OutputDir)
	multiCSVOutput.WriteGenerationalCumulative(evolutionParams, fmt.Sprintf("%s/%s", s.OutputDir,
		"cumulative.csv"))
	if err != nil {
		return err
	}

	multiCSVOutput.WriteAverages(evolutionParams, fmt.Sprintf("%s/%s", s.OutputDir,
		"averages.csv"))
	if err != nil {
		return err
	}

	multiCSVOutput.WriteBestIndividuals(evolutionParams, s.OutputDir, fmt.Sprintf("%s/%s", s.OutputDir,
		"best.csv"))
	if err != nil {
		return err
	}

	return err
}

func (s *Simulation) RunRScript(absolutePath string, filePath string, topLevelDir string, subInfoDir string,
	subSubNameDir string, run int) chan error {

	workingDir := strings.ReplaceAll(absolutePath, filePath, "")
	epochalPath := fmt.Sprintf("%s%s/%s/%s/%s-%d.csv", workingDir, topLevelDir, subInfoDir, subSubNameDir, "epochal",
		run)
	statsPath := fmt.Sprintf("%s%s/%s/%s/%s", workingDir, topLevelDir, subInfoDir, subSubNameDir, "stats")
	RLaunchPath := fmt.Sprintf("%s%s", workingDir, "R/runScript.R")

	errChan := make(chan error)

	//go func() {
	cmd := exec.Command("Rscript",
		RLaunchPath,
		absolutePath,
		epochalPath,
		statsPath)
	log.Println(fmt.Sprintf("Rscript: %s", cmd.String()))
	err := cmd.Start()
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		errChan <- err
	}
	//}()

	return errChan
}

func (s *Simulation) StartEngine(engine *evolution.EvolutionEngine) error {
	evolutionResult, err := engine.Start()
	if err != nil {
		return err
	}

	err = evolutionResult.Analyze(engine.Generations, true, engine.Parameters)
	if err != nil {
		return err
	}

	topAnt, err := evolutionResult.TopAntagonist.Clone()
	topProt, err := evolutionResult.TopProtagonist.Clone()
	finAnt, err := evolutionResult.FinalAntagonist.Clone()
	finProt, err := evolutionResult.FinalProtagonist.Clone()


	s.SimulationStats[engine.Parameters.InternalCount] = SimulationRunStats{
		TopAntagonist:    topAnt,
		TopProtagonist:   topProt,
		FinalAntagonist:  finAnt,
		FinalProtagonist: finProt,
	}

	return nil
}

func (s *Simulation) ToRunStats() ([]RunBasedStatistics, error) {
	if s.SimulationStats == nil {
		return nil, fmt.Errorf("ToRunStats | simulationStats is nil")
	}
	if len(s.SimulationStats) < 0 {
		return nil, fmt.Errorf("ToRunStats | simulationStats is empty")
	}
	runStats := make([]RunBasedStatistics, s.NumberOfRunsPerState)
	for i, simulationStat := range s.SimulationStats {
		topAntEq, _ := simulationStat.TopAntagonist.Program.T.ToMathematicalString()
		topProEq, _ := simulationStat.TopProtagonist.Program.T.ToMathematicalString()
		finAntEq, _ := simulationStat.FinalAntagonist.Program.T.ToMathematicalString()
		finProEq, _ := simulationStat.FinalProtagonist.Program.T.ToMathematicalString()
		runStats[i] = RunBasedStatistics{
			TopAntagonist:            simulationStat.TopAntagonist.TotalFitness,
			TopProtagonist:           simulationStat.TopProtagonist.TotalFitness,
			TopAntagonistDelta:       simulationStat.TopAntagonist.FitnessDelta,
			TopProtagonistDelta:      simulationStat.TopProtagonist.FitnessDelta,
			TopAntagonistStrategy:    evolution.StrategiesToString(simulationStat.TopAntagonist),
			TopProtagonistStrategy:   evolution.StrategiesToString(simulationStat.TopProtagonist),
			TopAntagonistEquation:    topAntEq,
			TopProtagonistEquation:   topProEq,
			FinalAntagonist:          simulationStat.FinalAntagonist.TotalFitness,
			FinalProtagonist:         simulationStat.FinalProtagonist.TotalFitness,
			FinalAntagonistDelta:     simulationStat.FinalAntagonist.FitnessDelta,
			FinalProtagonistDelta:    simulationStat.FinalProtagonist.FitnessDelta,
			FinalAntagonistStrategy:  evolution.StrategiesToString(simulationStat.FinalAntagonist),
			FinalProtagonistStrategy: evolution.StrategiesToString(simulationStat.FinalAntagonist),
			FinalAntagonistEquation:  finAntEq,
			FinalProtagonistEquation: finProEq,
			Run:                      i,
		}
	}

	return runStats, nil
}

// ToStrategyStats
func (s *Simulation) ToStrategyStats(dirPath string) (statistics []RunStrategyStatistics, err error) {
	if s.SimulationStats == nil {
		return  nil, fmt.Errorf("ToRunStats | simulationStats is nil")
	}
	if len(s.SimulationStats) < 0 {
		return  nil, fmt.Errorf("ToRunStats | simulationStats is empty")
	}
	statistics = make([]RunStrategyStatistics, len(s.SimulationStats[0].TopAntagonist.Strategy))

	for i, run := range s.SimulationStats {
		for j := range s.SimulationStats[0].TopAntagonist.Strategy {
			statistics[j] = RunStrategyStatistics {
				TopAntagonistStrategy:    string(run.TopAntagonist.Strategy[j]),
				TopProtagonistStrategy:   string(run.TopProtagonist.Strategy[j]),
				FinalAntagonistStrategy:  string(run.FinalAntagonist.Strategy[j]),
				FinalProtagonistStrategy: string(run.FinalProtagonist.Strategy[j]),
				StrategyNumber:           j+1,
			}
		}
		WriteRunStrategy(statistics, fmt.Sprintf("%s/%d/strategy-%d.csv", dirPath, i, i ))
	}

	return statistics, nil
}

type RunStrategyStatistics struct {
	StrategyNumber                       int     `csv:"stratNum"`
	TopAntagonistStrategy  string  `csv:"runTopAStrategy"`
	TopProtagonistStrategy string  `csv:"runTopPStrategy"`
	FinalAntagonistStrategy  string  `csv:"runFinAStrategy"`
	FinalProtagonistStrategy string  `csv:"runFinPStrategy"`
}

type RunBasedStatistics struct {
	TopAntagonist          float64 `csv:"runTopA"`
	TopProtagonist         float64 `csv:"runTopP"`
	TopAntagonistDelta     float64 `csv:"runTopADelta"`
	TopProtagonistDelta    float64 `csv:"runTopPDelta"`
	TopAntagonistStrategy  string  `csv:"runTopAStrategy"`
	TopProtagonistStrategy string  `csv:"runTopPStrategy"`
	TopAntagonistEquation     string  `csv:"runTopAEquation"`
	TopProtagonistEquation    string  `csv:"runTopPEquation"`

	FinalAntagonist          float64 `csv:"runFinalA"`
	FinalProtagonist         float64 `csv:"runFinalP"`
	FinalAntagonistDelta     float64 `csv:"runFinalADelta"`
	FinalProtagonistDelta    float64 `csv:"runFinalPDelta"`
	FinalAntagonistStrategy  string  `csv:"runFinalAStrategy"`
	FinalProtagonistStrategy string  `csv:"runFinalPStrategy"`
	FinalAntagonistEquation     string  `csv:"runFinalAEquation"`
	FinalProtagonistEquation    string  `csv:"runFinalPEquation"`

	Run                       int     `csv:"run"`
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

	fmt.Printf("ProtagonistEquation vs AntagonistEquation Competitive Coevolution:\nMathematical Expression: %s\nSpec:%s\n",
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
	fmt.Printf("Iteration Count: %d\n", count)

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
		fmt.Printf("Fitness Strategy: %s\n", "FitnessAbsolute")
		break
	case evolution.FitnessRatio:
		fmt.Printf("Fitness Strategy: %s\n", "FitnessRatio")
		break
	case evolution.FitnessProtagonistThresholdTally:
		fmt.Printf("Fitness Strategy: %s\n", "FitnessProtagonistThresholdTally")
		break
	case evolution.FitnessThresholdedAntagonistRatio:
		fmt.Printf("Fitness Strategy: %s\n", "FitnessThresholdedAntagonistRatio")
		break
	case evolution.FitnessMonoThresholdedRatio:
		fmt.Printf("Fitness Strategy: %s\n", "FitnessMonoThresholdedRatio")
		break
	case evolution.FitnessDualThresholdedRatio:
		fmt.Printf("Fitness Strategy: %s\n", "FitnessDualThresholdedRatio")
		break
	default:
		log.Printf("Fitness Strategy: %s\n", "Unknown")
	}
	fmt.Println()

	var outputPath string
	var folder string
	timeStr := engine.Parameters.InternalCount

	folder = fmt.Sprintf("data/%s/", engine.Parameters.ToString())
	outputPath = fmt.Sprintf("%s%d", folder,
		timeStr)

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
																					Type:               "SteadyState",
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
																		err := s.StartEngine(engine)
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
