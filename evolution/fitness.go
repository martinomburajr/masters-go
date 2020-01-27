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

	DivByZeroIgnore           = "Ignore"
	DivByZeroSteadyPenalize   = "SteadyPenalization"
	DivByZeroPenalize         = "Penalize"
	DivByZeroSetSpecValueZero = "SetSpecValueZero"
)

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
	divByZeroStrategy string) (antagonistFitness float64,
	protagonistFitness, antagonistFitnessDelta, protagonistFitnessDelta float64, err error) {
	err = fitnessParameterValidator(spec, antagonist, protagonist)
	if err != nil {
		return math.MaxInt64, math.MaxInt64, math.MaxInt64, math.MaxInt64, err
	}

	return thresholdedRatioFitness(spec, antagonist, protagonist, divByZeroStrategy)
}

// thresholdedRatioFitness performs fitness evaluation using the given antagonist and protagonist.
// It returns information regarding thresholds as well,
// they can be ignored if the function does not require information on the thresholds.
// Furthermore these values are averaged based on the length of the spec.
// A nil or empty spec will throw an error. It takes advantage of RMSE
func thresholdedRatioFitness(spec SpecMulti, antagonist, protagonist *Program,
	divByZeroStrategy string) (antagonistFitness,
	protagonistFitness, antagonistFitnessError, protagonistFitnessError float64, err error) {

	fitnessPenalization := spec[0].DivideByZeroPenalty
	badDeltaValue := math.Inf(1)

	antagonistExpression, protagonistExpression, err := generateExpressions(antagonist, protagonist)
	if err != nil {
		return fitnessPenalization, fitnessPenalization, fitnessPenalization, fitnessPenalization, err
	}
	deltaProtagonist := 0.0
	deltaAntagonist := 0.0
	deltaAntagonistThreshold := 0.0
	deltaProtagonistThreshold := 0.0
	antagonistDividedByZeroCount := 0
	protagonistDividedByZeroCount := 0
	isAntagonistValid := true
	isProtagonistValid := true

	for i := range spec {
		independentX := spec[i].Independents
		independentXVal := spec[i].Independents["x"]
		if isAntagonistValid {
			dependentAntagonistVar, err := antagonist.EvalMulti(independentX, antagonistExpression)
			if err != nil {
				switch divByZeroStrategy {
				case DivByZeroIgnore:
					dependentAntagonistVar = 0

				case DivByZeroPenalize:
					isAntagonistValid = false
					antagonistFitness = fitnessPenalization

				case DivByZeroSetSpecValueZero:
					dependentAntagonistVar = 0

				case DivByZeroSteadyPenalize:
					if independentXVal != 0 {
						// If the spec does not contain a zero,
						// yet you still divide by zero. Give maximum penalty!
						if math.IsNaN(dependentAntagonistVar) || dependentAntagonistVar == 0 {
							isAntagonistValid = false
							antagonistFitness = fitnessPenalization
						} else {
							antagonistDividedByZeroCount++
						}
					} else {
						// Unlikely to ever reach here
						antagonistDividedByZeroCount++
					}
				}
			} else {
				diff := (spec[i].Dependent - dependentAntagonistVar)
				deltaAntagonist += diff * diff
			}
		}
		deltaAntagonistThreshold += (math.Abs(spec[i].AntagonistThreshold) * math.Abs(spec[i].AntagonistThreshold))
		if isProtagonistValid {
			dependentProtagonistVar, err := protagonist.EvalMulti(independentX, protagonistExpression)
			if err != nil {
				switch divByZeroStrategy {
				case DivByZeroIgnore:
					dependentProtagonistVar = 0

				case DivByZeroPenalize:
					isProtagonistValid = false
					protagonistFitness = fitnessPenalization

				case DivByZeroSetSpecValueZero:
					dependentProtagonistVar = 0

				case DivByZeroSteadyPenalize:
					if independentXVal != 0 {
						if math.IsNaN(dependentProtagonistVar) || dependentProtagonistVar == 0 {
							isProtagonistValid = false
							protagonistFitness = fitnessPenalization
						} else {
							protagonistDividedByZeroCount++
						}
					} else {
						// Unlikely to ever reach here
						protagonistDividedByZeroCount++
					}
				}
			} else {
				diff := (spec[i].Dependent - dependentProtagonistVar)
				deltaProtagonist += diff * diff
			}
		}
		deltaProtagonistThreshold += (spec[i].ProtagonistThreshold * spec[i].ProtagonistThreshold)
	}

	specLen := float64(len(spec))

	deltaAntagonist = math.Sqrt(deltaAntagonist / specLen)
	deltaProtagonist = math.Sqrt(deltaProtagonist / specLen)
	deltaAntagonistThreshold = math.Sqrt(deltaAntagonistThreshold / specLen)
	deltaProtagonistThreshold = math.Sqrt(deltaProtagonistThreshold / specLen)

	if !isProtagonistValid && !isAntagonistValid {
		// TODO is math.Nan the best alternative?
		return fitnessPenalization, fitnessPenalization, math.NaN(), math.NaN(), nil
	} else if !isProtagonistValid && isAntagonistValid {
		if deltaAntagonist >= deltaAntagonistThreshold { // good thing
			if deltaAntagonist == 0 { // This is to punish deltaAntagonist for coalescing near the spec
				antagonistFitness = -1
				antagonistDividedByZeroCount = -1
			} else {
				// Award fitness if it did not cluster around the spec
				antagonistFitness = (deltaAntagonist - deltaAntagonistThreshold) / deltaAntagonist
			}
		} else {
			antagonistFitness = -1 * ((deltaAntagonistThreshold - deltaAntagonist) / deltaAntagonistThreshold)
		}
		protagonistFitness = fitnessPenalization
		deltaProtagonist = badDeltaValue
		return antagonistFitness, protagonistFitness, deltaAntagonist, deltaProtagonist, nil
	} else if !isAntagonistValid && isProtagonistValid {
		if deltaProtagonist <= deltaProtagonistThreshold {
			if deltaProtagonist == 0 {
				protagonistFitness = 1
				protagonistDividedByZeroCount = -1
			} else {
				protagonistFitness = (deltaProtagonistThreshold - deltaProtagonist) / deltaProtagonistThreshold
			}
		} else {
			protagonistFitness = -1 * ((deltaProtagonist - deltaProtagonistThreshold) / deltaProtagonist)
			//protagonistFitness = -1 * (deltaProtagonist / deltaAntagonist)
		}
		antagonistFitness = fitnessPenalization
		deltaAntagonist = badDeltaValue
		return antagonistFitness, protagonistFitness, deltaAntagonist, deltaProtagonist, nil
	} else {
		//antagonists
		if deltaAntagonist >= deltaAntagonistThreshold {
			if deltaAntagonist == 0 {
				antagonistFitness = -1 // This is to punish deltaAntagonist for coalescing near the spec
				antagonistDividedByZeroCount = -1
			} else {
				antagonistFitness = (deltaAntagonist - deltaAntagonistThreshold) / deltaAntagonist
			}
		} else {
			antagonistFitness = -1 * ((deltaAntagonistThreshold - deltaAntagonist) / deltaAntagonistThreshold)
		}

		//protagonist
		if deltaProtagonist <= deltaProtagonistThreshold {
			if deltaProtagonist == 0 {
				protagonistFitness = 1
				protagonistDividedByZeroCount = -1
			} else {
				protagonistFitness = (deltaProtagonistThreshold - deltaProtagonist) / deltaProtagonistThreshold
			}
		} else {
			protagonistFitness = -1 * ((deltaProtagonist - deltaProtagonistThreshold) / deltaProtagonist)
		}

		if antagonistDividedByZeroCount > 0 {
			if antagonistFitness > 0 {
				antagonistFitness = antagonistFitness - (antagonistFitness * 0.1 * float64(
					antagonistDividedByZeroCount))
			}
			// No else statement as if the antagonist is already less than 0, it should remain there.
		}
		if protagonistDividedByZeroCount > 0 {
			if protagonistFitness > 0 {
				protagonistFitness = protagonistFitness - (protagonistFitness * 0.1 * float64(
					protagonistDividedByZeroCount))
			}
		}
		return antagonistFitness, protagonistFitness, deltaAntagonist, deltaProtagonist, nil
	}
}

// RMSE (Root Mean Squared Error) Root Mean Square Error (RMSE) is the standard deviation of the residuals
// (prediction errors). Residuals are a measure of how far from the regression line data points are;
// RMSE is a measure of how spread out these residuals are. In other words,
// it tells you how concentrated the data is around the line of best fit.
// Root mean square error is commonly used in climatology,
// forecasting, and regression analysis to verify experimental results.
func RMSE(forecast, observed []float64) float64 {
	size := len(observed)

	sumDiffSquared := 0.0
	for i := 0; i < size; i++ {
		diff := forecast[i] - observed[i]
		diffSquared := diff * diff
		sumDiffSquared += diffSquared
	}

	normalized := sumDiffSquared / float64(size)

	return math.Sqrt(normalized)
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
