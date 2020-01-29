package evolution

import (
	"fmt"
	"github.com/martinomburajr/masters-go/utils"
	"math/rand"
	"sync"
	"time"
)

type SingleEliminationTournamentTopology struct {
	Engine *EvolutionEngine
}

func (s SingleEliminationTournamentTopology) Topology(currentGeneration *Generation,
	params EvolutionParams) (*Generation,
	error) {

	fittestAntagonists := make([]*Individual, 0)
	fittestProtagonists := make([]*Individual, 0)

	wgAntagonist := sync.WaitGroup{}

	setNoOfTournaments := int(params.Topology.SETNoOfTournaments * float64(params.EachPopulationSize))
	if params.Topology.SETNoOfTournaments == 0 {
		setNoOfTournaments = int(0.1 * float64(params.EachPopulationSize))
	}
	if params.Topology.SETNoOfTournaments > 1 {
		setNoOfTournaments = 1
	}

	for i := 0; i < setNoOfTournaments ; i++ {
		wgAntagonist.Add(1)
		go func(wgAntagonist *sync.WaitGroup, individuals []*Individual) {
			defer wgAntagonist.Done()
			clonedIndividuals, err := CloneIndividualsLinkParent(individuals)
			if err != nil {
				params.ErrorChan <- err
			}
			topAntagonist, err := singleETCompeteAntagonists(clonedIndividuals, params)
			if err != nil {
				params.ErrorChan <- err
			}

			currentGeneration.Mutex.Lock()
				fittestAntagonists = append(fittestAntagonists, topAntagonist.Parent)
			currentGeneration.Mutex.Unlock()
		}(&wgAntagonist, currentGeneration.Antagonists)
	}
	wgAntagonist.Wait()

	if len(fittestAntagonists) != params.EachPopulationSize {
		diff := params.EachPopulationSize - len(fittestAntagonists)
		perm := rand.Perm(diff)
		for i := 0; i < diff; i++ {
			fittestAntagonists = append(fittestAntagonists, currentGeneration.Antagonists[perm[i]])
		}
	}

		wgProtagonist := sync.WaitGroup{}
	for i := 0; i < setNoOfTournaments; i++ {
		wgProtagonist.Add(1)
		go func(wgAntagonist *sync.WaitGroup, individuals []*Individual, antagonists []*Individual, i int) {
			defer wgProtagonist.Done()
			clonedIndividuals, err := CloneIndividualsLinkParent(individuals)
			if err != nil {
				params.ErrorChan <- err
			}
			topProtagonist, err := singleETCompeteProtagonists(clonedIndividuals, *antagonists[i].Program.T,
				params)
			if err != nil {
				params.ErrorChan <- err
			}
			currentGeneration.Mutex.Lock()
			fittestProtagonists = append(fittestProtagonists, topProtagonist.Parent)
			currentGeneration.Mutex.Unlock()
		}(&wgProtagonist, currentGeneration.Protagonists, fittestAntagonists, i)
	}
	wgProtagonist.Wait()

	if len(fittestProtagonists) != params.EachPopulationSize {
		diff := params.EachPopulationSize - len(fittestProtagonists)
		perm := rand.Perm(diff)
		for i := 0; i < diff; i++ {
			fittestProtagonists = append(fittestProtagonists, currentGeneration.Protagonists[perm[i]])
		}
	}

	for i := 0; i < len(currentGeneration.Protagonists); i++ {
		anttagAvgFitness := CoalesceFitnessStatistics(fittestAntagonists[i])
		protagAvgFitness := CoalesceFitnessStatistics(fittestProtagonists[i])

		currentGeneration.AntagonistAvgFitness = append(currentGeneration.AntagonistAvgFitness, anttagAvgFitness)
		currentGeneration.ProtagonistAvgFitness = append(currentGeneration.ProtagonistAvgFitness, protagAvgFitness)
	}

	currentGeneration.Antagonists = fittestAntagonists
	currentGeneration.Protagonists = fittestProtagonists

	antagonistSurvivors, protagonistSurvivors := currentGeneration.ApplySelection(currentGeneration.Antagonists, currentGeneration.Protagonists, params.ErrorChan)

	newGeneration := &Generation{
		GenerationID: GenerateGenerationID(currentGeneration.count+1,
			TopologySingleEliminationTournament),
		Protagonists:                 protagonistSurvivors,
		Antagonists:                  antagonistSurvivors,
		engine:                       currentGeneration.engine,
		isComplete:                   true,
		hasParentSelectionHappened:   true,
		hasSurvivorSelectionHappened: true,
		count:                        currentGeneration.count,
	}

	return newGeneration, nil
}

