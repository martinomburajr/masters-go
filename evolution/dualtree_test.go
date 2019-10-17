package evolution

import (
	"fmt"
	"log"
	"reflect"
	"sync"
	"testing"
)

//func TestDualTree_FromSymbolicExpressionSet(t *testing.T) {
//	tests := []struct {
//		name    string
//		fields  *DualTree
//		args    []SymbolicExpression
//		wantErr bool
//	}{
//		{"nil terminalSet", &DualTree{}, nil, true},
//		{"empty terminalSet", &DualTree{}, make([]SymbolicExpression, 0), true},
//		{"T", &DualTree{}, []SymbolicExpression{X1}, false},
//		{"NT", &DualTree{}, []SymbolicExpression{Mult}, true},
//		{"T-T", &DualTree{}, []SymbolicExpression{X1, Const4}, true},
//		{"NT-NT", &DualTree{}, []SymbolicExpression{Mult, Sub}, true},
//		{"T-NT(1)", &DualTree{}, []SymbolicExpression{X1, Sin}, false},
//		{"T-NT(2)", &DualTree{}, []SymbolicExpression{X1, Sub}, true},
//		{"T-NT(2)-T", &DualTree{}, []SymbolicExpression{X1, Add, Const4}, false},
//		{"T-NT(2)-T-NT(2)-T", &DualTree{}, []SymbolicExpression{X1, Add, Const8, Mult, Const4}, false},
//		{"T-NT(2)-T-NT(1)-T", &DualTree{}, []SymbolicExpression{X1, Add, Const4, Sin, Const4}, true},
//		{"T-NT(2)-T-NT(1)", &DualTree{}, []SymbolicExpression{X1, Mult, Const4, Sub, Const8, Sin}, false},
//		{"T-NT(1)-NT(1)-NT(1)-NT(1)", &DualTree{}, []SymbolicExpression{X1, Sin, Sin, Sin, Sin}, false},
//		{"T-NT(1)-NT(2)-T-NT(1)", &DualTree{}, []SymbolicExpression{X1, Sin, Add, Const8, Sin}, false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := tt.fields.FromSymbolicExpressionSet(tt.args); (err != nil) != tt.wantErr {
//				t.Errorf("DualTree.FromNodeTypes() error = %v, wantErr %v", err, tt.wantErr)
//			}
//			if !tt.wantErr {
//				expected := make([]string, 0)
//				for i := range tt.args {
//					if len(tt.args) == 1 && tt.args[0].Kind >= 1 {
//						continue
//					}
//					expected = append(expected, tt.args[i].value)
//				}
//
//				got := make([]string, 0)
//				tt.fields.InOrderTraverse(func(node *DualTreeNode) {
//					got = append(got, node.value)
//				})
//
//				if !reflect.DeepEqual(expected, got) {
//					t.Errorf("expected: %#v ||| got: %#v", expected, got)
//				}
//
//				tt.fields.Print()
//			}
//
//		})
//	}
//}

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
		{"T", TreeT_0(), "(x)", false},
		{"T-NT-T", TreeT_NT_T_0(), fmt.Sprintf("((%s)%s(%s))", X1.value, Mult.value, Const4.value),
			false},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_0(),
			fmt.Sprintf("(((%s)%s(%s))%s(%s))", X1.value, Sub.value, X1.value, Mult.value,
				Const4.value), false},

		{"T-NT-T-NT-T", TreeXAddXMult4Sub9Mult0(),
			fmt.Sprintf("((((%s)%s((%s)%s(%s)))%s(%s))%s(%s))", X1.value, Add.value,
				X1.value,
				Mult.value, Const4.value, Sub.value, Const9.value, Mult.value, Const0.value), false},

		//{"NT(1)", Tree5(), "", true},
		//{"T - NT(2)", Tree6(), "", true},
		//{"T - NT(2)", Tree7(), "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.ToMathematicalString()
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.ToMathematicalString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DualTree.ToMathematicalString() = %v, isEqual %v", got, tt.want)
			}
			log.Print(got)
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
		depth        int
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
		{"err-nil-terminals", args{2, nil, nil}, nil, true},
		{"err-nil-non-terminals", args{2, []SymbolicExpression{X1}, nil}, nil, true},
		{"err-nil-empty-terminals", args{2, []SymbolicExpression{}, []SymbolicExpression{}}, nil, true},
		{"err-nil-empty-nonterminals", args{2, []SymbolicExpression{X1}, []SymbolicExpression{}}, nil, true},
		{"err-nil-empty-nonterminals", args{2, []SymbolicExpression{X1}, []SymbolicExpression{}}, nil, true},
		{"T", args{0, []SymbolicExpression{X1}, []SymbolicExpression{}}, TreeT_1(), false},
		{"T", args{0, []SymbolicExpression{X1, Const0, Const1, Const2, Const3, Const4, Const5, Const6,
			Const7, Const8, Const9}, []SymbolicExpression{}}, TreeT_1(), false},
		{"err-depth-1-no-NT", args{1, []SymbolicExpression{X1}, []SymbolicExpression{}}, TreeT_1(), true},
		{"depth-1", args{1, []SymbolicExpression{X1}, []SymbolicExpression{Add}}, TreeT_NT_T_0(), false},
		{"depth-2", args{2, []SymbolicExpression{X1}, []SymbolicExpression{Add}}, TreeT_NT_T_NT_T_NT_T_0(), false},
		{"depth-3-diverse", args{2, []SymbolicExpression{X1, Const0, Const1, Const2, Const3, Const4, Const5, Const6,
			Const7, Const8, Const9},
			[]SymbolicExpression{Add,
				Mult,
				Sub}},
			TreeT_NT_T_NT_T_NT_T_0(),
			false},
		{"depth-2-diverse", args{3, []SymbolicExpression{X1, Const0, Const1, Const2, Const3, Const4, Const5, Const6,
			Const7, Const8, Const9},
			[]SymbolicExpression{Add,
				Mult,
				Sub}},
			TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0(),
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			got, err := GenerateRandomTree(tt.args.depth, tt.args.terminals, tt.args.nonTerminals)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRandomTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				wantExpressionSet, err := tt.want.ToSymbolicExpressionSet()
				if err != nil {
					t.Error(err)
				}
				gotExpressionSet, err := got.ToSymbolicExpressionSet()
				if err != nil {
					t.Error(err)
				}

				if len(wantExpressionSet) != len(gotExpressionSet) {
					t.Errorf("They are not the same length error = %v, wantErr %v", gotExpressionSet, wantExpressionSet)
				}
				got.Print()
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
		{"size 1", -1, []SymbolicExpression{X1}},
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
				t.Errorf("GenerateRandomSymbolicExpressionSet() = %v, isEqual %v", got, tt.want)
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
		{"T", TreeT_0(), []*DualTreeNode{X1.ToDualTreeNode(RandString(5))}, false},
		{"T-NT-T", TreeT_NT_T_0(), []*DualTreeNode{X1.ToDualTreeNode(RandString(5)), Const4.ToDualTreeNode(RandString(5))}, false},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_0(), []*DualTreeNode{X1.ToDualTreeNode(RandString(5)), X1.ToDualTreeNode(RandString(5)),
			Const4.ToDualTreeNode(RandString(5))}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tree.Terminals()
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.Terminals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if !got[i].IsValEqual(tt.want[i]) {
					t.Errorf("DualTree.Terminals() = %v, isEqual %v", got[i].value, tt.want[i].value)
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
				t.Errorf("DualTree.Size() = %v, isEqual %v", got, tt.want)
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
		{"T", TreeT_0(), X1.ToDualTreeNode(RandString(5)), false},
		{"T-NT-T", TreeT_NT_T_0(), X1.ToDualTreeNode(RandString(5)), false},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_0(), Const4.ToDualTreeNode(RandString(5)), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tree.RandomTerminal()
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.RandomTerminal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			nodes, err := tt.tree.Terminals()
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
		{"T-NT-T", TreeT_NT_T_0(), []*DualTreeNode{Mult.ToDualTreeNode(RandString(5))}, false},
		{"T-NT-T", Tree8(), []*DualTreeNode{Mult.ToDualTreeNode(RandString(5))}, false},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_0(), []*DualTreeNode{Sub.ToDualTreeNode(RandString(5)), Mult.ToDualTreeNode(RandString(5))}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tree.NonTerminals()
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.NonTerminals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if !got[i].IsValEqual(tt.want[i]) {
					t.Errorf("DualTree.Terminals() = %v, isEqual %v", got[i].value, tt.want[i].value)
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
		{"nil-treeNode-T", TreeNil(), TreeT_NT_T_0(), true},
		{"treeNode-T", TreeT_0(), TreeT_NT_T_0(), false},
		{"T-NT-T + T-NT-T", TreeT_NT_T_0(), TreeT_NT_T_1(), false},
		{"T-NT-T + T-NT-T-NT-T", TreeT_NT_T_0(), TreeT_NT_T_NT_T_3(), false},
		{"T-NT-T + T-NT-T-NT-T", TreeT_NT_T_NT_T_3(), TreeT_NT_T_NT_T_3(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err1 := tt.tree.AddSubTree(tt.subTree); (err1 != nil) != tt.wantErr {
				t.Errorf("DualTree.StrategyAddSubTree() error = %v, wantErr %v", err1, tt.wantErr)
			} else {
				if err1 == nil {
					contains, err := tt.tree.ContainsSubTree(tt.subTree)
					if err != nil {
						t.Error(err)
					}
					if !contains {
						t.Errorf("The main treeNode does not contain elements of the subTree")
					}
					tt.tree.Print()
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
		{"nil-treeNode-T", TreeNil(), TreeT_NT_T_0(), false, true},
		{"same - T in T", TreeT_0(), TreeT_0(), true, false},
		{"diff - T in T", TreeT_0(), TreeT_1(), false, false},
		{"diff sizes", TreeT_0(), TreeT_NT_T_NT_T_1(), false, false},
		{"same - T in T-NT-T", TreeT_NT_T_0(), TreeT_0(), false, false},
		{"diff - T in T-NT-T", TreeT_NT_T_0(), TreeT_1(), false, false},
		{"same - T-NT-T in T-NT-T", TreeT_NT_T_0(), TreeT_NT_T_0(), true, false},
		{"same - T-NT-T in T-NT-T-NT-T", TreeT_NT_T_NT_T_0(), TreeT_NT_T_0(), false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tree.ContainsSubTree(tt.subTree)
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.ContainsSubTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DualTree.ContainsSubTree() = %v, isEqual %v", got, tt.want)
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
		{"T-NT-T", TreeT_NT_T_0(), X1.ToDualTreeNode(RandString(5)), false},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_3(), X1.ToDualTreeNode(RandString(5)), false},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_3(), Const8.ToDualTreeNode(RandString(5)), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tree.RandomNonTerminal()
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.RandomNonTerminal() error = %v, wantErr %v", err, tt.wantErr)
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
		{"same | T in T", TreeT_0(), X1.ToDualTreeNode("0"), true, false},
		{"same | T in T-NT-T", TreeT_NT_T_0(), X1.ToDualTreeNode("2"), true, false},
		{"same | T in T-NT-T", TreeT_NT_T_0(), Const4.ToDualTreeNode("3"), true, false},
		{"same | T in T-NT-T", TreeT_NT_T_0(), Const8.ToDualTreeNode("2"), false, false},
		{"same | T in T-NT-T-NT-T", TreeT_NT_T_NT_T_3(), Sub.ToDualTreeNode("1"), true, false},
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
				t.Errorf("DualTree.ContainsNode() = %v, isEqual %v", got, tt.want)
			}
		})
	}
}

func TestDualTree_DeleteSubTree(t *testing.T) {
	tests := []struct {
		name             string
		tree             *DualTree
		deletionStrategy int
		startingSize     int
		wantErr          bool
	}{
		{"nil", TreeNil(), 0, 0, true},
		{"T", TreeT_0(), 0, 1, false},
		{"T", TreeT_0(), 1, 1, false},
		{"T", TreeT_0(), -1, 1, false},
		{"T", TreeT_0(), 20, 1, false},
		{"T-NT-T", TreeT_NT_T_0(), 0, TreeT_NT_T_0().Size(), false},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_3(), 0, TreeT_NT_T_NT_T_3().Size(), false},
		{"TreeT_NT_T_NT_T_NT_T_0", TreeT_NT_T_NT_T_NT_T_0(), 0, TreeT_NT_T_NT_T_NT_T_0().Size(), false},
		{"TreeT_NT_T_NT_T_NT_T_NT_T_0", TreeT_NT_T_NT_T_NT_T_NT_T_0(), 0, TreeT_NT_T_NT_T_NT_T_NT_T_0().Size(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tt.tree.Print()
			if err = tt.tree.DeleteSubTree(tt.deletionStrategy); (err != nil) != tt.wantErr {
				t.Errorf("DualTree.StrategyDeleteSubTree() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				size := tt.tree.Size()
				if tt.startingSize == 1 {
					if size < tt.startingSize {
						t.Errorf("DualTree.StrategyDeleteSubTree() cannot be same size after delete = %d, wantErr %d", tt.startingSize,
							size)
					}
				} else if tt.startingSize <= size {
					t.Errorf("DualTree.StrategyDeleteSubTree() cannot be same size after delete = %d, wantErr %d", tt.startingSize,
						size)
				}
			}
			tt.tree.Print()
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
				t.Errorf("DualTree.StrategySoftDeleteSubTree() error = %v, wantErr %v", err, tt.wantErr)
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
				oldTreeLeafs, err := tt.oldTree.Terminals()
				if err != nil {
					t.Error(err)
				}
				newTreeLeafs, err := tt.tree.Terminals()
				if err != nil {
					t.Error(err)
				}
				if len(oldTreeLeafs) != len(newTreeLeafs) {
					t.Errorf("len of oldTree is not equal to new treeNode %d | got: %d", len(oldTreeLeafs), len(newTreeLeafs))
				}

				diffCount := 0
				for i := 0; i < len(oldTreeLeafs); i++ {
					if !oldTreeLeafs[i].IsValEqual(newTreeLeafs[i]) {
						diffCount++
					}
					if diffCount > 1 {
						t.Errorf("old and new treeNode should by different by only a single node. "+
							"got: %#v |  original: %#v", newTreeLeafs, oldTreeLeafs)
					}
				}

				tt.oldTree.Print()
				tt.tree.Print()
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
				oldTreeNonTerminals, err := tt.oldTree.NonTerminals()
				if err != nil {
					t.Error(err)
				}
				newTreeNonTerminals, err := tt.tree.NonTerminals()
				if err != nil {
					t.Error(err)
				}
				if len(oldTreeNonTerminals) != len(newTreeNonTerminals) {
					t.Errorf("len of oldTree is not equal to new treeNode %d | got: %d", len(oldTreeNonTerminals), len(newTreeNonTerminals))
				}

				diffCount := 0
				for i := 0; i < len(oldTreeNonTerminals); i++ {
					if !oldTreeNonTerminals[i].IsValEqual(newTreeNonTerminals[i]) {
						diffCount++
					}
					if diffCount > 1 {
						t.Errorf("old and new treeNode should by different by only a single node. "+
							"got: %#v |  original: %#v", newTreeNonTerminals, oldTreeNonTerminals)
					}
				}

				tt.oldTree.Print()
				tt.tree.Print()
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
			got, err := tt.tree.hasDiverseNonTerminalSet()
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.hasDiverseNonTerminalSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DualTree.hasDiverseNonTerminalSet() = %v, isEqual %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSplitter(b *testing.B) {

}

func TestSplitter(t *testing.T) {

	symEx0 := []SymbolicExpression{X1, Add, Const4}
	node0 := symEx0[1].ToDualTreeNode(RandString(5))
	node0.left = symEx0[0].ToDualTreeNode(RandString(5))
	node0.right = symEx0[2].ToDualTreeNode(RandString(5))
	node0Arr := []*DualTreeNode{node0}

	symEx1 := []SymbolicExpression{X1, Add, Const4, Mult, Const8}
	node1 := symEx1[1].ToDualTreeNode(RandString(5))
	node1.left = symEx1[0].ToDualTreeNode(RandString(5))
	node2 := symEx1[3].ToDualTreeNode(RandString(5))
	node2.left = symEx1[2].ToDualTreeNode(RandString(5))
	node3 := symEx1[4].ToDualTreeNode(RandString(5))
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
		tree        *DualTree
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
			X1, Mult, X1, Mult, X1, Mult, X1}, false},
		{"TreeT_NT_T_NT_T_NT_T_NT_T_1", TreeT_NT_T_NT_T_NT_T_NT_T_1(), []SymbolicExpression{X1, Mult,
			X1, Mult, X1, Mult, X1, Add, Const4}, false},
		{"TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0()", TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0(), []SymbolicExpression{X1, Mult,
			X1, Mult, X1, Mult, X1, Add, Const4}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if err = tt.tree.FromSymbolicExpressionSet2(tt.terminalSet); (err != nil) != tt.wantErr {
				t.Errorf("DualTree.FromSymbolicExpressionSet2() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				expressionSet, err := tt.tree.ToSymbolicExpressionSet()
				if err != nil {
					t.Error(err)
				}
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
		tt.tree.Print()
	}
}

func TestDualTree_ToSymbolicExpressionSet(t *testing.T) {
	tests := []struct {
		name   string
		fields *DualTree
		want   []SymbolicExpression
	}{
		{"treeNode-nil", TreeNil(), []SymbolicExpression{}},
		{"T", TreeT_1(), []SymbolicExpression{Const0}},
		{"T-NT-T", TreeT_NT_T_0(), []SymbolicExpression{X1, Mult, Const4}},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_0(), []SymbolicExpression{X1, Sub, X1, Mult, Const4}},
		{"T-NT-T-NT-T-NT-T", TreeT_NT_T_NT_T_NT_T_1(), []SymbolicExpression{X1, Mult, X1, Add, Const4, Add, Const8}},
		{"T_NT_T_NT_T_NT_T_NT_T", TreeT_NT_T_NT_T_NT_T_NT_T_0(), []SymbolicExpression{X1, Mult, X1, Mult, X1, Mult, X1,
			Mult, X1}},
		{"T_NT_T_NT_T_NT_T_NT_T_NT_T", TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_0(), []SymbolicExpression{Const0, Add, Const1, Add,
			Const2, Add, Const3, Add, Const4, Add, Const5}},
		{"T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T", TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0(), []SymbolicExpression{Const0, Add, Const1, Add,
			Const2, Add, Const3, Add, Const4, Add, Const5, Add, Const6}},
		{"T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T", TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0(), []SymbolicExpression{Const0, Add, Const1, Add,
			Const2, Add, Const3, Add, Const4, Add, Const5, Add, Const6, Add, Const7}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := &DualTree{
				root: tt.fields.root,
				lock: tt.fields.lock,
			}
			got, err := bst.ToSymbolicExpressionSet()
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DualTree.ToSymbolicExpressionSet() = %v, \n isEqual %v", got, tt.want)
			}
		})
	}
}

func Test_weaver(t *testing.T) {
	type args struct {
		terminals    []SymbolicExpression
		nonTerminals []SymbolicExpression
	}
	tests := []struct {
		name string
		args args
		want []SymbolicExpression
	}{
		{"empty", args{[]SymbolicExpression{}, []SymbolicExpression{}}, []SymbolicExpression{}},
		{"T-1 | NT-0", args{[]SymbolicExpression{X1}, []SymbolicExpression{}}, []SymbolicExpression{X1}},
		{"T-1 | NT-1", args{[]SymbolicExpression{X1}, []SymbolicExpression{Add}}, []SymbolicExpression{X1, Add}},
		{"T-2 | NT-1", args{[]SymbolicExpression{X1, Const0}, []SymbolicExpression{Add}}, []SymbolicExpression{X1,
			Add, Const0}},
		{"T-2 | NT-2", args{[]SymbolicExpression{X1, Const0}, []SymbolicExpression{Add, Sub}}, []SymbolicExpression{X1,
			Add, Const0, Sub}},
		{"T-3 | NT-2", args{[]SymbolicExpression{X1, Const0, Const1}, []SymbolicExpression{Add, Sub}},
			[]SymbolicExpression{X1,
				Add, Const0, Sub, Const1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := weaver(tt.args.terminals, tt.args.nonTerminals); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("weaver() = %v, isEqual %v", got, tt.want)
			}
		})
	}
}

func TestDualTree_GetNode(t *testing.T) {
	tests := []struct {
		name       string
		fields     *DualTree
		value      string
		wantNode   *DualTreeNode
		wantParent *DualTreeNode
		wantErr    bool
	}{
		{"nil", TreeNil(), "", nil, nil, true},
		{"empty-value", TreeT_0(), "", nil, nil, true},
		{"wrong-value", TreeT_0(), "y", nil, nil, true},
		{"correct-value", TreeT_0(), "x", TreeT_0().root, TreeT_0().root, false},
		{"T-NT-T-wrong-value", TreeT_NT_T_0(), "3", nil, nil, true},
		{"T-NT-T-correct-value", TreeT_NT_T_0(), "x", TreeT_NT_T_0().root.left, TreeT_NT_T_0().root, false},
		{"T-NT-T-correct-value-other-T", TreeT_NT_T_0(), "4", TreeT_NT_T_0().root.right, TreeT_NT_T_0().root, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNode, gotParent, err := tt.fields.GetNode(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.GetNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNode, tt.wantNode) {
				t.Errorf("DualTree.GetNode() gotNode = %v, want %v", gotNode, tt.wantNode)
			}
			if !reflect.DeepEqual(gotParent, tt.wantParent) {
				t.Errorf("DualTree.GetNode() gotParent = %v, want %v", gotParent, tt.wantParent)
			}
		})
	}
}

func TestDualTree_Depth(t *testing.T) {
	tests := []struct {
		name    string
		fields  *DualTree
		want    int
		wantErr bool
	}{
		{"nil", TreeNil(), -1, true},
		{"T", TreeT_0(), 0, false},
		{"T", TreeT_1(), 0, false},
		{"T-NT-T", TreeT_NT_T_0(), 1, false},
		{"TreeT_NT_T_NT_T_0", TreeT_NT_T_NT_T_0(), 2, false},
		{"TreeT_NT_T_NT_T_0", TreeT_NT_T_NT_T_0(), 2, false},
		{"TreeVine_D3", TreeVine_D3(), 3, false},
		{"TreeVine_D4", TreeVine_D4(), 4, false},
		{"TreeVine_D5_R", TreeVine_D5_R(), 5, false},
		{"TreeVine_D6_R", TreeVine_D6_R(), 6, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DualTree{
				root: tt.fields.root,
				lock: tt.fields.lock,
			}
			d.Print()
			got, err := d.Depth()
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.Depth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DualTree.Depth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDualTree_DepthTo(t *testing.T) {
	type args struct {
		depth int
	}
	tests := []struct {
		name    string
		fields  *DualTree
		depth   int
		want    []*DualTreeNode
		wantErr bool
	}{
		{"nil root", TreeNil(), 0, nil, true},
		{"negative depth", TreeT_0(), -1, nil, true},
		{"T", TreeT_0(), 0, []*DualTreeNode{X1.ToDualTreeNode(RandString(5))}, false},
		{"T-NT-T", TreeT_NT_T_1(), 0, []*DualTreeNode{Add.ToDualTreeNode(RandString(5))}, false},
		{"T-NT-T", TreeT_NT_T_1(), 1, []*DualTreeNode{Add.ToDualTreeNode(RandString(5)), Const8.ToDualTreeNode(RandString(5)),
			X1.ToDualTreeNode(RandString(5))}, false},
		{"T-NT-NT-NT-T", TreeT_NT_T_NT_T_0(), 1, []*DualTreeNode{Sub.ToDualTreeNode(RandString(5)), Const4.ToDualTreeNode(RandString(5)),
			Mult.ToDualTreeNode(RandString(5))}, false},
		{"T-NT-NT-NT-T", TreeT_NT_T_NT_T_0(), 2, []*DualTreeNode{Sub.ToDualTreeNode(RandString(5)), Const4.ToDualTreeNode(RandString(5)),
			Mult.ToDualTreeNode(RandString(5)), X1.ToDualTreeNode(RandString(5)), X1.ToDualTreeNode(RandString(5))}, false},
		{"Vine-NT-NT-NT-T", TreeVine_D3(), 2, []*DualTreeNode{Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)),
			Sin.ToDualTreeNode(RandString(5))}, false},
		{"Vine-NT-NT-NT-T", TreeVine_D6_R(), 2, []*DualTreeNode{Add.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)),
			Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5))}, false},
		{"Vine-6", TreeVine_D6_R(), 5, []*DualTreeNode{Add.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)),
			Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)),
			Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)), X1.ToDualTreeNode(RandString(5))}, false},
		{"Vine-6v2", TreeVine_D6_R(), 6, []*DualTreeNode{Add.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)),
			Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)),
			Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)), X1.ToDualTreeNode(RandString(5)), X1.ToDualTreeNode(RandString(5))}, false},
		{"Vine-5", TreeVine_D5_R(), 5, []*DualTreeNode{Add.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)),
			Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)),
			Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)), X1.ToDualTreeNode(RandString(5))}, false},
		{"Vine-4", TreeVine_D4(), 4, []*DualTreeNode{Sin.ToDualTreeNode(RandString(5)),
			Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)), Sin.ToDualTreeNode(RandString(5)), X1.ToDualTreeNode(RandString(5))}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DualTree{
				root: tt.fields.root,
				lock: tt.fields.lock,
			}
			got, err := d.DepthTo(tt.depth)
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.DepthTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			contain, err := contains(got, tt.want)
			if err != nil {
				t.Error(err)
			}
			if !contain {
				t.Errorf("DualTree.DepthTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func contains(a, b []*DualTreeNode) (bool, error) {
	if len(a) != len(b) {
		return false, fmt.Errorf("not the same size a %v | b %v", a, b)
	}
	count := 0
	for i := range a {
		for j := range b {
			if a[i].IsValEqual(b[j]) {
				count++
				break
			}
		}
	}

	return count == len(a), nil
}

func TestDualTree_GetRandomSubTreeAtDepth(t *testing.T) {
	tests := []struct {
		name    string
		fields  *DualTree
		depth   int
		want    DualTree
		wantErr bool
	}{
		{"nil", TreeNil(), 0, DualTree{}, true},
		{"negative depth", TreeT_0(), -1, DualTree{}, true},
		{"T", TreeT_0(), 0, DualTree{}, false},
		{"TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_0", TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_0(), 2, DualTree{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DualTree{
				root: tt.fields.root,
				lock: tt.fields.lock,
			}
			got, err := d.GetRandomSubTreeAtDepth(tt.depth)
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.GetRandomSubTreeAtDepth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			b, err := tt.fields.ContainsSubTree(&got)
			if err != nil {
				t.Errorf("DualTree.GetRandomSubTreeAtDepth() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !b {
				t.Errorf("DualTree.GetRandomSubTreeAtDepth() error = %v, wantErr %v", err, tt.wantErr)
			}
			got.Print()
		})
	}
}

func TestDualTree_DepthAt(t *testing.T) {
	tests := []struct {
		name    string
		fields  *DualTree
		depth   int
		want    []*DualTreeNode
		wantErr bool
	}{
		{"nil root", TreeNil(), 0, nil, true},
		{"negative depth", TreeT_0(), -1, nil, true},
		{"T", TreeT_0(), 0, []*DualTreeNode{X1.ToDualTreeNode(RandString(5))}, false},
		{"T-NT-T", TreeT_NT_T_1(), 0, []*DualTreeNode{Add.ToDualTreeNode(RandString(5))}, false},
		{"T-NT-T", TreeT_NT_T_1(), 1, []*DualTreeNode{Const8.ToDualTreeNode(RandString(5)),
			X1.ToDualTreeNode(RandString(5))}, false},
		{"T-NT-NT-NT-T", TreeT_NT_T_NT_T_0(), 1, []*DualTreeNode{Sub.ToDualTreeNode(RandString(5)), Const4.ToDualTreeNode(RandString(5))}, false},
		{"T-NT-NT-NT-T", TreeT_NT_T_NT_T_0(), 2, []*DualTreeNode{X1.ToDualTreeNode(RandString(5)), X1.ToDualTreeNode(RandString(5))}, false},
		{"Vine-NT-NT-NT-T", TreeVine_D3(), 2, []*DualTreeNode{
			Sin.ToDualTreeNode(RandString(5))}, false},
		{"Vine-NT-NT-NT-T", TreeVine_D3(), 3, []*DualTreeNode{
			X1.ToDualTreeNode(RandString(5))}, false},
		{"Vine-NT-NT-NT-T", TreeVine_D6_R(), 0, []*DualTreeNode{Add.ToDualTreeNode(RandString(5))}, false},
		{"Vine-6", TreeVine_D6_R(), 4, []*DualTreeNode{Sin.ToDualTreeNode(RandString(5)), X1.ToDualTreeNode(RandString(5))}, false},
		{"Vine-6v2", TreeVine_D6_R(), 6, []*DualTreeNode{X1.ToDualTreeNode(RandString(5))}, false},
		{"Vine-5", TreeVine_D5_R(), 5, []*DualTreeNode{X1.ToDualTreeNode(RandString(5))}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DualTree{
				root: tt.fields.root,
				lock: tt.fields.lock,
			}
			got, err := d.DepthAt(tt.depth)
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.DepthAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			contain, err := contains(got, tt.want)
			if err != nil {
				t.Error(err)
			}
			if !contain {
				t.Errorf("DualTree.DepthAt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDualTree_Search(t *testing.T) {
	tests := []struct {
		name       string
		fields     *DualTree
		key        string
		wantNode   *DualTreeNode
		wantParent *DualTreeNode
		wantErr    bool
	}{
		{"nil", TreeNil(), "", nil, nil, true},
		{"T", TreeT_0(), "0", X1.ToDualTreeNode("0"), nil, false},
		{"T-NT-T", TreeT_NT_T_0(), "2", X1.ToDualTreeNode("2"), Mult.ToDualTreeNode("1"), false},
		{"TreeT_NT_T_NT_T_0", TreeT_NT_T_NT_T_0(), "4", X1.ToDualTreeNode("4"), Sub.ToDualTreeNode("2"), false},
		{"TreeVine_D6_R", TreeVine_D6_R(), "10", X1.ToDualTreeNode("10"), Sin.ToDualTreeNode("9"), false},
		{"TreeT_NT_T_NT_T_0", TreeVine_D6_R(), "14", nil, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := &DualTree{
				root: tt.fields.root,
				lock: tt.fields.lock,
			}
			gotNode, gotParent, err := bst.Search(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !gotNode.IsEqual(tt.wantNode) {
				t.Errorf("DualTree.Search() gotNode = %v, want %v", gotNode, tt.wantNode)
			}
			if !gotParent.IsEqual(tt.wantParent) {
				t.Errorf("DualTree.Search() gotParent = %v, want %v", gotParent, tt.wantParent)
			}
		})
	}
}

func TestDualTree_Replace(t *testing.T) {
	//t0 := TreeT_0()
	type fields struct {
		root *DualTreeNode
		lock sync.RWMutex
	}
	type args struct {
		node     *DualTreeNode
		replacer DualTreeNode
	}
	tests := []struct {
		name       string
		fields     *DualTree
		args       args
		wantHobo   DualTreeNode
		wantParent *DualTreeNode
		wantErr    bool
	}{
		{"nil-Tree", TreeNil(), args{nil, DualTreeNode{}}, DualTreeNode{}, nil, true},
		{"nil-node", TreeT_0(), args{nil, DualTreeNode{}}, DualTreeNode{}, nil, true},
		{"nil-replacer", TreeT_0(), args{X1.ToDualTreeNode("0"), DualTreeNode{}}, DualTreeNode{}, nil, true},
		{"cannot find node to swap", TreeT_0(), args{X1.ToDualTreeNode("30"), DualTreeNode{value: "3"}},
			DualTreeNode{}, nil,
			true},
		{"T", TreeT_0(), args{X1.ToDualTreeNode("0"), DualTreeNode{value: "3", key: "123414"}},
			*X1.ToDualTreeNode("0"), &DualTreeNode{value: "3"},
			false},
		{"T-NT-T", TreeT_NT_T_0(), args{X1.ToDualTreeNode("2"), DualTreeNode{value: "y", arity: 0, key: "123414"}},
			*X1.ToDualTreeNode("2"), &DualTreeNode{value: "*", arity: 2},
			false},
		{"T-NT-T + T-NT-T", TreeT_NT_T_0(), args{X1.ToDualTreeNode("2"), *TreeT_NT_T_0().root},
			*X1.ToDualTreeNode("2"), &DualTreeNode{value: "*", arity: 2},
			false},
		{"T-NT-T + T at root", TreeT_NT_T_0(), args{Mult.ToDualTreeNode("1"), *X1.ToDualTreeNode("4")},
			*Mult.ToDualTreeNode("1"), X1.ToDualTreeNode("4"),
			false},
		{"T-NT-T + T-NT-T at root", TreeT_NT_T_0(), args{Mult.ToDualTreeNode("1"), *TreeT_NT_T_1().root},
			*Mult.ToDualTreeNode("1"), Add.ToDualTreeNode("1234"),
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHobo, gotParent, err := tt.fields.Replace(tt.args.node, tt.args.replacer)
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.Replace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				equal := gotHobo.IsValEqual(tt.args.node)
				if !equal {
					t.Errorf("DualTree.Replace() Hobo node doesnt have the same value is nil | goParentKey:  gotHobo"+
						" = %v, want %v", gotParent,
						tt.wantParent)
				}
				if gotParent == nil {
					t.Errorf("DualTree.Replace() Parent is nil | goParentKey:  gotParent = %v, want %v", gotParent,
						tt.wantParent)
					return
				}
				if gotParent.key == "" {
					t.Errorf("DualTree.Replace() | goParentKey:  gotParent = %v, want %v", gotParent, tt.wantParent)
				}
				if gotParent.value != tt.wantParent.value {
					t.Errorf("DualTree.Replace() | goParentKey:  gotParent = %v, want %v", gotParent, tt.wantParent)
				}
				if gotParent.arity != tt.wantParent.arity {
					t.Errorf("DualTree.Replace() | goParentKey:  gotParent = %v, want %v", gotParent, tt.wantParent)
				}
			}
		})
	}
}

func TestDualTree_RandomLeafAware(t *testing.T) {
	tests := []struct {
		name       string
		fields     *DualTree
		wantNode   *DualTreeNode
		wantParent *DualTreeNode
		wantErr    bool
	}{
		//{"nil", TreeNil(), nil,nil, true },
		{"T", TreeT_0(), X1.ToDualTreeNode("0"), nil, false},
		{"T-NT-T", TreeT_0(), X1.ToDualTreeNode("0"), nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := &DualTree{
				root: tt.fields.root,
				lock: tt.fields.lock,
			}
			gotNode, gotParent, err := bst.RandomTerminalAware()
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.RandomTerminalAware() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNode, tt.wantNode) {
				t.Errorf("DualTree.RandomTerminalAware() gotNode = %v, want %v", gotNode, tt.wantNode)
			}
			if !reflect.DeepEqual(gotParent, tt.wantParent) {
				t.Errorf("DualTree.RandomTerminalAware() gotParent = %v, want %v", gotParent, tt.wantParent)
			}
		})
	}
}

func TestDualTree_Clone(t *testing.T) {
	type fields struct {
		root *DualTreeNode
		lock sync.RWMutex
	}
	tests := []struct {
		name   string
		fields *DualTree
	}{
		{"nil", TreeNil()},
		{"T", TreeT_0()},
		{"T-NT-T", TreeT_NT_T_0()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := DualTree{
				root: tt.fields.root,
				lock: tt.fields.lock,
			}
			got, err := bst.Clone()
			if err != nil {
				t.Error(err)
			}
			if tt.fields.root == nil && got.root == nil {
				return
			}
			b, err := got.ContainsSubTree(tt.fields)
			if err != nil {
				t.Error(err)
			}
			if b {
				t.Errorf("DualTree.Clone() They contain = %#v, want %#v", got, tt.fields)
			}
			if &got == tt.fields {
				t.Errorf("DualTree.Clone() = %#v, want %#v", got, tt.fields)
			}
		})
	}
}

func TestDualTree_GetShortestBranch(t *testing.T) {
	tests := []struct {
		name                   string
		fields                 *DualTree
		minAcceptableDepth     int
		wantShortestNode       *DualTreeNode
		wantShortestNodeParent *DualTreeNode
		wantShortestDepth      int
		wantErr                bool
	}{
		{"nil", TreeNil(), 2, nil, nil, -1, true},
		{"negative minAcceptableDepth", TreeT_0(), -1, nil, nil, -1, true},
		{"T", TreeT_0(), 0, TreeT_0().root, nil, 0, false},
		{"T-NT-T", TreeT_NT_T_1(), 1, TreeT_NT_T_1().root.left, TreeT_NT_T_1().root, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := &DualTree{
				root: tt.fields.root,
				lock: tt.fields.lock,
			}
			gotShortestNode, gotShortestNodeParent, gotShortestDepth, err := bst.GetShortestBranch(tt.minAcceptableDepth)
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.GetShortestBranch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotShortestNode, tt.wantShortestNode) {
				t.Errorf("DualTree.GetShortestBranch() gotShortestNode = %v, want %v", gotShortestNode, tt.wantShortestNode)
			}
			if !reflect.DeepEqual(gotShortestNodeParent, tt.wantShortestNodeParent) {
				t.Errorf("DualTree.GetShortestBranch() gotShortestNodeParent = %v, want %v", gotShortestNodeParent, tt.wantShortestNodeParent)
			}
			if gotShortestDepth != tt.wantShortestDepth {
				t.Errorf("DualTree.GetShortestBranch() gotShortestDepth = %v, want %v", gotShortestDepth, tt.wantShortestDepth)
			}
		})
	}
}

