package evolution

import (
	"fmt"
	"math"
)

const (
	FitnessAbsolute                        = 0
	FitnessImproverTally                   = 1
	FitnessProtagonistThresholdTally       = 2
	FitnessStrictProtagonistThresholdTally = 3
	FitnessRelativeTally                   = 4
)


//// CalculateFitness
//func CalculateFitness(spec Spec, protagonistAntagonistProgramPair *Program, threshold, minThreshold float64) (antagonistFitness int, protagonistFitness int, err error) {
//	switch i.fitnessCalculationMethod {
//	case FitnessProtagonistThresholdTally:
//		antagonistFitness, protagonistFitness, err := ProtagonistThresholdTally(spec, pro)
//	}
//	return 0
//}

// FitnessAbsolute calculates fitness for two competing individuals.
// The absolute value from the spec is obtained summed, and given to each individual. Smaller values are better.
// 0 being the absolute best.
func AbsoluteFitness(antagonist *Antagonist, protagonist *Protagonist) {}

// ProtagonistThresholdTally takes only the protagonist and checks to see if the Protagonist(Antagonist(
// InitialProgram)) (indicated by protagonistAntagonistProgramPair) lay within the threshold of the spec.
// If not the antagonist receives a fitness of -1 (
// which is better) and the protagonist receives a fitness of (1) and vice versa
func ProtagonistThresholdTally(spec Spec, protagonistAntagonistProgramPair *Program, threshold,
	minThreshold float64) (antagonistFitness int,
	protagonistFitness int, err error) {

	if spec == nil {
		return 0, 0, fmt.Errorf("spec cannot be nil when calculating fitness")
	}
	if protagonistAntagonistProgramPair == nil {
		return 0, 0, fmt.Errorf("protagonistAntagonistProgramPair cannot be nil when calculating fitness")
	}
	if minThreshold <= 0 {
		return 0, 0, fmt.Errorf("minThreshold cannot be less than or equal to 0")
	}
	if threshold <= minThreshold {
		return 0, 0, fmt.Errorf("threshold cannot be less than or equal to minThreshold of %f", minThreshold)
	}

	protagonistDiffSum := 0.0
	for _, equationPairing := range spec {
		dependentVal, err := protagonistAntagonistProgramPair.Eval(equationPairing.Independent)
		if err != nil {
			return 0, 0, err
		}
		abs := math.Abs(float64(dependentVal - equationPairing.Dependent))
		protagonistDiffSum += abs
	}

	avgProtagonist := protagonistDiffSum / float64(len(spec))
	if avgProtagonist <= threshold {
		return 1, -1, nil
	}
	return -1, 1, nil
}
