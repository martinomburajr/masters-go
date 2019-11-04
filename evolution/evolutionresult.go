package evolution

import (
	"bufio"
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
	"strconv"
	"strings"
)

type EvolutionResult struct {
	HasBeenAnalyzed     bool
	TopAntagonist       *Individual
	TopProtagonist      *Individual
	IsMoreFitnessBetter bool
	FinalAntagonist *Individual
	FinalProtagonist *Individual

	CoevolutionaryAverages []generationalCoevolutionaryAverages

	SortedGenerationIndividuals []*Generation
	OutputFile                  string
}

type multiIndividualsPerGeneration struct {
	Generation  *Generation
	Individuals []*Individual
}

type generationalCoevolutionaryAverages struct {
	Generation        *Generation
	AntagonistResult  float64
	ProtagonistResult float64
}

func (e *EvolutionResult) Analyze(generations []*Generation, isMoreFitnessBetter bool,
	params EvolutionParams) error {
	sortedFinalAntagonists, err := SortIndividuals(generations[len(generations)-1].Antagonists, true)
	if err != nil {
		return err
	}
	sortedFinalProtagonists, err := SortIndividuals(generations[len(generations)-1].Protagonists, true)
	if err != nil {
		return err
	}

	e.FinalAntagonist = sortedFinalAntagonists[0]
	e.FinalProtagonist = sortedFinalProtagonists[0]

	// Perform all sorting functions on each generation for each kind of individual
	e.IsMoreFitnessBetter = isMoreFitnessBetter
	sortedGenerations, err := SortGenerationsThoroughly(generations, isMoreFitnessBetter)
	if err != nil {
		return err
	}
	e.SortedGenerationIndividuals = sortedGenerations

	// Calculate Top Individuals
	topAntagonist, topProtagonist, err := GetTopIndividualInAllGenerations(sortedGenerations, isMoreFitnessBetter)
	if err != nil {
		return err
	}
	e.TopAntagonist = topAntagonist
	e.TopProtagonist = topProtagonist

	// Calculate GenerationalStatistics Averages
	coevolutionaryAverages, err := GetGenerationalFitnessAverage(sortedGenerations)
	if err != nil {
		return err
	}
	e.CoevolutionaryAverages = coevolutionaryAverages
	e.HasBeenAnalyzed = true

	outputFile, err := WritetoFile(params.StatisticsOutput.OutputPath, e, params, params.InternalCount)
	e.OutputFile = outputFile
	return err
}

func (e *EvolutionResult) PrintAverageGenerationSummary() (strings.Builder, error) {
	if e.CoevolutionaryAverages == nil {
		return strings.Builder{},
			fmt.Errorf("PrintAverageGenerationSummary | cannot format as protagonist average field is nil | Run" +
				" analyze")
	}
	if e.CoevolutionaryAverages == nil {
		return strings.Builder{},
			fmt.Errorf("PrintAverageGenerationSummary | cannot format as antagonist average field is nil | Run" +
				" analyze")
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("" +
		"########## AVERAGE ANTAGONISTS VS PROTAGONISTS PER GENERATION" +
		" ############\n\n"))
	sb.WriteString("        ANT | PRO\n")
	for i := range e.CoevolutionaryAverages {
		res := e.CoevolutionaryAverages[i].AntagonistResult
		resPro := e.CoevolutionaryAverages[i].ProtagonistResult
		float := strconv.FormatFloat(res, 'g', 03, 64)
		floatPro := strconv.FormatFloat(resPro, 'g', 03, 64)

		gen := strconv.Itoa(i)
		sb.WriteString("gen" + gen + ": " + float + " | " + floatPro + "\n")
	}
	sb.WriteString("\n")
	return sb, nil
}