//func TestDualTree_InOrderTraverseDepthAware(t *testing.T) {
//	tests := []struct {
//		name   string
//		fields *DualTree
//		f func(node *DualTreeNode, parentNode *DualTreeNode, depth *int)
//	}{
//		{"nil", TreeNil(), },
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			bst := &DualTree{
//				root: tt.fields.root,
//				lock: tt.fields.lock,
//			}
//			bst.InOrderTraverseDepthAware(tt.f)
//		})
//	}
//}

func TestGenerateRandomTreeEnforceIndependentVariable(t *testing.T) {
	type args struct {
		depth          int
		independentVar SymbolicExpression
		terminals      []SymbolicExpression
		nonTerminals   []SymbolicExpression
	}
	tests := []struct {
		name    string
		args    args
		want    *DualTree
		wantErr bool
	}{
		{"err-lowMaxDepth", args{-1, SymbolicExpression{}, nil, nil}, nil, true},
		{"err-empty-indepvar", args{-1, SymbolicExpression{}, nil, nil}, nil, true},
		{"err-nil-terminals", args{2, X1, nil, nil}, nil, true},
		{"err-nil-non-terminals", args{2, X1, []SymbolicExpression{X1}, nil}, nil, true},
		{"err-nil-empty-terminals", args{2, X1, []SymbolicExpression{}, []SymbolicExpression{}}, nil, true},
		{"err-nil-empty-nonterminals", args{2, X1, []SymbolicExpression{X1}, []SymbolicExpression{}}, nil, true},
		{"err-nil-empty-nonterminals", args{2, X1, []SymbolicExpression{X1}, []SymbolicExpression{}}, nil, true},
		{"T", args{0, X1, []SymbolicExpression{X1}, []SymbolicExpression{}}, TreeT_1(), false},
		{"T", args{0, X1, []SymbolicExpression{X1, Const0, Const1, Const2, Const3, Const4, Const5, Const6,
			Const7, Const8, Const9}, []SymbolicExpression{}}, TreeT_1(), false},
		{"err-depth-1-no-NT", args{1, X1, []SymbolicExpression{X1}, []SymbolicExpression{}}, TreeT_1(), true},
		{"depth-1", args{1, X1, []SymbolicExpression{X1}, []SymbolicExpression{Add}}, TreeT_NT_T_0(), false},
		{"depth-2", args{2, X1, []SymbolicExpression{X1}, []SymbolicExpression{Add}}, TreeT_NT_T_NT_T_NT_T_0(), false},
		{"depth-3-diverse", args{2, X1, []SymbolicExpression{X1, Const0, Const1, Const2, Const3, Const4, Const5, Const6,
			Const7, Const8, Const9},
			[]SymbolicExpression{Add,
				Mult,
				Sub}},
			TreeT_NT_T_NT_T_NT_T_0(),
			false},
		{"depth-2-diverse", args{3, X1, []SymbolicExpression{X1, Const0, Const1, Const2, Const3, Const4, Const5, Const6,
			Const7, Const8, Const9},
			[]SymbolicExpression{Add,
				Mult,
				Sub}},
			TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0(),
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			got, err := GenerateRandomTreeEnforceIndependentVariable(tt.args.depth, tt.args.independentVar,
				tt.args.terminals,
				tt.args.nonTerminals)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRandomTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				wantExpressionSet, err := tt.want.ToSymbolicExpressionSet()
				if err != nil {
					t.Error(err)
				}
				gotExpressionSet, err := got.ToSymbolicExpressionSet()
				if err != nil {
					t.Error(err)
				}

				if len(wantExpressionSet) != len(gotExpressionSet) {
					t.Errorf("They are not the same length error = %v, wantErr %v", gotExpressionSet, wantExpressionSet)
				}

				for e := range gotExpressionSet {
					if gotExpressionSet[e].value == tt.args.independentVar.value {
						break
					}
					t.Errorf("Does not contain independent variable = %v, "+
						"wantErr %v", gotExpressionSet,
						wantExpressionSet)
				}
			}
		})
	}
}

