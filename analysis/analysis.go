package analysis

import (
	"encoding/json"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/martinomburajr/masters-go/evolution"
	"github.com/martinomburajr/masters-go/simulation"
	"os"
	"path/filepath"
	"strings"
)

func CombineBest(baseFolder string) error {
	accCSV := make([]CSVBestAll, 0)
	if baseFolder == "" {
		return fmt.Errorf("baseFolder cannot be empty")
	}

	dataFolders, err := RetrieveDataFolders(baseFolder)
	if err != nil {
		return err
	}

	for _, dataFolder := range dataFolders {
		params, err := GetParams(dataFolder)
		if err != nil {
			return err
		}

		err = filepath.Walk(dataFolder, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				if strings.Contains(path, "best-all.csv") && !strings.Contains(path, ".png") {
					bestAllFile, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
					if err != nil {
						return err
					}
					defer bestAllFile.Close()

					bestIndividualStatistics := []*simulation.RunBestIndividualStatistic{}
					err = gocsv.Unmarshal(bestAllFile, &bestIndividualStatistics)
					if err != nil {
						return err
					}

					if len(bestIndividualStatistics) > 0 {
						bst := *bestIndividualStatistics[0]
						idSplit := strings.Split(dataFolder, "/")

						// #### COMBINE
						csvBest := CSVBestAll{
							FileID:       idSplit[len(idSplit)-1],
							SpecEquation: bst.SpecEquation,

							// INDIVIDUAL
							Antagonist:                 bst.Antagonist,
							Protagonist:                bst.Protagonist,
							AntagonistAge:              bst.AntagonistAge,
							AntagonistAverageDelta:     bst.AntagonistAverageDelta,
							AntagonistBestDelta:        bst.AntagonistBestDelta,
							AntagonistBestFitness:      bst.AntagonistBestFitness,
							AntagonistBirthGen:         bst.AntagonistBirthGen,
							AntagonistDominantStrategy: bst.AntagonistDominantStrategy,
							AntagonistEquation:         bst.AntagonistEquation,
							AntagonistGeneration:       bst.AntagonistGeneration,
							AntagonistID:               bst.AntagonistID,
							AntagonistStdDev:           bst.AntagonistStdDev,
							AntagonistStrategy:         bst.AntagonistStrategy,
							AntagonistNoOComp:          bst.AntagonistNoOfCompetitions,

							ProtagonistAge:              bst.ProtagonistAge,
							ProtagonistAverageDelta:     bst.ProtagonistAverageDelta,
							ProtagonistBestDelta:        bst.ProtagonistBestDelta,
							ProtagonistBestFitness:      bst.ProtagonistBestFitness,
							ProtagonistBirthGen:         bst.ProtagonistBirthGen,
							ProtagonistDominantStrategy: bst.ProtagonistDominantStrategy,
							ProtagonistEquation:         bst.ProtagonistEquation,
							ProtagonistGeneration:       bst.ProtagonistGeneration,
							ProtagonistID:               bst.ProtagonistID,
							ProtagonistStdDev:           bst.ProtagonistStdDev,
							ProtagonistStrategy:         bst.ProtagonistStrategy,
							ProtagonistNoOComp:          bst.ProtagonistNoOfCompetitions,


							FinalAntagonist:                 bst.FinalAntagonist,
							FinalAntagonistAge:              bst.FinalAntagonistAge,
							FinalAntagonistAverageDelta:     bst.FinalAntagonistAverageDelta,
							FinalAntagonistBestDelta:        bst.FinalAntagonistBestDelta,
							FinalAntagonistBestFitness:      bst.FinalAntagonistBestFitness,
							FinalAntagonistBirthGen:         bst.FinalAntagonistBirthGen,
							FinalAntagonistDominantStrategy: bst.FinalAntagonistDominantStrategy,
							FinalAntagonistEquation:         bst.FinalAntagonistEquation,
							FinalAntagonistStdDev:           bst.FinalAntagonistStdDev,
							FinalAntagonistStrategy:         bst.FinalAntagonistStrategy,
							FinalAntagonistNoOComp:          bst.FinalAntagonistNoOfCompetitions,

							FinalProtagonist:                 bst.FinalProtagonist,
							FinalProtagonistAge:              bst.FinalProtagonistAge,
							FinalProtagonistAverageDelta:     bst.FinalProtagonistAverageDelta,
							FinalProtagonistBestDelta:        bst.FinalProtagonistBestDelta,
							FinalProtagonistBestFitness:      bst.FinalProtagonistBestFitness,
							FinalProtagonistBirthGen:         bst.FinalProtagonistBirthGen,
							FinalProtagonistDominantStrategy: bst.FinalProtagonistDominantStrategy,
							FinalProtagonistEquation:         bst.FinalProtagonistEquation,
							FinalProtagonistStdDev:           bst.FinalProtagonistStdDev,
							FinalProtagonistStrategy:         bst.FinalProtagonistStrategy,
							FinalProtagonistNoOComp:          bst.FinalProtagonistNoOfCompetitions,

							// PARAMS
							SpecRange:          params.SpecParam.Range,
							SpecSeed:           params.SpecParam.Seed,
							TopologyType:       params.Topology.Type,
							GenerationCount:    params.GenerationsCount,
							EachPopulationSize: params.EachPopulationSize,
							ParentSelect:       params.Selection.Parent.Type,
							SurvivorSelect: params.Selection.Survivor.Type,
							CrossPercent: params.Reproduction.CrossoverPercentage,
							ProbMutation: params.Reproduction.ProbabilityOfMutation,
							AntStratCount: params.Strategies.AntagonistStrategyCount,
							AntStrat: evolution.StrategiesToStringArr(evolution.ConvertStrategiesToString(params.Strategies.
								AntagonistAvailableStrategies)),
							AntThreshMult: params.FitnessStrategy.AntagonistThresholdMultiplier,

							ProThresMult:  params.FitnessStrategy.ProtagonistThresholdMultiplier,
							ProStratCount: params.Strategies.ProtagonistStrategyCount,
							ProStrat: evolution.StrategiesToStringArr(evolution.ConvertStrategiesToString(params.Strategies.
								ProtagonistAvailableStrategies)),

							RandTreeDepth: params.Strategies.DepthOfRandomNewTrees,
							DivByZero:     params.SpecParam.DivideByZeroStrategy,
							DivByZeroPen:  params.SpecParam.DivideByZeroPenalty,
						}
						(accCSV) = append(accCSV, csvBest)
						return err
					}
				}
			}
			return err
		})
		if err != nil {
			return err
		}
	}

	finalCSV := make([]CSVBestAll, 0)
	for i := range accCSV {
		if i%2 == 0 {
			finalCSV = append(finalCSV, accCSV[i])
		}
	}

	outputFilePath := fmt.Sprintf("%s/%s", baseFolder, "coalescedBest.csv")
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	err = gocsv.MarshalFile(finalCSV, outputFile)
	if err != nil {
		return err
	}

	return nil
}

