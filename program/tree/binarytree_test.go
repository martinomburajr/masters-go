package tree

import (
	"reflect"
	"testing"
)

var (
	nodeEmpty  = &DualTreeNode{}
	tree1node  = &DualTreeNode{right: tree1node2, left: tree1node1, parent: nil, item:"*"}
	tree1node1 = &DualTreeNode{right: nil, left: nil, parent: tree1node, item:"1"}
	tree1node2 = &DualTreeNode{right: nil, left: nil, parent: tree1node, item:"2"}
)

func TestInorderDFS(t *testing.T) {
	type args struct {
		root          *DualTreeNode
		dualTreeNodes []*DualTreeNode
	}

	tt := []struct {
		name     string
		args     args
		expected []*DualTreeNode
	}{
		//{"nil root", args{nil, nil}, nil},
		//{"nil root init arr", args{nil, make([]*DualTreeNode, 0)},[]*DualTreeNode{}},
		////{"single root", args{nodeEmpty, make([]*DualTreeNode, 0)},[]*DualTreeNode{nodeEmpty}},
		//{"basic bin tree", args{tree1node, make([]*DualTreeNode, 0)},[]*DualTreeNode{tree1node1, tree1node,
		//	tree1node2}},
	}

	for _, v := range tt {
		t.Run(v.name, func(t *testing.T) {
			InorderDFS(v.args.root, v.args.dualTreeNodes)
			if !reflect.DeepEqual(v.expected, v.args.dualTreeNodes) {
				t.Errorf("\n\nwant: %#v\ngot %v", v.expected, v.args.dualTreeNodes)
			}
		})
	}
}
