package simulation

import (
	"encoding/json"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/martinomburajr/masters-go/evolution"
	"log"
	"math"
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

func (s *Simulation) Begin(params evolution.EvolutionParams) (evolution.EvolutionParams, error) {
	os.Mkdir("data", 0755)

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

			StartEngine(engine)
		//}(params, i)
	}
	//wg.Wait()
	log.Println("SYNCHRONIZED!")

	s.OutputDir = folder
	return params, nil
}

// CoalesceFiles will coalesce all files into a single document that can be analyzed
func (s *Simulation) CoalesceFiles(evolutionParams evolution.EvolutionParams) (string, error) {
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
		if strings.Contains(files[i], "generational") {
			newFiles = append(newFiles, files[i])
		}
	}
	if newFiles == nil {
		return "", fmt.Errorf("CoalesceFiles | could not coalesce files - nil")
	}
	if len(newFiles) == 0 {
		return "", fmt.Errorf("CoalesceFiles | no files to coalesce")
	}

	csvOutputs := make([]evolution.CSVOutput, 0)
	generationalStatistics := make([][]evolution.GenerationalStatistics, 0)

	for i := 0; i < len(newFiles); i++ {
		filePath := fmt.Sprintf("%s", newFiles[i])
		split := strings.Split(filePath, "/")
		topLevelDir := split[0]
		subInfoDir := split[1]
		subSubNameDir := split[2]

		absolutePath, err := filepath.Abs(filePath)
		if err != nil {
			return "", err
		}

		s.RunRScript(absolutePath, filePath, topLevelDir, subInfoDir, subSubNameDir,
			evolutionParams.InternalCount)

		openFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			return "", err
		}
		defer openFile.Close()

		var generationalStatistic []evolution.GenerationalStatistics
		err = gocsv.UnmarshalFile(openFile, &generationalStatistic)
		if err != nil {
			return "", err
		}

		generationalStatistics = append(generationalStatistics, generationalStatistic)

		csvOutput := evolution.CSVOutput{
			Generational: generationalStatistic,
		}
		csvOutputs = append(csvOutputs, csvOutput)
	}

	// write params to file
	paramsPath := fmt.Sprintf("%s/%s", s.OutputDir, "params.json")
	paramsFile, err := os.Create(paramsPath)
	if err != nil {
		return "", err
	}
	err = json.NewEncoder(paramsFile).Encode(evolutionParams)
	if err != nil {
		return "", err
	}
	fmt.Println("Wrote Params file to: " + paramsPath)
	// do epochal
	cumulative, err := coalesceCumulative(csvOutputs, evolutionParams, s.OutputDir, "cumulative-generational.csv")
	if err != nil {
		return cumulative, err
	}
	averages, err := coalesceAverages(csvOutputs, evolutionParams, s.OutputDir, "coalesced-generational.csv")
	if err != nil {
		return averages, err
	}
	best, err := coalesceBest(csvOutputs, evolutionParams, s.OutputDir, "coalesced-best.csv")
	if err != nil {
		return best, err
	}
	// do epochal
	return "", err
}

func coalesceCumulative(csvFiles []evolution.CSVOutput, evolutionParams evolution.EvolutionParams, outputDir,
	cumulativeOutputFilePath string) (string, error) {
	if csvFiles == nil {
		return "", fmt.Errorf("coalesce | json csvFiles cannot be nil")
	}
	if len(csvFiles) < 1 {
		return "", fmt.Errorf("coalesce | json csvFiles cannot be empty")
	}
	if outputDir == "" {
		return "", fmt.Errorf("outputDir empty")
	}

	baseCSV := csvFiles[0]
	for i := 1; i < len(csvFiles); i++ {
		baseCSV.Generational = append(baseCSV.Generational, csvFiles[i].Generational...)
	}

	// cumulative
	path := fmt.Sprintf("%s%s", outputDir, cumulativeOutputFilePath)
	err := os.Mkdir(outputDir, 0755)

	outputFileCSV, err := os.Create(path)
	if err != nil {
		return path, err
	}
	defer outputFileCSV.Close()

	writer := gocsv.DefaultCSVWriter(outputFileCSV)
	if writer.Error() != nil {
		return path, writer.Error()
	}
	err = gocsv.Marshal(baseCSV.Generational, outputFileCSV)
	if err != nil {
		return path, err
	}
	fmt.Printf("\nWrote Cumulative to file: %s", path)
	return path, err
}

