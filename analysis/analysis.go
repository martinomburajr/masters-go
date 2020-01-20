package analysis

import (
	"encoding/json"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/martinomburajr/masters-go/evolution"
	"github.com/martinomburajr/masters-go/simulation"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type CSVBestAll struct {
	FileID                  string `csv:"ID"`
	bestIndividualStatistic simulation.RunBestIndividualStatistic `csv:"bestIndividualStatistic"`
	params                  evolution.EvolutionParams `csv:"evolutionaryParams"`

	//BEST INDIVIDUAL
	SpecEquation string `csv:"specEquation"`
	SpecRange    int    `csv:"range"`
	SpecSeed     int    `csv:"seed"`

	AntagonistID string `csv:"AID"`
	ProtagonistID string `csv:"PID"`
	Antagonist                  float64 `csv:"AAvg"`
	Protagonist                 float64 `csv:"PAvg"`
	AntagonistBestFitness       float64 `csv:"ABestFit"`
	ProtagonistBestFitness      float64 `csv:"PBestFit"`
	AntagonistStdDev            float64 `csv:"AStdDev"`
	ProtagonistStdDev           float64 `csv:"PStdDev"`
	AntagonistAverageDelta      float64 `csv:"AAvgDelta"`
	ProtagonistAverageDelta     float64 `csv:"PAvgDelta"`
	AntagonistBestDelta         float64 `csv:"ABestDelta"`
	ProtagonistBestDelta        float64 `csv:"PBestDelta"`
	AntagonistEquation          string  `csv:"AEquation"`
	ProtagonistEquation         string  `csv:"PEquation"`
	AntagonistStrategy          string  `csv:"AStrat"`
	ProtagonistStrategy         string  `csv:"PStrat"`
	AntagonistDominantStrategy  string  `csv:"ADomStrat"`
	ProtagonistDominantStrategy string  `csv:"PDomStrat"`
	AntagonistGeneration        int     `csv:"AGen"`
	ProtagonistGeneration       int     `csv:"PGen"`
	AntagonistBirthGen          int     `csv:"ABirthGen"`
	ProtagonistBirthGen         int     `csv:"PBirthGen"`
	AntagonistAge               int     `csv:"AAge"`
	ProtagonistAge              int     `csv:"PAge"`

	FinalAntagonist                  float64 `csv:"finAAvg"`
	FinalProtagonist                 float64 `csv:"finPAvg"`
	FinalAntagonistBestFitness       float64 `csv:"finABestFit"`
	FinalProtagonistBestFitness      float64 `csv:"finPBestFit"`
	FinalAntagonistStdDev            float64 `csv:"finAStdDev"`
	FinalProtagonistStdDev           float64 `csv:"finPStdDev"`
	FinalAntagonistAverageDelta      float64 `csv:"finAAvgDelta"`
	FinalProtagonistAverageDelta     float64 `csv:"finPAvgDelta"`
	FinalAntagonistBestDelta         float64 `csv:"finABestDelta"`
	FinalProtagonistBestDelta        float64 `csv:"finPBestDelta"`
	FinalAntagonistEquation          string  `csv:"finAEquation"`
	FinalProtagonistEquation         string  `csv:"finPEquation"`
	FinalAntagonistStrategy          string  `csv:"finAStrat"`
	FinalProtagonistStrategy         string  `csv:"finPStrat"`
	FinalAntagonistDominantStrategy  string  `csv:"finADomStrat"`
	FinalProtagonistDominantStrategy string  `csv:"finPDomStrat"`
	FinalAntagonistBirthGen          int     `csv:"finABirthGen"`
	FinalProtagonistBirthGen         int     `csv:"finPBirthGen"`
	FinalAntagonistAge               int     `csv:"finAAge"`
	FinalProtagonistAge              int     `csv:"finPAge"`

	Run int `csv:"run"`

	// PARAMS
	GenerationCount int `csv:"genCount"`
	EachPopulationSize int `csv:"popCount"`
	AntStratCount int `csv:"antStratCount"`
	ProStratCount int `csv:"proStratCount"`
	AntStrat string `csv:"antStrat"`
	ProStrat string `csv:"proStrat"`
	RandTreeDepth int `csv:"randTreeDepth"`
	AntThreshMult float64 `csv:"antThreshMult"`
	ProThresMult float64 `csv:"proThresMult"`
	CrossPercent float64 `csv:"crossPercent"`
	ProbMutation float64 `csv:"probMutation"`
	ParentSelect string `csv:"parentSelect"`
	TournamentSize int `csv:"tournamentSize"`
	SurvivorSelect string `csv:"survivorSelect"`
	SurvivorPercent float64 `csv:"survivorPercent"`
	DivByZero string `csv:"d0"`
	DivByZeroPen float64 `csv:"d0Pen"`
}

// ReadFile will read in a .csv file. The baseFolder argument must be the path to the id folder e.g.
// ~/home/masters-go/_dataBackup/1222
func ReadCSVFile(baseFolder string) ([]CSVBestAll, error) {
	accCSV := make([]CSVBestAll, 0)

	// REQUIRE
	if baseFolder == "" {
		return nil, fmt.Errorf("baseFolder cannot be empty")
	}

	// Get all dirs
	totalDirsCount := 0
	err := filepath.Walk(baseFolder, func(path2 string, info os.FileInfo, err error) error {
		if info.IsDir() &&  path2 != baseFolder {
			totalDirsCount++
			// DO
			id := ""
			bestAllPath := ""
			err = filepath.Walk(path2, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					if strings.Contains(path, "best-") {
						//check for run
						split := strings.Split(path, "best-")
						csvNumber := split[1]
						numString := strings.ReplaceAll(csvNumber, ".csv", "")
						_, err := strconv.ParseInt(numString, 10, 64)
						if err != nil {
							return nil
						}

						fmt.Printf("%v", split)
						bestAllPath = path


						splitStrings := strings.Split(bestAllPath, string(filepath.Separator))
						id = splitStrings[len(splitStrings)-3]

						bestAllFile, err := os.OpenFile(bestAllPath, os.O_RDONLY, os.ModePerm)
						if err != nil {
							return err
						}
						defer bestAllFile.Close()

						bestIndividualStatistics := []*simulation.RunBestIndividualStatistic{}
						err = gocsv.Unmarshal(bestAllFile, &bestIndividualStatistics)
						if err != nil {
							return err
						}

						// ######################################## GET PARAMS ############################################
						paramsJsonPath := ""
						err = filepath.Walk(path2, func(path string, info os.FileInfo, err error) error {
							if !info.IsDir() {
								if strings.Contains(path, "_params.json") {
									paramsJsonPath = path
									return err
								}
							}
							return err
						})
						if err != nil {
							return err
						}
						paramsJsonFile, err := os.OpenFile(paramsJsonPath, os.O_RDONLY, os.ModePerm)
						var params evolution.EvolutionParams
						err = json.NewDecoder(paramsJsonFile).Decode(&params)
						if err != nil {
							return  err
						}

						bst := *bestIndividualStatistics[0]

						// #### COMBINE
						csvBest := CSVBestAll{
							FileID:                  id,

							SpecEquation: bst.SpecEquation,
							SpecRange: params.SpecParam.Range,
							SpecSeed: params.SpecParam.Seed,

							// PARAMS
							GenerationCount: params.GenerationsCount,
							EachPopulationSize: params.EachPopulationSize,
							ParentSelect: params.Selection.Parent.Type,

							SurvivorSelect: params.Selection.Survivor.Type,

							CrossPercent: params.Reproduction.CrossoverPercentage,
							ProbMutation: params.Reproduction.ProbabilityOfMutation,

							AntStratCount: params.Strategies.AntagonistStrategyCount,
							AntStrat:  evolution.StrategiesToStringArr(evolution.ConvertStrategiesToString(params.Strategies.
								AntagonistAvailableStrategies)),
							AntThreshMult: params.FitnessStrategy.AntagonistThresholdMultiplier,

							ProThresMult:params.FitnessStrategy.ProtagonistThresholdMultiplier,
							ProStratCount: params.Strategies.ProtagonistStrategyCount,
							ProStrat: evolution.StrategiesToStringArr(evolution.ConvertStrategiesToString(params.Strategies.
								ProtagonistAvailableStrategies)),

							RandTreeDepth:params.Strategies.DepthOfRandomNewTrees,
							DivByZero: params.SpecParam.DivideByZeroStrategy,
							DivByZeroPen: params.SpecParam.DivideByZeroPenalty,


							// INDIVIDUAL
							Antagonist: bst.Antagonist,
							Protagonist: bst.Protagonist,
							AntagonistAge: bst.AntagonistAge,
							AntagonistAverageDelta:bst.AntagonistAverageDelta,
							AntagonistBestDelta:bst.AntagonistBestDelta,
							AntagonistBestFitness:bst.AntagonistBestFitness,
							AntagonistBirthGen:bst.AntagonistBirthGen,
							AntagonistDominantStrategy:bst.AntagonistDominantStrategy,
							AntagonistEquation:bst.AntagonistEquation,
							AntagonistGeneration:bst.AntagonistGeneration,
							AntagonistID:bst.AntagonistID,
							AntagonistStdDev:bst.AntagonistStdDev,
							AntagonistStrategy:bst.AntagonistStrategy,

							ProtagonistAge: bst.ProtagonistAge,
							ProtagonistAverageDelta:bst.ProtagonistAverageDelta,
							ProtagonistBestDelta:bst.ProtagonistBestDelta,
							ProtagonistBestFitness:bst.ProtagonistBestFitness,
							ProtagonistBirthGen:bst.ProtagonistBirthGen,
							ProtagonistDominantStrategy:bst.ProtagonistDominantStrategy,
							ProtagonistEquation:bst.ProtagonistEquation,
							ProtagonistGeneration:bst.ProtagonistGeneration,
							ProtagonistID:bst.ProtagonistID,
							ProtagonistStdDev:bst.ProtagonistStdDev,
							ProtagonistStrategy:bst.ProtagonistStrategy,

							FinalAntagonist:bst.FinalAntagonist,
							FinalAntagonistAge: bst.FinalAntagonistAge,
							FinalAntagonistAverageDelta:bst.FinalAntagonistAverageDelta,
							FinalAntagonistBestDelta:bst.FinalAntagonistBestDelta,
							FinalAntagonistBestFitness:bst.FinalAntagonistBestFitness,
							FinalAntagonistBirthGen:bst.FinalAntagonistBirthGen,
							FinalAntagonistDominantStrategy:bst.FinalAntagonistDominantStrategy,
							FinalAntagonistEquation:bst.FinalAntagonistEquation,
							FinalAntagonistStdDev:bst.FinalAntagonistStdDev,
							FinalAntagonistStrategy:bst.FinalAntagonistStrategy,


							FinalProtagonist:bst.FinalProtagonist,
							FinalProtagonistAge: bst.FinalProtagonistAge,
							FinalProtagonistAverageDelta:bst.FinalProtagonistAverageDelta,
							FinalProtagonistBestDelta:bst.FinalProtagonistBestDelta,
							FinalProtagonistBestFitness:bst.FinalProtagonistBestFitness,
							FinalProtagonistBirthGen:bst.FinalProtagonistBirthGen,
							FinalProtagonistDominantStrategy:bst.FinalProtagonistDominantStrategy,
							FinalProtagonistEquation:bst.FinalProtagonistEquation,
							FinalProtagonistStdDev:bst.FinalProtagonistStdDev,
							FinalProtagonistStrategy:bst.FinalProtagonistStrategy,
						}

						mut := sync.Mutex{}
						mut.Lock()
						(accCSV) = append(accCSV, csvBest)
						mut.Unlock()

						return err
					}
				}
				return err
			})
			if err != nil {
				return  err
			}

					}
		return err
	})
	if err != nil {
		return nil, err
	}

	finalCSV := make([]CSVBestAll, 0)
	for i := range accCSV {
		if i % 2 == 0 {
			finalCSV = append(finalCSV, accCSV[i])
		}
	}

	fmt.Printf("TOTAL DIRS: %d", totalDirsCount/2)
	return finalCSV, nil
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
