package evolution

import (
	"fmt"
	"github.com/martinomburajr/masters-go/evolog"
	"github.com/martinomburajr/masters-go/utils"
	"gonum.org/v1/gonum/stat"
	"time"
)

type RoundRobin struct {
	Engine *EvolutionEngine
}

func (r RoundRobin) Topology(currentGeneration *Generation,
	params EvolutionParams) (*Generation,
	error) {

	err := r.Compete(currentGeneration)
	if err != nil {
		return nil, err
	}

	antagonistSurvivors, protagonistSurvivors := currentGeneration.ApplySelection(currentGeneration.Antagonists, currentGeneration.Protagonists, params.ErrorChan)

	newGeneration := &Generation{
		GenerationID:                 GenerateGenerationID(currentGeneration.count+1, TopologyRoundRobin),
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

func (s *RoundRobin) Evolve(params EvolutionParams, topology ITopology) (*EvolutionResult,
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

// setupEpochs takes in the Generation individuals (
// protagonists and antagonists) and creates a set of uninitialized epochs
func (r *RoundRobin) setupEpochs(g *Generation) ([]Epoch, error) {
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
	Program          *Program
	BestFitnessValue float64
	BestFitnessDelta float64
}

// runEpoch begins the run of a single epoch
func (r *RoundRobin) runEpochs(g *Generation, epochs []Epoch) ([]Epoch, error) {
	if epochs == nil {
		return nil, fmt.Errorf("epochs have not been initialized | epochs is nil")
	}
	if len(epochs) < 1 {
		return nil, fmt.Errorf("epochs slice is empty")
	}

	perfectFitnessMap := map[string]PerfectTree{}
	for i := 0; i < len(epochs); i++ {
		err := epochs[i].Start(perfectFitnessMap, r.Engine.Parameters)
		if err != nil {
			g.engine.Parameters.ErrorChan <- err
			return nil, err
		}

		if len(epochs) > 10 {
			if i%(len(epochs)/10) == 0 {
				if g.engine.Parameters.EnableLogging {
					msg := fmt.Sprintf("\n  ==> Run: %d | Epoch: (%d/%d)",
						g.engine.Parameters.InternalCount,
						i+1,
						len(epochs))
					g.engine.Parameters.LoggingChan <- evolog.Logger{Type: evolog.LoggerEpoch, Message: msg,
						Timestamp: time.Now()}
				}
			}
		}
	}

	// Set individuals with the best representation of their tree
	for i := 0; i < len(g.Antagonists); i++ {
		perfectAntagonistTree := perfectFitnessMap[g.Antagonists[i].Id]
		g.Antagonists[i].Program = perfectAntagonistTree.Program
		g.Antagonists[i].BestDelta = perfectAntagonistTree.BestFitnessDelta
		g.Antagonists[i].BestFitness = perfectAntagonistTree.BestFitnessValue
	}
	for i := 0; i < len(g.Protagonists); i++ {
		perfectProtagonistTree := perfectFitnessMap[g.Protagonists[i].Id]
		g.Protagonists[i].Program = perfectProtagonistTree.Program
		g.Protagonists[i].BestDelta = perfectProtagonistTree.BestFitnessDelta
		g.Protagonists[i].BestFitness = perfectProtagonistTree.BestFitnessValue
	}

	return epochs, nil
}

// Compete gives protagonist and anatagonists the chance to compete. A competition involves an epoch,
// that returns the Individuals of the epoch.
func (r *RoundRobin) Compete(g *Generation) error {
	setupEpochs, err := r.setupEpochs(g)
	if err != nil {
		return err
	}

	// Runs the epochs and returns completed epochs that contain Fitness information within each individual.
	_, err = r.runEpochs(g, setupEpochs)
	if err != nil {
		return err
	}

	// TODO Ensure Children of Antagonists are being created, i.e different IDs during crossover
	// TODO use penalization when SPEc is 0

	// Calculate the Fitness for individuals in the Generation
	for i := 0; i < len(g.Protagonists); i++ {
		deltaAntMean := stat.Mean(g.Antagonists[i].Deltas, nil)
		antMean, antStd := stat.MeanStdDev(g.Antagonists[i].Fitness, nil)
		antVariance := stat.Variance(g.Antagonists[i].Fitness, nil)
		g.Antagonists[i].AverageFitness = antMean
		g.Antagonists[i].FitnessStdDev = antStd
		g.Antagonists[i].FitnessVariance = antVariance
		g.Antagonists[i].HasCalculatedFitness = true
		g.Antagonists[i].HasAppliedStrategy = true
		g.Antagonists[i].Age += 1
		g.Antagonists[i].AverageDelta = deltaAntMean
		g.AntagonistAvgFitness = append(g.AntagonistAvgFitness, antMean)

		deltaMean := stat.Mean(g.Protagonists[i].Deltas, nil)
		mean, std := stat.MeanStdDev(g.Protagonists[i].Fitness, nil)
		variance := stat.Variance(g.Protagonists[i].Fitness, nil)
		g.Protagonists[i].AverageFitness = mean
		g.Protagonists[i].FitnessStdDev = std
		g.Protagonists[i].FitnessVariance = variance
		g.Protagonists[i].HasCalculatedFitness = true
		g.Protagonists[i].HasAppliedStrategy = true
		g.Protagonists[i].Age += 1
		g.Protagonists[i].AverageDelta = deltaMean
		g.ProtagonistAvgFitness = append(g.ProtagonistAvgFitness, mean)
	}

	return err
}


func CoalesceFitnessStatistics(individual *Individual) (fitnessToBeAppendedToGenerationAvgFitness float64) {
	deltaMean := stat.Mean(individual.Deltas, nil)
	mean, std := stat.MeanStdDev(individual.Fitness, nil)
	variance := stat.Variance(individual.Fitness, nil)
	individual.AverageFitness = mean
	individual.FitnessStdDev = std
	individual.FitnessVariance = variance
	individual.HasCalculatedFitness = true
	individual.HasAppliedStrategy = true
	individual.Age += 1
	individual.AverageDelta = deltaMean
	return mean
}
