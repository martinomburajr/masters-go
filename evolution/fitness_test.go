package evolution

import (
	"reflect"
	"testing"
)

func TestProtagonistThresholdTally(t *testing.T) {
	type args struct {
		spec                             SpecMulti
		protagonistAntagonistProgramPair *Program
		threshold                        float64
		minthreshold                     float64
	}
	tests := []struct {
		name            string
		args            args
		wantAntagonist  int
		wantProtagonist int
		wantErr         bool
	}{
		{"nil-spec", args{nil, nil, 0, 0}, 0, 0, true},
		{"nil-antagonist", args{Spec0, nil, 0, 0}, 0, 0, true},
		{"nil-protagonist", args{Spec0, nil, 0, 0}, 0, 0, true},
		{"nil-papair", args{Spec0, nil, 0, 0}, 0, 0, true},
		{"nil-minthreshold0", args{Spec0, &ProgNil, 0, 0}, 0, 0,
			true},
		{"nil-threshold<minthreshold", args{Spec0, &ProgNil, 0.0001,
			0.001}, 0, 0,
			true},
		{"empty-antagonist", args{Spec0, &ProgNil, 1,
			0.001}, 0, 0,
			true},
		{"empty-protagonist", args{Spec0, &ProgNil, 1,
			0.001}, 0, 0,
			true},
		{"empty-papair", args{Spec0, &ProgNil, 1,
			0.001}, 0, 0,
			true},
		{"spec0", args{Spec0, &Prog1, 1,
			0.001}, 1, -1,
			false},
		{"spec0", args{Spec0, &ProgX, 1,
			0.001}, -1, 1,
			false},
		//{"specX", args{SpecX, &ProgX, 1,
		//	0.001}, 1, -1,
		//	false},
		//{"specXX", args{SpecXX, &ProgTreeT_NT_T_4, 0.1,
		//	0.001}, 1, -1,
		//	false},
		//{"specXXXX", args{SpecXXXX, &ProgXXXX, 0.01,
		//	0.001}, 1, -1,
		//	false},
		//{"specXXXXAdd4 - small threshold", args{SpecXXXXAdd4, &ProgTreeXXXX4, 0.01,
		//	0.001}, -1, 1,
		//	false},
		//{"specXXXXAdd4 - large threshold", args{SpecXXXXAdd4, &ProgTreeXXXX4, 50,
		//	0.001}, 1, -1,
		//	false},
		//{"specXXXXAdd4", args{SpecXXXXAdd4, &ProgTreeXXXX4, 1,
		//	0.001}, 1, -1,
		//	false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ProtagonistThresholdTally(tt.args.spec,
				tt.args.protagonistAntagonistProgramPair, tt.args.threshold)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProtagonistThresholdTally() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantAntagonist) {
				t.Errorf("ProtagonistThresholdTally() got = %v, isEqual %v", got, tt.wantAntagonist)
			}
			if !reflect.DeepEqual(got1, tt.wantProtagonist) {
				t.Errorf("ProtagonistThresholdTally() got1 = %v, isEqual %v", got1, tt.wantProtagonist)
			}
		})
	}
}

