package evolution

import (
	"fmt"
	"sort"
)

type EvolutionResult struct{
	hasBeenAnalyzed bool
	TopProtagonist  ResultTopIndividuals
	TopAntagonist   ResultTopIndividuals

	TopProtagonistsPerGeneration        ResultInfo1DPerGeneration
	TopAntagonistPerGeneration          ResultInfo1DPerGeneration

	TopNProtagonistsPerGeneration       ResultInfo2DPerGeneration
	TopNAntagonistsPerGeneration        ResultInfo2DPerGeneration

	ProtagonistAverageAcrossGenerations ResultInfo1DAveragesPerGeneration
	AntagonistAverageAcrossGenerations  ResultInfo1DAveragesPerGeneration

	SortedProtagonistsPerGeneration     ResultInfo2DPerGeneration
	SortedAntagonistsPerGeneration      ResultInfo2DPerGeneration
}

type ResultInfo2DPerGeneration struct {
	generation *Generation
	result [][]Individual
}

type ResultInfo1DPerGeneration struct {
	generation *Generation
	result []Individual
}

type ResultTopIndividuals struct {
	generation *Generation
	result Individual
	tree string
}

type ResultInfo1DAveragesPerGeneration struct {
	generation *Generation
	result []float64
}

func (e *EvolutionResult) Analyze(generations []*Generation, topN int) (EvolutionSummary, error) {
	if generations == nil {
		return EvolutionSummary{}, fmt.Errorf("generations cannot be nil")
	}
	if len(generations) < 1 {
		return EvolutionSummary{}, fmt.Errorf("generations cannot be empty")
	}

	e.TopAntagonist = ResultTopIndividuals{result: Individual{}}
	e.TopProtagonist = ResultTopIndividuals{}

	e.TopAntagonistPerGeneration = ResultInfo1DPerGeneration{result: make([]Individual, 0)}
	e.TopProtagonistsPerGeneration = ResultInfo1DPerGeneration{result: make([]Individual, 0)}

	e.SortedAntagonistsPerGeneration = ResultInfo2DPerGeneration{result:  make([][]Individual, 0)}
	e.SortedProtagonistsPerGeneration = ResultInfo2DPerGeneration{result:  make([][]Individual, 0)}

	e.TopNAntagonistsPerGeneration = ResultInfo2DPerGeneration{result: make([][]Individual, 0)}
	e.TopNProtagonistsPerGeneration = ResultInfo2DPerGeneration{result: make([][]Individual, 0)}

	e.AntagonistAverageAcrossGenerations = ResultInfo1DAveragesPerGeneration{ result: make([]float64, 0)}
	e.ProtagonistAverageAcrossGenerations = ResultInfo1DAveragesPerGeneration{ result: make([]float64, 0)}

	//wg := sync.WaitGroup{}
	//wg.Add(2)
	//go func() {
	//	defer wg.Done()

		for i := range generations {
			sortedAntagonistsInGeneration := SortIndividuals(generations[i].Antagonists)
			averageAntagonists := CalculateAverage(generations[i].Antagonists)

			// Handle Top Individual
			if sortedAntagonistsInGeneration[0].totalFitness < e.TopAntagonist.result.totalFitness {
				e.TopAntagonist.result = sortedAntagonistsInGeneration[0]
				e.TopAntagonist.generation = generations[i]
			}

			// Handle Averages
			e.AntagonistAverageAcrossGenerations.generation = generations[i]
			e.AntagonistAverageAcrossGenerations.result = append(e.AntagonistAverageAcrossGenerations.result, averageAntagonists)

			// Handle Top Individual in Generation
			e.TopAntagonistPerGeneration.result = append(e.TopAntagonistPerGeneration.result, sortedAntagonistsInGeneration[0])
			e.TopAntagonistPerGeneration.generation = generations[i]

			// Handle Top N
			if topN < 1 {
				topN = 1
			} else if topN > len(sortedAntagonistsInGeneration) {
				topN = len(sortedAntagonistsInGeneration)
			}

			e.TopNAntagonistsPerGeneration.result[i] = sortedAntagonistsInGeneration[:topN]
			e.TopNAntagonistsPerGeneration.generation = generations[i]

			// Handle Sorted Individuals
			e.SortedAntagonistsPerGeneration.generation = generations[i]
			e.SortedAntagonistsPerGeneration.result[i] = make([]Individual, len(generations[i].Antagonists))
			var individuals = e.SortedAntagonistsPerGeneration.result[i]
			e.SortedAntagonistsPerGeneration.result[i] = append(individuals, sortedAntagonistsInGeneration...)
		}
		e.TopAntagonist.tree = e.TopAntagonist.result.Program.T.ToString()
	//}()

	//go func() {
	//	defer wg.Done()
		for i := range generations {
			sortedProtagonistsInGeneration := SortIndividuals(generations[i].Protagonists)
			averageProtagonists := CalculateAverage(generations[i].Protagonists)

			// Handle Top Individual
			if sortedProtagonistsInGeneration[0].totalFitness < e.TopProtagonist.result.totalFitness {
				e.TopProtagonist.result = sortedProtagonistsInGeneration[0]
				e.TopProtagonist.generation = generations[i]
			}

			// Handle Averages
			e.ProtagonistAverageAcrossGenerations.generation = generations[i]
			e.ProtagonistAverageAcrossGenerations.result = append(e.ProtagonistAverageAcrossGenerations.result, averageProtagonists)

			// Handle Top Individual in Generation
			e.TopProtagonistsPerGeneration.result = append(e.TopProtagonistsPerGeneration.result, sortedProtagonistsInGeneration[0])
			e.TopProtagonistsPerGeneration.generation = generations[i]

			// Handle Top N
			if topN < 1 {
				topN = 1
			} else if topN > len(sortedProtagonistsInGeneration) {
				topN = len(sortedProtagonistsInGeneration)
			}
			e.TopNProtagonistsPerGeneration.result[i] = sortedProtagonistsInGeneration[:topN]
			e.TopNProtagonistsPerGeneration.generation = generations[i]

			// Handle Sorted Individuals
			e.SortedProtagonistsPerGeneration.generation = generations[i]
			e.SortedProtagonistsPerGeneration.result[i] = make([]Individual, len(generations[i].Protagonists))
			var individuals = e.SortedAntagonistsPerGeneration.result[i]
			e.SortedProtagonistsPerGeneration.result[i] = append(individuals, sortedProtagonistsInGeneration...)
		}
		e.TopProtagonist.tree = e.TopProtagonist.result.Program.T.ToString()
	//}()
	//wg.Done()
	return EvolutionSummary{}, nil
}


type EvolutionSummary struct{}

func (e *EvolutionSummary) PrintSummary() *Program {
	return nil
}

// SortIndividuals returns the Top N-1 individuals. In this application less is more,
// so they are sorted in ascending order, with smaller indices representing better individuals.
// It is for the user to specify the kind of individual to pass in be it antagonist or protagonist.
func SortIndividuals(individuals []Individual) []Individual {
	sort.Slice(individuals, func(i, j int) bool {
		return individuals[i].totalFitness < individuals[j].totalFitness
	})
	return individuals
}

func CalculateAverage(individuals []Individual) float64 {
	sum := 0
	for i := range individuals {
		sum += individuals[i].totalFitness
	}
	return float64(sum/len(individuals))
}
