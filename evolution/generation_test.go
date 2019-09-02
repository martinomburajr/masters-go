package evolution

import (
	"testing"
)

func TestGeneration_setupEpochs(t *testing.T) {
	tests := []struct {
		name    string
		fields  *Generation
		want    []*Epoch
		wantErr bool
	}{
		{"nil-antagonists", &Generation{}, nil, true},
		{"nil-protagonists", &Generation{Antagonists: []*Individual{&IndividualProg1Kind1}}, nil, true},
		{"empty-antagonists", &Generation{Antagonists: []*Individual{}, Protagonists: []*Individual{}}, nil, true},
		{"empty-protagonists", &Generation{Antagonists: []*Individual{&IndividualProg1Kind1}, Protagonists: []*Individual{}}, nil, true},
		{"one-of-each", &GenerationTest0, []*Epoch{{}, {}, {}, {}}, false},
		{"1-less-protagonist", &GenerationTest1, []*Epoch{{}, {}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.setupEpochs()
			if (err != nil) != tt.wantErr {
				t.Errorf("Generation.setupEpochs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("Generation.setupEpochs() = %v, want %v", got, tt.want)
			}
		})
	}
}