func (e *EvolutionResult) PrintTopIndividualSummary(kind int) (strings.Builder, error) {
	sb := strings.Builder{}
	var name string
	if kind == IndividualAntagonist {
		if e.TopAntagonist == nil {
			return strings.Builder{},
				fmt.Errorf("PrintTopIndividualSummary | cannot format as field is nil | Run analyze")
		}
		name = "ANTAGONIST"
		sb.WriteString(fmt.Sprintf("############### TOP %s IN ALL GENERATIONS"+" #######################\n", name))
		toString := e.TopAntagonist.ToString()
		sb.WriteString(toString.String())
	} else if kind == IndividualProtagonist {
		if e.TopProtagonist == nil {
			return strings.Builder{},
				fmt.Errorf("PrintTopIndividualSummary | cannot format as field is nil | Run analyze")
		}
		name = "PROTAGONIST"
		sb.WriteString(fmt.Sprintf("############### TOP %s IN ALL GENERATIONS"+" #######################\n", name))
		toString := e.TopProtagonist.ToString()
		sb.WriteString(toString.String())
	}
	return sb, nil
}

// StartInteractiveTerminal will start a Command Line Interface CLI that allows the user to interact with the presented
// data. The user has the ability to select from a given set of options using their keyboard after the user has been
// prompted. Here are a list of interactive elements that this module can print
// 0. Top Protagonist
// 1. Top Antagonist
// 2. Average GenerationalStatistics Fitness
// 3. Top Antagonist in Gen(X)
// 4. Top Protagonist in Gen(X)
// 5. Top N Antagonists in Gen(x)
// 6. Top N Protagonists in Gen(x)
// 7. Individual (X) in Generation (Y)
// 8. Output Results to File
func (e *EvolutionResult) StartInteractiveTerminal(params EvolutionParams) error {
	fmt.Println()
	fmt.Println()
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("        Welcome to the Evolutionary Interactive Module")
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("Select Your Choices Below | Type !q to exit!")

	reader := bufio.NewReader(os.Stdin)
	isRunning := true
	for isRunning {
		fmt.Println("0. Top Antagonist")
		fmt.Println("1. Top Protagonist")
		fmt.Println("2. Average GenerationalStatistics Fitness")
		fmt.Println("3. Top Antagonist in Gen(X)")
		fmt.Println("4. Top Protagonist in Gen(X)")
		fmt.Println("5. Top N Antagonists in Gen(x)")
		fmt.Println("6. Top N Protagonists in Gen(x)")
		fmt.Println("7. Individual (X) in Generation (Y)")
		fmt.Println("8. Search By Expression")
		fmt.Println("9. Output Results to File")
		fmt.Println("\nType !q to exit!")
		fmt.Print("->")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		if strings.Compare("!q", text) == 0 {
			isRunning = false
		}

		switch text {
		case "0":
			strBuilder := e.TopAntagonist.ToString()
			bannerStr := banner("Top Antagonist")
			fmt.Println(bannerStr)
			fmt.Println(strBuilder.String())
			fmt.Println()
		case "1":
			strBuilder := e.TopProtagonist.ToString()
			bannerStr := banner("Top Protagonist")
			fmt.Println(bannerStr)
			fmt.Println(strBuilder.String())
			fmt.Println()
		case "2":
			strBuilder, err := e.PrintAverageGenerationSummary()
			if err != nil {
				return err
			}
			fmt.Println(strBuilder.String())
			fmt.Println()
		case "3":
			str, err := interactiveGetTopIndividualInGenX(reader, e.SortedGenerationIndividuals, IndividualAntagonist,
				e.IsMoreFitnessBetter)
			if err != nil {
				return err
			}
			fmt.Println(str)
			fmt.Println()
		case "4":
			str, err := interactiveGetTopIndividualInGenX(reader, e.SortedGenerationIndividuals, IndividualProtagonist,
				e.IsMoreFitnessBetter)
			if err != nil {
				return err
			}
			fmt.Println(str)
			fmt.Println()
		case "5":
			str, err := interactiveGetTopNIndividualInGenX(reader, e.SortedGenerationIndividuals, IndividualAntagonist,
				e.IsMoreFitnessBetter)
			if err != nil {
				return err
			}
			fmt.Println(str)
			fmt.Println()
		case "6":
			str, err := interactiveGetTopNIndividualInGenX(reader, e.SortedGenerationIndividuals, IndividualProtagonist,
				e.IsMoreFitnessBetter)
			if err != nil {
				return err
			}
			fmt.Println(str)
			fmt.Println()
		case "7":
			str, err := interactiveGetIndividualXInGenY(reader, e.SortedGenerationIndividuals,
				e.IsMoreFitnessBetter)
			if err != nil {
				return err
			}
			fmt.Println(str)
			fmt.Println()
		case "8":
			str, err := interactiveSearchForTreeShape(reader, e.SortedGenerationIndividuals, params)
			if err != nil {
				return err
			}
			fmt.Println(str)
			fmt.Println()
		case "9":
			str, err := interactiveWriteToFile(reader, e, params)
			if err != nil {
				return err
			}
			fmt.Println(str)
			fmt.Println()
		default:
			fmt.Println("Invalid Option! To exit type !q ")
		}
	}
	return nil
}

