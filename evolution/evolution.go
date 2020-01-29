package evolution

import (
	"fmt"
	"github.com/gosuri/uiprogress"
	"github.com/martinomburajr/masters-go/evolog"
	"gonum.org/v1/gonum/stat"
	"math"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type EvolutionEngine struct {
	Generations []*Generation   `json:"generations"`
	Parameters  EvolutionParams `json:"parameters"`

	successfulGenerations                       int
	successfulGenerationsByAvg                  int
	minimumTopProtagonistThreshold              int
	minimumMeanProtagonistInGenerationThreshold int

	ProgressBar *uiprogress.Bar
}

// EvaluateTerminationCriteria looks at the current state of the Generation and checks to see if the current
// termination criteria have been achieved. If so it returns true, if not the evolution can move on to the next step
func (engine *EvolutionEngine) EvaluateTerminationCriteria(generation *Generation,
	params EvolutionParams) (shouldTerminateEvolution bool) {
	meanPro := stat.Mean(generation.ProtagonistAvgFitness, nil)
	bestProtagonist := &Individual{}
	bestProtagonist.AverageFitness = -2.0
	for j := range generation.Protagonists {
		currProtagonist := generation.Protagonists[j]
		if currProtagonist.AverageFitness >= bestProtagonist.AverageFitness {
			bestProtagonist = currProtagonist
		}
	}
	if bestProtagonist.AverageFitness >= params.ProtagonistMinGenAvgFit {
		engine.successfulGenerations++
	} else {
		engine.successfulGenerations = 0
	}
	if meanPro >= params.ProtagonistMinGenAvgFit {
		engine.successfulGenerationsByAvg++
	} else {
		engine.successfulGenerationsByAvg = 0
	}
	// If number of successful Generations has been hit, break
	if engine.successfulGenerations >= engine.minimumTopProtagonistThreshold {
		msg := fmt.Sprintf("COMPLETED CYCLE AT GENERATION: %d \n", generation.count)
		params.FinalGeneration = generation.count
		params.LoggingChan <- evolog.Logger{Type: evolog.LoggerGeneration, Message: msg,
			Timestamp: time.Now()}
		params.FinalGenerationReason = "BestIndividual"
		return true
	}
	if engine.successfulGenerations >= engine.minimumMeanProtagonistInGenerationThreshold {
		msg := fmt.Sprintf("COMPLETED CYCLE AT GENERATION: %d \n", generation.count)
		params.FinalGeneration = generation.count
		params.LoggingChan <- evolog.Logger{Type: evolog.LoggerGeneration, Message: msg,
			Timestamp: time.Now()}
		params.FinalGenerationReason = "AvgGeneration"
		return true
	}
	return false
}

func WriteGenerationToLog(e *EvolutionEngine, i int, elapsed time.Duration) {
	numGoroutine := runtime.NumGoroutine()
	msg := fmt.Sprintf("\nFile: %s\t | Spec: %s\t | Run: %d | Gen: (%d/%d) | TSz: %d | numG#: %d | Elapsed: %s",
		e.Parameters.ParamFile,
		e.Parameters.SpecParam.ExpressionParsed,
		e.Parameters.InternalCount,
		i+1,
		e.Parameters.MaxGenerations,
		e.Parameters.Strategies.DepthOfRandomNewTrees,
		numGoroutine,
		elapsed.String())
	e.Parameters.LoggingChan <- evolog.Logger{Type: evolog.LoggerGeneration, Message: msg, Timestamp: time.Now()}
}

func (engine *EvolutionEngine) ValidateGenerationTerminationMinimums() (minimumTopProtagonistThreshold int,
	minimumMeanProtagonistInGenerationThreshold int) {
	minimumTopProtagonistThreshold = int(engine.Parameters.MinimumTopProtagonistMeanBeforeTerminate * float64(engine.Parameters.
		MaxGenerations))
	minimumMeanProtagonistInGenerationThreshold = int(engine.Parameters.MinimumGenerationMeanBeforeTerminate * float64(engine.Parameters.
		MaxGenerations))

	if minimumTopProtagonistThreshold < MinAllowableGenerationsToTerminate {
		engine.Parameters.MinimumTopProtagonistMeanBeforeTerminate = MinAllowableGenerationsToTerminate + 1
		engine.Parameters.LoggingChan <- evolog.Logger{
			Type:      evolog.LoggerGeneration,
			Message:   fmt.Sprintf("NOTE: Set MinimumTopProtagonistMeanBeforeTerminate: %d", engine.Parameters.MinimumTopProtagonistMeanBeforeTerminate),
			Timestamp: time.Now(),
		}
	}
	if minimumMeanProtagonistInGenerationThreshold < MinAllowableGenerationsToTerminate {
		engine.Parameters.MinimumGenerationMeanBeforeTerminate = MinAllowableGenerationsToTerminate + 1
		engine.Parameters.LoggingChan <- evolog.Logger{
			Type:      evolog.LoggerGeneration,
			Message:   fmt.Sprintf("NOTE: Set MinimumGenerationMeanBeforeTerminate: %d", engine.Parameters.MinimumGenerationMeanBeforeTerminate),
			Timestamp: time.Now(),
		}
	}

	return minimumTopProtagonistThreshold, minimumMeanProtagonistInGenerationThreshold
}

// InitializeGenerations starts the first generation as a building block for the evolutionary process.
// It will embedd the antagonists and protagonists created into its Generations slice at index [0]
func (engine *EvolutionEngine) InitializeGenerations(params EvolutionParams) (antagonists []*Individual, protagonists []*Individual, err error) {
	engine.Generations = make([]*Generation, 1)
	engine.Generations[0] = &Generation{}

	antagonists, protagonists, err = engine.Generations[0].InitializePopulation(params)
	if err != nil {
		return nil, nil, err
	}
	engine.Generations[0].GenerationID = GenerateGenerationID(0, params.Topology.Type)
	engine.Generations[0].Antagonists = antagonists
	engine.Generations[0].Protagonists = protagonists
	engine.Generations[0].engine = engine
	engine.Generations[0].AntagonistAvgFitness = make([]float64, 0)
	engine.Generations[0].ProtagonistAvgFitness = make([]float64, 0)

	engine.successfulGenerations = 0
	engine.successfulGenerationsByAvg = 0
	engine.minimumTopProtagonistThreshold, engine.minimumMeanProtagonistInGenerationThreshold = engine.ValidateGenerationTerminationMinimums()

	return antagonists, protagonists, err
}

func (engine *EvolutionEngine) RunGenerationStatistics(currentGeneration *Generation) {
	correlation := stat.Correlation(currentGeneration.AntagonistAvgFitness,
		currentGeneration.ProtagonistAvgFitness, nil)
	covariance := stat.Covariance(currentGeneration.AntagonistAvgFitness,
		currentGeneration.ProtagonistAvgFitness, nil)
	antMean, antStd := stat.MeanStdDev(currentGeneration.AntagonistAvgFitness, nil)
	proMean, proStd := stat.MeanStdDev(currentGeneration.ProtagonistAvgFitness, nil)

	antVar := stat.Variance(currentGeneration.AntagonistAvgFitness, nil)
	antSkew := stat.Skew(currentGeneration.AntagonistAvgFitness, nil)
	antExKurtosis := stat.ExKurtosis(currentGeneration.AntagonistAvgFitness, nil)

	proVar := stat.Variance(currentGeneration.ProtagonistAvgFitness, nil)
	proSkew := stat.Skew(currentGeneration.ProtagonistAvgFitness, nil)
	proExKurtosis := stat.ExKurtosis(currentGeneration.ProtagonistAvgFitness, nil)

	currentGeneration.AntagonistAverage = antMean
	currentGeneration.AntagonistStdDev = antStd
	currentGeneration.AntagonistVariance = antVar
	currentGeneration.ProtagonistAverage = proMean
	currentGeneration.ProtagonistStdDev = proStd
	currentGeneration.ProtagonistVariance = proVar
	currentGeneration.AntagonistExKurtosis = antExKurtosis
	currentGeneration.ProtagonistExKurtosis = proExKurtosis
	currentGeneration.AntagonistSkew = antSkew
	currentGeneration.ProtagonistSkew = proSkew
	currentGeneration.Correlation = correlation
	currentGeneration.Covariance = covariance

	if currentGeneration.BestAntagonist.Id == "" {
		bestAnt := &Individual{AverageFitness: math.MinInt64}
		for _, currAnt := range currentGeneration.Antagonists {
			if currAnt.AverageFitness > bestAnt.AverageFitness {
				bestAnt = currAnt
			}
		}
		bestAntClone, err := bestAnt.Clone()
		if err != nil {
			engine.Parameters.ErrorChan <- err
		}
		currentGeneration.BestAntagonist = bestAntClone
	}

	if currentGeneration.BestProtagonist.Id == "" {
		bestPro := &Individual{AverageFitness: math.MinInt64}
		for _, currPro := range currentGeneration.Protagonists {
			if currPro.AverageFitness > bestPro.AverageFitness {
				bestPro = currPro
			}
		}
		bestProClone, err := bestPro.Clone()
		if err != nil {
			engine.Parameters.ErrorChan <- err
		}
		currentGeneration.BestProtagonist = bestProClone
	}

	statsString := currentGeneration.ToString()
	engine.Parameters.LoggingChan <- evolog.Logger{Timestamp: time.Now(), Type: evolog.LoggerGeneration, Message: statsString}
}

func WriteToDataFolders(folderPercentages []float64, currentGeneration, generationTotalSize int,
	params EvolutionParams) {
	for _, folderPercentage := range folderPercentages {
		if float64(currentGeneration) == math.Floor(float64(generationTotalSize)*folderPercentage) {
			fileName := fmt.Sprintf("%d.txt", int(folderPercentage*100))
			go WriteToDataFolder(params.StatisticsOutput.OutputPath, fileName,
				time.Now().Format(time.RFC3339),
				params.LoggingChan,
				params.ErrorChan)
		}
	}
}

func WriteToDataFolder(dataFolderPath string, fileName string, fileValue string, logChan chan evolog.Logger,
	errChan chan error) {
	mut := sync.Mutex{}
	mut.Lock()

	filepath := fmt.Sprintf("%s/%s", dataFolderPath, fileName)
	os.Mkdir(dataFolderPath, 0775)

	file, err := os.Create(filepath)
	if err != nil {
		errChan <- err
	}

	n, err := fmt.Fprintf(file, "%s", fileValue)
	if err != nil {
		errChan <- err
	} else {
		msg := fmt.Sprintf("25 PERCENT: => Wrote %d bytes to file %s", n, filepath)
		logChan <- evolog.Logger{Type: evolog.LoggerGeneration, Message: msg, Timestamp: time.Now()}
	}

	file.Close()
	mut.Unlock()
}

// Todo Implement EvolutionProcess validate
func (engine *EvolutionEngine) validate() error {
	if engine.Parameters.GenerationsCount < 1 {
		return fmt.Errorf("set number of generationCount by calling e.GenerationsCount(x)")
	}
	if engine.Parameters.EachPopulationSize%4 != 0 {
		return fmt.Errorf("set number of EachPopulationSize to a number that is divisible by 2^x e.g. 8, 16, 32, 64, " +
			"128")
	}
	//if e.Parameters.SetEqualStrategyLength == true && e.Parameters.EqualStrategiesLength < 1 {
	//	return fmt.Errorf("cannot SetEqualStrategyLength to true and EqualStrategiesLength less than 1")
	//}
	if engine.Parameters.StartIndividual.T == nil {
		return fmt.Errorf("start individual cannot have a nil Tree")
	}
	if engine.Parameters.Spec == nil {
		return fmt.Errorf("spec cannot be nil")
	}
	if len(engine.Parameters.Spec) < 1 {
		return fmt.Errorf("spec cannot be empty")
	}
	if engine.Parameters.Selection.Survivor.SurvivorPercentage > 1 || engine.Parameters.Selection.Survivor.
		SurvivorPercentage < 0 {
		return fmt.Errorf("SurvivorPercentage cannot be less than 0 or greater than 1. It is a percent value")
	}
	if engine.Parameters.Selection.Parent.TournamentSize >= engine.Parameters.EachPopulationSize {
		return fmt.Errorf("Tournament Size should not be greater than the population size.")
	}
	//err := e.StartIndividual.Validate()
	//if err != nil {
	//	return err
	//}

	if len(engine.Parameters.Spec) < 3 {
		return fmt.Errorf("a small spec will hamper evolutionary accuracy")
	}
	return nil
}

func TruncShort(s []Strategy) string {
	sb := strings.Builder{}

	for _, str := range s {
		sb.WriteByte(str[0])
	}

	return sb.String()
}