func TestDualTree_DeleteNonTerminal(t *testing.T) {
	tests := []struct {
		name    string
		fields  *DualTree
		wantErr bool
	}{
		{"nil", TreeNil(), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := &DualTree{
				root: tt.fields.root,
				lock: tt.fields.lock,
			}
			if err := bst.DeleteNonTerminal(); (err != nil) != tt.wantErr {
				t.Errorf("DualTree.DeleteNonTerminal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDualTree_RandomBranchAware(t *testing.T) {
	type fields struct {
		root *DualTreeNode
		lock sync.RWMutex
	}
	tests := []struct {
		name       string
		fields     fields
		wantNode   *DualTreeNode
		wantParent *DualTreeNode
		wantErr    bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := &DualTree{
				root: tt.fields.root,
				lock: tt.fields.lock,
			}
			gotNode, gotParent, err := bst.RandomNonTerminalAware()
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.RandomNonTerminalAware() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNode, tt.wantNode) {
				t.Errorf("DualTree.RandomNonTerminalAware() gotNode = %v, want %v", gotNode, tt.wantNode)
			}
			if !reflect.DeepEqual(gotParent, tt.wantParent) {
				t.Errorf("DualTree.RandomNonTerminalAware() gotParent = %v, want %v", gotParent, tt.wantParent)
			}
		})
	}
}

func TestDualTree_GetBranchesAware(t *testing.T) {
	tests := []struct {
		name    string
		fields  *DualTree
		want    []AwareTree
		wantErr bool
	}{
		{"nil", TreeNil(), nil, true},
		{"T", TreeT_0(), []AwareTree{}, false},
		{"T-NT-T", TreeT_NT_T_0(), []AwareTree{{parent: nil, node: Mult.ToDualTreeNode("1")}}, false},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_0(), []AwareTree{
			{parent: Mult.ToDualTreeNode("0"), node: Sub.ToDualTreeNode("2")},
			{parent: nil, node: Mult.ToDualTreeNode("0")}},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := &DualTree{
				root: tt.fields.root,
				lock: tt.fields.lock,
			}
			got, err := bst.GetNonTerminalsAware()
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.GetNonTerminalsAware() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				for o := range got {
					if !got[o].node.IsEqual(tt.want[o].node) {
						t.Errorf("DualTree.GetNonTerminalsAware() = %v, want %v", got[o].node, tt.want[o].node)
					}
					if got[o].parent == nil && tt.want[o].parent != nil {
						t.Errorf("DualTree.GetNonTerminalsAware() = %v, want %v", got[o].parent, tt.want[o].parent)
					}
					if got[o].parent != nil && tt.want[o].parent == nil {
						t.Errorf("DualTree.GetNonTerminalsAware() = %v, want %v", got[o].parent, tt.want[o].parent)
					}
					if got[o].parent != nil && tt.want[o].parent != nil {
						if !got[o].parent.IsEqual(tt.want[o].parent) || !got[o].parent.IsEqual(tt.want[o].
							parent) {
							t.Errorf("DualTree.GetNonTerminalsAware() = %v, want %v", got[o].parent, tt.want[o].parent)
						}
					}
				}
			}
		})
	}
}

