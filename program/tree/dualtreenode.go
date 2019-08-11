package tree

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
//// IsEqual checks to see if all aspects of a DualTreeNode are equivalent. This includes value as well as pointers
//func (b *DualTreeNode) IsEqual(t *DualTreeNode) bool {
//	if b.item != t.item {
//		return false
//	}
//	if b.left != t.left {
//		return false
//	}
//	if b.right != t.right {
//		return false
//	}
//	if b.parent != t.parent {
//		return false
//	}
//	return true
//}
//
//// IsValEqual is a simple check to see if values of strings in the nodes are equal
//func (b *DualTreeNode) IsValEqual(t *DualTreeNode) bool {
//	if b.item == t.item {
//		return true
//	}
//	return false
//}
//
//// ParentDualTreeNode represents a DualTree, but with no parent
//type ParentDualTreeNode struct {
//	left  *DualTreeNode
//	right *DualTreeNode
//	item  string
//}


// DualTreeNode a single node that composes the tree
type DualTreeNode struct {
	key   int
	value string
	left  *DualTreeNode //left
	right *DualTreeNode //right
}