package evolution

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
)

type EvolutionResult struct {
	hasBeenAnalyzed bool
	TopProtagonist  ResultTopIndividuals
	TopAntagonist   ResultTopIndividuals

	TopProtagonistsPerGeneration ResultInfo1DPerGeneration
	TopAntagonistPerGeneration   ResultInfo1DPerGeneration

	TopNProtagonistsPerGeneration []ResultInfo2DPerGeneration
	TopNAntagonistsPerGeneration  []ResultInfo2DPerGeneration

	ProtagonistAverageAcrossGenerations []ResultInfo1DAveragesPerGeneration
	AntagonistAverageAcrossGenerations  []ResultInfo1DAveragesPerGeneration

	ProtagonistCumAcrossGenerations []ResultInfo1DAveragesPerGeneration
	AntagonistCumAcrossGenerations  []ResultInfo1DAveragesPerGeneration

	SortedProtagonistsPerGeneration ResultInfo2DPerGeneration
	SortedAntagonistsPerGeneration  ResultInfo2DPerGeneration
}

type ResultInfo2DPerGeneration struct {
	Generation *Generation
	Result     []*Individual
}

type ResultInfo1DPerGeneration struct {
	Generation *Generation
	Result     []*Individual
}

type ResultInfo1DAveragesPerGeneration struct {
	Generation *Generation
	Result     float64
}

const (
	FitnessMoreIsBetter = 0
	FitnessLessIsBetter = 1
)

// CalcTopIndividualAllGenerations calculates the top individual in all generations. It uses a fitness comparator to
// determine if a fitness value that is less is better than that of which is more
func CalcTopIndividual(individuals []*Individual,
	fitnessComparator int) (*Individual, error) {
	if individuals == nil {
		return nil, fmt.Errorf("CalcTopIndividual | Individuals cannot be nil")
	}
	if len(individuals) < 1 {
		return nil, fmt.Errorf("CalcTopIndividual | Individuals cannot be empty")
	}

	var individual = &Individual{}
	switch fitnessComparator {
	case FitnessMoreIsBetter:
		individual = &Individual{TotalFitness: math.MinInt64}
	case FitnessLessIsBetter:
		individual = &Individual{TotalFitness: math.MaxInt64}
	}

	for i := range individuals {
		switch fitnessComparator {
		case FitnessMoreIsBetter:
			if individuals[i].TotalFitness >= individual.TotalFitness {
				individual = individuals[i]
			}
		case FitnessLessIsBetter:
			if individuals[i].TotalFitness <= individual.TotalFitness {
				individual = individuals[i]
			}
		default:
			if individuals[i].TotalFitness >= individual.TotalFitness {
				individual = individuals[i]
			}
		}
	}
	return individual, nil
}