func banner(title string) string {
	builder := strings.Builder{}
	builder.WriteString("############### ")
	builder.WriteString(strings.ToUpper(title))
	builder.WriteString("  ###############")
	return builder.String()
}

func interactiveGetTopIndividualInGenX(reader *bufio.Reader, sortedIndividuals []*Generation, kind int,
	isMoreFitnessBetter bool) (string, error) {
	var generationInt int
	isNotValidInt := true
	for isNotValidInt {
		fmt.Print(fmt.Sprintf("Input a Generation Number [0,%d)", len(sortedIndividuals)))
		fmt.Print("---->")
		generation, _ := reader.ReadString('\n')
		// convert CRLF to LF
		generation = strings.Replace(generation, "\n", "", -1)
		genInt, err := strconv.ParseInt(generation, 10, 64)
		if err != nil {
			isNotValidInt = true
			fmt.Print(fmt.Sprintf("Please enter a NUMBER between [0,%d)",
				len(sortedIndividuals)))
		} else {
			generationInt = int(genInt)
			isNotValidInt = false
		}
	}
	topIndividual, err := GetTopNIndividualInGenerationX(sortedIndividuals, kind,
		isMoreFitnessBetter, 1, int(generationInt))
	if err != nil {
		return "", err
	}
	bannerBuilder := strings.Builder{}
	switch kind {
	case IndividualAntagonist:
		name := "Antagonist"
		bannerStr := banner(fmt.Sprintf("Top %s in Gen %d", name, generationInt))
		bannerBuilder.WriteString(bannerStr)
	case IndividualProtagonist:
		name := "Protagonist"
		bannerStr := banner(fmt.Sprintf("Top %s in Gen %d", name, generationInt))
		bannerBuilder.WriteString(bannerStr)
	}

	builder := topIndividual.Individuals[0].ToString()
	bannerBuilder.WriteString(builder.String())
	return bannerBuilder.String(), nil
}

func interactiveGetTopNIndividualInGenX(reader *bufio.Reader, sortedIndividuals []*Generation, kind int,
	isMoreFitnessBetter bool) (string, error) {
	var generationInt int
	isNotValidInt := true

	for isNotValidInt {
		fmt.Print(fmt.Sprintf("Input a Generation Number [0,%d)", len(sortedIndividuals)))
		fmt.Print("---->")
		generation, _ := reader.ReadString('\n')
		// convert CRLF to LF
		generation = strings.Replace(generation, "\n", "", -1)
		genInt, err := strconv.ParseInt(generation, 10, 64)
		if err != nil {
			isNotValidInt = true
			fmt.Print(fmt.Sprintf("Please enter a NUMBER between [0,%d)",
				len(sortedIndividuals)))
		} else {
			generationInt = int(genInt)
			isNotValidInt = false
		}
	}
	topN := 3
	isNotValidTopN := true
	for isNotValidTopN {
		fmt.Print(fmt.Sprintf("Input Top(N) Individual Number [0,%d)", len(sortedIndividuals[0].Antagonists)))
		fmt.Print("---->")
		topStr, _ := reader.ReadString('\n')
		// convert CRLF to LF
		topStr = strings.Replace(topStr, "\n", "", -1)
		topStrInt, err := strconv.ParseInt(topStr, 10, 64)
		topN = int(topStrInt)
		if err != nil {
			isNotValidTopN = true
			fmt.Print(fmt.Sprintf("Please enter a NUMBER between [0,%d)",
				len(sortedIndividuals)))
		} else {
			if topN > len(sortedIndividuals[0].Antagonists) {
				topN = len(sortedIndividuals[0].Antagonists) - 1
			}
			if topN < 0 {
				topN = 0
			}
			generationInt = int(topStrInt)
			isNotValidTopN = false
		}
	}

	topIndividuals, err := GetTopNIndividualInGenerationX(sortedIndividuals, kind,
		isMoreFitnessBetter, topN, int(generationInt))
	if err != nil {
		return "", err
	}
	bannerBuilder := strings.Builder{}
	switch kind {
	case IndividualAntagonist:
		name := "Antagonist"
		bannerStr := banner(fmt.Sprintf("Top %d %s in Gen %d", topN, name, generationInt))
		bannerBuilder.WriteString(bannerStr)
	case IndividualProtagonist:
		name := "Protagonist"
		bannerStr := banner(fmt.Sprintf("Top %d %s in Gen %d", topN, name, generationInt))
		bannerBuilder.WriteString(bannerStr)
	}
	for i := range topIndividuals.Individuals {
		builder := topIndividuals.Individuals[i].ToString()
		bannerBuilder.WriteString(builder.String())
		bannerBuilder.WriteString("\n")
	}
	return bannerBuilder.String(), nil
}

