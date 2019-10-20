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
	err := SortGenerationsThoroughly(generations, isMoreFitnessBetter)
	if err != nil {
		return err
	}
	e.SortedGenerationIndividuals = generations

	// Calculate Top Individuals
	topAntagonist, topProtagonist, err := GetTopIndividualInAllGenerations(generations, isMoreFitnessBetter)
	if err != nil {
		return err
	}
	e.TopAntagonist = topAntagonist
	e.TopProtagonist = topProtagonist

	// Calculate Generational Averages
	coevolutionaryAverages, err := GetGenerationalFitnessAverage(generations)
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
	sb.WriteString("ANT | PRO\n")
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

	fmt.Println("0. Top Protagonist")
	fmt.Println("1. Top Antagonist")
	fmt.Println("2. Average Generational Fitness")
	fmt.Println("3. Top Antagonist in Gen(X)")
	fmt.Println("4. Top Protagonist in Gen(X)")
	fmt.Println("5. Top N Antagonists in Gen(x)")
	fmt.Println("6. Top N Protagonists in Gen(x)")
	fmt.Println("7. Individual (X) in Generation (Y)")
	fmt.Println("8. Output Results to File")
	fmt.Println("\nType !q to exit!")

	reader := bufio.NewReader(os.Stdin)
	isRunning := true
	for isRunning {
		fmt.Print("->")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		if strings.Compare("!q", text) == 0 {
			isRunning = false
		}

		switch text {
		case "0":
			strBuilder := e.TopProtagonist.ToString()
			bannerStr := banner("Top Protagonist")
			fmt.Println(bannerStr)
			fmt.Println(strBuilder.String())
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
