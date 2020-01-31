package evolution

// CleansePopulation removes the trees from the population and refits them with the starter Tree.
func CleansePopulation(individuals []*Individual, treeReplacer DualTree) ([]*Individual, error) {
	for i := range individuals {
		if individuals[i].Kind == IndividualAntagonist {
			tree, err := treeReplacer.Clone()
			if err != nil {
				return nil, err
			}
			individual := individuals[i]
			if individual.Program == nil {
				individual.Program = &Program{}
			}
			newIndividual := individual.CloneWithTree(tree)
			newIndividual.Fitness = make([]float64, 0)
			newIndividual.Deltas = make([]float64, 0)
			newIndividual.HasCalculatedFitness = false
			newIndividual.HasAppliedStrategy = false
			newIndividual.AverageFitness = -1
			newIndividual.AverageDelta = -1
			newIndividual.BestFitness = -1
			newIndividual.BestDelta = -1
			newIndividual.FitnessVariance = 0
			newIndividual.FitnessStdDev = 0
			newIndividual.Program.T = &tree
			newIndividual.Strategy = individuals[i].Strategy
			individuals[i] = &newIndividual
		} else {
			newIndividual, err := individuals[i].Clone()
			if err != nil {
				return nil, err
			}
			newIndividual.Fitness = make([]float64, 0)
			newIndividual.Deltas = make([]float64, 0)
			newIndividual.FitnessVariance = 0
			newIndividual.FitnessStdDev = 0
			newIndividual.HasCalculatedFitness = false
			newIndividual.HasAppliedStrategy = false
			newIndividual.AverageFitness = -1
			newIndividual.AverageFitness = -1
			newIndividual.AverageDelta = -1
			newIndividual.BestFitness = -1
			newIndividual.BestDelta = -1
			newIndividual.Strategy = individuals[i].Strategy
			individuals[i] = &newIndividual
			individuals[i].Program = &Program{}
		}
	}
	return individuals, nil
}
