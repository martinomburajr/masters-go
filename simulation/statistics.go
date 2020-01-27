package simulation

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/martinomburajr/masters-go/evolution"
	"math"
	"os"
)

//########################################### BEST INDIVIDUAL ##############
type SimulationBestActualIndividual struct {
	Antagonist            evolution.Individual
	Protagonist           evolution.Individual
	AntagonistGeneration  int
	ProtagonistGeneration int
	AntagonistRun         int
	ProtagonistRun        int
	Correlation           float64
	Covariance            float64
}

// SimulationBestActualIndividuals returns the best actual individuals in the entire simulation
func (s *Simulation) SimulationBestActualIndividuals(params evolution.EvolutionParams) (bestActualIndividuals SimulationBestActualIndividual, err error) {
	if s.SimulationStats == nil {
		return SimulationBestActualIndividual{}, fmt.Errorf(
			"ToRunStats | simulationStats is nil")
	}
	if len(s.SimulationStats) < 0 {
		return SimulationBestActualIndividual{}, fmt.Errorf("ToRunStats | simulationStats is empty")
	}

	antagonist := evolution.Individual{AverageFitness: math.MinInt64}
	protagonist := evolution.Individual{AverageFitness: math.MinInt64}
	topAntGen := 0
	topProGen := 0
	topAntRun := 0
	topProRun := 0
	for i, ss := range s.SimulationStats {
		if ss.TopAntagonist.AverageFitness > antagonist.AverageFitness {
			antagonist = ss.TopAntagonist
			topAntGen = ss.TopAntagonistGeneration
			topAntRun = i
		}
		if ss.TopProtagonist.AverageFitness > protagonist.AverageFitness {
			protagonist = ss.TopProtagonist
			topProGen = ss.TopProtagonistGeneration
			topProRun = i
		}
	}

	bestActualIndividuals =
		SimulationBestActualIndividual{
			Antagonist:            antagonist,
			Protagonist:           protagonist,
			AntagonistGeneration:  topAntGen,
			ProtagonistGeneration: topProGen,
			AntagonistRun:         topAntRun,
			ProtagonistRun:        topProRun,
		}

	return bestActualIndividuals, err
}

// BestIndividualInRun returns a CSV type of the best antagonist and protagonist in the given run,
// as well as the final individuals produced
func (s *Simulation) BestIndividualsInRun(params evolution.EvolutionParams) (runBest RunBestIndividualStatistics, err error) {
	if s.SimulationStats == nil {
		return nil, fmt.Errorf(
			"ToRunStats | simulationStats is nil")
	}
	if len(s.SimulationStats) < 0 {
		return nil, fmt.Errorf("ToRunStats | simulationStats is empty")
	}

	runIndex := params.InternalCount
	runBest = make([]RunBestIndividualStatistic, 1)
	run := s.SimulationStats[runIndex]

	topAntagonistEq, err := run.TopAntagonist.Program.T.ToMathematicalString()
	topProtagonistEq, err := run.TopProtagonist.Program.T.ToMathematicalString()
	finalAntagonistEq, err := run.FinalAntagonist.Program.T.ToMathematicalString()
	finalProtagonistEq, err := run.FinalProtagonist.Program.T.ToMathematicalString()

	runBest[0] = RunBestIndividualStatistic{
		SpecEquation: params.SpecParam.ExpressionParsed,
		SpecRange:    params.SpecParam.Range,
		SpecSeed:     params.SpecParam.Seed,

		Antagonist:                  run.TopAntagonist.AverageFitness,
		Protagonist:                 run.TopProtagonist.AverageFitness,
		AntagonistBestFitness:       run.TopAntagonist.BestFitness,
		ProtagonistBestFitness:      run.TopProtagonist.BestFitness,
		AntagonistStdDev:            run.TopAntagonist.FitnessStdDev,
		ProtagonistStdDev:           run.TopProtagonist.FitnessStdDev,
		AntagonistAverageDelta:      run.TopAntagonist.AverageDelta,
		ProtagonistAverageDelta:     run.TopProtagonist.AverageDelta,
		AntagonistBestDelta:         run.TopAntagonist.BestDelta,
		ProtagonistBestDelta:        run.TopProtagonist.BestDelta,
		AntagonistEquation:          topAntagonistEq,
		ProtagonistEquation:         topProtagonistEq,
		AntagonistStrategy:          evolution.StrategiesToString(run.TopAntagonist),
		ProtagonistStrategy:         evolution.StrategiesToString(run.TopProtagonist),
		AntagonistDominantStrategy:  evolution.DominantStrategy(run.TopAntagonist),
		ProtagonistDominantStrategy: evolution.DominantStrategy(run.TopProtagonist),
		AntagonistBirthGen:          run.TopAntagonist.BirthGen,
		ProtagonistBirthGen:         run.TopProtagonist.BirthGen,
		AntagonistAge:               run.TopAntagonist.Age,
		ProtagonistAge:              run.TopProtagonist.Age,
		AntagonistGeneration:        run.TopAntagonistGeneration,
		ProtagonistGeneration:       run.TopProtagonistGeneration,

		MeanCorrelation: run.MeanCorrelation,
		MeanCovariance:  run.MeanCovariance,

		FinalAntagonist:                  run.FinalAntagonist.AverageFitness,
		FinalProtagonist:                 run.FinalProtagonist.AverageFitness,
		FinalAntagonistBestFitness:       run.FinalAntagonist.BestFitness,
		FinalProtagonistBestFitness:      run.FinalProtagonist.BestFitness,
		FinalAntagonistStdDev:            run.FinalAntagonist.FitnessStdDev,
		FinalProtagonistStdDev:           run.FinalProtagonist.FitnessStdDev,
		FinalAntagonistAverageDelta:      run.FinalAntagonist.AverageDelta,
		FinalProtagonistAverageDelta:     run.FinalProtagonist.AverageDelta,
		FinalAntagonistBestDelta:         run.FinalAntagonist.BestDelta,
		FinalProtagonistBestDelta:        run.FinalProtagonist.BestDelta,
		FinalAntagonistEquation:          finalAntagonistEq,
		FinalProtagonistEquation:         finalProtagonistEq,
		FinalAntagonistStrategy:          evolution.StrategiesToString(run.FinalAntagonist),
		FinalProtagonistStrategy:         evolution.StrategiesToString(run.FinalProtagonist),
		FinalAntagonistDominantStrategy:  evolution.DominantStrategy(run.FinalAntagonist),
		FinalProtagonistDominantStrategy: evolution.DominantStrategy(run.FinalProtagonist),
		FinalAntagonistBirthGen:          run.FinalAntagonist.BirthGen,
		FinalProtagonistBirthGen:         run.FinalProtagonist.BirthGen,
		FinalAntagonistAge:               run.FinalAntagonist.Age,
		FinalProtagonistAge:              run.FinalProtagonist.Age,
		Run:                              runIndex,
	}

	return runBest, err
}

