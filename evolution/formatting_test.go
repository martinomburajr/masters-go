package evolution

import (
	"reflect"
	"testing"
)

func TestFormatStrategiesTotal(t *testing.T) {
	type args struct {
		strategies []Strategy
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		//{"", args{[]Strategy{"1"}}, "---- 1: 1\n"},
		//{"", args{[]Strategy{"1", "1"}}, "---- 1: 2\n"},
		//{"", args{[]Strategy{"1", "2"}}, "---- 1: 1\n---- 2: 1\n"},
		//{"", args{[]Strategy{"1", "1", "2"}}, "---- 1: 2\n---- 2: 1\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatStrategiesTotal(tt.args.strategies); !reflect.DeepEqual(got.String(), tt.want) {
				t.Errorf("FormatStrategiesTotal() =\n%v, want %v", got.String(), tt.want)
			}
		})
	}
}
