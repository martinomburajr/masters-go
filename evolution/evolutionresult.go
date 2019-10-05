package evolution

import (
	"fmt"
	"math"
	"sort"
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

	SortedProtagonistsPerGeneration ResultInfo2DPerGeneration
	SortedAntagonistsPerGeneration  ResultInfo2DPerGeneration
}

type ResultInfo2DPerGeneration struct {
	generation *Generation
	result     []*Individual
}

type ResultInfo1DPerGeneration struct {
	generation *Generation
	result     []*Individual
}

type ResultTopIndividuals struct {
	generation *Generation
	result     *Individual
	tree       string
}

type ResultInfo1DAveragesPerGeneration struct {
	generation *Generation
	result     float64
}

func CalcTopIndividual(individuals []*Individual) (*Individual, error) {
	if individuals == nil {
		return nil, fmt.Errorf("CalcTopIndividual | Individuals cannot be nil")
	}
	if len(individuals) < 1 {
		return nil, fmt.Errorf("CalcTopIndividual | Individuals cannot be empty")
	}

	fitness := math.MaxInt64
	individual := &Individual{}
	for i := range individuals {
		if individuals[i].totalFitness < fitness {
			individual = individuals[i]
		}
	}
	return individual, nil
}

func CalcTopIndividualAllGenerations(generations []*Generation, individualKind int) (ResultTopIndividuals,
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

	topIndividual := ResultTopIndividuals{
		generation: nil,
		tree:       "",
		result:     &Individual{totalFitness: math.MaxInt64},
	}

	if individualKind == IndividualAntagonist {
		for i := range generations {
			individual, err := CalcTopIndividual(generations[i].Antagonists)
			if err != nil {
				return ResultTopIndividuals{}, err
			}
			if individual.totalFitness < topIndividual.result.totalFitness {
				topIndividual.result = individual
				topIndividual.generation = generations[i]
				topIndividual.tree = topIndividual.result.Program.T.ToString()
			}
		}

	} else {
		for i := range generations {
			individual, err := CalcTopIndividual(generations[i].Protagonists)
			if err != nil {
				return ResultTopIndividuals{}, err
			}
			if individual.totalFitness < topIndividual.result.totalFitness {
				topIndividual.result = individual
				topIndividual.generation = generations[i]
				topIndividual.tree = topIndividual.result.Program.T.ToString()
			}
		}
	}

	return topIndividual, nil
}

func CalcGenerationalFitnessAverage(generations []*Generation,
	individualKind int) ([]ResultInfo1DAveragesPerGeneration, error) {
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

	result := make([]ResultInfo1DAveragesPerGeneration, len(generations))
	if individualKind == IndividualAntagonist {
		for i := range generations {
			average := CalculateAverage(generations[i].Antagonists)
			result[i] = ResultInfo1DAveragesPerGeneration{
				result: average,
				generation: generations[i],
			}
		}

	} else {
		for i := range generations {
			average := CalculateAverage(generations[i].Protagonists)
			result[i] = ResultInfo1DAveragesPerGeneration{
				result: average,
				generation: generations[i],
			}
		}
	}
	return result, nil
}


func CalcTopNIndividualAllGenerations(generations []*Generation, individualKind int,
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
	} else if topN > len(generations) {
		topN = len(generations)
	}

	resultInfo2DPerGenerations := make([]ResultInfo2DPerGeneration, len(generations))

	if individualKind == IndividualAntagonist {
		for i := range generations {
			sortIndividuals := SortIndividuals(generations[i].Antagonists)
			resultInfo2DPerGenerations[i].generation = generations[i]
			resultInfo2DPerGenerations[i].result = sortIndividuals[:topN]
		}
	} else {
		for i := range generations {
			sortIndividuals := SortIndividuals(generations[i].Protagonists)
			resultInfo2DPerGenerations[i].generation = generations[i]
			resultInfo2DPerGenerations[i].result = sortIndividuals[:topN]
		}
	}

	return resultInfo2DPerGenerations, nil
}

