package evolution

import (
	"sync"
)


const (
	ProgressCountersEvolutionResult= 10
)

type EvolutionResult struct {
	HasBeenAnalyzed        bool
	TopAntagonist          *Individual
	TopProtagonist         *Individual
	TopAntagonistDelta     *Individual
	TopProtagonistDelta    *Individual
	TopAntagonistAvgDelta  *Individual
	TopProtagonistAvgDelta *Individual
	IsMoreFitnessBetter    bool
	FinalAntagonist        *Individual
	FinalProtagonist       *Individual

	CoevolutionaryAverages []GenerationalCoevolutionaryAverages
	Generational           Generational

	SortedGenerationIndividuals           []*Generation
	SortedGenerationIndividualsByDelta    []*Generation
	SortedGenerationIndividualsByDeltaAvg []*Generation
	OutputFile                            string

	Mutex sync.Mutex
}

type multiIndividualsPerGeneration struct {
	Generation  *Generation
	Individuals []*Individual
}

type GenerationalCoevolutionaryAverages struct {
	Generation                     *Generation
	AntagonistFitnessAverages      float64
	ProtagonistFitnessAverages     float64
	AntagonistBestFitnessAverages  float64
	ProtagonistBestFitnessAverages float64
	AntagonistDeltaAverages        float64
	ProtagonistDeltaAverages       float64
	AntagonistBestDeltaAverages    float64
	ProtagonistBestDeltaAverages   float64
}

//Generational averages contain slices of length of the generations in a given run
type Generational struct {
	Antagonists                    []Individual
	Protagonists                   []Individual
	AntagonistsByDelta             []Individual
	ProtagonistsByDelta            []Individual
	AntagonistsByDeltaAvg          []Individual
	ProtagonistsByDeltaAvg         []Individual
	AntagonistFitnessAverages      []float64
	ProtagonistFitnessAverages     []float64
	AntagonistBestFitnessAverages  []float64
	ProtagonistBestFitnessAverages []float64
	AntagonistDeltaAverages        []float64
	ProtagonistDeltaAverages       []float64
	AntagonistBestDeltaAverages    []float64
	ProtagonistBestDeltaAverages   []float64
}

