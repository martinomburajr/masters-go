package evolution

import (
	"fmt"
	"math/rand"
)

type Generation struct {
	GenerationID                 string
	Protagonists                 []*Individual //Protagonists in a given Generation
	Antagonists                  []*Individual //Antagonists in a given Generation
	engine                       *EvolutionEngine
	isComplete                   bool
	hasParentSelectionHappened   bool
	hasSurvivorSelectionHappened bool
	count                        int
}

// Start begins the generational evolutionary cycle.
// It creates a new Generation that it links the {nextGeneration} field to. Similar to the way a LinkedList works
func (g *Generation) Start(generationCount int) (*Generation, error) {
	setupEpochs, err := g.setupEpochs()
	if err != nil {
		return nil, err
	}

	// Runs the epochs and returns completed epochs that contain Fitness information within each individual.
	_, err = g.runEpochs(setupEpochs)
	if err != nil {
		return nil, err
	}

	// Calculate the Fitness for individuals in the Generation
	for i := range g.Protagonists {
		protagonistFitness, err := AggregateFitness(*g.Protagonists[i])
		if err != nil {
			return nil, err
		}
		g.Protagonists[i].TotalFitness = float64(protagonistFitness / float64(len(g.Protagonists[i].Fitness)))
		g.Protagonists[i].HasCalculatedFitness = true
		g.Protagonists[i].Age++

		antagonistFitness, err := AggregateFitness(*g.Antagonists[i])
		if err != nil {
			return nil, err
		}
		g.Antagonists[i].TotalFitness = float64(antagonistFitness / float64(len(g.Antagonists[i].Fitness)))
		g.Antagonists[i].HasCalculatedFitness = true
		g.Antagonists[i].Age++
	}

	nextGenAntagonists, err := JudgementDay(g.Antagonists, generationCount, g.engine.Parameters)
	if err != nil {
		return nil, err
	}

	nextGenProtagonists, err := JudgementDay(g.Protagonists, generationCount, g.engine.Parameters)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("#################### GEN: %d ANTAGONISTS ####################### \n", generationCount)
	//for _, g := range nextGenAntagonists {
	//	ss := g.ToString()
	//	fmt.Println(ss.String())
	//}
	//fmt.Printf("#################### GEN: %d PROTAGONISTS ####################### \n", generationCount)
	//for _, g := range nextGenProtagonists {
	//	ss := g.ToString()
	//	fmt.Println(ss.String())
	//}

	nextGenID := GenerateGenerationID(g.count + 1)
	nextGen := &Generation{
		Antagonists:                  nextGenAntagonists,
		Protagonists:                 nextGenProtagonists,
		engine:                       g.engine,
		GenerationID:                 nextGenID,
		hasSurvivorSelectionHappened: false,
		isComplete:                   false,
		hasParentSelectionHappened:   false,
		count:                        g.count + 1,
	}
	//return new Generation
	return nextGen, nil
}

func GenerateGenerationID(count int) string {
	return fmt.Sprintf("GEN-%d", count)
}

// setupEpochs takes in the Generation individuals (
// protagonists and antagonists) and creates a set of uninitialized epochs
func (g *Generation) setupEpochs() ([]Epoch, error) {
	if g.Antagonists == nil {
		return nil, fmt.Errorf("antagonists cannot be nil in Generation")
	}
	if g.Protagonists == nil {
		return nil, fmt.Errorf("protagonists cannot be nil in Generation")
	}
	if len(g.Antagonists) < 1 {
		return nil, fmt.Errorf("antagonists cannot be empty")
	}
	if len(g.Protagonists) < 1 {
		return nil, fmt.Errorf("protagonists cannot be empty")
	}

	epochs := make([]Epoch, len(g.Antagonists)*len(g.Protagonists))
	count := 0
	for i, _ := range g.Antagonists {
		for j, _ := range g.Protagonists {
			epochs[count] = Epoch{
				isComplete:                       false,
				protagonistBegins:                false,
				terminalSet:                      g.engine.Parameters.TerminalSet,
				nonTerminalSet:                   g.engine.Parameters.NonTerminalSet,
				hasAntagonistApplied:             false,
				hasProtagonistApplied:            false,
				probabilityOfMutation:            g.engine.Parameters.ProbabilityOfMutation,
				probabilityOfNonTerminalMutation: g.engine.Parameters.ProbabilityOfNonTerminalMutation,
				antagonist:                       g.Antagonists[i],
				protagonist:                      g.Protagonists[j],
				generation:                       g,
				program:                          g.engine.Parameters.StartIndividual,
				id:                               CreateEpochID(count, g.GenerationID, g.Antagonists[i].Id,
					g.Protagonists[j].Id),
			}
			count++
		}
	}
	return epochs, nil
}