// SimulationBestIndividual returns the best antagonist and protagonist in all runs
func (s *Simulation) SimulationBestIndividuals(params evolution.EvolutionParams) (
	simulationBestIndividuals SimulationBestIndividuals,
	err error) {
	if s.SimulationStats == nil {
		return nil, fmt.Errorf(
			"ToRunStats | simulationStats is nil")
	}
	if len(s.SimulationStats) < 0 {
		return nil, fmt.Errorf("ToRunStats | simulationStats is empty")
	}

	simulationBestIndividuals = make([]SimulationBestIndividual, len(s.SimulationStats))

	for i, run := range s.SimulationStats {
		topAntagonist := run.TopAntagonist
		topProtagonist := run.TopProtagonist
		topAntGen := run.TopAntagonistGeneration
		topProGen := run.TopProtagonistGeneration
		topAntRun := i
		topProRun := i

		topAntagonistEq, _ := topAntagonist.Program.T.ToMathematicalString()
		topProtagonistEq, _ := topProtagonist.Program.T.ToMathematicalString()
		simulationBestIndividuals[i] = SimulationBestIndividual{
			SpecEquation:                params.SpecParam.ExpressionParsed,
			SpecRange:                   params.SpecParam.Range,
			SpecSeed:                    params.SpecParam.Seed,
			Covariance:                  0,
			AntagonistID:                topAntagonist.Id,
			ProtagonistID:               topProtagonist.Id,
			Antagonist:                  topAntagonist.AverageFitness,
			Protagonist:                 topProtagonist.AverageFitness,
			AntagonistBestFitness:       topAntagonist.BestFitness,
			ProtagonistBestFitness:      topProtagonist.BestFitness,
			AntagonistStdDev:            topAntagonist.FitnessStdDev,
			ProtagonistStdDev:           topProtagonist.FitnessStdDev,
			AntagonistAverageDelta:      topAntagonist.AverageDelta,
			ProtagonistAverageDelta:     topProtagonist.AverageDelta,
			AntagonistBestDelta:         topAntagonist.BestDelta,
			ProtagonistBestDelta:        topProtagonist.BestDelta,
			AntagonistEquation:          topAntagonistEq,
			ProtagonistEquation:         topProtagonistEq,
			AntagonistStrategy:          evolution.StrategiesToString(topAntagonist),
			ProtagonistStrategy:         evolution.StrategiesToString(topProtagonist),
			AntagonistDominantStrategy:  evolution.DominantStrategy(topAntagonist),
			ProtagonistDominantStrategy: evolution.DominantStrategy(topProtagonist),
			AntagonistGeneration:        topAntGen,
			ProtagonistGeneration:       topProGen,
			AntagonistBirthGen:          topAntagonist.BirthGen,
			ProtagonistBirthGen:         topProtagonist.BirthGen,
			AntagonistAge:               topAntagonist.Age,
			ProtagonistAge:              topProtagonist.Age,
			AntagonistRun:               topAntRun,
			ProtagonistRun:              topProRun,
		}
	}

	return simulationBestIndividuals, err
}

