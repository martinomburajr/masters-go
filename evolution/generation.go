package evolution

import (
	"fmt"
	"math/rand"
	"time"
)

type Generation struct {
	GenerationID                 string
	Protagonists                 []Individual //Protagonists in a given generation
	Antagonists                  []Individual //Antagonists in a given generation
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
	completeEpochs, err := g.runEpochs(setupEpochs)
	if err != nil {
		return nil, err
	}

	// Set Individuals back to generation
	for e := range completeEpochs {
		for i := range g.Protagonists {
			if completeEpochs[e].protagonist.id == g.Protagonists[i].id {
				g.Protagonists[i].fitness = append(completeEpochs[e].protagonist.fitness)
				g.Protagonists[i].totalFitness = completeEpochs[e].protagonist.totalFitness
				g.Protagonists[i].hasCalculatedFitness = completeEpochs[e].protagonist.hasCalculatedFitness
			}
		}
	}




	for i := range g.Antagonists {
		for e := range completeEpochs {
			if completeEpochs[e].antagonist.id == g.Antagonists[i].id {
				g.Antagonists[i] = completeEpochs[e].antagonist
				break
			}
		}
	}

	for e := range g.Protagonists {
		protagonistFitness, err := AggregateFitness(g.Protagonists[e])
		if err != nil {
			return nil, err
		}
		g.Protagonists[e].totalFitness = protagonistFitness
		g.Protagonists[e].hasCalculatedFitness = true
	}

	for e := range g.Antagonists {
		antagonistFitness, err := AggregateFitness(g.Antagonists[e])
		if err != nil {
			return nil, err
		}
		g.Antagonists[e].totalFitness = antagonistFitness
		g.Antagonists[e].hasCalculatedFitness = true
	}


	err = g.CollectBruteStatistics()
	if err != nil {
		return nil, err
	}

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
		count:                        g.count + 1,
	}, nil
}

func GenerateGenerationID(count int) string {
	return fmt.Sprintf("GEN-%d", count)
}

func CalculateCumulativeFitnessOfAllIndividuals(individuals []Individual) (int, error) {
	fitness := 0
	for e := range individuals {
		indFit, err := CalculateCumulativeFitness(individuals[e])
		if err != nil {
			return -1, err
		}
		fitness += indFit
	}
	return fitness, nil
}

func CalculateCumulativeFitness(individual Individual) (int, error) {
	if !individual.hasCalculatedFitness {
		return -1, fmt.Errorf("CalculateCumulativeFitness | cannot get fitness as individual not applied fitness" +
			" %s", individual.id)
	}

	fitness := 0
	for e := range individual.fitness {
		fitness += individual.fitness[e]
	}
	return fitness, nil
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
	for i := range g.Antagonists {
		for j := range g.Protagonists {
			epochs[count] = Epoch{
				isComplete:                       false,
				protagonistBegins:                false,
				terminalSet:                      g.engine.Parameters.TerminalSet,
				nonTerminalSet:                   g.engine.Parameters.NonTerminalSet,
				hasAntagonistApplied:             false,
				hasProtagonistApplied:            false,
				probabilityOfMutation:            g.engine.Parameters.ProbabilityOfMutation,
				probabilityOfNonTerminalMutation: g.engine.Parameters.ProbabilityOfNonTerminalMutation,
				antagonist:                       &g.Antagonists[i],
				protagonist:                      &g.Protagonists[j],
				generation:                       g,
				program:                          g.engine.Parameters.StartIndividual,
				id:                               CreateEpochID(count, g.GenerationID, g.Antagonists[i].id, g.Protagonists[j].id),
			}
			count++
		}
	}
	return epochs, nil
}

