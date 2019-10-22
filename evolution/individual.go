package evolution

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

const (
	IndividualAntagonist  = 0
	IndividualProtagonist = 1
)

type Individual struct {
	Id                       string
	Strategy                 []Strategy
	Fitness                  []float64
	HasAppliedStrategy       bool
	HasCalculatedFitness     bool
	FitnessCalculationMethod int
	Kind                     int
	Age                      int
	TotalFitness             float64
	// BirthGen represents the generation where this individual was spawned
	BirthGen int
	Program  *Program
}

func (i Individual) Clone() (Individual, error) {
	i.Id = i.Id +"c"
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

// Crossover will perform crossover on the strategies of a given  set of individuals
func Crossover(individual Individual, individual2 Individual, params EvolutionParams) (Individual, Individual,
	error) {

	if individual.Id == "" {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual1 - individual Id cannot be empty")
	}
	if individual.Strategy == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual1 - Strategy array cannot be nil")
	}
	if len(individual.Strategy) == 0 {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual1 - Strategy array cannot be empty")
	}
	if individual.HasCalculatedFitness == false {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual1 - HasCalculatedFitness should be true")
	}
	if individual.HasAppliedStrategy == false {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual1 - HasAppliedStrategy should be true")
	}
	//if individual.Program == nil {
	//	return Individual{}, Individual{}, fmt.Errorf("crossover | individual1 - program cannot be nil")
	//}
	//if individual.Program.T == nil {
	//	return Individual{}, Individual{}, fmt.Errorf("crossover | individual1 - program Tree cannot be nil")
	//}
	if individual2.Id == "" {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual2 - individual Id cannot be empty")
	}
	if individual2.Strategy == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual2 - Strategy array cannot be nil")
	}
	if len(individual2.Strategy) == 0 {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual2 - Strategy array cannot be empty")
	}
	if individual2.HasCalculatedFitness == false {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual2 - HasCalculatedFitness should be true")
	}
	if individual2.HasAppliedStrategy == false {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual2 - HasAppliedStrategy should be true")
	}
	//if individual2.Program == nil {
	//	return Individual{}, Individual{}, fmt.Errorf("crossover | individual2 - program cannot be nil")
	//}
	//if individual2.Program.T == nil {
	//	return Individual{}, Individual{}, fmt.Errorf("crossover | individual2 - program Tree cannot be nil")
	//}
	if params.StrategyLengthLimit < 1 {
		return Individual{}, Individual{}, fmt.Errorf("crossover | params.StrategyLengthLimit must be greater than 0")
	}

	individual1Len := len(individual.Strategy)
	individual2Len := len(individual2.Strategy)

	child1, err := individual.Clone()
	if err != nil {
		return Individual{}, Individual{}, err
	}
	child2, err := individual2.Clone()
	if err != nil {
		return Individual{}, Individual{}, err
	}

	crossoverPercentage := params.CrossoverPercentage
	if crossoverPercentage == 0 {
		return child1, child2, err
	}
	if crossoverPercentage == 1 {
		return child2, child1, err
	}

	individual1ChunkSize := int(math.Ceil(float64(individual1Len) * float64(crossoverPercentage)))
	individual2ChunkSize := int(float64(individual2Len) * crossoverPercentage)

	if individual1ChunkSize >= individual2ChunkSize {
		if params.MaintainCrossoverGeneTransferEquality {
			var ind1StartIndex int
			if individual1Len == individual1ChunkSize {
				ind1StartIndex = 0
			} else {
				ind1StartIndex = rand.Intn((individual1Len + 1) - individual1ChunkSize)
			}
			c1, c2 := StrategySwapper(individual.Strategy, individual2.Strategy, individual1ChunkSize, ind1StartIndex)
			child1.Strategy = c1
			child2.Strategy = c2
			return child1, child2, nil
		} else {

		}
	} else {
		if params.MaintainCrossoverGeneTransferEquality {
			var ind2StartIndex int
			if individual2Len == individual2ChunkSize {
				ind2StartIndex = 0
			} else {
				ind2StartIndex = rand.Intn(individual1Len + 1 - individual1ChunkSize)
			}
			c1, c2 := StrategySwapper(individual.Strategy, individual2.Strategy, individual1ChunkSize, ind2StartIndex)
			child1.Strategy = c1
			child2.Strategy = c2
			return child1, child2, nil
		} else {

		}
	}

	//ind1Copy := make([]Strategy, individual1Len)
	//copy(ind1Copy, individual.Strategy)
	//ind2Copy := make([]Strategy, individual2Len)
	//copy(ind2Copy, individual2.Strategy)
	//
	//if individual1Len <= individual2Len {
	//	individual1ChunkSize := int(math.Ceil(float64(individual1Len) * float64(crossoverPercentage)))
	//	if params.MaintainCrossoverGeneTransferEquality {
	//
	//		var ind1StartIndex int
	//		if individual1Len == individual1ChunkSize {
	//			ind1StartIndex = 0
	//		} else {
	//			ind1StartIndex = rand.Intn(individual1Len+1 - individual1ChunkSize)
	//		}
	//
	//		ind2StartIndex := rand.Intn(individual2Len+1 - individual1ChunkSize)
	//
	//		for i := 0; i < individual1ChunkSize; i++ {
	//			child2.Strategy[ind2StartIndex+i] = ind1Copy[ind1StartIndex+i]
	//			child1.Strategy[ind1StartIndex+i] = ind2Copy[ind2StartIndex+i]
	//		}
	//	} else {
	//		individual2ChunkSize := int(float32(individual2Len) * crossoverPercentage)
	//
	//		ind1StartIndex := rand.Intn(individual1Len - individual1ChunkSize)
	//		ind1EndIndex := ind1StartIndex + individual1ChunkSize
	//		ind2StartIndex := rand.Intn(individual2Len - individual2ChunkSize)
	//		ind2EndIndex := ind2StartIndex + individual2ChunkSize
	//
	//		ind1Chunk := make([]Strategy, individual1ChunkSize)
	//		ind2Chunk := make([]Strategy, individual2ChunkSize)
	//		for i := 0; i < individual1ChunkSize; i++ {
	//			ind1Chunk[i] = child1.Strategy[ind1StartIndex+i]
	//		}
	//		child1.Strategy = append(child1.Strategy[:ind1StartIndex],
	//			child1.Strategy[:ind1EndIndex]...) // REMOVE ITEMS COPIED TO CHUNK
	//		for i := 0; i < individual2ChunkSize; i++ {
	//			ind2Chunk[i] = child2.Strategy[ind2StartIndex+i]
	//		}
	//		child2.Strategy = append(child2.Strategy[:ind2StartIndex], child2.Strategy[:ind2EndIndex]...) // REMOVE ITEMS COPIED TO CHUNK
	//
	//		child1.Strategy = append(child1.Strategy[:ind1StartIndex], append(child2.Strategy,
	//			child1.Strategy[ind1StartIndex:]...)...) // INSERT TO CHILD1
	//		child2.Strategy = append(child2.Strategy[:ind2StartIndex], append(child1.Strategy,
	//			child2.Strategy[ind2StartIndex:]...)...)
	//		log.Print()
	//	}
	//} else {
	//	individual2ChunkSize := int(math.Ceil(float64(individual2Len) * float64(crossoverPercentage)))
	//	if params.MaintainCrossoverGeneTransferEquality {
	//
	//		var ind2StartIndex int
	//		if individual2Len == individual2ChunkSize {
	//			ind2StartIndex = 0
	//		} else {
	//			ind2StartIndex = rand.Intn(individual2Len+1 - individual2ChunkSize)
	//		}
	//
	//		ind1StartIndex := rand.Intn(individual1Len+1 - individual2ChunkSize)
	//
	//		for i := 0; i < individual2ChunkSize; i++ {
	//			child1.Strategy[ind1StartIndex+i] = ind1Copy[ind2StartIndex+i]
	//			child2.Strategy[ind2StartIndex+i] = ind2Copy[ind1StartIndex+i]
	//		}
	//	} else {
	//		individual1ChunkSize := int(float32(individual1Len) * crossoverPercentage)
	//
	//		ind2StartIndex := rand.Intn(individual2Len - individual2ChunkSize)
	//		ind2EndIndex := ind2StartIndex + individual2ChunkSize
	//		ind1StartIndex := rand.Intn(individual1Len - individual1ChunkSize)
	//		ind1EndIndex := ind1StartIndex + individual1ChunkSize
	//
	//		ind2Chunk := make([]Strategy, individual2ChunkSize)
	//		ind1Chunk := make([]Strategy, individual1ChunkSize)
	//		for i := 0; i < individual2ChunkSize; i++ {
	//			ind2Chunk[i] = child2.Strategy[ind2StartIndex+i]
	//		}
	//		child2.Strategy = append(child2.Strategy[:ind2StartIndex],
	//			child2.Strategy[:ind2EndIndex]...) // REMOVE ITEMS COPIED TO CHUNK
	//		for i := 0; i < individual1ChunkSize; i++ {
	//			ind1Chunk[i] = child1.Strategy[ind1StartIndex+i]
	//		}
	//		child1.Strategy = append(child1.Strategy[:ind1StartIndex],
	//			child1.Strategy[:ind1EndIndex]...) // REMOVE ITEMS COPIED TO CHUNK
	//
	//		child2.Strategy = append(child2.Strategy[:ind2StartIndex], append(child1.Strategy,
	//			child2.Strategy[ind2StartIndex:]...)...) // INSERT TO CHILD1
	//		child1.Strategy = append(child1.Strategy[:ind1StartIndex], append(child2.Strategy,
	//			child1.Strategy[ind1StartIndex:]...)...)
	//	}
	//}
	return child1, child2, err
}

