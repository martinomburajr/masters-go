package evolution

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
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
	NoOfCompetitions         int
	Mutex sync.Mutex
	// BirthGen represents the generation where this individual was spawned

	Program *Program // The best program generated
}

func (individual Individual) Clone() (Individual, error) {
	if individual.Program != nil {
		programClone, err := individual.Program.Clone()
		if err != nil {
			return Individual{}, err
		}
		individual.Program = &programClone
	}
	return individual, nil
}

// CloneCleanse removes performance based information but keeps the strategy intact.
func (individual Individual) CloneCleanse() (Individual, error) {
	if individual.Program != nil {
		programClone, err := individual.Program.Clone()
		if err != nil {
			return Individual{}, err
		}
		individual.Program = &programClone
	}

	individual.Id += "**"
	individual.Fitness = nil
	individual.Deltas = nil
	individual.AverageFitness = 0
	individual.AverageDelta = 0
	individual.Program = &Program{}
	individual.FitnessStdDev = 0
	individual.FitnessVariance = 0

	return individual, nil
}

// ApplyAntagonistStrategy applies the AntagonistEquation strategies to program.
func (individual *Individual) ApplyAntagonistStrategy(params EvolutionParams) error {
	if individual.Kind == IndividualProtagonist {
		return fmt.Errorf("ApplyAntagonistStrategy | cannot apply Antagonist Strategy to Protagonist")
	}
	if individual.Strategy == nil {
		return fmt.Errorf("antagonist stategy cannot be nil")
	}
	if len(individual.Strategy) < 1 {
		return fmt.Errorf("antagonist Strategy cannot be empty")
	}
	program, err := params.StartIndividual.Clone()
	if err != nil {
		return err
	}
	individual.Program = &program
	for _, strategy := range individual.Strategy {
		err := individual.Program.ApplyStrategy(strategy,
			params.SpecParam.AvailableSymbolicExpressions.Terminals,
			params.SpecParam.AvailableSymbolicExpressions.NonTerminals,
			params.Strategies.DepthOfRandomNewTrees)
		if err != nil {
			return err
		}
	}
	individual.HasAppliedStrategy = true
	if individual.Parent != nil {
		individual.Parent.NoOfCompetitions++
	}else {
		individual.NoOfCompetitions++
	}

	return nil
}

// ApplyProtagonistStrategy applies the AntagonistEquation strategies to program.
func (individual *Individual) ApplyProtagonistStrategy(antagonistTree DualTree, params EvolutionParams) error {
	if individual.Kind == IndividualAntagonist {
		return fmt.Errorf("ApplyProtagonistStrategy | cannot apply Protagonist Strategy to Antagonist")
	}
	if individual.Strategy == nil {
		return fmt.Errorf("protagonist stategy cannot be nil")
	}
	if len(individual.Strategy) < 1 {
		return fmt.Errorf("protagonist Strategy cannot be empty")
	}
	if antagonistTree.root == nil {
		return fmt.Errorf("applyProtagonistStrategy | antagonist supplied to protagonist has a nill root Tree")
	}

	tree, err := antagonistTree.Clone()
	if err != nil {
		return err
	}
	individual.Program.T = &tree

	for _, strategy := range individual.Strategy {
		err := individual.Program.ApplyStrategy(strategy,
			params.SpecParam.AvailableSymbolicExpressions.Terminals,
			params.SpecParam.AvailableSymbolicExpressions.NonTerminals,
			params.Strategies.DepthOfRandomNewTrees)
		if err != nil {
			return err
		}
	}
	individual.HasAppliedStrategy = true
	individual.HasAppliedStrategy = true
	if individual.Parent != nil {
		individual.Parent.NoOfCompetitions++
	}else {
		individual.NoOfCompetitions++
	}
	return nil
}

func (individual Individual) CloneWithTree(tree DualTree) Individual {
	individual.Id = GenerateIndividualID("", individual.Kind)

	programClone := individual.Program.CloneWithTree(tree)
	individual.Program = &programClone
	return individual
}

type Antagonist Individual
type Protagonist Individual

