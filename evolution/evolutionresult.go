package evolution

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type EvolutionResult struct {
	HasBeenAnalyzed        bool
	TopAntagonist          *Individual
	TopProtagonist         *Individual
	TopAntagonistDelta     *Individual
	TopProtagonistDelta    *Individual
	TopAntagonistAvgDelta  *Individual
	TopProtagonistAvgDelta *Individual
	IsMoreFitnessBetter    bool
	FinalAntagonist        *Individual
	FinalProtagonist       *Individual

	CoevolutionaryAverages []GenerationalCoevolutionaryAverages
	Generational           Generational

	SortedGenerationIndividuals           []*Generation
	SortedGenerationIndividualsByDelta    []*Generation
	SortedGenerationIndividualsByDeltaAvg []*Generation
	OutputFile                            string

	Mutex sync.Mutex
}

type multiIndividualsPerGeneration struct {
	Generation  *Generation
	Individuals []*Individual
}

type GenerationalCoevolutionaryAverages struct {
	Generation                     *Generation
	AntagonistFitnessAverages      float64
	ProtagonistFitnessAverages     float64
	AntagonistBestFitnessAverages  float64
	ProtagonistBestFitnessAverages float64
	AntagonistDeltaAverages        float64
	ProtagonistDeltaAverages       float64
	AntagonistBestDeltaAverages    float64
	ProtagonistBestDeltaAverages   float64
}

//Generational averages contain slices of length of the generations in a given run
type Generational struct {
	Antagonists                    []Individual
	Protagonists                   []Individual
	AntagonistsByDelta             []Individual
	ProtagonistsByDelta            []Individual
	AntagonistsByDeltaAvg          []Individual
	ProtagonistsByDeltaAvg         []Individual
	AntagonistFitnessAverages      []float64
	ProtagonistFitnessAverages     []float64
	AntagonistBestFitnessAverages  []float64
	ProtagonistBestFitnessAverages []float64
	AntagonistDeltaAverages        []float64
	ProtagonistDeltaAverages       []float64
	AntagonistBestDeltaAverages    []float64
	ProtagonistBestDeltaAverages   []float64
}

