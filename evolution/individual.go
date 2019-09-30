package evolution

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	IndividualAntagonist  = 0
	IndividualProtagonist = 1
)

type Individual struct {
	id                       string
	strategy                 []Strategy
	fitness                  []int
	hasAppliedStrategy       bool
	hasCalculatedFitness     bool
	fitnessCalculationMethod int
	kind                     int
	age                      int
	totalFitness             int
	Program                  *Program
}

func (i Individual) Clone() (Individual, error) {
	i.id = GenerateIndividualID("", i.kind)
	i.fitness = make([]int, 0)
	i.totalFitness = 0
	i.hasCalculatedFitness = false
	i.age = 0

	programClone, err := i.Program.Clone()
	if err != nil {
		return Individual{}, err
	}
	i.Program = &programClone
	return i, nil
}

type Antagonist Individual
type Protagonist Individual

func GenerateIndividualID(identifier string, individualKind int) string {
	return fmt.Sprintf("%s-%s-%s%s", "individual", KindToString(individualKind), RandString(3), identifier)
}

func GenerateRandomIndividuals(number int, idTemplate string, kind int, strategyLength int,
	maxNumberOfStrategies int, availableStrategies []Strategy, depth int, terminals SymbolicExpressionSet,
	nonTerminals SymbolicExpressionSet) ([]*Individual, error) {
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

	individuals := make([]*Individual, number)

	for i := 0; i < number; i++ {
		rand.Seed(time.Now().UnixNano())
		numberOfStrategies := rand.Intn(maxNumberOfStrategies)
		randomStrategies := GenerateRandomStrategy(numberOfStrategies, strategyLength, availableStrategies)
		id := fmt.Sprintf("%s-%s-%d", KindToString(kind), "", i)

		program := Program{}

		programID := GenerateProgramID(i)
		tree, err := GenerateRandomTree(depth, terminals, nonTerminals)
		if err != nil {
			return nil, err
		}
		program.T = tree
		program.ID = programID

		individual := &Individual{
			kind:     kind,
			id:       id,
			strategy: randomStrategies,
			fitness:  make([]int, 0),
			Program:  &program,
		}
		individuals[i] = individual
	}
	return individuals, nil
}

// GenerateRandomStrategy creates a random strategy list that contains some or all of the availableStrategies.
// They are randomly selected and populated.
func GenerateRandomStrategy(number int, strategyLength int, availableStrategies []Strategy) []Strategy {
	if number < 1 {
		number = 1
	}
	if strategyLength < 1 {
		strategyLength = 1
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
