package evolution

import (
	"reflect"
	"testing"
)

func TestGenerateSpec(t *testing.T) {
	type args struct {
		mathematicalExpression        string
		independentVars               []string
		count                         int
		initialSeed                   int
		antagonistTresholdMultiplier  float64
		protagonistTresholdMultiplier float64
	}
	tests := []struct {
		name    string
		args    args
		want    SpecMulti
		wantErr bool
	}{
		//{"empty", args{mathematicalExpression: "", independentVars:nil, count: 0, initialSeed: 0}, nil, true},
		//{"negative count", args{mathematicalExpression: "x", count: -1, initialSeed: 0}, nil, true},
		//{"0 countt", args{mathematicalExpression: "1", count: 0, initialSeed: 0}, nil, true},
		//{"bad string", args{mathematicalExpression: "$", count: 1, initialSeed: 0}, nil, true},
		//{"1", args{mathematicalExpression: "1", count: 1, initialSeed: 0}, SpecMulti{
		//	EquationPairing{
		//		Independents:nil,
		//		Dependent: 1}},
		//false},
		//{"1", args{mathematicalExpression: "1", count: 1, initialSeed: 1}, SpecMulti{
		//	EquationPairing{
		//		Independents:nil,
		//		Dependent: 1}},
		//	false},
		//{"x | ", args{mathematicalExpression: "x", count: 1, initialSeed: 0}, nil,
		//	true},
		//{"x | bad indp", args{mathematicalExpression: "x", independentVars:[]string{"a"}, count: 1, initialSeed: 0},
		//	nil,
		//	true},
		//
		//	{"x | indp", args{mathematicalExpression: "x", independentVars:[]string{"x"}, count: 1, initialSeed: 0},
		//		SpecMulti{
		//			EquationPairing{
		//				Independents:IndependentVariableMap{"x": 0},
		//				Dependent: 0}},
		//false},
		//	{"x | indp", args{mathematicalExpression: "x", independentVars:[]string{"x"}, count: 2, initialSeed: 0},
		//		SpecMulti{
		//			EquationPairing{Independents:IndependentVariableMap{"x": 0}, Dependent: 0},
		//			EquationPairing{Independents:IndependentVariableMap{"x": 1}, Dependent: 1},
		//		},
		//		false},

		//{"x + y | count 1", args{mathematicalExpression: "x + y", independentVars: []string{"x", "y"}, count: 1,
		//	initialSeed: 0},
		//	SpecMulti{
		//		EquationPairing{Independents: IndependentVariableMap{"x": 0, "y": 0}, Dependent: 0},
		//		//EquationPairing{Independents:IndependentVariableMap{"x": 1}, Dependent: 1},
		//	},
		//	false},

		//{"x + y | count 2", args{mathematicalExpression: "x + y", independentVars: []string{"x", "y"}, count: 2,
		//	initialSeed: 0},
		//	SpecMulti{
		//		EquationPairing{Independents: IndependentVariableMap{"x": 0, "y": 0}, Dependent: 0},
		//		EquationPairing{Independents: IndependentVariableMap{"x": 1, "y": 0}, Dependent: 1},
		//		EquationPairing{Independents: IndependentVariableMap{"x": 0, "y": 1}, Dependent: 1},
		//		EquationPairing{Independents: IndependentVariableMap{"x": 1, "y": 1}, Dependent: 2},
		//	},
		//	false},
		//
		//{"x + y | count 2", args{mathematicalExpression: "x + y + 1", independentVars: []string{"x", "y"}, count: 2,
		//	initialSeed: 0},
		//	SpecMulti{
		//		EquationPairing{Independents: IndependentVariableMap{"x": 0, "y": 0}, Dependent: 1},
		//		EquationPairing{Independents: IndependentVariableMap{"x": 1, "y": 0}, Dependent: 2},
		//		EquationPairing{Independents: IndependentVariableMap{"x": 0, "y": 1}, Dependent: 2},
		//		EquationPairing{Independents: IndependentVariableMap{"x": 1, "y": 1}, Dependent: 3},
		//	},
		//	false},
		//
		//{"x + y + z | count 2", args{mathematicalExpression: "x + y + z", independentVars: []string{"x", "y", "z"},
		//	count: 2,
		//	initialSeed: 0},
		//	SpecMulti{
		//		EquationPairing{Independents: IndependentVariableMap{"x": 0, "y": 0}, Dependent: 0},
		//		EquationPairing{Independents: IndependentVariableMap{"x": 1, "y": 0}, Dependent: 1},
		//		EquationPairing{Independents: IndependentVariableMap{"x": 0, "y": 1}, Dependent: 1},
		//		EquationPairing{Independents: IndependentVariableMap{"x": 1, "y": 1}, Dependent: 2},
		//	},
		//	false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateSpec(tt.args.mathematicalExpression, tt.args.independentVars, tt.args.count,
				tt.args.initialSeed, tt.args.antagonistTresholdMultiplier, tt.args.protagonistTresholdMultiplier)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSpec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func Test_fillMap(t *testing.T) {
//	type args struct {
//		terminals []SymbolicExpression
//		count     int
//		seed      float64
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    []map[string]float64
//		wantErr bool
//	}{
//		//{"nil", args{nil, 0, 0}, nil, true},
//		//{"empty", args{[]SymbolicExpression{}, 0, 0}, nil, true},
//		//{"val-1", args{[]SymbolicExpression{{value:"1", arity:0}}, 0, 0},
//		//	[]map[string]float64{{"1": 0}}, false},
//		//{"val-1 | seed -3", args{[]SymbolicExpression{{value:"1", arity:0}}, 0, -3},
//		//	[]map[string]float64{{"1": -3}}, false},
//		//{"val-x", args{[]SymbolicExpression{{value:"x", arity:0}}, 0, 0},
//		//	[]map[string]float64{{"x": 0}}, false},
//		//{"val-x | count 1", args{[]SymbolicExpression{{value:"x", arity:0}}, 1, 0},
//		//	[]map[string]float64{{"x": 0}}, false},
//		//{"val-x | count 2", args{[]SymbolicExpression{{value:"x", arity:0}}, 2, 0},
//		//	[]map[string]float64{{"x": 0}, {"x": 1}}, false},
//		//{"val-x | val- y | count 1", args{[]SymbolicExpression{{value:"x", arity:0}, {value:"y", arity:0}}, 1, 0},
//		//	[]map[string]float64{{"x": 0, "y": 0}}, false},
//		{"val-x | val- y | count 2", args{[]SymbolicExpression{{value: "x", arity: 0}, {value: "y", arity: 0}}, 2, 0},
//			[]map[string]float64{
//				{"x": 0, "y": 0},
//				{"x": 0, "y": 1},
//				{"x": 1, "y": 0},
//				{"x": 1, "y": 1},
//			}, false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := fillMap(tt.args.terminals, tt.args.count, tt.args.seed)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("fillMap() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("fillMap() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestGenerateSpecSimple(t *testing.T) {
	type args struct {
		mathematicalExpression        string
		count                         int
		initialSeed                   int
		antagonistTresholdMultiplier  float64
		protagonistTresholdMultiplier float64
	}
	tests := []struct {
		name    string
		args    args
		want    SpecMulti
		wantErr bool
	}{
		{"nil", args{mathematicalExpression: "", count: -1, initialSeed: -1}, nil, true},
		{"invalid symbol", args{mathematicalExpression: "#", count: -1, initialSeed: -1}, nil, true},
		{"count < 1", args{mathematicalExpression: "1", count: -1, initialSeed: -1}, nil, true},
		{"bad symbol", args{mathematicalExpression: "$", count: 0, initialSeed: -1}, nil, true},
		{"1 count = 0", args{mathematicalExpression: "1", count: 0, initialSeed: -1}, nil, true},
		{"1 count = 0", args{mathematicalExpression: "1", count: 1, initialSeed: -1}, SpecMulti{
			EquationPairing{Independents: map[string]float64{"x": -1}, Dependent: 1},
		}, false},
		{"1 count = 0", args{mathematicalExpression: "x", count: 1, initialSeed: -1}, SpecMulti{
			EquationPairing{Independents: map[string]float64{"x": -1}, Dependent: -1},
		}, false},
		{"1 count = 0", args{mathematicalExpression: "x", count: 2, initialSeed: -1}, SpecMulti{
			EquationPairing{Independents: map[string]float64{"x": -1}, Dependent: -1},
			EquationPairing{Independents: map[string]float64{"x": 0}, Dependent: -0},
		}, false},
		{"1 count = 0", args{mathematicalExpression: "x * x", count: 3, initialSeed: -1}, SpecMulti{
			EquationPairing{Independents: map[string]float64{"x": -1}, Dependent: 1},
			EquationPairing{Independents: map[string]float64{"x": 0}, Dependent: 0},
			EquationPairing{Independents: map[string]float64{"x": 1}, Dependent: 1},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateSpecSimple(tt.args.mathematicalExpression, tt.args.count, tt.args.initialSeed,
				tt.args.antagonistTresholdMultiplier, tt.args.protagonistTresholdMultiplier)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSpecSimple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateSpecSimple() = %v, want %v", got, tt.want)
			}
		})
	}
}
