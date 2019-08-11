package dualtree

import (
	"reflect"
	"testing"
)

//
//var (
//	nodeEmpty  = &DualTreeNode{}
//	tree1node  = &DualTreeNode{right: tree1node2, left: tree1node1, parent: nil, item:"*"}
//	tree1node1 = &DualTreeNode{right: nil, left: nil, parent: tree1node, item:"1"}
//	tree1node2 = &DualTreeNode{right: nil, left: nil, parent: tree1node, item:"2"}
//)
//
//func TestInorderDFS(t *testing.T) {
//	type tree struct {
//		root          *DualTreeNode
//		dualTreeNodes []*DualTreeNode
//	}
//
//	tt := []struct {
//		name     string
//		tree     tree
//		expected []*DualTreeNode
//	}{
//		//{"nil root", tree{nil, nil}, nil},
//		//{"nil root init arr", tree{nil, make([]*DualTreeNode, 0)},[]*DualTreeNode{}},
//		////{"single root", tree{nodeEmpty, make([]*DualTreeNode, 0)},[]*DualTreeNode{nodeEmpty}},
//		//{"basic bin tree", tree{tree1node, make([]*DualTreeNode, 0)},[]*DualTreeNode{tree1node1, tree1node,
//		//	tree1node2}},
//	}
//
//	for _, v := range tt {
//		t.Run(v.name, func(t *testing.T) {
//			InorderDFS(v.tree.root, v.tree.dualTreeNodes)
//			if !reflect.DeepEqual(v.expected, v.tree.dualTreeNodes) {
//				t.Errorf("\n\nwant: %#v\ngot %v", v.expected, v.tree.dualTreeNodes)
//			}
//		})
//	}
//}

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

func TestDualTree_FromNodeTypes(t *testing.T) {
	//TERMINALS
	x1 := NodeType{kind: 0, value: "x"}
	//x2 := NodeType{kind: 0, value: "x"}
	const1 := NodeType{kind: 0, value: "4"}
	//const2 := NodeType{kind: 0, value: "8"}
	invalidT := NodeType{kind: -3, value: "12"}

	// NON-TERMINALS
	//sub := NodeType{kind: 1, value: "-", arity:2}
	//add := NodeType{kind: 1, value: "+", arity:2}
	mult := NodeType{kind: 1, value: "*", arity: 2}
	//sinx := NodeType{kind: 1, value: "sin"}
	invalidNT := NodeType{kind: -3, value: "="}

	type args struct {
		E []NodeType
	}
	tests := []struct {
		name    string
		fields  *DualTree
		args    []NodeType
		wantErr bool
	}{
		{"nil nodetypes", &DualTree{}, nil, true},
		{"empty nodetypes", &DualTree{}, make([]NodeType, 0), true},
		{"T", &DualTree{}, []NodeType{x1}, false},
		{"NT", &DualTree{}, []NodeType{mult}, true},
		{"Invalid-NT", &DualTree{}, []NodeType{invalidNT}, true},
		{"Invalid-T", &DualTree{}, []NodeType{invalidT}, true},
		{"T-T", &DualTree{}, []NodeType{x1, const1}, true},
		//{"NT-NT", &DualTree{}, []NodeType{mult, sub}, true},
		//{"T-NT(1)", &DualTree{}, []NodeType{x1, sinx}, false},
		//{"T-NT(2)", &DualTree{}, []NodeType{x2, sub}, true},
		//{"T-NT(2)-T", &DualTree{}, []NodeType{x1, add, x2}, false},
		//{"T-NT(2)-T-NT(2)-T", &DualTree{}, []NodeType{x1, add, x2, mult, const1} ,false},
		//{"T-NT(2)-T-NT(1)-T", &DualTree{}, []NodeType{x1, add, x2, sinx, const1}, true},
		//{"T-NT(2)-T-NT(1)", &DualTree{}, []NodeType{x1, add, x2, sinx}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.FromNodeTypes(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("DualTree.FromNodeTypes() error = %v, wantErr %v", err, tt.wantErr)
			}

			expected := make([]string, 0)
			for i := range tt.args {
				if len(tt.args) == 1 && tt.args[0].kind > 0 {
					continue
				}
				if len(tt.args) == 1 && tt.args[0].kind < 0 {
					continue
				}
				expected = append(expected, tt.args[i].value)
			}

			got := make([]string, 0)
			tt.fields.InOrderTraverse(func(s string) {
				got = append(got, s)
			})

			if !reflect.DeepEqual(expected, got) {
				t.Errorf("expected: %#v ||| got: %#v", expected, got)
			}
		})
	}

	//tree := DualTree{}
	//
	//
	//nodeTypes := []NodeType{x1, mult, x2, add, const1}
	//
	//err := tree.FromNodeTypes(nodeTypes)
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//tree.String()
	//
	//tree.InOrderTraverse(func(s string) {
	//	log.Print(s)
	//})
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
		{"T", &DualTree{}, []NodeType{x1}, false},
		{"NT", &DualTree{}, []NodeType{mult}, true},
		{"T-T", &DualTree{}, []NodeType{x1, const1}, true},
		{"NT-NT", &DualTree{}, []NodeType{mult, sub}, true},
		{"T-NT(1)", &DualTree{}, []NodeType{x1, sin}, false},
		{"T-NT(2)", &DualTree{}, []NodeType{x1, sub}, true},
		{"T-NT(2)-T", &DualTree{}, []NodeType{x1, add, const1}, false},
		{"T-NT(2)-T-NT(2)-T", &DualTree{}, []NodeType{x1, add, const2, mult, const1} ,false},
		{"T-NT(2)-T-NT(1)-T", &DualTree{}, []NodeType{x1, add, const1, sin, const1}, true},
		{"T-NT(2)-T-NT(1)", &DualTree{}, []NodeType{x1, mult, const1, sub, const2, sin}, false},
		{"T-NT(1)-NT(1)-NT(1)-NT(1)", &DualTree{}, []NodeType{x1, sin, sin, sin, sin}, false},
		{"T-NT(1)-NT(2)-T-NT(1)", &DualTree{}, []NodeType{x1, sin, add, const2, sin}, false},
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
				tt.fields.InOrderTraverse(func(s string) {
					got = append(got, s)
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
		{"full - NT(2)", tree1(), 0},
		{"full - NT(2)", tree2(), 0},
		{"full - NT(1)", tree3(), 0},
		{"full - NT(1)", tree4(), 0},
		{"full - NT(1)", tree5(), 1},
		{"half - NT(1)", tree6(), 1},
		{"empty - NT(2)", tree7(), 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := arityRemainder(tt.tree.root); got != tt.want {
				t.Errorf("arityRemainder() = %v, want %v", got, tt.want)
			}
		})
	}
}