// SimulationBestIndividual returns the best antagonist and protagonist in all runs
func (s *Simulation) SimulationBestIndividual(params evolution.EvolutionParams) (simulationBestIndividuals SimulationBestIndividuals,
	err error) {
	if s.SimulationStats == nil {
		return nil, fmt.Errorf(
			"ToRunStats | simulationStats is nil")
	}
	if len(s.SimulationStats) < 0 {
		return nil, fmt.Errorf("ToRunStats | simulationStats is empty")
	}

	simulationBestIndividuals = make([]SimulationBestIndividual, 1)

	topAntagonist := evolution.Individual{AverageFitness: math.MinInt64}
	topProtagonist := evolution.Individual{AverageFitness: math.MinInt64}
	topAntGen := 0
	topProGen := 0
	topAntRun := 0
	topProRun := 0
	for i, ss := range s.SimulationStats {
		if ss.TopAntagonist.AverageFitness > topAntagonist.AverageFitness {
			topAntagonist = ss.TopAntagonist
			topAntGen = ss.TopAntagonistGeneration
			topAntRun = i
		}
		if ss.TopProtagonist.AverageFitness > topProtagonist.AverageFitness {
			topProtagonist = ss.TopProtagonist
			topProGen = ss.TopProtagonistGeneration
			topProRun = i
		}
	}

	topAntagonistEq, _ := topAntagonist.Program.T.ToMathematicalString()
	topProtagonistEq, _ := topProtagonist.Program.T.ToMathematicalString()

	simulationBestIndividuals[0] = SimulationBestIndividual{
		SpecEquation:                params.SpecParam.ExpressionParsed,
		SpecRange:                   params.SpecParam.Range,
		SpecSeed:                    params.SpecParam.Seed,
		AntagonistID:                topAntagonist.Id,
		ProtagonistID:               topProtagonist.Id,
		Antagonist:                  topAntagonist.AverageFitness,
		Protagonist:                 topProtagonist.AverageFitness,
		AntagonistBestFitness:       topAntagonist.BestFitness,
		ProtagonistBestFitness:      topProtagonist.BestFitness,
		AntagonistStdDev:            topAntagonist.FitnessStdDev,
		ProtagonistStdDev:           topProtagonist.FitnessStdDev,
		AntagonistBestDelta:         topAntagonist.BestDelta,
		ProtagonistBestDelta:        topProtagonist.BestDelta,
		AntagonistAverageDelta:      topAntagonist.AverageDelta,
		ProtagonistAverageDelta:     topProtagonist.AverageDelta,
		AntagonistEquation:          topAntagonistEq,
		ProtagonistEquation:         topProtagonistEq,
		AntagonistStrategy:          evolution.StrategiesToString(topAntagonist),
		ProtagonistStrategy:         evolution.StrategiesToString(topProtagonist),
		AntagonistDominantStrategy:  evolution.DominantStrategy(topAntagonist),
		ProtagonistDominantStrategy: evolution.DominantStrategy(topProtagonist),
		AntagonistBirthGen:          topAntagonist.BirthGen,
		ProtagonistBirthGen:         topProtagonist.BirthGen,
		AntagonistAge:               topAntagonist.Age,
		ProtagonistAge:              topProtagonist.Age,
		AntagonistGeneration:        topAntGen,
		ProtagonistGeneration:       topProGen,
		AntagonistRun:               topAntRun,
		ProtagonistRun:              topProRun,
	}

	return simulationBestIndividuals, err
}

//########################################### GENERATIONAL ##############

// BestIndividualInRun returns a CSV type of the best antagonist and protagonist in the given run,
// as well as the final individuals produced
func (s *Simulation) GenerationalInRun(params evolution.EvolutionParams) (runGen RunGenerationalStatistics, err error) {
	runIndex := params.InternalCount
	if s.SimulationStats == nil {
		return nil, fmt.Errorf(
			"ToRunStats | simulationStats is nil")
	}
	if len(s.SimulationStats) < 0 {
		return nil, fmt.Errorf("ToRunStats | simulationStats is empty")
	}
	if runIndex > len(s.SimulationStats) {
		runIndex = len(s.SimulationStats) - 1
	}
	if runIndex < 0 {
		runIndex = 0
	}

	genCount := evolution.CalculateGenerationSize(params)

	runGen = make([]RunGenerationalStatistic, genCount)
	run := s.SimulationStats[runIndex]

	for i := 0; i < genCount; i++ {
		topAntagonistInGenerationByAvgFitness := run.Generational.BestAntagonistInEachGenerationByAvgFitness[i]
		topProtagonistInGenerationByAvgFitness := run.Generational.BestProtagonistInEachGenerationByAvgFitness[i]
		AntagonistEq, _ := topAntagonistInGenerationByAvgFitness.Program.T.ToMathematicalString()
		ProtagonistEq, _ := topProtagonistInGenerationByAvgFitness.Program.T.ToMathematicalString()

		runGen[i] = RunGenerationalStatistic{
			Generation:          i,
			SpecEquation:        params.SpecParam.ExpressionParsed,
			SpecRange:           params.SpecParam.Range,
			SpecSeed:            params.SpecParam.Seed,
			AntagonistEquation:  AntagonistEq,
			ProtagonistEquation: ProtagonistEq,

			// Generational Stats
			Correlation:           run.Generational.CorrelationInEachGeneration[i],
			Covariance:            run.Generational.CovarianceInEachGeneration[i],
			AntagonistMean:        run.Generational.AntagonistAverageInEachGeneration[i],
			ProtagonistMean:       run.Generational.ProtagonistAverageInEachGeneration[i],
			AntagonistVariance:    run.Generational.AntagonistVarianceInEachGeneration[i],
			ProtagonistVariance:   run.Generational.ProtagonistVarianceInEachGeneration[i],
			AntagonistSkew:        run.Generational.AntagonistSkewInEachGeneration[i],
			ProtagonistSkew:       run.Generational.ProtagonistSkewInEachGeneration[i],
			AntagonistExKurtosis:  run.Generational.AntagonistExKurtosisInEachGeneration[i],
			ProtagonistExKurtosis: run.Generational.ProtagonistExKurtosisInEachGeneration[i],

			// Best Individual in Generation Stats
			TopAntagonistMeanFitness:  topAntagonistInGenerationByAvgFitness.AverageFitness,
			TopProtagonistMeanFitness: topProtagonistInGenerationByAvgFitness.AverageFitness,
			AntagonistBestFitness:     topAntagonistInGenerationByAvgFitness.BestFitness,
			ProtagonistBestFitness:    topProtagonistInGenerationByAvgFitness.BestFitness,
			AntagonistStdDev:          topAntagonistInGenerationByAvgFitness.FitnessStdDev,
			ProtagonistStdDev:         topProtagonistInGenerationByAvgFitness.FitnessStdDev,
			AntagonistAverageDelta:    topAntagonistInGenerationByAvgFitness.AverageDelta,
			ProtagonistAverageDelta:   topProtagonistInGenerationByAvgFitness.AverageDelta,
			AntagonistBestDelta:       topAntagonistInGenerationByAvgFitness.BestDelta,
			ProtagonistBestDelta:      topProtagonistInGenerationByAvgFitness.BestDelta,

			AntagonistStrategy:          evolution.StrategiesToString(topAntagonistInGenerationByAvgFitness),
			ProtagonistStrategy:         evolution.StrategiesToString(topProtagonistInGenerationByAvgFitness),
			AntagonistDominantStrategy:  evolution.DominantStrategy(topAntagonistInGenerationByAvgFitness),
			ProtagonistDominantStrategy: evolution.DominantStrategy(topProtagonistInGenerationByAvgFitness),
			AntagonistBirthGen:          topAntagonistInGenerationByAvgFitness.BirthGen,
			ProtagonistBirthGen:         topProtagonistInGenerationByAvgFitness.BirthGen,
			AntagonistAge:               topAntagonistInGenerationByAvgFitness.Age,
			ProtagonistAge:              topProtagonistInGenerationByAvgFitness.Age,

			Run: runIndex,
		}
	}

	return runGen, err
}

