package evolution

import (
	"reflect"
	"sync"
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
		{"T", TreeT_0(), "x", false},
		{"T-NT-T", TreeT_NT_T_0(), X1.value + " " + Mult.value + " " + Const4.value, false},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_0(), X1.value + " " + Sub.value + " " + X1.value + " " + Mult.value + " " + Const4.value, false},
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
		{"T", TreeT_0(), false},
		{"T-NT-T", TreeT_NT_T_0(), false},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_0(), false},
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
		{"T", TreeT_0(), []*DualTreeNode{X1.ToDualTreeNode(0)}, false},
		{"T-NT-T", TreeT_NT_T_0(), []*DualTreeNode{X1.ToDualTreeNode(0), Const4.ToDualTreeNode(1)}, false},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_0(), []*DualTreeNode{X1.ToDualTreeNode(0), X1.ToDualTreeNode(0),
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
		{"T", TreeT_0(), 1},
		{"T-NT-T", TreeT_NT_T_0(), 3},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_0(), 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tree.Size(); got != tt.want {
				t.Errorf("DualTree.Size() = %v, want %v", got, tt.want)
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
		{"T", TreeT_0(), X1.ToDualTreeNode(0), false},
		{"T-NT-T", TreeT_NT_T_0(), X1.ToDualTreeNode(0), false},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_0(), Const4.ToDualTreeNode(0), false},
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
		{"T", TreeT_0(), nil, true},
		{"T-NT-T", TreeT_NT_T_0(), []*DualTreeNode{Mult.ToDualTreeNode(0)}, false},
		{"T-NT-T", Tree8(), []*DualTreeNode{Mult.ToDualTreeNode(0)}, false},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_0(), []*DualTreeNode{Sub.ToDualTreeNode(0), Mult.ToDualTreeNode(0)}, false},
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

func TestDualTree_AddSubTree(t *testing.T) {
	tests := []struct {
		name    string
		tree    *DualTree
		subTree *DualTree
		wantErr bool
	}{
		{"nil-subTree", TreeNil(), nil, true},
		{"nil-subTree-root", TreeNil(), TreeNil(), true},
		{"err-subTree-T", TreeNil(), TreeT_0(), true},
		{"nil-tree-T", TreeNil(), TreeT_NT_T_0(), true},
		{"tree-T", TreeT_0(), TreeT_NT_T_0(), true},
		{"T-NT-T + T-NT-T", TreeT_NT_T_0(), TreeT_NT_T_1(), false},
		{"T-NT-T + T-NT-T-NT-T", TreeT_NT_T_0(), TreeT_NT_T_NT_T_3(), false},
		{"T-NT-T + T-NT-T-NT-T", TreeT_NT_T_NT_T_3(), TreeT_NT_T_NT_T_3(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err1 := tt.tree.AddSubTree(tt.subTree); (err1 != nil) != tt.wantErr {
				t.Errorf("DualTree.AddSubTree() error = %v, wantErr %v", err1, tt.wantErr)
			} else {
				if err1 == nil {
					contains, err := tt.tree.Contains(tt.subTree)
					if err != nil {
						t.Error(err)
					}
					if !contains {
						t.Errorf("The main tree does not contain elements of the subTree")
					}
					tt.tree.String()
				}
			}

		})
	}
}