func CombineGenerations(baseFolder string) error {
	// REQUIRE
	if baseFolder == "" {
		return fmt.Errorf("baseFolder cannot be empty")
	}

	dataFolders, err := RetrieveDataFolders(baseFolder)
	if err != nil {
		return err
	}

	totalDirs := -1
	AllCSVGens := make([]CSVCombinedGenerations, 0)
	for _, dataFolder := range dataFolders {
		generationInFolder := make([][]*simulation.RunGenerationalStatistic, 0)
		totalDirs++
		params, err := GetParams(dataFolder)
		if err != nil {
			return err
		}
		csvGens := make([]CSVCombinedGenerations, 0)

		err = filepath.Walk(dataFolder, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				if strings.Contains(path, "generational") && !strings.Contains(path, ".png") {
					generationPath, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
					if err != nil {
						//return err
					}
					defer generationPath.Close()

					generationalStatistics := []*simulation.RunGenerationalStatistic{}
					err = gocsv.Unmarshal(generationPath, &generationalStatistics)
					if err != nil {
						//return err
					}

					for _, gen := range generationalStatistics {
						csvGen := CSVCombinedGenerations{
							FileID: dataFolder,

							SpecEquation: gen.SpecEquation,
							Generation:   gen.Generation,

							// PARAMS
							SpecRange:          params.SpecParam.Range,
							SpecSeed:           params.SpecParam.Seed,
							TopologyType:       params.Topology.Type,
							GenerationCount:    params.GenerationsCount,
							EachPopulationSize: params.EachPopulationSize,
							ParentSelect:       params.Selection.Parent.Type,

							SurvivorSelect: params.Selection.Survivor.Type,

							CrossPercent: params.Reproduction.CrossoverPercentage,
							ProbMutation: params.Reproduction.ProbabilityOfMutation,

							AntStratCount: params.Strategies.AntagonistStrategyCount,
							AntStrat: evolution.StrategiesToStringArr(evolution.ConvertStrategiesToString(params.Strategies.
								AntagonistAvailableStrategies)),
							AntThreshMult: params.FitnessStrategy.AntagonistThresholdMultiplier,

							ProThresMult:  params.FitnessStrategy.ProtagonistThresholdMultiplier,
							ProStratCount: params.Strategies.ProtagonistStrategyCount,
							ProStrat: evolution.StrategiesToStringArr(evolution.ConvertStrategiesToString(params.Strategies.
								ProtagonistAvailableStrategies)),

							RandTreeDepth: params.Strategies.DepthOfRandomNewTrees,
							DivByZero:     params.SpecParam.DivideByZeroStrategy,
							DivByZeroPen:  params.SpecParam.DivideByZeroPenalty,

							Correlation:  gen.Correlation,
							Covariance:   gen.Covariance,
							TopAEquation: gen.AntagonistEquation,
							TopPEquation: gen.ProtagonistEquation,

							Antagonist:                gen.AntagonistMean,
							Protagonist:               gen.ProtagonistMean,
							TopAntagonistStdDev:       gen.AntagonistStdDev,
							TopProtagonistStdDev:      gen.ProtagonistStdDev,
							TopAntagonistVar:          gen.AntagonistVariance,
							TopProtagonistVar:         gen.ProtagonistVariance,
							TopAntagonistSkew:         gen.AntagonistSkew,
							TopProtagonistSkew:        gen.ProtagonistSkew,
							TopAntagonistKurtosis:     gen.AntagonistExKurtosis,
							TopProtagonistKurtosis:    gen.ProtagonistExKurtosis,
							TopAntagonistMean:         gen.TopAntagonistMeanFitness,
							TopProtagonistMean:        gen.TopProtagonistMeanFitness,
							TopAntagonistBestFitness:  gen.AntagonistBestFitness,
							TopProtagonistBestFitness: gen.ProtagonistBestFitness,

							TopAntagonistAverageDelta:      gen.AntagonistAverageDelta,
							TopProtagonistAverageDelta:     gen.ProtagonistAverageDelta,
							TopAntagonistBestDelta:         gen.AntagonistBestDelta,
							TopProtagonistBestDelta:        gen.ProtagonistBestDelta,
							TopAntagonistStrategy:          gen.AntagonistStrategy,
							TopProtagonistStrategy:         gen.ProtagonistStrategy,
							TopAntagonistDominantStrategy:  gen.AntagonistDominantStrategy,
							TopProtagonistDominantStrategy: gen.ProtagonistDominantStrategy,
							TopAntagonistBirthGen:          gen.AntagonistBirthGen,
							TopProtagonistBirthGen:         gen.ProtagonistBirthGen,
							TopAntagonistAge:               gen.AntagonistAge,
							TopProtagonistAge:              gen.ProtagonistAge,

							Run: gen.Run,
						}
						csvGens = append(csvGens, csvGen)
						AllCSVGens = append(AllCSVGens, csvGen)
					}

					generationInFolder = append(generationInFolder, generationalStatistics)
				}
			}
			return err
		})
		if err != nil {
			//return err
		}

		print(totalDirs)
		//err = writeToFolder(dataFolder, "coalescedGenerations.csv", csvGens)
		//if err != nil {
		//	//return err
		//}
	}

	err = writeToBaseFolder(baseFolder, "coalescedAllGenerations.csv", AllCSVGens)
	if err != nil {
		return err
	}

	fmt.Printf("TOTAL DIRS: %d", totalDirs)

	return err
}

