package evolution

// EquationPairing refers to a set dependent and independent values for a given equation.
// For example the equation x^2 + 1 has an equation pairing of {1, 0}, {2, 1}, {5,
// 2} for dependent and independent pairs respectively
type EquationPairing struct {
	Independent float32
	Dependent   float32
}

type Spec []EquationPairing

//
//func Spec(spec Spec) *Program {
//	p.spec = spec
//	return p
//}
//
//
//func  Validate() error {
//	return nil
//}