// ######################################## EPOCHAL ################

func (s *Simulation) EpochalInRun(params evolution.EvolutionParams) (runEpochal RunEpochalStatistics, err error) {
	runIndex := params.InternalCount
	antagonist := s.SimulationStats[runIndex].TopAntagonist
	protagonist := s.SimulationStats[runIndex].TopProtagonist
	finalAntagonist := s.SimulationStats[runIndex].FinalAntagonist
	finalProtagonist := s.SimulationStats[runIndex].FinalProtagonist

	topAntagonistEq, err := antagonist.Program.T.ToMathematicalString()
	topProtagonistEq, err := protagonist.Program.T.ToMathematicalString()
	finalAntagonistEq, err := finalAntagonist.Program.T.ToMathematicalString()
	finalProtagonistEq, err := finalProtagonist.Program.T.ToMathematicalString()

	epochLength := params.EachPopulationSize
	runEpochal = make([]RunEpochalStatistic, epochLength)
	for i := 0; i < epochLength; i++ {
		runEpochal[i] = RunEpochalStatistic{
			Epoch:        i,
			Run:          runIndex,
			SpecEquation: params.SpecParam.ExpressionParsed,
			SpecRange:    params.SpecParam.Range,
			SpecSeed:     params.SpecParam.Seed,

			AntagonistID:                antagonist.Id,
			ProtagonistID:               protagonist.Id,
			Antagonist:                  antagonist.Fitness[i],
			Protagonist:                 protagonist.Fitness[i],
			AntagonistDelta:             antagonist.Deltas[i],
			ProtagonistDelta:            protagonist.Deltas[i],
			AntagonistStdDev:            antagonist.FitnessStdDev,
			ProtagonistStdDev:           protagonist.FitnessStdDev,
			AntagonistAverageDelta:      antagonist.AverageDelta,
			ProtagonistAverageDelta:     protagonist.AverageDelta,
			AntagonistBestDelta:         antagonist.BestDelta,
			ProtagonistBestDelta:        protagonist.BestDelta,
			AntagonistEquation:          topAntagonistEq,
			ProtagonistEquation:         topProtagonistEq,
			AntagonistStrategy:          evolution.StrategiesToString(antagonist),
			ProtagonistStrategy:         evolution.StrategiesToString(protagonist),
			AntagonistDominantStrategy:  evolution.DominantStrategy(antagonist),
			ProtagonistDominantStrategy: evolution.DominantStrategy(protagonist),
			AntagonistGeneration:        s.SimulationStats[runIndex].TopAntagonistGeneration,
			ProtagonistGeneration:       s.SimulationStats[runIndex].TopProtagonistGeneration,

			FinalAntagonist:                  finalAntagonist.Fitness[i],
			FinalProtagonist:                 finalProtagonist.Fitness[i],
			FinalAntagonistStdDev:            finalAntagonist.FitnessStdDev,
			FinalProtagonistStdDev:           finalProtagonist.FitnessStdDev,
			FinalAntagonistDelta:             finalAntagonist.Deltas[i],
			FinalProtagonistDelta:            finalProtagonist.Deltas[i],
			FinalAntagonistBestDelta:         finalAntagonist.BestDelta,
			FinalProtagonistBestDelta:        finalProtagonist.BestDelta,
			FinalAntagonistEquation:          finalAntagonistEq,
			FinalProtagonistEquation:         finalProtagonistEq,
			FinalAntagonistStrategy:          evolution.StrategiesToString(finalAntagonist),
			FinalProtagonistStrategy:         evolution.StrategiesToString(finalProtagonist),
			FinalAntagonistDominantStrategy:  evolution.DominantStrategy(finalAntagonist),
			FinalProtagonistDominantStrategy: evolution.DominantStrategy(finalProtagonist),
		}
	}

	return runEpochal, err
}

func (s *Simulation) SimulationBestEpochal(params evolution.EvolutionParams) (bestEpochs SimulationBestEpochs, err error) {
	if s.SimulationStats == nil {
		return nil, fmt.Errorf(
			"ToRunStats | simulationStats is nil")
	}
	if len(s.SimulationStats) < 0 {
		return nil, fmt.Errorf("ToRunStats | simulationStats is empty")
	}

	//epochLength := params.InternalCount
	bestEpochs = make([]SimulationBestEpoch, params.InternalCount)

	//for j, ss := range s.SimulationStats {
	//	antagonist := evolution.Individual{AverageFitness: math.MinInt64}
	//	protagonist := evolution.Individual{AverageFitness: math.MinInt64}
	//	topAntGen := 0
	//	topProGen := 0
	//	topAntRun := 0
	//	topProRun := 0
	//
	//	topAntagonistEq, _ := antagonist.Program.T.ToMathematicalString()
	//	topProtagonistEq, _ := protagonist.Program.T.ToMathematicalString()
	//
	//	for i := 0; i < epochLength; i++ {
	//
	//		bestEpochs[i] = SimulationBestEpoch{
	//			Epoch:        i,
	//			SpecEquation: params.SpecParam.ParsedExpression,
	//			SpecRange:    params.SpecParam.Range,
	//			SpecSeed:     params.SpecParam.Seed,
	//
	//			Antagonist:                  antagonist.AverageFitness,
	//			Protagonist:                 protagonist.AverageFitness,
	//			AntagonistStdDevInGeneration:            antagonist.FitnessStdDev,
	//			ProtagonistStdDevInGeneration:           protagonist.FitnessStdDev,
	//			AntagonistBestDelta:         antagonist.BestDelta,
	//			ProtagonistBestDelta:        protagonist.BestDelta,
	//			AntagonistEquation:          topAntagonistEq,
	//			ProtagonistEquation:         topProtagonistEq,
	//			AntagonistStrategy:          evolution.StrategiesToString(antagonist),
	//			ProtagonistStrategy:         evolution.StrategiesToString(protagonist),
	//			AntagonistDominantStrategy:  evolution.DominantStrategy(antagonist),
	//			ProtagonistDominantStrategy: evolution.DominantStrategy(protagonist),
	//			AntagonistGeneration:        topAntGen,
	//			ProtagonistGeneration:       topProGen,
	//			AntagonistRun:               topAntRun,
	//			ProtagonistRun:              topProRun,
	//		}
	//	}
	//}

	return bestEpochs, err
}