func (e *EvolutionResult) Analyze(generations []*Generation, isMoreFitnessBetter bool,
	params EvolutionParams) error {

	wg := sync.WaitGroup{}
	wg.Add(5)
	go func(generations []*Generation, e *EvolutionResult, wg *sync.WaitGroup) {
		defer wg.Done()
		sortedFinalAntagonists, err := SortIndividuals(generations[len(generations)-1].Antagonists, true)
		if err != nil {
			params.ErrorChan <- err
		}
		e.Mutex.Lock()
		e.FinalAntagonist = sortedFinalAntagonists[0]
		e.Mutex.Unlock()
	}(generations, e, &wg)

	go func(generations []*Generation, e *EvolutionResult, wg *sync.WaitGroup) {
		defer wg.Done()
		sortedFinalProtagonists, err := SortIndividuals(generations[len(generations)-1].Protagonists, true)
		if err != nil {
			params.ErrorChan <- err
		}
		e.Mutex.Lock()
		e.FinalProtagonist = sortedFinalProtagonists[0]
		e.Mutex.Unlock()
	}(generations, e, &wg)

	go func(generations []*Generation, e *EvolutionResult, wg *sync.WaitGroup) {
		defer wg.Done()
		sortedGenerations, err := SortGenerationsThoroughly(generations, isMoreFitnessBetter)
		if err != nil {
			params.ErrorChan <- err
		}
		e.Mutex.Lock()
		e.SortedGenerationIndividuals = sortedGenerations
		e.Mutex.Unlock()
	}(generations, e, &wg)

	go func(generations []*Generation, e *EvolutionResult, wg *sync.WaitGroup) {
		defer wg.Done()
		sortedGenerationsByDelta, err := SortGenerationsThoroughlyByDelta(generations, true, false)
		if err != nil {
			params.ErrorChan <- err
		}
		e.Mutex.Lock()
		e.SortedGenerationIndividualsByDelta = sortedGenerationsByDelta
		e.Mutex.Unlock()
	}(generations, e, &wg)

	go func(generations []*Generation, e *EvolutionResult, wg *sync.WaitGroup) {
		defer wg.Done()
		sortedGenerationsByAvgDelta, err := SortGenerationsThoroughlyByAvgDelta(generations, true, false)
		if err != nil {
			params.ErrorChan <- err
		}
		e.Mutex.Lock()
		e.SortedGenerationIndividualsByDeltaAvg = sortedGenerationsByAvgDelta
		e.Mutex.Unlock()
	}(generations, e, &wg)
	wg.Wait()


	e.Generational.Antagonists = make([]Individual, params.GenerationsCount)
	e.Generational.AntagonistsByDeltaAvg = make([]Individual, params.GenerationsCount)
	e.Generational.AntagonistsByDelta = make([]Individual, params.GenerationsCount)

	for i, v := range e.SortedGenerationIndividuals {
		antagonist := v.Antagonists[0]
		antagonistByDelta := e.SortedGenerationIndividualsByDelta[i].Antagonists[0]
		antagonistByAvgDelta := e.SortedGenerationIndividualsByDeltaAvg[i].Antagonists[0]
		antClone, _ := antagonist.Clone()
		antagonistByDeltaClone, _ := antagonistByDelta.Clone()
		antagonistByAvgDeltaClone, _ := antagonistByAvgDelta.Clone()

		e.Generational.Antagonists[i] = antClone
		e.Generational.AntagonistsByDelta[i] = antagonistByAvgDeltaClone
		e.Generational.AntagonistsByDeltaAvg[i] = antagonistByDeltaClone
	}
	e.Generational.Protagonists = make([]Individual, params.GenerationsCount)
	e.Generational.ProtagonistsByDelta = make([]Individual, params.GenerationsCount)
	e.Generational.ProtagonistsByDeltaAvg = make([]Individual, params.GenerationsCount)

	for i, v := range e.SortedGenerationIndividuals {
		protagonist := v.Protagonists[0]
		rotagonistByDelta := e.SortedGenerationIndividualsByDelta[i].Protagonists[0]
		rotagonistByAvgDelta := e.SortedGenerationIndividualsByDeltaAvg[i].Protagonists[0]
		protagonistClone, _ := protagonist.Clone()
		protagonistByDeltaClone, _ := rotagonistByDelta.Clone()
		protagonistByAvgDeltaClone, _ := rotagonistByAvgDelta.Clone()

		e.Generational.ProtagonistsByDelta[i] = protagonistByDeltaClone
		e.Generational.ProtagonistsByDeltaAvg[i] = protagonistByAvgDeltaClone
		e.Generational.Protagonists[i] = protagonistClone
	}

	// Calculate Top Individuals
	topAntagonist, topProtagonist, err := GetTopIndividualInAllGenerations(e.SortedGenerationIndividuals, isMoreFitnessBetter)
	if err != nil {
		return err
	}
	e.TopAntagonist = topAntagonist
	e.TopProtagonist = topProtagonist

	// Calculate GenerationalStatistics Averages
	coevolutionaryAverages, err := GetGenerationalAverages(e.SortedGenerationIndividuals)
	if err != nil {
		return err
	}
	e.CoevolutionaryAverages = coevolutionaryAverages

	e.Generational.AntagonistFitnessAverages = make([]float64, params.GenerationsCount)
	e.Generational.ProtagonistFitnessAverages = make([]float64, params.GenerationsCount)
	e.Generational.AntagonistBestFitnessAverages = make([]float64, params.GenerationsCount)
	e.Generational.ProtagonistBestFitnessAverages = make([]float64, params.GenerationsCount)
	e.Generational.AntagonistBestDeltaAverages = make([]float64, params.GenerationsCount)
	e.Generational.ProtagonistBestDeltaAverages = make([]float64, params.GenerationsCount)
	e.Generational.AntagonistDeltaAverages = make([]float64, params.GenerationsCount)
	e.Generational.ProtagonistDeltaAverages = make([]float64, params.GenerationsCount)

	for i, v := range coevolutionaryAverages {
		e.Generational.AntagonistFitnessAverages[i] = v.AntagonistFitnessAverages
		e.Generational.ProtagonistFitnessAverages[i] = v.ProtagonistFitnessAverages
		e.Generational.AntagonistBestFitnessAverages[i] = v.AntagonistBestFitnessAverages
		e.Generational.ProtagonistBestFitnessAverages[i] = v.ProtagonistBestFitnessAverages
		e.Generational.AntagonistBestDeltaAverages[i] = v.AntagonistBestDeltaAverages
		e.Generational.ProtagonistBestDeltaAverages[i] = v.ProtagonistBestDeltaAverages
		e.Generational.AntagonistDeltaAverages[i] = v.AntagonistDeltaAverages
		e.Generational.ProtagonistDeltaAverages[i] = v.ProtagonistDeltaAverages
		//e.Generational.Antagonists[i] = v.
	}
	e.HasBeenAnalyzed = true

	return err
}

