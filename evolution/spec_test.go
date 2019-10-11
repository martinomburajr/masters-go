package evolution

import (
	"reflect"
	"testing"
)

func TestGenerateSpec(t *testing.T) {
	type args struct {
		mathematicalExpression string
		count                  int
		initialSeed            int
	}
	tests := []struct {
		name    string
		args    args
		want    SpecMulti
		wantErr bool
	}{
		{"empty", args{mathematicalExpression: "", count: 0, initialSeed: 0}, SpecMulti{}, true},
		{"negative count", args{mathematicalExpression: "x", count: -1, initialSeed: 0}, SpecMulti{}, true},
		{"0 countt", args{mathematicalExpression: "1", count: 0, initialSeed: 0}, SpecMulti{}, true},
		{"bad string", args{mathematicalExpression: "$", count: 1, initialSeed: 0},
			SpecMulti{}, true},
		{"1", args{mathematicalExpression: "1", count: 1, initialSeed: 0}, SpecMulti{EquationPairing{Independent: 0,
			Dependent: 1}}, false},
		{"1", args{mathematicalExpression: "1", count: 1, initialSeed: 1}, SpecMulti{EquationPairing{Independent: 1,
			Dependent: 1}}, false},
		{"x", args{mathematicalExpression: "x", count: 1, initialSeed: 1}, SpecMulti{EquationPairing{Independent: 1,
			Dependent: 1}}, false},
		{"x", args{mathematicalExpression: "x", count: 2, initialSeed: 1},
			SpecMulti{
				EquationPairing{Independent: 1, Dependent: 1},
				EquationPairing{Independent: 2, Dependent: 2},
			}, false,
		},
		{"x * x", args{mathematicalExpression: "x * x", count: 3, initialSeed: 1},
			SpecMulti{
				EquationPairing{Independent: 1, Dependent: 1},
				EquationPairing{Independent: 2, Dependent: 4},
				EquationPairing{Independent: 3, Dependent: 9},
			}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateSpec(tt.args.mathematicalExpression, tt.args.count, tt.args.initialSeed)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSpec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}
