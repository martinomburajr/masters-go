package evolution

import (
	"testing"
)

func TestEpoch_Start(t *testing.T) {
	tests := []struct {
		name    string
		fields  *Epoch
		wantErr bool
	}{
		{"nil-protagonist", &Epoch{}, true},
		//{"nil-antagonist", &Epoch{protagonist: &IndividualProg0Kind0}, true},
		//{"err-applyAntagonistStrategy", &Epoch{protagonist: &IndividualProg0Kind0, antagonist: &IndividualProg0Kind0},
		//	true},
		//{"err-applyAntagonistStrategy", &Epoch{protagonist: &IndividualProg0Kind0, antagonist: &IndividualProg0Kind0},
		//	true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Epoch{
				id:                               tt.fields.id,
				protagonist:                      tt.fields.protagonist,
				antagonist:                       tt.fields.antagonist,
				generation:                       tt.fields.generation,
				program:                          tt.fields.program,
				protagonistBegins:                tt.fields.protagonistBegins,
				isComplete:                       tt.fields.isComplete,
				probabilityOfMutation:            tt.fields.probabilityOfMutation,
				probabilityOfNonTerminalMutation: tt.fields.probabilityOfNonTerminalMutation,
				terminalSet:                      tt.fields.terminalSet,
				nonTerminalSet:                   tt.fields.nonTerminalSet,
				hasAntagonistApplied:             tt.fields.hasAntagonistApplied,
				hasProtagonistApplied:            tt.fields.hasProtagonistApplied,
			}
			if err := e.Start(); (err != nil) != tt.wantErr {
				t.Errorf("Epoch.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEpoch_applyProtagonistStrategy(t *testing.T) {
	tests := []struct {
		name    string
		fields  *Epoch
		wantErr bool
	}{
		{"nil-strategy", &EpochNil, true},
		//{"nil-strategy", &Epoch0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Epoch{
				id:                               tt.fields.id,
				protagonist:                      tt.fields.protagonist,
				antagonist:                       tt.fields.antagonist,
				generation:                       tt.fields.generation,
				program:                          tt.fields.program,
				protagonistBegins:                tt.fields.protagonistBegins,
				isComplete:                       tt.fields.isComplete,
				probabilityOfMutation:            tt.fields.probabilityOfMutation,
				probabilityOfNonTerminalMutation: tt.fields.probabilityOfNonTerminalMutation,
				terminalSet:                      tt.fields.terminalSet,
				nonTerminalSet:                   tt.fields.nonTerminalSet,
				hasAntagonistApplied:             tt.fields.hasAntagonistApplied,
				hasProtagonistApplied:            tt.fields.hasProtagonistApplied,
			}
			if err := e.applyProtagonistStrategy(); (err != nil) != tt.wantErr {
				t.Errorf("Epoch.applyProtagonistStrategy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