// ######################################## STRATEGICAL ################

func (s *Simulation) StrategyInRun(params evolution.EvolutionParams) (runStrategy RunStrategyStatistics, err error) {

	runIndex := params.InternalCount
	if s.SimulationStats == nil {
		return nil, fmt.Errorf("ToRunStats | simulationStats is nil")
	}
	if len(s.SimulationStats) < 0 {
		return nil, fmt.Errorf("ToRunStats | simulationStats is empty")
	}
	if runIndex > len(s.SimulationStats) {
		runIndex = len(s.SimulationStats) - 1
	}
	if runIndex < 0 {
		runIndex = 0
	}

	run := s.SimulationStats[runIndex]
	protagonist := run.TopProtagonist
	antagonist := run.TopAntagonist
	finalAntagonist := run.FinalAntagonist
	finalProtagonist := run.FinalProtagonist

	strategyLength := len(antagonist.Strategy)
	if strategyLength < len(protagonist.Strategy) {
		strategyLength = len(protagonist.Strategy)
	}
	if strategyLength < len(finalAntagonist.Strategy) {
		strategyLength = len(finalAntagonist.Strategy)
	}
	if strategyLength < len(finalProtagonist.Strategy) {
		strategyLength = len(finalProtagonist.Strategy)
	}
	runStrategy = make([]RunStrategyStatistic, strategyLength)

	for j := 0; j < strategyLength; j++ {
		antStrat := ""
		proStrat := ""
		finAntStrat := ""
		finProStrat := ""
		if j < len(antagonist.Strategy) {
			antStrat = string(antagonist.Strategy[j])
		}
		if j < len(protagonist.Strategy) {
			proStrat = string(protagonist.Strategy[j])
		}
		if j < len(finalAntagonist.Strategy) {
			finAntStrat = string(finalAntagonist.Strategy[j])
		}
		if j < len(finalProtagonist.Strategy) {
			finProStrat = string(finalProtagonist.Strategy[j])
		}
		runStrategy[j] = RunStrategyStatistic{
			Antagonist:            antStrat,
			Protagonist:           proStrat,
			FinalAntagonist:       finAntStrat,
			FinalProtagonist:      finProStrat,
			AntagonistGeneration:  run.TopAntagonistGeneration,
			ProtagonistGeneration: run.TopProtagonistGeneration,
			StrategyCount:         j,
			Run:                   runIndex,
		}
	}

	return runStrategy, err
}

func (s *Simulation) SimulationBestStrategy(params evolution.EvolutionParams) (simulationStrategy SimulationStrategyStatistics,
	err error) {
	if s.SimulationStats == nil {
		return nil, fmt.Errorf("ToRunStats | simulationStats is nil")
	}
	if len(s.SimulationStats) < 0 {
		return nil, fmt.Errorf("ToRunStats | simulationStats is empty")
	}

	bestActualIndividuals, err := s.SimulationBestActualIndividuals(params)
	if err != nil {
		return nil, err
	}

	protagonist := bestActualIndividuals.Protagonist
	antagonist := bestActualIndividuals.Antagonist

	strategyLength := len(antagonist.Strategy)
	if strategyLength < len(protagonist.Strategy) {
		strategyLength = len(protagonist.Strategy)
	}
	simulationStrategy = make([]SimulationStrategyStatistic, strategyLength)

	for j := 0; j < strategyLength; j++ {
		antStrat := ""
		proStrat := ""
		if j < len(antagonist.Strategy) {
			antStrat = string(antagonist.Strategy[j])
		}
		if j < len(protagonist.Strategy) {
			proStrat = string(protagonist.Strategy[j])
		}
		antEq, _ := antagonist.Program.T.ToMathematicalString()
		proEq, _ := protagonist.Program.T.ToMathematicalString()
		simulationStrategy[j] = SimulationStrategyStatistic{
			Antagonist:            antStrat,
			Protagonist:           proStrat,
			AntagonistEquation:    antEq,
			ProtagonistEquation:   proEq,
			ProtagonistRun:        bestActualIndividuals.ProtagonistRun,
			AntagonistRun:         bestActualIndividuals.AntagonistRun,
			ProtagonistGeneration: bestActualIndividuals.ProtagonistGeneration,
			AntagonistGeneration:  bestActualIndividuals.AntagonistGeneration,
			StrategyCount:         j,
		}
	}

	return simulationStrategy, err
}

