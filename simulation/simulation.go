package simulation

import (
	"encoding/json"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/gosuri/uiprogress"
	"github.com/martinomburajr/masters-go/evolog"
	"github.com/martinomburajr/masters-go/evolution"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	ProgressCounterSimulation = 10
	ProgressRunStats          = 5
)

type Simulation struct {
	CurrentEvolutionState evolution.EvolutionParams `json:"currentEvolutionState"`
	NumberOfRunsPerState  int                       `json:"numberOfRunsPerState"`
	Name                  string                    `json:"name"`
	StatsFiles            []string                  `json:"statsFiles"`

	ProgressBar *uiprogress.Bar

	// Output-Only
	OutputDir       string               `json:"outputDir"`
	RPath           string               `json:"rPath"`
	SimulationStats []SimulationRunStats `json:"simulationStats"`
	DataPath        string
}

func (s *Simulation) Begin(params evolution.EvolutionParams) (evolution.EvolutionParams, error) {
	os.Mkdir("data", 0755)
	s.SimulationStats = make([]SimulationRunStats, s.NumberOfRunsPerState)
	newParamsChan := make(chan evolution.EvolutionParams, s.NumberOfRunsPerState)

	mutex := sync.Mutex{}
	uiprogress.Start()

	if params.EnableParallelism {
		wg := sync.WaitGroup{}
		for i := 0; i < s.NumberOfRunsPerState; i++ {
			time.Sleep(100 * time.Millisecond)
			wg.Add(1)
			progCount := ProgressCounterSimulation + evolution.ProgressCountersEvolutionResult

			if params.RunStats {
				progCount += ProgressRunStats
			}
			if params.MaxGenerations > evolution.MinAllowableGenerationsToTerminate {
				progCount += params.MaxGenerations

			} else {
				progCount += params.GenerationsCount
			}

			go func(i int, params evolution.EvolutionParams, newParamsChan chan evolution.EvolutionParams,
				s *Simulation, mutex *sync.Mutex, wg *sync.WaitGroup, progCount int) {
				defer wg.Done()

				var bar *uiprogress.Bar
				bar = uiprogress.AddBar(progCount).AppendCompleted().PrependElapsed().PrependFunc(func(b *uiprogress.
					Bar) string {

					folderNumber := strings.Split(params.ParamFile, "/")[0]
					sb := fmt.Sprintf("Folder: %s \t| Run: %d/%d \t| ==>", folderNumber, i+1, s.NumberOfRunsPerState)
					return sb
				})
				s.ProgressBar = bar

				mutex.Lock()
				params.InternalCount = i
				engine := PrepareSimulation(params, i)
				engine.Parameters.StatisticsOutput = params.StatisticsOutput

				params = engine.Parameters
				newParamsChan <- params
				s.OutputDir = engine.Parameters.StatisticsOutput.OutputDir
				mutex.Unlock()

				engine.ProgressBar = s.ProgressBar
				err := s.StartEngine(engine)
				if err != nil {
					params.ErrorChan <- err
				}
				engine.Generations = nil // FREE UP MEMORY
			}(i, params, newParamsChan, s, &mutex, &wg, progCount)
		}
		wg.Wait()
	} else {
		for i := 0; i < s.NumberOfRunsPerState; i++ {
			params.InternalCount = i
			engine := PrepareSimulation(params, i)
			params = engine.Parameters
			s.OutputDir = engine.Parameters.StatisticsOutput.OutputDir

			engine.ProgressBar = s.ProgressBar
			s.StartEngine(engine)
		}
	}

	close(newParamsChan)
	for i := range newParamsChan {
		params = i
	}

	s.ProgressBar.Incr()

	// CUMULATIVE STATISTICS
	simulationBestIndividuals, err := s.SimulationBestIndividuals(params)
	if err != nil {
		return params, nil
	}
	err = simulationBestIndividuals.ToCSV(s.generateSimulationPathCSV("best-combined"))
	s.ProgressBar.Incr()
	simulationStrategy, err := s.SimulationBestStrategy(params)
	if err != nil {
		return params, nil
	}
	err = simulationStrategy.ToCSV(s.generateSimulationPathCSV("strategybest"))
	s.ProgressBar.Incr()
	simulationBestIndividual, err := s.SimulationBestIndividual(params)
	if err != nil {
		return params, nil
	}
	err = simulationBestIndividual.ToCSV(s.generateSimulationPathCSV("best-all"))
	s.ProgressBar.Incr()
	abs, _ := filepath.Abs(s.DataPath)
	if params.RunStats {
		s.RunRScript(s.RPath, abs, s.StatsFiles, params.LoggingChan, params.ErrorChan)
	}

	//msg := fmt.Sprintf("SIMULATION COMPLETE:\nFile: %s", params.ToString())
	s.ProgressBar.Incr()

	fmt.Printf("SIMUlATION COMPLETE: Number of Goroutines %d", runtime.NumGoroutine())
	return params, nil
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
	TopAntagonist            evolution.Individual
	TopProtagonist           evolution.Individual
	TopAntagonistByDelta     evolution.Individual
	TopProtagonistByDelta    evolution.Individual
	TopAntagonistByDeltaAvg  evolution.Individual
	TopProtagonistByDeltaAvg evolution.Individual
	TopAntagonistGeneration  int
	TopProtagonistGeneration int
	FinalAntagonist          evolution.Individual
	FinalProtagonist         evolution.Individual
	Generational             evolution.Generational
	MeanCorrelation          float64
	MeanCovariance           float64
}

