package evolution

import (
	"fmt"
	"math"
	"strings"
)

// EquationPairing refers to a set dependent and independent values for a given equation.
// For example the equation x^2 + 1 has an equation pairing of {1, 0}, {2, 1}, {5,
// 2} for dependent and independent pairs respectively
type EquationPairings struct {
	Independents IndependentVariableMap
	Dependent    float64
}

type IndependentVariableMap map[string]float64

func (e *EquationPairings) ToString() string {
	return fmt.Sprintf("  %#v  |  %.2f  \n", e.Independents, e.Dependent)
}

type SpecMulti []EquationPairings

// GenerateSpec will create a spec given a valid mathematical expression.
// It is advised the mathematical expression contain an independent variable e.g. x. or multiple x + a = y
// The initialSeed will be the starting value to evaluate.
// It is the callers responsibility to properly format the mathematical expression.
// Here are a few examples of valid mathematical expressions the function takes (Note the spacing between items)
// Example 1: x => x
// Example 2: ( x ) => x
// Example 3: ( x ) * ( x ) => x ^ 2
// Example 4: ( x ) + ( 2 ) => x + 2

func GenerateSpec(mathematicalExpression string, count int, initialSeed int) (SpecMulti,
	error) {
	if mathematicalExpression == "" {
		return nil, fmt.Errorf("GenerateSpec | cannot containe empty mathematical expression")
	}
	if count < 1 {
		return nil, fmt.Errorf("GenerateSpec | count cannot be less than 0")
	}
	spec := make([]EquationPairings, count)
	terminals, _, _, err := ParseStringLiberal(mathematicalExpression)
	if err != nil {
		return nil, err
	}

	for i := 0; i < count; i++ {
		for j, s := range spec {
			var x = map[string]float64{}
			for _, terminal := range terminals {
				x[terminal.value] = float64(initialSeed + i + j)
				s.Independents = x
				dependentVariable, err := EvaluateMathematicalExpression(mathematicalExpression, s.Independents)
				if err != nil {
					return nil, err
				}
				s.Dependent = dependentVariable
			}
		}
	}

	return spec, nil
}

// fillMap takes in a set of terminals and runs through various permutations to ensure each independent variable is
// populated with an item. A count of 0 or 1 will simply initialize the return value to ensure all the independent
// variables start off at the seed
func fillMap(terminals []SymbolicExpression, count int, seed float64) ([]map[string]float64, error) {
	// require
	if terminals == nil {
		return nil, fmt.Errorf("terminals cannot be nil")
	}
	if len(terminals) < 1 {
		return nil, fmt.Errorf("terminals cannot be empty")
	}
	if count < 0 {
		count = 0
	}
	if count > 5 {
		count = 5
	}

	// do
	x := map[string]float64{}

	// initialize
	for i := range terminals {
		x[terminals[i].value] = seed
	}

	combinations := int(math.Pow(float64(count), float64(len(terminals))))
	response := make([]map[string]float64, combinations)

	if count == 0 {
		response = append(response,  x)
		return response, nil
	}

	noOfPerms := permutation(rangeSlice(int(seed), int(seed) + count))
	return response, nil
}


func rangeSlice(start, stop int) []map[string]float64 {
	if start > stop {
		panic("Slice ends before it started")
	}
	xs := make([]int, stop-start)
	for i := 0; i < len(xs); i++ {
		xs[i] = i + 1 + start
	}
	return xs
}

func permutation(xs []int) (permuts [][]int) {
	var rc func([]int, int)
	rc = func(a []int, k int) {
		if k == len(a) {
			permuts = append(permuts, append([]int{}, a...))
		} else {
			for i := k; i < len(xs); i++ {
				a[k], a[i] = a[i], a[k]
				rc(a, k+1)
				a[k], a[i] = a[i], a[k]
			}
		}
	}
	rc(xs, 0)

	return permuts
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
