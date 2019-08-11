package program

// EquationPairing refers to a set dependent and independent values for a given equation.
// For example the equation x^2 + 1 has an equation pairing of {1, 0}, {2, 1}, {5,
// 2} for dependent and independent pairs respectively
type EquationPairing struct {
	Dependent   float32
	Independent float32
}

type Spec []*EquationPairing
