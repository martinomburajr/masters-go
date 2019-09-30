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

func TestTournamentSelection(t *testing.T) {
	type args struct {
		population     []*Individual
		tournamentSize int
	}
	tests := []struct {
		name    string
		args    args
		want    []*Individual
		wantErr bool
	}{
		{"nil-population", args{population: nil, tournamentSize: 0}, nil, true},
		{"empty population", args{population: []*Individual{}, tournamentSize: 0}, nil, true},
		{"tournamentSize=0", args{population: []*Individual{&IndividualProg0Kind1}, tournamentSize: 0}, nil, true},
		{"tournamentSize=0", args{population: []*Individual{&IndividualProg0Kind1}, tournamentSize: 0}, nil, true},
		{"pop=1 tournamentSize=1", args{population: []*Individual{&IndividualProg0Kind1}, tournamentSize: 1},
			[]*Individual{&IndividualProg0Kind1}, false},
		{"pop=2 tournamentSize=1", args{population: []*Individual{&IndividualProg0Kind1, &IndividualProg1Kind1},
			tournamentSize: 1},
			[]*Individual{&IndividualProg0Kind1, &IndividualProg1Kind1}, false},
		{"pop=2 tournamentSize=2", args{population: []*Individual{&IndividualProg0Kind1, &IndividualProg1Kind1},
			tournamentSize: 2},
			[]*Individual{&IndividualProg0Kind1, &IndividualProg1Kind1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TournamentSelection(tt.args.population, tt.args.tournamentSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("TournamentSelection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("TournamentSelection() = %v, want %v", got, tt.want)
			}
		})
	}
}
