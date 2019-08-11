package dualtree

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
func (b *DualTreeNode) IsEqual(t *DualTreeNode) bool {
	if b.value != t.value {
		return false
	}
	if b.left != t.left {
		return false
	}
	if b.right != t.right {
		return false
	}
	return true
}

// IsValEqual is a simple check to see if values of strings in the nodes are equal
func (b *DualTreeNode) IsValEqual(t *DualTreeNode) bool {
	if b.value == t.value {
		return true
	}
	return false
}

//type Terminal DualTreeNode
//
//type NodeType DualTreeNode
//
