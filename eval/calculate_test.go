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
		{"5+4+3+2+ 1 spacing", args{"5+4+3+2+ 1"}, 15, false},
		{"(5+4)", args{"(5+4)"}, 9, false},
		{"(5+4", args{"(5+4"}, 9, false},
		{"5+4)", args{"5+4)"}, 9, false},
		{"(5+(4+(3+(2+ 1))))", args{"(5+(4+(3+(2+ 1))))"}, 15, false},
		{"(4/2)", args{"(4/2)"}, 2, false},
		{"(4-2)", args{"(4-2)"}, 2, false},
		{"(4-9)", args{"(4-9)"}, -5, false},
		{"(4--9)", args{"(4--9)"}, 13, false},
		{"((4)-(-9))", args{"((4) - (-9))"}, 13, false},
		{"(4-(-9))", args{"(4-(-9))"}, 13, false},
		{"0-(0-50)+(2-0-2)", args{"0-(0-50)+(2-0-2)"}, 50, false},
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
		{"x + y (y no map)", args{"x + y", map[string]float64{"x": 1}}, 0, true},
		{"x + y", args{"x + y", map[string]float64{"x": 1, "y": 2}}, 3, false},
		{"x * y", args{"x * y", map[string]float64{"x": 1, "y": 2}}, 2, false},
		{"x * y * a* b", args{"x * y *a *b", map[string]float64{"x": 1, "y": 2, "a": 3, "b": 4}}, 24, false},
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

func TestNegativeNumberParser(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{""}, ""},
		{"$", args{"$"}, "$"},
		{"1", args{"1"}, "1"},
		{"1+2", args{"1+2"}, "1+2"},
		{"(1+2)", args{"(1+2)"}, "(1+2)"},
		{"(1+2--3)", args{"(1+2--3)"}, "(1+2+3)"},
		{"-1", args{"-1"}, "0-1"},
		{"-1+2", args{"-1+2"}, "0-1+2"},
		{"1-1+2", args{"1-1+2"}, "1-1+2"},
		{"0--1+2", args{"0--1+2"}, "0+1+2"},
		{"0--1+2--2", args{"0--1+2--2"}, "0+1+2+2"},
		{"0--1+2--2--3", args{"0--1+2--2--3"}, "0+1+2+2+3"},
		{"-0--1+2--2--3", args{"-0--1+2--2--3"}, "0-0+1+2+2+3"},
		{"(-50)+(2)", args{"(-50)+(2)"}, "(0-50)+(2)"},
		{"-(-50)+(2)", args{"-(-50)+(2)"}, "0-(0-50)+(2)"},
		{"-(-50)+(2--2)", args{"-(-50)+(2--2)"}, "0-(0-50)+(2+2)"},
		{"(4-(-9))", args{"(4-(-9))"}, "(4-(0-9))"},
		{"(4-(-3-9))", args{"(4-(-3-9))"}, "(4-(0-3-9))"},
		{"(4-(-3-9)-(-9))", args{"(4-(-3-9)-(-9))"}, "(4-(0-3-9)-(0-9)))"},
		{"-50*-50", args{"-50*-50"}, "0-50*(0-50)"},
		{"-50*50", args{"-50*50"}, "0-50*50"},
		{"-50+-50", args{"-50+-50"}, "0-50+(0-50)"},
		{"-50/-50", args{"-50/-50"}, "0-50/(0-50)"},
		//{"(4-(-3-9)-(-9)*-2)", args{"(4-(-3-9)-(-9)*-2)"}, "(4-(0-3-9)-(0-9)*(0-2))"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NegativeNumberParser(tt.args.str)
			if got != tt.want {
				t.Errorf("NegativeNumberParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMartinsReplace(t *testing.T) {
	type args struct {
		str string
		old string
		new string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"x * x * x", args{"x * x * x", " ", ""}, "x*x*x"},
		{"x * x *                                               x", args{"x * x *                                               x", " ", ""}, "x*x*x"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MartinsReplace(tt.args.str, tt.args.old, tt.args.new); got != tt.want {
				t.Errorf("MartinsReplace() = %v, want %v", got, tt.want)
			}
		})
	}
}