// CalcTopIndividualAllGenerations calculates the top individual in all generations. It uses a fitness comparator to
// determine if a fitness value that is less is better than that of which is more
func CalcTopIndividualAllGenerations(generations []*Generation, individualKind int, fitnessComparator int) (ResultTopIndividuals,
	error) {
	if generations == nil {
		return ResultTopIndividuals{}, fmt.Errorf("CalcTopIndividualAllGenerations | Generation cannot be nil")
	}
	if len(generations) < 1 {
		return ResultTopIndividuals{}, fmt.Errorf("CalcTopIndividualAllGenerations | Generation cannot be empty")
	}
	if individualKind < 0 {
		individualKind = 0
	}
	if individualKind > 1 {
		individualKind = 1
	}

	var topIndividual = ResultTopIndividuals{}
	switch fitnessComparator {
	case FitnessMoreIsBetter:
		topIndividual = ResultTopIndividuals{
			Generation: nil,
			Result:     &Individual{TotalFitness: math.MinInt64},
		}
	case FitnessLessIsBetter:
		topIndividual = ResultTopIndividuals{
			Generation: nil,
			Result:     &Individual{TotalFitness: math.MaxInt64},
		}
	}

	if individualKind == IndividualAntagonist {
		for i := range generations {
			individual, err := CalcTopIndividual(generations[i].Antagonists, fitnessComparator)
			if err != nil {
				return ResultTopIndividuals{}, err
			}
			switch fitnessComparator {
			case FitnessMoreIsBetter:
				if individual.TotalFitness > topIndividual.Result.TotalFitness {
					topIndividual.Result = individual
					topIndividual.Generation = generations[i]
				}
			case FitnessLessIsBetter:
				if individual.TotalFitness < topIndividual.Result.TotalFitness {
					topIndividual.Result = individual
					topIndividual.Generation = generations[i]
				}
			default:
				if individual.TotalFitness > topIndividual.Result.TotalFitness {
					topIndividual.Result = individual
					topIndividual.Generation = generations[i]
				}
			}

		}

	} else {
		for i := range generations {
			individual, err := CalcTopIndividual(generations[i].Protagonists, fitnessComparator)
			if err != nil {
				return ResultTopIndividuals{}, err
			}
			switch fitnessComparator {
			case FitnessMoreIsBetter:
				if individual.TotalFitness > topIndividual.Result.TotalFitness {
					topIndividual.Result = individual
					topIndividual.Generation = generations[i]
				}
			case FitnessLessIsBetter:
				if individual.TotalFitness < topIndividual.Result.TotalFitness {
					topIndividual.Result = individual
					topIndividual.Generation = generations[i]
				}
			default:
				if individual.TotalFitness > topIndividual.Result.TotalFitness {
					topIndividual.Result = individual
					topIndividual.Generation = generations[i]
				}
			}
		}
	}

	return topIndividual, nil
}

func CalcGenerationalFitnessAverage(generations []*Generation,
	individualKind int) ([]ResultInfo1DAveragesPerGeneration, error) {
	if generations == nil {
		return nil, fmt.Errorf("CalcGenerationalFitnessAverage | Generation cannot be nil")
	}
	if len(generations) < 1 {
		return nil, fmt.Errorf("CalcGenerationalFitnessAverage | Generation cannot be empty")
	}
	if individualKind < 0 {
		individualKind = 0
	}
	if individualKind > 1 {
		individualKind = 1
	}

	result := make([]ResultInfo1DAveragesPerGeneration, len(generations))
	if individualKind == IndividualAntagonist {
		for i := range generations {
			average := CalculateAverage(generations[i].Antagonists)
			result[i] = ResultInfo1DAveragesPerGeneration{
				Result:     average,
				Generation: generations[i],
			}
		}

	} else {
		for i := range generations {
			average := CalculateAverage(generations[i].Protagonists)
			result[i] = ResultInfo1DAveragesPerGeneration{
				Result:     average,
				Generation: generations[i],
			}
		}
	}
	return result, nil
}

func CalcGenerationalFitnessCum(generations []*Generation,
	individualKind int) ([]ResultInfo1DAveragesPerGeneration, error) {
	if generations == nil {
		return nil, fmt.Errorf("CalcGenerationalFitnessCum | Generation cannot be nil")
	}
	if len(generations) < 1 {
		return nil, fmt.Errorf("CalcGenerationalFitnessCum | Generation cannot be empty")
	}
	if individualKind < 0 {
		individualKind = 0
	}
	if individualKind > 1 {
		individualKind = 1
	}

	result := make([]ResultInfo1DAveragesPerGeneration, len(generations))
	if individualKind == IndividualAntagonist {
		for i := range generations {
			cum := CalculateCum(generations[i].Antagonists)
			result[i] = ResultInfo1DAveragesPerGeneration{
				Result:     cum,
				Generation: generations[i],
			}
		}

	} else {
		for i := range generations {
			cum := CalculateCum(generations[i].Protagonists)
			result[i] = ResultInfo1DAveragesPerGeneration{
				Result:     cum,
				Generation: generations[i],
			}
		}
	}
	return result, nil
}

