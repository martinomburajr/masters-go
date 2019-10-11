package evolution

import (
	"math"
	"reflect"
	"testing"
)

var i1 = &Individual{TotalFitness: 1, Id: "i1"}
var i2 = &Individual{TotalFitness: 2, Id: "i2"}
var iMax = &Individual{TotalFitness: math.MaxInt64, Id: "iMaxInt"}
var imin1 = &Individual{TotalFitness: -1, Id: "iMin1"}

func TestCalcTopIndividual(t *testing.T) {
	type args struct {
		individuals []*Individual
	}
	tests := []struct {
		name    string
		args    args
		want    *Individual
		wantErr bool
	}{
		{"nil", args{nil}, nil, true},
		{"empty", args{[]*Individual{}}, nil, true},
		{"1", args{[]*Individual{i1}}, i1, false},
		{"2", args{[]*Individual{i1, i2}}, i1, false},
		{"3", args{[]*Individual{i1, i2, iMax}}, i1, false},
		{"4", args{[]*Individual{i1, i2, iMax, imin1}}, imin1, false},
		{"4", args{[]*Individual{imin1, i2, iMax, i1}}, imin1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalcTopIndividual(tt.args.individuals)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalcTopIndividual() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got.TotalFitness, tt.want.TotalFitness) {
					t.Errorf("CalcTopIndividual() = %v, want %v", got.TotalFitness, tt.want.TotalFitness)
				}
			}
		})
	}
}

var g0 = &Generation{GenerationID: "g0", Protagonists: []*Individual{i1}}
var g1 = &Generation{GenerationID: "g1", Protagonists: []*Individual{i1, i2}}
var g2 = &Generation{GenerationID: "g2", Protagonists: []*Individual{iMax, i1}}
var g4 = &Generation{GenerationID: "g4", Protagonists: []*Individual{iMax, i2}}
var g3 = &Generation{GenerationID: "g3", Protagonists: []*Individual{imin1, iMax}}

func TestCalcTopIndividualAllGenerations(t *testing.T) {
	type args struct {
		generations    []*Generation
		individualKind int
	}
	tests := []struct {
		name    string
		args    args
		want    ResultTopIndividuals
		wantErr bool
	}{
		{"nil", args{nil, 1}, ResultTopIndividuals{}, true},
		{"empty", args{[]*Generation{}, 1}, ResultTopIndividuals{}, true},
		{"1", args{[]*Generation{g0}, 1}, ResultTopIndividuals{Result: i1, Generation: g0, Tree: ""}, false},
		{"2", args{[]*Generation{g0, g1}, 1}, ResultTopIndividuals{Result: i1, Generation: g1, Tree: ""}, false},
		{"3", args{[]*Generation{g0, g1, g2}, 1}, ResultTopIndividuals{Result: i1, Generation: g2, Tree: ""}, false},
		{"3", args{[]*Generation{g0, g1, g2, g3}, 1}, ResultTopIndividuals{Result: imin1, Generation: g3, Tree: ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalcTopIndividualAllGenerations(tt.args.generations, tt.args.individualKind)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalcTopIndividualAllGenerations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalcTopIndividualAllGenerations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortIndividuals(t *testing.T) {
	type args struct {
		individuals []*Individual
	}
	tests := []struct {
		name string
		args args
		want []*Individual
	}{
		{"", args{[]*Individual{i1, i2, iMax, imin1}}, []*Individual{imin1, i1, i2, iMax}},
		{"", args{[]*Individual{i1, i2, iMax, i1, imin1}}, []*Individual{imin1, i1, i1, i2, iMax}},
		{"", args{[]*Individual{i1, i2, iMax, i1, imin1, i1}}, []*Individual{imin1, i1, i1, i1, i2, iMax}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortIndividuals(tt.args.individuals)
			if len(got) != len(tt.want) {
				t.Errorf("not same lengthgot = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i].TotalFitness != tt.want[i].TotalFitness {
					t.Errorf("got = %v, want %v", got[i].TotalFitness, tt.want[i].TotalFitness)
					return
				}
			}
		})
	}
}
