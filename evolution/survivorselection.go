package evolution

const (
	SurvivorSelectionFitnessBased = "SurvivorSelectionFitnessBased"
	SurvivorSelectionRandom = "SurvivorSelectionRandom"
)

// FitnessBasedSurvivorSelection returns a set of survivors proportionate to the survivor percentage.
// It orders some of the best parents and some of the best children based on the ratio
func FitnessBasedSurvivorSelection(selectedParents, selectedChildren []*Individual,
	params EvolutionParams) ([]*Individual, error) {
	survivors := make([]*Individual, params.EachPopulationSize)

	parentPopulationSize := int(params.Selection.Survivor.SurvivorPercentage * float64(params.EachPopulationSize))
	childPopulationSize := params.EachPopulationSize - parentPopulationSize

	sortedParents, err := SortIndividuals(selectedParents, true)
	if err != nil {
		return nil, err
	}
	sortedChildren, err := SortIndividuals(selectedChildren, true)
	if err != nil {
		return nil, err
	}

	for i := 0; i < parentPopulationSize; i++ {
		survivors[i] = sortedParents[i]
	}
	for i := 0; i < childPopulationSize; i++ {
		survivors[i+parentPopulationSize] = sortedChildren[childPopulationSize]
	}

	return survivors, nil
}

// RandomSurvivorSelection selects a random set of parents and a random set of children. The numbers are based on 
func RandomSurvivorSelection(selectedParents, selectedChildren []*Individual,
	params EvolutionParams) ([]*Individual, error) {
	survivors := make([]*Individual, params.EachPopulationSize)

	parentPopulationSize := int(params.Selection.Survivor.SurvivorPercentage * float64(params.EachPopulationSize))
	childPopulationSize := params.EachPopulationSize - parentPopulationSize

	sortedParents, err := SortIndividuals(selectedParents, true)
	if err != nil {
		return nil, err
	}
	sortedChildren, err := SortIndividuals(selectedChildren, true)
	if err != nil {
		return nil, err
	}

	for i := 0; i < parentPopulationSize; i++ {
		survivors[i] = sortedParents[i]
	}
	for i := 0; i < childPopulationSize; i++ {
		survivors[i+parentPopulationSize] = sortedChildren[childPopulationSize]
	}

	return survivors, nil
}

// GenerationalSurvivorSelection is a process where the entire input population gets replaced by their offspring.
// The returned individuals do not exist with their parents as they have been totally annihilated.
// These new individuals will go on into the next Generation
func GenerationalSurvivorSelection(population *Generation) ([]*Individual, error) {
	return nil, nil
}

// SteadyStateSurvivorSelection is a process where a select amount of individuals make it through.
// Some parents may make it through based on their Fitness or other compounding parameters.
// These new individuals will go on into the next Generation
func SteadyStateSurvivorSelection(population *Generation) ([]*Individual, error) {
	return nil, nil
}
