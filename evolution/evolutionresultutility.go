package evolution

import (
	"fmt"
	"math"
	"sort"
)

// GetTopIndividualInAllGenerations returns the best protagonist and antagonist in the entire evolutionary process
func GetTopIndividualInAllGenerations(sortedGenerations []*Generation, isMoreFitnessBetter bool) (topAntagonist *Individual, topProtagonist *Individual, err error) {
	if sortedGenerations == nil {
		return nil, nil, fmt.Errorf("GetGenerationalFitnessAverage | Generation cannot be nil")
	}
	if len(sortedGenerations) < 1 {
		return nil, nil, fmt.Errorf("GetGenerationalFitnessAverage | Generation cannot be empty")
	}

	if isMoreFitnessBetter {
		topAntagonist = &Individual{TotalFitness: math.MinInt64}
		topProtagonist = &Individual{TotalFitness: math.MinInt64}
		for i := 0; i < len(sortedGenerations); i++ {
			// This ensures it picks more recent individuals
			if sortedGenerations[i].Antagonists[0].TotalFitness >= topAntagonist.TotalFitness {
				topAntagonist = sortedGenerations[i].Antagonists[0]
			}
			if sortedGenerations[i].Protagonists[0].TotalFitness >= topProtagonist.TotalFitness {
				topProtagonist = sortedGenerations[i].Protagonists[0]
			}
		}
		return topAntagonist, topProtagonist, nil
	} else {
		topAntagonist = &Individual{TotalFitness: math.MaxInt64}
		topProtagonist = &Individual{TotalFitness: math.MaxInt64}
		for i := 0; i < len(sortedGenerations); i++ {
			// This ensures it picks more recent individuals
			if sortedGenerations[i].Antagonists[0].TotalFitness >= topAntagonist.TotalFitness {
				topAntagonist = sortedGenerations[i].Antagonists[0]
			}
			if sortedGenerations[i].Protagonists[0].TotalFitness >= topProtagonist.TotalFitness {
				topProtagonist = sortedGenerations[i].Protagonists[0]
			}
		}
	}

	return topAntagonist, topProtagonist, nil
}

// GetGenerationalFitnessAverage returns the average for antagonists and protagonists of individual for each generation
func GetGenerationalFitnessAverage(sortedGenerations []*Generation) ([]generationalCoevolutionaryAverages, error) {
	if sortedGenerations == nil {
		return nil, fmt.Errorf("GetGenerationalFitnessAverage | Generation cannot be nil")
	}
	if len(sortedGenerations) < 1 {
		return nil, fmt.Errorf("GetGenerationalFitnessAverage | Generation cannot be empty")
	}

	result := make([]generationalCoevolutionaryAverages, len(sortedGenerations))
	for i := range sortedGenerations {
		antagonistAverage, err := CalculateAverage(sortedGenerations[i].Antagonists)
		if err != nil {
			return nil, err
		}
		protagonistAverage, err := CalculateAverage(sortedGenerations[i].Protagonists)
		if err != nil {
			return nil, err
		}
		result[i] = generationalCoevolutionaryAverages{
			AntagonistResult:  antagonistAverage,
			ProtagonistResult: protagonistAverage,
			Generation:        sortedGenerations[i],
		}
	}

	return result, nil
}

// GetTopNIndividualsPerGeneration returns the top n individuals in each generation
func GetTopNIndividualsPerGeneration(sortedGenerations []*Generation, individualKind int, topN int) ([]multiIndividualsPerGeneration,
	error) {
	if sortedGenerations == nil {
		return nil, fmt.Errorf("CalcNthPlaceIndividualAllGenerations | Generation cannot be nil")
	}
	if len(sortedGenerations) < 1 {
		return nil, fmt.Errorf("CalcNthPlaceIndividualAllGenerations | Generation cannot be empty")
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
	} else if topN >= len(sortedGenerations[0].Antagonists) {
		topN = len(sortedGenerations[0].Antagonists)
	}

	resultInfo2DPerGenerations := make([]multiIndividualsPerGeneration, len(sortedGenerations))
	if individualKind == IndividualAntagonist {
		for i := range sortedGenerations {
			resultInfo2DPerGenerations[i].Generation = sortedGenerations[i]
			resultInfo2DPerGenerations[i].Individuals = sortedGenerations[i].Antagonists[:topN]
		}
	} else {
		for i := range sortedGenerations {
			resultInfo2DPerGenerations[i].Generation = sortedGenerations[i]
			resultInfo2DPerGenerations[i].Individuals = sortedGenerations[i].Protagonists[:topN]
		}
	}

	return resultInfo2DPerGenerations, nil
}

