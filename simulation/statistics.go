package simulation

import (
	"fmt"
	"github.com/martinomburajr/masters-go/evolution"
	"math"
)

// BestIndividualsAllRuns returns the best antagonist and protagonist in all runs
func (s *Simulation) BestIndividualsAllRuns(params evolution.EvolutionParams) (topAntagonist,
	topProtagonist evolution.Individual, sim SimulationBestIndividualsAllTime, err error) {
	if s.SimulationStats == nil {
		return evolution.Individual{}, evolution.Individual{}, SimulationBestIndividualsAllTime{}, fmt.Errorf(
			"ToRunStats | simulationStats is nil")
	}
	if len(s.SimulationStats) < 0 {
		return evolution.Individual{},evolution.Individual{},SimulationBestIndividualsAllTime{}, fmt.Errorf("ToRunStats | simulationStats is empty")
	}

	topAntagonist = evolution.Individual{AverageFitness: math.MinInt64}
	topProtagonist = evolution.Individual{AverageFitness: math.MinInt64}
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

	sim = SimulationBestIndividualsAllTime{
		SpecEquation:                params.SpecParam.Expression,
		SpecRange:                   params.SpecParam.Range,
		SpecSeed:                    params.SpecParam.Seed,
		AntagonistEquation:          topAntagonistEq,
		ProtagonistEquation:         topProtagonistEq,
		AntagonistDelta:             topAntagonist.BestFitnessDelta,
		ProtagonistDelta:            topProtagonist.BestFitnessDelta,
		AntagonistGeneration:        topAntGen,
		ProtagonistGeneration:       topProGen,
		AntagonistRun:               topAntRun,
		ProtagonistRun:              topProRun,
		AntagonistAge:               topAntagonist.Age,
		ProtagonistAge:              topProtagonist.Age,
		AntagonistBirthGen:          topAntagonist.BirthGen,
		ProtagonistBirthGen:         topAntagonist.BirthGen,
		AntagonistDominantStrategy:  evolution.DominantStrategy(topAntagonist),
		ProtagonistDominantStrategy: evolution.DominantStrategy(topProtagonist),
		AntagonistStrategyList:      evolution.StrategiesToString(topAntagonist),
		ProtagonistStrategyList:     evolution.StrategiesToString(topProtagonist),
	}

	return topAntagonist, topProtagonist, sim, err
}

// GenerationalAveragesAllRuns returns the best antagonist and protagonist in all runs
func (s *Simulation) GenerationalAveragesAllRuns(params evolution.EvolutionParams) (generationalAverages []GenerationalAverages,
	err error) {
	if s.SimulationStats == nil {
		return nil, fmt.Errorf("ToRunStats | simulationStats is nil")
	}
	if len(s.SimulationStats) < 0 {
		return nil, fmt.Errorf("ToRunStats | simulationStats is empty")
	}

	length := len(s.SimulationStats[0].GenerationalAverages)
	generationalAverages = make([]GenerationalAverages, length)
	for i := range generationalAverages {
		for j := 0; j < length; j++ {
			sumAnt := 0.0
			sumProt := 0.0
			for a := 0; a < len(s.SimulationStats); a++ {
				runStats := s.SimulationStats[a]
				sumAnt += runStats.GenerationalAverages[i].AntagonistResult
				sumProt += runStats.GenerationalAverages[i].ProtagonistResult
			}
			generationalAverages[i].AntagonistResult = sumAnt/float64(len(s.SimulationStats))
			generationalAverages[i].ProtagonistResult = sumProt/float64(len(s.SimulationStats))
		}
		generationalAverages[i].Generation = i
	}

	return generationalAverages, err
}


