package evolution

import (
	"fmt"
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


func (e *EvolutionResult) Analyze(generations []*Generation, isMoreFitnessBetter bool, topN int) error {
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
		"####################################### AVERAGE ANTAGONISTS VS PROTAGONISTS PER GENERATION" +
		" #####################################################\n\n"))
	sb.WriteString("ANT | PRO\n")
	for i := range e.CoevolutionaryAverages {
		res := e.CoevolutionaryAverages[i].AntagonistResult
		resPro := e.ProtagonistAverageAcrossGenerations[i].AntagonistResult
		float := strconv.FormatFloat(res, 'g', 03, 64)
		floatPro := strconv.FormatFloat(resPro, 'g', 03, 64)

		gen := strconv.Itoa(i)
		sb.WriteString("gen" + gen + ": " + float + " | " + floatPro + "\n")
	}
	sb.WriteString("\n")
	return sb, nil
}

//func (e *EvolutionResult) PrintCumGenerationSummary() (strings.Builder, error) {
//	if e.ProtagonistCumAcrossGenerations == nil {
//		return strings.Builder{},
//			fmt.Errorf("PrintCumGenerationSummary | cannot format as protagonist average field is nil | Run" +
//				" analyze")
//	}
//	if e.AntagonistCumAcrossGenerations == nil {
//		return strings.Builder{},
//			fmt.Errorf("PrintCumGenerationSummary | cannot format as antagonist average field is nil | Run" +
//				" analyze")
//	}
//
//	sb := strings.Builder{}
//	sb.WriteString(fmt.Sprintf("" +
//		"####################################### CUMULATIVE ANTAGONISTS VS PROTAGONISTS PER GENERATION" +
//		" #####################################################\n\n"))
//	sb.WriteString("ANT | PRO\n")
//	for i := range e.AntagonistCumAcrossGenerations {
//		res := e.AntagonistCumAcrossGenerations[i].Result
//		resPro := e.ProtagonistCumAcrossGenerations[i].Result
//		float := strconv.FormatFloat(res, 'g', 03, 64)
//		floatPro := strconv.FormatFloat(resPro, 'g', 03, 64)
//
//		gen := strconv.Itoa(i)
//		sb.WriteString("gen" + gen + ": " + float + " | " + floatPro + "\n")
//	}
//	sb.WriteString("\n")
//	return sb, nil
//}


func (e *EvolutionResult) PrintTopIndividualSummary(kind int) (strings.Builder, error) {
	sb := strings.Builder{}
	var name string
	if kind == IndividualAntagonist {
		if e.TopAntagonist.Individual == nil {
			return strings.Builder{},
				fmt.Errorf("PrintTopIndividualSummary | cannot format as field is nil | Run analyze")
		}
		name = "ANTAGONIST"
		sb.WriteString(fmt.Sprintf("############### TOP %s IN ALL GENERATIONS"+" #######################\n", name))
		toString := e.TopAntagonist.Individual.ToString()
		sb.WriteString(toString.String())
	} else if kind == IndividualProtagonist {
		if e.TopProtagonist.Individual == nil {
			return strings.Builder{},
				fmt.Errorf("PrintTopIndividualSummary | cannot format as field is nil | Run analyze")
		}
		name = "PROTAGONIST"
		sb.WriteString(fmt.Sprintf("############### TOP %s IN ALL GENERATIONS"+" #######################\n", name))
		toString := e.TopProtagonist.Individual.ToString()
		sb.WriteString(toString.String())
	}
	return sb, nil
}

func (e *EvolutionResult) PrintTopNInFinalGeneration(kind int) (strings.Builder, error) {
	if !e.HasBeenAnalyzed {
		return strings.Builder{}, fmt.Errorf("PrintTopNInFinalGeneration| Evolution Results Have Not Been Analyzed")
	}

}

//func (e *EvolutionResult) PrintTopNIndividualSummary(kind int) (strings.Builder, error) {
//	sb := strings.Builder{}
//
//	if kind == IndividualProtagonist {
//		sb.WriteString("TopN Protagonists\n\n")
//		for i, topNIndividual := range e.TopNProtagonistsPerGeneration {
//			topIndividualSummary, err := e.PrintTopIndividualSummary(kind)
//			if err != nil {
//				return strings.Builder{}, fmt.Errorf("PrintTopIndividualSummary | %s", err.Error())
//			}
//			s := topIndividualSummary.String()
//			sb.WriteString(fmt.Sprintf("Individual %d\n%s", i+1,s))
//		}
//	} else if kind == IndividualAntagonist {
//		sb.WriteString("TopN Antagonists\n\n")
//		for i, topNIndividual := range e.TopNAntagonistsPerGeneration {
//			topIndividualSummary, err := topNIndividual  e.PrintTopIndividualSummary(kind)
//			if err != nil {
//				return strings.Builder{}, fmt.Errorf("PrintTopIndividualSummary | %s", err.Error())
//			}
//			s := topIndividualSummary.String()
//			sb.WriteString(fmt.Sprintf("Individual %d\n%s", i+1,s))
//		}
//	}
//}
