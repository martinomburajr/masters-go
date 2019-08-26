package evolution

import (
	"reflect"
	"testing"
)

func TestEpochSimulator_Start(t *testing.T) {
	//epochNil := &Epoch{antagonist: nil, protagonist: nil, program: nil, protagonistBegins: false}
	//epochNilProtaProg := &Epoch{antagonist: nil, protagonist: nil, program: nil, protagonistBegins: false}
	//epochNilPr := &Epoch{antagonist: nil, protagonist: nil, program: p, protagonistBegins: false}
	tests := []struct {
		name    string
		fields  *EpochSimulator
		want    *EpochResult
		wantErr bool
	}{
		//{"nil"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Start()
			if (err != nil) != tt.wantErr {
				t.Errorf("EpochSimulator.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EpochSimulator.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEpochSimulator_ApplyAntagonistStrategy(t *testing.T) {
	type fields struct {
		epoch                 *Epoch
		hasAntagonistApplied  bool
		hasProtagonistApplied bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    *EpochSimulator
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EpochSimulator{
				epoch:                 tt.fields.epoch,
				hasAntagonistApplied:  tt.fields.hasAntagonistApplied,
				hasProtagonistApplied: tt.fields.hasProtagonistApplied,
			}
			got, err := e.applyAntagonistStrategy()
			if (err != nil) != tt.wantErr {
				t.Errorf("EpochSimulator.applyAntagonistStrategy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EpochSimulator.applyAntagonistStrategy() = %v, want %v", got, tt.want)
			}
		})
	}
}