func (s *Simulation) ToRunStats() ([]RunBasedStatistics, error) {
	if s.SimulationStats == nil {
		return nil, fmt.Errorf("ToRunStats | simulationStats is nil")
	}
	if len(s.SimulationStats) < 0 {
		return nil, fmt.Errorf("ToRunStats | simulationStats is empty")
	}
	runStats := make([]RunBasedStatistics, s.NumberOfRunsPerState)
	for i, simulationStat := range s.SimulationStats {
		topAntEq, _ := simulationStat.TopAntagonist.Program.T.ToMathematicalString()
		topProEq, _ := simulationStat.TopProtagonist.Program.T.ToMathematicalString()
		finAntEq, _ := simulationStat.FinalAntagonist.Program.T.ToMathematicalString()
		finProEq, _ := simulationStat.FinalProtagonist.Program.T.ToMathematicalString()
		runStats[i] = RunBasedStatistics{
			TopAntagonist:            simulationStat.TopAntagonist.AverageFitness,
			TopProtagonist:           simulationStat.TopProtagonist.AverageFitness,
			TopAntagonistDelta:       simulationStat.TopAntagonist.BestFitnessDelta,
			TopProtagonistDelta:      simulationStat.TopProtagonist.BestFitnessDelta,
			TopAntagonistStrategy:    evolution.StrategiesToString(simulationStat.TopAntagonist),
			TopProtagonistStrategy:   evolution.StrategiesToString(simulationStat.TopProtagonist),
			TopAntagonistEquation:    topAntEq,
			TopProtagonistEquation:   topProEq,
			FinalAntagonist:          simulationStat.FinalAntagonist.AverageFitness,
			FinalProtagonist:         simulationStat.FinalProtagonist.AverageFitness,
			FinalAntagonistDelta:     simulationStat.FinalAntagonist.BestFitnessDelta,
			FinalProtagonistDelta:    simulationStat.FinalProtagonist.BestFitnessDelta,
			FinalAntagonistStrategy:  evolution.StrategiesToString(simulationStat.FinalAntagonist),
			FinalProtagonistStrategy: evolution.StrategiesToString(simulationStat.FinalAntagonist),
			FinalAntagonistEquation:  finAntEq,
			FinalProtagonistEquation: finProEq,
			Run:                      i,
		}
	}

	return runStats, nil
}

//
//func (s *Simulation) ToFinalStats(runStats []RunBasedStatistics) (error) {
//	if s.SimulationStats == nil {
//		return fmt.Errorf("ToRunStats | simulationStats is nil")
//	}
//	if len(s.SimulationStats) < 0 {
//		return fmt.Errorf("ToRunStats | simulationStats is empty")
//	}
//
//	// best
//	type BestIndividuals struct {
//		SpecEquation                string  `csv:"specEquation"`
//		SpecRange                   int     `csv:"range"`
//		SpecSeed                    int     `csv:"seed"`
//		AntagonistEquation          string  `csv:"A"`
//		ProtagonistEquation         string  `csv:"P"`
//		AntagonistDelta             float64 `csv:"ADelta"`
//		ProtagonistDelta            float64 `csv:"PDelta"`
//		AntagonistGeneration        int     `csv:"AGeneration"`
//		ProtagonistGeneration       int     `csv:"PGeneration"`
//		AntagonistRun               int     `csv:"ARun"`
//		ProtagonistRun              int     `csv:"PRun"`
//		AntagonistBirthGen          int     `csv:"ABirthGen"`
//		ProtagonistBirthGen         int     `csv:"PBirthGen"`
//		AntagonistDominantStrategy  string  `csv:"AFaveStrategy"`
//		ProtagonistDominantStrategy string  `csv:"PFaveStrategy"`
//		AntagonistStrategyList      string  `csv:"AStrategies"`
//		ProtagonistStrategyList     string  `csv:"PStrategies"`
//	}
//
//	for _, run := range runStats {
//		bestAnt := 0
//		if run.TopAntagonist > bestAnt {
//
//		}
//		bestProt := 0
//	}
//
//	//average
//
//	return runStats, nil
//}

