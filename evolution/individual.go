package evolution

import (
	"fmt"
	"math/rand"
	"strings"
)

const (
	IndividualAntagonist  = 0
	IndividualProtagonist = 1
)

type Individual struct {
	Id                       string
	Parent                   *Individual
	Strategy                 []Strategy
	Fitness                  []float64
	Deltas                   []float64
	FitnessVariance          float64
	FitnessStdDev            float64
	HasAppliedStrategy       bool
	HasCalculatedFitness     bool
	FitnessCalculationMethod string
	Kind                     int
	BirthGen                 int
	Age                      int
	BestFitness              float64 // Best fitness from all epochs
	AverageFitness           float64 // Measures average fitness throughout epoch
	BestDelta                float64
	AverageDelta             float64
	// BirthGen represents the generation where this individual was spawned

	Program *Program // The best program generated
}

func (i Individual) Clone() (Individual, error) {
	if i.Program != nil {
		programClone, err := i.Program.Clone()
		if err != nil {
			return Individual{}, err
		}
		i.Program = &programClone
	}
	return i, nil
}

func (i Individual) CloneWithTree(tree DualTree) Individual {
	i.Id = GenerateIndividualID("", i.Kind)

	programClone := i.Program.CloneWithTree(tree)
	i.Program = &programClone
	return i
}

type Antagonist Individual
type Protagonist Individual


func (i *Individual) ToString() strings.Builder {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("####   %s   ####\n", i.Id))
	sb.WriteString(fmt.Sprintf("AGE:  %d\n", i.Age))
	sb.WriteString(fmt.Sprintf("FITNESS:  %f\n", i.AverageFitness))
	sb.WriteString(fmt.Sprintf("FITNESS-ARR:  %v\n", i.Fitness))
	sb.WriteString(fmt.Sprintf("SPEC-DELTA:  %v\n", i.BestDelta))
	sb.WriteString(fmt.Sprintf("BIRTH GEN:  %d\n", i.BirthGen))
	strategiesSummary := FormatStrategiesTotal(i.Strategy)
	sb.WriteString(fmt.Sprintf("Strategy Summary:\n%s\n", strategiesSummary.String()))
	strategiesList := FormatStrategiesList(i.Strategy)
	sb.WriteString(fmt.Sprintf("Strategy Summary:%s\n", strategiesList.String()))
	if i.Program != nil {
		dualTree := i.Program.T
		if dualTree != nil {
			toString := dualTree.ToString()
			sb.WriteString(fmt.Sprintf("TREE:  \n%s", toString.String()))
			mathematicalString, err := dualTree.ToMathematicalString()
			if err != nil {

			} else {
				sb.WriteString(fmt.Sprintf("Mathematical Expression: %s\n", mathematicalString))
			}
		}
	}

	return sb
}

func GenerateIndividualID(identifier string, individualKind int) string {
	return fmt.Sprintf("%s-%s%s", KindToString(individualKind), RandString(4), identifier)
}

// GenerateRandomStrategy creates a random Strategy list that contains some or all of the availableStrategies.
// They are randomly selected and populated.
func GenerateRandomStrategy(number int, availableStrategies []Strategy) []Strategy {
	if number < 1 {
		number = 1
	}
	if availableStrategies == nil || len(availableStrategies) < 1 {
		return []Strategy{}
	}

	strategies := make([]Strategy, number)

	for i := 0; i < number; i++ {
		strategyIndex := rand.Intn(len(availableStrategies))
		strategies[i] = availableStrategies[strategyIndex]
	}

	return strategies
}

// KindToString checks the Kind and returns the appropriate string representation
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

func DominantStrategy(individual Individual) string {
	domStrat := map[string]int{}
	for i := range individual.Strategy {
		strategy := string(individual.Strategy[i])

		stratCount := domStrat[strategy]
		domStrat[strategy] = stratCount + 1
	}

	var topStrategy string
	counter := 0
	for k, v := range domStrat {
		if v > counter {
			counter = v
			topStrategy = k
		}
	}
	return topStrategy
}

func DominantStrategyStr(str string) string {
	strategies := strings.Split(str, "|")

	domStrat := map[string]int{}
	for i := range strategies {
		strategy := string(strategies[i])
		stratCount := domStrat[strategy]
		if domStrat[strategy] > -1 {
			domStrat[strategy] = stratCount + 1
		}
	}

	var topStrategy string
	counter := 0
	for k, v := range domStrat {
		if v > counter {
			counter = v
			topStrategy = k
		}
	}
	return topStrategy
}

func StrategiesToString(individual Individual) string {
	sb := strings.Builder{}
	for _, strategy := range individual.Strategy {
		sb.WriteString(string(strategy))
		sb.WriteString("|")
	}

	final := sb.String()
	return final[:len(final)-1]
}


func StrategiesToStringArr(strategies []string) string {
	sb := strings.Builder{}
	for _, strategy := range strategies {
		sb.WriteString(string(strategy))
		sb.WriteString("|")
	}

	final := sb.String()
	return final[:len(final)-1]
}

func ConvertStrategiesToString(strategies []Strategy) (stringStrategies []string) {
	for i := range strategies {
		stringStrategies = append(stringStrategies, string(strategies[i]))
	}
	return stringStrategies
}