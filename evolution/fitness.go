package evolution

import (
	"fmt"
	"math"
)

const (
	FitnessAbsolute                        = 0
	FitnessImproverTally                   = 1
	FitnessProtagonistThresholdTally       = 2
	FitnessRatio                           = 3
	FitnessRatioThresholder                = 4
	FitnessStrictProtagonistThresholdTally = 5
	FitnessRelativeTally                   = 6
)

//// CalculateFitness
//func CalculateFitness(spec SpecMulti, protagonistAntagonistProgramPair *Program, threshold, minThreshold float64) (antagonistFitness int, protagonistFitness int, err error) {
//	switch i.FitnessCalculationMethod {
//	case FitnessProtagonistThresholdTally:
//		antagonistFitness, protagonistFitness, err := ProtagonistThresholdTally(spec, pro)
//	}
//	return 0
//}

// FitnessAbsolute calculates Fitness for two competing individuals.
// The absolute value from the spec is obtained summed, and given to each individual. Smaller values are better.
// 0 being the absolute best.
func AbsoluteFitness(spec SpecMulti, protagonistExpression, antagonistExpression string,
	protagonistProgram, antagonistProgram *Program) (protagonistFitness,
	antagonistFitness float64, err error) {
	deltaProtagonist := 0.0
	deltaAntagonist := 0.0
	for _, s := range spec {
		dependentProtagonistVar, err := protagonistProgram.EvalMulti(s.Independents, protagonistExpression)
		if err != nil {
			return math.MaxFloat32, math.MaxFloat32, err
		}
		dP := calculateDelta(float64(s.Dependent), float64(dependentProtagonistVar))
		deltaProtagonist += dP

		dependentAntagonistVar, err := antagonistProgram.EvalMulti(s.Independents, antagonistExpression)
		if err != nil {
			return math.MaxFloat32, math.MaxFloat32, err
		}
		dA := calculateDelta(float64(s.Dependent), float64(dependentAntagonistVar))
		deltaAntagonist += dA
	}

	return deltaProtagonist, deltaAntagonist, nil
}

// ProtagonistThresholdTally takes only the protagonist and checks to see if the Protagonist(Antagonist(
// InitialProgram)) (indicated by protagonistAntagonistProgramPair) lay within the threshold of the spec.
// If not the antagonist receives a Fitness of -1 (
// which is better) and the protagonist receives a Fitness of (1) and vice versa
func ProtagonistThresholdTally(spec SpecMulti, protagonistAntagonistProgramPair *Program,
	threshold float64) (antagonistFitness float64,
	protagonistFitness float64, err error) {

	if spec == nil {
		return math.MaxInt64, math.MaxInt64, fmt.Errorf("spec cannot be nil when calculating Fitness")
	}
	if protagonistAntagonistProgramPair == nil {
		return math.MaxInt64, math.MaxInt64, fmt.Errorf("protagonistAntagonistProgramPair cannot be nil when calculating Fitness")
	}
	if threshold <= 0 {
		return math.MaxInt64, math.MaxInt64, fmt.Errorf("minThreshold cannot be less than or equal to 0")
	}

	err = protagonistAntagonistProgramPair.T.Validate()
	if err != nil {
		return math.MaxInt64, math.MaxInt64, err
	}

	_, err = protagonistAntagonistProgramPair.T.ToMathematicalString()
	if err != nil {
		return math.MaxInt64, math.MaxInt64, err
	}

	protagonistDiffSum := 0.0
	//for _, equationPairing := range spec {
	//	dependentVal, err := protagonistAntagonistProgramPair.EvalMulti(equationPairing.Independents, expressionString)
	//	if err != nil {
	//		return math.MaxInt64, math.MaxInt64, err
	//	}
	//	abs := math.Abs(float64(dependentVal - equationPairing.Dependent))
	//	protagonistDiffSum += abs
	//}

	avgProtagonist := protagonistDiffSum / float64(len(spec))
	if avgProtagonist <= threshold {
		return 1, -1, nil
	}
	return -1, 1, nil
}