func interactiveGetIndividualXInGenY(reader *bufio.Reader, sortedIndividuals []*Generation,
	isMoreFitnessBetter bool) (string, error) {
	var generationInt int
	isNotValidInt := true

	for isNotValidInt {
		fmt.Print(fmt.Sprintf("Input a Generation Number [0,%d)", len(sortedIndividuals)))
		fmt.Print("---->")
		generation, _ := reader.ReadString('\n')
		// convert CRLF to LF
		generation = strings.Replace(generation, "\n", "", -1)
		genInt, err := strconv.ParseInt(generation, 10, 64)
		if err != nil {
			isNotValidInt = true
			fmt.Print(fmt.Sprintf("Please enter a NUMBER between [0,%d)",
				len(sortedIndividuals)))
		} else {
			generationInt = int(genInt)
			isNotValidInt = false
		}
	}
	individualIndex := 0
	isNotValidTopN := true
	for isNotValidTopN {
		fmt.Print(fmt.Sprintf("Input Top(N) Generation Number [0,%d)", len(sortedIndividuals[0].Antagonists)))
		fmt.Print("------->")
		individualIndexStr, _ := reader.ReadString('\n')
		// convert CRLF to LF
		individualIndexStr = strings.Replace(individualIndexStr, "\n", "", -1)
		topStrInt, err := strconv.ParseInt(individualIndexStr, 10, 64)
		individualIndex = int(topStrInt)
		if err != nil {
			isNotValidTopN = true
			fmt.Print(fmt.Sprintf("Please enter a NUMBER between [0,%d)",
				len(sortedIndividuals)))
		} else {
			if individualIndex > len(sortedIndividuals[0].Antagonists) {
				individualIndex = len(sortedIndividuals[0].Antagonists) - 1
			}
			if individualIndex < 0 {
				individualIndex = 0
			}
			generationInt = int(topStrInt)
			isNotValidTopN = false
		}
	}

	kind := 0
	isValidKind := true
	for isValidKind {
		fmt.Print(fmt.Sprintf("Antagonist(0) or Protagonist(1)"))
		fmt.Print("--------->")
		kindStr, _ := reader.ReadString('\n')
		// convert CRLF to LF
		kindStr = strings.Replace(kindStr, "\n", "", -1)
		kindInt, err := strconv.ParseInt(kindStr, 10, 64)
		kind = int(kindInt)
		if err != nil {
			isValidKind = true
			fmt.Print(fmt.Sprintf("Please enter a NUMBER between [0,1]"))
		} else {
			if kind > 1 {
				kind = 1
			}
			if kind < 0 {
				kind = 0
			}
			isValidKind = false
		}
	}

	topIndividuals, err := GetTopNIndividualInGenerationX(sortedIndividuals, kind,
		isMoreFitnessBetter, individualIndex, int(generationInt))
	if err != nil {
		return "", err
	}
	bannerBuilder := strings.Builder{}
	switch kind {
	case IndividualAntagonist:
		name := "Antagonist"
		bannerStr := banner(fmt.Sprintf("%s[%d] in Gen %d", name, individualIndex, generationInt))
		bannerBuilder.WriteString(bannerStr)
	case IndividualProtagonist:
		name := "Protagonist"
		bannerStr := banner(fmt.Sprintf("%s[%d] in Gen %d", name, individualIndex, generationInt))
		bannerBuilder.WriteString(bannerStr)
	}
	builder := topIndividuals.Individuals[individualIndex].ToString()
	bannerBuilder.WriteString(builder.String())
	bannerBuilder.WriteString("\n")

	return bannerBuilder.String(), nil
}