//TERMINALS
var x1 = NodeType{kind: 0, value: "x", arity: 0}
var const1 = NodeType{kind: 0, value: "4", arity: 0}
var const2 = NodeType{kind: 0, value: "8", arity: 0}

// NON-TERMINALS
var sub = NodeType{kind: 1, value: "-", arity: 2}
var add = NodeType{kind: 1, value: "+", arity: 2}
var mult = NodeType{kind: 1, value: "*", arity: 2}
var sin = NodeType{kind: 1, value: "sin", arity: 1}

// SAMPLE TREES
// tree0 = x
var tree0 = func() *DualTree {
	t := DualTree{}
	t.root = x1.ToDualTreeNode(0)
	return &t
}

// tree1 = x * 4
var tree1 = func() *DualTree {
	t := DualTree{}
	t.root = mult.ToDualTreeNode(0)
	t.root.left = x1.ToDualTreeNode(1)
	t.root.right = const1.ToDualTreeNode(2)
	return &t
}

// tree2 = x - x * 4
var tree2 = func() *DualTree {
	t := DualTree{}
	t.root = mult.ToDualTreeNode(0)
	t.root.right = const1.ToDualTreeNode(4)
	t.root.left = sub.ToDualTreeNode(1)
	t.root.left.left = x1.ToDualTreeNode(2)
	t.root.left.right = x1.ToDualTreeNode(3)
	return &t
}

// tree3 = sin(4 - x)
var tree3 = func() *DualTree {
	t := DualTree{}
	t.root = sin.ToDualTreeNode(0)
	t.root.left = sub.ToDualTreeNode(1)
	t.root.left.left = x1.ToDualTreeNode(2)
	t.root.left.right = const1.ToDualTreeNode(3)
	return &t
}

// tree4 = sin(x)
var tree4 = func() *DualTree {
	t := DualTree{}
	t.root = sin.ToDualTreeNode(0)
	t.root.left = x1.ToDualTreeNode(1)
	return &t
}

// tree5 = sin
var tree5 = func() *DualTree {
	t := DualTree{}
	t.root = sin.ToDualTreeNode(0)
	return &t
}

// tree6 = x +
var tree6 = func() *DualTree {
	t := DualTree{}
	t.root = add.ToDualTreeNode(0)
	t.root.left = x1.ToDualTreeNode(1)
	return &t
}

// tree7 =  +
var tree7 = func() *DualTree {
	t := DualTree{}
	t.root = add.ToDualTreeNode(0)
	return &t
}