func CombineBestCombinedAll(baseFolder string) error {
	accCSV := make([]CSVBestAll, 0)
	if baseFolder == "" {
		return fmt.Errorf("baseFolder cannot be empty")
	}

	dataFolders, err := RetrieveDataFolders(baseFolder)
	if err != nil {
		return err
	}

	for _, dataFolder := range dataFolders {
		params, err := GetParams(dataFolder)
		if err != nil {
			return err
		}

		err = filepath.Walk(dataFolder, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				if strings.Contains(path, "best-combined.csv") && !strings.Contains(path, ".png") {
					bestAllFile, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
					if err != nil {
						return err
					}
					defer bestAllFile.Close()

					bestIndividualStatistics := []*simulation.RunBestIndividualStatistic{}
					err = gocsv.Unmarshal(bestAllFile, &bestIndividualStatistics)
					if err != nil {
						return err
					}

					if len(bestIndividualStatistics) > 0 {
						for _, bst := range bestIndividualStatistics {
							idSplit := strings.Split(dataFolder, "/")

							// #### COMBINE
							csvBest := CSVBestAll{
								FileID:       idSplit[len(idSplit)-1],
								SpecEquation: bst.SpecEquation,

								// INDIVIDUAL
								Antagonist:                 bst.Antagonist,
								Protagonist:                bst.Protagonist,
								AntagonistAge:              bst.AntagonistAge,
								AntagonistAverageDelta:     bst.AntagonistAverageDelta,
								AntagonistBestDelta:        bst.AntagonistBestDelta,
								AntagonistBestFitness:      bst.AntagonistBestFitness,
								AntagonistBirthGen:         bst.AntagonistBirthGen,
								AntagonistDominantStrategy: bst.AntagonistDominantStrategy,
								AntagonistEquation:         bst.AntagonistEquation,
								AntagonistGeneration:       bst.AntagonistGeneration,
								AntagonistID:               bst.AntagonistID,
								AntagonistStdDev:           bst.AntagonistStdDev,
								AntagonistStrategy:         bst.AntagonistStrategy,
								AntagonistNoOComp:          bst.AntagonistNoOfCompetitions,

								ProtagonistAge:              bst.ProtagonistAge,
								ProtagonistAverageDelta:     bst.ProtagonistAverageDelta,
								ProtagonistBestDelta:        bst.ProtagonistBestDelta,
								ProtagonistBestFitness:      bst.ProtagonistBestFitness,
								ProtagonistBirthGen:         bst.ProtagonistBirthGen,
								ProtagonistDominantStrategy: bst.ProtagonistDominantStrategy,
								ProtagonistEquation:         bst.ProtagonistEquation,
								ProtagonistGeneration:       bst.ProtagonistGeneration,
								ProtagonistID:               bst.ProtagonistID,
								ProtagonistStdDev:           bst.ProtagonistStdDev,
								ProtagonistStrategy:         bst.ProtagonistStrategy,
								ProtagonistNoOComp:          bst.ProtagonistNoOfCompetitions,


								FinalAntagonist:                 bst.FinalAntagonist,
								FinalAntagonistAge:              bst.FinalAntagonistAge,
								FinalAntagonistAverageDelta:     bst.FinalAntagonistAverageDelta,
								FinalAntagonistBestDelta:        bst.FinalAntagonistBestDelta,
								FinalAntagonistBestFitness:      bst.FinalAntagonistBestFitness,
								FinalAntagonistBirthGen:         bst.FinalAntagonistBirthGen,
								FinalAntagonistDominantStrategy: bst.FinalAntagonistDominantStrategy,
								FinalAntagonistEquation:         bst.FinalAntagonistEquation,
								FinalAntagonistStdDev:           bst.FinalAntagonistStdDev,
								FinalAntagonistStrategy:         bst.FinalAntagonistStrategy,
								FinalAntagonistNoOComp:          bst.FinalAntagonistNoOfCompetitions,

								FinalProtagonist:                 bst.FinalProtagonist,
								FinalProtagonistAge:              bst.FinalProtagonistAge,
								FinalProtagonistAverageDelta:     bst.FinalProtagonistAverageDelta,
								FinalProtagonistBestDelta:        bst.FinalProtagonistBestDelta,
								FinalProtagonistBestFitness:      bst.FinalProtagonistBestFitness,
								FinalProtagonistBirthGen:         bst.FinalProtagonistBirthGen,
								FinalProtagonistDominantStrategy: bst.FinalProtagonistDominantStrategy,
								FinalProtagonistEquation:         bst.FinalProtagonistEquation,
								FinalProtagonistStdDev:           bst.FinalProtagonistStdDev,
								FinalProtagonistStrategy:         bst.FinalProtagonistStrategy,
								FinalProtagonistNoOComp:          bst.FinalProtagonistNoOfCompetitions,

								// PARAMS
								SpecRange:          params.SpecParam.Range,
								SpecSeed:           params.SpecParam.Seed,
								TopologyType:       params.Topology.Type,
								GenerationCount:    params.GenerationsCount,
								EachPopulationSize: params.EachPopulationSize,
								ParentSelect:       params.Selection.Parent.Type,
								SurvivorSelect: params.Selection.Survivor.Type,
								CrossPercent: params.Reproduction.CrossoverPercentage,
								ProbMutation: params.Reproduction.ProbabilityOfMutation,
								AntStratCount: params.Strategies.AntagonistStrategyCount,
								AntStrat: evolution.StrategiesToStringArr(evolution.ConvertStrategiesToString(params.Strategies.
									AntagonistAvailableStrategies)),
								AntThreshMult: params.FitnessStrategy.AntagonistThresholdMultiplier,

								ProThresMult:  params.FitnessStrategy.ProtagonistThresholdMultiplier,
								ProStratCount: params.Strategies.ProtagonistStrategyCount,
								ProStrat: evolution.StrategiesToStringArr(evolution.ConvertStrategiesToString(params.Strategies.
									ProtagonistAvailableStrategies)),

								RandTreeDepth: params.Strategies.DepthOfRandomNewTrees,
								DivByZero:     params.SpecParam.DivideByZeroStrategy,
								DivByZeroPen:  params.SpecParam.DivideByZeroPenalty,
							}
							(accCSV) = append(accCSV, csvBest)
						}
						return err
					}
				}
			}
			return err
		})
		if err != nil {
			return err
		}
	}

	finalCSV := make([]CSVBestAll, 0)
	for i := range accCSV {
		if i%2 == 0 {
			finalCSV = append(finalCSV, accCSV[i])
		}
	}

	outputFilePath := fmt.Sprintf("%s/%s", baseFolder, "coalescedBestCombined.csv")
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	err = gocsv.MarshalFile(finalCSV, outputFile)
	if err != nil {
		return err
	}

	return nil
}

