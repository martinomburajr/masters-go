package evolution

import (
	"fmt"
	"math"
	"sort"
)

// GetTopIndividualInRun returns the best protagonist and antagonist in the entire evolutionary process
func GetTopIndividualInRun(sortedGenerations []*Generation, isMoreFitnessBetter bool) (topAntagonist *Individual, topProtagonist *Individual, err error) {
	if sortedGenerations == nil {
		return nil, nil, fmt.Errorf("GetGenerationalFitnessAverage | Generation cannot be nil")
	}
	if len(sortedGenerations) < 1 {
		return nil, nil, fmt.Errorf("GetGenerationalFitnessAverage | Generation cannot be empty")
	}

	if isMoreFitnessBetter {
		topAntagonist = &Individual{AverageFitness: math.MinInt64}
		topProtagonist = &Individual{AverageFitness: math.MinInt64}
		for i := 0; i < len(sortedGenerations); i++ {
			// This ensures it picks more recent individuals
			if sortedGenerations[i].Antagonists[0].AverageFitness >= topAntagonist.AverageFitness {
				topAntagonist = sortedGenerations[i].Antagonists[0]
			}
			if sortedGenerations[i].Protagonists[0].AverageFitness >= topProtagonist.AverageFitness {
				topProtagonist = sortedGenerations[i].Protagonists[0]
			}
		}
		return topAntagonist, topProtagonist, nil
	} else {
		topAntagonist = &Individual{AverageFitness: math.MaxInt64}
		topProtagonist = &Individual{AverageFitness: math.MaxInt64}
		for i := 0; i < len(sortedGenerations); i++ {
			// This ensures it picks more recent individuals
			if sortedGenerations[i].Antagonists[0].AverageFitness >= topAntagonist.AverageFitness {
				topAntagonist = sortedGenerations[i].Antagonists[0]
			}
			if sortedGenerations[i].Protagonists[0].AverageFitness >= topProtagonist.AverageFitness {
				topProtagonist = sortedGenerations[i].Protagonists[0]
			}
		}
	}

	return topAntagonist, topProtagonist, nil
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
			return individuals[i].AverageFitness > individuals[j].AverageFitness
		})
		break
	case false:
		sort.Slice(individuals, func(i, j int) bool {
			return individuals[i].AverageFitness < individuals[j].AverageFitness
		})
		break
	default:
		// Default to More is better
		sort.Slice(individuals, func(i, j int) bool {
			return individuals[i].AverageFitness > individuals[j].AverageFitness
		})
	}
	return individuals, nil
}

// SortIndividuals returns the Top N-1 individuals. In this application less is more,
// so they are sorted in ascending order, with smaller indices representing better individuals.
// It is for the user to specify the Kind of individual to pass in be it antagonist or protagonist.
func SortIndividualsByAvgDelta(individuals []*Individual, isMoreFitnessBetter bool) ([]*Individual, error) {
	if individuals == nil {
		return nil, fmt.Errorf("SortIndividuals | individuals cannot be nil")
	}
	if len(individuals) < 1 {
		return nil, fmt.Errorf("SortIndividuals | individuals cannot be empty")
	}

	switch isMoreFitnessBetter {
	case true:
		sort.Slice(individuals, func(i, j int) bool {
			return individuals[i].AverageDelta > individuals[j].AverageDelta
		})
	case false:
		sort.Slice(individuals, func(i, j int) bool {
			return individuals[i].AverageDelta < individuals[j].AverageDelta
		})
	default:
		// Default to More is better
		sort.Slice(individuals, func(i, j int) bool {
			return individuals[i].AverageDelta > individuals[j].AverageDelta
		})
	}
	return individuals, nil
}

// SortIndividuals returns the Top N-1 individuals. In this application less is more,
// so they are sorted in ascending order, with smaller indices representing better individuals.
// It is for the user to specify the Kind of individual to pass in be it antagonist or protagonist.
func SortIndividualsByDelta(individuals []*Individual, isMoreFitnessBetter bool) ([]*Individual, error) {
	if individuals == nil {
		return nil, fmt.Errorf("SortIndividuals | individuals cannot be nil")
	}
	if len(individuals) < 1 {
		return nil, fmt.Errorf("SortIndividuals | individuals cannot be empty")
	}

	switch isMoreFitnessBetter {
	case true:
		sort.Slice(individuals, func(i, j int) bool {
			return individuals[i].AverageDelta > individuals[j].AverageDelta
		})
	case false:
		sort.Slice(individuals, func(i, j int) bool {
			return individuals[i].AverageDelta < individuals[j].AverageDelta
		})
	default:
		// Default to More is better
		sort.Slice(individuals, func(i, j int) bool {
			return individuals[i].AverageDelta > individuals[j].AverageDelta
		})
	}
	return individuals, nil
}

