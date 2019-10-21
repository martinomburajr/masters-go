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
		individuals       []*Individual
		fitnessComparator int
	}
	tests := []struct {
		name    string
		args    args
		want    *Individual
		wantErr bool
	}{
		{"nil", args{nil, 0}, nil, true},
		{"empty", args{[]*Individual{}, 0}, nil, true},
		{"1", args{[]*Individual{i1}, 0}, i1, false},
		{"2", args{[]*Individual{i1, i2}, 0}, i1, false},
		{"3", args{[]*Individual{i1, i2, iMax}, 0}, i1, false},
		{"4", args{[]*Individual{i1, i2, iMax, imin1}, 0}, imin1, false},
		{"4", args{[]*Individual{imin1, i2, iMax, i1}, 0}, imin1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNthPlaceIndividual(tt.args.individuals, tt.args.fitnessComparator)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNthPlaceIndividual() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got.TotalFitness, tt.want.TotalFitness) {
					t.Errorf("GetNthPlaceIndividual() = %v, want %v", got.TotalFitness, tt.want.TotalFitness)
				}
			}
		})
	}
}

var g0Pro = &Generation{GenerationID: "g0Pro", Protagonists: []*Individual{i1}}
var g0Ant = &Generation{GenerationID: "g0Pro", Antagonists: []*Individual{i1}}
var g0 = &Generation{GenerationID: "g0", Antagonists: []*Individual{i1}, Protagonists:[]*Individual{i1}}

var g1Pro = &Generation{GenerationID: "g1Pro", Protagonists: []*Individual{i1, i2}}
var g1 = &Generation{GenerationID: "g1",
	Protagonists: []*Individual{i2, i1},
	Antagonists: []*Individual{i1, iMax},}
var g1SortedMoreBetter = &Generation{GenerationID: "g1",
	Protagonists: []*Individual{i2, i1},
	Antagonists: []*Individual{iMax, i1}}
var g1SortedLessBetter = &Generation{GenerationID: "g1",
	Protagonists: []*Individual{i1, i2},
	Antagonists: []*Individual{i1, iMax}}

var g2Pro = &Generation{GenerationID: "g2Pro", Protagonists: []*Individual{iMax, i1}}
var g2 = &Generation{GenerationID: "g2Pro",
		Protagonists: []*Individual{iMax, i1},
		Antagonists: []*Individual{i1, iMax},}
var g2SortedMoreBetter = &Generation{GenerationID: "g2",
	Protagonists: []*Individual{iMax, i1},
	Antagonists: []*Individual{iMax, i1}}
var g2SortedLessBetter = &Generation{GenerationID: "g2",
	Protagonists: []*Individual{i1, iMax},
	Antagonists: []*Individual{i1, iMax}}

var g4Pro = &Generation{GenerationID: "g4Pro", Protagonists: []*Individual{iMax, i2}}
var g3Pro = &Generation{GenerationID: "g3Pro", Protagonists: []*Individual{imin1, iMax}}
//
//func TestCalcTopIndividualAllGenerations(t *testing.T) {
//	type args struct {
//		generations       []*Generation
//		individualKind    int
//		fitnessComparator int
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    ResultTopIndividuals
//		wantErr bool
//	}{
//		{"nil", args{nil, 1, 0}, ResultTopIndividuals{}, true},
//		{"empty", args{[]*Generation{}, 1, 0}, ResultTopIndividuals{}, true},
//		{"1", args{[]*Generation{g0Pro}, 1, 0}, ResultTopIndividuals{Individual: i1, Generation: g0Pro, Tree: ""}, false},
//		{"2", args{[]*Generation{g0Pro, g1Pro}, 1, 0}, ResultTopIndividuals{Individual: i1, Generation: g1Pro, Tree: ""}, false},
//		{"3", args{[]*Generation{g0Pro, g1Pro, g2Pro}, 1, 0}, ResultTopIndividuals{Individual: i1, Generation: g2Pro, Tree: ""}, false},
//		{"3", args{[]*Generation{g0Pro, g1Pro, g2Pro, g3Pro}, 1, 0}, ResultTopIndividuals{Individual: imin1, Generation: g3Pro, Tree: ""}, false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := CalcNthPlaceIndividualAllGenerations(tt.args.generations, tt.args.individualKind, tt.args.fitnessComparator)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("CalcNthPlaceIndividualAllGenerations() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("CalcNthPlaceIndividualAllGenerations() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestSortIndividuals(t *testing.T) {
//	type args struct {
//		individuals       []*Individual
//		fitnessComparator int
//	}
//	tests := []struct {
//		name string
//		args args
//		want []*Individual
//	}{
//		{"", args{[]*Individual{i1, i2, iMax, imin1}, 0}, []*Individual{imin1, i1, i2, iMax}},
//		{"", args{[]*Individual{i1, i2, iMax, i1, imin1}, 0}, []*Individual{imin1, i1, i1, i2, iMax}},
//		{"", args{[]*Individual{i1, i2, iMax, i1, imin1, i1}, 0}, []*Individual{imin1, i1, i1, i1, i2, iMax}},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got := SortIndividuals(tt.args.individuals, tt.args.fitnessComparator)
//			if len(got) != len(tt.want) {
//				t.Errorf("not same lengthgot = %v, want %v", got, tt.want)
//			}
//			for i := range got {
//				if got[i].TotalFitness != tt.want[i].TotalFitness {
//					t.Errorf("got = %v, want %v", got[i].TotalFitness, tt.want[i].TotalFitness)
//					return
//				}
//			}
//		})
//	}
//}
