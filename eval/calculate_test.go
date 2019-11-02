package eval

import (
	"testing"
)

func TestCalculate(t *testing.T) {
	type args struct {
		substitutedExpression string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{"empty", args{""}, 0.0, true},
		{"foreign symbol", args{"r"}, 0.0, true},
		{"valid number | 7", args{"7"}, 7, false},
		{"valid number | 73213", args{"73213"}, 73213, false},
		{"5+4+3+2+1", args{"5+4+3+2+1"}, 15, false},
		{"5+4+3+2+ 1 spacing", args{"5+4+3+2+ 1"}, 0.0, true},
		{"(5+4)", args{"(5+4)"}, 9, false},
		{"(5+4", args{"(5+4"}, 9, false},
		{"5+4)", args{"5+4)"}, 9, false},
		{"(5+(4+(3+(2+ 1))))", args{"(5+(4+(3+(2+ 1))))"}, 0.0, true},
		{"(4/2)", args{"(4/2)"}, 2, false},
		{"-1+2", args{"-1+2"}, 1, false},
		{"((((6)*(2))-((5)+((5)+(6))))*(1))", args{"((((6)*(2))-((5)+((5)+(6))))*(1))"}, -4, false},
		{"(((6)*(((8)+(2))-(5)))+(((1)+(9))-(((1)*(2))*((1)+(6)))))", args{"(((6)*(((8)+(2))-(5)))+(((1)+(9))-(((1)*(" +
			"2))*((1)+(6)))))"}, 26, false},
		{"(4-2)", args{"(4-2)"}, 2, false},
		{"(4-9)", args{"(4-9)"}, -5, false},
		{"(4--9)", args{"(4--9)"}, 13, false},
		{"((4)-(-9))", args{"((4)-(-9))"}, 13, false},
		{"(4-(-9))", args{"(4-(-9))"}, 13, false},
		{"0-(0-50)+(2-0-2)", args{"0-(0-50)+(2-0-2)"}, 50, false},
		{"1.23*-1.34+2.3442", args{"1.23*-1.34+2.3442"}, 0.6959999999999997, false},
		{"(1.23)*(-1.34)+(2.3442)", args{"(1.23)*(-1.34)+(2.3442)"}, 0.6959999999999997, false},
		{"1.23*1.34+2.3442", args{"1.23*1.34+2.3442"}, 3.9924, false},
		{"(1.23)*(1.34)+(2.3442)", args{"(1.23)*(1.34)+(2.3442)"}, 3.9924, false},
		{"((1.23)*(1.34)+(2.3442))", args{"((1.23)*(1.34)+(2.3442))"}, 3.9924, false},
		{"((1.23)*(1.34)+(2.3442))", args{"((1.23)*(1.34)+(2.3442)*(1))"}, 3.9924, false},
		{"((1.23)*(1.34)+(2.3442))", args{"((1.23)*(1.34)+(2.3442)*(1)*(-1))"}, -0.6959999999999997, false},
		{"((((6)*(2))-((5)+((5)+(6))))*(0))", args{"((((6)*(2))-((5)+((5)+(6))))*(0))"}, 0, false},
		{"(((6)*(((8)+(2))-(5)))+(((0)+(9))-(((0)*(2))*((0)+(6)))))", args{"(((6)*(((8)+(2))-(5)))+(((0)+(9))-(((0)*(" +
			"2))*((0)+(6)))))"}, 39, false},
		{"1000*1000*1000", args{"1000*1000*1000"}, 1000000000, false},
		{"1000/0", args{"1000/0"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.args.substitutedExpression)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateWithVar(t *testing.T) {
	type args struct {
		substitutedExpression string
		variables             map[string]float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{"empty", args{"", nil}, 0.0, true},
		{"1", args{"1", nil}, 1, false},
		{"x without map", args{"x", nil}, 0.0, true},
		{"x with map", args{"x", map[string]float64{"x": 1}}, 1, false},
		{"x + y (y no map)", args{"x + y", map[string]float64{"x": 1}}, 1, false},
		{"x + y", args{"x+y", map[string]float64{"x": 1, "y": 2}}, 3, false},
		{"x * y", args{"x*y", map[string]float64{"x": 1, "y": 2}}, 2, false},
		{"x * y * a* b", args{"x*y*a*b", map[string]float64{"x": 1, "y": 2, "a": 3, "b": 4}}, 24, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateWithVar(tt.args.substitutedExpression, tt.args.variables)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateWithVar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CalculateWithVar() = %v, want %v", got, tt.want)
			}
		})
	}
}