// SortIndividuals returns the Top N-1 individuals. In this application less is more,
// so they are sorted in ascending order, with smaller indices representing better individuals.
// It is for the user to specify the kind of individual to pass in be it antagonist or protagonist.
func SortIndividuals(individuals []*Individual) []*Individual {
	sort.Slice(individuals, func(i, j int) bool {
		return individuals[i].totalFitness < individuals[j].totalFitness
	})
	return individuals
}

func CalculateAverage(individuals []*Individual) float64 {
	sum := 0
	for i := range individuals {
		sum += individuals[i].totalFitness
	}
	return float64(sum / len(individuals))
}

func (e *EvolutionResult) Analyze(generations []*Generation, topN int) (EvolutionSummary, error) {
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

	topNAntagonistsAllGenerations, err := CalcTopNIndividualAllGenerations(generations, IndividualAntagonist, topN)
	if err != nil {
		return EvolutionSummary{}, nil
	}
	e.TopNAntagonistsPerGeneration = topNAntagonistsAllGenerations
	topNProtagonistsAllGenerations, err := CalcTopNIndividualAllGenerations(generations, IndividualProtagonist, topN)
	if err != nil {
		return EvolutionSummary{}, nil
	}
	e.TopNProtagonistsPerGeneration = topNProtagonistsAllGenerations

	topAntagonistAllGenerations, err := CalcTopIndividualAllGenerations(generations, IndividualAntagonist)
	if err != nil {
		return EvolutionSummary{}, nil
	}
	e.TopAntagonist = topAntagonistAllGenerations

	topProtagonistAllGenerations, err := CalcTopIndividualAllGenerations(generations, IndividualProtagonist)
	if err != nil {
		return EvolutionSummary{}, nil
	}
	e.TopProtagonist = topProtagonistAllGenerations

	//if generations == nil {
	//	return EvolutionSummary{}, fmt.Errorf("generations cannot be nil")
	//}
	//if len(generations) < 1 {
	//	return EvolutionSummary{}, fmt.Errorf("generations cannot be empty")
	//}
	//
	//e.TopAntagonist = ResultTopIndividuals{result: &Individual{totalFitness: math.MaxInt64}}
	//e.TopProtagonist = ResultTopIndividuals{result: &Individual{totalFitness: math.MaxInt64}}
	//
	//e.TopAntagonistPerGeneration = ResultInfo1DPerGeneration{result: make([]*Individual, 0)}
	//e.TopProtagonistsPerGeneration = ResultInfo1DPerGeneration{result: make([]*Individual, 0)}
	//
	//e.SortedAntagonistsPerGeneration = ResultInfo2DPerGeneration{result: make([][]*Individual, len(generations))}
	//e.SortedProtagonistsPerGeneration = ResultInfo2DPerGeneration{result: make([][]*Individual, len(generations))}
	//
	//e.TopNAntagonistsPerGeneration = ResultInfo2DPerGeneration{result: make([][]*Individual, len(generations))}
	//e.TopNProtagonistsPerGeneration = ResultInfo2DPerGeneration{result: make([][]*Individual, len(generations))}
	//
	//e.AntagonistAverageAcrossGenerations = ResultInfo1DAveragesPerGeneration{result: make([]float64, 0)}
	//e.ProtagonistAverageAcrossGenerations = ResultInfo1DAveragesPerGeneration{result: make([]float64, 0)}
	//
	////wg := sync.WaitGroup{}
	////wg.Add(2)
	////go func() {
	////	defer wg.Done()
	//
	//for i := range generations {
	//	sortedAntagonistsInGeneration := SortIndividuals(generations[i].Antagonists)
	//	averageAntagonists := CalculateAverage(generations[i].Antagonists)
	//
	//	// Handle Top Individual
	//	if sortedAntagonistsInGeneration[0].totalFitness < e.TopAntagonist.result.totalFitness {
	//		e.TopAntagonist.result = sortedAntagonistsInGeneration[0]
	//		e.TopAntagonist.generation = generations[i]
	//	}
	//
	//	// Handle Averages
	//	e.AntagonistAverageAcrossGenerations.generation = generations[i]
	//	e.AntagonistAverageAcrossGenerations.result = append(e.AntagonistAverageAcrossGenerations.result, averageAntagonists)
	//
	//	// Handle Top Individual in Generation
	//	e.TopAntagonistPerGeneration.result = append(e.TopAntagonistPerGeneration.result, sortedAntagonistsInGeneration[0])
	//	e.TopAntagonistPerGeneration.generation = generations[i]
	//
	//	// Handle Top N
	//	if topN < 1 {
	//		topN = 1
	//	} else if topN > len(sortedAntagonistsInGeneration) {
	//		topN = len(sortedAntagonistsInGeneration)
	//	}
	//
	//	e.TopNAntagonistsPerGeneration.result[i] = make([]*Individual, 0)
	//	e.TopNAntagonistsPerGeneration.result[i] = append(e.TopNAntagonistsPerGeneration.result[i],
	//		sortedAntagonistsInGeneration[:topN]...)
	//	e.TopNAntagonistsPerGeneration.generation = generations[i]
	//
	//	// Handle Sorted Individuals
	//	//e.SortedAntagonistsPerGeneration.generation = generations[i]
	//	//e.SortedAntagonistsPerGeneration.result[i] = make([]*Individual, len(generations[i].Antagonists))
	//	//var individuals = e.SortedAntagonistsPerGeneration.result[i]
	//	//e.SortedAntagonistsPerGeneration.result[i] = append(individuals, sortedAntagonistsInGeneration...)
	//}
	//e.TopAntagonist.tree = e.TopAntagonist.result.Program.T.ToString()
	////}()
	//
	////go func() {
	////	defer wg.Done()
	//for i := range generations {
	//	sortedProtagonistsInGeneration := SortIndividuals(generations[i].Protagonists)
	//	averageProtagonists := CalculateAverage(generations[i].Protagonists)
	//
	//	// Handle Top Individual
	//	if sortedProtagonistsInGeneration[0].totalFitness < e.TopProtagonist.result.totalFitness {
	//		e.TopProtagonist.result = sortedProtagonistsInGeneration[0]
	//		e.TopProtagonist.generation = generations[i]
	//	}
	//
	//	// Handle Averages
	//	e.ProtagonistAverageAcrossGenerations.generation = generations[i]
	//	e.ProtagonistAverageAcrossGenerations.result = append(e.ProtagonistAverageAcrossGenerations.result, averageProtagonists)
	//
	//	// Handle Top Individual in Generation
	//	e.TopProtagonistsPerGeneration.result = append(e.TopProtagonistsPerGeneration.result, sortedProtagonistsInGeneration[0])
	//	e.TopProtagonistsPerGeneration.generation = generations[i]
	//
	//	// Handle Top N
	//	if topN < 1 {
	//		topN = 1
	//	} else if topN > len(sortedProtagonistsInGeneration) {
	//		topN = len(sortedProtagonistsInGeneration)
	//	}
	//	e.TopNProtagonistsPerGeneration.result[i] = make([]*Individual, 0)
	//	e.TopNProtagonistsPerGeneration.result[i] = append(e.TopNProtagonistsPerGeneration.result[i],
	//		sortedProtagonistsInGeneration[:topN]...)
	//	e.TopNProtagonistsPerGeneration.generation = generations[i]
	//
	//	// Handle Sorted Individuals
	//	//e.SortedProtagonistsPerGeneration.generation = generations[i]
	//	//e.SortedProtagonistsPerGeneration.result[i] = make([]*Individual, len(generations[i].Protagonists))
	//	//var individuals = e.SortedAntagonistsPerGeneration.result[i]
	//	//e.SortedProtagonistsPerGeneration.result[i] = append(individuals, sortedProtagonistsInGeneration...)
	//}
	//e.TopProtagonist.tree = e.TopProtagonist.result.Program.T.ToString()
	//}()
	//wg.Done()
	return EvolutionSummary{}, nil
}

type EvolutionSummary struct{}

func (e *EvolutionSummary) PrintSummary() *Program {
	return nil
}
