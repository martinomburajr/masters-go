package evolution

/**
EvolutionProcess represents the state of an evolutionary process once the evolution engine starts
*/
type EvolutionProcess struct {
	currentGeneration *Generation
	engine            *EvolutionEngine
}