type RunBasedStatistics struct {
	TopAntagonist          float64 `csv:"runTopA"`
	TopProtagonist         float64 `csv:"runTopP"`
	TopAntagonistDelta     float64 `csv:"runTopADelta"`
	TopProtagonistDelta    float64 `csv:"runTopPDelta"`
	TopAntagonistStrategy  string  `csv:"runTopAStrategy"`
	TopProtagonistStrategy string  `csv:"runTopPStrategy"`
	TopAntagonistEquation  string  `csv:"runTopAEquation"`
	TopProtagonistEquation string  `csv:"runTopPEquation"`

	FinalAntagonist          float64 `csv:"runFinalA"`
	FinalProtagonist         float64 `csv:"runFinalP"`
	FinalAntagonistDelta     float64 `csv:"runFinalADelta"`
	FinalProtagonistDelta    float64 `csv:"runFinalPDelta"`
	FinalAntagonistStrategy  string  `csv:"runFinalAStrategy"`
	FinalProtagonistStrategy string  `csv:"runFinalPStrategy"`
	FinalAntagonistEquation  string  `csv:"runFinalAEquation"`
	FinalProtagonistEquation string  `csv:"runFinalPEquation"`

	Run int `csv:"run"`
}

type GenerationalAverages struct {
	Antagonists  []evolution.Individual `csv:"A"`
	Protagonists []evolution.Individual `csv:"A"`
}

// ############################################# DATA TYPES ##########################################################

// RunGenerationalStatistics refer to statistics per generation.
// So Top or Bottom refer to the best or worst in the given generation and not a cumulative of the evolutionary process.
type RunGenerationalStatistic struct {
	Generation int `csv:"gen"`

	SpecEquation string `csv:"specEquation"`
	SpecRange    int    `csv:"range"`
	SpecSeed     int    `csv:"seed"`

	AntagonistEquation  string `csv:"topAEquation"`
	ProtagonistEquation string `csv:"topPEquation"`

	// Generation Stats
	Correlation float64 `csv:"correlation"`
	Covariance  float64 `csv:"covariance"`

	AntagonistMean        float64 `csv:"AMean"`
	ProtagonistMean       float64 `csv:"PMean"`
	AntagonistStdDev      float64 `csv:"AStd"`
	ProtagonistStdDev     float64 `csv:"PStd"`
	AntagonistVariance    float64 `csv:"AVar"`
	ProtagonistVariance   float64 `csv:"PVar"`
	AntagonistSkew        float64 `csv:"ASkew"`
	ProtagonistSkew       float64 `csv:"PSkew"`
	AntagonistExKurtosis  float64 `csv:"AExKur"`
	ProtagonistExKurtosis float64 `csv:"PExKur"`

	// Top Individual In Generation Stats
	TopAntagonistMeanFitness    float64 `csv:"topAMean"`
	TopProtagonistMeanFitness   float64 `csv:"topPMean"`
	AntagonistBestFitness       float64 `csv:"topABest"`
	ProtagonistBestFitness      float64 `csv:"topPBest"`
	AntagonistAverageDelta      float64 `csv:"topAMeanDelta"`
	ProtagonistAverageDelta     float64 `csv:"topPMeanDelta"`
	AntagonistBestDelta         float64 `csv:"topABestDelta"`
	ProtagonistBestDelta        float64 `csv:"topPBestDelta"`
	AntagonistStrategy          string  `csv:"topAStrat"`
	ProtagonistStrategy         string  `csv:"topPStrat"`
	AntagonistDominantStrategy  string  `csv:"topADomStrat"`
	ProtagonistDominantStrategy string  `csv:"topPDomStrat"`
	AntagonistBirthGen          int     `csv:"topABirthGen"`
	ProtagonistBirthGen         int     `csv:"topPBirthGen"`
	AntagonistAge               int     `csv:"topAAge"`
	ProtagonistAge              int     `csv:"topPAge"`

	Run int `csv:"run"`
}
type RunGenerationalStatistics []RunGenerationalStatistic

func (g *RunGenerationalStatistics) ToCSV(outputPath string) error {
	outputFileCSV, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFileCSV.Close()

	writer := gocsv.DefaultCSVWriter(outputFileCSV)
	if writer.Error() != nil {
		return writer.Error()
	}
	err = gocsv.Marshal(g, outputFileCSV)
	if err != nil {
		return err
	}
	return nil
}

type RunEpochalStatistic struct {
	SpecEquation string `csv:"specEquation"`
	SpecRange    int    `csv:"range"`
	SpecSeed     int    `csv:"seed"`

	AntagonistID      string  `csv:"AID"`
	ProtagonistID     string  `csv:"PID"`
	Antagonist        float64 `csv:"A"`
	Protagonist       float64 `csv:"P"`
	AntagonistDelta   float64 `csv:"ADelta"`
	ProtagonistDelta  float64 `csv:"PDelta"`
	AntagonistStdDev  float64 `csv:"AStdDev"`
	ProtagonistStdDev float64 `csv:"PStdDev"`

	AntagonistAverageDelta      float64 `csv:"AAvgDelta"`
	ProtagonistAverageDelta     float64 `csv:"PAvgDelta"`
	AntagonistBestDelta         float64 `csv:"ABestDelta"`
	ProtagonistBestDelta        float64 `csv:"PBestDelta"`
	AntagonistEquation          string  `csv:"AEquation"`
	ProtagonistEquation         string  `csv:"PEquation"`
	AntagonistStrategy          string  `csv:"AStrat"`
	ProtagonistStrategy         string  `csv:"PStrat"`
	AntagonistDominantStrategy  string  `csv:"ADomStrat"`
	ProtagonistDominantStrategy string  `csv:"PDomStrat"`
	AntagonistGeneration        int     `csv:"AGen"`
	ProtagonistGeneration       int     `csv:"PGen"`

	FinalAntagonist                  float64 `csv:"finA"`
	FinalProtagonist                 float64 `csv:"finP"`
	FinalAntagonistStdDev            float64 `csv:"finAStdDev"`
	FinalProtagonistStdDev           float64 `csv:"finPStdDev"`
	FinalAntagonistDelta             float64 `csv:"finADelta"`
	FinalProtagonistDelta            float64 `csv:"finPDelta"`
	FinalAntagonistBestDelta         float64 `csv:"finABestDelta"`
	FinalProtagonistBestDelta        float64 `csv:"finPBestDelta"`
	FinalAntagonistEquation          string  `csv:"finAEquation"`
	FinalProtagonistEquation         string  `csv:"finPEquation"`
	FinalAntagonistStrategy          string  `csv:"finAStrat"`
	FinalProtagonistStrategy         string  `csv:"finPStrat"`
	FinalAntagonistDominantStrategy  string  `csv:"finADomStrat"`
	FinalProtagonistDominantStrategy string  `csv:"finPDomStrat"`

	Epoch int `csv:"epoch"`
	Run   int `csv:"run"`
}

