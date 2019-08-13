package program

import (
	"testing"
)

func TestProgram_Eval(t *testing.T) {

	tests := []struct {
		name   string
		fields *Program
		want   float32
	}{
		//{"nil", },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Program{
				ID:                   tt.fields.ID,
				T:                    tt.fields.T,
				Strategies:           tt.fields.Strategies,
				hasAppliedStrategies: tt.fields.hasAppliedStrategies,
				generation:           tt.fields.generation,
			}
			if got := p.Eval(); got != tt.want {
				t.Errorf("Program.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
