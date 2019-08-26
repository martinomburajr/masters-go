package evolution

import "fmt"

type Generation struct {
	GenerationID       int
	Epochs             []*Epoch
	PreviousGeneration *Generation
	nextGeneration     *Generation
	Protagonists       []*Individual //Protagonists in a given generation
	Antagonists        []*Individual //Antagonists in a given generation
	FittestProtagonist *Program
	FittestAntagonist  *Program
	engine *EvolutionEngine
}

// Next returns the next generation
func (g *Generation) Next() *Generation {
	return g.nextGeneration
}


// Previous returns the previous generation
func (g *Generation) Previous() *Generation {
	return g.PreviousGeneration
}

// Start begins the generational evolutionary cycle.
// It creates a new generation that it links the {nextGeneration} field to. Similar to the way a LinkedList works
func (g *Generation) Start() (*GenerationResult, error) {
	g.setupEpochs()

	return g.runEpochs()
}

// setupEpochs initializes the epochs based on the information retrieved from the EvolutionEngine.
// Each epoch has an M x N pairing of individuals.
func (g *Generation) setupEpochs()  {
	g.Epochs = make([]*Epoch, len(g.Protagonists) * len(g.Antagonists))

	for antagonistIndex := range g.Antagonists {
		for protagonistIndex := range g.Protagonists {
			epoch := Epoch{}.
				SetProbabilityOfMutation(g.engine.probabilityOfMutation).
				SetProbabilityOfNonTerminalMutation(g.engine.probabilityOfNonTerminalMutation).
				SetAntagonist(g.Antagonists[antagonistIndex]).
				SetProtagonist(g.Protagonists[protagonistIndex]).
				SetProgram(g.engine.StartIndividual())

			g.Epochs = append(g.Epochs, epoch)
		}
	}
}

// runEpoch begins the run of a single epoch
func (g *Generation) runEpochs() (*GenerationResult, error)  {
	if g.Epochs == nil {
		return nil, fmt.Errorf("epochs have not been initialized | g.Epochs is nil")
	}

	epochResults := make([]*EpochResult, len(g.Epochs))
	for i := range g.Epochs {
		epochResult, err := g.Epochs[i].InitSimulator().Start()
		if err != nil {
			return nil, err
		}
		epochResults[i] = epochResult
	}

	return &GenerationResult{
		epochResults: epochResults,
		generation: g,
	}, nil
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

type GenerationResult struct {
	epochResults []*EpochResult
	generation *Generation
}