func (s *Simulation) RunRScript(RPath, dirPath string, RFiles []string, logChan chan evolog.Logger, errChan chan error) {
	wg := sync.WaitGroup{}
	for _, rFile := range RFiles {
		wg.Add(1)
		go func(group *sync.WaitGroup, rFile string, logChan chan evolog.Logger, errChan chan error) {
			defer group.Done()

			fqdn := fmt.Sprintf("%s/%s", RPath, rFile)
			cmd := exec.Command("Rscript", fqdn, dirPath)
			//msg := fmt.Sprintf("Rscript: \n%s\n", cmd.String())
			//logChan <- msg

			err := cmd.Run()
			if err != nil {
				errChan <- err
			}
			s.ProgressBar.Incr()
		}(&wg, rFile, logChan, errChan)
	}
	wg.Wait()
	fmt.Println("COMPLETED RSCRIPTS")
}

func (s *Simulation) StartEngine(engine *evolution.EvolutionEngine) error {
	//mut := sync.Mutex{}
	evolutionResult, err := engine.TopologyRoundRobin()
	if err != nil {
		return err
	}

	engine.ProgressBar.Incr()

	topAnt, err := evolutionResult.TopAntagonistInRun.Clone()
	topProt, err := evolutionResult.TopProtagonistInRun.Clone()
	finAnt, err := evolutionResult.FinalAntagonist.Clone()
	finProt, err := evolutionResult.FinalProtagonist.Clone()

	s.SimulationStats[engine.Parameters.InternalCount] = SimulationRunStats{
		TopAntagonist:            topAnt,
		TopProtagonist:           topProt,
		TopAntagonistGeneration:  0,
		TopProtagonistGeneration: 0,
		FinalAntagonist:          finAnt,
		FinalProtagonist:         finProt,
		Generational:             evolutionResult.Generational,
		MeanCorrelation:          evolutionResult.Correlation,
		MeanCovariance:           evolutionResult.Covariance,
	}

	os.Mkdir(fmt.Sprintf("%s%d", s.OutputDir, s.CurrentEvolutionState.InternalCount), 0755)

	runEpochalStatistics, err := s.EpochalInRun(engine.Parameters)
	if err != nil {
		return err
	}
	err = runEpochalStatistics.ToCSV(s.generateRunPathCSV("epochal", engine.Parameters.InternalCount))

	engine.ProgressBar.Incr()
	runGenerationalStatistics, err := s.GenerationalInRun(engine.Parameters)
	if err != nil {
		return err
	}
	err = runGenerationalStatistics.ToCSV(s.generateRunPathCSV("generational", engine.Parameters.InternalCount))
	engine.ProgressBar.Incr()
	runStrategyStatistics, err := s.StrategyInRun(engine.Parameters)
	if err != nil {
		return err
	}
	err = runStrategyStatistics.ToCSV(s.generateRunPathCSV("strategy", engine.Parameters.InternalCount))

	engine.ProgressBar.Incr()
	runBestIndividualStatistics, err := s.BestIndividualsInRun(engine.Parameters)
	if err != nil {
		return err
	}
	err = runBestIndividualStatistics.ToCSV(s.generateRunPathCSV("best", engine.Parameters.InternalCount))
	if err != nil {
		return err
	}

	engine.ProgressBar.Incr()

	return nil
}

