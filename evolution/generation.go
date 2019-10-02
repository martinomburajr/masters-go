package evolution

import "fmt"

type Generation struct {
	GenerationID                 string
	Protagonists                 []*Individual //Protagonists in a given generation
	Antagonists                  []*Individual //Antagonists in a given generation
	engine                       *EvolutionEngine
	isComplete                   bool
	hasParentSelectionHappened   bool
	hasSurvivorSelectionHappened bool
	count                        int
}

// Start begins the generational evolutionary cycle.
// It creates a new generation that it links the {nextGeneration} field to. Similar to the way a LinkedList works
func (g *Generation) Start() (*Generation, error) {
	setupEpochs, err := g.setupEpochs()
	if err != nil {
		return nil, err
	}

	// Runs the epochs and returns completed epochs that contain fitness information within each individual.
	_, err = g.runEpochs(setupEpochs)
	if err != nil {
		return nil, err
	}

	// Calculate the fitness for individuals in the generation
	for i := range g.Protagonists {
		fitness, err := AggregateFitness(*g.Protagonists[i])
		if err != nil {
			return nil, err
		}
		g.Protagonists[i].totalFitness = fitness
		g.Protagonists[i].hasCalculatedFitness = true
	}

	// Calculate the fitness for individuals in the generation
	for i := range g.Antagonists {
		fitness, err := AggregateFitness(*g.Antagonists[i])
		if err != nil {
			return nil, err
		}
		g.Antagonists[i].totalFitness = fitness
		g.Antagonists[i].hasCalculatedFitness = true
	}

	//protagonists := make([]Individual, len(g.Protagonists))
	//for i := range g.Protagonists {
	//	protagonists[i] = *g.Protagonists[i]
	//}
	//
	//antagonists := make([]Individual, len(g.Antagonists))
	//for i := range g.Antagonists {
	//	antagonists[i] = *g.Antagonists[i]
	//}

	nextGenAntagonists, err := JudgementDay(g.Antagonists, g.engine.Parameters)
	if err != nil {
		return nil, err
	}

	nextGenProtagonists, err := JudgementDay(g.Protagonists, g.engine.Parameters)
	if err != nil {
		return nil, err
	}

	//protagonists := make([]*Individual, len(nextGenProtagonists))
	//for i := range nextGenProtagonists {
	//	protagonists[i] = &nextGenProtagonists[i]
	//}
	//
	//antagonists := make([]*Individual, len(nextGenAntagonists))
	//for i := range nextGenAntagonists {
	//	antagonists[i] = &nextGenAntagonists[i]
	//}

	nextGenID := GenerateGenerationID(g.count + 1)
	//return new generation
	return &Generation{
		Antagonists:                  nextGenAntagonists,
		Protagonists:                 nextGenProtagonists,
		engine:                       g.engine,
		GenerationID:                 nextGenID,
		hasSurvivorSelectionHappened: false,
		isComplete:                   false,
		hasParentSelectionHappened:   false,
		count:                        (g.count + 1),
	}, nil
}

func GenerateGenerationID(count int) string {
	return fmt.Sprintf("GEN-%d", count)
}

// setupEpochs takes in the generation individuals (
// protagonists and antagonists) and creates a set of uninitialized epochs
func (g *Generation) setupEpochs() ([]Epoch, error) {
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

	epochs := make([]Epoch, len(g.Antagonists)*len(g.Protagonists))
	count := 0
	for _, antagonist := range g.Antagonists {
		for _, protagonist := range g.Protagonists {
			epochs[count] = Epoch{
				isComplete:                       false,
				protagonistBegins:                false,
				terminalSet:                      g.engine.Parameters.TerminalSet,
				nonTerminalSet:                   g.engine.Parameters.NonTerminalSet,
				hasAntagonistApplied:             false,
				hasProtagonistApplied:            false,
				probabilityOfMutation:            g.engine.Parameters.ProbabilityOfMutation,
				probabilityOfNonTerminalMutation: g.engine.Parameters.ProbabilityOfNonTerminalMutation,
				antagonist:                       antagonist,
				protagonist:                      protagonist,
				generation:                       g,
				program:                          g.engine.StartIndividual,
				id:                               CreateEpochID(count, g.GenerationID, antagonist.id, protagonist.id),
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
// case you do not isEqual the population to grow in size.
func (g *Generation) ApplyParentSelection() ([]*Individual, error) {
	if !g.isComplete {
		return nil, fmt.Errorf("generation #id: %s has not competed, ", g.GenerationID)
	}

	currentPopulation, err := g.CurrentPopulation()
	if err != nil {
		return nil, err
	}

	switch g.engine.ParentSelection {
	case ParentSelectionTournament:
		selectedInvididuals, err := TournamentSelection(currentPopulation, g.engine.Parameters.TournamentSize)
		if err != nil {
			return nil, err
		}
		g.hasParentSelectionHappened = true
		return selectedInvididuals, nil
	case ParentSelectionElitism:
		selectedInvididuals, err := Elitism(currentPopulation, g.engine.ElitismPercentage)
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