func (e *EvolutionResult) Clean() {
	e.Generational.Antagonists = nil
	e.Generational.Protagonists = nil
	e.Generational.ProtagonistFitnessAverages = nil
	e.Generational.AntagonistFitnessAverages = nil
	e.FinalAntagonist = nil
	e.FinalProtagonist = nil
	e.TopProtagonist = nil
	e.TopAntagonist = nil
	e.CoevolutionaryAverages = nil
	e.SortedGenerationIndividuals = nil
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
		res := e.CoevolutionaryAverages[i].AntagonistFitnessAverages
		resPro := e.CoevolutionaryAverages[i].ProtagonistFitnessAverages
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
// 0. Top ProtagonistEquation
// 1. Top AntagonistEquation
// 2. Average GenerationalStatistics Fitness
// 3. Top AntagonistEquation in Gen(X)
// 4. Top ProtagonistEquation in Gen(X)
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
		fmt.Println("0. Top AntagonistEquation")
		fmt.Println("1. Top ProtagonistEquation")
		fmt.Println("2. Average GenerationalStatistics Fitness")
		fmt.Println("3. Top AntagonistEquation in Gen(X)")
		fmt.Println("4. Top ProtagonistEquation in Gen(X)")
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
			bannerStr := banner("Top AntagonistEquation")
			fmt.Println(bannerStr)
			fmt.Println(strBuilder.String())
			fmt.Println()
		case "1":
			strBuilder := e.TopProtagonist.ToString()
			bannerStr := banner("Top ProtagonistEquation")
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
		name := "AntagonistEquation"
		bannerStr := banner(fmt.Sprintf("Top %s in Gen %d", name, generationInt))
		bannerBuilder.WriteString(bannerStr)
	case IndividualProtagonist:
		name := "ProtagonistEquation"
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
		name := "AntagonistEquation"
		bannerStr := banner(fmt.Sprintf("Top %d %s in Gen %d", topN, name, generationInt))
		bannerBuilder.WriteString(bannerStr)
	case IndividualProtagonist:
		name := "ProtagonistEquation"
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
		fmt.Print(fmt.Sprintf("AntagonistEquation(0) or ProtagonistEquation(1)"))
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
		name := "AntagonistEquation"
		bannerStr := banner(fmt.Sprintf("%s[%d] in Gen %d", name, individualIndex, generationInt))
		bannerBuilder.WriteString(bannerStr)
	case IndividualProtagonist:
		name := "ProtagonistEquation"
		bannerStr := banner(fmt.Sprintf("%s[%d] in Gen %d", name, individualIndex, generationInt))
		bannerBuilder.WriteString(bannerStr)
	}
	builder := topIndividuals.Individuals[individualIndex].ToString()
	bannerBuilder.WriteString(builder.String())
	bannerBuilder.WriteString("\n")

	return bannerBuilder.String(), nil
}

func interactiveWriteToFile(reader *bufio.Reader, evolutionResult *EvolutionResult, params EvolutionParams) (string, error) {
	//	var fileName string
	//	isNotValidFileName := true
	//
	//	for isNotValidFileName {
	//		fmt.Print("---->File name to output statistics: ")
	//		filePath, _ := reader.ReadString('\n')
	//		// convert CRLF to LF
	//		filePath = strings.Replace(filePath, "\n", "", -1)
	//		_, err := os.Create(filePath)
	//		if err != nil {
	//			isNotValidFileName = true
	//			fmt.Print(fmt.Sprintf("Cannot write to %s | %s", filePath, err.Error()))
	//		} else {
	//			fileName = filePath
	//			isNotValidFileName = false
	//		}
	//	}
	//	//path, err := evolutionResult.WriteToFile(fileName, params)
	//	return "path", err
	//return fmt.Sprintf("Feature not yet implemented. Will not write to %s", fileName), nil
	return "", nil
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

		params.SpecParam.ExpressionParsed = starterTreeAsMathematicalExpression

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