func (s *Simulation) generateRunPathCSV(fileName string, run int) string {
	path := fmt.Sprintf("%s/%s-%d.csv", s.DataPath, fileName, run)

	return path
}

func (s *Simulation) generateSimulationPathCSV(fileName string) string {
	path := fmt.Sprintf("%s/%s.csv", s.DataPath, fileName)

	return path
}

// PrepareSimulation takes in the given evolution parameters and a count variable and returns the engine that can be
// started run the simulation. The evolution engine will run count times.
func PrepareSimulation(params evolution.EvolutionParams, count int) *evolution.EvolutionEngine {
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

	params.SpecParam.ExpressionParsed = starterTreeAsMathematicalExpression
	startProgram := evolution.Program{
		T: &starterTree,
	}
	var spec evolution.SpecMulti

	switch params.FitnessStrategy.Type {
	case evolution.FitnessMonoThresholdedRatio:
		spec, err = evolution.GenerateSpecSimple(params.SpecParam, params.FitnessStrategy)
		if err != nil {
			log.Fatalf("MAIN | failed to create a valid spec | %s", err.Error())
		}
	case evolution.FitnessDualThresholdedRatio:
		spec, err = evolution.GenerateSpecSimple(params.SpecParam, params.FitnessStrategy)
		if err != nil {
			log.Fatalf("MAIN | failed to create a valid spec | %s", err.Error())
		}
	default:
		spec, err = evolution.GenerateSpecSimple(params.SpecParam, params.FitnessStrategy)
		if err != nil {
			log.Fatalf("MAIN | failed to create a valid spec | %s", err.Error())
		}
	}

	fmt.Printf(
		"Simulation:\n"+
			"Mathematical Expression: %s",
		starterTreeAsMathematicalExpression,
	)

	// Set extra params
	params.Spec = spec
	params.StartIndividual = startProgram
	params.SpecParam.AvailableSymbolicExpressions.Terminals = append(variableTerminals, constantTerminals...)
	params.SpecParam.AvailableSymbolicExpressions.NonTerminals = nonTerminals

	genCount := 0

	if params.MaxGenerations > evolution.MinAllowableGenerationsToTerminate {
		genCount = params.MaxGenerations
	} else {
		genCount = params.GenerationsCount
	}

	engine := &evolution.EvolutionEngine{
		Parameters:  params,
		Generations: make([]*evolution.Generation, genCount),
	}

	// ########################### OUTPUT STATISTICS  #######################################################3
	fmt.Printf("\nGeneration Count: %d\n", genCount)
	fmt.Printf("Each Individual Count: %d\n", engine.Parameters.EachPopulationSize)
	fmt.Printf("Iteration Count: (%d)\n", count)

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
	//timeStr := engine.Parameters.InternalCount

	folder = fmt.Sprintf("data/%s/", engine.Parameters.ToString())
	outputPath = fmt.Sprintf("%s", folder)

	engine.Parameters.StatisticsOutput.OutputDir = folder
	engine.Parameters.StatisticsOutput.OutputPath = outputPath
	return engine
	// ########################### START THE EVOLUTION PROCESS ##################################################3
}