func CombineBestCombinedAll2(baseFolder string) error {
	accCSV := make([]CSVBestAll, 0)
	if baseFolder == "" {
		return fmt.Errorf("baseFolder cannot be empty")
	}

	dataFolders, err := RetrieveDataFolders(baseFolder)
	if err != nil {
		return err
	}

	for _, dataFolder := range dataFolders {
		params, err := GetParams(dataFolder)
		if err != nil {
			return err
		}

		err = filepath.Walk(dataFolder, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				if strings.Contains(path, "best-") && !strings.Contains(path, ".png") && !strings.Contains(path,
					"all") && !strings.Contains(path, "combined"){
					bestAllFile, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
					if err != nil {
						return err
					}
					defer bestAllFile.Close()

					bestIndividualStatistics := []*simulation.RunBestIndividualStatistic{}
					err = gocsv.Unmarshal(bestAllFile, &bestIndividualStatistics)
					if err != nil {
						return err
					}

					if len(bestIndividualStatistics) > 0 {
						for _, bst := range bestIndividualStatistics {
							idSplit := strings.Split(dataFolder, "/")

							// #### COMBINE
							csvBest := CSVBestAll{
								FileID:       idSplit[len(idSplit)-1],
								SpecEquation: bst.SpecEquation,

								// INDIVIDUAL
								Antagonist:                 bst.Antagonist,
								Protagonist:                bst.Protagonist,
								AntagonistAge:              bst.AntagonistAge,
								AntagonistAverageDelta:     bst.AntagonistAverageDelta,
								AntagonistBestDelta:        bst.AntagonistBestDelta,
								AntagonistBestFitness:      bst.AntagonistBestFitness,
								AntagonistBirthGen:         bst.AntagonistBirthGen,
								AntagonistDominantStrategy: bst.AntagonistDominantStrategy,
								AntagonistEquation:         bst.AntagonistEquation,
								AntagonistGeneration:       bst.AntagonistGeneration,
								AntagonistID:               bst.AntagonistID,
								AntagonistStdDev:           bst.AntagonistStdDev,
								AntagonistStrategy:         bst.AntagonistStrategy,
								AntagonistNoOComp:          bst.AntagonistNoOfCompetitions,

								ProtagonistAge:              bst.ProtagonistAge,
								ProtagonistAverageDelta:     bst.ProtagonistAverageDelta,
								ProtagonistBestDelta:        bst.ProtagonistBestDelta,
								ProtagonistBestFitness:      bst.ProtagonistBestFitness,
								ProtagonistBirthGen:         bst.ProtagonistBirthGen,
								ProtagonistDominantStrategy: bst.ProtagonistDominantStrategy,
								ProtagonistEquation:         bst.ProtagonistEquation,
								ProtagonistGeneration:       bst.ProtagonistGeneration,
								ProtagonistID:               bst.ProtagonistID,
								ProtagonistStdDev:           bst.ProtagonistStdDev,
								ProtagonistStrategy:         bst.ProtagonistStrategy,
								ProtagonistNoOComp:          bst.ProtagonistNoOfCompetitions,


								FinalAntagonist:                 bst.FinalAntagonist,
								FinalAntagonistAge:              bst.FinalAntagonistAge,
								FinalAntagonistAverageDelta:     bst.FinalAntagonistAverageDelta,
								FinalAntagonistBestDelta:        bst.FinalAntagonistBestDelta,
								FinalAntagonistBestFitness:      bst.FinalAntagonistBestFitness,
								FinalAntagonistBirthGen:         bst.FinalAntagonistBirthGen,
								FinalAntagonistDominantStrategy: bst.FinalAntagonistDominantStrategy,
								FinalAntagonistEquation:         bst.FinalAntagonistEquation,
								FinalAntagonistStdDev:           bst.FinalAntagonistStdDev,
								FinalAntagonistStrategy:         bst.FinalAntagonistStrategy,
								FinalAntagonistNoOComp:          bst.FinalAntagonistNoOfCompetitions,

								FinalProtagonist:                 bst.FinalProtagonist,
								FinalProtagonistAge:              bst.FinalProtagonistAge,
								FinalProtagonistAverageDelta:     bst.FinalProtagonistAverageDelta,
								FinalProtagonistBestDelta:        bst.FinalProtagonistBestDelta,
								FinalProtagonistBestFitness:      bst.FinalProtagonistBestFitness,
								FinalProtagonistBirthGen:         bst.FinalProtagonistBirthGen,
								FinalProtagonistDominantStrategy: bst.FinalProtagonistDominantStrategy,
								FinalProtagonistEquation:         bst.FinalProtagonistEquation,
								FinalProtagonistStdDev:           bst.FinalProtagonistStdDev,
								FinalProtagonistStrategy:         bst.FinalProtagonistStrategy,
								FinalProtagonistNoOComp:          bst.FinalProtagonistNoOfCompetitions,

								// PARAMS
								SpecRange:          params.SpecParam.Range,
								SpecSeed:           params.SpecParam.Seed,
								TopologyType:       params.Topology.Type,
								GenerationCount:    params.GenerationsCount,
								EachPopulationSize: params.EachPopulationSize,
								ParentSelect:       params.Selection.Parent.Type,
								SurvivorSelect: params.Selection.Survivor.Type,
								CrossPercent: params.Reproduction.CrossoverPercentage,
								ProbMutation: params.Reproduction.ProbabilityOfMutation,
								AntStratCount: params.Strategies.AntagonistStrategyCount,
								AntStrat: evolution.StrategiesToStringArr(evolution.ConvertStrategiesToString(params.Strategies.
									AntagonistAvailableStrategies)),
								AntThreshMult: params.FitnessStrategy.AntagonistThresholdMultiplier,

								ProThresMult:  params.FitnessStrategy.ProtagonistThresholdMultiplier,
								ProStratCount: params.Strategies.ProtagonistStrategyCount,
								ProStrat: evolution.StrategiesToStringArr(evolution.ConvertStrategiesToString(params.Strategies.
									ProtagonistAvailableStrategies)),

								RandTreeDepth: params.Strategies.DepthOfRandomNewTrees,
								DivByZero:     params.SpecParam.DivideByZeroStrategy,
								DivByZeroPen:  params.SpecParam.DivideByZeroPenalty,
							}
							(accCSV) = append(accCSV, csvBest)
						}
						return err
					}
				}
			}
			return err
		})
		if err != nil {
			return err
		}
	}

	finalCSV := make([]CSVBestAll, 0)
	for i := range accCSV {
		if i%2 == 0 {
			finalCSV = append(finalCSV, accCSV[i])
		}
	}

	outputFilePath := fmt.Sprintf("%s/%s", baseFolder, "coalescedBestCombined.csv")
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	err = gocsv.MarshalFile(finalCSV, outputFile)
	if err != nil {
		return err
	}

	return nil
}

