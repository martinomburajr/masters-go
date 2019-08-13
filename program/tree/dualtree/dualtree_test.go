package dualtree

import (
	"reflect"
	"testing"
)

func TestDualTree_FromString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		fields  *DualTree
		args    args
		wantErr bool
	}{
		{"empty string", &DualTree{}, args{""}, true},
		{"empty string", &DualTree{}, args{""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.FromString(tt.args.str); (err != nil) != tt.wantErr {
				t.Errorf("DualTree.FromString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDualTree_GetEquationSlice(t *testing.T) {
	type args struct {
		equationString string
	}
	tests := []struct {
		name    string
		fields  *DualTree
		args    args
		want    []string
		wantErr bool
	}{
		{"empty string", &DualTree{}, args{""}, nil, true},
		{"1 T", &DualTree{}, args{"x"}, []string{"x"}, false},
		{"1 NT", &DualTree{}, args{"*"}, []string{"*"}, false},
		{"7 items", &DualTree{}, args{"1,*,x,*,2,/,3"}, []string{"1", "*", "x", "*", "2", "/", "3"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := &DualTree{
				root: tt.fields.root,
				lock: tt.fields.lock,
			}
			got, err := bst.GetEquationSlice(tt.args.equationString)
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.GetEquationSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DualTree.GetEquationSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDualTree_FromTerminalSet(t *testing.T) {

	type args struct {
		E []NodeType
	}
	tests := []struct {
		name    string
		fields  *DualTree
		args    []NodeType
		wantErr bool
	}{
		{"nil terminalSet", &DualTree{}, nil, true},
		{"empty terminalSet", &DualTree{}, make([]NodeType, 0), true},
		{"T", &DualTree{}, []NodeType{X1}, false},
		{"NT", &DualTree{}, []NodeType{Mult}, true},
		{"T-T", &DualTree{}, []NodeType{X1, Const1}, true},
		{"NT-NT", &DualTree{}, []NodeType{Mult, Sub}, true},
		{"T-NT(1)", &DualTree{}, []NodeType{X1, Sin}, false},
		{"T-NT(2)", &DualTree{}, []NodeType{X1, Sub}, true},
		{"T-NT(2)-T", &DualTree{}, []NodeType{X1, Add, Const1}, false},
		{"T-NT(2)-T-NT(2)-T", &DualTree{}, []NodeType{X1, Add, Const2, Mult, Const1}, false},
		{"T-NT(2)-T-NT(1)-T", &DualTree{}, []NodeType{X1, Add, Const1, Sin, Const1}, true},
		{"T-NT(2)-T-NT(1)", &DualTree{}, []NodeType{X1, Mult, Const1, Sub, Const2, Sin}, false},
		{"T-NT(1)-NT(1)-NT(1)-NT(1)", &DualTree{}, []NodeType{X1, Sin, Sin, Sin, Sin}, false},
		{"T-NT(1)-NT(2)-T-NT(1)", &DualTree{}, []NodeType{X1, Sin, Add, Const2, Sin}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.FromTerminalSet(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("DualTree.FromNodeTypes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				expected := make([]string, 0)
				for i := range tt.args {
					if len(tt.args) == 1 && tt.args[0].kind >= 1 {
						continue
					}
					expected = append(expected, tt.args[i].value)
				}

				got := make([]string, 0)
				tt.fields.InOrderTraverse(func(node *DualTreeNode) {
					got = append(got, node.value)
				})

				if !reflect.DeepEqual(expected, got) {
					t.Errorf("expected: %#v ||| got: %#v", expected, got)
				}

				tt.fields.String()
			}

		})
	}

}

func Test_arityRemainder(t *testing.T) {
	tests := []struct {
		name string
		tree *DualTree
		want int
	}{
		{"full - NT(2)", Tree1(), 0},
		{"full - NT(2)", Tree2(), 0},
		{"full - NT(1)", Tree3(), 0},
		{"full - NT(1)", Tree4(), 0},
		{"full - NT(1)", Tree5(), 1},
		{"half - NT(1)", Tree6(), 1},
		{"empty - NT(2)", Tree7(), 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := arityRemainder(tt.tree.root); got != tt.want {
				t.Errorf("arityRemainder() = %v, want %v", got, tt.want)
			}
		})
	}
}

/**
THIS DOES NOT TEST OR CORRECT FOR TRIG OPERATORS YET
*/
func TestDualTree_ToMathematicalString(t *testing.T) {
	tests := []struct {
		name    string
		fields  *DualTree
		want    string
		wantErr bool
	}{
		{"nil", TreeNil(), "", true},
		{"T", Tree0(), "x", false},
		{"T-NT-T", Tree1(), X1.value + Mult.value + Const1.value, false},
		{"T-NT-T-NT-T", Tree2(), X1.value + Sub.value + X1.value + Mult.value + Const1.value, false},
		{"NT(1)", Tree5(), "", true},
		{"T - NT(2)", Tree6(), "", true},
		{"T - NT(2)", Tree7(), "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.ToMathematicalString()
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.ToMathematicalString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DualTree.ToMathematicalString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDualTree_Validate(t *testing.T) {
	tests := []struct {
		name    string
		tree    *DualTree
		wantErr bool
	}{
		{"nil", TreeNil(), true},
		{"T", Tree0(), false},
		{"T-NT-T", Tree1(), false},
		{"T-NT-T-NT-T", Tree2(), false},
		{"T-NT-T-NT(1)", Tree3(), false},
		{"T-NT(1)", Tree4(), false},
		{"NT(1)", Tree5(), true},
		{"T-NT(2)", Tree6(), true},
		{"NT", Tree7(), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tree.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("DualTree.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
