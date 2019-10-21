package evolution

import (
	"testing"
)

func TestSortGenerationsThoroughly(t *testing.T) {
	type args struct {
		generations         []*Generation
		isMoreFitnessBetter bool
	}
	tests := []struct {
		name    string
		args    args
		want    []*Generation
		wantErr bool
	}{
		{"nil generations", args{nil, true}, nil, true},
		{"empty generations", args{nil, true}, nil, true},
		{"Gen No antagonists", args{[]*Generation{g0Pro}, true}, nil, true},
		{"Gen No protagonists", args{[]*Generation{g0Ant}, true}, nil, true},
		{"Gen Count 1", args{[]*Generation{g0}, true}, []*Generation{g0}, false},
		{"Gen Count 1", args{[]*Generation{g1}, true}, []*Generation{g1SortedMoreBetter}, false},
		{"Gen Count 1 Less is Better", args{[]*Generation{g1}, false}, []*Generation{g1SortedLessBetter}, false},
		{"Gen Count 2 Less is Better", args{[]*Generation{g1, g2}, false}, []*Generation{g1SortedLessBetter, g2SortedLessBetter},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SortGenerationsThoroughly(tt.args.generations, tt.args.isMoreFitnessBetter)
			if (err != nil) != tt.wantErr {
				t.Errorf("SortGenerationsThoroughly() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(tt.want) != len(got) {
					t.Errorf("SortGenerationsThoroughly() = %v, want %v", got, tt.want)
				}
				for i := 0; i < len(got); i++ {
					if len(tt.want[i].Antagonists) != len(got[i].Antagonists) {
						t.Errorf("SortGenerationsThoroughly() Antagonists not Same Length = %v, want %v", got, tt.want)
					}
					if len(tt.want[i].Protagonists) != len(got[i].Protagonists) {
						t.Errorf("SortGenerationsThoroughly() Protagonists not Same Length = %v, want %v", got, tt.want)
					}
					for e := range got[i].Antagonists {
						if got[i].Antagonists[e].TotalFitness != tt.want[i].Antagonists[e].TotalFitness {
							t.Errorf("SortGenerationsThoroughly() Antagonists = %v, want %v", got, tt.want)
						}
					}
					for e := range got[i].Protagonists {
						if got[i].Protagonists[e].TotalFitness != tt.want[i].Protagonists[e].TotalFitness {
							t.Errorf("SortGenerationsThoroughly() Protagonists = %v, want %v", got, tt.want)
						}
					}
				}
			}
		})
	}
}