func interactiveWriteToFile(reader *bufio.Reader, evolutionResult *EvolutionResult, params EvolutionParams) (string, error) {
	var fileName string
	isNotValidFileName := true

	for isNotValidFileName {
		fmt.Print("---->File name to output statistics: ")
		filePath, _ := reader.ReadString('\n')
		// convert CRLF to LF
		filePath = strings.Replace(filePath, "\n", "", -1)
		_, err := os.Create(filePath)
		if err != nil {
			isNotValidFileName = true
			fmt.Print(fmt.Sprintf("Cannot write to %s | %s", filePath, err.Error()))
		} else {
			fileName = filePath
			isNotValidFileName = false
		}
	}
	path, err := WritetoFile(fileName, evolutionResult, params, 0)
	return path, err
	//return fmt.Sprintf("Feature not yet implemented. Will not write to %s", fileName), nil
}

func interactiveSearchForTreeShape(reader *bufio.Reader, sortedGenerations []*Generation, params EvolutionParams) (
	string,
	error) {
	//var fileName string
	isNotValidSearch := true
	builder := strings.Builder{}

	for isNotValidSearch {
		fmt.Print("---->Input a mathematical expression to search for. No parentheses needed: ")
		mathExpression, _ := reader.ReadString('\n')
		// convert CRLF to LF
		mathExpression = strings.Replace(mathExpression, "\n", "", -1)

		_, _, mathematicalExpression, err := ParseString(params.SpecParam.Expression, params.SpecParam.AvailableVariablesAndOperators.Operators, params.SpecParam.AvailableVariablesAndOperators.Variables)
		if err != nil {
			return "", fmt.Errorf(err.Error())
		}
		queryTree := DualTree{}
		err = queryTree.FromSymbolicExpressionSet2(mathematicalExpression)
		if err != nil {

			return "", fmt.Errorf("interactiveSearchForTreeShap | cannot parse symbolic expression tree to convert starter tree to a" +
				" mathematical" +
				" expression")
		}
		starterTreeAsMathematicalExpression, err := queryTree.ToMathematicalString()
		builder.WriteString(fmt.Sprintf("Searching for: %s\n", starterTreeAsMathematicalExpression))

		if err != nil {
			return "", fmt.Errorf("main | failed to convert starter tree to a mathematical expression")
		}
		isNotValidSearch = false
		results := make([]*Individual, 0)
		for i := range sortedGenerations {
			for a := range sortedGenerations[i].Antagonists {
				antagonistMathString, err := sortedGenerations[i].Antagonists[a].Program.T.ToMathematicalString()
				if err != nil {
					return "", fmt.Errorf(err.Error())
				}
				if strings.Contains(antagonistMathString, starterTreeAsMathematicalExpression) {
					results = append(results, sortedGenerations[i].Antagonists[a])
				}
			}
			for p := range sortedGenerations[i].Protagonists {
				protagonistMathString, err := sortedGenerations[i].Protagonists[p].Program.T.ToMathematicalString()
				if err != nil {
					return "", fmt.Errorf(err.Error())
				}
				if strings.Contains(starterTreeAsMathematicalExpression, protagonistMathString) {
					results = append(results, sortedGenerations[i].Protagonists[p])
				}
			}
		}
		if len(results) == 0 {
			builder.WriteString(fmt.Sprintf("No match found. Searched %d individuals on both sides.\n",
				(params.EachPopulationSize * 2 * params.GenerationsCount)))
		} else {
			for i := range results {
				individualString := results[i].ToString()
				builder.WriteString(individualString.String() + "\n")
			}
		}

	}
	return builder.String(), nil
}

