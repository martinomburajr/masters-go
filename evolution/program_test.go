package evolution

import (
	"testing"
)

func TestProgram_Eval(t *testing.T) {
	tests := []struct {
		name    string
		Program *Program
		args    float32
		want    float32
		wantErr bool
	}{
		{"nil", &Program{}, 0, -1, true},
		{"nil", &Program{T: TreeNil()}, 0, -1, true},
		{"T", &Program{T: TreeT_0()}, 5, 5, false},
		{"T", &Program{T: Tree5()}, 5, -1, true},
		{"T-NT-T", &Program{T: TreeT_NT_T_0()}, 5, 20, false},
		{"T-NT-T-NT-T", &Program{T: TreeT_NT_T_NT_T_0()}, 5, -15, false},
		{"T-NT-T", &Program{T: Tree8()}, 7, 49, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.Program.Eval(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Program.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Program.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProgram_ApplyStrategy(t *testing.T) {
	type fields struct {
		ID    string
		T     *DualTree
		Depth int
	}
	type args struct {
		strategy                       Strategy
		terminals                      []SymbolicExpression
		nonTerminals                   []SymbolicExpression
		mutationProbability            float32
		nonTerminalMutationProbability float32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Program{
				ID:    tt.fields.ID,
				T:     tt.fields.T,
				Depth: tt.fields.Depth,
			}
			if err := p.ApplyStrategy(tt.args.strategy, tt.args.terminals, tt.args.nonTerminals, tt.args.mutationProbability, tt.args.nonTerminalMutationProbability); (err != nil) != tt.wantErr {
				t.Errorf("Program.ApplyStrategy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
