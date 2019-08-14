package program

import (
	"github.com/martinomburajr/masters-go/program/tree/dualtree"
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
		{"nil", &Program{}, 0, -1, true },
		{"nil", &Program{T: dualtree.TreeNil()}, 0, -1, true },
		{"T", &Program{T: dualtree.Tree0()}, 5, 5, false },
		{"T", &Program{T: dualtree.Tree5()}, 5, -1, true },
		{"T-NT-T", &Program{T: dualtree.Tree1()}, 5, 20, false },
		{"T-NT-T-NT-T", &Program{T: dualtree.Tree2()}, 5, -15, false },
		{"T-NT-T", &Program{T: dualtree.Tree8()}, 7, 49, false },
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