func CalcTopNIndividualAllGenerations(generations []*Generation, individualKind int,
	fitnessComparator int,
	topN int) ([]ResultInfo2DPerGeneration,
	error) {
	if generations == nil {
		return nil, fmt.Errorf("CalcTopIndividualAllGenerations | Generation cannot be nil")
	}
	if len(generations) < 1 {
		return nil, fmt.Errorf("CalcTopIndividualAllGenerations | Generation cannot be empty")
	}
	if individualKind < 0 {
		individualKind = 0
	}
	if individualKind > 1 {
		individualKind = 1
	}

	// Handle Top N
	if topN < 1 {
		topN = 1
	} else if topN >= len(generations[0].Antagonists) {
		topN = len(generations[0].Antagonists)
	}

	resultInfo2DPerGenerations := make([]ResultInfo2DPerGeneration, len(generations))

	if individualKind == IndividualAntagonist {
		for i := range generations {
			sortIndividuals := SortIndividuals(generations[i].Antagonists, fitnessComparator)
			resultInfo2DPerGenerations[i].Generation = generations[i]
			resultInfo2DPerGenerations[i].Result = sortIndividuals[:topN]
		}
	} else {
		for i := range generations {
			sortIndividuals := SortIndividuals(generations[i].Protagonists, fitnessComparator)
			resultInfo2DPerGenerations[i].Generation = generations[i]
			resultInfo2DPerGenerations[i].Result = sortIndividuals[:topN]
		}
	}

	return resultInfo2DPerGenerations, nil
}

// SortIndividuals returns the Top N-1 individuals. In this application less is more,
// so they are sorted in ascending order, with smaller indices representing better individuals.
// It is for the user to specify the Kind of individual to pass in be it antagonist or protagonist.
// TODO CHECK NULL
func SortIndividuals(individuals []*Individual, fitnessComparator int) []*Individual {
	switch fitnessComparator {
	case FitnessMoreIsBetter:
		sort.Slice(individuals, func(i, j int) bool {
			return individuals[i].TotalFitness > individuals[j].TotalFitness
		})
	case FitnessLessIsBetter:
		sort.Slice(individuals, func(i, j int) bool {
			return individuals[i].TotalFitness < individuals[j].TotalFitness
		})
	default:
		sort.Slice(individuals, func(i, j int) bool {
			return individuals[i].TotalFitness > individuals[j].TotalFitness
		})
	}
	return individuals
}

func CalculateAverage(individuals []*Individual) float64 {
	sum := 0.0
	for i := range individuals {
		sum += individuals[i].TotalFitness
	}
	return float64(sum / float64(len(individuals)))
}

func CalculateCum(individuals []*Individual) float64 {
	sum := 0.0
	for i := range individuals {
		sum += individuals[i].TotalFitness
	}
	return float64(sum)
}

func (e *EvolutionResult) Analyze(generations []*Generation, fitnessStrategy int, topN int) (EvolutionSummary, error) {
	switch fitnessStrategy {
	case FitnessRatio:
		err := e.CalculateTop(generations, FitnessMoreIsBetter, topN)
		if err != nil {
			return EvolutionSummary{}, err
		}
		break
	case FitnessProtagonistThresholdTally:
		err := e.CalculateTop(generations, FitnessLessIsBetter, topN)
		if err != nil {
			return EvolutionSummary{}, err
		}
		break
	default:
		log.Print("Unknown Fitness Strategy")
	}

	generationalAntagonistFitnessCum, err := CalcGenerationalFitnessCum(generations, IndividualAntagonist)
	if err != nil {
		return EvolutionSummary{}, nil
	}
	e.AntagonistCumAcrossGenerations = generationalAntagonistFitnessCum

	generationalProtagonistFitnessCum, err := CalcGenerationalFitnessCum(generations, IndividualProtagonist)
	if err != nil {
		return EvolutionSummary{}, nil
	}
	e.ProtagonistCumAcrossGenerations = generationalProtagonistFitnessCum

	averageAntagonists, err := CalcGenerationalFitnessAverage(generations, IndividualAntagonist)
	if err != nil {
		return EvolutionSummary{}, nil
	}
	e.AntagonistAverageAcrossGenerations = averageAntagonists

	averageProtagonists, err := CalcGenerationalFitnessAverage(generations, IndividualProtagonist)
	if err != nil {
		return EvolutionSummary{}, nil
	}
	e.ProtagonistAverageAcrossGenerations = averageProtagonists
	return EvolutionSummary{}, nil
}

