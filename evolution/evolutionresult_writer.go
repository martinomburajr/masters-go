package evolution

//// WriteToFile will output the results of an evolutionResult to a specified filepath.
//func (evolutionResult *EvolutionResult) WriteToFile(path string, params EvolutionParams) (string, error) {
//	longestStrategy := params.Strategies.AntagonistStrategyCount
//	if longestStrategy < params.Strategies.ProtagonistStrategyCount {
//		longestStrategy = params.Strategies.ProtagonistStrategyCount
//	}
//
//	csvOutput := CSVOutput{
//		Generational: make([]GenerationalStatistics, len(evolutionResult.ThoroughlySortedGenerations)),
//		Epochal: make([]EpochalStatistics, len(evolutionResult.ThoroughlySortedGenerations[0].Protagonists[0].
//			Fitness)),
//		Strategy: make([]RunStrategyStatistics, longestStrategy),
//	}
//
//	coevolutionaryAverages := evolutionResult.CoevolutionaryAverages
//
//	// GENERATIONAL
//	for i := range coevolutionaryAverages {
//		csvOutput.Generational[i].Generation = i + 1
//		csvOutput.Generational[i].Run = params.InternalCount + 1
//		csvOutput.Generational[i].Spec = params.SpecParam.Expression
//
//		// ########################################## ANTAGONISTS ###################################################
//		topAntagonist := evolutionResult.ThoroughlySortedGenerations[i].Antagonists[0]
//		topAntagonistEquation, _ := topAntagonist.Program.T.ToMathematicalString()
//
//		csvOutput.Generational[i].AverageAntagonist = coevolutionaryAverages[i].AntagonistFitnessAverages
//		csvOutput.Generational[i].TopAntagonistInRun = topAntagonist.AverageFitness
//		csvOutput.Generational[i].TopAntagonistBirthGen = topAntagonist.BirthGen
//		csvOutput.Generational[i].TopAntagonistDelta = topAntagonist.BestDelta
//		csvOutput.Generational[i].TopAntagonistEquation = topAntagonistEquation
//		csvOutput.Generational[i].TopAntagonistFavoriteStrategy = DominantStrategy(*topAntagonist)
//		csvOutput.Generational[i].TopAntagonistStrategies = StrategiesToString(*topAntagonist)
//		csvOutput.Generational[i].TopAntagonistSD = topAntagonist.FitnessStdDev
//		csvOutput.Generational[i].TopAntagonistBestFitness = topAntagonist.BestFitness
//
//		// ########################################## PROTAGONISTS ###################################################
//		topProtagonist := evolutionResult.ThoroughlySortedGenerations[i].Protagonists[0]
//		topProtagonistEquation, _ := topProtagonist.Program.T.ToMathematicalString()
//
//		csvOutput.Generational[i].AverageProtagonist = coevolutionaryAverages[i].ProtagonistFitnessAverages
//		csvOutput.Generational[i].TopProtagonistInRun = topProtagonist.AverageFitness
//		csvOutput.Generational[i].TopProtagonistBirthGen = topProtagonist.BirthGen
//		csvOutput.Generational[i].TopProtagonistDelta = topProtagonist.BestDelta
//		csvOutput.Generational[i].TopProtagonistEquation = topProtagonistEquation
//		csvOutput.Generational[i].TopProtagonistFavoriteStrategy = DominantStrategy(*topProtagonist)
//		csvOutput.Generational[i].TopProtagonistStrategies = StrategiesToString(*topProtagonist)
//		csvOutput.Generational[i].TopAntagonistSD = topProtagonist.FitnessStdDev
//		csvOutput.Generational[i].TopAntagonistBestFitness = topProtagonist.BestFitness
//	}
//
//	topProtagonist := evolutionResult.ThoroughlySortedGenerations[0].Protagonists[0]
//	topProtagonistEq, _ := topProtagonist.Program.T.ToMathematicalString()
//	topAntagonist := evolutionResult.ThoroughlySortedGenerations[0].Antagonists[0]
//	topAntagonistEq, _ := topAntagonist.Program.T.ToMathematicalString()
//	finalProtagonist := evolutionResult.FinalProtagonist
//	finalProtagonistEq, _ := finalProtagonist.Program.T.ToMathematicalString()
//	finalAntagonist := evolutionResult.FinalAntagonist
//	finalAntagonistEq, _ := finalAntagonist.Program.T.ToMathematicalString()
//
//	// Epochal
//	for i := 0; i < len(csvOutput.Epochal); i++ {
//		csvOutput.Epochal[i].Epoch = i + 1
//
//		csvOutput.Epochal[i].TopAntagonistInRun = topAntagonist.Fitness[i]
//		csvOutput.Epochal[i].TopAntagonistBirthGen = topAntagonist.BirthGen
//		csvOutput.Epochal[i].TopAntagonistDelta = topAntagonist.BestDelta
//		csvOutput.Epochal[i].TopAntagonistEquation = topAntagonistEq
//		csvOutput.Epochal[i].TopAntagonistStrategy = StrategiesToString(*topAntagonist)
//		csvOutput.Epochal[i].TopAntagonistDominantStrategy = DominantStrategy(*topAntagonist)
//
//		csvOutput.Epochal[i].TopProtagonistInRun = topProtagonist.Fitness[i]
//		csvOutput.Epochal[i].TopProtagonistBirthGen = topProtagonist.BirthGen
//		csvOutput.Epochal[i].TopProtagonistDelta = topProtagonist.BestDelta
//		csvOutput.Epochal[i].TopProtagonistEquation = topProtagonistEq
//		csvOutput.Epochal[i].TopProtagonistStrategy = StrategiesToString(*topProtagonist)
//		csvOutput.Epochal[i].TopProtagonistDominantStrategy = DominantStrategy(*topProtagonist)
//
//		csvOutput.Epochal[i].FinalAntagonist = finalAntagonist.Fitness[i]
//		csvOutput.Epochal[i].FinalAntagonistBirthGen = finalAntagonist.BirthGen
//		csvOutput.Epochal[i].FinalAntagonistBestDelta = finalAntagonist.BestDelta
//		csvOutput.Epochal[i].FinalAntagonistEquation = finalAntagonistEq
//		csvOutput.Epochal[i].FinalAntagonistStrategy = StrategiesToString(*finalAntagonist)
//		csvOutput.Epochal[i].FinalAntagonistDominantStrategy = DominantStrategy(*finalAntagonist)
//
//		csvOutput.Epochal[i].FinalProtagonist = finalProtagonist.Fitness[i]
//		csvOutput.Epochal[i].FinalProtagonistBirthGen = finalProtagonist.BirthGen
//		csvOutput.Epochal[i].FinalProtagonistBestDelta = finalProtagonist.BestDelta
//		csvOutput.Epochal[i].FinalProtagonistEquation = finalProtagonistEq
//		csvOutput.Epochal[i].FinalProtagonistStrategy = StrategiesToString(*finalProtagonist)
//		csvOutput.Epochal[i].FinalProtagonistDominantStrategy = DominantStrategy(*finalProtagonist)
//	}
//
//	// Strategy
//
//	// Internal Variance of Ultimate Individuals
//	err := os.Mkdir(params.StatisticsOutput.OutputDir, 0755)
//	innerFolder := strings.ReplaceAll(path, ".json", "")
//	err = os.Mkdir(innerFolder, 0755)
//	g := strings.SplitAfter(path, "/")
//
//	mainDir := g[0]
//	subDirInfo := g[1]
//	subsubDirName := strings.ReplaceAll(g[2], ".json", "")
//
//	csvMap := map[string]interface{}{
//		"generational": csvOutput.Generational,
//		"epochal":      csvOutput.Epochal,
//	}
//	err = WriteCSVWithMap(csvMap, mainDir, subDirInfo, subsubDirName, params.InternalCount)
//	if err != nil {
//		return path, err
//	}
//	return path, nil
//}
//

