package evolution

// DualTreeNode represents a a tree with a maximum of two children.
// It is not technically a binary tree as it DOES not place any ordering on left and right children as binary trees
// prototypically do.
//type DualTreeNode struct {
//	left   *DualTreeNode
//	right  *DualTreeNode
//	parent *DualTreeNode
//	item   string
//}
//

//
//// ParentDualTreeNode represents a DualTree, but with no parent
//type ParentDualTreeNode struct {
//	left  *DualTreeNode
//	right *DualTreeNode
//	item  string
//}

// IsEqual checks to see if all aspects of a DualTreeNode are equivalent. This includes value as well as pointers
//func (b *DualTreeNode) IsEqual(t *DualTreeNode) bool {
//	if b.value != t.value {
//		return false
//	}
//	if b.arity != b.arity {
//		return false
//	}
//	if b.left  != nil && t.left == nil {
//		return false
//	}
//	if b.left  == nil && t.left != nil {
//		return false
//	}
//	if b.left  == nil && t.left == nil {
//		return true
//	}
//	if b.left.IsEqual(t.left) {
//		return false
//	}
//	if b.right.IsEqual(t.right) {
//		return false
//	}
//	return true
//}

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