// CurrentPopulation retrieves the current population of the given generation.
// Individuals may have competed or may have been altered in a variety of ways.
// This will return a list of references Individuals
func (g *Generation) CurrentPopulation() ([]Individual, error) {
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

type JudementDayStatistics struct {
	Top3Antagonists  []Individual
	Top3Protagonists []Individual
}

func (g Generation) CollectBruteStatistics() error {
	//if !g.isComplete {
	//	return fmt.Errorf("generation | Statistics | cannot collect statistics until generation is complete (" +
	//		"marked complete)")
	//}
	if g.Antagonists == nil {
		return fmt.Errorf("generation | Statistics | antagonists cannot be nil")
	}
	if len(g.Antagonists) < 1 {
		return fmt.Errorf("generation | Statistics | antagonists cannot be empty")
	}
	if g.Protagonists == nil {
		return fmt.Errorf("generation | Statistics | protagonists cannot be nil")
	}
	if len(g.Protagonists) < 1 {
		return fmt.Errorf("generation | Statistics | protagonists cannot be empty")
	}
	if g.engine == nil {
		return fmt.Errorf("generation | Statistics | evolution engine cannot be nil")
	}
	if g.hasParentSelectionHappened {
		return fmt.Errorf("generation | Statistics | cannot collect statistics after parent selection. " +
			"Brute statistics will be collected before any form of selection")
	}
	if g.hasSurvivorSelectionHappened {
		return fmt.Errorf("generation | Statistics | cannot collect statistics after survivor selection. " +
			"Brute statistics will be collected before any form of selection")
	}

	return nil
}

func (g Generation) CollectRefinedStatistics() error {
	if !g.isComplete {
		return fmt.Errorf("generation | Statistics | cannot collect statistics until generation is complete (" +
			"marked complete)")
	}
	if g.Antagonists == nil {
		return fmt.Errorf("generation | Statistics | antagonists cannot be nil")
	}
	if len(g.Antagonists) < 1 {
		return fmt.Errorf("generation | Statistics | antagonists cannot be empty")
	}
	if g.Protagonists == nil {
		return fmt.Errorf("generation | Statistics | protagonists cannot be nil")
	}
	if len(g.Protagonists) < 1 {
		return fmt.Errorf("generation | Statistics | protagonists cannot be empty")
	}
	if g.engine == nil {
		return fmt.Errorf("generation | Statistics | evolution engine cannot be nil")
	}
	if !g.hasParentSelectionHappened {
		return fmt.Errorf("generation | Statistics | Refined | cannot collect statistics after parent selection. " +
			"Brute statistics will be collected before any form of selection")
	}
	if !g.hasSurvivorSelectionHappened {
		return fmt.Errorf("generation | Statistics | Refined | cannot collect statistics after survivor selection. " +
			"Brute statistics will be collected before any form of selection")
	}

	return nil
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

// ApplyParentSelection - INCOMPLETE!!!! takes in a given generation and returns a set of individuals once the
// preselected parent
// selection strategy has been applied to the generation.
// These individuals are ready to be taken to either a new generation or preferably through survivor selection in the
// case you do not isEqual the population to grow in size.
func (g *Generation) ApplyParentSelection() ([]Individual, error) {
	//if !g.isComplete {
	//	return nil, fmt.Errorf("generation #id: %s has not competed, ", g.GenerationID)
	//}
	//
	//currentPopulation, err := g.CurrentPopulation()
	//if err != nil {
	//	return nil, err
	//}
	//
	//switch g.engine.Parameters.ParentSelection {
	//case ParentSelectionTournament:
	//	selectedInvididuals, err := TournamentSelection(currentPopulation, g.engine.Parameters.TournamentSize)
	//	if err != nil {
	//		return nil, err
	//	}
	//	g.hasParentSelectionHappened = true
	//	return selectedInvididuals, nil
	//case ParentSelectionElitism:
	//	selectedInvididuals, err := Elitism(currentPopulation, g.engine.Parameters.ElitismPercentage)
	//	if err != nil {
	//		return nil, err
	//	}
	//	g.hasParentSelectionHappened = true
	//	return selectedInvididuals, nil
	//default:
	//	return nil, fmt.Errorf("no appropriate parent selection strategy selected. See parentselection." +
	//		"go file for information on integer values that represent different parent selection strategies")
	//}
	return nil, nil
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

// GenerateRandomAntagonists creates a a random set of antagonists based on the parameters passed into the
// evolution engine. Antagonists are by default set with the StartIndividuals Program as their own program.
func (g *Generation) GenerateRandomAntagonists(idTemplate string) ([]Individual, error) {
	kind := IndividualAntagonist
	if g.engine.Parameters.EachPopulationSize < 1 {
		return nil, fmt.Errorf("number should at least be 1")
	}
	if g.engine.Parameters.AntagonistMaxStrategies < 1 {
		return nil, fmt.Errorf("maxNumberOfStrategies should at least be 1")
	}
	if len(g.engine.Parameters.AntagonistAvailableStrategies) < 1 {
		return nil, fmt.Errorf("availableStrategies should at least have one strategy")
	}
	if idTemplate == "" {
		return nil, fmt.Errorf("idTemplate cannot be empty")
	}

	individuals := make([]Individual, g.engine.Parameters.EachPopulationSize)

	for i := 0; i < g.engine.Parameters.EachPopulationSize; i++ {
		rand.Seed(time.Now().UnixNano())
		var numberOfStrategies int
		var randomStrategies []Strategy
		if g.engine.Parameters.SetEqualStrategyLength {
			numberOfStrategies = rand.Intn(g.engine.Parameters.EqualStrategiesLength)
			randomStrategies = GenerateRandomStrategy(numberOfStrategies, g.engine.Parameters.EqualStrategiesLength,
				g.engine.Parameters.AntagonistAvailableStrategies)
		} else {
			numberOfStrategies = rand.Intn(g.engine.Parameters.ProtagonistMaxStrategies)
			randomStrategies = GenerateRandomStrategy(numberOfStrategies, numberOfStrategies, g.engine.Parameters.AntagonistAvailableStrategies)
		}
		id := fmt.Sprintf("%s-%s-%d", KindToString(kind), "", i)

		program := g.engine.Parameters.StartIndividual
		programID := GenerateProgramID(i)
		program.ID = programID
		clone, err := program.Clone()

		if err != nil {
			return nil, err
		}

		individual := Individual{
			kind:     kind,
			id:       id,
			strategy: randomStrategies,
			fitness:  make([]int, 0),
			Program:  clone,
		}
		individuals[i] = individual
	}
	return individuals, nil
}

// GenerateRandomProtagonists creates a a random set of protagonists based on the parameters passed into the
// evolution engine.
func (g *Generation) GenerateRandomProtagonists(idTemplate string) ([]Individual, error) {
	kind := IndividualProtagonist
	if g.engine.Parameters.EachPopulationSize < 1 {
		return nil, fmt.Errorf("number should at least be 1")
	}
	if g.engine.Parameters.ProtagonistMaxStrategies < 1 {
		return nil, fmt.Errorf("maxNumberOfStrategies should at least be 1")
	}
	if len(g.engine.Parameters.ProtagonistAvailableStrategies) < 1 {
		return nil, fmt.Errorf("availableStrategies should at least have one strategy")
	}
	if idTemplate == "" {
		return nil, fmt.Errorf("idTemplate cannot be empty")
	}

	individuals := make([]Individual, g.engine.Parameters.EachPopulationSize)

	for i := 0; i < g.engine.Parameters.EachPopulationSize; i++ {
		rand.Seed(time.Now().UnixNano())
		var numberOfStrategies int
		var randomStrategies []Strategy
		if g.engine.Parameters.SetEqualStrategyLength {
			numberOfStrategies = rand.Intn(g.engine.Parameters.EqualStrategiesLength)
			randomStrategies = GenerateRandomStrategy(numberOfStrategies, g.engine.Parameters.EqualStrategiesLength, g.engine.Parameters.ProtagonistAvailableStrategies)
		} else {
			numberOfStrategies = rand.Intn(g.engine.Parameters.ProtagonistMaxStrategies)
			randomStrategies = GenerateRandomStrategy(numberOfStrategies, numberOfStrategies, g.engine.Parameters.ProtagonistAvailableStrategies)
		}
		id := fmt.Sprintf("%s-%s-%d", KindToString(kind), "", i)

		program := Program{}
		programID := GenerateProgramID(i)
		program.ID = programID

		individual := Individual{
			kind:     kind,
			id:       id,
			strategy: randomStrategies,
			fitness:  make([]int, 0),
			Program:  program,
		}
		individuals[i] = individual
	}
	return individuals, nil
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