func writeToFolder(folderpath, filename string, data []CSVCombinedGenerations) error {
	split := strings.Split(folderpath, "/")
	folderpath = fmt.Sprintf("%s/%s", "/home/martinomburajr/Desktop/Simulations", split[len(split)-1])
	os.MkdirAll(folderpath, 0777)
	outputFilePath := fmt.Sprintf("%s/%s", folderpath, filename)
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	err = gocsv.MarshalFile(data, outputFile)
	if err != nil {
		return err
	}
	return nil
}

func writeToBaseFolder(folderpath, filename string, data []CSVCombinedGenerations) error {
	os.MkdirAll(folderpath, 0777)
	outputFilePath := fmt.Sprintf("%s/%s", folderpath, filename)
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	err = gocsv.MarshalFile(data, outputFile)
	if err != nil {
		return err
	}
	return nil
}

func RetrieveDataFolders(baseFolder string) ([]string, error) {
	allFolders := make([]string, 0)
	err := filepath.Walk(baseFolder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			allFolders = append(allFolders, path)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	allFolders = allFolders[1:]
	outputFolders := make([]string, 0)
	for i := 0; i < len(allFolders); i += 2 {
		outputFolders = append(outputFolders, allFolders[i])
	}
	return outputFolders, nil
}

func GetParams(dataFolderPath string) (evolution.EvolutionParams, error) {
	paramsJsonPath := ""
	err := filepath.Walk(dataFolderPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if strings.Contains(path, "_params.json") {
				paramsJsonPath = path
				return err
			}
		}
		return err
	})
	if err != nil {
		return evolution.EvolutionParams{}, err
	}
	paramsJsonFile, err := os.OpenFile(paramsJsonPath, os.O_RDONLY, os.ModePerm)
	var params evolution.EvolutionParams
	err = json.NewDecoder(paramsJsonFile).Decode(&params)
	if err != nil {
		return evolution.EvolutionParams{}, err
	}
	return params, err
}

//func CreateFinalCSV(sampleFilePath, baseSimulationFolder string) error {
//
//}

// ReadFile will read in a .csv file. The filepath argument must be the path to the id folder e.g.
// ~/home/masters-go/_dataBackup/1222/sometopfolder/
// paramsFileName is the filename of the params.json e.g _params.json
//func ReadJSONParamsFile(filepath, paramsFileName string) error {
//
//}