func TestDualTree_Contains(t *testing.T) {
	tests := []struct {
		name    string
		tree    *DualTree
		subTree *DualTree
		want    bool
		wantErr bool
	}{
		{"nil-subTree", TreeNil(), nil, false, true},
		{"nil-subTree-root", TreeNil(), TreeNil(), false, true},
		{"err-subTree-T", TreeNil(), TreeT_0(), false, true},
		{"nil-tree-T", TreeNil(), TreeT_NT_T_0(), false, true},
		{"same - T in T", TreeT_0(), TreeT_0(), true, false},
		{"diff - T in T", TreeT_0(), TreeT_1(), false, false},
		{"diff sizes", TreeT_0(), TreeT_NT_T_NT_T_1(), false, false},
		{"same - T in T-NT-T", TreeT_NT_T_0(), TreeT_0(), true, false},
		{"diff - T in T-NT-T", TreeT_NT_T_0(), TreeT_1(), false, false},
		{"same - T-NT-T in T-NT-T", TreeT_NT_T_0(), TreeT_NT_T_0(), true, false},
		{"same - T-NT-T in T-NT-T-NT-T", TreeT_NT_T_NT_T_0(), TreeT_NT_T_0(), true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tree.Contains(tt.subTree)
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.Contains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DualTree.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ############### REVIEW! IT COULD BE INCONSISTENT
func TestDualTree_RandomBranch(t *testing.T) {
	tests := []struct {
		name    string
		tree    *DualTree
		want    *DualTreeNode
		wantErr bool
	}{
		{"nil", TreeNil(), nil, true},
		{"T", TreeT_0(), nil, true},
		{"T-NT-T", TreeT_NT_T_0(), X1.ToDualTreeNode(0), false},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_3(), X1.ToDualTreeNode(0), false},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_3(), Const8.ToDualTreeNode(0), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tree.RandomBranch()
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.RandomBranch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				b, err := tt.tree.ContainsNode(got)
				if err != nil {
					t.Error(err)
				}
				if !b {
					t.Error("could not find object inside")
				}
			}
		})
	}
}