func TestDualTree_InOrderTraverseAware(t *testing.T) {
	type args struct {
		f func(node *DualTreeNode, parentNode *DualTreeNode)
	}
	tests := []struct {
		name   string
		fields *DualTree
		args   args
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := &DualTree{
				root: tt.fields.root,
				lock: tt.fields.lock,
			}
			bst.InOrderTraverseAware(tt.args.f)
		})
	}
}

func Test_inOrderTraverseAwareDepth(t *testing.T) {
	type args struct {
		n            *DualTreeNode
		parent       *DualTreeNode
		depth        *int
		shouldReturn *bool
		f            func(node *DualTreeNode, parentNode *DualTreeNode, depth *int, shouldReturn *bool)
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inOrderTraverseAwareDepth(tt.args.n, tt.args.parent, tt.args.depth, tt.args.shouldReturn, tt.args.f)
		})
	}
}

func Test_inOrderTraverseAware(t *testing.T) {
	type args struct {
		n      *DualTreeNode
		parent *DualTreeNode
		f      func(node *DualTreeNode, parentNode *DualTreeNode)
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inOrderTraverseAware(tt.args.n, tt.args.parent, tt.args.f)
		})
	}
}

func TestDualTree_GetLeafsAware(t *testing.T) {
	tests := []struct {
		name    string
		fields  *DualTree
		want    []AwareTree
		wantErr bool
	}{
		{"nil", TreeNil(), nil, true},
		{"T", TreeT_0(), []AwareTree{{node:X1.ToDualTreeNode("0"), parent: nil}}, false},
		{"T-NT-T", TreeT_NT_T_0(), []AwareTree{
			{node:X1.ToDualTreeNode("2"), parent: Mult.ToDualTreeNode("1")},
			{node:Const4.ToDualTreeNode("3"), parent: Mult.ToDualTreeNode("1")},
		}, false},
		{"T-NT-T", TreeT_NT_T_NT_T_0(), []AwareTree{
			{node:X1.ToDualTreeNode("3"), parent: Sub.ToDualTreeNode("2")},
			{node:X1.ToDualTreeNode("4"), parent: Sub.ToDualTreeNode("2")},
			{node:Const4.ToDualTreeNode("1"), parent: Mult.ToDualTreeNode("0")},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bst := &DualTree{
				root: tt.fields.root,
				lock: tt.fields.lock,
			}
			got, err := bst.GetTerminalsAware()
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTree.GetTerminalsAware() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				for o := range got {
					if !got[o].node.IsEqual(tt.want[o].node) {
						t.Errorf("DualTree.GetNonTerminalsAware() = %v, want %v", got[o].node, tt.want[o].node)
					}
					if got[o].parent == nil && tt.want[o].parent != nil {
						t.Errorf("DualTree.GetNonTerminalsAware() = %v, want %v", got[o].parent, tt.want[o].parent)
					}
					if got[o].parent != nil && tt.want[o].parent == nil {
						t.Errorf("DualTree.GetNonTerminalsAware() = %v, want %v", got[o].parent, tt.want[o].parent)
					}
					if got[o].parent != nil && tt.want[o].parent != nil {
						if !got[o].parent.IsEqual(tt.want[o].parent) || !got[o].parent.IsEqual(tt.want[o].
							parent) {
							t.Errorf("DualTree.GetNonTerminalsAware() = %v, want %v", got[o].parent, tt.want[o].parent)
						}
					}
				}
				if len(got) < 1 {
					t.Errorf("DualTree.GetNonTerminalsAware() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
