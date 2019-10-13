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
		{"spec0", args{Spec0, &Prog0, 1,
			0.001}, -1, 1,
			false},
		//{"specX", args{SpecX, &Prog0, 1,
		//	0.001}, 1, -1,
		//	false},
		//{"specXX", args{SpecXX, &ProgTreeT_NT_T_4, 0.1,
		//	0.001}, 1, -1,
		//	false},
		//{"specXXXX", args{SpecXXXX, &ProgTreeTreeT_NT_T_NT_T_NT_T_2, 0.01,
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