func coalesceAverages(csvFiles []evolution.CSVOutput, evolutionParams evolution.EvolutionParams, outputDir, averagesOutputFilepath string) (
	string, error) {

	// averages
	type AveragedGenerationalStatistics struct {
		AverageAntagonist      float64 `csv:"averageAntagonist"`
		AverageProtagonist     float64 `csv:"averageProtagonist"`
		TopAntagonist          float64 `csv:"topAntagonist"`
		TopProtagonist         float64 `csv:"topProtagonist"`
		TopAntagonistEquation  string  `csv:"topAntagonistEquation"`
		TopProtagonistEquation string  `csv:"topProtagonistEquation"`
	}

	listLength := len(csvFiles[0].Generational)
	type AveragedStatistics struct {
		AveragedGenerationalStatistics []AveragedGenerationalStatistics `csv:"averagedGenerational"`
	}

	coalesced := AveragedStatistics{
		AveragedGenerationalStatistics: make([]AveragedGenerationalStatistics, listLength),
	}

	for i := 0; i < len(csvFiles[0].Generational); i++ {
		sumAverageProtagonists := 0.0
		sumAverageAntagonists := 0.0
		sumTopAntagonist := 0.0
		sumTopProtagonist := 0.0
		for _, csvFile := range csvFiles {
			sumAverageAntagonists += csvFile.Generational[i].AverageAntagonist
			sumAverageProtagonists += csvFile.Generational[i].AverageProtagonist
			sumTopAntagonist += csvFile.Generational[i].TopAntagonist
			sumTopProtagonist += csvFile.Generational[i].TopProtagonist
		}
		coalesced.AveragedGenerationalStatistics[i].AverageAntagonist = sumAverageAntagonists / float64(len(
			csvFiles))
		coalesced.AveragedGenerationalStatistics[i].AverageProtagonist = sumAverageProtagonists / float64(len(csvFiles))
		coalesced.AveragedGenerationalStatistics[i].TopAntagonist = sumTopAntagonist / float64(len(csvFiles))
		coalesced.AveragedGenerationalStatistics[i].TopProtagonist = sumTopProtagonist / float64(len(csvFiles))
	}

	// BEST EQUATIONS
	bestEquations, err := BestEquationPerGeneration(csvFiles, evolutionParams.Spec)
	if err != nil {
		return "", err
	}
	for i, bestEquation := range bestEquations {
		coalesced.AveragedGenerationalStatistics[i].TopProtagonistEquation = bestEquation.ProtagonistEquation
		coalesced.AveragedGenerationalStatistics[i].TopAntagonistEquation = bestEquation.AntagonistEquation
	}

	// WRITE TO FILE
	path := fmt.Sprintf("%s%s", outputDir, averagesOutputFilepath)
	err = os.Mkdir(outputDir, 0755)
	outputFileCSV, err := os.Create(path)
	if err != nil {
		return path, err
	}
	defer outputFileCSV.Close()

	writer := gocsv.DefaultCSVWriter(outputFileCSV)
	if writer.Error() != nil {
		return path, writer.Error()
	}
	err = gocsv.Marshal(coalesced.AveragedGenerationalStatistics, outputFileCSV)
	if err != nil {
		return path, err
	}
	fmt.Printf("\nWrote Averages to file: %s", path)

	return path, nil
}