func TestAggregateFitness(t *testing.T) {
	tests := []struct {
		name    string
		args    *Individual
		want    float64
		wantErr bool
	}{
		//{"nil Fitness", &Individual{}, math.MaxInt8, true},
		//{"empty Fitness", &Individual{Fitness: []int{}}, math.MaxInt8, true},
		//{"input | 1,2", &Individual{Fitness: []int{1, 2}}, 3, false},
		//{"input | 0", &Individual{Fitness: []int{0}}, 0, false},
		//{"input | -1,1", &Individual{Fitness: []int{-1, 1}}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AggregateFitness(*tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("AggregateFitness() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AggregateFitness() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestThresholdedRatioFitness(t *testing.T) {
	xx25Spec, _ := GenerateSpec("x*x", []string{"x"}, 5, -2, 2, 5)
	type args struct {
		spec        SpecMulti
		antagonist  *Program
		protagonist *Program
	}
	tests := []struct {
		name                   string
		args                   args
		wantAntagonistFitness  float64
		wantProtagonistFitness float64
		wantErr                bool
	}{
		{"perfect-protagonist", args{xx25Spec, &ProgTreeT_NT_T_0, &ProgTreeT_NT_T_4}, 0, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAntagonistFitness, gotProtagonistFitness, err := ThresholdedRatioFitness(tt.args.spec, tt.args.antagonist, tt.args.protagonist)
			if (err != nil) != tt.wantErr {
				t.Errorf("ThresholdedRatioFitness() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAntagonistFitness != tt.wantAntagonistFitness {
				t.Errorf("ThresholdedRatioFitness() gotAntagonistFitness = %v, want %v", gotAntagonistFitness, tt.wantAntagonistFitness)
			}
			if gotProtagonistFitness != tt.wantProtagonistFitness {
				t.Errorf("ThresholdedRatioFitness() gotProtagonistFitness = %v, want %v", gotProtagonistFitness, tt.wantProtagonistFitness)
			}
		})
	}
}

var xSpec, _ = GenerateSpecSimple("x", 5, -1*(5/2), 5, 2)
var xSpecMono, _ = GenerateSpecSimple("x", 5, -1*(5/2), 3, 3)
var xSpecCount10, _ = GenerateSpecSimple("x", 10, -1*(10/2), 5, 2)
var xSpecCount10Mono, _ = GenerateSpecSimple("x", 10, -1*(10/2), 3, 3)
var xSpecCount100, _ = GenerateSpecSimple("x", 100, -1*(10/2), 5, 2)
var xSpecCount100Mono, _ = GenerateSpecSimple("x", 100, -1*(10/2), 3, 3)
var xxSpecCount100, _ = GenerateSpecSimple("x*x", 100, -1*(10/2), 5, 2)
var xxSpecCount100Mono, _ = GenerateSpecSimple("x*x", 100, -1*(10/2), 3, 3)

func Test_evaluateFitnessThresholded(t *testing.T) {

	type args struct {
		spec        SpecMulti
		antagonist  *Program
		protagonist *Program
	}
	tests := []struct {
		name                   string
		args                   args
		wantAntagonistFitness  float64
		wantProtagonistFitness float64
		wantErr                bool
	}{
		{"bad antagonist shape", args{xSpec, &ProgNil, &ProgNil}, 10000.001, 10000.001, true},
		{"bad protagonist shape", args{xSpec, &ProgX, &ProgNil}, 10000.001, 10000.001, true},
		{"DualThreshold > Perfect Pro | Pefect Ant", args{xSpec, &ProgX, &ProgX},
			-1, 1, false},
		{"DualThreshold > Perfect Pro | Amazing Ant", args{xSpec, &ProgTreeXXXX4, &ProgX},
			0.545, 1, false},
		{"DualThreshold > Perfect Pro | Amazing Ant", args{xSpecCount10, &ProgTreeXXXX4, &ProgX},
			0.863, 1, false},
		{"DualThreshold > Bad Pro | Amazing Ant", args{xSpec, &ProgTreeXXXX4, &ProgTreeXXXX4},
			0.69, -0.090909, false},
		{"DualThreshold > Bad Pro | Amazing Ant", args{xSpecCount10, &ProgTreeXXXX4, &ProgTreeXXXX4},
			0.863, -0.945, false},
		{"DualThreshold > Bad Pro | Amazing Ant", args{xSpecCount10, &ProgTreeXXXX4, &ProgXXXX},
			0.863, -0.962, false},
		{"DualThreshold > Bad Pro | Amazing Ant", args{xSpecCount10, &ProgTreeXXXX4, &ProgXXXX},
			0.863, -0.962, false},
		{"DualThreshold > OK Pro | Bad Ant", args{xSpecCount10, &ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0,
			&ProgXXXX},
			0.561, -0.962, false},
		{"DualThreshold > OK Pro | Bad Ant", args{xSpec, &ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0,
			&ProgXXXX},
			0.785, -0.647, false},
		{"DualThreshold > OK Pro | Bad Ant", args{xSpecCount100, &ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0,
			&ProgXXXX},
			-0.87625, -0.999, false},
		{"DualThreshold > OK Pro | Bad Ant", args{xxSpecCount100, &ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0,
			&ProgXXXX},
			-0.801, -0.999, false},

		{"MonoThreshold > Perfect Pro | Pefect Ant", args{xSpecMono, &ProgX, &ProgX},
			-1, 1, false},
		{"MonoThreshold > Perfect Pro | Amazing Ant", args{xSpecMono, &ProgTreeXXXX4, &ProgX},
			0.727, 1, false},
		{"MonoThreshold > Perfect Pro | Amazing Ant", args{xSpecCount10Mono, &ProgTreeXXXX4, &ProgX},
			0.918, 1, false},
		{"MonoThreshold > Bad Pro | Amazing Ant", args{xSpec, &ProgTreeXXXX4, &ProgTreeXXXX4},
			0.545, -0.818, false},
		{"MonoThreshold > Bad Pro | Amazing Ant", args{xSpecCount10Mono, &ProgTreeXXXX4, &ProgTreeXXXX4},
			0.918, -0.918, false},
		{"MonoThreshold > Bad Pro | Amazing Ant", args{xSpecCount10Mono, &ProgTreeXXXX4, &ProgXXXX},
			0.918, -0.943, false},
		{"MonoThreshold > Bad Pro | Amazing Ant", args{xSpecCount10Mono, &ProgTreeXXXX4, &ProgXXXX},
			0.918, -0.943, false},
		{"MonoThreshold > OK Pro | Bad Ant", args{xSpecCount10Mono, &ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0,
			&ProgXXXX},
			0.763, -0.943, false},
		{"MonoThreshold > OK Pro | Bad Ant", args{xSpecMono, &ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0,
			&ProgXXXX},
			0.871, -0.470, false},
		{"MonoThreshold > OK Pro | Bad Ant", args{xSpecCount100Mono, &ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0,
			&ProgXXXX},
			-0.793, -0.999, false},
		{"MonoThreshold > OK Pro | Bad Ant", args{xxSpecCount100Mono, &ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0,
			&ProgXXXX},
			-0.669, -0.999, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAntagonistFitness, gotProtagonistFitness, err := evaluateFitnessThresholded(tt.args.spec, tt.args.antagonist, tt.args.protagonist)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateFitnessThresholded() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if floatEquals(gotAntagonistFitness, tt.wantAntagonistFitness) {
				t.Errorf("evaluateFitnessThresholded() gotAntagonistFitness = %v, want %v", gotAntagonistFitness, tt.wantAntagonistFitness)
			}
			if floatEquals(gotProtagonistFitness, tt.wantProtagonistFitness) {
				t.Errorf("evaluateFitnessThresholded() gotProtagonistFitness = %v, want %v", gotProtagonistFitness, tt.wantProtagonistFitness)
			}
		})
	}
}

