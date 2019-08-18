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
