package evolution

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
)

const (
	CrossoverSinglePoint = "CrossoverSinglePoint"
	CrossoverFixedPoint = "CrossoverFixedPoint"
	CrossoverKPoint = "CrossoverKPoint"
	CrossoverUniform = "CrossoverUniform"
)

// CrossoverSinglePoint performs a single-point crossover that is dictated by the crossover percentage float.
// Both parent chromosomes are split at the percentage section specified by crossoverPercentage
func SinglePointCrossover(parentA, parentB *Individual) (childA Individual,
	childB Individual,
	err error) {
	// Require
	if parentA.Strategy == nil {
		return Individual{}, Individual{}, fmt.Errorf("parentA strategy cannot be nil")
	}
	if len(parentA.Strategy) < 1 {
		return Individual{}, Individual{}, fmt.Errorf("parentA strategy cannot be empty")
	}
	if parentB.Strategy == nil {
		return Individual{}, Individual{}, fmt.Errorf("parentB strategy cannot be nil")
	}
	if len(parentB.Strategy) < 1 {
		return Individual{}, Individual{}, fmt.Errorf("parentB strategy cannot be empty")
	}

	// DO
	switch parentA.Kind {
	case IndividualAntagonist:
		childA.Id = GenerateIndividualID("", IndividualAntagonist)
		childB.Id = GenerateIndividualID("", IndividualAntagonist)
	case IndividualProtagonist:
		childA.Id = GenerateIndividualID("", IndividualProtagonist)
		childB.Id = GenerateIndividualID("", IndividualProtagonist)
	}

	childA.Strategy = parentA.Strategy
	childB.Strategy = parentB.Strategy
	childA.Age = 0
	childB.Age = 0

	mut := sync.Mutex{}
	mut.Lock()
	if len(parentA.Strategy) >= len(parentB.Strategy) {
		prob := rand.Intn(len(parentB.Strategy))
		for i := 0; i < prob; i++ {
			childA.Strategy[i] = parentB.Strategy[i]
			childB.Strategy[i] = parentA.Strategy[i]
		}
	}  else {
		prob := rand.Intn(len(parentA.Strategy))
		for i := 0; i < prob; i++ {
			childA.Strategy[i] = parentB.Strategy[i]
			childB.Strategy[i] = parentA.Strategy[i]
		}
	}
	mut.Unlock()

	return childA, childB, nil
}

// CrossoverSinglePoint performs a single-point crossover that is dictated by the crossover percentage float.
// Both parent chromosomes are split at the percentage section specified by crossoverPercentage
func KPointCrossover(parentA, parentB *Individual, kPoint int) (childA Individual,
	childB Individual,
	err error) {

	// Require
	if parentA.Strategy == nil {
		return Individual{}, Individual{}, fmt.Errorf("parentA strategy cannot be nil")
	}
	if len(parentA.Strategy) < 1 {
		return Individual{}, Individual{}, fmt.Errorf("parentA strategy cannot be empty")
	}
	if parentB.Strategy == nil {
		return Individual{}, Individual{}, fmt.Errorf("parentB strategy cannot be nil")
	}
	if len(parentB.Strategy) < 1 {
		return Individual{}, Individual{}, fmt.Errorf("parentB strategy cannot be empty")
	}

	//DO
	switch parentA.Kind {
	case IndividualAntagonist:
		childA.Id = GenerateIndividualID("", IndividualAntagonist)
		childB.Id = GenerateIndividualID("", IndividualAntagonist)
	case IndividualProtagonist:
		childA.Id = GenerateIndividualID("", IndividualProtagonist)
		childB.Id = GenerateIndividualID("", IndividualProtagonist)
	}

	childA.Strategy = parentA.Strategy
	childB.Strategy = parentB.Strategy
	childA.Age = 0
	childB.Age = 0


	mut := sync.Mutex{}
	mut.Lock()
	// Swap every element
	if kPoint < 1 || (kPoint > len(parentA.Strategy) && kPoint > len(parentB.Strategy)){
		if len(parentA.Strategy) >= len(parentB.Strategy) {
			for i := 0; i < len(parentB.Strategy); i++ {
				if i % 2 == 0 {
					childA.Strategy[i] = parentA.Strategy[i]
					childB.Strategy[i] = parentB.Strategy[i]
				}else {
					childA.Strategy[i] = parentB.Strategy[i]
					childB.Strategy[i] = parentA.Strategy[i]
				}
			}
		} else {
			for i := 0; i < len(parentA.Strategy); i++ {
				if i % 2 == 0 {
					childA.Strategy[i] = parentA.Strategy[i]
					childB.Strategy[i] = parentB.Strategy[i]
				}else {
					childA.Strategy[i] = parentB.Strategy[i]
					childB.Strategy[i] = parentA.Strategy[i]
				}
			}
		}
	}else {
		// USe the smaller chromosome as reference for K. Randomly select K points on the smaller one.
		if len(parentA.Strategy) >= len(parentB.Strategy) {
			kPoints := rand.Perm(kPoint)
			sort.Ints(kPoints)

			shouldSwap := true
			for i := 0; i < len(parentB.Strategy); i++ {
				for j := range kPoints {
					if i == kPoints[j] {
						shouldSwap = !shouldSwap
					}
					if shouldSwap {
						childA.Strategy[i] = parentB.Strategy[i]
						childB.Strategy[i] = parentA.Strategy[i]
					}else {
						childA.Strategy[i] = parentA.Strategy[i]
						childB.Strategy[i] = parentB.Strategy[i]
					}
				}
			}
		} else {
			kPoints := rand.Perm(kPoint)
			sort.Ints(kPoints)

			shouldSwap := true
			for i := 0; i < len(parentA.Strategy); i++ {
				for j := range kPoints {
					if i == kPoints[j] {
						shouldSwap = !shouldSwap
					}
					if shouldSwap {
						childA.Strategy[i] = parentB.Strategy[i]
						childB.Strategy[i] = parentA.Strategy[i]
					}else {
						childA.Strategy[i] = parentA.Strategy[i]
						childB.Strategy[i] = parentB.Strategy[i]
					}
				}
			}
		}
	}

	mut.Unlock()
	return childA, childB, nil
}