func coalesceBest(csvFiles []evolution.CSVOutput, evolutionParams evolution.EvolutionParams, outputDir,
	averagesOutputFilepath string) (
	string, error) {

	type TotalGenerationalStatistics struct {
		SpecEquation                string  `csv:"specEquation"`
		SpecRange                   int     `csv:"range"`
		SpecSeed                    int     `csv:"seed"`
		AntagonistEquation          string  `csv:"A"`
		ProtagonistEquation         string  `csv:"P"`
		AntagonistDelta             float64 `csv:"ADelta"`
		ProtagonistDelta            float64 `csv:"PDelta"`
		AntagonistGeneration        int     `csv:"AGeneration"`
		ProtagonistGeneration       int     `csv:"PGeneration"`
		AntagonistRun               int     `csv:"ARun"`
		ProtagonistRun              int     `csv:"PRun"`
		AntagonistBirthGen          int     `csv:"ABirthGen"`
		ProtagonistBirthGen         int     `csv:"PBirthGen"`
		AntagonistDominantStrategy  string  `csv:"AFaveStrategy"`
		ProtagonistDominantStrategy string  `csv:"PFaveStrategy"`
		AntagonistStrategyList      string  `csv:"AStrategies"`
		ProtagonistStrategyList     string  `csv:"PStrategies"`
	}

	type TotalStatistics struct {
		TotalStatistics []TotalGenerationalStatistics `csv:"averagedGenerational"`
	}

	coalesced := TotalStatistics{
		TotalStatistics: make([]TotalGenerationalStatistics, 1),
	}

	// BEST EQUATIONS
	bestEquation, err := BestEquationAllGenerations(csvFiles, evolutionParams.Spec)
	if err != nil {
		return "", err
	}

	coalesced.TotalStatistics[0] = TotalGenerationalStatistics{
		SpecEquation:                evolutionParams.SpecParam.Expression,
		AntagonistEquation:          bestEquation.AntagonistEquation,
		ProtagonistEquation:         bestEquation.ProtagonistEquation,
		AntagonistDelta:             bestEquation.AntagonistDelta,
		ProtagonistDelta:            bestEquation.ProtagonistDelta,
		AntagonistGeneration:        bestEquation.AntagonistGeneration,
		ProtagonistGeneration:       bestEquation.ProtagonistGeneration,
		AntagonistBirthGen:          bestEquation.AntagonistBirthGen,
		ProtagonistBirthGen:         bestEquation.ProtagonistBirthGen,
		AntagonistRun:               bestEquation.AntagonistRun,
		ProtagonistRun:              bestEquation.ProtagonistRun,
		AntagonistStrategyList:      bestEquation.AntagonistStrategy,
		ProtagonistStrategyList:     bestEquation.ProtagonistStrategy,
		AntagonistDominantStrategy:  evolution.DominantStrategyStr(bestEquation.AntagonistStrategy),
		ProtagonistDominantStrategy: evolution.DominantStrategyStr(bestEquation.ProtagonistStrategy),
		SpecRange:                   evolutionParams.SpecParam.Range,
		SpecSeed:                    evolutionParams.SpecParam.Seed,
	}

	// WRITE TO FILE
	path := fmt.Sprintf("%s%s", outputDir, averagesOutputFilepath)
	err = os.Mkdir(outputDir, 0755)
	outputFileCSV, err := os.Create(path)
	if err != nil {
		return path, err
	}
	defer outputFileCSV.Close()

	writer := gocsv.DefaultCSVWriter(outputFileCSV)
	if writer.Error() != nil {
		return path, writer.Error()
	}
	err = gocsv.Marshal(coalesced.TotalStatistics, outputFileCSV)
	if err != nil {
		return path, err
	}
	fmt.Printf("\nWrote Averages to file: %s", path)

	return path, nil
}

type BestEquation struct {
	AntagonistEquation   string
	AntagonistDelta      float64
	AntagonistStrategy   string
	AntagonistGeneration int
	AntagonistBirthGen   int
	AntagonistRun        int

	ProtagonistEquation   string
	ProtagonistDelta      float64
	ProtagonistStrategy   string
	ProtagonistGeneration int
	ProtagonistBirthGen   int
	ProtagonistRun        int

	SpecDelta float64
}

