package evolution

import (
	"fmt"
	"strings"
)

// EquationPairing refers to a set dependent and independent values for a given equation.
// For example the equation x^2 + 1 has an equation pairing of {1, 0}, {2, 1}, {5,
// 2} for dependent and independent pairs respectively
type EquationPairing struct {
	Independent float64
	Dependent   float64
}

func (e *EquationPairing) ToString() string {
	return fmt.Sprintf("  %.2f  |  %.2f  \n", e.Independent, e.Dependent)
}

// EquationPairing refers to a set dependent and independent values for a given equation.
// For example the equation x^2 + 1 has an equation pairing of {1, 0}, {2, 1}, {5,
// 2} for dependent and independent pairs respectively
type EquationPairings struct {
	Independents IndependentVariableMap
	Dependent    float64
}

type IndependentVariableMap map[string]float64

func (e *EquationPairings) ToString() string {
	return fmt.Sprintf("  %.2f  |  %.2f  \n", e.Independents, e.Dependent)
}

type SpecMulti []EquationPairings

// GenerateSpec will create a spec given a valid mathematical expression.
// It is advised the mathematical expression contain an independent variable e.g. x.
// The initialSeed will be the starting value to evaluate.
// It is the callers responsibility to properly format the mathematical expression.
// Here are a few examples of valid mathematical expressions the function takes (Note the spacing between items)
// Example 1: x => x
// Example 2: ( x ) => x
// Example 3: ( x ) * ( x ) => x ^ 2
// Example 4: ( x ) + ( 2 ) => x + 2

func GenerateSpec(mathematicalExpression string, count int, initialSeed int) (SpecMulti, error) {
	if mathematicalExpression == "" {
		return SpecMulti{}, fmt.Errorf("GenerateSpec | cannot containe empty mathematical expression")
	}
	if count < 1 {
		return SpecMulti{}, fmt.Errorf("GenerateSpec | count cannot be less than 0")
	}
	spec := make([]EquationPairing, count)
	for i := range spec {
		independentVar := float32(initialSeed) + float32(i)
		dependentVariable, err := EvaluateMathematicalExpression(mathematicalExpression, independentVar)
		if err != nil {
			return SpecMulti{}, err
		}
		spec[i].Independent = independentVar
		spec[i].Dependent = dependentVariable
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

//
//func SpecMulti(spec SpecMulti) *Program {
//	p.spec = spec
//	return p
//}
//
//
//func  Validate() error {
//	return nil
//}
