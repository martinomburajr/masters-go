package evolution

// Spec0 x * 0 = y
var Spec0 = SpecMulti{
	EquationPairing{Independents: IndependentVariableMap{"x": 0}, Dependent: 0},
	EquationPairing{Independents: IndependentVariableMap{"x": 1}, Dependent: 0},
	EquationPairing{Independents: IndependentVariableMap{"x": 2}, Dependent: 0},
	EquationPairing{Independents: IndependentVariableMap{"x": 3}, Dependent: 0},
	EquationPairing{Independents: IndependentVariableMap{"x": 4}, Dependent: 0},
}

//// SpecX x = y
//var SpecX = SpecMulti{
//	EquationPairing{0, 0},
//	EquationPairing{1, 1},
//	EquationPairing{2, 2},
//	EquationPairing{3, 3},
//	EquationPairing{4, 4},
//}
//
//// SpecXX 2x = y
//var Spec2X = SpecMulti{
//	EquationPairing{0, 0},
//	EquationPairing{1, 2},
//	EquationPairing{2, 4},
//	EquationPairing{3, 6},
//	EquationPairing{4, 8},
//}
//
//// SpecXX 2x = y
//var SpecXBy4 = SpecMulti{
//	EquationPairing{0, 0},
//	EquationPairing{1, 4},
//	EquationPairing{2, 8},
//	EquationPairing{3, 12},
//	EquationPairing{4, 16},
//}
//
//// SpecXX 2x = y
//var SpecXBy5 = SpecMulti{
//	EquationPairing{0, 0},
//	EquationPairing{1, 5},
//	EquationPairing{2, 10},
//	EquationPairing{3, 15},
//	EquationPairing{4, 20},
//}
//
//// SpecXX 2x+1 = y
//var Spec2XAdd1 = SpecMulti{
//	EquationPairing{-2, -3},
//	EquationPairing{-1, -1},
//	EquationPairing{0, 1},
//	EquationPairing{1, 3},
//	EquationPairing{2, 5},
//}
//
//// SpecXX x*x = y
//var SpecXX = SpecMulti{
//	EquationPairing{0, 0},
//	EquationPairing{1, 1},
//	EquationPairing{2, 4},
//	EquationPairing{3, 9},
//	EquationPairing{4, 16},
//}
//
//// SpecXXX x*x*x = y
//var SpecXXX = SpecMulti{
//	EquationPairing{-2, -8},
//	EquationPairing{-1, -1},
//	EquationPairing{0, 0},
//	EquationPairing{1, 1},
//	EquationPairing{2, 8},
//	EquationPairing{3, 27},
//	EquationPairing{4, 64},
//}
//
//// SpecXXXX x*x*x*x = y
//var SpecXXXX = SpecMulti{
//	EquationPairing{-2, 16},
//	EquationPairing{-1, 1},
//	EquationPairing{0, 0},
//	EquationPairing{1, 1},
//	EquationPairing{2, 16},
//	EquationPairing{3, 81},
//	EquationPairing{4, 256},
//}
//
//// SpecXXXXAdd4 x*x*x*x+4 = y
//var SpecXXXXAdd4 = SpecMulti{
//	EquationPairing{-2, 20},
//	EquationPairing{-1, 5},
//	EquationPairing{0, 0},
//	EquationPairing{1, 5},
//	EquationPairing{2, 20},
//	EquationPairing{3, 85},
//	EquationPairing{4, 260},
//}
