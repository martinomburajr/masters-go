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
				t.Errorf("ArityRemainder() = %v, wantAntagonist %v", got, tt.want)
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
				t.Errorf("IsLeaf() = %v, wantAntagonist %v", got, tt.want)
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
				t.Errorf("DualTreeNode.IsValEqual() = %v, wantAntagonist %v", got, tt.want)
			}
		})
	}
}

func TestDualTreeNode_IsEqual(t *testing.T) {
	tests := []struct {
		name           string
		treeNode       *DualTreeNode
		subTree        *DualTreeNode
		wantAntagonist bool
	}{
		{"", &DualTreeNode{}, &DualTreeNode{}, true},
		{"value", &DualTreeNode{value: "x"}, &DualTreeNode{value: "x"}, true},
		{"same-val-same-left", &DualTreeNode{value: "x", left: Add.ToDualTreeNode("1")}, &DualTreeNode{value: "x",
			left: Add.ToDualTreeNode("1"), key:"123"},
			true},
		{"same-val-same-right", &DualTreeNode{value: "x", right: Add.ToDualTreeNode("1")}, &DualTreeNode{value: "x",
			right: Add.ToDualTreeNode("1")},
			true},
		{"diff-arity-same-val", &DualTreeNode{value: "x", arity: 2}, &DualTreeNode{value: "x", arity: 1}, false},
		{"diff-arity-same-val", &DualTreeNode{value: "*", arity: 2}, &DualTreeNode{value: "x", arity: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.treeNode.IsEqual(tt.subTree); got != tt.wantAntagonist {
				t.Errorf("DualTreeNode.IsEqual() = %v, wantAntagonist %v", got, tt.wantAntagonist)
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
		{"const1", *Const1.ToDualTreeNode(RandString(5)), *Const1.ToDualTreeNode(RandString(5))},
		{"const4", *Const4.ToDualTreeNode(RandString(5)), *Const4.ToDualTreeNode(RandString(5))},
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
			var got DualTreeNode
			if got = d.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DualTreeNode.Clone() = %v, want %v", got, tt.want)
			}
			if &got == &tt.want {
				t.Errorf("DualTreeNode.Clone() cannot reference the same address = %p, want %p", &got, &tt.want)
			}
		})
	}
}