func (e *EvolutionResult) Analyze(evolutionEngine *EvolutionEngine, generations []*Generation, isMoreFitnessBetter bool,
	params EvolutionParams) error {

	evolutionEngine.ProgressBar.Incr()
	wg := sync.WaitGroup{}
	wg.Add(5)
	go func(generations []*Generation, e *EvolutionResult, wg *sync.WaitGroup) {
		defer wg.Done()

		e.Mutex.Lock()
		sortedFinalAntagonists, err := SortIndividuals(generations[len(generations)-1].Antagonists, true)
		if err != nil {
			params.ErrorChan <- err
		}
		e.FinalAntagonist = sortedFinalAntagonists[0]
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

		e.FinalProtagonist = sortedFinalProtagonists[0]
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
		e.SortedGenerationIndividuals = sortedGenerations
		e.Mutex.Unlock()
		evolutionEngine.ProgressBar.Incr()
	}(generations, e, &wg)

	go func(generations []*Generation, e *EvolutionResult, wg *sync.WaitGroup) {
		defer wg.Done()
		e.Mutex.Lock()
		sortedGenerationsByDelta, err := SortGenerationsThoroughlyByDelta(generations, true, false)
		if err != nil {
			params.ErrorChan <- err
		}
		e.SortedGenerationIndividualsByDelta = sortedGenerationsByDelta
		e.Mutex.Unlock()
		evolutionEngine.ProgressBar.Incr()
	}(generations, e, &wg)

	go func(generations []*Generation, e *EvolutionResult, wg *sync.WaitGroup) {
		defer wg.Done()
		e.Mutex.Lock()
		sortedGenerationsByAvgDelta, err := SortGenerationsThoroughlyByAvgDelta(generations, true, false)
		if err != nil {
			params.ErrorChan <- err
		}
		e.SortedGenerationIndividualsByDeltaAvg = sortedGenerationsByAvgDelta
		e.Mutex.Unlock()
		evolutionEngine.ProgressBar.Incr()
	}(generations, e, &wg)
	wg.Wait()



	if params.MaxGenerations > MinAllowableGenerationsForContinuous {
		e.Generational.Antagonists = make([]Individual, params.MaxGenerations)
		e.Generational.AntagonistsByDeltaAvg = make([]Individual, params.MaxGenerations)
		e.Generational.AntagonistsByDelta = make([]Individual, params.MaxGenerations)

		e.Generational.Protagonists = make([]Individual, params.MaxGenerations)
		e.Generational.ProtagonistsByDelta = make([]Individual, params.MaxGenerations)
		e.Generational.ProtagonistsByDeltaAvg = make([]Individual, params.MaxGenerations)
	} else {
		e.Generational.Antagonists = make([]Individual, params.GenerationsCount)
		e.Generational.AntagonistsByDeltaAvg = make([]Individual, params.GenerationsCount)
		e.Generational.AntagonistsByDelta = make([]Individual, params.GenerationsCount)

		e.Generational.Protagonists = make([]Individual, params.GenerationsCount)
		e.Generational.ProtagonistsByDelta = make([]Individual, params.GenerationsCount)
		e.Generational.ProtagonistsByDeltaAvg = make([]Individual, params.GenerationsCount)
	}


	for i, v := range e.SortedGenerationIndividuals {
		antagonist := v.Antagonists[0]
		antagonistByDelta := e.SortedGenerationIndividualsByDelta[i].Antagonists[0]
		antagonistByAvgDelta := e.SortedGenerationIndividualsByDeltaAvg[i].Antagonists[0]
		antClone, _ := antagonist.Clone()
		antagonistByDeltaClone, _ := antagonistByDelta.Clone()
		antagonistByAvgDeltaClone, _ := antagonistByAvgDelta.Clone()

		e.Generational.Antagonists[i] = antClone
		e.Generational.AntagonistsByDelta[i] = antagonistByAvgDeltaClone
		e.Generational.AntagonistsByDeltaAvg[i] = antagonistByDeltaClone
	}
	evolutionEngine.ProgressBar.Incr()


	for i, v := range e.SortedGenerationIndividuals {
		protagonist := v.Protagonists[0]
		rotagonistByDelta := e.SortedGenerationIndividualsByDelta[i].Protagonists[0]
		rotagonistByAvgDelta := e.SortedGenerationIndividualsByDeltaAvg[i].Protagonists[0]
		protagonistClone, _ := protagonist.Clone()
		protagonistByDeltaClone, _ := rotagonistByDelta.Clone()
		protagonistByAvgDeltaClone, _ := rotagonistByAvgDelta.Clone()

		e.Generational.ProtagonistsByDelta[i] = protagonistByDeltaClone
		e.Generational.ProtagonistsByDeltaAvg[i] = protagonistByAvgDeltaClone
		e.Generational.Protagonists[i] = protagonistClone
	}
	evolutionEngine.ProgressBar.Incr()

	// Calculate Top Individuals
	topAntagonist, topProtagonist, err := GetTopIndividualInAllGenerations(e.SortedGenerationIndividuals, isMoreFitnessBetter)
	if err != nil {
		return err
	}
	e.TopAntagonist = topAntagonist
	e.TopProtagonist = topProtagonist

	// Calculate GenerationalStatistics Averages
	coevolutionaryAverages, err := GetGenerationalAverages(e.SortedGenerationIndividuals)
	if err != nil {
		return err
	}
	e.CoevolutionaryAverages = coevolutionaryAverages

	if params.MaxGenerations > MinAllowableGenerationsForContinuous {
		e.Generational.AntagonistFitnessAverages = make([]float64, params.MaxGenerations)
		e.Generational.ProtagonistFitnessAverages = make([]float64, params.MaxGenerations)
		e.Generational.AntagonistBestFitnessAverages = make([]float64, params.MaxGenerations)
		e.Generational.ProtagonistBestFitnessAverages = make([]float64, params.MaxGenerations)
		e.Generational.AntagonistBestDeltaAverages = make([]float64, params.MaxGenerations)
		e.Generational.ProtagonistBestDeltaAverages = make([]float64, params.MaxGenerations)
		e.Generational.AntagonistDeltaAverages = make([]float64, params.MaxGenerations)
		e.Generational.ProtagonistDeltaAverages = make([]float64, params.MaxGenerations)
	} else {
		e.Generational.AntagonistFitnessAverages = make([]float64, params.GenerationsCount)
		e.Generational.ProtagonistFitnessAverages = make([]float64, params.GenerationsCount)
		e.Generational.AntagonistBestFitnessAverages = make([]float64, params.GenerationsCount)
		e.Generational.ProtagonistBestFitnessAverages = make([]float64, params.GenerationsCount)
		e.Generational.AntagonistBestDeltaAverages = make([]float64, params.GenerationsCount)
		e.Generational.ProtagonistBestDeltaAverages = make([]float64, params.GenerationsCount)
		e.Generational.AntagonistDeltaAverages = make([]float64, params.GenerationsCount)
		e.Generational.ProtagonistDeltaAverages = make([]float64, params.GenerationsCount)
	}

	evolutionEngine.ProgressBar.Incr()

	for i, v := range coevolutionaryAverages {
		e.Generational.AntagonistFitnessAverages[i] = v.AntagonistFitnessAverages
		e.Generational.ProtagonistFitnessAverages[i] = v.ProtagonistFitnessAverages
		e.Generational.AntagonistBestFitnessAverages[i] = v.AntagonistBestFitnessAverages
		e.Generational.ProtagonistBestFitnessAverages[i] = v.ProtagonistBestFitnessAverages
		e.Generational.AntagonistBestDeltaAverages[i] = v.AntagonistBestDeltaAverages
		e.Generational.ProtagonistBestDeltaAverages[i] = v.ProtagonistBestDeltaAverages
		e.Generational.AntagonistDeltaAverages[i] = v.AntagonistDeltaAverages
		e.Generational.ProtagonistDeltaAverages[i] = v.ProtagonistDeltaAverages
		//e.Generational.Antagonists[i] = v.
	}
	e.HasBeenAnalyzed = true
	evolutionEngine.ProgressBar.Incr()
	return err
}

func (e *EvolutionResult) Clean() {
	e.Generational.Antagonists = nil
	e.Generational.Protagonists = nil
	e.Generational.ProtagonistFitnessAverages = nil
	e.Generational.AntagonistFitnessAverages = nil
	e.FinalAntagonist = nil
	e.FinalProtagonist = nil
	e.TopProtagonist = nil
	e.TopAntagonist = nil
	e.CoevolutionaryAverages = nil
	e.SortedGenerationIndividuals = nil
}