// StrategySwapper takes two slices containing variable length strategies.
// The swapLength must be smaller than the length of the largest, but less than the length of the smallest.
// A swap length of 0 will return the same arrays a and b untouched.
func StrategySwapper(a []Strategy, b []Strategy, swapLength int, startIndex int) ([]Strategy, []Strategy) {
	if a == nil || b == nil {
		return nil, nil
	}
	if len(a) == 0 || len(b) == 0 {
		return nil, nil
	}
	if swapLength == 0 {
		return a, b
	}
	if swapLength < 0 {
		swapLength = 0
	}
	if startIndex < 0 {
		startIndex = 0
	}

	aCopy := make([]Strategy, len(a))
	bCopy := make([]Strategy, len(b))

	copy(aCopy, a)
	copy(bCopy, b)

	if len(a) >= len(b) {
		if swapLength > len(b) {
			swapLength = len(b)
		}
		if (swapLength + startIndex) > len(b) {
			startIndex = 0
		}
	} else {
		if swapLength > len(a) {
			swapLength = len(a)
		}
		if (swapLength + startIndex) > len(a) {
			startIndex = 0
		}
	}

	aHolder := make([]Strategy, swapLength)
	bHolder := make([]Strategy, swapLength)

	for i := 0; i < swapLength; i++ {
		aHolder[i] = a[i+startIndex]
		bHolder[i] = b[i+startIndex]
	}

	for i := 0; i < swapLength; i++ {
		aCopy[startIndex+i] = bHolder[i]
		bCopy[startIndex+i] = aHolder[i]
	}

	return aCopy, bCopy
}