// TODO Return Ultimate Statistics
func BestEquationAllGenerations(csvFiles []evolution.CSVOutput, spec evolution.SpecMulti) (bestEquation BestEquation,
	err error) {
	if csvFiles == nil {
		return BestEquation{}, fmt.Errorf("coalesce | json csvFiles cannot be nil")
	}
	if len(csvFiles) < 1 {
		return BestEquation{}, fmt.Errorf("coalesce | json csvFiles cannot be empty")
	}

	numberOfGenerations := len(csvFiles[0].Generational)
	bestEquation = BestEquation{}

	bestAntagonistDelta := -math.MaxFloat64
	bestProtagonistDelta := math.MaxFloat64
	bestAntagonistEquation := ""
	bestProtagonistEquation := ""
	for i := 0; i < numberOfGenerations; i++ {

		for _, csvFile := range csvFiles {
			antagonistDelta := 0.0
			protagonistDelta := 0.0
			//specAntagonistDelta := 0.0
			//specProtagonistDelta := 0.0
			antagonistEquation := csvFile.Generational[i].TopAntagonistEquation
			protagonistEquation := csvFile.Generational[i].TopProtagonistEquation
			for s := range spec {
				independentX := spec[s].Independents
				dependentVarAntagonist, err := evolution.EvaluateMathematicalExpression(antagonistEquation,
					independentX, 0)
				if err != nil {
					// Handle Divide By Zero
					//return nil, err
					
				}
				antagonistDelta += math.Abs(dependentVarAntagonist - spec[s].Dependent)
				dependentVarProagonist, err := evolution.EvaluateMathematicalExpression(protagonistEquation,
					independentX, 0)
				if err != nil {
					// Handle Divide By Zero
					//return nil, err
					
				}
				protagonistDelta += math.Abs(dependentVarProagonist - spec[s].Dependent)
			}
			if antagonistDelta >= bestAntagonistDelta {
				bestAntagonistDelta = antagonistDelta
				bestAntagonistEquation = antagonistEquation

				bestEquation.AntagonistDelta = bestAntagonistDelta
				bestEquation.AntagonistEquation = bestAntagonistEquation
				bestEquation.AntagonistBirthGen = csvFile.Generational[i].TopAntagonistBirthGen
				bestEquation.AntagonistGeneration = i
				bestEquation.AntagonistRun = csvFile.Generational[i].Run
				bestEquation.AntagonistStrategy = csvFile.Generational[i].TopAntagonistStrategies
			}
			if protagonistDelta <= bestProtagonistDelta {
				bestProtagonistDelta = protagonistDelta
				bestProtagonistEquation = protagonistEquation

				bestEquation.ProtagonistDelta = protagonistDelta
				bestEquation.ProtagonistEquation = bestProtagonistEquation
				bestEquation.ProtagonistBirthGen = csvFile.Generational[i].TopProtagonistBirthGen
				bestEquation.ProtagonistGeneration = i
				bestEquation.ProtagonistRun = csvFile.Generational[i].Run
				bestEquation.ProtagonistStrategy = csvFile.Generational[i].TopProtagonistStrategies
			}
		}
	}

	return bestEquation, nil
}

// TODO Return Ultimate Statistics
func BestEquationPerGeneration(csvFiles []evolution.CSVOutput, spec evolution.SpecMulti) (bestEquation []BestEquation,
	err error) {
	if csvFiles == nil {
		return nil, fmt.Errorf("coalesce | json csvFiles cannot be nil")
	}
	if len(csvFiles) < 1 {
		return nil, fmt.Errorf("coalesce | json csvFiles cannot be empty")
	}

	numberOfGenerations := len(csvFiles[0].Generational)
	bestEquation = make([]BestEquation, numberOfGenerations)

	bestAntagonistDelta := -math.MaxFloat64
	bestProtagonistDelta := math.MaxFloat64
	bestAntagonistEquation := ""
	bestProtagonistEquation := ""
	for i := 0; i < numberOfGenerations; i++ {
		for _, csvFile := range csvFiles {
			antagonistDelta := 0.0
			protagonistDelta := 0.0
			//specAntagonistDelta := 0.0
			//specProtagonistDelta := 0.0
			antagonistEquation := csvFile.Generational[i].TopAntagonistEquation
			protagonistEquation := csvFile.Generational[i].TopProtagonistEquation
			for s := range spec {
				independentX := spec[s].Independents
				dependentVarAntagonist, err := evolution.EvaluateMathematicalExpression(antagonistEquation,
					independentX, 0)
				if err != nil {
					// Handle Divide By Zero
					//return nil, err
					//
				}
				antagonistDelta += math.Abs(dependentVarAntagonist - spec[s].Dependent)
				dependentVarProagonist, err := evolution.EvaluateMathematicalExpression(protagonistEquation,
					independentX, 0)
				if err != nil {
					// Handle Divide By Zero
					//return nil, err
					//
				}
				protagonistDelta += math.Abs(dependentVarProagonist - spec[s].Dependent)
			}

			if antagonistDelta > bestAntagonistDelta {
				bestAntagonistDelta = antagonistDelta
				bestAntagonistEquation = antagonistEquation
			}
			if protagonistDelta < bestProtagonistDelta {
				bestProtagonistDelta = protagonistDelta
				bestProtagonistEquation = protagonistEquation
			}
		}
		bestEquation[i] = BestEquation{
			AntagonistEquation:  bestAntagonistEquation,
			ProtagonistEquation: bestProtagonistEquation,
		}
	}

	return bestEquation, nil
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

func StartEngine(engine *evolution.EvolutionEngine) error {
	evolutionResult, err := engine.Start()
	if err != nil {
		return err
	}

	err = evolutionResult.Analyze(engine.Generations, true, engine.Parameters)
	if err != nil {
		return err
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
