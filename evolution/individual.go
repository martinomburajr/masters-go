package evolution

import (
	"fmt"
	"math"
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

	programClone, err := i.Program.Clone()
	if err != nil {
		return Individual{}, err
	}
	i.Program = &programClone
	return i, nil
}

type Antagonist Individual
type Protagonist Individual

// Crossover will perform crossover on the strategies of a given  set of individuals
func Crossover(individual Individual, individual2 Individual, params EvolutionParams) (Individual, Individual,
	error) {

	if individual.id == "" {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual1 - individual id cannot be empty")
	}
	if individual.strategy == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual1 - strategy array cannot be nil")
	}
	if len(individual.strategy) == 0 {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual1 - strategy array cannot be empty")
	}
	if individual.hasCalculatedFitness == false {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual1 - hasCalculatedFitness should be true")
	}
	if individual.hasAppliedStrategy == false {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual1 - hasAppliedStrategy should be true")
	}
	if individual.Program == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual1 - program cannot be nil")
	}
	if individual.Program.T == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual1 - program Tree cannot be nil")
	}
	if individual2.id == "" {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual2 - individual id cannot be empty")
	}
	if individual2.strategy == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual2 - strategy array cannot be nil")
	}
	if len(individual2.strategy) == 0 {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual2 - strategy array cannot be empty")
	}
	if individual2.hasCalculatedFitness == false {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual2 - hasCalculatedFitness should be true")
	}
	if individual2.hasAppliedStrategy == false {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual2 - hasAppliedStrategy should be true")
	}
	if individual2.Program == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual2 - program cannot be nil")
	}
	if individual2.Program.T == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual2 - program Tree cannot be nil")
	}
	if params.StrategyLengthLimit < 1 {
		return Individual{}, Individual{}, fmt.Errorf("crossover | params.StrategyLengthLimit must be greater than 0")
	}

	individual1Len := len(individual.strategy)
	individual2Len := len(individual2.strategy)

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
	individual2ChunkSize := int(float32(individual2Len) * crossoverPercentage)

	if individual1ChunkSize >= individual1ChunkSize {
		//if params.MaintainCrossoverGeneTransferEquality {
			rand.Seed(time.Now().UnixNano())
			var ind1StartIndex int
			if individual1Len == individual1ChunkSize {
				ind1StartIndex = 0
			} else {
				ind1StartIndex = rand.Intn(individual1Len+1 - individual1ChunkSize)
			}
			c1, c2 := StrategySwapper(individual.strategy, individual2.strategy, individual1ChunkSize, ind1StartIndex)
			child1.strategy = c1
			child2.strategy = c2
			return child1, child2, nil
		//} else {
		//
		//}
	} else {
		rand.Seed(time.Now().UnixNano())
		var ind2StartIndex int
		if individual2Len == individual2ChunkSize {
			ind2StartIndex = 0
		} else {
			ind2StartIndex = rand.Intn(individual1Len+1 - individual1ChunkSize)
		}
		c1, c2 := StrategySwapper(individual.strategy, individual2.strategy, individual1ChunkSize, ind2StartIndex)
		child1.strategy = c1
		child2.strategy = c2
		return child1, child2, nil
	}


	//ind1Copy := make([]Strategy, individual1Len)
	//copy(ind1Copy, individual.strategy)
	//ind2Copy := make([]Strategy, individual2Len)
	//copy(ind2Copy, individual2.strategy)
	//
	//if individual1Len <= individual2Len {
	//	individual1ChunkSize := int(math.Ceil(float64(individual1Len) * float64(crossoverPercentage)))
	//	if params.MaintainCrossoverGeneTransferEquality {
	//		rand.Seed(time.Now().UnixNano())
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
	//			child2.strategy[ind2StartIndex+i] = ind1Copy[ind1StartIndex+i]
	//			child1.strategy[ind1StartIndex+i] = ind2Copy[ind2StartIndex+i]
	//		}
	//	} else {
	//		individual2ChunkSize := int(float32(individual2Len) * crossoverPercentage)
	//		rand.Seed(time.Now().UnixNano())
	//		ind1StartIndex := rand.Intn(individual1Len - individual1ChunkSize)
	//		ind1EndIndex := ind1StartIndex + individual1ChunkSize
	//		ind2StartIndex := rand.Intn(individual2Len - individual2ChunkSize)
	//		ind2EndIndex := ind2StartIndex + individual2ChunkSize
	//
	//		ind1Chunk := make([]Strategy, individual1ChunkSize)
	//		ind2Chunk := make([]Strategy, individual2ChunkSize)
	//		for i := 0; i < individual1ChunkSize; i++ {
	//			ind1Chunk[i] = child1.strategy[ind1StartIndex+i]
	//		}
	//		child1.strategy = append(child1.strategy[:ind1StartIndex],
	//			child1.strategy[:ind1EndIndex]...) // REMOVE ITEMS COPIED TO CHUNK
	//		for i := 0; i < individual2ChunkSize; i++ {
	//			ind2Chunk[i] = child2.strategy[ind2StartIndex+i]
	//		}
	//		child2.strategy = append(child2.strategy[:ind2StartIndex], child2.strategy[:ind2EndIndex]...) // REMOVE ITEMS COPIED TO CHUNK
	//
	//		child1.strategy = append(child1.strategy[:ind1StartIndex], append(child2.strategy,
	//			child1.strategy[ind1StartIndex:]...)...) // INSERT TO CHILD1
	//		child2.strategy = append(child2.strategy[:ind2StartIndex], append(child1.strategy,
	//			child2.strategy[ind2StartIndex:]...)...)
	//		log.Print()
	//	}
	//} else {
	//	individual2ChunkSize := int(math.Ceil(float64(individual2Len) * float64(crossoverPercentage)))
	//	if params.MaintainCrossoverGeneTransferEquality {
	//		rand.Seed(time.Now().UnixNano())
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
	//			child1.strategy[ind1StartIndex+i] = ind1Copy[ind2StartIndex+i]
	//			child2.strategy[ind2StartIndex+i] = ind2Copy[ind1StartIndex+i]
	//		}
	//	} else {
	//		individual1ChunkSize := int(float32(individual1Len) * crossoverPercentage)
	//		rand.Seed(time.Now().UnixNano())
	//		ind2StartIndex := rand.Intn(individual2Len - individual2ChunkSize)
	//		ind2EndIndex := ind2StartIndex + individual2ChunkSize
	//		ind1StartIndex := rand.Intn(individual1Len - individual1ChunkSize)
	//		ind1EndIndex := ind1StartIndex + individual1ChunkSize
	//
	//		ind2Chunk := make([]Strategy, individual2ChunkSize)
	//		ind1Chunk := make([]Strategy, individual1ChunkSize)
	//		for i := 0; i < individual2ChunkSize; i++ {
	//			ind2Chunk[i] = child2.strategy[ind2StartIndex+i]
	//		}
	//		child2.strategy = append(child2.strategy[:ind2StartIndex],
	//			child2.strategy[:ind2EndIndex]...) // REMOVE ITEMS COPIED TO CHUNK
	//		for i := 0; i < individual1ChunkSize; i++ {
	//			ind1Chunk[i] = child1.strategy[ind1StartIndex+i]
	//		}
	//		child1.strategy = append(child1.strategy[:ind1StartIndex],
	//			child1.strategy[:ind1EndIndex]...) // REMOVE ITEMS COPIED TO CHUNK
	//
	//		child2.strategy = append(child2.strategy[:ind2StartIndex], append(child1.strategy,
	//			child2.strategy[ind2StartIndex:]...)...) // INSERT TO CHILD1
	//		child1.strategy = append(child1.strategy[:ind1StartIndex], append(child2.strategy,
	//			child1.strategy[ind1StartIndex:]...)...)
	//	}
	//}
	return child1, child2, err
}