func (individual *Individual) ToString() strings.Builder {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("####   %s   ####\n", individual.Id))
	sb.WriteString(fmt.Sprintf("AGE:  %d\n", individual.Age))
	sb.WriteString(fmt.Sprintf("FITNESS:  %f\n", individual.AverageFitness))
	sb.WriteString(fmt.Sprintf("FITNESS-ARR:  %v\n", individual.Fitness))
	sb.WriteString(fmt.Sprintf("SPEC-DELTA:  %v\n", individual.BestDelta))
	sb.WriteString(fmt.Sprintf("BIRTH GEN:  %d\n", individual.BirthGen))
	strategiesSummary := FormatStrategiesTotal(individual.Strategy)
	sb.WriteString(fmt.Sprintf("Strategy Summary:\n%s\n", strategiesSummary.String()))
	strategiesList := FormatStrategiesList(individual.Strategy)
	sb.WriteString(fmt.Sprintf("Strategy Summary:%s\n", strategiesList.String()))
	if individual.Program != nil {
		dualTree := individual.Program.T
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

func (individual *Individual) CalculateProtagonistThresholdedFitness(params EvolutionParams) (
	protagonistFitness float64,
	delta float64, err error) {
	if !individual.HasAppliedStrategy {
		return 0, 0, fmt.Errorf(" CalculateProtagonistThresholdedFitness | has not applied strategies")
	}
	if individual.Kind == IndividualAntagonist {
		return 0, 0, fmt.Errorf(" CalculateProtagonistThresholdedFitness | cannot apply protagonist antagonist" +
			" fitness to" +
			" antagonist")
	}

	fitnessPenalization := params.Spec[0].DivideByZeroPenalty
	badDeltaValue := math.Inf(1)
	divByZeroStrategy := params.SpecParam.DivideByZeroStrategy

	protagonistExpression, err := individual.Program.T.ToMathematicalString()
	if err != nil {
		return 0, 0, err
	}

	deltaProtagonist := 0.0
	deltaProtagonistThreshold := 0.0
	protagonistDividedByZeroCount := 0
	isProtagonistValid := true
	spec := params.Spec

	for i := range spec {
		independentX := spec[i].Independents
		independentXVal := spec[i].Independents["x"]
		if isProtagonistValid {
			dependentProtagonistVar, err := individual.Program.EvalMulti(independentX, protagonistExpression)
			if err != nil {
				switch divByZeroStrategy {
				case DivByZeroIgnore:
					dependentProtagonistVar = 0

				case DivByZeroPenalize:
					isProtagonistValid = false
					protagonistFitness = fitnessPenalization

				case DivByZeroSetSpecValueZero:
					dependentProtagonistVar = 0

				case DivByZeroSteadyPenalize:
					if independentXVal != 0 {
						if math.IsNaN(dependentProtagonistVar) || dependentProtagonistVar == 0 {
							isProtagonistValid = false
							protagonistFitness = fitnessPenalization
						} else {
							protagonistDividedByZeroCount++
						}
					} else {
						// Unlikely to ever reach here
						protagonistDividedByZeroCount++
					}
				}
			} else {
				diff := spec[i].Dependent - dependentProtagonistVar
				deltaProtagonist += diff * diff
			}
		}
		deltaProtagonistThreshold += spec[i].ProtagonistThreshold * spec[i].ProtagonistThreshold
	}

	specLen := float64(len(spec))

	deltaProtagonist = math.Sqrt(deltaProtagonist / specLen)
	deltaProtagonistThreshold = math.Sqrt(deltaProtagonistThreshold / specLen)

	if isProtagonistValid {
		if deltaProtagonist <= deltaProtagonistThreshold {
			if deltaProtagonist == 0 {
				protagonistFitness = 1
				protagonistDividedByZeroCount = -1
			} else {
				protagonistFitness = (deltaProtagonistThreshold - deltaProtagonist) / deltaProtagonistThreshold
			}
		} else {
			protagonistFitness = -1 * ((deltaProtagonist - deltaProtagonistThreshold) / deltaProtagonist)
		}

		return protagonistFitness, deltaProtagonist, nil
	} else {
		protagonistFitness = fitnessPenalization
		deltaProtagonist = badDeltaValue
	}

	if protagonistDividedByZeroCount > 0 {
		if protagonistFitness > 0 {
			protagonistFitness = protagonistFitness - (protagonistFitness * 0.1 * float64(
				protagonistDividedByZeroCount))
		}
	}
	return protagonistFitness, deltaProtagonist, nil
}

func (individual *Individual) CalculateAntagonistThresholdedFitness(params EvolutionParams) (antagonistFitness float64,
	delta float64, err error) {
	if !individual.HasAppliedStrategy {
		return 0, 0, fmt.Errorf(" CalculateAntagonistThresholdedFitness | has not applied strategies")
	}
	if individual.Kind == IndividualProtagonist {
		return 0, 0, fmt.Errorf(" CalculateAntagonistThresholdedFitness | cannot apply antagonist fitness to" +
			" protagonist")
	}

	fitnessPenalization := params.Spec[0].DivideByZeroPenalty
	badDeltaValue := math.Inf(1)
	divByZeroStrategy := params.SpecParam.DivideByZeroStrategy

	antagonistExpression, err := individual.Program.T.ToMathematicalString()
	if err != nil {
		return 0, 0, err
	}

	deltaAntagonist := 0.0
	deltaAntagonistThreshold := 0.0
	antagonistDividedByZeroCount := 0
	isAntagonistValid := true

	spec := params.Spec

	for i := range spec {
		independentX := spec[i].Independents
		independentXVal := spec[i].Independents["x"]
		if isAntagonistValid {
			dependentAntagonistVar, err := individual.Program.EvalMulti(independentX, antagonistExpression)
			if err != nil {
				switch divByZeroStrategy {
				case DivByZeroIgnore:
					dependentAntagonistVar = 0

				case DivByZeroPenalize:
					isAntagonistValid = false
					antagonistFitness = fitnessPenalization

				case DivByZeroSetSpecValueZero:
					dependentAntagonistVar = 0

				case DivByZeroSteadyPenalize:
					if independentXVal != 0 {
						// If the spec does not contain a zero,
						// yet you still divide by zero. Give maximum penalty!
						if math.IsNaN(dependentAntagonistVar) || dependentAntagonistVar == 0 {
							isAntagonistValid = false
							antagonistFitness = fitnessPenalization
						} else {
							antagonistDividedByZeroCount++
						}
					} else {
						// Unlikely to ever reach here
						antagonistDividedByZeroCount++
					}
				}
			} else {
				diff := spec[i].Dependent - dependentAntagonistVar
				deltaAntagonist += diff * diff
			}
		}
		deltaAntagonistThreshold += math.Abs(spec[i].AntagonistThreshold) * math.Abs(spec[i].AntagonistThreshold)
	}

	specLen := float64(len(spec))

	deltaAntagonist = math.Sqrt(deltaAntagonist / specLen)
	deltaAntagonistThreshold = math.Sqrt(deltaAntagonistThreshold / specLen)

	if !isAntagonistValid {
		// TODO is math.Nan the best alternative?
		return fitnessPenalization, badDeltaValue, nil
	} else {
		if deltaAntagonist >= deltaAntagonistThreshold { // good thing
			if deltaAntagonist == 0 { // This is to punish deltaAntagonist for coalescing near the spec
				antagonistFitness = -1
				antagonistDividedByZeroCount = -1
			} else {
				// Award fitness if it did not cluster around the spec
				antagonistFitness = (deltaAntagonist - deltaAntagonistThreshold) / deltaAntagonist
			}
		} else {
			antagonistFitness = -1 * ((deltaAntagonistThreshold - deltaAntagonist) / deltaAntagonistThreshold)
		}

		if antagonistDividedByZeroCount > 0 {
			if antagonistFitness > 0 {
				antagonistFitness = antagonistFitness - (antagonistFitness * 0.1 * float64(
					antagonistDividedByZeroCount))
			}
			// No else statement as if the antagonist is already less than 0, it should remain there.
		}
		return antagonistFitness, deltaAntagonist, nil
	}
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
