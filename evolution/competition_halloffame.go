package evolution

import (
	"github.com/martinomburajr/masters-go/utils"
	"math"
	"math/rand"
	"time"
)

type HallOfFame struct {
	Engine *EvolutionEngine

	AntagonistArchive  []Individual
	ProtagonistArchive []Individual

	GenerationIntervals int
}

func (s *HallOfFame) Topology(currentGeneration *Generation,
	params EvolutionParams) (*Generation,
	error) {
	roundRobin := RoundRobin{Engine: s.Engine}
	nextGeneration, err := roundRobin.Topology(currentGeneration, params)
	if err != nil {
		return nil, err
	}

	bestAntagonist := &Individual{AverageFitness: math.MinInt16}
	bestProtagonist := &Individual{AverageFitness: math.MinInt16}
	for i := 0; i < len(currentGeneration.Antagonists); i++ {
		if currentGeneration.Antagonists[i].AverageFitness >= bestAntagonist.AverageFitness {
			bestAntagonist = currentGeneration.Antagonists[i]
		}
		if currentGeneration.Protagonists[i].AverageFitness >= bestProtagonist.AverageFitness {
			bestProtagonist = currentGeneration.Protagonists[i]
		}
	}

	currentGeneration.BestAntagonist, err = bestAntagonist.Clone()
	currentGeneration.BestProtagonist, err = bestProtagonist.Clone()

	s.AntagonistArchive = append(s.AntagonistArchive, currentGeneration.BestAntagonist)
	s.ProtagonistArchive = append(s.ProtagonistArchive, currentGeneration.BestProtagonist)

	return nextGeneration, nil
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

	s.GenerationIntervals = int(engine.Parameters.Topology.HoFGenerationInterval * float64(genCount))
	if s.GenerationIntervals >= int(float64(params.EachPopulationSize)*0.1) {
		for s.GenerationIntervals >= int(float64(params.EachPopulationSize)*0.1) {
			if s.GenerationIntervals < MinAllowableGenerationsToTerminate {
				s.GenerationIntervals = params.EachPopulationSize / 2
				break
			}
			s.GenerationIntervals /= 2
			if s.GenerationIntervals == 0 {
				s.GenerationIntervals = 4
			}
		}
	}

	for i := 0; i < genCount; i++ {
		started := time.Now()
		// 1. CLEANSE
		engine.Generations[i].CleansePopulations(engine.Parameters)

		// REINSERT HALL OF FAME

		if i%s.GenerationIntervals == 0 && i != 0 {
			//Reinsert
			perm := rand.Perm(s.GenerationIntervals)
			permProtagonist := rand.Perm(s.GenerationIntervals)
			for j := 0; j < s.GenerationIntervals; j++ {
				antagonistClone, err := s.AntagonistArchive[perm[j]].CloneCleanse()
				if err != nil {
					return nil, err
				}

				protagonistClone, err := s.ProtagonistArchive[perm[j]].CloneCleanse()
				if err != nil {
					return nil, err
				}

				engine.Generations[i].Antagonists[perm[j]] = &antagonistClone
				engine.Generations[i].Protagonists[permProtagonist[j]] = &protagonistClone
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
