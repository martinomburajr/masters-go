package evolution

import (
	"reflect"
	"testing"
)

func TestCrossover(t *testing.T) {
	type args struct {
		individual1 *Individual
		individual2 *Individual
		maxDepth    int
	}
	tests := []struct {
		name       string
		args       args
		wantChild1 Individual
		wantChild2 Individual
		wantErr    bool
	}{
		{"nil individual1", args{nil, nil, 0}, Individual{}, Individual{}, true},
		{"nil individual1.prog", args{&Individual{Program:nil}, nil, 0}, Individual{}, Individual{}, true},
		{"nil individual1.prog.T", args{&Individual{Program:&Program{T:nil}}, nil, 0}, Individual{}, Individual{}, true},
		{"nil individual1.prog.T.root", args{&Individual{Program:&Program{T:TreeNil()}}, nil, 0}, Individual{}, Individual{},
			true},
		{"nil individual2", args{&IndividualProg1Kind1, nil, 0}, Individual{}, Individual{},
			true},
		{"nil individual2.prog", args{&IndividualProg1Kind1, &Individual{Program:nil}, 0}, Individual{},
			Individual{},
			true},
		{"nil individual2.prog.T", args{&IndividualProg1Kind1, &Individual{Program:&Program{T:nil}}, 0}, Individual{},
			Individual{},
			true},
		{"nil individual2.prog.T.root", args{&IndividualProg1Kind1, &Individual{Program:&Program{T:TreeNil()}}, 0}, Individual{},
			Individual{},
			true},
		{"depth < 0", args{&IndividualProg1Kind1, &IndividualProg1Kind1, -1},  Individual{},
			Individual{},
			true},
		{"depth < 0", args{&IndividualProg0Kind1, &IndividualProg1Kind1, 0},  Individual{},
			Individual{},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChild1, gotChild2, err := Crossover(tt.args.individual1, tt.args.individual2, tt.args.maxDepth)
			if (err != nil) != tt.wantErr {
				t.Errorf("Crossover() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotChild1, tt.wantChild1) {
				t.Errorf("Crossover() gotChild1 = %v, want %v", gotChild1, tt.wantChild1)
			}
			if !reflect.DeepEqual(gotChild2, tt.wantChild2) {
				t.Errorf("Crossover() gotChild2 = %v, want %v", gotChild2, tt.wantChild2)
			}
		})
	}
}