// CalculateTop calculates all the relevant top individual metrics e.g Top in Generationand Top across generations
func (e *EvolutionResult) CalculateTop(generations []*Generation, fitnessStrategy int, topN int) error {
	topAntagonistAllGenerations, err := CalcTopIndividualAllGenerations(generations, IndividualAntagonist, fitnessStrategy)
	if err != nil {
		return err
	}
	e.TopAntagonist = topAntagonistAllGenerations

	topProtagonistAllGenerations, err := CalcTopIndividualAllGenerations(generations, IndividualProtagonist, fitnessStrategy)
	if err != nil {
		return err
	}
	e.TopProtagonist = topProtagonistAllGenerations

	topNAntagonistsAllGenerations, err := CalcTopNIndividualAllGenerations(generations, IndividualAntagonist, fitnessStrategy, topN)
	if err != nil {
		return err
	}
	e.TopNAntagonistsPerGeneration = topNAntagonistsAllGenerations
	topNProtagonistsAllGenerations, err := CalcTopNIndividualAllGenerations(generations, IndividualProtagonist, fitnessStrategy, topN)
	if err != nil {
		return err
	}
	e.TopNProtagonistsPerGeneration = topNProtagonistsAllGenerations
	return nil
}

type EvolutionSummary struct{}

func (e *EvolutionResult) PrintAverageGenerationSummary() (strings.Builder, error) {
	if e.ProtagonistAverageAcrossGenerations == nil {
		return strings.Builder{},
			fmt.Errorf("PrintAverageGenerationSummary | cannot format as protagonist average field is nil | Run" +
				" analyze")
	}
	if e.AntagonistAverageAcrossGenerations == nil {
		return strings.Builder{},
			fmt.Errorf("PrintAverageGenerationSummary | cannot format as antagonist average field is nil | Run" +
				" analyze")
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("" +
		"####################################### AVERAGE ANTAGONISTS VS PROTAGONISTS PER GENERATION" +
		" #####################################################\n\n"))
	sb.WriteString("ANT | PRO\n")
	for i := range e.AntagonistAverageAcrossGenerations {
		res := e.AntagonistAverageAcrossGenerations[i].Result
		resPro := e.ProtagonistAverageAcrossGenerations[i].Result
		float := strconv.FormatFloat(res, 'g', 03, 64)
		floatPro := strconv.FormatFloat(resPro, 'g', 03, 64)

		gen := strconv.Itoa(i)
		sb.WriteString("gen" + gen + ": " + float + " | " + floatPro + "\n")
	}
	sb.WriteString("\n")
	return sb, nil
}

func (e *EvolutionResult) PrintCumGenerationSummary() (strings.Builder, error) {
	if e.ProtagonistCumAcrossGenerations == nil {
		return strings.Builder{},
			fmt.Errorf("PrintCumGenerationSummary | cannot format as protagonist average field is nil | Run" +
				" analyze")
	}
	if e.AntagonistCumAcrossGenerations == nil {
		return strings.Builder{},
			fmt.Errorf("PrintCumGenerationSummary | cannot format as antagonist average field is nil | Run" +
				" analyze")
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("" +
		"####################################### CUMULATIVE ANTAGONISTS VS PROTAGONISTS PER GENERATION" +
		" #####################################################\n\n"))
	sb.WriteString("ANT | PRO\n")
	for i := range e.AntagonistCumAcrossGenerations {
		res := e.AntagonistCumAcrossGenerations[i].Result
		resPro := e.ProtagonistCumAcrossGenerations[i].Result
		float := strconv.FormatFloat(res, 'g', 03, 64)
		floatPro := strconv.FormatFloat(resPro, 'g', 03, 64)

		gen := strconv.Itoa(i)
		sb.WriteString("gen" + gen + ": " + float + " | " + floatPro + "\n")
	}
	sb.WriteString("\n")
	return sb, nil
}

