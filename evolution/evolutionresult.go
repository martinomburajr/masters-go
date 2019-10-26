package evolution

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type EvolutionResult struct {
	HasBeenAnalyzed     bool
	TopAntagonist       *Individual
	TopProtagonist      *Individual
	IsMoreFitnessBetter bool

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

	// Calculate Generational Averages
	coevolutionaryAverages, err := GetGenerationalFitnessAverage(sortedGenerations)
	if err != nil {
		return err
	}
	e.CoevolutionaryAverages = coevolutionaryAverages
	e.HasBeenAnalyzed = true

	err = writetoFile(params.StatisticsOutput.OutputPath, e, params)
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
// 2. Average Generational Fitness
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
		fmt.Println("2. Average Generational Fitness")
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
	err := writetoFile(fileName, evolutionResult, params)
	return "", err
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

func writetoFile(path string, evolutionResult *EvolutionResult, params EvolutionParams) error {
	jsonOutput := JSONOutput{
		Averages: JSONGeneric{
			Antagonist: Coordinates{
				DependentCoordinates: make([]float64,
					len(evolutionResult.SortedGenerationIndividuals)),
				IndependentCoordinates: make([]float64,
					len(evolutionResult.SortedGenerationIndividuals)),
			},
			Protagonist: Coordinates{
				DependentCoordinates: make([]float64,
					len(evolutionResult.SortedGenerationIndividuals)),
				IndependentCoordinates: make([]float64,
					len(evolutionResult.SortedGenerationIndividuals)),
			},
		},
		UltimateIndividuals: JSONGeneric{
			Title: "Variance in Ultimate Individuals Internal Fitness",
			Antagonist: Coordinates{
				DependentCoordinates: make([]float64,
					len(evolutionResult.SortedGenerationIndividuals[0].Antagonists[0].Fitness)),
				IndependentCoordinates: make([]float64,
					len(evolutionResult.SortedGenerationIndividuals[0].Antagonists[0].Fitness)),
			},
			Protagonist: Coordinates{
				DependentCoordinates: make([]float64,
					len(evolutionResult.SortedGenerationIndividuals[0].Antagonists[0].Fitness)),
				IndependentCoordinates: make([]float64,
					len(evolutionResult.SortedGenerationIndividuals[0].Antagonists[0].Fitness)),
			},
		},
		TopPerGeneration: JSONGeneric{
			Antagonist: Coordinates{
				DependentCoordinates: make([]float64,
					len(evolutionResult.SortedGenerationIndividuals)),
				IndependentCoordinates: make([]float64,
					len(evolutionResult.SortedGenerationIndividuals)),
			},
			Protagonist: Coordinates{
				DependentCoordinates: make([]float64,
					len(evolutionResult.SortedGenerationIndividuals)),
				IndependentCoordinates: make([]float64,
					len(evolutionResult.SortedGenerationIndividuals)),
			},
		},
		BottomPerGeneration: JSONGeneric{
			Antagonist: Coordinates{
				DependentCoordinates: make([]float64,
					len(evolutionResult.SortedGenerationIndividuals)),
				IndependentCoordinates: make([]float64,
					len(evolutionResult.SortedGenerationIndividuals)),
			},
			Protagonist: Coordinates{
				DependentCoordinates: make([]float64,
					len(evolutionResult.SortedGenerationIndividuals)),
				IndependentCoordinates: make([]float64,
					len(evolutionResult.SortedGenerationIndividuals)),
			},
		},
	}

	jsonOutput.Averages.Title = "Generational Averages for Antagonists vs Protagonists"
	jsonOutput.Averages.SubTitle = fmt.Sprintf("Fitness (is more better?): %t", evolutionResult.IsMoreFitnessBetter)

	coevolutionaryAverages := evolutionResult.CoevolutionaryAverages
	for i := range coevolutionaryAverages {
		// ##################### ANTAGONISTS #########################
		// Averages
		jsonOutput.Averages.Antagonist.IndependentCoordinates[i] = float64(i)
		jsonOutput.Averages.Antagonist.DependentCoordinates[i] = coevolutionaryAverages[i].AntagonistResult
		// Top Per Generation
		jsonOutput.TopPerGeneration.Antagonist.IndependentCoordinates[i] = float64(i)
		jsonOutput.TopPerGeneration.Antagonist.DependentCoordinates[i] = evolutionResult.
			SortedGenerationIndividuals[i].Antagonists[0].TotalFitness
		//Bottom Per Generation
		jsonOutput.BottomPerGeneration.Antagonist.IndependentCoordinates[i] = float64(i)
		jsonOutput.BottomPerGeneration.Antagonist.DependentCoordinates[i] = evolutionResult.
			SortedGenerationIndividuals[i].Antagonists[len(evolutionResult.
			SortedGenerationIndividuals[i].Antagonists)-1].TotalFitness

		// ##################### PROTAGONISTS #########################
		// Averages
		jsonOutput.Averages.Protagonist.IndependentCoordinates[i] = float64(i)
		jsonOutput.Averages.Protagonist.DependentCoordinates[i] = coevolutionaryAverages[i].ProtagonistResult
		// Top Per Generation
		jsonOutput.TopPerGeneration.Protagonist.IndependentCoordinates[i] = float64(i)
		jsonOutput.TopPerGeneration.Protagonist.DependentCoordinates[i] = evolutionResult.
			SortedGenerationIndividuals[i].Protagonists[0].TotalFitness
		//Bottom Per Generation
		jsonOutput.BottomPerGeneration.Protagonist.IndependentCoordinates[i] = float64(i)
		jsonOutput.BottomPerGeneration.Protagonist.DependentCoordinates[i] = evolutionResult.
			SortedGenerationIndividuals[i].Protagonists[len(evolutionResult.
			SortedGenerationIndividuals[i].Protagonists)-1].TotalFitness
	}

	// Internal Variance of Ultimate Individuals
	for i := 0; i < len(evolutionResult.TopAntagonist.Fitness); i++ {
		jsonOutput.UltimateIndividuals.Antagonist.IndependentCoordinates[i] = float64(i)
		jsonOutput.UltimateIndividuals.Antagonist.DependentCoordinates[i] = evolutionResult.TopAntagonist.Fitness[i]

		jsonOutput.UltimateIndividuals.Protagonist.IndependentCoordinates[i] = float64(i)
		jsonOutput.UltimateIndividuals.Protagonist.DependentCoordinates[i] = evolutionResult.TopProtagonist.Fitness[i]
	}

	topProtagonistMathExpression, err := evolutionResult.TopProtagonist.Program.T.ToMathematicalString()
	if err != nil {
		return err
	}
	topAntagonistMathExpression, err := evolutionResult.TopAntagonist.Program.T.ToMathematicalString()
	if err != nil {
		return err
	}
	// Equations
	jsonOutput.Equations = JSONEquations{
		Spec: JSONEquation{
			Title:      "Spec",
			Seed:       params.SpecParam.Seed,
			Range:      params.SpecParam.Range,
			Expression: params.SpecParam.Expression,
		},
		UltimateAntagonist: JSONEquation{
			Title:      "Ult-Antagonist",
			Seed:       params.SpecParam.Seed,
			Range:      params.SpecParam.Range,
			Expression: topAntagonistMathExpression,
		},
		UltimateProtagonist: JSONEquation{
			Title:      "Ult-AProtagonist",
			Seed:       params.SpecParam.Seed,
			Range:      params.SpecParam.Range,
			Expression: topProtagonistMathExpression,
		},
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	fmt.Printf("Wrote to file: %s", path)
	return json.NewEncoder(file).Encode(jsonOutput)
}

type JSONCoalescedOutput struct {
	Name string `json:"name"`
	TopProtagonists []JSONGeneric `json:"coalescedTopProtagonists"`
	TopAntagonists []JSONGeneric `json:"coalescedTopAntagonists"`
	Averages []JSONGeneric `json:"coalescedAverages"`
	Equations []JSONEquations `json:"coalescedEquations"`
	CoalescedOutput []JSONOutput `json:"coalescedOutput"`
}

type JSONOutput struct {
	Averages JSONGeneric `json:"averages"`

	// UltimateIndividuals returns the internal variance of the best individuals in all generations
	UltimateIndividuals JSONGeneric `json:"ultimateIndividuals"`
	// BottomPerGeneration returns the best kind of individual in each generation
	TopPerGeneration JSONGeneric `json:"topPerGeneration"`
	// BottomPerGeneration returns the worst kind of individual in each generation
	BottomPerGeneration JSONGeneric `json:"bottomPerGeneration"`

	Equations JSONEquations `json:"equations"`

	//UltimateIndividualsDelta JSONGeneric `json:"averages"`
	//
	//FinalGenIndividuals JSONGeneric `json:"averages"`
}

type JSONGeneric struct {
	Title       string      `json:"title"`
	SubTitle    string      `json:"subTitle"`
	Description string      `json:"description"`
	Protagonist Coordinates `json:"protagonistCoordinates"`
	Antagonist  Coordinates `json:"antagonistCoordinates"`
}

type JSONEquations struct {
	Spec                JSONEquation `json:"spec"`
	UltimateAntagonist  JSONEquation `json:"ultimateAntagonist"`
	UltimateProtagonist JSONEquation `json:"ultimateProtagonist"`
}

type JSONEquation struct {
	Title      string `json:"title"`
	Expression string `json:"expression"`
	Seed       int    `json:"seed"`
	Range      int    `json:"range"`
}

type Coordinates struct {
	IndependentCoordinates []float64 `json:"independentCoordinates"`
	DependentCoordinates   []float64 `json:"dependentCoordinates"`
}
