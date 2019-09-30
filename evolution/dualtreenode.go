package evolution

import "fmt"

// DualTreeNode represents a a treeNode with a maximum of two children.
// It is not technically a binary treeNode as it DOES not place any ordering on left and right children as binary trees
// prototypically do.
type DualTreeNode struct {
	key   string
	value string
	left  *DualTreeNode //left
	right *DualTreeNode //right
	arity int
}

// IsEqual checks to see if all aspects of a DualTreeNode are equivalent. This includes value as well as pointers
func (b *DualTreeNode) IsEqual(t *DualTreeNode) bool {
	if b == nil && t == nil {
		return true
	}
	if b.key != t.key {
		return false
	}
	if b.value != t.value {
		return false
	}
	if b.arity != t.arity {
		return false
	}
	return true
}

// IsValEqual is a simple check to see if values of strings in the nodes are equal
func (d *DualTreeNode) IsValEqual(t *DualTreeNode) bool {
	if d.value == t.value {
		return true
	}
	return false
}

// IsLeaf checks to see if a given node is a leaf
func (d *DualTreeNode) IsLeaf() bool {
	if d.arity == 0 {
		if d.right == nil || d.left == nil {
			return true
		}
	}
	return false
}

// ArityRemainder calculates the remaining available node connections based on arity for a given root node.
// This is used to balance the NonTerminals and the Terminals depending on their requirements.
func (d *DualTreeNode) ArityRemainder() int {
	available := d.arity
	if d.arity == 2 {
		if d.right != nil {
			available--
		}
		if d.left != nil {
			available--
		}
		return available
	} else if d.arity == 1 {
		if d.left != nil {
			available--
		}
		return available
	}
	return 0
}

// IsLeaf checks to see if a given node is a leaf
func (d *DualTreeNode) ToSymbolicExpression() (SymbolicExpression, error) {
	err := d.isValid()
	if err != nil {
		return SymbolicExpression{}, err
	}

	kind := 0
	if d.arity == 2 {
		kind = 1
	}
	return SymbolicExpression{
		arity: d.arity,
		value: d.value,
		kind:  kind,
	}, err
}

// ToDualTree takes a given node and returns a treeNode from it by following the path.
func (d *DualTreeNode) ToDualTree() (DualTree, error) {
	err := d.isValid()
	if err != nil {
		return DualTree{}, err
	}
	return DualTree{
		root: d,
	}, err
}

func (d *DualTreeNode) isValid() error {
	if d.key == "" {
		return fmt.Errorf("ToDualTree | key is empty")
	}
	if d.value == "" {
		return fmt.Errorf("ToDualTree | value is empty")
	}
	return nil
}

// Clone performs an O(N) deep clone of a given DualTreeNode and returns a new DualTreeNode,
// granted no errors are present.
func (d DualTreeNode) Clone() DualTreeNode {
	d.key = RandString(5)
	return d
}
