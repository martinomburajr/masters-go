package evolution

import (
	"fmt"
	"math/rand"
)

/**
Selection is the stage of a genetic algorithm in which individual genomes are chosen from a population for later breeding (using the crossover operator).

	A generic selection procedure may be implemented as follows:

	1. The fitness function is evaluated for each individual, providing fitness values,
	which are then normalized. Normalization means dividing the fitness value of each individual by the sum of all fitness values, so that the sum of all resulting fitness values equals 1.
	2. The population is sorted by descending fitness values.
	3. Accumulated normalized fitness values are computed: the accumulated fitness value of an individual is the sum of its
	own fitness value plus the fitness values of all the previous individuals; the accumulated fitness of the last individual should be 1, otherwise something went wrong in the normalization step.
	4. A random number R between 0 and 1 is chosen.
	5. The selected individual is the last one whose accumulated normalized value is greater than or equal to R.
	For a large number of individuals the above algorithm might be computationally quite demanding. A simpler and faster alternative uses the so-called stochastic acceptance.

	//https://en.wikipedia.org/wiki/Selection_(genetic_algorithm)
*/

const (
	ParentSelectionTournament           = 0 // ID for Tournament Selection
	ParentSelectionElitism              = 1 //ID for elitism
	ParentSelectionFitnessProportionate = 2
)

// TournamentSelection is a process whereby a random set of individuals from the population are selected,
// and the best in that sample succeed onto the next generation
func TournamentSelection(population []*Individual, tournamentSize int) ([]*Individual, error) {
	if population == nil {
		return nil, fmt.Errorf("tournament population cannot be nil")
	}
	if len(population) < 1 {
		return nil, fmt.Errorf("tournament population cannot have len < 1")
	}
	if tournamentSize < 1 {
		return nil, fmt.Errorf("tournament size cannot be less than 1")
	}

	// do
	newPop := make([]*Individual, len(population))

	for i := 0; i < len(population); i++ {
		randSelectedIndividuals := getNRandom(population, tournamentSize)
		fittest := tournamentSelect(randSelectedIndividuals)

		newPop[i] = fittest
	}

	return newPop, nil
}

// getNRandom selects  a random group of individiduals equivalent to the tournamentSize
func getNRandom(population []*Individual, tournamentSize int) []*Individual {
	newPop := make([]*Individual, tournamentSize)
	for i := 0; i < tournamentSize; i++ {
		randIndex := rand.Intn(len(population))
		newPop[i] = population[randIndex]
	}

	return newPop
}

//tournamentSelect returns the fittest individual in a given tournament
func tournamentSelect(selectedIndividuals []*Individual) *Individual {
	fittest := selectedIndividuals[0]
	for i := range selectedIndividuals {
		if selectedIndividuals[i].totalFitness > fittest.totalFitness {
			fittest = selectedIndividuals[i]
		}
	}
	return fittest
}

// Elitism is an evolutionary process where only the top (
// n) individuals based on eliteCount are selected based on their fitness.
// In essence it ranks the individuals based on fitness, then returns the top (n)
func Elitism(population []*Individual, elitePercentage float32) ([]*Individual, error) {
	return nil, nil
}

// Fitness Proportionate Selection is one of the most popular ways of parent selection.
// In this every individual can become a parent with a probability which is proportional to its fitness.
// Therefore, fitter individuals have a higher chance of mating and propagating their features to the next generation.
// Therefore, such a selection strategy applies a selection pressure to the more fit individuals in the population, evolving better individuals over time.
func FitnessProportionateSelection(population []*Individual) ([]*Individual, error) {
	return nil, nil
}
