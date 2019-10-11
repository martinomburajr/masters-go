package evolution

const (
	SurvivorSelectionSteadyState  = 0
	SurvivorSelectionGenerational = 1
)

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
