package evolution

import (
	"log"
	"testing"
)

func TestCrossover(t *testing.T) {
	type args struct {
		individual1 *Individual
		individual2 *Individual
		maxDepth    int
		params      EvolutionParams
	}
	tests := []struct {
		name       string
		args       args
		wantChild1 Individual
		wantChild2 Individual
		wantErr    bool
	}{
		//{"nil individual1", args{nil, nil, 0}, Individual{}, Individual{}, true},
		//{"nil individual1.prog", args{&Individual{Program:nil}, nil, 0}, Individual{}, Individual{}, true},
		//{"nil individual1.prog.T", args{&Individual{Program:&Program{T:nil}}, nil, 0}, Individual{}, Individual{}, true},
		//{"nil individual1.prog.T.root", args{&Individual{Program:&Program{T:TreeNil()}}, nil, 0}, Individual{}, Individual{},
		//	true},
		//{"nil individual2", args{&IndividualProg1Kind1, nil, 0}, Individual{}, Individual{},
		//	true},
		//{"nil individual2.prog", args{&IndividualProg1Kind1, &Individual{Program:nil}, 0}, Individual{},
		//	Individual{},
		//	true},
		//{"nil individual2.prog.T", args{&IndividualProg1Kind1, &Individual{Program:&Program{T:nil}}, 0}, Individual{},
		//	Individual{},
		//	true},
		//{"nil individual2.prog.T.root", args{&IndividualProg1Kind1, &Individual{Program:&Program{T:TreeNil()}}, 0}, Individual{},
		//	Individual{},
		//	true},
		//{"depth < 0", args{&IndividualProg1Kind1, &IndividualProg1Kind1, -1},  Individual{},
		//	Individual{},
		//	true},
		//{"depth == 0", args{&IndividualProg0Kind1, &IndividualProg1Kind1, 0},  IndividualProg0Kind1,
		//	IndividualProg1Kind1,
		//	false},
		//{"depth == 1", args{&IndividualProg0Kind1, &IndividualProgTreeT_NT_T_0, 1,
		//	EvolutionParams{DepthPenaltyStrategyPenalization: 3}},
		//	IndividualProg0Kind1,
		//	IndividualProg1Kind1,
		//	false},
		//{"depth == 1", args{&IndividualProgTreeT_NT_T_1, &IndividualProgTreeT_NT_T_0, 1,
		//	EvolutionParams{DepthPenaltyStrategyPenalization: 3}},
		//	IndividualProg0Kind1,
		//	IndividualProg1Kind1,
		//	false},
		{"depth == 2 & 1", args{&IndividualProgTreeT_NT_T_NT_T_0, &IndividualProgTreeT_NT_T_0, 2,
			EvolutionParams{DepthPenaltyStrategyPenalization: 3}},
			IndividualProg0Kind1,
			IndividualProg1Kind1,
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChild1, gotChild2, err := Crossover(tt.args.individual1, tt.args.individual2, tt.args.maxDepth, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Crossover() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//containsSubTree, err := gotChild1.Program.T.ContainsSubTree(tt.args.individual1.Program.T)
			//if err != nil {
			//	t.Errorf("Crossover() error = %v, wantErr %v", err, tt.wantErr)
			//}
			//if containsSubTree {
			//	t.Errorf("Crossover() | child1 is identical to parent1 gotChild1 = %v, want %v", gotChild1,
			//		tt.args.individual1)
			//}
			//
			//containsSubTree2, err := gotChild2.Program.T.ContainsSubTree(tt.args.individual2.Program.T)
			//if err != nil {
			//	t.Errorf("Crossover() error = %v, wantErr %v", err, tt.wantErr)
			//}
			//if containsSubTree2 {
			//	t.Errorf("Crossover() | child1 is identical to parent1 gotChild1 = %v, individual 2 %v", gotChild2,
			//		tt.args.individual2)
			//}

			if gotChild1.Program.T != nil && gotChild2.Program != nil {
				log.Print("Individual 1")
				tt.args.individual1.Program.T.Print()
				log.Print("Individual 2")
				tt.args.individual2.Program.T.Print()
				log.Print("Child 1")
				gotChild1.Program.T.Print()
				log.Print("Child 2")
				gotChild2.Program.T.Print()
			}

		})
	}
}

func Test_depthPenaltyIgnore(t *testing.T) {
	type args struct {
		maxDepth         int
		individual1Depth int
		individual2Depth int
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 int
	}{
		{"", args{3, 0, 0}, 3, 3},
		{"", args{3, 4, 1}, 0, 2},
		{"", args{0, 0, 0}, 0, 0},
		{"", args{-1, 0, 0}, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := depthPenaltyIgnore(tt.args.maxDepth, tt.args.individual1Depth, tt.args.individual2Depth)
			if got != tt.want {
				t.Errorf("depthPenaltyIgnore() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("depthPenaltyIgnore() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
