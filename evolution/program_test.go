package evolution

import "testing"

//func TestProgram_Eval(t *testing.T) {
//	tests := []struct {
//		name    string
//		Program *Program
//		args    float32
//		want    float32
//		wantErr bool
//	}{
//		{"nil", &Program{}, 0, -1, true},
//		{"nil", &Program{T: TreeNil()}, 0, -1, true},
//		{"T", &Program{T: TreeT_X()}, 5, 5, false},
//		{"T", &Program{T: Tree5()}, 5, -1, true},
//		{"T-NT-T", &Program{T: TreeT_NT_T_0()}, 5, 20, false},
//		{"T-NT-T-NT-T", &Program{T: TreeT_NT_T_NT_T_0()}, 5, -15, false},
//		{"T-NT-T", &Program{T: Tree8()}, 7, 49, false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			expressionString, err := tt.Program.T.ToMathematicalString()
//			if tt.Program.T == nil && tt.wantErr {
//				t.Errorf("Program.Eval() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Program.Eval() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			got, err := tt.Program.Eval(tt.args, expressionString)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Program.Eval() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("Program.Eval() = %v, isEqual %v", got, tt.want)
//			}
//		})
//	}
//}

func TestEvaluateMathematicalExpression(t *testing.T) {
	type args struct {
		expressionString     string
		independentVariables IndependentVariableMap
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{"nil", args{"", nil}, -1, true},
		{"empty map", args{"x", map[string]float64{}}, -1, true},
		{"x -> 0", args{"x", map[string]float64{"x": 0}}, 0, false},
		{"x -> 0 | y -> 1", args{"x + y", map[string]float64{"x": 0, "y": 1}}, 1, false},
		{"x -> 2 | y -> 3", args{"(x * y) * y", map[string]float64{"x": 2, "y": 3}}, 18, false},
		{"1", args{"1", nil}, 1, false},
		//{"(4-(-3--9)-(-9)*-2*-3)", args{"(4-(-3--9)-(-9)*-2*-3)", nil}, 41, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateMathematicalExpression(tt.args.expressionString, tt.args.independentVariables, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateMathematicalExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EvaluateMathematicalExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}

var mathematicalExpression float64

func BenchmarkEvaluateMathematicalExpression(b *testing.B) {
	expression := "x*x+5*x+10"
	varMap := map[string]float64{"x": 0}

	for i := 0; i < b.N; i++ {
		mathematicalExpression1, err := EvaluateMathematicalExpression(expression, varMap, 0)
		if err != nil {
			b.Error(err)
		}
		mathematicalExpression = mathematicalExpression1
	}
	b.Log(mathematicalExpression)
}