type RunEpochalStatistics []RunEpochalStatistic

func (e *RunEpochalStatistics) ToCSV(outputPath string) error {
	outputFileCSV, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFileCSV.Close()

	writer := gocsv.DefaultCSVWriter(outputFileCSV)
	if writer.Error() != nil {
		return writer.Error()
	}
	err = gocsv.Marshal(e, outputFileCSV)
	if err != nil {
		return err
	}
	return nil
}

type RunStrategyStatistic struct {
	Antagonist            string `csv:"A"`
	Protagonist           string `csv:"P"`
	FinalAntagonist       string `csv:"AFinal"`
	FinalProtagonist      string `csv:"PFinal"`
	AntagonistGeneration  int    `csv:"AGen"`
	ProtagonistGeneration int    `csv:"PGen"`
	StrategyCount         int    `csv:"count"`
	Run                   int    `csv:"run"`
}

type RunStrategyStatistics []RunStrategyStatistic

func (e *RunStrategyStatistics) ToCSV(outputPath string) error {
	outputFileCSV, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFileCSV.Close()

	writer := gocsv.DefaultCSVWriter(outputFileCSV)
	if writer.Error() != nil {
		return writer.Error()
	}
	err = gocsv.Marshal(e, outputFileCSV)
	if err != nil {
		return err
	}
	return nil
}

type SimulationStrategyStatistic struct {
	Antagonist            string `csv:"A"`
	Protagonist           string `csv:"P"`
	AntagonistEquation    string `csv:"AEquation"`
	ProtagonistEquation   string `csv:"PEquation"`
	AntagonistGeneration  int    `csv:"AGen"`
	ProtagonistGeneration int    `csv:"PGen"`
	AntagonistRun         int    `csv:"ARun"`
	ProtagonistRun        int    `csv:"PRun"`
	StrategyCount         int    `csv:"count"`
	Run                   int    `csv:"run"`
}

type SimulationStrategyStatistics []SimulationStrategyStatistic

func (e *SimulationStrategyStatistics) ToCSV(outputPath string) error {
	outputFileCSV, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFileCSV.Close()

	writer := gocsv.DefaultCSVWriter(outputFileCSV)
	if writer.Error() != nil {
		return writer.Error()
	}
	err = gocsv.Marshal(e, outputFileCSV)
	if err != nil {
		return err
	}
	return nil
}

type RunBestIndividualStatistic struct {
	SpecEquation string `csv:"specEquation"`
	SpecRange    int    `csv:"range"`
	SpecSeed     int    `csv:"seed"`

	AntagonistID                string  `csv:"AID"`
	ProtagonistID               string  `csv:"PID"`
	Antagonist                  float64 `csv:"AAvg"`
	Protagonist                 float64 `csv:"PAvg"`
	MeanCorrelation             float64 `csv:"meanCorrInRun"`
	MeanCovariance              float64 `csv:"meanCovInRun"`
	AntagonistBestFitness       float64 `csv:"ABestFit"`
	ProtagonistBestFitness      float64 `csv:"PBestFit"`
	AntagonistStdDev            float64 `csv:"AStdDev"`
	ProtagonistStdDev           float64 `csv:"PStdDev"`
	AntagonistAverageDelta      float64 `csv:"AAvgDelta"`
	ProtagonistAverageDelta     float64 `csv:"PAvgDelta"`
	AntagonistBestDelta         float64 `csv:"ABestDelta"`
	ProtagonistBestDelta        float64 `csv:"PBestDelta"`
	AntagonistEquation          string  `csv:"AEquation"`
	ProtagonistEquation         string  `csv:"PEquation"`
	AntagonistStrategy          string  `csv:"AStrat"`
	ProtagonistStrategy         string  `csv:"PStrat"`
	AntagonistDominantStrategy  string  `csv:"ADomStrat"`
	ProtagonistDominantStrategy string  `csv:"PDomStrat"`
	AntagonistGeneration        int     `csv:"AGen"`
	ProtagonistGeneration       int     `csv:"PGen"`
	AntagonistBirthGen          int     `csv:"ABirthGen"`
	ProtagonistBirthGen         int     `csv:"PBirthGen"`
	AntagonistAge               int     `csv:"AAge"`
	ProtagonistAge              int     `csv:"PAge"`

	FinalAntagonist                  float64 `csv:"finAAvg"`
	FinalProtagonist                 float64 `csv:"finPAvg"`
	FinalAntagonistBestFitness       float64 `csv:"finABestFit"`
	FinalProtagonistBestFitness      float64 `csv:"finPBestFit"`
	FinalAntagonistStdDev            float64 `csv:"finAStdDev"`
	FinalProtagonistStdDev           float64 `csv:"finPStdDev"`
	FinalAntagonistAverageDelta      float64 `csv:"finAAvgDelta"`
	FinalProtagonistAverageDelta     float64 `csv:"finPAvgDelta"`
	FinalAntagonistBestDelta         float64 `csv:"finABestDelta"`
	FinalProtagonistBestDelta        float64 `csv:"finPBestDelta"`
	FinalAntagonistEquation          string  `csv:"finAEquation"`
	FinalProtagonistEquation         string  `csv:"finPEquation"`
	FinalAntagonistStrategy          string  `csv:"finAStrat"`
	FinalProtagonistStrategy         string  `csv:"finPStrat"`
	FinalAntagonistDominantStrategy  string  `csv:"finADomStrat"`
	FinalProtagonistDominantStrategy string  `csv:"finPDomStrat"`
	FinalAntagonistBirthGen          int     `csv:"finABirthGen"`
	FinalProtagonistBirthGen         int     `csv:"finPBirthGen"`
	FinalAntagonistAge               int     `csv:"finAAge"`
	FinalProtagonistAge              int     `csv:"finPAge"`

	Run int `csv:"run"`
}