type ResultTopIndividuals struct {
	Generation *Generation
	Result     *Individual
	Tree       string
}

func (e *EvolutionResult) PrintTopIndividualSummary(kind int) (strings.Builder, error) {
	sb := strings.Builder{}
	var name string

	if kind == IndividualProtagonist {
		if e.TopProtagonist.Result == nil {
			return strings.Builder{},
				fmt.Errorf("PrintTopIndividualSummary | cannot format as field is nil | Run analyze")
		}

		name = "PROTAGONIST"
		sb.WriteString(fmt.Sprintf("############### TOP %s IN ALL GENERATIONS"+" #######################\n", name))
		sb.WriteString(fmt.Sprintf("ID: %s\n", e.TopProtagonist.Result.Id))
		sb.WriteString(fmt.Sprintf("GENERATION:  %s\n", e.TopProtagonist.Generation.GenerationID))
		sb.WriteString(fmt.Sprintf("AGE:  %d\n", e.TopProtagonist.Result.Age))
		sb.WriteString(fmt.Sprintf("FITNESS:  %f\n", e.TopProtagonist.Result.TotalFitness))
		sb.WriteString(fmt.Sprintf("BIRTH GEN:  %d\n", e.TopProtagonist.Result.BirthGen))

		strategiesSummary := FormatStrategiesTotal(e.TopProtagonist.Result.Strategy)
		sb.WriteString(fmt.Sprintf("Strategy Summary:\n%s\n", strategiesSummary.String()))

		sb.WriteString("Tree Shape:\n")
		treeBuilder := e.TopProtagonist.Result.Program.T.ToString()
		sb.WriteString(treeBuilder.String())

		mathematicalString, err := e.TopProtagonist.Result.Program.T.ToMathematicalString()
		if err != nil {
			return strings.Builder{}, err
		}
		sb.WriteString(fmt.Sprintf("Mathematical Expression: %s\n", mathematicalString))

		strategiesList := FormatStrategiesList(e.TopProtagonist.Result.Strategy)
		sb.WriteString(fmt.Sprintf("Strategy Summary:\n  %s\n", strategiesList.String()))
	} else if kind == IndividualAntagonist {
		if e.TopAntagonist.Result == nil {
			return strings.Builder{},
				fmt.Errorf("PrintTopIndividualSummary | cannot format as field is nil | Run analyze")
		}

		name = "ANTAGONIST"
		sb.WriteString(fmt.Sprintf("############### TOP %s IN ALL GENERATIONS"+" #######################\n", name))
		sb.WriteString(fmt.Sprintf("ID: %s\n", e.TopAntagonist.Result.Id))
		sb.WriteString(fmt.Sprintf("GENERATION:  %s\n", e.TopAntagonist.Generation.GenerationID))
		sb.WriteString(fmt.Sprintf("AGE:  %d\n", e.TopAntagonist.Result.Age))
		sb.WriteString(fmt.Sprintf("FITNESS:  %f\n", e.TopAntagonist.Result.TotalFitness))
		sb.WriteString(fmt.Sprintf("BIRTH GEN:  %d\n", e.TopAntagonist.Result.BirthGen))

		strategiesSummary := FormatStrategiesTotal(e.TopAntagonist.Result.Strategy)
		sb.WriteString(fmt.Sprintf("Strategy Summary:\n%s\n", strategiesSummary.String()))

		sb.WriteString("Tree Shape:\n")
		treeBuilder := e.TopAntagonist.Result.Program.T.ToString()
		sb.WriteString(treeBuilder.String())

		mathematicalString, err := e.TopAntagonist.Result.Program.T.ToMathematicalString()
		if err != nil {
			return strings.Builder{}, err
		}
		sb.WriteString(fmt.Sprintf("Mathematical Expression: %s\n", mathematicalString))

		strategiesList := FormatStrategiesList(e.TopAntagonist.Result.Strategy)
		sb.WriteString(fmt.Sprintf("Strategy Summary:\n%s\n", strategiesList.String()))
	}
	return sb, nil
}

func (e *EvolutionResult) PrintTopNInFinalGeneration(kind int) {

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