func CloneIndividualsLinkParent(individuals []*Individual) (outgoing []*Individual, err error) {
	outgoing = make([]*Individual, len(individuals))

	perm := rand.Perm(len(individuals))
	for i:= 0; i < len(individuals); i++ {
		individual, err := individuals[perm[i]].Clone()
		if err != nil {
			return nil, err
		}
		individual.Parent = individuals[perm[i]]
		outgoing[i] = &individual
	}
	return outgoing, nil
}

func (s *SingleEliminationTournamentTopology) Evolve(params EvolutionParams, topology ITopology) (*EvolutionResult,
	error) {
	engine := s.Engine
	err := engine.validate()
	if err != nil {
		return nil, err
	}

	_, _, err = engine.InitializeGenerations(engine.Parameters)
	if err != nil {
		return nil, err
	}

	genCount := CalculateGenerationSize(engine.Parameters)

	for i := 0; i < genCount; i++ {
		started := time.Now()
		// 1. CLEANSE
		engine.Generations[i].CleansePopulations(engine.Parameters)

		// 2. START
		nextGeneration, err := topology.Topology(engine.Generations[i], params)
		if err != nil {
			return nil, err
		}
		// 3. EVALUATE
		if genCount == params.GenerationsCount && params.MaxGenerations < MinAllowableGenerationsToTerminate {
			shouldTerminateEvolution := engine.EvaluateTerminationCriteria(engine.Generations[i], engine.Parameters)
			if shouldTerminateEvolution {
				break
			}
		}
		go engine.RunGenerationStatistics(engine.Generations[i])

		if i == engine.Parameters.MaxGenerations-1 {
			break
		}
		engine.Generations = append(engine.Generations, nextGeneration)
		engine.ProgressBar.Incr()

		// 4. LOG
		elapsed := utils.TimeTrack(started)
		go WriteGenerationToLog(engine, i, elapsed)
		go WriteToDataFolders(engine.Parameters.FolderPercentages, i, engine.Parameters.GenerationsCount, engine.Parameters)
	}

	evolutionResult := &EvolutionResult{}
	err = evolutionResult.Analyze(engine, engine.Generations, true,
		engine.Parameters)
	if err != nil {
		return nil, err
	}

	return evolutionResult, nil
}

type bracket struct {
	individualA *Individual
	individualB *Individual
}