// CrossoverSinglePoint performs a single-point crossover that is dictated by the crossover percentage float.
// Both parent chromosomes are split at the percentage section specified by crossoverPercentage
func UniformCrossover(parentA, parentB *Individual) (childA Individual,
	childB Individual,
	err error) {
	// Require
	if parentA.Strategy == nil {
		return Individual{}, Individual{}, fmt.Errorf("parentA strategy cannot be nil")
	}
	if len(parentA.Strategy) < 1 {
		return Individual{}, Individual{}, fmt.Errorf("parentA strategy cannot be empty")
	}
	if parentB.Strategy == nil {
		return Individual{}, Individual{}, fmt.Errorf("parentB strategy cannot be nil")
	}
	if len(parentB.Strategy) < 1 {
		return Individual{}, Individual{}, fmt.Errorf("parentB strategy cannot be empty")
	}

	// DO
	switch parentA.Kind {
	case IndividualAntagonist:
		childA.Id = GenerateIndividualID("", IndividualAntagonist)
		childB.Id = GenerateIndividualID("", IndividualAntagonist)
	case IndividualProtagonist:
		childA.Id = GenerateIndividualID("", IndividualProtagonist)
		childB.Id = GenerateIndividualID("", IndividualProtagonist)
	}

	childA.Strategy = parentA.Strategy
	childB.Strategy = parentB.Strategy
	childA.Age = 0
	childB.Age = 0

	mut := sync.Mutex{}
	mut.Lock()
	if len(parentA.Strategy) >= len(parentB.Strategy) {
		for i := 0; i < len(parentB.Strategy); i++ {
			prob := rand.Intn(2)
			if prob == 0 {
				childA.Strategy[i] = parentA.Strategy[i]
				childB.Strategy[i] = parentB.Strategy[i]
			}else {
				childA.Strategy[i] = parentB.Strategy[i]
				childB.Strategy[i] = parentA.Strategy[i]
			}
		}
	}  else {
		for i := 0; i < len(parentA.Strategy); i++ {
			prob := rand.Intn(2)
			if prob == 0 {
				childA.Strategy[i] = parentA.Strategy[i]
				childB.Strategy[i] = parentB.Strategy[i]
			}else {
				childA.Strategy[i] = parentB.Strategy[i]
				childB.Strategy[i] = parentA.Strategy[i]
			}
		}
	}
	mut.Unlock()

	return childA, childB, nil
}

// FixedPointCrossover will perform crossover on the strategies of a given  set of individuals
func FixedPointCrossover(individual Individual, individual2 Individual, params EvolutionParams) (Individual, Individual,
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
	if individual.Program == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual1 - program cannot be nil")
	}
	if individual.Program.T == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual1 - program Tree cannot be nil")
	}
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
	if individual2.Program == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual2 - program cannot be nil")
	}
	if individual2.Program.T == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover | individual2 - program Tree cannot be nil")
	}

	individual1Len := len(individual.Strategy)
	individual2Len := len(individual2.Strategy)

	child1, err := individual.Clone()
	child1.Id = child1.Id + "c1"
	if err != nil {
		return Individual{}, Individual{}, err
	}
	child2, err := individual2.Clone()
	child2.Id = child2.Id + "c2"
	if err != nil {
		return Individual{}, Individual{}, err
	}

	crossoverPercentage := params.Reproduction.CrossoverPercentage
	if crossoverPercentage == 0 {
		return child1, child2, err
	}
	if crossoverPercentage == 1 {
		return child2, child1, err
	}

	individual1ChunkSize := int(math.Ceil(float64(individual1Len) * float64(crossoverPercentage)))
	individual2ChunkSize := int(float64(individual2Len) * crossoverPercentage)

	if individual1ChunkSize >= individual2ChunkSize {
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
	}
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



func depthPenaltyIgnore(maxDepth int, individual1Depth int, individual2Depth int) (int, int) {
	if maxDepth < 0 {
		maxDepth = 0
	}
	var individual1DepthRemainderFromMaX, individual2DepthRemainderFromMax int
	if individual1Depth >= maxDepth {
		individual1DepthRemainderFromMaX = 0
	} else {
		individual1DepthRemainderFromMaX = maxDepth - individual1Depth
	}
	if individual2Depth >= maxDepth {
		individual2DepthRemainderFromMax = 0
	} else {
		individual2DepthRemainderFromMax = maxDepth - individual2Depth
	}
	return individual1DepthRemainderFromMaX, individual2DepthRemainderFromMax
}
