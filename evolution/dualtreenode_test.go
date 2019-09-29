package evolution

import (
	"reflect"
	"testing"
)

func Test_arityRemainder(t *testing.T) {
	tests := []struct {
		name string
		tree *DualTree
		want int
	}{
		{"full - NT(2)", TreeT_NT_T_0(), 0},
		{"full - NT(2)", TreeT_NT_T_NT_T_0(), 0},
		{"full - NT(1)", Tree3(), 0},
		{"full - NT(1)", Tree4(), 0},
		{"full - NT(1)", Tree5(), 1},
		{"half - NT(1)", Tree6(), 1},
		{"empty - NT(2)", Tree7(), 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tree.root.ArityRemainder(); got != tt.want {
				t.Errorf("ArityRemainder() = %v, isEqual %v", got, tt.want)
			}
		})
	}
}

func TestDualTreeNode_IsLeaf(t *testing.T) {
	tests := []struct {
		name string
		node *DualTreeNode
		want bool
	}{
		{"nil", &DualTreeNode{RandString(5), "", nil, nil, 0}, true},
		{"root", TreeT_NT_T_0().root, false},
		{"x in x * 4", TreeT_NT_T_0().root.left, true},
		{"4 in x * 4", TreeT_NT_T_0().root.right, true},
		{"4 in x - x * 4", TreeT_NT_T_NT_T_0().root.right, true},
		{"x in x - x * 4", TreeT_NT_T_NT_T_0().root.left.left, true},
		{"sin in sin", Tree5().root, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.node.IsLeaf(); got != tt.want {
				t.Errorf("IsLeaf() = %v, isEqual %v", got, tt.want)
			}
		})
	}
}

func TestDualTreeNode_IsValEqual(t *testing.T) {
	tests := []struct {
		name   string
		fields *DualTreeNode
		args   *DualTreeNode
		want   bool
	}{
		{"", &DualTreeNode{}, &DualTreeNode{}, true},
		{"value", &DualTreeNode{value: "x"}, &DualTreeNode{value: "x"}, true},
		{"diff-arity-same-val", &DualTreeNode{value: "x", arity: 2}, &DualTreeNode{value: "x", arity: 1}, true},
		{"diff-arity-same-val", &DualTreeNode{value: "*", arity: 2}, &DualTreeNode{value: "x", arity: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.IsValEqual(tt.args); got != tt.want {
				t.Errorf("DualTreeNode.IsValEqual() = %v, isEqual %v", got, tt.want)
			}
		})
	}
}

func TestDualTreeNode_IsEqual(t *testing.T) {
	tests := []struct {
		name     string
		treeNode *DualTreeNode
		subTree  *DualTreeNode
		isEqual  bool
	}{
		{"nil", &DualTreeNode{}, &DualTreeNode{}, true},
		{"value", &DualTreeNode{value: "x"}, &DualTreeNode{value: "x"}, true},
		{"same-val-same-left", &DualTreeNode{value: "x"}, &DualTreeNode{value: "x",
			left: Add.ToDualTreeNode("1"), key: "123"},
			false},
		{"same-val-same-right", &DualTreeNode{value: "x", right: Add.ToDualTreeNode("1")}, &DualTreeNode{value: "x",
			right: Add.ToDualTreeNode("1")},
			true},
		{"diff-arity-same-val", &DualTreeNode{value: "x", arity: 2}, &DualTreeNode{value: "x", arity: 1}, false},
		{"diff-arity-same-val", &DualTreeNode{value: "*", arity: 2}, &DualTreeNode{value: "x", arity: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.treeNode.IsEqual(tt.subTree); got != tt.isEqual {
				t.Errorf("DualTreeNode.IsEqual() = %v, isEqual %v", got, tt.isEqual)
			}
		})
	}
}

func TestDualTreeNode_Clone(t *testing.T) {
	tests := []struct {
		name   string
		fields DualTreeNode
		want   DualTreeNode
	}{
		{"nil", DualTreeNode{}, DualTreeNode{}},
		{"const1", *Const1.ToDualTreeNode("123"), *Const1.ToDualTreeNode("123")},
		{"const4", *Const4.ToDualTreeNode("234"), *Const4.ToDualTreeNode("234")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DualTreeNode{
				key:   tt.fields.key,
				value: tt.fields.value,
				left:  tt.fields.left,
				right: tt.fields.right,
				arity: tt.fields.arity,
			}
			var got DualTreeNode = d.Clone()

			if got.key == d.key {
				t.Errorf("DualTreeNode.Clone() = %v, want %v", got, tt.want)
			}
			if &got == &tt.fields {
				t.Errorf("DualTreeNode.Clone() = %v, want %v | same address", got, tt.want)
			}
		})
	}
}

func TestDualTreeNode_ToDualTree(t *testing.T) {

	tests := []struct {
		name    string
		fields  *DualTreeNode
		want    DualTree
		wantErr bool
	}{
		{"T", TreeT_0().root, *TreeT_0(), false},
		{"T-NT-T", TreeT_NT_T_0().root, *TreeT_NT_T_0(), false},
		{"T-NT-T-NT-T", TreeT_NT_T_NT_T_0().root.left, *TreeT_NT_T_5(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DualTreeNode{
				key:   tt.fields.key,
				value: tt.fields.value,
				left:  tt.fields.left,
				right: tt.fields.right,
				arity: tt.fields.arity,
			}
			got, err := d.ToDualTree()
			if (err != nil) != tt.wantErr {
				t.Errorf("DualTreeNode.ToDualTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DualTreeNode.ToDualTree() = %v, want %v", got, tt.want)
			}
		})
	}
}