func singleETCompeteAntagonists(individuals []*Individual, params EvolutionParams) (topIndividual *Individual,
	err error) {
	if len(individuals) < 1 {
		return nil, fmt.Errorf("singleETCompeteAntagonists | input individuals cannot be empty")
	}
	if len(individuals) == 0 {
		return nil, fmt.Errorf("singleETCompeteAntagonists | input individuals cannot be null")
	}

	brackets, err := setCreateTournamentBrackets(individuals)
	if err != nil {
		return nil, err
	}

	var winner *Individual
	for len(brackets) >= 1 {
		winners := make([]*Individual, 0)
		for i := range brackets {
			err := brackets[i].individualA.ApplyAntagonistStrategy(params)
			if err != nil {
				return nil, err
			}
			s, _ := brackets[i].individualA.Program.T.ToMathematicalString()
			fmt.Println(s)
			err = brackets[i].individualB.ApplyAntagonistStrategy(params)
			if err != nil {
				return nil, err
			}
			s2, _ := brackets[i].individualB.Program.T.ToMathematicalString()
			fmt.Println(s2)

			individualAFitness, individualADelta, err := brackets[i].individualA.CalculateAntagonistThresholdedFitness(
				params)
			if err != nil {
				return nil, err
			}

			individualBFitness, individualBDelta, err := brackets[i].individualB.
				CalculateAntagonistThresholdedFitness(params)
			if err != nil {
				return nil, err
			}
			brackets[i].individualA.Fitness = append(brackets[i].individualA.Fitness, individualAFitness)
			brackets[i].individualA.Deltas = append(brackets[i].individualA.Deltas, individualADelta)
			brackets[i].individualB.Fitness = append(brackets[i].individualB.Fitness, individualBFitness)
			brackets[i].individualB.Deltas = append(brackets[i].individualB.Deltas, individualBDelta)

			brackets[i].individualA.Parent.Fitness = append(brackets[i].individualA.Parent.Fitness, individualAFitness)
			brackets[i].individualA.Parent.Deltas = append(brackets[i].individualA.Parent.Deltas, individualADelta)
			brackets[i].individualB.Parent.Fitness = append(brackets[i].individualB.Parent.Fitness, individualBFitness)
			brackets[i].individualB.Parent.Deltas = append(brackets[i].individualB.Parent.Deltas, individualBDelta)


			if individualAFitness >= individualBFitness {
				if len(brackets) == 1 {
					winner = brackets[i].individualA
					break
				}
				winners = append(winners, brackets[i].individualA)
			} else {
				if len(brackets) == 1 {
					winner = brackets[i].individualB
					break
				}
				winners = append(winners, brackets[i].individualB)
			}
		}
		if len(brackets) > 1 {
			brackets, err = setCreateTournamentBrackets(winners)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	for i := range individuals {
		individuals[i] = nil
	}

	return winner, err
}

func singleETCompeteProtagonists(individuals []*Individual, bestAntagonistTree DualTree,
	params EvolutionParams) (
	topIndividual *Individual,
	err error) {

	if len(individuals) < 1 {
		return nil, fmt.Errorf("singleETCompeteProtagonists | input individuals cannot be empty")
	}
	if len(individuals) == 0 {
		return nil, fmt.Errorf("singleETCompeteProtagonists | input individuals cannot be null")
	}

	brackets, err := setCreateTournamentBrackets(individuals)
	if err != nil {
		return nil, err
	}

	var winner *Individual
	for len(brackets) >= 1 {
		winners := make([]*Individual, 0)
		for i := range brackets {
			err := brackets[i].individualA.ApplyProtagonistStrategy(bestAntagonistTree, params)
			if err != nil {
				return nil, err
			}
			err = brackets[i].individualB.ApplyProtagonistStrategy(bestAntagonistTree, params)
			if err != nil {
				return nil, err
			}

			individualAFitness, individualADelta, err := brackets[i].individualA.CalculateProtagonistThresholdedFitness(params)
			if err != nil {
				return nil, err
			}
			individualBFitness, individualBDelta, err := brackets[i].individualB.CalculateProtagonistThresholdedFitness(params)
			if err != nil {
				return nil, err
			}

			brackets[i].individualA.Fitness = append(brackets[i].individualA.Fitness, individualAFitness)
			brackets[i].individualA.Deltas = append(brackets[i].individualA.Deltas, individualADelta)
			brackets[i].individualB.Fitness = append(brackets[i].individualB.Fitness, individualBFitness)
			brackets[i].individualB.Deltas = append(brackets[i].individualB.Deltas, individualBDelta)

			brackets[i].individualA.Parent.Fitness = append(brackets[i].individualA.Parent.Fitness, individualAFitness)
			brackets[i].individualA.Parent.Deltas = append(brackets[i].individualA.Parent.Deltas, individualADelta)
			brackets[i].individualB.Parent.Fitness = append(brackets[i].individualB.Parent.Fitness, individualBFitness)
			brackets[i].individualB.Parent.Deltas = append(brackets[i].individualB.Parent.Deltas, individualBDelta)

			if individualAFitness >= individualBFitness {
				if len(brackets) == 1 {
					winner = brackets[i].individualA
					break
				}
				winners = append(winners, brackets[i].individualA)
			} else {
				if len(brackets) == 1 {
					winner = brackets[i].individualB
					break
				}
				winners = append(winners, brackets[i].individualB)
			}
		}
		if len(brackets) > 1 {
			brackets, err = setCreateTournamentBrackets(winners)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	for i := range individuals {
		individuals[i] = nil
	}

	return winner, err
}

// setCreateTournamentBrackets create the tournament bracket. individuals should be of one kind.
func setCreateTournamentBrackets(individuals []*Individual) ([]bracket, error) {
	if len(individuals) < 1 {
		return nil, fmt.Errorf("setCreateTournamentBrackets | input individuals cannot be empty")
	}
	if len(individuals) == 0 {
		return nil, fmt.Errorf("setCreateTournamentBrackets | input individuals cannot be null")
	}
	brackets := make([]bracket, len(individuals)/2)
	counter := 0
	for i := 0; i < len(individuals); i += 2 {
		brackets[counter].individualA = individuals[i]
		brackets[counter].individualB = individuals[i+1]
		counter++
	}
	return brackets, nil
}