// ToStrategyStats
func (s *Simulation) ToStrategyStats(dirPath string) (statistics []RunStrategyStatistics, err error) {
	if s.SimulationStats == nil {
		return  nil, fmt.Errorf("ToRunStats | simulationStats is nil")
	}
	if len(s.SimulationStats) < 0 {
		return  nil, fmt.Errorf("ToRunStats | simulationStats is empty")
	}
	statistics = make([]RunStrategyStatistics, len(s.SimulationStats[0].TopAntagonist.Strategy))

	for i, run := range s.SimulationStats {
		for j := range s.SimulationStats[0].TopAntagonist.Strategy {
			statistics[j] = RunStrategyStatistics {
				TopAntagonistStrategy:    string(run.TopAntagonist.Strategy[j]),
				TopProtagonistStrategy:   string(run.TopProtagonist.Strategy[j]),
				FinalAntagonistStrategy:  string(run.FinalAntagonist.Strategy[j]),
				FinalProtagonistStrategy: string(run.FinalProtagonist.Strategy[j]),
				StrategyNumber:           j+1,
			}
		}
		WriteRunStrategy(statistics, fmt.Sprintf("%s/%d/strategy-%d.csv", dirPath, i, i ))
	}

	return statistics, nil
}

type RunStrategyStatistics struct {
	StrategyNumber                       int     `csv:"stratNum"`
	TopAntagonistStrategy  string  `csv:"runTopAStrategy"`
	TopProtagonistStrategy string  `csv:"runTopPStrategy"`
	FinalAntagonistStrategy  string  `csv:"runFinAStrategy"`
	FinalProtagonistStrategy string  `csv:"runFinPStrategy"`
}

type RunBasedStatistics struct {
	TopAntagonist          float64 `csv:"runTopA"`
	TopProtagonist         float64 `csv:"runTopP"`
	TopAntagonistDelta     float64 `csv:"runTopADelta"`
	TopProtagonistDelta    float64 `csv:"runTopPDelta"`
	TopAntagonistStrategy  string  `csv:"runTopAStrategy"`
	TopProtagonistStrategy string  `csv:"runTopPStrategy"`
	TopAntagonistEquation     string  `csv:"runTopAEquation"`
	TopProtagonistEquation    string  `csv:"runTopPEquation"`

	FinalAntagonist          float64 `csv:"runFinalA"`
	FinalProtagonist         float64 `csv:"runFinalP"`
	FinalAntagonistDelta     float64 `csv:"runFinalADelta"`
	FinalProtagonistDelta    float64 `csv:"runFinalPDelta"`
	FinalAntagonistStrategy  string  `csv:"runFinalAStrategy"`
	FinalProtagonistStrategy string  `csv:"runFinalPStrategy"`
	FinalAntagonistEquation     string  `csv:"runFinalAEquation"`
	FinalProtagonistEquation    string  `csv:"runFinalPEquation"`

	Run                       int     `csv:"run"`
}


// GenerationalAveragesAllRuns returns the best antagonist and protagonist in all runs
//func (s *Simulation) AllGenerations(params evolution.EvolutionParams) (
//	generationalAverages []GenerationalAverages,
//	err error) {
//	if s.SimulationStats == nil {
//		return nil, fmt.Errorf("ToRunStats | simulationStats is nil")
//	}
//	if len(s.SimulationStats) < 0 {
//		return nil, fmt.Errorf("ToRunStats | simulationStats is empty")
//	}
//
//	length := len(s.SimulationStats[0].GenerationalAverages)
//	generationalAverages = make([]GenerationalAverages, length)
//	for i := range generationalAverages {
//		for j := 0; j < length; j++ {
//			sumAnt := 0.0
//			sumProt := 0.0
//			for a := 0; a < len(s.SimulationStats); a++ {
//				runStats := s.SimulationStats[a]
//				sumAnt += runStats.GenerationalAverages[i].AntagonistResult
//				sumProt += runStats.GenerationalAverages[i].ProtagonistResult
//			}
//			generationalAverages[i].AntagonistResult = sumAnt/float64(len(s.SimulationStats))
//			generationalAverages[i].ProtagonistResult = sumProt/float64(len(s.SimulationStats))
//		}
//		generationalAverages[i].Generation = i
//	}
//
//	return generationalAverages, err
//}

type GenerationalAverages struct {
	AntagonistResult  float64 `csv:"A"`
	ProtagonistResult float64 `csv:"P"`
	Generation int `csv:"gen"`
}