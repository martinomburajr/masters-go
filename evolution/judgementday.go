package evolution

import (
	"fmt"
	"math/rand"
)

// JudgementDay represents a moment where all individuals have completed their epoch phase and are waiting a decision
// onto who proceeds to the next Generation.
// Parent Selection has already been performed Judgement day is a compound function or abstraction that includes the
// following processes.
// 2. Reproduction (via CrossoverTree)
// 3. Mutation (low probability)
// 4. Survivor Selection
// 5. Statistical Output
// 6. FinalPopulation configuration (incrementing Age, clearing Fitness values for old worthy individuals)
func JudgementDay(incomingPopulation []*Individual, kind int, generationCount int,
	opts EvolutionParams) ([]*Individual, error) {
	survivors := make([]*Individual, len(incomingPopulation))

	// Reproduction
	// CrossoverTree
	children := make([]*Individual, opts.EachPopulationSize)
	for i := 0; i < len(incomingPopulation); i += 2 {
		child1, child2, err := Crossover(*incomingPopulation[i], *incomingPopulation[i+1], opts)
		if err != nil {
			return nil, err
		}
		child1.BirthGen = generationCount + 1
		child2.BirthGen = generationCount + 1
		children[i] = &child1
		children[i+1] = &child2
	}

	// Reproduction
	// Mutation

	parentPopulationSize := int(opts.SurvivorPercentage * float64(opts.EachPopulationSize))
	childPopulationSize := opts.EachPopulationSize - parentPopulationSize

	// Reproduction
	// Mutation
	incomingPopulation, children, err := Mutate(incomingPopulation, children, kind, opts)
	if err != nil {
		return nil, err
	}

	//TODO CHANGE - This only selects the first N parents
	for i := 0; i < parentPopulationSize; i++ {
		survivors[i] = incomingPopulation[i]
	}
	for i := 0; i < childPopulationSize; i++ {
		survivors[i+parentPopulationSize] = children[i]
	}

	// Survivor Selection
	return survivors, nil
}

func Mutate(outgoingParents []*Individual, children []*Individual, kind int,
	opts EvolutionParams) (parents []*Individual, childs []*Individual, err error) {
	if kind == IndividualAntagonist {
		for i := 0; i < (len(outgoingParents)); i++ {
			probabilityOfMutation := rand.Float64()
			if probabilityOfMutation < opts.ProbabilityOfMutation {
				err := outgoingParents[i].Mutate(opts.AntagonistAvailableStrategies)
				if err != nil {
					return nil, nil, err
				}
			}
		}
		// childs
		for i := 0; i < (len(children)); i++ {
			probabilityOfMutation := rand.Float64()
			if probabilityOfMutation < opts.ProbabilityOfMutation {
				err := children[i].Mutate(opts.AntagonistAvailableStrategies)
				if err != nil {
					return nil, nil, err
				}
			}
		}
	} else if kind == IndividualProtagonist {
		for i := 0; i < (len(outgoingParents)); i++ {
			probabilityOfMutation := rand.Float64()
			if probabilityOfMutation < opts.ProbabilityOfMutation {
				err := outgoingParents[i].Mutate(opts.ProtagonistAvailableStrategies)
				if err != nil {
					return nil, nil, err
				}
			}
		}
		// childs
		for i := 0; i < (len(children)); i++ {
			probabilityOfMutation := rand.Float64()
			if probabilityOfMutation < opts.ProbabilityOfMutation {
				err := children[i].Mutate(opts.ProtagonistAvailableStrategies)
				if err != nil {
					return nil, nil, err
				}
			}
		}
	} else {
		return nil, nil, fmt.Errorf("Judgement Day | Invalid kind")
	}
	return outgoingParents, children, nil
}

// CleansePopulation removes the trees from the population and refits them with the starter Tree.
func CleansePopulation(individuals []*Individual, treeReplacer DualTree) ([]*Individual, error) {
	for i := range individuals {
		if individuals[i].Kind == IndividualAntagonist {
			tree, err := treeReplacer.Clone()
			if err != nil {
				return nil, err
			}
			newIndividual := individuals[i].CloneWithTree(tree)
			newIndividual.Fitness = make([]float64, 0)
			newIndividual.HasCalculatedFitness = false
			newIndividual.HasAppliedStrategy = false
			newIndividual.TotalFitness = 0
			newIndividual.Program.T = &tree
			newIndividual.Strategy = individuals[i].Strategy
			individuals[i] = &newIndividual
		} else {
			newIndividual, err := individuals[i].Clone()
			if err != nil {
				return nil, err
			}
			newIndividual.Fitness = make([]float64, 0)
			newIndividual.HasCalculatedFitness = false
			newIndividual.HasAppliedStrategy = false
			newIndividual.TotalFitness = 0
			newIndividual.Strategy = individuals[i].Strategy
			individuals[i] = &newIndividual
			individuals[i].Program = &Program{}
		}
	}
	return individuals, nil
}