// CalculateAverageFitnessAverage averages the fitness values for each individual
func CalculateAverageFitnessAverage(individuals []*Individual) (float64, error) {
	if individuals == nil {
		return -1, fmt.Errorf("SortIndividuals | individuals cannot be nil")
	}
	if len(individuals) < 1 {
		return -1, fmt.Errorf("SortIndividuals | individuals cannot be empty")
	}

	sum := 0.0
	for i := range individuals {
		sum += individuals[i].AverageFitness
	}
	return sum / float64(len(individuals)), nil
}

// CalculateAverage averages the fitness values for each individual
func CalculateAverage(items []float64) (float64, error) {
	if items == nil {
		return -1, fmt.Errorf("CalculateAverage | items cannot be nil")
	}
	if len(items) < 1 {
		return -1, fmt.Errorf("CalculateAverage | items cannot be empty")
	}

	sum := 0.0
	for i := range items {
		sum += items[i]
	}
	return sum / float64(len(items)), nil
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
		sum += individuals[i].AverageFitness
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
		generations[i].Mutex.Lock()
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
		generations[i].Mutex.Unlock()
	}
	return sortedGenerations, nil
}

// SortGenerationsThoroughlyByDelta sorts each kind of individual in each generation for every generation.
// This allows for easy querying in later phases.
func SortGenerationsThoroughlyByDelta(generations []*Generation, shouldAntagonistDeltaBig,
	shouldProtagonistDeltaBig bool) ([]*Generation, error) {
	if generations == nil {
		return nil, fmt.Errorf("SortGenerationsThoroughly | generations cannot be nil")
	}
	if len(generations) < 1 {
		return nil, fmt.Errorf("SortGenerationsThoroughly | generations cannot be empty")
	}

	sortedGenerations := make([]*Generation, len(generations))
	for i := 0; i < len(generations); i++ {
		sortedGenerations[i] = generations[i]
		sortedAntagonists, err := SortIndividualsByDelta(generations[i].Antagonists, shouldAntagonistDeltaBig)
		if err != nil {
			return nil, err
		}
		sortedProtagonists, err := SortIndividualsByDelta(generations[i].Protagonists, shouldProtagonistDeltaBig)
		if err != nil {
			return nil, err
		}
		sortedGenerations[i].Protagonists = sortedProtagonists
		sortedGenerations[i].Antagonists = sortedAntagonists
	}
	return sortedGenerations, nil
}

// SortGenerationsThoroughlyByAvgDelta sorts each kind of individual in each generation for every generation.
// This allows for easy querying in later phases.
func SortGenerationsThoroughlyByAvgDelta(generations []*Generation, shouldAntagonistDeltaBig,
	shouldProtagonistDeltaBig bool) ([]*Generation, error) {
	if generations == nil {
		return nil, fmt.Errorf("SortGenerationsThoroughly | generations cannot be nil")
	}
	if len(generations) < 1 {
		return nil, fmt.Errorf("SortGenerationsThoroughly | generations cannot be empty")
	}

	sortedGenerations := make([]*Generation, len(generations))
	for i := 0; i < len(generations); i++ {
		sortedGenerations[i] = generations[i]
		sortedAntagonists, err := SortIndividualsByAvgDelta(generations[i].Antagonists, shouldAntagonistDeltaBig)
		if err != nil {
			return nil, err
		}
		sortedProtagonists, err := SortIndividualsByAvgDelta(generations[i].Protagonists, shouldProtagonistDeltaBig)
		if err != nil {
			return nil, err
		}
		sortedGenerations[i].Protagonists = sortedProtagonists
		sortedGenerations[i].Antagonists = sortedAntagonists
	}
	return sortedGenerations, nil
}
