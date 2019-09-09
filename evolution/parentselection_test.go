package evolution

import (
	"testing"
)

func Test_getNRandom(t *testing.T) {
	type args struct {
		population     []*Individual
		tournamentSize int
	}
	tests := []struct {
		name string
		args args
		want []*Individual
	}{
		{"ok", args{[]*Individual{&IndividualProg0Kind1, &IndividualNilProgNil}, 1},
			[]*Individual{&IndividualProg0Kind1}},
		{"ok", args{[]*Individual{&IndividualProg0Kind1, &IndividualNilProgNil}, 2},
			[]*Individual{&IndividualProg0Kind1, &IndividualNilProgNil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNRandom(tt.args.population, tt.args.tournamentSize); len(tt.want) != len(got) {
				t.Errorf("getNRandom() = %v, want %v", got, tt.want)
			}
		})
	}
}
