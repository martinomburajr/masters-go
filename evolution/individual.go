package evolution

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	IndividualAntagonist = 0
	IndividualProtagonist = 1
)

type Individual struct {
	id                       string
	strategy                 []Strategy
	fitness                  []float32
	hasAppliedStrategy       bool
	hasCalculatedFitness     bool
	fitnessCalculationMethod int
	kind                     int
	age                      int
	Program                  *Program
}

type Antagonist Individual
type Protagonist Individual

// CalculateFitness
func (i *Individual) CalculateFitness() float32 {
	switch i {

	}
	return 0
}

// GenerateRandomIndividuals creates a random number of individuals
func GenerateRandomIndividuals(number int, idTemplate string, kind int, strategyLength int,
	maxNumberOfStrategies int, availableStrategies []Strategy) ([]Individual, error) {
	if number < 1 {
		return nil, fmt.Errorf("number should at least be 1")
	}
	if kind < 0 || kind > 1 {
		return nil, fmt.Errorf("kind should be in bounds of [0,2)")
	}
	if strategyLength < 1 {
		return nil, fmt.Errorf("strategyLength should at least be 1")
	}
	if maxNumberOfStrategies < 1 {
		return nil, fmt.Errorf("maxNumberOfStrategies should at least be 1")
	}
	if availableStrategies == nil {
		return nil, fmt.Errorf("availableStrategies cannot be nil")
	}
	if len(availableStrategies) < 1 {
		return nil, fmt.Errorf("availableStrategies should at least have one strategy")
	}
	if idTemplate == "" {
		return nil, fmt.Errorf("idTemplate cannot be empty")
	}

	individuals := make([]Individual, number)

	for i := 0; i < number; i++ {
		rand.Seed(time.Now().UnixNano())
		numberOfStrategies := rand.Intn(maxNumberOfStrategies)
		randomStrategies := GenerateRandomStrategy(numberOfStrategies, strategyLength, availableStrategies)
		id := fmt.Sprintf("%s-%s-%d", KindToString(kind), idTemplate, i)

		individual := Individual{
			kind:     kind,
			id:       id,
			strategy: randomStrategies,
		}
		individuals[i] = individual
	}
	return individuals, nil
}

// GenerateRandomStrategy creates a random strategy list that contains some or all of the availableStrategies.
// They are randomly selected and populated.
func GenerateRandomStrategy(number int, strategyLength int, availableStrategies []Strategy) []Strategy {
	if number < 1 {
		number = 0
	}
	if strategyLength < 1 {
		strategyLength = 0
	}
	if availableStrategies == nil || len(availableStrategies) < 1 {
		return []Strategy{}
	}

	strategies := make([]Strategy, number)

	for i := 0; i < number; i++ {
		for j := 0; j < strategyLength; j++ {
			rand.Seed(time.Now().UnixNano())
			strategyIndex := rand.Intn(len(availableStrategies))
			strategies[i] = availableStrategies[strategyIndex]
		}
	}

	return strategies
}

// KindToString checks the kind and returns the appropriate string representation
func KindToString(kind int) string {
	switch kind {
	case IndividualAntagonist:
		return "ANTAGONIST"
	case IndividualProtagonist:
		return "PROTAGONIST"
	default:
		return "UNKNOWN"
	}
}