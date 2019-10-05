package evolution

import (
	"math/rand"
)

// JudgementDay represents a moment where all individuals have completed their epoch phase and are waiting a decision
// onto who proceeds to the next generation. Judgement day is a compound function or abstraction that includes the
// following processes.
// 1. Parent Selection
// 2. Reproduction (via CrossoverTree)
// 3. Mutation (low probability)
// 4. Survivor Selection
// 5. Statistical Output
// 6. FinalPopulation configuration (incrementing age, clearing fitness values for old worthy individuals)
func JudgementDay(incomingPopulation []*Individual, opts EvolutionParams) ([]*Individual, error) {
	survivors := make([]*Individual, len(incomingPopulation))
	// Parent Selection
	// Tournament Selection
	outgoingParents, err := TournamentSelection(incomingPopulation, opts.TournamentSize)
	if err != nil {
		return nil, err
	}

	// Reproduction
	// CrossoverTree
	children := make([]Individual, opts.EachPopulationSize)
	for i := 0; i < len(outgoingParents); i += 2 {
		child1, child2, err := Crossover(*outgoingParents[i], *outgoingParents[i+1], opts)
		if err != nil {
			return nil, err
		}
		children[i] = child1
		children[i+1] = child2
	}

	// Reproduction
	// Mutation

	parentPopulationSize := int(opts.SurvivorPercentage * float32(opts.EachPopulationSize))
	childPopulationSize := opts.EachPopulationSize - parentPopulationSize

	// Reproduction
	// Mutation
	dualStrategies := append(opts.AntagonistAvailableStrategies, opts.ProtagonistAvailableStrategies...)
	// parents
	for i := 0; i < (len(outgoingParents)); i++ {

		probabilityOfMutation := rand.Float32()
		if probabilityOfMutation < opts.ProbabilityOfMutation {
			err := outgoingParents[i].Mutate(dualStrategies)
			if err != nil {
				return nil, err
			}
		}
	}

	// childs
	for i := 0; i < (len(children)); i++ {

		probabilityOfMutation := rand.Float32()
		if probabilityOfMutation < opts.ProbabilityOfMutation {
			err := children[i].Mutate(dualStrategies)
			if err != nil {
				return nil, err
			}
		}
	}

	// CHANGE - This only selects the first N parents

	for i := 0; i < parentPopulationSize; i++ {
		survivors[i] = outgoingParents[i]
	}
	for i := parentPopulationSize; i < parentPopulationSize+childPopulationSize; i++ {
		survivors[i] = outgoingParents[i]
	}

	// Survivor Selection

	// Statistical Output

	// Anointing Final Population and Return
	//individuals, err := CleansePopulation(survivors, *opts.StartIndividual.T)
	//if err != nil {
	//	return nil, err
	//}

	return survivors, nil
}

// CleansePopulation removes the trees from the population and refits them with the starter tree.
func CleansePopulation(individuals []*Individual, treeReplacer DualTree) []*Individual {
	for i := range individuals {
		if individuals[i].kind == IndividualAntagonist {
			newIndividual := individuals[i].CloneWithTree(treeReplacer)
			individuals[i] = &newIndividual
		} else {
			newIndividual := individuals[i].CloneWithTree(treeReplacer)
			individuals[i] = &newIndividual
			individuals[i].Program.T = nil
		}
	}
	return individuals
}