var EPSILON float64 = 0.00000001

func floatEquals(a, b float64) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

func Test_fitnessParameterValidator(t *testing.T) {
	type args struct {
		spec        SpecMulti
		antagonist  *Program
		protagonist *Program
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"nil-spec", args{nil, nil, nil}, true},
		{"empty-spec", args{SpecMulti{}, nil, nil}, true},
		{"nil antagonist", args{SpecMulti{EquationPairing{}}, nil, nil}, true},
		{"nil antagonist tree", args{SpecMulti{EquationPairing{}}, &Program{}, nil}, true},
		{"nil antagonist tree root", args{SpecMulti{EquationPairing{}}, &ProgNil, nil}, true},
		{"nil protagonist", args{SpecMulti{EquationPairing{}}, &ProgX, nil}, true},
		{"nil protagonist tree", args{SpecMulti{EquationPairing{}}, &ProgX, &ProgNil}, true},
		{"nil protagonist tree root", args{SpecMulti{EquationPairing{}}, &ProgX, &Program{}}, true},
		{"bad antagonist tree ", args{SpecMulti{EquationPairing{}}, &ProgBadTree, &ProgX}, true},
		{"bad protagonist tree ", args{SpecMulti{EquationPairing{}}, &ProgX, &ProgBadTree}, true},
		{"validation OK", args{SpecMulti{EquationPairing{}}, &ProgX, &ProgX}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fitnessParameterValidator(tt.args.spec, tt.args.antagonist, tt.args.protagonist); (err != nil) != tt.wantErr {
				t.Errorf("fitnessParameterValidator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_generateExpressions(t *testing.T) {
	type args struct {
		antagonist  *Program
		protagonist *Program
	}
	tests := []struct {
		name                      string
		args                      args
		wantAntagonistExpression  string
		wantProtagonistExpression string
		wantErr                   bool
	}{
		{"antagonist nil root", args{&ProgNil, &ProgNil}, "", "", true},
		{"protagonist nil root", args{&ProgX, &ProgNil}, "", "", true},
		{"Valid", args{&ProgX, &ProgTreeT_NT_T_0}, "(x)", "((x)*(4))", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAntagonistExpression, gotProtagonistExpression, err := generateExpressions(tt.args.antagonist, tt.args.protagonist)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateExpressions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAntagonistExpression != tt.wantAntagonistExpression {
				t.Errorf("generateExpressions() gotAntagonistExpression = %v, want %v", gotAntagonistExpression, tt.wantAntagonistExpression)
			}
			if gotProtagonistExpression != tt.wantProtagonistExpression {
				t.Errorf("generateExpressions() gotProtagonistExpression = %v, want %v", gotProtagonistExpression, tt.wantProtagonistExpression)
			}
		})
	}
}

func Test_calculateDelta(t *testing.T) {
	type args struct {
		truth float64
		value float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"0|1", args{0, 1}, 1},
		{"0| -1", args{0, -1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateDelta(tt.args.truth, tt.args.value); got != tt.want {
				t.Errorf("calculateDelta() = %v, want %v", got, tt.want)
			}
		})
	}
}
