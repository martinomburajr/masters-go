package evolution

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type EvolutionResult struct {
	HasBeenAnalyzed bool
	TopAntagonist   *Individual
	TopProtagonist  *Individual
	IsMoreFitnessBetter bool

	CoevolutionaryAverages []generationalCoevolutionaryAverages

	SortedGenerationIndividuals []*Generation
}

type multiIndividualsPerGeneration struct {
	Generation  *Generation
	Individuals []*Individual
}

type generationalCoevolutionaryAverages struct {
	Generation       *Generation
	AntagonistResult float64
	ProtagonistResult float64
}


func (e *EvolutionResult) Analyze(generations []*Generation, isMoreFitnessBetter bool) error {
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

	return nil
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
func (e *EvolutionResult) StartInteractiveTerminal() error {
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
		fmt.Println("8. Output Results to File")
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
			str, err := interactiveWriteToFile(reader, e)
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
		name:= "Antagonist"
		bannerStr := banner(fmt.Sprintf("Top %s in Gen %d", name, generationInt))
		bannerBuilder.WriteString(bannerStr)
	case IndividualProtagonist:
		name:= "Protagonist"
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
		fmt.Print(fmt.Sprintf("Input Top(N) Generation Number [0,%d)", len(sortedIndividuals[0].Antagonists)))
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
				topN = len(sortedIndividuals[0].Antagonists)-1
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
		name:= "Antagonist"
		bannerStr := banner(fmt.Sprintf("Top %d %s in Gen %d", topN, name, generationInt))
		bannerBuilder.WriteString(bannerStr)
	case IndividualProtagonist:
		name:= "Protagonist"
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
				individualIndex = len(sortedIndividuals[0].Antagonists)-1
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
		name:= "Antagonist"
		bannerStr := banner(fmt.Sprintf("%s[%d] in Gen %d", name, individualIndex, generationInt))
		bannerBuilder.WriteString(bannerStr)
	case IndividualProtagonist:
		name:= "Protagonist"
		bannerStr := banner(fmt.Sprintf("%s[%d] in Gen %d", name, individualIndex, generationInt))
		bannerBuilder.WriteString(bannerStr)
	}
	builder := topIndividuals.Individuals[individualIndex].ToString()
	bannerBuilder.WriteString(builder.String())
	bannerBuilder.WriteString("\n")

	return bannerBuilder.String(), nil
}

func interactiveWriteToFile(reader *bufio.Reader, evolutionResult *EvolutionResult) (string, error) {
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
	return fmt.Sprintf("Feature not yet implemented. Will not write to %s", fileName), nil
}