// SpewJSON enables the creation of multiple JSON files containing parameter information.
// baseRelDir is the relative directory to the parameter folder. Should be within the project.
// Split should be the number of files per folder. It will try split them evenly
func (s *Simulation) SpewJSONNoSplit(projectAbsolutePath, paramsDir string) error {
	s.NumberOfRunsPerState = 20
	//os.Mkdir(baseRelDir, 0775)

	counter := 0

	for expressionIndex := 0; expressionIndex < len(AllExpressions); expressionIndex++ {
		for allMaxGenIndex := 0; allMaxGenIndex < len(AllMaxGenerations); allMaxGenIndex++ {
			for rangesIndex := 0; rangesIndex < len(AllRanges); rangesIndex++ {
				for seedIndex := 0; seedIndex < len(AllSeed); seedIndex++ {
					for generationsCountIndex := 0; generationsCountIndex < len(AllGenerationsCount); generationsCountIndex++ {
						for eachPopulationIndex := 0; eachPopulationIndex < len(AllEachPopulationSize); eachPopulationIndex++ {
							for reproductionIndex := 0; reproductionIndex < len(AllReproduction); reproductionIndex++ {
								for newTreeIndex := 0; newTreeIndex < len(AllDepthOfRandomNewTree); newTreeIndex++ {
									for antStratIndex := 0; antStratIndex < len(AllAntagonistStrategyCount); antStratIndex++ {
										for minProIndex := 0; minProIndex < len(
											AllMinimumTopProtagonistMeanBeforeTerminate); minProIndex++ {
											for proStratIndex := 0; proStratIndex < len(
												AllProtagonistStrategyCount); proStratIndex++ {
												for minGenIndex := 0; minGenIndex < len(
													AllMinimumGenerationMeanBeforeTerminate); minGenIndex++ {
													//for availableStrategyIndex = 0; availableStrategyIndex < len(AllAvailableStrategy); availableStrategyIndex++ {
													for fitnessStrategyTypeIndex := 0; fitnessStrategyTypeIndex < len(AllFitnessStrategyType); fitnessStrategyTypeIndex++ {
														for fitStratAntThresMultIndex := 0; fitStratAntThresMultIndex < len(AllFitStratAntThreshMult); fitStratAntThresMultIndex++ {
															for fitStratProtThreshMultIndex := 0; fitStratProtThreshMultIndex < len(AllFitStratProThreshMult); fitStratProtThreshMultIndex++ {
																for strategiesAntagonistIndex := 0; strategiesAntagonistIndex < len(AllPossibleStrategies); strategiesAntagonistIndex++ {
																	for strategiesProtagonistIndex := 0; strategiesProtagonistIndex < len(AllPossibleStrategies); strategiesProtagonistIndex++ {
																		for divByZeroPenaltyIndex := 0; divByZeroPenaltyIndex < len(AllDivByZeroPenalty); divByZeroPenaltyIndex++ {
																			for divByZeroStrategyIndex := 0; divByZeroStrategyIndex < len(AllDivByZeroStrategy); divByZeroStrategyIndex++ {
																				for survivorIndex := 0; survivorIndex < len(AllSurvivorSelection); survivorIndex++ {
																					for parentIndex := 0; parentIndex < len(AllSelectionParentType); parentIndex++ {
																						for ProtagonistMinGenAvgFitIndex := 0; ProtagonistMinGenAvgFitIndex < len(AllProtagonistMinGenAvgFit); ProtagonistMinGenAvgFitIndex++ {
																							// TODO Add Parallelism
																							params := evolution.EvolutionParams{
																								SpecParam: evolution.SpecParam{
																									Seed:       AllSeed[seedIndex],
																									Expression: AllExpressions[expressionIndex],
																									Range:      AllRanges[rangesIndex],
																									AvailableVariablesAndOperators: evolution.AvailableVariablesAndOperators{
																										Constants: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
																										Variables: []string{"x"},
																										Operators: []string{"*", "+", "-", "/"},
																									},
																									DivideByZeroStrategy: AllDivByZeroStrategy[divByZeroStrategyIndex],
																									DivideByZeroPenalty:  AllDivByZeroPenalty[divByZeroPenaltyIndex],
																								},
																								MaxGenerations:                           AllMaxGenerations[allMaxGenIndex],
																								MinimumTopProtagonistMeanBeforeTerminate: AllMinimumTopProtagonistMeanBeforeTerminate[minProIndex],
																								MinimumGenerationMeanBeforeTerminate:     AllMinimumGenerationMeanBeforeTerminate[minGenIndex],
																								ProtagonistMinGenAvgFit:                  AllProtagonistMinGenAvgFit[ProtagonistMinGenAvgFitIndex],
																								GenerationsCount:                         AllGenerationsCount[generationsCountIndex],
																								EachPopulationSize:                       AllEachPopulationSize[eachPopulationIndex],
																								EnableParallelism:                        false,
																								Strategies: evolution.Strategies{
																									AntagonistAvailableStrategies:  AllPossibleStrategies[strategiesAntagonistIndex],
																									ProtagonistAvailableStrategies: AllPossibleStrategies[strategiesProtagonistIndex],
																									AntagonistStrategyCount:        AllAntagonistStrategyCount[antStratIndex],
																									ProtagonistStrategyCount:       AllProtagonistStrategyCount[proStratIndex],
																									DepthOfRandomNewTrees:          AllDepthOfRandomNewTree[newTreeIndex],
																								},
																								FitnessStrategy: evolution.FitnessStrategy{
																									Type:                           AllFitnessStrategyType[fitnessStrategyTypeIndex],
																									AntagonistThresholdMultiplier:  AllFitStratAntThreshMult[fitStratAntThresMultIndex],
																									ProtagonistThresholdMultiplier: AllFitStratProThreshMult[fitStratProtThreshMultIndex],
																								},
																								Reproduction: AllReproduction[reproductionIndex],
																								Selection: evolution.Selection{
																									Survivor: AllSurvivorSelection[survivorIndex],
																									Parent:   AllSelectionParentType[parentIndex],
																								},
																								FolderPercentages: []float64{0.01, 0.25, 0.5, 0.75, 0.95},
																							}

																							fmt.Printf("Loaded: %s\n", params.ToString())

																							engine := PrepareSimulation(params, 0)
																							outputPath := fmt.Sprintf("%s%s",
																								engine.Parameters.
																									StatisticsOutput.
																									OutputPath[:len(engine.Parameters.StatisticsOutput.OutputPath)-1], ".json")

																							folder := fmt.Sprintf("%s/%d",
																								paramsDir, counter)
																							absFolderPath := fmt.Sprintf("%s/%s", projectAbsolutePath, folder)
																							os.Mkdir(absFolderPath, 0755)
																							outputFilepath := strings.ReplaceAll(
																								outputPath, "data/", folder+"/")

																							file, err := os.Create(outputFilepath)
																							if err != nil {
																								return fmt.Errorf(err.Error())
																							}
																							err = json.NewEncoder(file).Encode(params)
																							if err != nil {
																								return fmt.Errorf(err.Error())
																							}
																							file.Close()

																							counter++
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
				}
			}
		}
	}

	fmt.Printf("WROTE %d files", counter)
	return nil
}

func WriteIndexProgressToFile(indexProgress IndexProgress, indexFile string) error {
	file, err := os.Create(indexFile)
	if err != nil {
		return err
	}
	return json.NewEncoder(file).Encode(indexProgress)
	file.Close()
	return err
}

type IndexProgress struct {
	ExpressionIndex              int `expressionIndex`
	RangesIndex                  int `json:"rangesIndex"`
	SeedIndex                    int `json:"seedIndex"`
	GenerationsCountIndex        int `json:"generationsCountIndex"`
	EachPopulationIndex          int `json:"eachPopulationIndex"`
	ReproductionIndex            int `json:"reproductionIndex"`
	AllDepthOfRandomNewTreeIndex int `json:"allDepthOfRandomNewTreeIndex"`
	AntagonistStrategyCountIndex int `json:"antagonistStrategyCountIndex"`
	ProStratCountInd             int `json:"protagonisttStrategyCountIndex"`
	FitnessStrategyTypeIndex     int `json:"fitnessStrategyTypeIndex"`
	FitStratAntThresMultIndex    int `json:"fitStratAntThresMultIndex"`
	FitStratProtThreshMultIndex  int `json:"fitStratProtThreshMultIndex"`
	SelectParentTypeIndex        int `json:"selectParentTypeIndex"`
	SelectSurvivorPercentIndex   int `json:"selectSurvivorPercentIndex"`
	StrategiesAntagonistIndex    int `json:"strategiesAntagonistIndex"`
	StrategiesProtagonistIndex   int `json:"strategiesProtagonistIndex"`
	DivByZeroStrategy            int `json:"divByZeroStrategy"`
	DivByZeroPenalty             int `json:"divByZeroPenalty"`
	TournamentSelectionIndex     int `json:"tournamentSelectionIndex"`
	NumberOfRunsPerState         int `json:"numberOfRunsPerState"`
}