func TestDualTree_ContainsNode(t *testing.T) {
	tests := []struct {
		name     string
		tree     *DualTree
		treeNode *DualTreeNode
		want     bool
		wantErr  bool
	}{
		{"nil", TreeNil(), nil, false, true},
		{"nil node", TreeT_0(), nil, false, true},
		{"same | T in T", TreeT_0(), X1.ToDualTreeNode(0), true, false},
		{"same | T in T-NT-T", TreeT_NT_T_0(), X1.ToDualTreeNode(0), true, false},
		{"same | T in T-NT-T", TreeT_NT_T_0(), Const4.ToDualTreeNode(0), true, false},
		{"same | T in T-NT-T", TreeT_NT_T_0(), Const8.ToDualTreeNode(0), false, false},
		{"same | T in T-NT-T-NT-T", TreeT_NT_T_NT_T_3(), Sub.ToDualTreeNode(0), true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := &DualTree{
				root: tt.tree.root,
				lock: tt.tree.lock,
			}
			got, err := bst.ContainsNode(tt.treeNode)
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.ContainsNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DualTree.ContainsNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDualTree_DeleteSubTree(t *testing.T) {
	tests := []struct {
		name         string
		tree         *DualTree
		startingSize int
		wantErr      bool
	}{
		{"nil", TreeNil(), 0, true},
		{"T", TreeT_0(), 0, true},
		{"T-NT-T", TreeT_NT_T_0(), TreeT_NT_T_0().Size(), false},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_3(), TreeT_NT_T_NT_T_3().Size(), false},
		{"TreeT_NT_T_NT_T_NT_T_0", TreeT_NT_T_NT_T_NT_T_0(), TreeT_NT_T_NT_T_NT_T_0().Size(), false},
		{"TreeT_NT_T_NT_T_NT_T_NT_T_0", TreeT_NT_T_NT_T_NT_T_NT_T_0(), TreeT_NT_T_NT_T_NT_T_NT_T_0().Size(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tt.tree.String()
			if err = tt.tree.DeleteSubTree(); (err != nil) != tt.wantErr {
				t.Errorf("DualTree.DeleteSubTree() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				size := tt.tree.Size()
				if tt.startingSize <= size {
					t.Errorf("DualTree.DeleteSubTree() cannot be same size after delete = %d, wantErr %d", tt.startingSize,
						size)
				}
			}
			tt.tree.String()
		})
	}
}

func TestDualTree_SoftDeleteSubTree(t *testing.T) {
	type fields struct {
		root *DualTreeNode
		lock sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := &DualTree{
				root: tt.fields.root,
				lock: tt.fields.lock,
			}
			if err := bst.SoftDeleteSubTree(); (err != nil) != tt.wantErr {
				t.Errorf("DualTree.SoftDeleteSubTree() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// NEEDS MORE TESTS
func TestDualTree_MutateTerminal(t *testing.T) {
	tests := []struct {
		name    string
		tree    *DualTree
		oldTree *DualTree
		args    []SymbolicExpression
		wantErr bool
	}{
		{"nil", TreeNil(), TreeNil(), nil, true},
		{"err-nil-symbExpressSet", TreeT_0(), TreeT_0(), nil, true},
		{"err-empty-symbExpressSet", TreeT_0(), TreeT_0(), make([]SymbolicExpression, 0), true},
		{"T", TreeT_0(), TreeT_0(), []SymbolicExpression{Const4}, false},
		{"T-Same", TreeT_0(), TreeT_0(), []SymbolicExpression{X1}, false},
		{"T", TreeT_1(), TreeT_1(), []SymbolicExpression{Const4}, false},
		{"T-NT-T", TreeT_NT_T_0(), TreeT_NT_T_0(), []SymbolicExpression{Const8, Const4, X1}, false},
		{"T-NT-T", TreeT_NT_T_1(), TreeT_NT_T_1(), []SymbolicExpression{Const8, Const4, X1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if err = tt.tree.MutateTerminal(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("DualTree.MutateTerminal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				oldTreeLeafs, err := tt.oldTree.Leafs()
				if err != nil {
					t.Error(err)
				}
				newTreeLeafs, err := tt.tree.Leafs()
				if err != nil {
					t.Error(err)
				}
				if len(oldTreeLeafs) != len(newTreeLeafs) {
					t.Errorf("len of oldTree is not equal to new tree %d | got: %d", len(oldTreeLeafs), len(newTreeLeafs))
				}

				diffCount := 0
				for i := 0; i < len(oldTreeLeafs); i++ {
					if !oldTreeLeafs[i].IsValEqual(newTreeLeafs[i]) {
						diffCount++
					}
					if diffCount > 1 {
						t.Errorf("old and new tree should by different by only a single node. "+
							"got: %#v |  original: %#v", newTreeLeafs, oldTreeLeafs)
					}
				}

				TestSectionDivider("MUTATE TERMINAL: BEFORE", t)
				tt.oldTree.String()
				TestSectionDivider("MUTATE TERMINAL: AFTER", t)
				tt.tree.String()
			}
		})
	}
}

func TestDualTree_MutateNonTerminal(t *testing.T) {
	tests := []struct {
		name    string
		tree    *DualTree
		oldTree *DualTree
		args    []SymbolicExpression
		wantErr bool
	}{
		{"nil", TreeNil(), TreeNil(), nil, true},
		{"err-nil-symbExpressSet", TreeT_0(), TreeT_0(), nil, true},
		{"err-empty-symbExpressSet", TreeT_0(), TreeT_0(), make([]SymbolicExpression, 0), true},
		{"T", TreeT_0(), TreeT_0(), []SymbolicExpression{Const4}, true},
		{"T-Same", TreeT_0(), TreeT_0(), []SymbolicExpression{X1}, true},
		{"T", TreeT_1(), TreeT_1(), []SymbolicExpression{Const4}, true},
		{"T-NT-T NT-Same", TreeT_NT_T_0(), TreeT_NT_T_0(), []SymbolicExpression{Mult}, false},
		{"T-NT-T", TreeT_NT_T_1(), TreeT_NT_T_1(), []SymbolicExpression{Sub}, false},
		{"T-NT-T-NT-T-NT-T", TreeT_NT_T_NT_T_NT_T_0(), TreeT_NT_T_NT_T_NT_T_0(), []SymbolicExpression{Sub}, false},
		{"T-NT-T-NT-T-NT-T-NT-T-SAME", TreeT_NT_T_NT_T_NT_T_NT_T_0(), TreeT_NT_T_NT_T_NT_T_NT_T_0(),
			[]SymbolicExpression{Mult},
			false},
		{"T-NT-T-NT-T-NT-T-NT-T", TreeT_NT_T_NT_T_NT_T_NT_T_0(), TreeT_NT_T_NT_T_NT_T_NT_T_0(),
			[]SymbolicExpression{Add},
			false},
		{"T-NT-T-NT-T-NT-T-NT-T-SAME", TreeT_NT_T_NT_T_NT_T_NT_T_0(), TreeT_NT_T_NT_T_NT_T_NT_T_0(),
			[]SymbolicExpression{Mult},
			false},
		{"T-NT-T-NT-T-NT-T-NT-T-SAME", TreeT_NT_T_NT_T_NT_T_NT_T_0(), TreeT_NT_T_NT_T_NT_T_NT_T_0(),
			[]SymbolicExpression{Mult, Add, Sub},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if err = tt.tree.MutateNonTerminal(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("DualTree.MutateNonTerminal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				oldTreeNonTerminals, err := tt.oldTree.Branches()
				if err != nil {
					t.Error(err)
				}
				newTreeNonTerminals, err := tt.tree.Branches()
				if err != nil {
					t.Error(err)
				}
				if len(oldTreeNonTerminals) != len(newTreeNonTerminals) {
					t.Errorf("len of oldTree is not equal to new tree %d | got: %d", len(oldTreeNonTerminals), len(newTreeNonTerminals))
				}

				diffCount := 0
				for i := 0; i < len(oldTreeNonTerminals); i++ {
					if !oldTreeNonTerminals[i].IsValEqual(newTreeNonTerminals[i]) {
						diffCount++
					}
					if diffCount > 1 {
						t.Errorf("old and new tree should by different by only a single node. "+
							"got: %#v |  original: %#v", newTreeNonTerminals, oldTreeNonTerminals)
					}
				}

				TestSectionDivider("MUTATE NON-TERMINAL: BEFORE", t)
				tt.oldTree.String()
				TestSectionDivider("MUTATE NON-TERMINAL: AFTER", t)
				tt.tree.String()
			}
		})
	}
}

func TestDualTree_HasDiverseNonTerminalSet(t *testing.T) {
	tests := []struct {
		name    string
		tree    *DualTree
		want    bool
		wantErr bool
	}{
		{"nil", TreeNil(), false, true},
		{"TreeT_0", TreeT_0(), false, true},
		{"TreeT_1", TreeT_1(), false, true},
		{"TreeT_NT_T_2", TreeT_NT_T_2(), false, false},
		{"TreeT_NT_T_3", TreeT_NT_T_3(), false, false},
		{"TreeT_NT_T_0", TreeT_NT_T_0(), false, false},
		{"TreeT_NT_T_NT_T_1", TreeT_NT_T_NT_T_1(), true, false},
		{"TreeT_NT_T_NT_T_2", TreeT_NT_T_NT_T_2(), true, false},
		{"TreeT_NT_T_NT_T_3", TreeT_NT_T_NT_T_3(), true, false},
		{"TreeT_NT_T_NT_T_NT_T_0", TreeT_NT_T_NT_T_NT_T_0(), true, false},
		{"TreeT_NT_T_NT_T_NT_T_1", TreeT_NT_T_NT_T_NT_T_1(), true, false},
		{"TreeT_NT_T_NT_T_NT_T_2", TreeT_NT_T_NT_T_NT_T_2(), false, false},
		{"TreeT_NT_T_NT_T_NT_T_NT_T_0", TreeT_NT_T_NT_T_NT_T_NT_T_0(), false, false},
		{"TreeT_NT_T_NT_T_NT_T_NT_T_1", TreeT_NT_T_NT_T_NT_T_NT_T_1(), true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tree.HasDiverseNonTerminalSet()
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.HasDiverseNonTerminalSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DualTree.HasDiverseNonTerminalSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSplitter(b *testing.B) {

}

func TestSplitter(t *testing.T) {

	symEx0 := []SymbolicExpression{X1, Add, Const4}
	node0 := symEx0[1].ToDualTreeNode(1)
	node0.left = symEx0[0].ToDualTreeNode(0)
	node0.right = symEx0[2].ToDualTreeNode(2)
	node0Arr := []*DualTreeNode{node0}

	symEx1 := []SymbolicExpression{X1, Add, Const4, Mult, Const8}
	node1 := symEx1[1].ToDualTreeNode(1)
	node1.left = symEx1[0].ToDualTreeNode(0)
	node2 := symEx1[3].ToDualTreeNode(1)
	node2.left = symEx1[2].ToDualTreeNode(1)
	node3 := symEx1[4].ToDualTreeNode(1)
	node1Arr := []*DualTreeNode{node1, node2, node3}

	tests := []struct {
		name          string
		expressionSet []SymbolicExpression
		want          []*DualTreeNode
		wantErr       bool
	}{
		{"err-T-NT", []SymbolicExpression{X1, Add}, nil,
			true},
		{"err-T-NT-T-NT", []SymbolicExpression{X1, Add, Const4, Sub}, nil, true},

		{"T-NT-T", symEx0, node0Arr, false},
		{"T-NT-T-NT-T", symEx1, node1Arr, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Splitter(tt.expressionSet)
			if (err != nil) != tt.wantErr {
				t.Errorf("Splitter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for e := range got {
				got[e].IsValEqual(tt.want[e])
			}
		})
	}
}

func TestDualTree_FromSymbolicExpressionSet2(t *testing.T) {
	tests := []struct {
		name        string
		fields      *DualTree
		terminalSet []SymbolicExpression
		wantErr     bool
	}{
		{"nil-terminalset", &DualTree{}, nil, true},
		{"err-terminalset<1", &DualTree{}, []SymbolicExpression{}, true},
		{"err-terminalset-only-NT", &DualTree{}, []SymbolicExpression{Add}, true},
		{"T", TreeT_0(), []SymbolicExpression{X1}, false},
		{"T-NT-T", TreeT_NT_T_0(), []SymbolicExpression{X1, Mult, Const4}, false},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_0(), []SymbolicExpression{X1, Sub, X1, Mult, Const4}, false},
		{"T-NT-T-NT-T-NT-T", TreeT_NT_T_NT_T_NT_T_0(), []SymbolicExpression{Const4, Sub, Const0, Add, Const4, Add,
			Const8}, false},
		{"T-NT-TreeT_NT_T_NT_T_NT_T_NT_T_0-NT-T-NT-T", TreeT_NT_T_NT_T_NT_T_NT_T_0(), []SymbolicExpression{X1, Mult,
			X1, Mult, X1, Mult, X1, Mult, X1},false},
		{"TreeT_NT_T_NT_T_NT_T_NT_T_1", TreeT_NT_T_NT_T_NT_T_NT_T_1(), []SymbolicExpression{X1, Mult,
			X1, Mult, X1, Mult, X1, Add, Const4},false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if err = tt.fields.FromSymbolicExpressionSet2(tt.terminalSet); (err != nil) != tt.wantErr {
				t.Errorf("DualTree.FromSymbolicExpressionSet2() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				expressionSet := tt.fields.ToSymbolicExpressionSet()
				if len(expressionSet) != len(tt.terminalSet) {
					t.Errorf("Generated Tree not the same LENGTH as input symbolic set error = %q, wantErr %q",
						expressionSet,
						tt.terminalSet)
				}
				for e := range tt.terminalSet {
					if tt.terminalSet[e].value != expressionSet[e].value {
						t.Errorf("Generated Tree not the same as input symbolic set error = %q, wantErr %q", expressionSet,
							tt.terminalSet)
					}
				}
			}

		})
		tt.fields.String()
	}
}
