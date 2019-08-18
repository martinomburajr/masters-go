package evolution

import (
	"reflect"
	"testing"
)

func TestDualTree_FromSymbolicExpressionSet(t *testing.T) {
	tests := []struct {
		name    string
		fields  *DualTree
		args    []SymbolicExpression
		wantErr bool
	}{
		{"nil terminalSet", &DualTree{}, nil, true},
		{"empty terminalSet", &DualTree{}, make([]SymbolicExpression, 0), true},
		{"T", &DualTree{}, []SymbolicExpression{X1}, false},
		{"NT", &DualTree{}, []SymbolicExpression{Mult}, true},
		{"T-T", &DualTree{}, []SymbolicExpression{X1, Const4}, true},
		{"NT-NT", &DualTree{}, []SymbolicExpression{Mult, Sub}, true},
		{"T-NT(1)", &DualTree{}, []SymbolicExpression{X1, Sin}, false},
		{"T-NT(2)", &DualTree{}, []SymbolicExpression{X1, Sub}, true},
		{"T-NT(2)-T", &DualTree{}, []SymbolicExpression{X1, Add, Const4}, false},
		{"T-NT(2)-T-NT(2)-T", &DualTree{}, []SymbolicExpression{X1, Add, Const8, Mult, Const4}, false},
		{"T-NT(2)-T-NT(1)-T", &DualTree{}, []SymbolicExpression{X1, Add, Const4, Sin, Const4}, true},
		{"T-NT(2)-T-NT(1)", &DualTree{}, []SymbolicExpression{X1, Mult, Const4, Sub, Const8, Sin}, false},
		{"T-NT(1)-NT(1)-NT(1)-NT(1)", &DualTree{}, []SymbolicExpression{X1, Sin, Sin, Sin, Sin}, false},
		{"T-NT(1)-NT(2)-T-NT(1)", &DualTree{}, []SymbolicExpression{X1, Sin, Add, Const8, Sin}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.FromSymbolicExpressionSet(tt.args); (err != nil) != tt.wantErr {
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

func BenchmarkDualTree_FromSymbolicExpressionSet_1(b *testing.B) {
	expressionSet1 := GenerateRandomSymbolicExpressionSet(1)
	tree1 := DualTree{}
	for i := 0; i < b.N; i++ {
		tree1.FromSymbolicExpressionSet(expressionSet1)
	}
}

func BenchmarkDualTree_FromSymbolicExpressionSet_10(b *testing.B) {
	expressionSet1 := GenerateRandomSymbolicExpressionSet(10)
	tree1 := DualTree{}
	for i := 0; i < b.N; i++ {
		tree1.FromSymbolicExpressionSet(expressionSet1)
	}
}

func BenchmarkDualTree_FromSymbolicExpressionSet_1000(b *testing.B) {
	expressionSet1 := GenerateRandomSymbolicExpressionSet(1000)
	tree1 := DualTree{}
	for i := 0; i < b.N; i++ {
		tree1.FromSymbolicExpressionSet(expressionSet1)
	}
}

func BenchmarkDualTree_FromSymbolicExpressionSet_100000(b *testing.B) {
	expressionSet1 := GenerateRandomSymbolicExpressionSet(100000)
	tree100000 := DualTree{}
	for i := 0; i < b.N; i++ {
		tree100000.FromSymbolicExpressionSet(expressionSet1)
	}
}

func BenchmarkDualTree_FromSymbolicExpressionSet_1000000(b *testing.B) {
	expressionSet1 := GenerateRandomSymbolicExpressionSet(1000000)
	tree1 := DualTree{}
	for i := 0; i < b.N; i++ {
		tree1.FromSymbolicExpressionSet(expressionSet1)
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
		{"T-NT-T", Tree1(), X1.value + " " + Mult.value + " " + Const4.value, false},
		{"T-NT-T-NT-T", Tree2(), X1.value + " " + Sub.value + " " + X1.value + " " + Mult.value + " " + Const4.value, false},
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

func TestGenerateRandomTree(t *testing.T) {
	type args struct {
		maxDepth     int
		terminals    []SymbolicExpression
		nonTerminals []SymbolicExpression
	}
	tests := []struct {
		name    string
		args    args
		want    *DualTree
		wantErr bool
	}{
		{"err-lowMaxDepth", args{-1, nil, nil}, nil, true},
		{"err-lowMaxDepth", args{-1, nil, nil}, nil, true},
		{"err-lowMaxDepth", args{-1, nil, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateRandomTree(tt.args.maxDepth, tt.args.terminals, tt.args.nonTerminals)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRandomTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateRandomTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateRandomSymbolicExpressionSet(t *testing.T) {
	tests := []struct {
		name string
		size int
		want []SymbolicExpression
	}{
		{"size 0", 0, []SymbolicExpression{X1}},
		{"size -1", -1, []SymbolicExpression{X1}},
		{"size 1", 1, []SymbolicExpression{X1}},
		{"size 2", 2, []SymbolicExpression{X1}},
		{"size 3", 3, []SymbolicExpression{X1, Add, X1}},
		{"size 4", 4, []SymbolicExpression{X1, Add, X1}},
		{"size 5", 5, []SymbolicExpression{X1, Add, X1, Add, X1}},
		{"size 6", 6, []SymbolicExpression{X1, Add, X1, Add, X1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateRandomSymbolicExpressionSet(tt.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateRandomSymbolicExpressionSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDualTree_Leafs(t *testing.T) {
	tests := []struct {
		name    string
		tree    *DualTree
		want    []*DualTreeNode
		wantErr bool
	}{
		{"nil", TreeNil(), nil, true},
		{"T", Tree0(), []*DualTreeNode{X1.ToDualTreeNode(0)}, false},
		{"T-NT-T", Tree1(), []*DualTreeNode{X1.ToDualTreeNode(0), Const4.ToDualTreeNode(1)}, false},
		{"T-NT-T-NT-T", Tree2(), []*DualTreeNode{X1.ToDualTreeNode(0), X1.ToDualTreeNode(0),
			Const4.ToDualTreeNode(1)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tree.Leafs()
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.Leafs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if !got[i].IsValEqual(tt.want[i]) {
					t.Errorf("DualTree.Leafs() = %v, want %v", got[i].value, tt.want[i].value)
				}
			}
		})
	}
}

func TestDualTree_Count(t *testing.T) {
	tests := []struct {
		name string
		tree *DualTree
		want int
	}{
		{"nil", TreeNil(), 0},
		{"T", Tree0(), 1},
		{"T-NT-T", Tree1(), 3},
		{"T-NT-T-NT-T", Tree2(), 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tree.Count(); got != tt.want {
				t.Errorf("DualTree.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDualTree_RandomLeaf(t *testing.T) {
	tests := []struct {
		name    string
		tree    *DualTree
		want    *DualTreeNode
		wantErr bool
	}{
		{"nil", TreeNil(), nil, true},
		{"T", Tree0(), X1.ToDualTreeNode(0), false},
		{"T-NT-T", Tree1(), X1.ToDualTreeNode(0), false},
		{"T-NT-T-NT-T", Tree2(), Const4.ToDualTreeNode(0), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tree.RandomLeaf()
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.RandomLeaf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			nodes, err := tt.tree.Leafs()
			if err != nil && tt.wantErr {
				return
			}
			for i := range nodes {
				if nodes[i].value == got.value {
					return
				}
			}
			t.Error("Didnt get that value")
		})
	}
}

func TestDualTree_Branches(t *testing.T) {
	tests := []struct {
		name    string
		tree    *DualTree
		want    []*DualTreeNode
		wantErr bool
	}{
		{"nil", TreeNil(), nil, true},
		{"T", Tree0(), nil, true},
		{"T-NT-T", Tree1(), []*DualTreeNode{Mult.ToDualTreeNode(0)}, false},
		{"T-NT-T", Tree8(), []*DualTreeNode{Mult.ToDualTreeNode(0)}, false},
		{"T-NT-T-NT-T", Tree2(), []*DualTreeNode{Sub.ToDualTreeNode(0), Mult.ToDualTreeNode(0)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tree.Branches()
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.Branches() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if !got[i].IsValEqual(tt.want[i]) {
					t.Errorf("DualTree.Leafs() = %v, want %v", got[i].value, tt.want[i].value)
				}
			}
		})
	}
}