// StrategySwapper takes two slices containing variable length strategies.
// The swapLength must be smaller than the length of the largest, but less than the length of the smallest.
// A swap length of 0 will return the same arrays a and b untouched.
func StrategySwapper(a []Strategy, b []Strategy, swapLength int, startIndex int) ([]Strategy, []Strategy) {
	if a == nil || b == nil{
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
		if (swapLength + startIndex) > len(b)  {
			startIndex = 0
		}
	}else {
		if swapLength > len(a) {
			swapLength = len(a)
		}
		if (swapLength + startIndex) > len(a)  {
			startIndex = 0
		}
	}


	aHolder := make([]Strategy, swapLength)
	bHolder := make([]Strategy, swapLength)

	for i:= 0; i < swapLength; i++ {
		aHolder[i] = a[i + startIndex]
		bHolder[i] = b[i + startIndex]
	}


	for i:= 0; i < swapLength; i++ {
		aCopy[startIndex + i] = bHolder[i]
		bCopy[startIndex + i] = aHolder[i]
	}

	return aCopy, bCopy
}

// StrategySwapperIgnorant will perform crossover regardless of size
func StrategySwapperIgnorant(a []Strategy, b []Strategy, swapLength int, startIndex int) ([]Strategy, []Strategy) {
	if a == nil || b == nil{
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
		if startIndex + swapLength >= len(a) {
			startIndex = 0
		}
		aCopy = make([]Strategy, len(a))
		bCopy = make([]Strategy, len(a))
		aHolder = make([]Strategy, swapLength)
		bHolder = make([]Strategy, swapLength)
		copy(aCopy, a)
		copy(bCopy, b)

		for i:= 0; i < swapLength; i++ {
			aHolder[i] = a[i + startIndex]
		}
		for i:= 0; i < swapLength; i++ {
			bHolder[i] = b[i + startIndex]
		}
	} else {
		if swapLength > len(b) {
			swapLength = len(b)
		}
		if startIndex + swapLength >= len(b) {
			startIndex = 0
		}
		aCopy = make([]Strategy, len(b))
		bCopy = make([]Strategy, len(b))
		aHolder = make([]Strategy, swapLength)
		bHolder = make([]Strategy, swapLength)
		copy(aCopy, a)
		copy(bCopy, b)

		for i:= 0; i < len(aCopy); i++ {
			aHolder[i] = a[i + startIndex]
		}
		for i:= 0; i < len(bCopy); i++ {
			bHolder[i] = b[i + startIndex]
		}
	}

	for i:= 0; i < swapLength; i++ {
		aCopy[startIndex+i] = bHolder[i]
		bCopy[startIndex+i] = aHolder[i]
	}

	return aCopy, bCopy
}