//
//func WriteCSVWithMap(csvFileMap map[string]interface{}, mainDir, subDirInfo, subsubDirName string, count int) (err error) {
//	for name := range csvFileMap {
//		pathCSV := fmt.Sprintf("%s%s%s/%s-%d%s", mainDir, subDirInfo, subsubDirName, name, count, ".csv")
//		fileCSV, err := os.Create(pathCSV)
//		if err != nil {
//			return err
//		}
//		defer fileCSV.Close()
//
//		writer := gocsv.DefaultCSVWriter(fileCSV)
//		if writer.Error() != nil {
//			return writer.Error()
//		}
//		err = gocsv.Marshal(csvFileMap[name], fileCSV)
//		if err != nil {
//			return err
//		}
//	}
//	return err
//}
//
//type CSVOutput struct {
//	Generational []GenerationalStatistics `csv:"generational"`
//	Epochal      []EpochalStatistics      `csv:"epochal"`
//	Strategy     []RunStrategyStatistics  `csv:"strategyStatistics"`
//	//Run          RunBasedStatistics     `csv:"runBased"`
//}
//
//type MultiOutput struct {
//	CSVOutputs []CSVOutput
//}
//
//// Load reads all the files from the specified dir and populates the MultiOutput struct.
//// It will overwrite existing data. kind can either be "generational" or "epochal"
//func (c *MultiOutput) Load(outputDir, filePath string, evolutionParams EvolutionParams) error {
//	files := make([]string, 0)
//	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
//		files = append(files, path)
//		return nil
//	})
//	if err != nil {
//		return err
//	}
//
//	generationalFiles := make([]string, 0)
//	epochalFiles := make([]string, 0)
//	for i := range files {
//		if strings.Contains(files[i], "generational") {
//			generationalFiles = append(generationalFiles, files[i])
//		}
//		if strings.Contains(files[i], "epochal") {
//			epochalFiles = append(epochalFiles, files[i])
//		}
//	}
//	if generationalFiles == nil {
//		return fmt.Errorf("CoalesceFiles | could not coalesce files - nil")
//	}
//	if len(generationalFiles) == 0 {
//		return fmt.Errorf("CoalesceFiles | no files to coalesce")
//	}
//
//	if epochalFiles == nil {
//		return fmt.Errorf("CoalesceFiles | epochalFiles could not coalesce files - nil")
//	}
//	if len(epochalFiles) == 0 {
//		return fmt.Errorf("CoalesceFiles |epochalFiles no files to coalesce")
//	}
//
//	maxLen := len(generationalFiles)
//	if maxLen < len(epochalFiles) {
//		maxLen = len(epochalFiles)
//	}
//	c.CSVOutputs = make([]CSVOutput, maxLen)
//
//	for i := 0; i < len(generationalFiles); i++ {
//		filePath := fmt.Sprintf("%s", generationalFiles[i])
//		openFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
//		if err != nil {
//			return err
//		}
//		defer openFile.Close()
//
//		var generationalStatistic []GenerationalStatistics
//		err = gocsv.UnmarshalFile(openFile, &generationalStatistic)
//		if err != nil {
//			return err
//		}
//
//		c.CSVOutputs[i].Generational = append(c.CSVOutputs[i].Generational, generationalStatistic...)
//	}
//
//	for i := 0; i < len(epochalFiles); i++ {
//		filePath := fmt.Sprintf("%s", epochalFiles[i])
//		openFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
//		if err != nil {
//			return err
//		}
//		defer openFile.Close()
//
//		var epochalStatistics []EpochalStatistics
//		err = gocsv.UnmarshalFile(openFile, &epochalStatistics)
//		if err != nil {
//			return err
//		}
//
//		c.CSVOutputs[i].Epochal = append(c.CSVOutputs[i].Epochal, epochalStatistics...)
//	}
//
//	err = os.Mkdir(outputDir, 0755)
//	return err
//}
//
//// WriteAverages
//func (c *MultiOutput) WriteAverages(evolutionParams EvolutionParams, outputPath string) error {
//	if c.CSVOutputs == nil {
//		return fmt.Errorf("cannot write as mutlicsvoutput is nil")
//	}
//	if len(c.CSVOutputs) < 1 {
//		return fmt.Errorf("cannot write as mutlicsvoutput is empty")
//	}
//
//	// averages
//	type AveragedGenerationalStatistics struct {
//		AverageAntagonist      float64 `csv:"avgA"`
//		AverageProtagonist     float64 `csv:"avgP"`
//		TopAntagonistInRun          float64 `csv:"topA"`
//		TopProtagonistInRun         float64 `csv:"topP"`
//		TopAntagonistEquation  string  `csv:"topAEquation"`
//		TopProtagonistEquation string  `csv:"topPEquation"`
//	}
//
//	listLength := len(c.CSVOutputs[0].Generational)
//	type AveragedStatistics struct {
//		AveragedGenerationalStatistics []AveragedGenerationalStatistics `csv:"averagedGenerational"`
//	}
//
//	coalesced := AveragedStatistics{
//		AveragedGenerationalStatistics: make([]AveragedGenerationalStatistics, listLength),
//	}
//
//	for i := 0; i < len(c.CSVOutputs[0].Generational); i++ {
//		sumAverageProtagonists := 0.0
//		sumAverageAntagonists := 0.0
//		sumTopAntagonist := 0.0
//		sumTopProtagonist := 0.0
//		for _, csvFile := range c.CSVOutputs {
//			sumAverageAntagonists += csvFile.Generational[i].AverageAntagonist
//			sumAverageProtagonists += csvFile.Generational[i].AverageProtagonist
//			sumTopAntagonist += csvFile.Generational[i].TopAntagonistInRun
//			sumTopProtagonist += csvFile.Generational[i].TopProtagonistInRun
//		}
//		coalesced.AveragedGenerationalStatistics[i].AverageAntagonist = sumAverageAntagonists / float64(len(
//			c.CSVOutputs))
//		coalesced.AveragedGenerationalStatistics[i].AverageProtagonist = sumAverageProtagonists / float64(len(c.CSVOutputs))
//		coalesced.AveragedGenerationalStatistics[i].TopAntagonistInRun = sumTopAntagonist / float64(len(c.CSVOutputs))
//		coalesced.AveragedGenerationalStatistics[i].TopProtagonistInRun = sumTopProtagonist / float64(len(c.CSVOutputs))
//	}
//
//	// BEST EQUATIONS
//	bestEquations, err := BestEquationPerGeneration(c, evolutionParams.Spec)
//	if err != nil {
//		return err
//	}
//	for i, bestEquation := range bestEquations {
//		coalesced.AveragedGenerationalStatistics[i].TopProtagonistEquation = bestEquation.ProtagonistEquation
//		coalesced.AveragedGenerationalStatistics[i].TopAntagonistEquation = bestEquation.AntagonistEquation
//	}
//
//	// WRITE TO FILE
//	outputFileCSV, err := os.Create(outputPath)
//	if err != nil {
//		return err
//	}
//	defer outputFileCSV.Close()
//
//	writer := gocsv.DefaultCSVWriter(outputFileCSV)
//	if writer.Error() != nil {
//		return writer.Error()
//	}
//	err = gocsv.Marshal(coalesced.AveragedGenerationalStatistics, outputFileCSV)
//	if err != nil {
//		return err
//	}
//	fmt.Printf("\nWrote Averages to file: %s", outputPath)
//	return err
//}
//
//// WriteBestIndividuals only returns a single row of the best individuals throughout all the runs ever. This is the creme dela creme.
//func (c *MultiOutput) WriteBestIndividuals(evolutionParams EvolutionParams, outputDir, outputPath string) error {
//	if c.CSVOutputs == nil {
//		return fmt.Errorf("cannot write as mutlicsvoutput is nil")
//	}
//	if len(c.CSVOutputs) < 1 {
//		return fmt.Errorf("cannot write as mutlicsvoutput is empty")
//	}
//
//	type TotalGenerationalStatistics struct {
//		SpecEquation                string  `csv:"specEquation"`
//		SpecRange                   int     `csv:"range"`
//		SpecSeed                    int     `csv:"seed"`
//		AntagonistEquation          string  `csv:"A"`
//		ProtagonistEquation         string  `csv:"P"`
//		AntagonistBestDelta             float64 `csv:"ADelta"`
//		ProtagonistBestDelta            float64 `csv:"PDelta"`
//		AntagonistGeneration        int     `csv:"AGeneration"`
//		ProtagonistGeneration       int     `csv:"PGeneration"`
//		AntagonistRun               int     `csv:"ARun"`
//		ProtagonistRun              int     `csv:"PRun"`
//		AntagonistBirthGen          int     `csv:"ABirthGen"`
//		ProtagonistBirthGen         int     `csv:"PBirthGen"`
//		AntagonistDominantStrategy  string  `csv:"AFaveStrategy"`
//		ProtagonistDominantStrategy string  `csv:"PFaveStrategy"`
//		AntagonistStrategyList      string  `csv:"AStrategies"`
//		ProtagonistStrategyList     string  `csv:"PStrategies"`
//	}
//
//	type TotalGenerationalStrategy struct {
//		strategy string `csv:"strategy"`
//	}
//
//	type TotalStatistics struct {
//		TotalStatistics []TotalGenerationalStatistics `csv:"averagedGenerational"`
//	}
//
//	coalesced := TotalStatistics{
//		TotalStatistics: make([]TotalGenerationalStatistics, 1),
//	}
//
//	// BEST EQUATIONS
//	bestEquation, err := BestEquationAllGenerations(c, evolutionParams.Spec)
//	if err != nil {
//		return err
//	}
//
//	coalesced.TotalStatistics[0] = TotalGenerationalStatistics{
//		SpecEquation:                evolutionParams.SpecParam.Expression,
//		AntagonistEquation:          bestEquation.AntagonistEquation,
//		ProtagonistEquation:         bestEquation.ProtagonistEquation,
//		AntagonistBestDelta:             bestEquation.AntagonistBestDelta,
//		ProtagonistBestDelta:            bestEquation.ProtagonistBestDelta,
//		AntagonistGeneration:        bestEquation.AntagonistGeneration,
//		ProtagonistGeneration:       bestEquation.ProtagonistGeneration,
//		AntagonistBirthGen:          bestEquation.AntagonistBirthGen,
//		ProtagonistBirthGen:         bestEquation.ProtagonistBirthGen,
//		AntagonistRun:               bestEquation.AntagonistRun,
//		ProtagonistRun:              bestEquation.ProtagonistRun,
//		AntagonistStrategyList:      bestEquation.AntagonistStrategy,
//		ProtagonistStrategyList:     bestEquation.ProtagonistStrategy,
//		AntagonistDominantStrategy:  DominantStrategyStr(bestEquation.AntagonistStrategy),
//		ProtagonistDominantStrategy: DominantStrategyStr(bestEquation.ProtagonistStrategy),
//		SpecRange:                   evolutionParams.SpecParam.Range,
//		SpecSeed:                    evolutionParams.SpecParam.Seed,
//	}
//
//	bestAntagonistStrategyList := coalesced.TotalStatistics[0].AntagonistStrategyList
//	bestAntagonistStrategy := strings.Split(bestAntagonistStrategyList, "|")
//
//	bestProtagonistStrategyList := coalesced.TotalStatistics[0].ProtagonistStrategyList
//	bestProtagonistStrategy := strings.Split(bestProtagonistStrategyList, "|")
//
//	stratLen := len(bestProtagonistStrategy)
//	if len(bestProtagonistStrategy) < len(bestAntagonistStrategy) {
//		stratLen = len(bestAntagonistStrategy)
//	}
//
//	bestStrategies := make([]BestStrategy, stratLen)
//	for i := 0; i < len(bestAntagonistStrategy); i++ {
//		strategy := BestStrategy{
//			AntagonistStrategy: bestAntagonistStrategy[i],
//		}
//		bestStrategies[i] = strategy
//	}
//	for i := 0; i < len(bestProtagonistStrategy); i++ {
//		bestStrategies[i].ProtagonistStrategy = bestProtagonistStrategy[i]
//	}
//
//	// WRITE TO FILE
//	outputFileCSV, err := os.Create(outputPath)
//	if err != nil {
//		return err
//	}
//	defer outputFileCSV.Close()
//
//	writer := gocsv.DefaultCSVWriter(outputFileCSV)
//	if writer.Error() != nil {
//		return writer.Error()
//	}
//	err = gocsv.Marshal(coalesced.TotalStatistics, outputFileCSV)
//	if err != nil {
//		return err
//	}
//
//	// STRATEGY FILE
//	outputStrategyPath := fmt.Sprintf("%s/%s", outputDir, "strategy.csv")
//	outputStratFileCSV, err := os.Create(outputStrategyPath)
//	if err != nil {
//		return err
//	}
//	defer outputFileCSV.Close()
//	writer2 := gocsv.DefaultCSVWriter(outputStratFileCSV)
//	if writer2.Error() != nil {
//		return writer.Error()
//	}
//	err = gocsv.Marshal(bestStrategies, outputStratFileCSV)
//	if err != nil {
//		return err
//	}
//	fmt.Printf("\nWrote Strategy to file: %s", outputStrategyPath)
//
//	return err
//}
//
//// WriteGenerationalCumulative takes all data in all files and concatenates them ontop of each other and writes to the specified
////output. This is great to compare runs side by side and havel all data collated into a single document
//func (c *MultiOutput) WriteGenerationalCumulative(evolutionParams EvolutionParams, outputPath string) error {
//	if c.CSVOutputs == nil {
//		return fmt.Errorf("cannot write as mutlicsvoutput is nil")
//	}
//	if len(c.CSVOutputs) < 1 {
//		return fmt.Errorf("cannot write as mutlicsvoutput is empty")
//	}
//	if outputPath == "" {
//		return fmt.Errorf("outputPath empty")
//	}
//
//	baseCSV := c.CSVOutputs[0]
//	for i := 1; i < len(c.CSVOutputs); i++ {
//		baseCSV.Generational = append(baseCSV.Generational, c.CSVOutputs[i].Generational...)
//	}
//
//	outputFileCSV, err := os.Create(outputPath)
//	if err != nil {
//		return err
//	}
//	defer outputFileCSV.Close()
//
//	writer := gocsv.DefaultCSVWriter(outputFileCSV)
//	if writer.Error() != nil {
//		return writer.Error()
//	}
//	err = gocsv.Marshal(baseCSV.Generational, outputFileCSV)
//	if err != nil {
//		return err
//	}
//	fmt.Printf("\nWrote Cumulative to file: %s", outputPath)
//	return err
//
//	return nil
//}
//
//type BestEquation struct {
//	AntagonistEquation   string
//	AntagonistBestDelta      float64
//	AntagonistStrategy   string
//	AntagonistGeneration int
//	AntagonistBirthGen   int
//	AntagonistRun        int
//
//	ProtagonistEquation   string
//	ProtagonistBestDelta      float64
//	ProtagonistStrategy   string
//	ProtagonistGeneration int
//	ProtagonistBirthGen   int
//	ProtagonistRun        int
//	SpecDelta             float64
//}
//
//type BestStrategy struct {
//	AntagonistStrategy  string
//	ProtagonistStrategy string
//}
//
//// TODO Return Ultimate Statistics
//func BestEquationAllGenerations(multiCSVOutput *MultiOutput, spec SpecMulti) (bestEquation BestEquation,
//	err error) {
//	if multiCSVOutput.CSVOutputs == nil {
//		return BestEquation{}, fmt.Errorf("BestEquation | json csvFiles cannot be nil")
//	}
//	if len(multiCSVOutput.CSVOutputs) < 1 {
//		return BestEquation{}, fmt.Errorf("BestEquation | json csvFiles cannot be empty")
//	}
//
//	numberOfGenerations := len(multiCSVOutput.CSVOutputs[0].Generational)
//	bestEquation = BestEquation{}
//
//	bestAntagonistDelta := -math.MaxFloat64
//	bestProtagonistDelta := math.MaxFloat64
//	bestAntagonistEquation := ""
//	bestProtagonistEquation := ""
//	for i := 0; i < numberOfGenerations; i++ {
//
//		for _, csvFile := range multiCSVOutput.CSVOutputs {
//			antagonistDelta := 0.0
//			protagonistDelta := 0.0
//			antagonistEquation := csvFile.Generational[i].TopAntagonistEquation
//			protagonistEquation := csvFile.Generational[i].TopProtagonistEquation
//			for s := range spec {
//				independentX := spec[s].Independents
//				dependentVarAntagonist, err := EvaluateMathematicalExpression(antagonistEquation,
//					independentX)
//				if err != nil {
//
//				}
//				antagonistDelta += math.Abs(dependentVarAntagonist - spec[s].Dependent)
//				dependentVarProagonist, err := EvaluateMathematicalExpression(protagonistEquation,
//					independentX)
//				if err != nil {
//				}
//				protagonistDelta += math.Abs(dependentVarProagonist - spec[s].Dependent)
//			}
//			if antagonistDelta >= bestAntagonistDelta {
//				bestAntagonistDelta = antagonistDelta
//				bestAntagonistEquation = antagonistEquation
//
//				bestEquation.AntagonistBestDelta = bestAntagonistDelta
//				bestEquation.AntagonistEquation = bestAntagonistEquation
//				bestEquation.AntagonistBirthGen = csvFile.Generational[i].TopAntagonistBirthGen
//				bestEquation.AntagonistGeneration = i
//				bestEquation.AntagonistRun = csvFile.Generational[i].Run
//				bestEquation.AntagonistStrategy = csvFile.Generational[i].TopAntagonistStrategies
//			}
//			if protagonistDelta <= bestProtagonistDelta {
//				bestProtagonistDelta = protagonistDelta
//				bestProtagonistEquation = protagonistEquation
//
//				bestEquation.ProtagonistBestDelta = protagonistDelta
//				bestEquation.ProtagonistEquation = bestProtagonistEquation
//				bestEquation.ProtagonistBirthGen = csvFile.Generational[i].TopProtagonistBirthGen
//				bestEquation.ProtagonistGeneration = i
//				bestEquation.ProtagonistRun = csvFile.Generational[i].Run
//				bestEquation.ProtagonistStrategy = csvFile.Generational[i].TopProtagonistStrategies
//			}
//		}
//	}
//
//	return bestEquation, nil
//}
//
//// TODO Return Ultimate Statistics
//func BestEquationPerGeneration(multiCSVOutput *MultiOutput, spec SpecMulti) (bestEquation []BestEquation,
//	err error) {
//	if multiCSVOutput.CSVOutputs == nil {
//		return nil, fmt.Errorf("coalesce | json csvFiles cannot be nil")
//	}
//	if len(multiCSVOutput.CSVOutputs) < 1 {
//		return nil, fmt.Errorf("coalesce | json csvFiles cannot be empty")
//	}
//
//	numberOfGenerations := len(multiCSVOutput.CSVOutputs[0].Generational)
//	bestEquation = make([]BestEquation, numberOfGenerations)
//
//	bestAntagonistDelta := -math.MaxFloat64
//	bestProtagonistDelta := math.MaxFloat64
//	bestAntagonistEquation := ""
//	bestProtagonistEquation := ""
//	for i := 0; i < numberOfGenerations; i++ {
//		for _, csvFile := range multiCSVOutput.CSVOutputs {
//			antagonistDelta := 0.0
//			protagonistDelta := 0.0
//			antagonistEquation := csvFile.Generational[i].TopAntagonistEquation
//			protagonistEquation := csvFile.Generational[i].TopProtagonistEquation
//			for s := range spec {
//				independentX := spec[s].Independents
//				dependentVarAntagonist, err := EvaluateMathematicalExpression(antagonistEquation,
//					independentX)
//				if err != nil {
//					// Handle Divide By Zero
//					//return nil, err
//					//
//				}
//				antagonistDelta += math.Abs(dependentVarAntagonist - spec[s].Dependent)
//				dependentVarProagonist, err := EvaluateMathematicalExpression(protagonistEquation,
//					independentX)
//				if err != nil {
//					// Handle Divide By Zero
//					//return nil, err
//					//
//				}
//				protagonistDelta += math.Abs(dependentVarProagonist - spec[s].Dependent)
//			}
//
//			if antagonistDelta > bestAntagonistDelta {
//				bestAntagonistDelta = antagonistDelta
//				bestAntagonistEquation = antagonistEquation
//			}
//			if protagonistDelta < bestProtagonistDelta {
//				bestProtagonistDelta = protagonistDelta
//				bestProtagonistEquation = protagonistEquation
//			}
//		}
//		bestEquation[i] = BestEquation{
//			AntagonistEquation:  bestAntagonistEquation,
//			ProtagonistEquation: bestProtagonistEquation,
//		}
//	}
//
//	return bestEquation, nil
//}