// RatioFitnessThresholded is a strategy where both antagonist and protagonist calculate their fitness given the spec.
// The SpecMulti needs to be broad enough to encompass a wide array of values e.g. > 50.
// The spec defines pure values represented by the input program.
// A threshold is defined that is used as a marker to judge the protagonists fitness.
// The spec also defines an outer threshold that if exceeded all fitness is attributed to the antagonists.
// This threshold is used to ensure the tests pass a given threshold to be deemed worthy.
// If the tests fail to cross the threshold (inwards i.
// e towards the pure values defined by the spec) then they cannot begin to accumulate fitness.
func RatioFitnessThresholded(spec SpecMulti, thresholdMultiplier float64, antagonist, protagonist *Program) (antagonistFitness float64,
	protagonistFitness float64, err error) {

	if spec == nil {
		return math.MaxInt64, math.MaxInt64, fmt.Errorf("RatioFitness | spec cannot be nil when calculating Fitness")
	}
	if thresholdMultiplier < 1 {
		return math.MaxFloat32, math.MaxFloat32, err
	}
	if antagonist == nil {
		return math.MaxInt64, math.MaxInt64,
			fmt.Errorf("RatioFitness | antagonist cannot be nil when calculating Fitness")
	}
	if antagonist.T == nil {
		return math.MaxInt64, math.MaxInt64,
			fmt.Errorf("RatioFitness | antagonist tree cannot be nil when calculating Fitness")
	}
	if antagonist.T.root == nil {
		return math.MaxInt64, math.MaxInt64,
			fmt.Errorf("RatioFitness | antagonist tree root cannot be nil when calculating Fitness")
	}
	if protagonist == nil {
		return math.MaxInt64, math.MaxInt64,
			fmt.Errorf("RatioFitness | protagonist cannot be nil when calculating Fitness")
	}
	if protagonist.T == nil {
		return math.MaxInt64, math.MaxInt64,
			fmt.Errorf("RatioFitness | protagonist tree cannot be nil when calculating Fitness")
	}
	if protagonist.T.root == nil {
		return math.MaxInt64, math.MaxInt64,
			fmt.Errorf("RatioFitness | protagonist tree root cannot be nil when calculating Fitness")
	}

	err = protagonist.T.Validate()
	if err != nil {
		return math.MaxInt64, math.MaxInt64, err
	}

	err = antagonist.T.Validate()
	if err != nil {
		return math.MaxInt64, math.MaxInt64, err
	}
	// Antagonist
	antagonistMathematicalExpression, err := antagonist.T.ToMathematicalString()
	if err != nil {
		return math.MaxInt64, math.MaxInt64, err
	}

	protagonistMathematicalExpression, err := protagonist.T.ToMathematicalString()
	if err != nil {
		return math.MaxInt64, math.MaxInt64, err
	}

	antagonistFitness, protagonistFitness, err = ratioFitnessThreshold(spec, thresholdMultiplier,
		antagonistMathematicalExpression,
		protagonistMathematicalExpression, antagonist, protagonist)
	if err != nil {
		return math.MaxInt64, math.MaxInt64, err
	}

	return antagonistFitness, protagonistFitness, nil
}

// ratioFitness assusmes the input variables have been checked for validity and nilness.
// If the protagonists delta is greater than that of the antagonist it automatically is given a value of 0.
// If the protagonist achieves identical spec values it obtains a value of 1 and the antagonists gets a value of 0
func ratioFitnessThreshold(spec SpecMulti, thresholdMultiplier float64, antagonistExpression,
	protagonistExpression string,
	antagonistProgram, protagonistProgram *Program) (antagonistFitness, protagonistFitness float64, err error) {
	deltaProtagonist, deltaAntagonist, thresholderSum := 0.0, 0.0, 0.0

	for _, s := range spec {
		thresholderSum += s.Dependent

		dependentProtagonistVar, err := protagonistProgram.EvalMulti(s.Independents, protagonistExpression)
		if err != nil {
			return math.MaxFloat32, math.MaxFloat32, err
		}
		dP := calculateDelta(float64(s.Dependent), float64(dependentProtagonistVar))
		deltaProtagonist += dP

		dependentAntagonistVar, err := antagonistProgram.EvalMulti(s.Independents, antagonistExpression)
		if err != nil {
			return math.MaxFloat32, math.MaxFloat32, err
		}
		dA := calculateDelta(float64(s.Dependent), float64(dependentAntagonistVar))
		deltaAntagonist += dA
	}
	threshold := thresholdMultiplier * thresholderSum

	if deltaProtagonist >= deltaAntagonist {
		return 0, 1, nil
	}
	if deltaProtagonist == 0 {
		return 1, 0, nil
	}

	var progFitness, antFitness float64
	if deltaProtagonist <= threshold {
		if deltaAntagonist <= threshold {
			antFitness = 0
			progFitness = 1 - ((deltaProtagonist) / deltaAntagonist)
		} else {
			progFitness = 1 - ((deltaProtagonist) / deltaAntagonist)
			antFitness = (deltaProtagonist) / deltaAntagonist
		}
	} else {
		progFitness = 0
		antFitness = 1
	}

	return progFitness, antFitness, nil
}