func WritetoFile(path string, evolutionResult *EvolutionResult, params EvolutionParams, count int) (string, error) {
	csvOutput := CSVOutput{
		Generational: make([]GenerationalStatistics, len(evolutionResult.SortedGenerationIndividuals)),
		Epochal: make([]EpochalStatistics, len(evolutionResult.SortedGenerationIndividuals[0].Protagonists[0].
			Fitness)),
	}

	coevolutionaryAverages := evolutionResult.CoevolutionaryAverages

	// GENERATIONAL
	for i := range coevolutionaryAverages {
		csvOutput.Generational[i].Generation = i + 1
		csvOutput.Generational[i].Run = count + 1
		csvOutput.Generational[i].Spec = params.SpecParam.Expression

		// ########################################## ANTAGONISTS ###################################################
		topAntagonist := evolutionResult.SortedGenerationIndividuals[i].Antagonists[0]
		topAntagonistEquation, _ := topAntagonist.Program.T.ToMathematicalString()

		csvOutput.Generational[i].AverageAntagonist = coevolutionaryAverages[i].AntagonistResult
		csvOutput.Generational[i].TopAntagonist = topAntagonist.TotalFitness
		csvOutput.Generational[i].TopAntagonistBirthGen = topAntagonist.BirthGen
		csvOutput.Generational[i].TopAntagonistDelta = topAntagonist.FitnessDelta
		csvOutput.Generational[i].TopAntagonistEquation = topAntagonistEquation
		csvOutput.Generational[i].TopAntagonistFavoriteStrategy = dominantStrategy(*topAntagonist)
		csvOutput.Generational[i].TopAntagonistStrategies = strategiesToString(*topAntagonist)

		// ########################################## PROTAGONISTS ###################################################
		topProtagonist := evolutionResult.SortedGenerationIndividuals[i].Protagonists[0]
		topProtagonistEquation, _ := topProtagonist.Program.T.ToMathematicalString()

		csvOutput.Generational[i].AverageProtagonist = coevolutionaryAverages[i].ProtagonistResult
		csvOutput.Generational[i].TopProtagonist = topProtagonist.TotalFitness
		csvOutput.Generational[i].TopProtagonistBirthGen = topProtagonist.BirthGen
		csvOutput.Generational[i].TopProtagonistDelta = topAntagonist.FitnessDelta
		csvOutput.Generational[i].TopProtagonistEquation = topProtagonistEquation
		csvOutput.Generational[i].TopProtagonistFavoriteStrategy = dominantStrategy(*topProtagonist)
		csvOutput.Generational[i].TopProtagonistStrategies = strategiesToString(*topProtagonist)
	}

	topProtagonist := evolutionResult.SortedGenerationIndividuals[0].Protagonists[0]
	topAntagonist:= evolutionResult.SortedGenerationIndividuals[0].Antagonists[0]
	finalProtagonist:= evolutionResult.FinalProtagonist
	finalAntagonist:= evolutionResult.FinalAntagonist

	for i := 0; i < len(csvOutput.Epochal); i++ {
		csvOutput.Epochal[i].Epoch = i+1
		csvOutput.Epochal[i].TopAntagonist = topAntagonist.Fitness[i]
		csvOutput.Epochal[i].TopProtagonist = topProtagonist.Fitness[i]

		csvOutput.Epochal[i].FinalAntagonist = finalAntagonist.Fitness[i]
		csvOutput.Epochal[i].FinalProtagonist= finalProtagonist.Fitness[i]
	}

	// Internal Variance of Ultimate Individuals
	err := os.Mkdir(params.StatisticsOutput.OutputDir, 0755)
	innerFolder := strings.ReplaceAll(path, ".json", "")
	err = os.Mkdir(innerFolder, 0755)
	g := strings.SplitAfter(path, "/")

	mainDir := g[0]
	subDirInfo := g[1]
	subsubDirName := strings.ReplaceAll(g[2], ".json", "")

	csvMap := map[string]interface{}{
		"generational": csvOutput.Generational,
		"epochal": csvOutput.Epochal,
		//"ultimateIndividuals": csvOutput.CSVUltimateIndividuals,
		//"topEquations": csvOutput.CSVTopEquations,
		//"topStrategies": csvOutput.CSVTopPerGeneration,
		//"topEquationVariation": csvOutput.CSVTopEquationVariation,
		//"strategiesPerGeneration": csvOutput.CSVStrategyPerGeneration,
		//"equationsPerGeneration": csvOutput.CSVEquationsPerGeneration,
		//"finalIndividuals": csvOutput.CSVFinalIndividuals,
	}
	err = WriteCSVWithMap(csvMap, mainDir, subDirInfo, subsubDirName, params.InternalCount)
	if err != nil {
		return path, err
	}
	return path, nil
}

