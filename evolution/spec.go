package evolution

import (
	"fmt"
	"strings"
)

// EquationPairing refers to a set dependent and independent values for a given equation.
// For example the equation x^2 + 1 has an equation pairing of {1, 0}, {2, 1}, {5,
// 2} for dependent and independent pairs respectively
type EquationPairing struct {
	Independents         IndependentVariableMap
	Dependent            float64
	ProtagonistThreshold float64
	AntagonistThreshold  float64
}

type IndependentVariableMap map[string]float64

func (e *EquationPairing) ToString() string {
	return fmt.Sprintf("  %#v  |  %.2f  \n", e.Independents, e.Dependent)
}

// SpecMulti is the underlying datascrutuctre that contains the spec as well as threshold information
type SpecMulti []EquationPairing

// GenerateSpec will create a spec given a valid mathematical expression.
// It is advised the mathematical expression contain an independent variable e.g. x. or multiple x + a = y
// The initialSeed will be the starting value to evaluate.
// It is the callers responsibility to properly format the mathematical expression.
// Here are a few examples of valid mathematical expressions the function takes (Note the spacing between items)
// Example 1: x => x
// Example 2: ( x ) => x
// Example 3: ( x ) * ( x ) => x ^ 2
// Example 4: ( x ) + ( 2 ) => x + 2
// This function should also work for multivariable elements. e.g. x+y+a+b=i where y,x,a,
// b are all independent variables i.e non-constants).
// antagonistThresholdMultiplier the threshold multipliers are important when using
// ControolledThresholdedAbsoluteRatioFitness.
// Esnure the antagonistThreshold is greater than that of the protagonist to demand antagonists drastically mutate
// trees to give more disparate values. For monoThresholding, you have to set both values to the same number
// CAN ONLY DO TWO DIFFERENT VARIABLES
//func GenerateSpec(mathematicalExpression string, independentVars []string, count int, initialSeed int,
//	antagonistThresholdMultiplier, protagonistThresholdMultiplier float64, params EvolutionParams) (SpecMulti,
//	error) {
//
//	if mathematicalExpression == "" {
//		return nil, fmt.Errorf("GenerateSpec | cannot containe empty mathematical expression")
//	}
//	if count < 1 {
//		return nil, fmt.Errorf("GenerateSpec | count cannot be less than 0")
//	}
//	if count >= 5 {
//		count = 3
//	}
//	if antagonistThresholdMultiplier < 1 {
//		antagonistThresholdMultiplier = 1
//	}
//	if protagonistThresholdMultiplier < 1 {
//		protagonistThresholdMultiplier = 1
//	}
//	spec := make([]EquationPairing, count)
//
//	//1. Determine number of unique independent variables and kings of independent variables e.g. x, y
//	// pass in unique independent variables as a slice of strings? [OK]
//
//	// determine the number permutationsCount we can create count^independentVars.len
//	numVars := len(independentVars)
//	if independentVars == nil || numVars < 1 {
//		for i, _ := range spec {
//			dependentVariable, err := EvaluateMathematicalExpression(mathematicalExpression, nil, params)
//			if err != nil {
//				return nil, err
//			}
//			spec[i].Dependent = dependentVariable
//			spec[i].AntagonistThreshold = dependentVariable * antagonistThresholdMultiplier
//			spec[i].ProtagonistThreshold = dependentVariable * protagonistThresholdMultiplier
//		}
//		return spec, nil
//	}
//
//	if numVars < 2 {
//		for i := range spec {
//			spec[i].Independents = map[string]float64{}
//			//for j := 0; j < count; j++ {
//			spec[i].Independents[independentVars[0]] = float64(i + initialSeed)
//			//}
//			dependentVariable, err := EvaluateMathematicalExpression(mathematicalExpression,
//				spec[i].Independents, params)
//			if err != nil {
//				return nil, err
//			}
//			spec[i].Dependent = dependentVariable
//		}
//		return spec, nil
//	}
//
//	g := make([]int, 0)
//	for i := 0; i < count; i++ {
//		g = append(g, i+initialSeed)
//	}
//
//	permutationsCount := int(math.Pow(float64(count), float64(numVars)))
//	spec = make([]EquationPairing, permutationsCount)
//	permutationsWithRepetitions := utils.PermutationsWithRepetitions(g)
//
//	for i := range spec {
//		spec[i].Independents = map[string]float64{}
//		for j := range independentVars {
//			spec[i].Independents[independentVars[j]] = float64(permutationsWithRepetitions[i][j])
//		}
//		dependentVariable, err := EvaluateMathematicalExpression(mathematicalExpression,
//			spec[i].Independents, params)
//		if err != nil {
//			return nil, err
//		}
//		spec[i].Dependent = dependentVariable
//	}
//	return spec, nil
//}

// GenerateSpecSimple assumes a single independent variable x with an unlimited count.
func GenerateSpecSimple(specParam SpecParam, fitnessStrategy FitnessStrategy, fitnessCalculatorType int) (SpecMulti,
	error) {

	if specParam.Expression == "" {
		return nil, fmt.Errorf("GenerateSpec | cannot containe empty mathematical expression")
	}
	if specParam.Range < 1 {
		return nil, fmt.Errorf("GenerateSpec | specParam.Range cannot be less than 0")
	}
	if fitnessStrategy.AntagonistThresholdMultiplier < 1 {
		fitnessStrategy.AntagonistThresholdMultiplier = 1
	}
	if fitnessStrategy.ProtagonistThresholdMultiplier < 1 {
		fitnessStrategy.ProtagonistThresholdMultiplier = 1
	}

	spec := make([]EquationPairing, specParam.Range)
	for i := range spec {
		spec[i].Independents = map[string]float64{}
		spec[i].Independents["x"] = float64(i + specParam.Seed)
		dependentVariable, err := EvaluateMathematicalExpression(specParam.Expression,
			spec[i].Independents, fitnessCalculatorType)
		if err != nil {
			return nil, err
		}
		spec[i].Dependent = dependentVariable
		spec[i].AntagonistThreshold = dependentVariable * fitnessStrategy.AntagonistThresholdMultiplier
		spec[i].ProtagonistThreshold = dependentVariable * fitnessStrategy.ProtagonistThresholdMultiplier
	}
	return spec, nil
}

func (spec SpecMulti) ToString() string {
	sb := strings.Builder{}
	if spec == nil {
		return sb.String()
	}

	sb.WriteString("  x  :  f(x)  \n")
	for i := range spec {
		s := spec[i].ToString()
		sb.WriteString(s)
	}
	return sb.String()
}