// RatioFitness see RatioFitnessThresholded. In this case there is no threshold and everything is evaluated to the
// pure spec. If the test is able to reduce the difference between the spec and that created by the bug then the the
// test gains positive fitness, if it worsens it,
// it gains negative ratio. The rations can be viewed as percentages e.g. SpecMulti value = 0 . Bug value: 10,
// Test value: 5. The test in this case has brought back the bug value to a value of 5. In this case the test gets 0.
// 5 e.g. 5/10 of the fitness. If the test brought it back to 0, it would get 100 where the bug would get 0.
// If the test worsened the result and got 15, the test would get xxx
func RatioFitness(spec SpecMulti, antagonist, protagonist *Program) (antagonistFitness float64,
	protagonistFitness float64, err error) {

	if spec == nil {
		return math.MaxInt64, math.MaxInt64, fmt.Errorf("RatioFitness | spec cannot be nil when calculating Fitness")
	}
	if antagonist == nil {
		return math.MaxInt64, math.MaxInt64,
			fmt.Errorf("RatioFitness | antagonist cannot be nil when calculating Fitness")
	}
	if antagonist.T == nil {
		return math.MaxInt64, math.MaxInt64,
			fmt.Errorf("RatioFitness | antagonist tree cannot be nil when calculating Fitness")
	}
	if antagonist.T.root == nil {
		return math.MaxInt64, math.MaxInt64,
			fmt.Errorf("RatioFitness | antagonist tree root cannot be nil when calculating Fitness")
	}
	if protagonist == nil {
		return math.MaxInt64, math.MaxInt64,
			fmt.Errorf("RatioFitness | protagonist cannot be nil when calculating Fitness")
	}
	if protagonist.T == nil {
		return math.MaxInt64, math.MaxInt64,
			fmt.Errorf("RatioFitness | protagonist tree cannot be nil when calculating Fitness")
	}
	if protagonist.T.root == nil {
		return math.MaxInt64, math.MaxInt64,
			fmt.Errorf("RatioFitness | protagonist tree root cannot be nil when calculating Fitness")
	}

	err = protagonist.T.Validate()
	if err != nil {
		return math.MaxInt64, math.MaxInt64, err
	}

	err = antagonist.T.Validate()
	if err != nil {
		return math.MaxInt64, math.MaxInt64, err
	}

	// Antagonist
	antagonistMathematicalExpression, err := antagonist.T.ToMathematicalString()
	if err != nil {
		return math.MaxInt64, math.MaxInt64, err
	}

	protagonistMathematicalExpression, err := protagonist.T.ToMathematicalString()
	if err != nil {
		return math.MaxInt64, math.MaxInt64, err
	}

	protagonistFitness, antagonistFitness, err = ratioFitness(spec, antagonistMathematicalExpression,
		protagonistMathematicalExpression, antagonist, protagonist)
	if err != nil {
		return math.MaxInt64, math.MaxInt64, err
	}

	return antagonistFitness, protagonistFitness, nil
}

// ratioFitness assusmes the input variables have been checked for validity and nilness.
// If the protagonists delta is greater than that of the antagonist it automatically is given a value of 0.
// If the protagonist achieves identical spec values it obtains a value of 1 and the antagonists gets a value of 0
func ratioFitness(spec SpecMulti, antagonistExpression, protagonistExpression string,
	antagonistProgram, protagonistProgram *Program) (protagonistFitness,
	antagonistFitness float64, err error) {
	deltaProtagonist := 0.0
	deltaAntagonist := 0.0
	for _, s := range spec {
		dependentProtagonistVar, err := protagonistProgram.EvalMulti(s.Independents, protagonistExpression)
		if err != nil {
			return math.MaxFloat32, math.MaxFloat32, err
		}
		dP := calculateDelta(float64(s.Dependent), float64(dependentProtagonistVar))
		deltaProtagonist += dP

		dependentAntagonistVar, err := antagonistProgram.EvalMulti(s.Independents, antagonistExpression)
		if err != nil {
			return math.MaxFloat32, math.MaxFloat32, err
		}
		dA := calculateDelta(float64(s.Dependent), float64(dependentAntagonistVar))
		deltaAntagonist += dA
	}

	if deltaProtagonist >= deltaAntagonist {
		return 0, 1, nil
	}
	if deltaProtagonist == 0 {
		return 1, 0, nil
	}

	progFitness := 1 - ((deltaProtagonist) / deltaAntagonist)
	antFitness := (deltaProtagonist) / deltaAntagonist

	return progFitness, antFitness, nil
}

func calculateDelta(truth float64, value float64) float64 {
	return math.Abs((truth - value))
}

// AggregateFitness simply adds all the Fitness values of a given individual to come up with a total number.
// If the Fitness array is nil or empty return MaxInt8 as values such as -1 or 0 have a differnt meaning
func AggregateFitness(individual Individual) (float64, error) {
	if individual.Fitness == nil {
		return math.MaxInt8, fmt.Errorf("individuals Fitness arr cannot be nil")
	}
	if len(individual.Fitness) == 0 {
		return math.MaxInt8, fmt.Errorf("individuals Fitness arr cannot be empty")
	}

	sum := 0.0
	for i := range individual.Fitness {
		sum += individual.Fitness[i]
	}
	return sum, nil
}
