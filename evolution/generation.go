package evolution

import (
	"fmt"
	"strings"
	"sync"
)

// TODO AGE
// TODO Calculate fitness average for GENERATIONS (seems off!)
type Generation struct {
	GenerationID string
	Protagonists []*Individual //Protagonists in a given Generation
	Antagonists  []*Individual //Antagonists in a given Generation

	engine                       *EvolutionEngine
	isComplete                   bool
	hasParentSelectionHappened   bool
	hasSurvivorSelectionHappened bool
	count                        int

	//Individuals
	BestAntagonist  Individual
	BestProtagonist Individual

	// Averages of all Antagonists and Protagonists in Generation
	Correlation float64
	Covariance  float64

	AntagonistAverage    float64
	AntagonistStdDev     float64
	AntagonistVariance   float64
	AntagonistAvgFitness []float64
	AntagonistSkew       float64
	AntagonistExKurtosis float64

	ProtagonistAverage    float64
	ProtagonistStdDev     float64
	ProtagonistVariance   float64
	ProtagonistSkew       float64
	ProtagonistExKurtosis float64
	ProtagonistAvgFitness []float64
}

func (g *Generation) ToString() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("\n%s\n", g.GenerationID))
	sb.WriteString(fmt.Sprintf("CorrelationInGeneration: : %.2f\n", g.Correlation))
	sb.WriteString(fmt.Sprintf("CovarianceInGeneration: : %.2f\n", g.Covariance))
	sb.WriteString("_______________________________\n")
	sb.WriteString(fmt.Sprintf("AntagonistStdDevInGeneration : %.2f\n", g.AntagonistStdDev))
	sb.WriteString(fmt.Sprintf("AntagonistAvg : %.2f\n", g.AntagonistAverage))
	sb.WriteString(fmt.Sprintf("AntagonistStdDevInGeneration : %.2f\n", g.AntagonistStdDev))
	sb.WriteString(fmt.Sprintf("AntagonistVarianceInGeneration : %.2f\n", g.AntagonistVariance))
	sb.WriteString(fmt.Sprintf("AntagonistSkewInGeneration : %.2f\n", g.AntagonistSkew))
	sb.WriteString(fmt.Sprintf("AntagonistExKurtosisInGeneration : %.2f\n", g.AntagonistExKurtosis))
	sb.WriteString("<===================================>\n")
	sb.WriteString(fmt.Sprintf("ProtagonistAverageInGeneration : %.2f\n", g.ProtagonistAverage))
	sb.WriteString(fmt.Sprintf("ProtagonistStdDevInGeneration : %.2f\n", g.ProtagonistStdDev))
	sb.WriteString(fmt.Sprintf("ProtagonistVarianceInGeneration : %.2f\n", g.ProtagonistVariance))
	sb.WriteString(fmt.Sprintf("ProtagonistSkewInGeneration : %.2f\n", g.ProtagonistSkew))
	sb.WriteString(fmt.Sprintf("ProtagonistExKurtosisInGeneration : %.2f\n\n\n", g.ProtagonistExKurtosis))

	return sb.String()
}

// initializePopulation randomly creates a set of antagonists and protagonists
func (g *Generation) InitializePopulation(params EvolutionParams) (antagonists []*Individual,
	protagonists []*Individual, err error) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(wg *sync.WaitGroup, params *EvolutionParams) {
		defer wg.Done()
		for i := 0; i < params.EachPopulationSize; i++ {
			g.Antagonists, err = g.GenerateRandomIndividuals(IndividualAntagonist, *params)
			if err != nil {
				params.ErrorChan <- err
			}
		}
	}(&wg, &params)

	go func(wg *sync.WaitGroup, params *EvolutionParams) {
		defer wg.Done()
		for i := 0; i < params.EachPopulationSize; i++ {
			g.Protagonists, err = g.GenerateRandomIndividuals(IndividualProtagonist, *params)
			if err != nil {
				params.ErrorChan <- err
			}
		}
	}(&wg, &params)
	wg.Wait()

	return g.Antagonists, g.Protagonists, nil
}

