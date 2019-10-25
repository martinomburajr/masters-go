package evolution

import (
	"fmt"
	"math"
)

const (
	FitnessAbsolute                   = "FitnessAbsolute"
	FitnessThresholdedAntagonistRatio = "FitnessThresholdedAntagonistRatio"
	FitnessProtagonistThresholdTally  = "FitnessProtagonistThresholdTally"
	FitnessRatio                      = "FitnessRatio"
	FitnessMonoThresholdedRatio       = "FitnessMonoThresholdedRatio"
	FitnessDualThresholdedRatio       = "FitnessDualThresholdedRatio"
)

//type IFitness interface {
//func (spec SpecMulti, antagonist, protagonist *Program) (err error)
//}

// FitnessAbsolute calculates Fitness for two competing individuals.
// The absolute value from the spec is obtained summed, and given to each individual. Smaller values are better.
// 0 being the absolute best.
func AbsoluteFitness(spec SpecMulti, protagonistExpression, antagonistExpression string,
	protagonistProgram, antagonistProgram *Program, fitnessCalculatorType int) (protagonistFitness,
	antagonistFitness float64, err error) {
	deltaProtagonist := 0.0
	deltaAntagonist := 0.0
	for _, s := range spec {
		dependentProtagonistVar, err := protagonistProgram.EvalMulti(s.Independents, protagonistExpression, fitnessCalculatorType)
		if err != nil {
			return math.MaxInt64, math.MaxInt64, err
		}
		dP := calculateDelta(float64(s.Dependent), float64(dependentProtagonistVar))
		deltaProtagonist += dP

		dependentAntagonistVar, err := antagonistProgram.EvalMulti(s.Independents, antagonistExpression, fitnessCalculatorType)
		if err != nil {
			return math.MaxInt64, math.MaxInt64, err
		}
		dA := calculateDelta(float64(s.Dependent), float64(dependentAntagonistVar))
		deltaAntagonist += dA
	}

	return deltaProtagonist, deltaAntagonist, nil
}

// ThresholdedAntagonistRatioFitness thresholds only the antagonist preventing it from playing the protagonists game.
// The protagonists fitness is ratio'd to the antagonists fitness and not threshold.
// Whereas the antagonists fitness is ratio'd to its threshold
func ThresholdedAntagonistRatioFitness(spec SpecMulti, antagonist, protagonist *Program,
	fitnessCalculatorType int) (antagonistFitness float64,
	protagonistFitness float64, err error) {

	err = fitnessParameterValidator(spec, antagonist, protagonist)
	if err != nil {
		return math.MaxInt64, math.MaxInt64, err
	}

	return evaluateFitnessAntagonistThresholded(spec, antagonist, protagonist, fitnessCalculatorType)
}

// evaluateFitnessAntagonistThresholded performs fitness evaluation using the given antagonist and protagonist.
// This strategy places a threshold on the antagonist and none on the protagonist. The results are ratiod,
// the antagonist has to exceed its ratio to gain fitness, the protagonist has to improve on the antagonists result
func evaluateFitnessAntagonistThresholded(spec SpecMulti, antagonist, protagonist *Program, fitnessCalculatorType int) (antagonistFitness,
	protagonistFitness float64, err error) {

	antagonistExpression, protagonistExpression, err := generateExpressions(antagonist, protagonist)
	if err != nil {
		return 10000.001, 10000.001, err
	}
	deltaProtagonist := 0.0
	deltaAntagonist := 0.0
	deltaAntagonistThreshold := 0.0
	for i := range spec {
		dependentAntagonistVar, err := antagonist.EvalMulti(spec[i].Independents, antagonistExpression, fitnessCalculatorType)
		if err != nil {
			return 10000.001, 10000.001, err
		}
		dA := calculateDelta(spec[i].Dependent, dependentAntagonistVar)
		deltaAntagonist += dA
		deltaAntagonistThreshold += math.Abs(spec[i].AntagonistThreshold)

		dependentProtagonistVar, err := protagonist.EvalMulti(spec[i].Independents, protagonistExpression, fitnessCalculatorType)
		if err != nil {
			return 10000.001, 10000.001, err
		}
		deltaTruthToProtagonistVar := calculateDelta(spec[i].Dependent, dependentProtagonistVar)
		deltaProtagonist += deltaTruthToProtagonistVar
	}

	specLen := float64(len(spec))
	deltaProtagonist = deltaProtagonist / specLen
	deltaAntagonist = deltaAntagonist / specLen
	deltaAntagonistThreshold = deltaAntagonistThreshold / specLen

	//antagonists
	if deltaAntagonist >= deltaAntagonistThreshold {
		antagonistFitness = (deltaAntagonist - deltaAntagonistThreshold) / deltaAntagonist
	} else {
		antagonistFitness = -1 * ((deltaAntagonistThreshold - deltaAntagonist) / deltaAntagonistThreshold)
	}

	if deltaProtagonist <= deltaAntagonist {
		protagonistFitness = deltaProtagonist / deltaAntagonist
	} else {
		protagonistFitness = -1 * ((deltaProtagonist) / deltaAntagonist)
	}

	return antagonistFitness, protagonistFitness, nil
}

