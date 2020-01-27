package evolution

import (
	"gonum.org/v1/gonum/stat"
	"sync"
)

const (
	ProgressCountersEvolutionResult = 7
)

type EvolutionResult struct {
	HasBeenAnalyzed     bool
	TopAntagonistInRun  Individual
	TopProtagonistInRun Individual
	FinalAntagonist     Individual
	FinalProtagonist    Individual

	Correlation    float64
	Covariance     float64
	CovarianceStd  float64
	CorrelationStd float64
	Generational   Generational

	ThoroughlySortedGenerations []*Generation
	OutputFile                  string

	Mutex sync.Mutex
}

//Generational averages contain slices of length of the generations in a given run
type Generational struct {
	BestAntagonistInEachGenerationByAvgFitness  []Individual
	BestProtagonistInEachGenerationByAvgFitness []Individual

	CorrelationInEachGeneration []float64
	CovarianceInEachGeneration  []float64

	AntagonistAverageInEachGeneration    []float64
	AntagonistStdDevInEachGeneration     []float64
	AntagonistVarianceInEachGeneration   []float64
	AntagonistAvgFitnessInEachGeneration []float64
	AntagonistSkewInEachGeneration       []float64
	AntagonistExKurtosisInEachGeneration []float64

	ProtagonistAverageInEachGeneration    []float64
	ProtagonistStdDevInEachGeneration     []float64
	ProtagonistVarianceInEachGeneration   []float64
	ProtagonistSkewInEachGeneration       []float64
	ProtagonistExKurtosisInEachGeneration []float64
	ProtagonistAvgFitnessInEachGeneration []float64
}