// CurrentPopulation retrieves the current population of the given Generation.
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
	if len(epochs) < 1 {
		return nil, fmt.Errorf("epochs slice is empty")
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
// that returns the Result of the epoch.
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

// ApplyParentSelection takes in a given Generation and returns a set of individuals once the preselected parent
// selection Strategy has been applied to the Generation.
// These individuals are ready to be taken to either a new Generation or preferably through survivor selection in the
// case you do not isEqual the population to grow in size.
func (g *Generation) ApplyParentSelection() ([]*Individual, error) {
	if !g.isComplete {
		return nil, fmt.Errorf("Generation #Id: %s has not competed, ", g.GenerationID)
	}

	currentPopulation, err := g.CurrentPopulation()
	if err != nil {
		return nil, err
	}

	switch g.engine.Parameters.ParentSelection {
	case ParentSelectionTournament:
		selectedInvididuals, err := TournamentSelection(currentPopulation, g.engine.Parameters.TournamentSize)
		if err != nil {
			return nil, err
		}
		g.hasParentSelectionHappened = true
		return selectedInvididuals, nil
	case ParentSelectionElitism:
		selectedInvididuals, err := Elitism(currentPopulation, g.engine.Parameters.ElitismPercentage)
		if err != nil {
			return nil, err
		}
		g.hasParentSelectionHappened = true
		return selectedInvididuals, nil
	default:
		return nil, fmt.Errorf("no appropriate parent selection Strategy selected. See parentselection." +
			"go file for information on integer values that represent different parent selection strategies")
	}
}

// ApplySurvivorSelection applies the preselected survivor selection Strategy.
// It DOES NOT check to see if the parent selection has already been applied,
// as in some cases evolutionary programs may choose to run without the parent selection phase.
// The onus is on the evolutionary architect to keep this consideration in mind.
func (g *Generation) ApplySurvivorSelection() ([]*Individual, error) {
	if !g.isComplete {
		return nil, fmt.Errorf("Generation #Id: %s has not competed, ", g.GenerationID)
	}

	return nil, nil
}

// GenerateRandomIndividual creates a a random set of individuals based on the parameters passed into the
// evolution engine. To pass a tree to an individual pass it via the formal parameters and not through the evolution
// engine
// parameter section
// Antagonists are by default
// set with the StartIndividuals Program as their own
// program.
func (g *Generation) GenerateRandomIndividual(kind int, prog Program) ([]*Individual, error) {
	if g.engine.Parameters.EachPopulationSize < 1 {
		return nil, fmt.Errorf("number should at least be 1")
	}
	if kind == IndividualAntagonist {
		if g.engine.Parameters.AntagonistMaxStrategies < 1 {
			return nil, fmt.Errorf("antagonist maxNumberOfStrategies should at least be 1")
		}
		if len(g.engine.Parameters.AntagonistAvailableStrategies) < 1 {
			return nil, fmt.Errorf("antagonist availableStrategies should at least have one Strategy")
		}
	} else if kind == IndividualProtagonist {
		if g.engine.Parameters.ProtagonistMaxStrategies < 1 {
			return nil, fmt.Errorf("protagonist maxNumberOfStrategies should at least be 1")
		}
		if len(g.engine.Parameters.ProtagonistAvailableStrategies) < 1 {
			return nil, fmt.Errorf("protagonist availableStrategies should at least have one Strategy")
		}
	} else {
		return nil, fmt.Errorf("unknown individual kind")
	}

	individuals := make([]*Individual, g.engine.Parameters.EachPopulationSize)

	for i := 0; i < g.engine.Parameters.EachPopulationSize; i++ {

		var numberOfStrategies int
		var randomStrategies []Strategy

		if kind == IndividualAntagonist {
			// TODO fix equal Strategy length issue
			if g.engine.Parameters.SetEqualStrategyLength {
				numberOfStrategies = g.engine.Parameters.EqualStrategiesLength
				randomStrategies = GenerateRandomStrategy(g.engine.Parameters.EqualStrategiesLength,
					g.engine.Parameters.AntagonistAvailableStrategies)
			} else {
				numberOfStrategies = rand.Intn(g.engine.Parameters.AntagonistMaxStrategies)
				randomStrategies = GenerateRandomStrategy(numberOfStrategies, g.engine.Parameters.AntagonistAvailableStrategies)
			}
		} else if kind == IndividualProtagonist {
			// TODO fix equal Strategy length issue
			if g.engine.Parameters.SetEqualStrategyLength {
				numberOfStrategies = g.engine.Parameters.EqualStrategiesLength
				randomStrategies = GenerateRandomStrategy(g.engine.Parameters.EqualStrategiesLength,
					g.engine.Parameters.ProtagonistAvailableStrategies)
			} else {
				numberOfStrategies = rand.Intn(g.engine.Parameters.ProtagonistMaxStrategies)
				randomStrategies = GenerateRandomStrategy(numberOfStrategies, g.engine.Parameters.ProtagonistAvailableStrategies)
			}
		}

		id := fmt.Sprintf("%s-%d", KindToString(kind), i)
		var individual *Individual

		if prog.T == nil {
			individual = &Individual{
				Kind:     kind,
				Id:       id,
				Strategy: randomStrategies,
				Fitness:  make([]float64, 0),
				Program:  nil,
				BirthGen: 0,
			}
		} else {
			prog.ID = GenerateProgramID(i)

			clone, err := prog.Clone()
			if err != nil {
				return nil, err
			}
			individual = &Individual{
				Kind:     kind,
				Id:       id,
				Strategy: randomStrategies,
				Fitness:  make([]float64, 0),
				Program:  &clone,
			}
		}

		individuals[i] = individual
	}
	return individuals, nil
}

type GenerationResult struct {
	generation *Generation
}

// RunNext takes in a current GenerationResult runs a set of parent and survivor selection mechanisms,
// and returns the new Generation
func (g *GenerationResult) RunNext() *GenerationResult {
	return nil
}

type GenerationEngine struct {
}