// #3 ProtagonistThresholdTally takes only the protagonist and checks to see if the Protagonist(Antagonist(
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

// RatioFitness see RatioFitnessThresholded. In this case there is no threshold and everything is evaluated to the
// pure spec. If the test is able to reduce the difference between the spec and that created by the bug then the the
// test gains positive fitness, if it worsens it,
// it gains negative ratio. The rations can be viewed as percentages e.g. SpecMulti value = 0 . Bug value: 10,
// Test value: 5. The test in this case has brought back the bug value to a value of 5. In this case the test gets 0.
// 5 e.g. 5/10 of the fitness. If the test brought it back to 0, it would get 100 where the bug would get 0.
// If the test worsened the result and got 15, the test would get xxx
func RatioFitness(spec SpecMulti, antagonist, protagonist *Program,
	fitnessCalculatorType int) (antagonistFitness float64,
	protagonistFitness float64, err error) {

	err = fitnessParameterValidator(spec, antagonist, protagonist)
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
		protagonistMathematicalExpression, antagonist, protagonist, fitnessCalculatorType)
	if err != nil {
		return math.MaxInt64, math.MaxInt64, err
	}

	return antagonistFitness, protagonistFitness, nil
}

// ratioFitness assusmes the input variables have been checked for validity and nilness.
// If the protagonists delta is greater than that of the antagonist it automatically is given a value of 0.
// If the protagonist achieves identical spec values it obtains a value of 1 and the antagonists gets a value of 0
func ratioFitness(spec SpecMulti, antagonistExpression, protagonistExpression string,
	antagonistProgram, protagonistProgram *Program, fitnessCalculatorType int) (protagonistFitness,
	antagonistFitness float64, err error) {
	deltaProtagonist := 0.0
	deltaAntagonist := 0.0

	//antString, err := antagonistProgram.T.ToMathematicalString()
	//fmt.Printf("##### ANTAGONIST ######\n%s\n", antString)
	//proString, err := protagonistProgram.T.ToMathematicalString()
	//fmt.Printf("##### PROTAGONIST ######\n%s\n", proString)

	for _, s := range spec {
		dependentAntagonistVar, err := antagonistProgram.EvalMulti(s.Independents, antagonistExpression, fitnessCalculatorType)
		if err != nil {
			return math.MaxInt64, math.MaxInt64, err
		}
		dA := calculateDelta(float64(s.Dependent), float64(dependentAntagonistVar))
		deltaAntagonist += dA

		dependentProtagonistVar, err := protagonistProgram.EvalMulti(s.Independents, protagonistExpression, fitnessCalculatorType)
		if err != nil {
			return math.MaxInt64, math.MaxInt64, err
		}
		dP := calculateDelta(float64(s.Dependent), float64(dependentProtagonistVar))
		deltaProtagonist += dP
	}

	specLen := float64(len(spec))
	deltaProtagonist = deltaProtagonist / specLen
	deltaAntagonist = deltaAntagonist / specLen

	if deltaProtagonist >= deltaAntagonist {
		return 0, 1, nil
	}
	if deltaProtagonist == 0 {
		return 1, 0.00, nil
	}
	if deltaAntagonist == 0 {
		return 0, 0, nil
	}

	progFitness := 1 - ((deltaProtagonist) / deltaAntagonist)
	antFitness := (deltaProtagonist) / deltaAntagonist

	return progFitness, antFitness, nil
}

// ThresholdedRatioFitness is a means to calculate fitness that spurs the protagonists and
// antagonists to do their best. It works by applying a threshold criteria that is based on the incoming spec.
// A mono threshold is applied by setting the protagonist and antagonist threshold values to the same value,
// this is done automatically when you select the fitness strategy at the start of the evolutionary process.
// Both individuals have to fall on their respective side and either edge closer to delta-0 for the protagonist or
// delta-infitinite for the antagonist. If they happen to fall on opposite sides they attain at most -1
// A dual threshold is used to punish both antagonist of protagonist for not performing as expected.
//// This fitness strategy works by comparing the average delta values of both protagonist and antagonist with a
//// specified threshold. Their deltas are not compared against each other as with other fitness strategies,
//// the thresholds act as markers of performance for each.
//// The porotagonist and antagonist each have their own threshold values that are embeded in the SpecMulti data
//// structure. Note this only compares the average and not the total deltas
func ThresholdedRatioFitness(spec SpecMulti, antagonist, protagonist *Program,
	fitnessCalculatorType int) (antagonistFitness float64,
	protagonistFitness float64, err error) {
	err = fitnessParameterValidator(spec, antagonist, protagonist)
	if err != nil {
		return math.MaxInt64, math.MaxInt64, err
	}

	return thresholdedRatioFitness(spec, antagonist, protagonist, fitnessCalculatorType)
}