func dominantStrategy(individual Individual) string {
	domStrat := map[string]int{}
	for i := range individual.Strategy {
		strategy := string(individual.Strategy[i])

		stratCount := domStrat[strategy]
		domStrat[strategy] = stratCount + 1
	}

	var topStrategy string
	counter := 0
	for e := range domStrat {
		if domStrat[e] > counter {
			topStrategy = e
		}
	}
	return topStrategy
}

func strategiesToString(individual Individual) string {
	sb := strings.Builder{}
	for _, strategy := range individual.Strategy {
		sb.WriteString(string(strategy))
		sb.WriteString("|")
	}
	return sb.String()
}

func WriteCSVWithMap(csvFileMap map[string]interface{}, mainDir, subDirInfo, subsubDirName string, count int) (err error) {
	for name := range csvFileMap {
		pathCSV := fmt.Sprintf("%s%s%s/%s-%d%s", mainDir, subDirInfo, subsubDirName, name, count, ".csv")
		fileCSV, err := os.Create(pathCSV)
		if err != nil {
			return err
		}
		defer fileCSV.Close()

		writer := gocsv.DefaultCSVWriter(fileCSV)
		if writer.Error() != nil {
			return writer.Error()
		}
		err = gocsv.Marshal(csvFileMap[name], fileCSV)
		if err != nil {
			return err
		}
	}
	return err
}

type JSONCoalescedOutput struct {
	Name            string       `json:"name" csv:"name"`
	CoalescedOutput []JSONOutput `json:"coalescedOutput" csv:"coalescedOutput"`
}

type CSVOutput struct {
	Generational []GenerationalStatistics `csv:"generational"`
	Epochal                   []EpochalStatistics `csv:"epochal"`
}

type CSVEquations struct {
	Protagonist string `csv:"protagonistEquation"`
	Antagonist  string `csv:"protagonistEquation"`
	Spec        string `csv:"protagonistEquation"`
}

type CSVStrategy struct {
	Protagonist string `csv:"protagonistStrategy"`
	Antagonist  string `csv:"antagonistStrategy"`
}

// GenerationalStatistics refer to statistics per generation.
// So Top or Bottom refer to the best or worst in the given generation and not a cumulative of the evolutionary process.
type GenerationalStatistics struct {
	Generation                     int     `csv:"generation"`
	AverageAntagonist              float64 `csv:"averageAntagonist"`
	AverageProtagonist             float64 `csv:"averageProtagonist"`
	TopAntagonist                  float64 `csv:"topAntagonist"`
	TopProtagonist                 float64 `csv:"topProtagonist"`
	TopAntagonistStrategies        string  `csv:"antagonistStrategies"`
	TopProtagonistStrategies       string  `csv:"topProtagonistStrategies"`
	TopAntagonistFavoriteStrategy  string  `csv:"topAntagonistDominantStrategy"`
	TopProtagonistFavoriteStrategy string  `csv:"topProtagonistDominantStrategy"`
	TopAntagonistBirthGen          int     `csv:"topAntagonistBirthGen"`
	TopProtagonistBirthGen         int     `csv:"topProtagonistBirthGen"`
	TopAntagonistDelta             float64 `csv:"topAntagonistDelta"`
	TopProtagonistDelta            float64 `csv:"topProtagonistDelta"`
	TopAntagonistEquation          string  `csv:"topAntagonistEquation"`
	TopProtagonistEquation         string  `csv:"topProtagonistEquation"`
	Spec                           string  `csv:"spec"`
	Run                            int     `csv:"runNumber"`
}

