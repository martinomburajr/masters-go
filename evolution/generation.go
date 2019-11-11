package evolution

import (
	"fmt"
	"gonum.org/v1/gonum/stat"
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
	err := g.Compete()

	// Parent Selection
	parentSelectionAntagonist, err := g.ApplyParentSelection(g.Antagonists)
	if err != nil {
		return nil, err
	}
	parentSelectionProtagonist, err := g.ApplyParentSelection(g.Protagonists)
	if err != nil {
		return nil, err
	}

	nextGenAntagonists, err := JudgementDay(parentSelectionAntagonist, IndividualAntagonist, generationCount, g.engine.Parameters)
	if err != nil {
		return nil, err
	}

	nextGenProtagonists, err := JudgementDay(parentSelectionProtagonist, IndividualProtagonist, generationCount,
		g.engine.Parameters)
	if err != nil {
		return nil, err
	}

	g.hasSurvivorSelectionHappened = true
	g.hasParentSelectionHappened = true

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
	for i := range g.Antagonists {
		for j := range g.Protagonists {

			cloneAntagonist, err := g.Antagonists[i].Clone()
			if err != nil {
				return nil, err
			}
			cloneAntagonist.Parent = g.Antagonists[i]

			cloneProtagonist, err := g.Protagonists[i].Clone()
			if err != nil {
				return nil, err
			}
			cloneProtagonist.Parent = g.Protagonists[i]

			epochs[count] = Epoch{
				isComplete:            false,
				terminalSet:           g.engine.Parameters.SpecParam.AvailableSymbolicExpressions.Terminals,
				nonTerminalSet:        g.engine.Parameters.SpecParam.AvailableSymbolicExpressions.NonTerminals,
				hasAntagonistApplied:  false,
				hasProtagonistApplied: false,
				antagonist:            &cloneAntagonist,
				protagonist:           &cloneProtagonist,
				generation:            g,
				program:               g.engine.Parameters.StartIndividual,
				id: CreateEpochID(count, g.GenerationID, g.Antagonists[i].Id,
					g.Protagonists[j].Id),
			}
			count++
		}
	}
	return epochs, nil
}

type PerfectTree struct {
	Program      *Program
	FitnessValue float64
	FitnessDelta float64
}

// runEpoch begins the run of a single epoch
func (g *Generation) runEpochs(epochs []Epoch) ([]Epoch, error) {
	if epochs == nil {
		return nil, fmt.Errorf("epochs have not been initialized | epochs is nil")
	}
	if len(epochs) < 1 {
		return nil, fmt.Errorf("epochs slice is empty")
	}

	perfectFitnessMap := map[string]PerfectTree{}
	for i := 0; i < len(epochs); i++ {
		err := epochs[i].Start(perfectFitnessMap)
		if err != nil {
			return nil, err
		}
	}

	// Set individuals with the best representation of their tree
	for i := 0; i < len(g.Antagonists); i++ {
		perfectAntagonistTree := perfectFitnessMap[g.Antagonists[i].Id]
		g.Antagonists[i].Program = perfectAntagonistTree.Program
		g.Antagonists[i].BestFitnessDelta = perfectAntagonistTree.FitnessDelta
		g.Antagonists[i].BestFitness = perfectAntagonistTree.FitnessValue
	}
	for i := 0; i < len(g.Protagonists); i++ {
		perfectProtagonistTree := perfectFitnessMap[g.Protagonists[i].Id]
		g.Protagonists[i].Program = perfectProtagonistTree.Program
		g.Protagonists[i].BestFitnessDelta = perfectProtagonistTree.FitnessDelta
		g.Antagonists[i].BestFitness = perfectProtagonistTree.FitnessValue
	}

	return epochs, nil
}

// Compete gives protagonist and anatagonists the chance to compete. A competition involves an epoch,
// that returns the Individuals of the epoch.
func (g *Generation) Compete() error {
	setupEpochs, err := g.setupEpochs()
	if err != nil {
		return err
	}

	// Runs the epochs and returns completed epochs that contain Fitness information within each individual.
	_, err = g.runEpochs(setupEpochs)
	if err != nil {
		return err
	}

	// Calculate the Fitness for individuals in the Generation
	for i := 0; i < len(g.Protagonists); i++ {
		mean, std := stat.MeanStdDev(g.Protagonists[i].Fitness, nil)
		variance := stat.Variance(g.Protagonists[i].Fitness, nil)
		g.Protagonists[i].AverageFitness = mean
		g.Protagonists[i].FitnessStdDev = std
		g.Protagonists[i].FitnessVariance = variance
		g.Protagonists[i].HasCalculatedFitness = true
		g.Protagonists[i].HasAppliedStrategy = true
		g.Protagonists[i].Age++

		antMean, antStd := stat.MeanStdDev(g.Antagonists[i].Fitness, nil)
		antVariance := stat.Variance(g.Antagonists[i].Fitness, nil)
		g.Antagonists[i].AverageFitness = antMean
		g.Antagonists[i].FitnessStdDev = antStd
		g.Antagonists[i].FitnessVariance = antVariance
		g.Antagonists[i].HasCalculatedFitness = true
		g.Antagonists[i].HasAppliedStrategy = true
		g.Antagonists[i].Age++
	}

	return err
}

// ApplyParentSelection takes in a given Generation and returns a set of individuals once the preselected parent
// selection Strategy has been applied to the Generation.
// These individuals are ready to be taken to either a new Generation or preferably through survivor selection in the
// case you do not isEqual the population to grow in size.
func (g *Generation) ApplyParentSelection(currentPopulation []*Individual) ([]*Individual, error) {
	switch g.engine.Parameters.Selection.Parent.Type {
	case ParentSelectionTournament:
		selectedInvididuals, err := TournamentSelection(currentPopulation, g.engine.Parameters.Selection.Parent.TournamentSize)
		if err != nil {
			return nil, err
		}
		g.hasParentSelectionHappened = true
		return selectedInvididuals, nil
	case ParentSelectionElitism:
		selectedInvididuals, err := Elitism(currentPopulation, true)
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
		if g.engine.Parameters.Strategies.AntagonistStrategyCount < 1 {
			return nil, fmt.Errorf("antagonist maxNumberOfStrategies should at least be 1")
		}
		if len(g.engine.Parameters.Strategies.AntagonistAvailableStrategies) < 1 {
			return nil, fmt.Errorf("antagonist availableStrategies should at least have one Strategy")
		}
	} else if kind == IndividualProtagonist {
		if g.engine.Parameters.Strategies.ProtagonistStrategyCount < 1 {
			return nil, fmt.Errorf("protagonist maxNumberOfStrategies should at least be 1")
		}
		if len(g.engine.Parameters.Strategies.ProtagonistAvailableStrategies) < 1 {
			return nil, fmt.Errorf("protagonist availableStrategies should at least have one Strategy")
		}
	} else {
		return nil, fmt.Errorf("unknown individual kind")
	}

	individuals := make([]*Individual, g.engine.Parameters.EachPopulationSize)

	for i := 0; i < g.engine.Parameters.EachPopulationSize; i++ {

		var randomStrategies []Strategy

		if kind == IndividualAntagonist {
			randomStrategies = GenerateRandomStrategy(g.engine.Parameters.Strategies.AntagonistStrategyCount,
				g.engine.Parameters.Strategies.AntagonistAvailableStrategies)
		} else if kind == IndividualProtagonist {
			randomStrategies = GenerateRandomStrategy(g.engine.Parameters.Strategies.ProtagonistStrategyCount,
				g.engine.Parameters.Strategies.ProtagonistAvailableStrategies)
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