func (e *EvolutionResult) Analyze(evolutionEngine *EvolutionEngine, generations []*Generation, isMoreFitnessBetter bool,
	params EvolutionParams) error {

	genCount := CalculateGenerationSize(params)

	evolutionEngine.ProgressBar.Incr()
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func(generations []*Generation, e *EvolutionResult, wg *sync.WaitGroup) {
		defer wg.Done()

		e.Mutex.Lock()
		sortedFinalAntagonists, err := SortIndividuals(generations[len(generations)-1].Antagonists, true)
		if err != nil {
			params.ErrorChan <- err
		}
		e.FinalAntagonist, err = sortedFinalAntagonists[0].Clone()
		if err != nil {
			params.ErrorChan <- err
		}
		e.Mutex.Unlock()
		evolutionEngine.ProgressBar.Incr()
	}(generations, e, &wg)

	go func(generations []*Generation, e *EvolutionResult, wg *sync.WaitGroup) {
		defer wg.Done()
		e.Mutex.Lock()
		sortedFinalProtagonists, err := SortIndividuals(generations[len(generations)-1].Protagonists, true)
		if err != nil {
			params.ErrorChan <- err
		}

		e.FinalProtagonist, err = sortedFinalProtagonists[0].Clone()
		if err != nil {
			params.ErrorChan <- err
		}
		e.Mutex.Unlock()
		evolutionEngine.ProgressBar.Incr()
	}(generations, e, &wg)

	go func(generations []*Generation, e *EvolutionResult, wg *sync.WaitGroup) {
		defer wg.Done()
		e.Mutex.Lock()
		sortedGenerations, err := SortGenerationsThoroughly(generations, isMoreFitnessBetter)
		if err != nil {
			params.ErrorChan <- err
		}
		e.ThoroughlySortedGenerations = sortedGenerations
		e.Mutex.Unlock()
		evolutionEngine.ProgressBar.Incr()
	}(generations, e, &wg)

	wg.Wait()

	for i := 0; i < genCount; i++ {
		correlations := make([]float64, genCount)
		covariances := make([]float64, genCount)

		correlations[i] = generations[i].Correlation
		covariances[i] = generations[i].Covariance

		corrMean, corrStd := stat.MeanStdDev(correlations, nil)
		covMean, covStd := stat.MeanStdDev(covariances, nil)

		e.Covariance = covMean
		e.CovarianceStd = covStd
		e.Correlation = corrMean
		e.CorrelationStd = corrStd
	}

	// Calculate Top Individuals
	topAntagonist, topProtagonist, err := GetTopIndividualInRun(e.ThoroughlySortedGenerations, isMoreFitnessBetter)
	if err != nil {
		return err
	}
	evolutionEngine.ProgressBar.Incr()
	e.TopAntagonistInRun, err = topAntagonist.Clone()
	if err != nil {
		return err
	}
	e.TopProtagonistInRun, err = topProtagonist.Clone()
	if err != nil {
		return err
	}

	e.Generational.BestAntagonistInEachGenerationByAvgFitness = make([]Individual, genCount)
	e.Generational.BestProtagonistInEachGenerationByAvgFitness = make([]Individual, genCount)
	e.Generational.AntagonistAverageInEachGeneration = make([]float64, genCount)
	e.Generational.AntagonistStdDevInEachGeneration = make([]float64, genCount)
	e.Generational.AntagonistVarianceInEachGeneration = make([]float64, genCount)
	e.Generational.AntagonistSkewInEachGeneration = make([]float64, genCount)
	e.Generational.AntagonistExKurtosisInEachGeneration = make([]float64, genCount)
	e.Generational.ProtagonistAverageInEachGeneration = make([]float64, genCount)
	e.Generational.ProtagonistStdDevInEachGeneration = make([]float64, genCount)
	e.Generational.ProtagonistVarianceInEachGeneration = make([]float64, genCount)
	e.Generational.ProtagonistSkewInEachGeneration = make([]float64, genCount)
	e.Generational.ProtagonistExKurtosisInEachGeneration = make([]float64, genCount)
	e.Generational.CorrelationInEachGeneration = make([]float64, genCount)
	e.Generational.CovarianceInEachGeneration = make([]float64, genCount)
	evolutionEngine.ProgressBar.Incr()

	for i := 0; i < genCount; i++ {
		e.Generational.BestAntagonistInEachGenerationByAvgFitness[i] = evolutionEngine.Generations[i].BestAntagonist
		e.Generational.BestProtagonistInEachGenerationByAvgFitness[i] = evolutionEngine.Generations[i].BestProtagonist
		e.Generational.AntagonistAverageInEachGeneration[i] = evolutionEngine.Generations[i].AntagonistAverage
		e.Generational.AntagonistStdDevInEachGeneration[i] = evolutionEngine.Generations[i].AntagonistStdDev
		e.Generational.AntagonistVarianceInEachGeneration[i] = evolutionEngine.Generations[i].AntagonistVariance
		e.Generational.AntagonistSkewInEachGeneration[i] = evolutionEngine.Generations[i].AntagonistSkew
		e.Generational.AntagonistExKurtosisInEachGeneration[i] = evolutionEngine.Generations[i].AntagonistExKurtosis
		e.Generational.ProtagonistAverageInEachGeneration[i] = evolutionEngine.Generations[i].ProtagonistAverage
		e.Generational.ProtagonistStdDevInEachGeneration[i] = evolutionEngine.Generations[i].ProtagonistStdDev
		e.Generational.ProtagonistVarianceInEachGeneration[i] = evolutionEngine.Generations[i].ProtagonistVariance
		e.Generational.ProtagonistSkewInEachGeneration[i] = evolutionEngine.Generations[i].ProtagonistSkew
		e.Generational.ProtagonistExKurtosisInEachGeneration[i] = evolutionEngine.Generations[i].ProtagonistExKurtosis
		e.Generational.CorrelationInEachGeneration[i] = evolutionEngine.Generations[i].Correlation
		e.Generational.CovarianceInEachGeneration[i] = evolutionEngine.Generations[i].Covariance
	}
	e.HasBeenAnalyzed = true
	evolutionEngine.ProgressBar.Incr()
	return err
}

func CalculateGenerationSize(params EvolutionParams) int {
	genCount := 0
	if params.MaxGenerations > MinAllowableGenerationsToTerminate {
		if params.FinalGeneration > 0 {
			genCount = params.FinalGeneration
		} else {
			genCount = params.MaxGenerations
		}
	} else {
		genCount = params.GenerationsCount
	}
	return genCount
}

func (e *EvolutionResult) Clean() {
	e.ThoroughlySortedGenerations = nil
}
