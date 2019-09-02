package evolution

import "fmt"

type Generation struct {
	GenerationID                 string
	PreviousGeneration           *Generation
	nextGeneration               *Generation
	Protagonists                 []*Individual //Protagonists in a given generation
	Antagonists                  []*Individual //Antagonists in a given generation
	engine                       *EvolutionEngine
	isComplete                   bool
	hasParentSelectionHappened   bool
	hasSurvivorSelectionHappened bool
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
func (g *Generation) Start() (*Generation, error) {
	setupEpochs, err := g.setupEpochs()
	if err != nil {
		return nil, err
	}

	completedEpochs, err := g.runEpochs(setupEpochs)
	if err != nil {
		return nil, err
	}

	for i := range completedEpochs {
		AggregateFitness(completedEpochs[i])
	}

	// perform parent selection
	switch g.engine.parentSelection {
	case ParentSelectionTournament:
		TournamentSelection(comp)
	}

	//perform survivor selection

	//return new generation
	g.nextGeneration = nil
	return g.nextGeneration, nil
}

// setupEpochs takes in the generation individuals (
// protagonists and antagonists) and creates a set of uninitialized epochs
func (g *Generation) setupEpochs() ([]Epoch, error){
	if g.Antagonists == nil {
		return nil, fmt.Errorf("antagonists cannot be nil in generation")
	}
	if g.Protagonists == nil {
		return nil, fmt.Errorf("protagonists cannot be nil in generation")
	}
	if len(g.Antagonists) < 1 {
		return nil, fmt.Errorf("antagonists cannot be empty")
	}
	if len(g.Protagonists) < 1 {
		return nil, fmt.Errorf("protagonists cannot be empty")
	}

	epochs := make([]Epoch, len(g.Antagonists) * len(g.Protagonists))
	count := 0
	for _, antagonist := range g.Antagonists {
		for _, protagonist := range g.Protagonists {
			epochs[count] = Epoch{
				isComplete:false,
				protagonistBegins:false,
				terminalSet:g.engine.availableTerminalSet,
				nonTerminalSet:g.engine.availableNonTerminalSet,
				hasAntagonistApplied:false,
				hasProtagonistApplied:false,
				probabilityOfMutation:g.engine.probabilityOfMutation,
				probabilityOfNonTerminalMutation:g.engine.probabilityOfNonTerminalMutation,
				antagonist:antagonist,
				protagonist:protagonist,
				generation:g,
				program:g.engine.startIndividual,
				id: CreateEpochID(count, g.GenerationID, antagonist.id, protagonist.id),
			}
			count++
		}
	}
	return epochs, nil
}


// CurrentPopulation retrieves the current population of the given generation.
// Individuals may have competed or may have been altered in a variety of ways.
// This will return a list of references Individuals
func (g *Generation) CurrentPopulation() ([]*Individual, error) {
	return nil, nil
}

// runEpoch begins the run of a single epoch
func (g *Generation) runEpochs(epochs []Epoch) ([]Epoch, error) {
	if epochs == nil {
		return nil, fmt.Errorf("epochs have not been initialized | epochs is nil")
	}

	for i := range epochs {
		err := epochs[i].Start()
		if err != nil {
			return nil, err
		}
	}

	return epochs, nil
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

// Compete gives protagonist and anatagonists the chance to compete. A competition involves an epoch,
// that returns the result of the epoch.
func (g *Generation) Compete() error {
	//for _, epoch := range g.Epochs {
	//	err := epoch.Start()
	//	if err != nil {
	//		return err
	//	}
	//	if epoch.hasAntagonistApplied && epoch.hasProtagonistApplied {
	//		continue
	//	} else {
	//		return fmt.Errorf("epoch completed but antagonist and/or protagonist not applied %#v, ", err)
	//	}
	//}
	//g.isComplete = true
	return nil
}

// ApplyParentSelection takes in a given generation and returns a set of individuals once the preselected parent
// selection strategy has been applied to the generation.
// These individuals are ready to be taken to either a new generation or preferably through survivor selection in the
// case you do not wantAntagonist the population to grow in size.
func (g *Generation) ApplyParentSelection() ([]*Individual, error) {
	if !g.isComplete {
		return nil, fmt.Errorf("generation #id: %s has not competed, ", g.GenerationID)
	}

	currentPopulation, err := g.CurrentPopulation()
	if err != nil {
		return nil, err
	}

	switch g.engine.parentSelection {
	case ParentSelectionTournament:
		selectedInvididuals, err := TournamentSelection(currentPopulation)
		if err != nil {
			return nil, err
		}
		g.hasParentSelectionHappened = true
		return selectedInvididuals, nil
	case ParentSelectionElitism:
		selectedInvididuals, err := Elitism(currentPopulation, g.engine.elitismPercentage)
		if err != nil {
			return nil, err
		}
		g.hasParentSelectionHappened = true
		return selectedInvididuals, nil
	default:
		return nil, fmt.Errorf("no appropriate parent selection strategy selected. See parentselection." +
			"go file for information on integer values that represent different parent selection strategies")
	}
}

// ApplySurvivorSelection applies the preselected survivor selection strategy.
// It DOES NOT check to see if the parent selection has already been applied,
// as in some cases evolutionary programs may choose to run without the parent selection phase.
// The onus is on the evolutionary architect to keep this consideration in mind.
func (g *Generation) ApplySurvivorSelection() ([]*Individual, error) {
	if !g.isComplete {
		return nil, fmt.Errorf("generation #id: %s has not competed, ", g.GenerationID)
	}

	return nil, nil
}

type GenerationResult struct {
	generation *Generation
}

// RunNext takes in a current GenerationResult runs a set of parent and survivor selection mechanisms,
// and returns the new generation
func (g *GenerationResult) RunNext() *GenerationResult {
	return nil
}

type GenerationEngine struct {
}
