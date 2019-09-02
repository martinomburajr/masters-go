package evolution

import (
	"math"
	"reflect"
	"testing"
)

func TestProtagonistThresholdTally(t *testing.T) {
	type args struct {
		spec                             Spec
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
		{"spec0", args{Spec0, &Prog0, 1,
			0.001}, -1, 1,
			false},
		{"specX", args{SpecX, &Prog0, 1,
			0.001}, 1, -1,
			false},
		{"specXX", args{SpecXX, &ProgTreeT_NT_T_4, 0.1,
			0.001}, 1, -1,
			false},
		{"specXXXX", args{SpecXXXX, &ProgTreeTreeT_NT_T_NT_T_NT_T_2, 0.01,
			0.001}, 1, -1,
			false},
		{"specXXXXAdd4 - small threshold", args{SpecXXXXAdd4, &ProgTreeT_NT_T_NT_T_NT_T_NT_T_1, 0.01,
			0.001}, -1, 1,
			false},
		{"specXXXXAdd4 - large threshold", args{SpecXXXXAdd4, &ProgTreeT_NT_T_NT_T_NT_T_NT_T_1, 50,
			0.001}, 1, -1,
			false},
		{"specXXXXAdd4", args{SpecXXXXAdd4, &ProgTreeT_NT_T_NT_T_NT_T_NT_T_1, 1,
			0.001}, 1, -1,
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ProtagonistThresholdTally(tt.args.spec,
				tt.args.protagonistAntagonistProgramPair, tt.args.threshold, tt.args.minthreshold)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProtagonistThresholdTally() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantAntagonist) {
				t.Errorf("ProtagonistThresholdTally() got = %v, wantAntagonist %v", got, tt.wantAntagonist)
			}
			if !reflect.DeepEqual(got1, tt.wantProtagonist) {
				t.Errorf("ProtagonistThresholdTally() got1 = %v, wantAntagonist %v", got1, tt.wantProtagonist)
			}
		})
	}
}

func TestAggregateFitness(t *testing.T) {
	tests := []struct {
		name    string
		args    *Individual
		want    int
		wantErr bool
	}{
		{"nil fitness", &Individual{}, math.MaxInt8, true},
		{"empty fitness", &Individual{fitness: []int{}}, math.MaxInt8, true},
		{"input | 1,2", &Individual{fitness: []int{1,2}}, 3, false},
		{"input | 0", &Individual{fitness: []int{0}}, 0, false},
		{"input | -1,1", &Individual{fitness: []int{-1,1}}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AggregateFitness(tt.args)
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