// Mutate will mutate the strategy in a given individual
func Mutate(individual Individual, params EvolutionParams) (Individual, error) {
	if individual.id == "" {
		return Individual{}, fmt.Errorf("crossover | individual1 - individual id cannot be empty")
	}
	if individual.strategy == nil {
		return Individual{}, fmt.Errorf("crossover | individual1 - strategy array cannot be nil")
	}
	if len(individual.strategy) == 0 {
		return Individual{}, fmt.Errorf("crossover | individual1 - strategy array cannot be empty")
	}
	if individual.hasCalculatedFitness == false {
		return Individual{}, fmt.Errorf("crossover | individual1 - hasCalculatedFitness should be true")
	}
	if individual.hasAppliedStrategy == false {
		return Individual{}, fmt.Errorf("crossover | individual1 - hasAppliedStrategy should be true")
	}
	if len(params.Strategies) < 1 {
		return Individual{}, fmt.Errorf("crossover | params.Strategies cannot be length 0")
	}

	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(len(individual.strategy))

	rand.Seed(time.Now().UnixNano())
	randomStrategy := rand.Intn(len(params.Strategies))
	individual.strategy[randIndex] = params.Strategies[randomStrategy]

	return individual, nil
}

func GenerateIndividualID(identifier string, individualKind int) string {
	return fmt.Sprintf("%s-%s-%s%s", "individual", KindToString(individualKind), RandString(3), identifier)
}

func GenerateRandomIndividuals(number int, idTemplate string, kind int, strategyLength int,
	maxNumberOfStrategies int, availableStrategies []Strategy, depth int, terminals SymbolicExpressionSet,
	nonTerminals SymbolicExpressionSet, enforceIndependentVariable bool) ([]*Individual, error) {
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
		var tree *DualTree
		var err error
		var retry bool = true
		if enforceIndependentVariable {
			for retry {
				tree, err = GenerateRandomTreeEnforceIndependentVariable(depth, terminals[0], terminals, nonTerminals)
				if err != nil {

				}
				retry = false
			}

		} else {
			for retry {
				tree, err = GenerateRandomTree(depth, terminals, nonTerminals)
				if err != nil {
				}
				retry = false
			}
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