// thresholdedRatioFitness performs fitness evaluation using the given antagonist and protagonist.
// It returns information regarding thresholds as well,
// they can be ignored if the function does not require information on the thresholds.
// Furthermore these values are averaged based on the length of the spec.
// A nil or empty spec will throw an error
func thresholdedRatioFitness(spec SpecMulti, antagonist, protagonist *Program,
	fitnessCalculatorType int) (antagonistFitness,
	protagonistFitness float64, err error) {

	antagonistExpression, protagonistExpression, err := generateExpressions(antagonist, protagonist)
	if err != nil {
		return 10000.001, 10000.001, err
	}
	deltaProtagonist := 0.0
	deltaAntagonist := 0.0
	deltaAntagonistThreshold := 0.0
	deltaProtagonistThreshold := 0.0
	for i := range spec {
		dependentAntagonistVar, err := antagonist.EvalMulti(spec[i].Independents, antagonistExpression, fitnessCalculatorType)
		if err != nil {
			return 10000.001, 10000.001, err
		}
		dA := calculateDelta(spec[i].Dependent, dependentAntagonistVar)
		deltaAntagonist += dA
		//dAT := calculateDelta(spec[i].AntagonistThreshold, dependentAntagonistVar)
		deltaAntagonistThreshold += math.Abs(spec[i].AntagonistThreshold)

		dependentProtagonistVar, err := protagonist.EvalMulti(spec[i].Independents, protagonistExpression, fitnessCalculatorType)
		if err != nil {
			return 10000.001, 10000.001, err
		}
		deltaTruthToProtagonistVar := calculateDelta(spec[i].Dependent, dependentProtagonistVar)
		deltaProtagonist += deltaTruthToProtagonistVar
		//dPT := calculateDelta(spec[i].ProtagonistThreshold, dependentProtagonistVar)
		deltaProtagonistThreshold += math.Abs(spec[i].ProtagonistThreshold)
	}

	specLen := float64(len(spec))
	deltaProtagonist = deltaProtagonist / specLen
	deltaAntagonist = deltaAntagonist / specLen
	deltaAntagonistThreshold = deltaAntagonistThreshold / specLen
	deltaProtagonistThreshold = deltaProtagonistThreshold / specLen

	//antagonists
	if deltaAntagonist >= deltaAntagonistThreshold {
		antagonistFitness = (deltaAntagonist - deltaAntagonistThreshold) / deltaAntagonist
	} else {
		antagonistFitness = -1 * ((deltaAntagonistThreshold - deltaAntagonist) / deltaAntagonistThreshold)
	}

	if deltaProtagonist <= deltaProtagonistThreshold {
		protagonistFitness = (deltaProtagonistThreshold - deltaProtagonist) / deltaProtagonistThreshold
	} else {
		protagonistFitness = -1 * ((deltaProtagonist - deltaProtagonistThreshold) / deltaProtagonist)
	}

	return antagonistFitness, protagonistFitness, nil
}

// fitnessParameterValidator is a convenience function that evaluates the input parameters to a fitness argument
func fitnessParameterValidator(spec SpecMulti, antagonist, protagonist *Program) (err error) {
	if spec == nil {
		return fmt.Errorf("RatioFitness | spec cannot be nil when calculating Fitness")
	}
	if len(spec) < 1 {
		return fmt.Errorf("RatioFitness | spec cannot be empty when calculating Fitness")
	}
	if antagonist == nil {
		return fmt.Errorf("RatioFitness | antagonist cannot be nil when calculating Fitness")
	}
	if antagonist.T == nil {
		return fmt.Errorf("RatioFitness | antagonist tree cannot be nil when calculating Fitness")
	}
	if antagonist.T.root == nil {
		return fmt.Errorf("RatioFitness | antagonist tree root cannot be nil when calculating Fitness")
	}
	if protagonist == nil {
		return fmt.Errorf("RatioFitness | protagonist cannot be nil when calculating Fitness")
	}
	if protagonist.T == nil {
		return fmt.Errorf("RatioFitness | protagonist tree cannot be nil when calculating Fitness")
	}
	if protagonist.T.root == nil {
		return fmt.Errorf("RatioFitness | protagonist tree root cannot be nil when calculating Fitness")
	}

	err = antagonist.T.Validate()
	if err != nil {
		return err
	}

	err = protagonist.T.Validate()
	if err != nil {
		return err
	}

	return nil
}

// generateExpressions returns a set of mathematical expressions of the antagonist and protagonist trees.
func generateExpressions(antagonist, protagonist *Program) (antagonistExpression, protagonistExpression string,
	err error) {
	antagonistMathematicalExpression, err := antagonist.T.ToMathematicalString()
	if err != nil {
		return "", "", err
	}

	protagonistMathematicalExpression, err := protagonist.T.ToMathematicalString()
	if err != nil {
		return "", "", err
	}

	return antagonistMathematicalExpression, protagonistMathematicalExpression, nil
}

// calculateDelta calculates the absolute value between the truth and the supplied value
func calculateDelta(truth float64, value float64) float64 {
	return math.Abs(truth - value)
}
