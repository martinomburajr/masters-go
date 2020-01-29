package evolution

import (
	"github.com/martinomburajr/masters-go/utils"
	"math/rand"
	"time"
)

type KRandom struct {
	Engine *EvolutionEngine
}

func (s KRandom) Topology(currentGeneration *Generation,
	params EvolutionParams) (*Generation,
	error) {

	tournamentLedger, err := s.createTournamentLedger(currentGeneration.Antagonists, currentGeneration.Protagonists, params)
	if err != nil {
		return nil, err
	}

	err = s.startTournaments(currentGeneration, tournamentLedger, params)
	if err != nil {
		return nil, err
	}

	antagonistSurvivors, protagonistSurvivors := currentGeneration.ApplySelection(currentGeneration.Antagonists, currentGeneration.Protagonists, params.ErrorChan)

	newGeneration := &Generation{
		GenerationID:                 GenerateGenerationID(currentGeneration.count+1, TopologyKRandom),
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

func (s *KRandom) createTournamentLedger(antagonists []*Individual, protagonists []*Individual,
	params EvolutionParams) (tournamentLedger map[*Individual][]*Individual, err error) {
	tournamentLedger = make(map[*Individual][]*Individual, params.EachPopulationSize)
	opponents := make([][]*Individual, params.EachPopulationSize)

	for i := 0; i < params.EachPopulationSize; i++ {
		opponents[i] = make([]*Individual, params.Topology.KRandomK)
		for j := 0; j < params.Topology.KRandomK; j++ {
			randIndex := rand.Intn(params.EachPopulationSize)
			clone, err := antagonists[randIndex].Clone()
			if err != nil {
				return nil, err
			}
			clone.Parent = antagonists[randIndex]
			opponents[i][j] = &clone
		}
	}

	for i := 0; i < params.EachPopulationSize; i++ {
		tournamentLedger[protagonists[i]] = opponents[i]
	}

	return tournamentLedger, nil
}

func (s *KRandom) contains(individual *Individual, individuals []*Individual) bool {
	for i := 0; i < len(individuals); i++ {
		if individual.Id == individuals[i].Id {
			return true
		}
	}
	return false
}

func (s *KRandom) startTournaments(currentGeneration *Generation, tournamentLedger map[*Individual][]*Individual,
	params EvolutionParams) (err error) {
	perfectFitnessMap := map[string]PerfectTree{}

	antagonists := make([]*Individual, params.Topology.KRandomK)
	protagonists := make([]*Individual, 0)

	for protagonist := range tournamentLedger {
		tournament := tournamentLedger[protagonist]
		for j, antagonist := range tournament {
			antagonists[j] = antagonist
			err := antagonist.ApplyAntagonistStrategy(params)
			if err != nil {
				return err
			}

			err = protagonist.ApplyProtagonistStrategy(*antagonist.Program.T, params)
			if err != nil {
				return err
			}

			antagonistFitness, protagonistFitness, antagonistFitnessDelta, protagonistFitnessDelta,
			err := ThresholdedRatioFitness(params.Spec, antagonist.Program, protagonist.Program,
				params.SpecParam.DivideByZeroStrategy)
			if err != nil {
				return err
			}

			//antagonist.Fitness = append(antagonist.Fitness, antagonistFitness)
			AntagonistFitnessResolver(perfectFitnessMap, antagonist, antagonistFitness, antagonistFitnessDelta)
			ProtagonistFitnessResolver(perfectFitnessMap, protagonist, protagonistFitness, protagonistFitnessDelta)
		}
		protagonists = append(protagonists, protagonist)
	}

	// Set individuals with the best representation of their tree
	for i := 0; i < len(antagonists); i++ {
		perfectAntagonistTree := perfectFitnessMap[antagonists[i].Id]
		currentGeneration.Antagonists[i] = antagonists[i].Parent
		currentGeneration.Antagonists[i].Program = perfectAntagonistTree.Program
		currentGeneration.Antagonists[i].BestDelta = perfectAntagonistTree.BestFitnessDelta
		currentGeneration.Antagonists[i].BestFitness = perfectAntagonistTree.BestFitnessValue
	}
	for i := 0; i < len(protagonists); i++ {
		perfectProtagonistTree := perfectFitnessMap[protagonists[i].Id]
		currentGeneration.Protagonists[i] = protagonists[i]
		currentGeneration.Protagonists[i].Program = perfectProtagonistTree.Program
		currentGeneration.Protagonists[i].BestDelta = perfectProtagonistTree.BestFitnessDelta
		currentGeneration.Protagonists[i].BestFitness = perfectProtagonistTree.BestFitnessValue
	}

	return nil
}

func (s *KRandom) Evolve(params EvolutionParams, topology ITopology) (*EvolutionResult,
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
