package evolution

import (
	"fmt"
	"github.com/martinomburajr/masters-go/utils"
	"time"
)

type SingleEliminationTournamentTopology struct {
	Engine *EvolutionEngine
}

func (s SingleEliminationTournamentTopology) Topology(currentGeneration *Generation,
	params EvolutionParams) (*Generation,
	error) {
	topAntagonist, err := singleETCompeteAntagonists(currentGeneration.Antagonists, params)
	if err != nil {
		return nil, err
	}

	topProtagonist, err := singleETCompeteProtagonists(currentGeneration.Protagonists, *topAntagonist.Program.T, params)
	if err != nil {
		return nil, err
	}

	antagonistSurvivors, protagonistSurvivors := currentGeneration.ApplySelection(currentGeneration.Antagonists, currentGeneration.Protagonists, params.ErrorChan)
	currentGeneration.BestProtagonist, err = topProtagonist.Clone()
	currentGeneration.BestAntagonist, err = topAntagonist.Clone()

	newGeneration := &Generation{
		GenerationID:                 GenerateGenerationID(currentGeneration.count + 1),
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
			err = brackets[i].individualB.ApplyAntagonistStrategy(params)
			if err != nil {
				return nil, err
			}

			individualAFitness, _, err := brackets[i].individualA.CalculateAntagonistThresholdedFitness(params)
			if err != nil {
				return nil, err
			}
			individualBFitness, _, err := brackets[i].individualB.CalculateAntagonistThresholdedFitness(params)
			if err != nil {
				return nil, err
			}

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
		brackets, err = setCreateTournamentBrackets(winners)
		if err != nil {
			return nil, err
		}
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

			individualAFitness, _, err := brackets[i].individualA.CalculateProtagonistThresholdedFitness(params)
			if err != nil {
				return nil, err
			}
			individualBFitness, _, err := brackets[i].individualB.CalculateProtagonistThresholdedFitness(params)
			if err != nil {
				return nil, err
			}

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
		brackets, err = setCreateTournamentBrackets(winners)
		if err != nil {
			return nil, err
		}
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
	for i := 0; i < len(individuals); i += 2 {
		brackets[i].individualA = individuals[i]
		brackets[i].individualB = individuals[i+1]
	}
	return brackets, nil
}