type CSVMetadata struct {
	Title           string          `json:"title" csv:"title"`
	SubTitle        string          `json:"subTitle" csv:"subTitle"`
	Description     string          `json:"description" csv:"description"`
	EvolutionParams EvolutionParams `json:"evolutionaryParams"`
}

type EpochalStatistics struct {
	TopAntagonist float64
	TopProtagonist float64
	FinalAntagonist float64
	FinalProtagonist float64
	Epoch int
}


type CSVUltimateIndividual struct {
	Protagonist float64 `csv:"protagonistCoordinates"`
	Antagonist  float64 `csv:"antagonistCoordinates"`
	Independent float64 `csv:"independent"`
}

type CSVTopPerGeneration struct {
	Protagonist float64 `csv:"protagonistCoordinates"`
	Antagonist  float64 `csv:"antagonistCoordinates"`
	Independent float64 `csv:"independent"`
}

type CSVGeneric struct {
	//Title       string      `json:"title" csv:"title"`
	//SubTitle    string      `json:"subTitle" csv:"subTitle"`
	//Description string      `json:"description" csv:"description"`
	Protagonist []float64 `json:"protagonistCoordinates" csv:"protagonistCoordinates"`
	Antagonist  []float64 `json:"antagonistCoordinates" csv:"antagonistCoordinates"`
}

type CSVCoalescedOutput struct {
	CSVAverages               []GenerationalStatistics `csv:"averages"`
	CSVUltimateIndividuals    []CSVUltimateIndividual  `csv:"ultimateIndividuals"`
	CSVTopPerGeneration       []CSVTopPerGeneration    `csv:"topPerGeneration"`
	CSVBottomPerGeneration    []GenerationalStatistics `csv:"bottomPerGeneration"`
	CSVTopEquations           CSVEquations             `csv:"bottomPerGeneration"`
	CSVTopStrategies          CSVEquations             `csv:"topStrategies"`
	CSVEquationsPerGeneration []CSVEquations           `csv:"equationesPerGeneration"`
	CSVStrategyPerGeneration  []CSVStrategy            `csv:"strategyPerGeneration"`
}

type JSONOutput struct {
	Averages JSONGeneric `json:"averages" csv:"averages"`

	// UltimateIndividuals returns the internal variance of the best individuals in all generations
	UltimateIndividuals JSONGeneric `json:"ultimateIndividuals" csv:"ultimateIndividuals"`
	// BottomPerGeneration returns the best kind of individual in each generation
	TopPerGeneration JSONGeneric `json:"topPerGeneration" csv:"topPerGeneration"`
	// BottomPerGeneration returns the worst kind of individual in each generation
	BottomPerGeneration JSONGeneric `json:"bottomPerGeneration" csv:"bottomPerGeneration"`

	Equations JSONEquations `json:"equations" csv:"equations"`

	//UltimateIndividualsDelta JSONGeneric `json:"averages"`
	//
	//FinalGenIndividuals JSONGeneric `json:"averages"`
}

type JSONGeneric struct {
	Title       string      `json:"title" csv:"title"`
	SubTitle    string      `json:"subTitle" csv:"subTitle"`
	Description string      `json:"description" csv:"description"`
	Protagonist Coordinates `json:"protagonistCoordinates" csv:"protagonistCoordinates"`
	Antagonist  Coordinates `json:"antagonistCoordinates" csv:"antagonistCoordinates"`
}

type JSONEquations struct {
	Spec                JSONEquation `json:"spec" csv:"spec"`
	UltimateAntagonist  JSONEquation `json:"ultimateAntagonist" csv:"ultimateAntagonist"`
	UltimateProtagonist JSONEquation `json:"ultimateProtagonist" csv:"ultimateProtagonist"`
}

type JSONEquation struct {
	Title      string `json:"title" csv:"title"`
	Expression string `json:"expression" csv:"expression"`
	Seed       int    `json:"seed" csv:"seed"`
	Range      int    `json:"range" csv:"range"`
}

type Coordinates struct {
	IndependentCoordinates []float64 `json:"independentCoordinates" csv:"independentCoordinates"`
	DependentCoordinates   []float64 `json:"dependentCoordinates" csv:"dependentCoordinates"`
}
