package utils

import (
	"reflect"
	"testing"
)

func TestPermutationsWithRepetitions(t *testing.T) {
	type args struct {
		values []int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{"1 item", args{values: []int{1}}, [][]int{{1}}},
		{"2 items", args{values: []int{0, 1}}, [][]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PermutationsWithRepetitions(tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PermutationsWithRepetitions() = %v, want %v", got, tt.want)
			}
		})
	}
}