// StrategySwapperIgnorant will perform crossover regardless of size
func StrategySwapperIgnorant(a []Strategy, b []Strategy, swapLength int, startIndex int) ([]Strategy, []Strategy) {
	if a == nil || b == nil {
		return nil, nil
	}
	if len(a) == 0 || len(b) == 0 {
		return nil, nil
	}
	if swapLength == 0 {
		return a, b
	}
	if swapLength < 0 {
		swapLength = 0
	}
	if startIndex < 0 {
		startIndex = 0
	}
	var aCopy, bCopy, aHolder, bHolder []Strategy

	if len(a) >= len(b) {
		if swapLength > len(a) {
			swapLength = len(a)
		}
		if startIndex+swapLength >= len(a) {
			startIndex = 0
		}
		aCopy = make([]Strategy, len(a))
		bCopy = make([]Strategy, len(a))
		aHolder = make([]Strategy, swapLength)
		bHolder = make([]Strategy, swapLength)
		copy(aCopy, a)
		copy(bCopy, b)

		for i := 0; i < swapLength; i++ {
			aHolder[i] = a[i+startIndex]
		}
		for i := 0; i < swapLength; i++ {
			bHolder[i] = b[i+startIndex]
		}
	} else {
		if swapLength > len(b) {
			swapLength = len(b)
		}
		if startIndex+swapLength >= len(b) {
			startIndex = 0
		}
		aCopy = make([]Strategy, len(b))
		bCopy = make([]Strategy, len(b))
		aHolder = make([]Strategy, swapLength)
		bHolder = make([]Strategy, swapLength)
		copy(aCopy, a)
		copy(bCopy, b)

		for i := 0; i < len(aCopy); i++ {
			aHolder[i] = a[i+startIndex]
		}
		for i := 0; i < len(bCopy); i++ {
			bHolder[i] = b[i+startIndex]
		}
	}

	for i := 0; i < swapLength; i++ {
		aCopy[startIndex+i] = bHolder[i]
		bCopy[startIndex+i] = aHolder[i]
	}

	return aCopy, bCopy
}

// Mutate will mutate the Strategy in a given individual
func (i *Individual) Mutate(availableStrategies []Strategy) error {
	if availableStrategies == nil {
		return fmt.Errorf("Mutate | availableStrategies param cannot be nil")
	}
	if i.Strategy == nil {
		return fmt.Errorf("Mutate | individual's strategies cannot be nil")
	}
	if len(i.Strategy) < 1 {
		return fmt.Errorf("Mutate | individual's strategies cannot empty")
	}

	randIndexToMutate := rand.Intn(len(i.Strategy))

	randIndexForStrategies := rand.Intn(len(availableStrategies))
	i.Strategy[randIndexToMutate] = availableStrategies[randIndexForStrategies]
	return nil
}

func (i *Individual) ToString() strings.Builder {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("####   %s   ####\n", i.Id))
	sb.WriteString(fmt.Sprintf("AGE:  %d\n", i.Age))
	sb.WriteString(fmt.Sprintf("FITNESS:  %f\n", i.TotalFitness))
	sb.WriteString(fmt.Sprintf("FITNESS-ARR:  %v\n", i.Fitness))
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
	return fmt.Sprintf("%s-%s%s", KindToString(individualKind), RandString(3), identifier)
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
