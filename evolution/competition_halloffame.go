package evolution

import (
	"github.com/martinomburajr/masters-go/utils"
	"math/rand"
	"time"
)

type HallOfFame struct {
	Engine *EvolutionEngine

	AntagonistArchive  []Individual
	ProtagonistArchive []Individual

	ReturnRounds int
}

func (s HallOfFame) Topology(currentGeneration *Generation,
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

func (s *HallOfFame) Evolve(params EvolutionParams, topology ITopology) (*EvolutionResult,
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

		// REINSERT HALL OF FAME
		if i%s.ReturnRounds == 0 && i != 0 {
			//Reinsert
			for j := 0; j < s.ReturnRounds; j++ {
				perm := rand.Perm(s.ReturnRounds)
				engine.Generations[i].Antagonists[perm[j]] = &s.AntagonistArchive[j]

				permProtagonist := rand.Perm(s.ReturnRounds)
				engine.Generations[i].Protagonists[permProtagonist[j]] = &s.ProtagonistArchive[j]
			}
		}

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