type RunBestIndividualStatistics []RunBestIndividualStatistic

func (e *RunBestIndividualStatistics) ToCSV(outputPath string) error {
	outputFileCSV, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFileCSV.Close()

	writer := gocsv.DefaultCSVWriter(outputFileCSV)
	if writer.Error() != nil {
		return writer.Error()
	}
	err = gocsv.Marshal(e, outputFileCSV)
	if err != nil {
		return err
	}
	return nil
}

type SimulationBestIndividual struct {
	SpecEquation string `csv:"specEquation"`
	SpecRange    int    `csv:"range"`
	SpecSeed     int    `csv:"seed"`

	Correlation int `csv:"corr"`
	Covariance  int `csv:"cov"`

	AntagonistID                string  `csv:"AID"`
	ProtagonistID               string  `csv:"PID"`
	Antagonist                  float64 `csv:"AAvg"`
	Protagonist                 float64 `csv:"PAvg"`
	AntagonistBestFitness       float64 `csv:"ABestFit"`
	ProtagonistBestFitness      float64 `csv:"PBestFit"`
	AntagonistStdDev            float64 `csv:"AStdDev"`
	ProtagonistStdDev           float64 `csv:"PStdDev"`
	AntagonistAverageDelta      float64 `csv:"AAverageDelta"`
	ProtagonistAverageDelta     float64 `csv:"PAverageDelta"`
	AntagonistBestDelta         float64 `csv:"ABestDelta"`
	ProtagonistBestDelta        float64 `csv:"PBestDelta"`
	AntagonistEquation          string  `csv:"AEquation"`
	ProtagonistEquation         string  `csv:"PEquation"`
	AntagonistStrategy          string  `csv:"AStrat"`
	ProtagonistStrategy         string  `csv:"PStrat"`
	AntagonistDominantStrategy  string  `csv:"ADomStrat"`
	ProtagonistDominantStrategy string  `csv:"PDomStrat"`
	AntagonistGeneration        int     `csv:"AGen"`
	ProtagonistGeneration       int     `csv:"PGen"`
	AntagonistBirthGen          int     `csv:"ABirthGen"`
	ProtagonistBirthGen         int     `csv:"PBirthGen"`
	AntagonistAge               int     `csv:"AAge"`
	ProtagonistAge              int     `csv:"PAge"`
	AntagonistRun               int     `csv:"ARun"`
	ProtagonistRun              int     `csv:"PRun"`
}

type SimulationBestIndividuals []SimulationBestIndividual

func (s *SimulationBestIndividuals) ToCSV(outputPath string) error {
	outputFileCSV, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFileCSV.Close()

	writer := gocsv.DefaultCSVWriter(outputFileCSV)
	if writer.Error() != nil {
		return writer.Error()
	}
	err = gocsv.Marshal(s, outputFileCSV)
	if err != nil {
		return err
	}

	return err
}

type SimulationBestEpoch struct {
	SpecEquation string `csv:"specEquation"`
	SpecRange    int    `csv:"range"`
	SpecSeed     int    `csv:"seed"`

	Epoch                       int     `csv:"epoch"`
	Antagonist                  float64 `csv:"AAvg"`
	Protagonist                 float64 `csv:"PAvg"`
	AntagonistBestDelta         float64 `csv:"ABestDelta"`
	ProtagonistBestDelta        float64 `csv:"PBestDelta"`
	AntagonistAverageDelta      float64 `csv:"PAverageDelta"`
	ProtagonistAverageDelta     float64 `csv:"PAverageDelta"`
	AntagonistBestFitness       float64 `csv:"ABestFit"`
	ProtagonistBestFitness      float64 `csv:"PBestFit"`
	AntagonistStdDev            float64 `csv:"AStdDev"`
	ProtagonistStdDev           float64 `csv:"PStdDev"`
	AntagonistEquation          string  `csv:"AEquation"`
	ProtagonistEquation         string  `csv:"PEquation"`
	AntagonistStrategy          string  `csv:"AStrat"`
	ProtagonistStrategy         string  `csv:"PStrat"`
	AntagonistDominantStrategy  string  `csv:"ADomStrat"`
	ProtagonistDominantStrategy string  `csv:"PDomStrat"`
	AntagonistGeneration        int     `csv:"AGen"`
	ProtagonistGeneration       int     `csv:"PGen"`
	AntagonistBirthGen          int     `csv:"ABirthGen"`
	ProtagonistBirthGen         int     `csv:"PBirthGen"`
	AntagonistAge               int     `csv:"AAge"`
	ProtagonistAge              int     `csv:"PAge"`
	AntagonistRun               int     `csv:"ARun"`
	ProtagonistRun              int     `csv:"PRun"`
}

type SimulationBestEpochs []SimulationBestEpoch

func (s *SimulationBestEpochs) ToCSV(outputPath string) error {
	outputFileCSV, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFileCSV.Close()

	writer := gocsv.DefaultCSVWriter(outputFileCSV)
	if writer.Error() != nil {
		return writer.Error()
	}
	err = gocsv.Marshal(s, outputFileCSV)
	if err != nil {
		return err
	}
	return err
}