// GetTopNIndividualInGenerationX calculates the top N individuals in a specified generation
func GetTopNIndividualInGenerationX(sortedGenerations []*Generation, individualKind int,
	isMoreFitnessBetter bool, topN int, generationN int) (multiIndividualsPerGeneration,
	error) {
	if sortedGenerations == nil {
		return multiIndividualsPerGeneration{}, fmt.Errorf("CalcNthPlaceIndividualAllGenerations | Generation cannot be nil")
	}
	if len(sortedGenerations) < 1 {
		return multiIndividualsPerGeneration{}, fmt.Errorf("CalcNthPlaceIndividualAllGenerations | Generation cannot be empty")
	}
	if individualKind < 0 {
		individualKind = 0
	}
	if individualKind > 1 {
		individualKind = 1
	}
	if generationN >= len(sortedGenerations) {
		generationN = len(sortedGenerations) - 1
	}
	if generationN < 0 {
		generationN = 0
	}

	// Handle Top N
	if topN < 1 {
		topN = 1
	} else if topN >= len(sortedGenerations[0].Antagonists) {
		topN = len(sortedGenerations[0].Antagonists)
	}

	resultInfo2DPerGenerations := multiIndividualsPerGeneration{}
	resultInfo2DPerGenerations.Generation = sortedGenerations[generationN]
	if individualKind == IndividualAntagonist {
		resultInfo2DPerGenerations.Individuals = sortedGenerations[generationN].Antagonists[:topN]
	} else {
		resultInfo2DPerGenerations.Individuals = sortedGenerations[generationN].Protagonists[:topN]
	}

	return resultInfo2DPerGenerations, nil
}

// GetNthPlaceIndividual returns an individual in the nth place. N must be an index and not an actual position e.g.
// 0 is the first individual
func GetNthPlaceIndividual(sortedIndividuals []*Individual, n int) (*Individual, error) {
	if sortedIndividuals == nil {
		return nil, fmt.Errorf("GetNthPlaceIndividual | Individuals cannot be nil")
	}
	if len(sortedIndividuals) < 1 {
		return nil, fmt.Errorf("GetNthPlaceIndividual | Individuals cannot be empty")
	}
	if n < 0 {
		n = 0
	}
	if n >= len(sortedIndividuals) {
		n = len(sortedIndividuals) - 1
	}

	return sortedIndividuals[n], nil
}

// SortIndividuals returns the Top N-1 individuals. In this application less is more,
// so they are sorted in ascending order, with smaller indices representing better individuals.
// It is for the user to specify the Kind of individual to pass in be it antagonist or protagonist.
func SortIndividuals(individuals []*Individual, isMoreFitnessBetter bool) ([]*Individual, error) {
	if individuals == nil {
		return nil, fmt.Errorf("SortIndividuals | individuals cannot be nil")
	}
	if len(individuals) < 1 {
		return nil, fmt.Errorf("SortIndividuals | individuals cannot be empty")
	}

	switch isMoreFitnessBetter {
	case true:
		sort.Slice(individuals, func(i, j int) bool {
			return individuals[i].TotalFitness > individuals[j].TotalFitness
		})
	case false:
		sort.Slice(individuals, func(i, j int) bool {
			return individuals[i].TotalFitness < individuals[j].TotalFitness
		})
	default:
		// Default to More is better
		sort.Slice(individuals, func(i, j int) bool {
			return individuals[i].TotalFitness > individuals[j].TotalFitness
		})
	}
	return individuals, nil
}

// CalculateAverage averages the fitness values for each individual
func CalculateAverage(individuals []*Individual) (float64, error) {
	if individuals == nil {
		return -1, fmt.Errorf("SortIndividuals | individuals cannot be nil")
	}
	if len(individuals) < 1 {
		return -1, fmt.Errorf("SortIndividuals | individuals cannot be empty")
	}

	sum := 0.0
	for i := range individuals {
		sum += individuals[i].TotalFitness
	}
	return float64(sum / float64(len(individuals))), nil
}

// CalculateCumulative accumulates all the averaged fitness values each individual has.
func CalculateCumulative(individuals []*Individual) (float64, error) {
	if individuals == nil {
		return -1, fmt.Errorf("SortIndividuals | individuals cannot be nil")
	}
	if len(individuals) < 1 {
		return -1, fmt.Errorf("SortIndividuals | individuals cannot be empty")
	}

	sum := 0.0
	for i := range individuals {
		sum += individuals[i].TotalFitness
	}

	return float64(sum), nil
}

// SortGenerationsThoroughly sorts each kind of individual in each generation for every generation.
// This allows for easy querying in later phases.
func SortGenerationsThoroughly(generations []*Generation, isMoreFitnessBetter bool) ([]*Generation, error) {
	if generations == nil {
		return nil, fmt.Errorf("SortGenerationsThoroughly | generations cannot be nil")
	}
	if len(generations) < 1 {
		return nil, fmt.Errorf("SortGenerationsThoroughly | generations cannot be empty")
	}

	sortedGenerations := make([]*Generation, len(generations))
	for i := 0; i < len(generations); i++ {
		sortedGenerations[i] = generations[i]
		sortedAntagonists, err := SortIndividuals(generations[i].Antagonists, isMoreFitnessBetter)
		if err != nil {
			return nil, err
		}
		sortedProtagonists, err := SortIndividuals(generations[i].Protagonists, isMoreFitnessBetter)
		if err != nil {
			return nil, err
		}
		sortedGenerations[i].Protagonists = sortedProtagonists
		sortedGenerations[i].Antagonists = sortedAntagonists
	}
	return sortedGenerations, nil
}