// ApplySelection applies all 3 selection methods, parent,
// reproduction and survivor to return a set of survivor antagonist and protagonists
func (g *Generation) ApplySelection(antagonists, protagonists []*Individual, errorChan chan error) (
	antagonistSurvivors []*Individual, protagonistSurvivors []*Individual) {
	antSurvivorChan := make(chan []*Individual)
	go func(g *Generation) {
		antWinnerParents, err := g.ApplyParentSelection(antagonists)
		if err != nil {
			errorChan <- err
		}
		antSelectedParents, antSelectedChildren, err := g.ApplyReproduction(antWinnerParents, IndividualAntagonist)
		if err != nil {
			errorChan <- err
		}
		antSurvivors, err := g.ApplySurvivorSelection(antSelectedParents, antSelectedChildren)
		if err != nil {
			errorChan <- err
		}
		antSurvivorChan <- antSurvivors
		close(antSurvivorChan)
	}(g)
	proSurvivorChan := make(chan []*Individual)
	go func(g *Generation) {
		proWinnerParents, err := g.ApplyParentSelection(protagonists)
		if err != nil {
			errorChan <- err
		}
		proSelectedParents, proSelectedChildren, err := g.ApplyReproduction(proWinnerParents, IndividualProtagonist)
		if err != nil {
			errorChan <- err
		}
		proSurvivors, err := g.ApplySurvivorSelection(proSelectedParents, proSelectedChildren)
		if err != nil {
			errorChan <- err
		}
		proSurvivorChan <- proSurvivors
		close(proSurvivorChan)
	}(g)

	for {
		select {
		case x := <-antSurvivorChan:
			if x != nil {
				antagonistSurvivors = x
			}
			break
		case y := <-proSurvivorChan:
			if y != nil {
				protagonistSurvivors = y
			}
			break
		default:
			if protagonistSurvivors != nil && antagonistSurvivors != nil {
				break
			}
		}
		if protagonistSurvivors != nil && antagonistSurvivors != nil {
			break
		}
	}
	g.hasSurvivorSelectionHappened = true
	g.hasParentSelectionHappened = true
	return antagonistSurvivors, protagonistSurvivors
}

func GenerateGenerationID(count int) string {
	return fmt.Sprintf("GEN-%d", count)
}

func (g *Generation) CleansePopulations(params EvolutionParams) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(wg *sync.WaitGroup, errChan chan error) {
		defer wg.Done()
		protagonists, err := CleansePopulation(g.Protagonists, *params.StartIndividual.T)
		if err != nil {
			errChan <- err
		}
		g.Protagonists = protagonists
	}(&wg, params.ErrorChan)

	go func(wg *sync.WaitGroup, errChan chan error) {
		defer wg.Done()
		protagonists, err := CleansePopulation(g.Antagonists, *params.StartIndividual.T)
		if err != nil {
			errChan <- err
		}
		g.Protagonists = protagonists
	}(&wg, params.ErrorChan)
	wg.Wait()
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
func (g *Generation) ApplyReproduction(incomingParents []*Individual, kind int) (outgoingParents []*Individual,
	children []*Individual,
	err error) {
	children = make([]*Individual, g.engine.Parameters.EachPopulationSize)

	switch g.engine.Parameters.Reproduction.CrossoverStrategy {
	case CrossoverSinglePoint:
		for i := 0; i < len(incomingParents); i += 2 {
			child1, child2, err := SinglePointCrossover(incomingParents[i], incomingParents[i+1])
			if err != nil {
				return nil, nil, err
			}
			child1.BirthGen = g.count
			child2.BirthGen = g.count
			child1.Age = 0
			child2.Age = 0
			children[i] = &child1
			children[i+1] = &child2
		}
	case CrossoverFixedPoint:
		for i := 0; i < len(incomingParents); i += 2 {
			child1, child2, err := FixedPointCrossover(*incomingParents[i], *incomingParents[i+1], g.engine.Parameters)
			if err != nil {
				return nil, nil, err
			}
			child1.BirthGen = g.count
			child2.BirthGen = g.count
			child1.Age = 0
			child2.Age = 0
			children[i] = &child1
			children[i+1] = &child2
		}
	case CrossoverKPoint:
		for i := 0; i < len(incomingParents); i += 2 {
			child1, child2, err := KPointCrossover(incomingParents[i], incomingParents[i+1],
				g.engine.Parameters.Reproduction.KPointCrossover)
			if err != nil {
				return nil, nil, err
			}
			child1.BirthGen = g.count
			child2.BirthGen = g.count
			child1.Age = 0
			child2.Age = 0
			children[i] = &child1
			children[i+1] = &child2
		}
	case CrossoverUniform:
		for i := 0; i < len(incomingParents); i += 2 {
			child1, child2, err := UniformCrossover(incomingParents[i], incomingParents[i+1])
			if err != nil {
				return nil, nil, err
			}
			child1.BirthGen = g.count
			child2.BirthGen = g.count
			child1.Age = 0
			child2.Age = 0
			children[i] = &child1
			children[i+1] = &child2
		}
	default:
		return nil, nil, fmt.Errorf("no appropriate FixedPointCrossover operation was selected")
	}

	return Mutate(incomingParents, children, kind, g.engine.Parameters)
}

// ApplySurvivorSelection applies the preselected survivor selection Strategy.
// It DOES NOT check to see if the parent selection has already been applied,
// as in some cases evolutionary programs may choose to run without the parent selection phase.
// The onus is on the evolutionary architect to keep this consideration in mind.
func (g *Generation) ApplySurvivorSelection(outgoingParents []*Individual,
	children []*Individual) ([]*Individual, error) {

	switch g.engine.Parameters.Selection.Survivor.Type {
	case SurvivorSelectionFitnessBased:
		return FitnessBasedSurvivorSelection(outgoingParents, children, g.engine.Parameters)
	case SurvivorSelectionRandom:
		return RandomSurvivorSelection(outgoingParents, children, g.engine.Parameters)
	default:
		return nil, fmt.Errorf("Invalid Survivor Selection Selected")
	}
}
