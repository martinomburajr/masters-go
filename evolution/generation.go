package evolution

import "github.com/martinomburajr/masters-go/program"

type Generation struct {
	GenerationID       int
	PreviousGeneration *Generation
	NextGeneration     *Generation
	Protagonists       []*program.Program //Protagonists in a given generation
	Antagonists        []*program.Program //Antagonists in a given generation
	FittestProtagonist *program.Program
	FittestAntagonist  *program.Program
	engine             *EvolutionEngine // Reference to Engine
}

// Next returns the next generation
func (g *Generation) Next() *Generation {
	return g.NextGeneration
}

// Engine returns a reference to the Evolution Engine in use
func (g *Generation) Engine() *EvolutionEngine {
	return g.Engine()
}

// Previous returns the previous generation
func (g *Generation) Previous() *Generation {
	return g.PreviousGeneration
}

// Start begins the generational evolutionary cycle.
// It creates a new generation that it links the {NextGeneration} field to. Similar to the way a LinkedList works
func (g *Generation) Start() *Generation {
	return g.PreviousGeneration
}

// Restart is similar to StartHOG but it restarts the evolutionary process from the selected Generation.
// All future generations are deleted to make way for this evolutionary process
func (g *Generation) Restart() *Generation {
	return g.PreviousGeneration
}

// StartHOG is a unique version of start. It clears future history and jumps to a given generation,
// inserts generational material into the generation, and creates a new evolutionary propagation from it.
func (g *Generation) StartHOG(gen Generation) *Generation {
	return g.PreviousGeneration
}

