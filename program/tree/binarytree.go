package tree

import (
	"fmt"
	"sync"
)

//type DualTree struct {
//	rootNode *DualTreeNode
//}
//
//func (t *DualTree) Traverse(traversalMethod treetraversal.TraversalStrategy) []*DualTreeNode {
//	return nil
//}
//
//// Get obtains a given DualTreeNode from a given index this assumes InorderDFS
//func (t *DualTree) Get(index int) (*DualTreeNode, error) {
//	return nil, fmt.Errorf("")
//}
//
//// Print simply prints the tree
//func (t *DualTree) Init(dualTreeNode *DualTreeNode) {
//	t.rootNode = dualTreeNode
//}
//
//func (t *DualTree) Insert(node DualTreeNode, index int) error {
//	return fmt.Errorf("")
//}
//
//// Print simply prints the tree
//func (t *DualTree) Print() string {
// 	return ""
//}
//
//// inorderDFS is a recursive implementation of the inorderDFS method. The dualTreeNodes array must be initialized
//func InorderDFS(root *DualTreeNode, dualTreeNodes []*DualTreeNode) {
//	if root == nil {
//		return
//	}
//
//	if dualTreeNodes == nil {
//		return
//	}
//
//	InorderDFS(root.left, dualTreeNodes)
//
//	dualTreeNodes = append(dualTreeNodes, root)
//
//	InorderDFS(root.right, dualTreeNodes)
//}

// Package binarysearchtree creates a DualTree data structure for the string type





// DualTree the binary search tree of Items
type DualTree struct {
	root *DualTreeNode
	lock sync.RWMutex
}

// Insert inserts the string t in the tree
func (bst *DualTree) Insert(key int, value string) {
	bst.lock.Lock()
	defer bst.lock.Unlock()
	n := &DualTreeNode{key, value, nil, nil}
	if bst.root == nil {
		bst.root = n
	} else {
		insertNode(bst.root, n)
	}
}

// internal function to find the correct place for a node in a tree
func insertNode(node, newNode *DualTreeNode) {
	if newNode.key < node.key {
		if node.left == nil {
			node.left = newNode
		} else {
			insertNode(node.left, newNode)
		}
	} else {
		if node.right == nil {
			node.right = newNode
		} else {
			insertNode(node.right, newNode)
		}
	}
}

// InOrderTraverse visits all nodes with in-order traversing
func (bst *DualTree) InOrderTraverse(f func(string)) {
	bst.lock.RLock()
	defer bst.lock.RUnlock()
	inOrderTraverse(bst.root, f)
}

// internal recursive function to traverse in order
func inOrderTraverse(n *DualTreeNode, f func(string)) {
	if n != nil {
		inOrderTraverse(n.left, f)
		f(n.value)
		inOrderTraverse(n.right, f)
	}
}

// PreOrderTraverse visits all nodes with pre-order traversing
func (bst *DualTree) PreOrderTraverse(f func(string)) {
	bst.lock.Lock()
	defer bst.lock.Unlock()
	preOrderTraverse(bst.root, f)
}

// internal recursive function to traverse pre order
func preOrderTraverse(n *DualTreeNode, f func(string)) {
	if n != nil {
		f(n.value)
		preOrderTraverse(n.left, f)
		preOrderTraverse(n.right, f)
	}
}

// PostOrderTraverse visits all nodes with post-order traversing
func (bst *DualTree) PostOrderTraverse(f func(string)) {
	bst.lock.Lock()
	defer bst.lock.Unlock()
	postOrderTraverse(bst.root, f)
}

// internal recursive function to traverse post order
func postOrderTraverse(n *DualTreeNode, f func(string)) {
	if n != nil {
		postOrderTraverse(n.left, f)
		postOrderTraverse(n.right, f)
		f(n.value)
	}
}

// Min returns the string with min value stored in the tree
func (bst *DualTree) Min() *string {
	bst.lock.RLock()
	defer bst.lock.RUnlock()
	n := bst.root
	if n == nil {
		return nil
	}
	for {
		if n.left == nil {
			return &n.value
		}
		n = n.left
	}
}

// Max returns the string with max value stored in the tree
func (bst *DualTree) Max() *string {
	bst.lock.RLock()
	defer bst.lock.RUnlock()
	n := bst.root
	if n == nil {
		return nil
	}
	for {
		if n.right == nil {
			return &n.value
		}
		n = n.right
	}
}

// Search returns true if the string t exists in the tree
func (bst *DualTree) Search(key int) bool {
	bst.lock.RLock()
	defer bst.lock.RUnlock()
	return search(bst.root, key)
}

// internal recursive function to search an item in the tree
func search(n *DualTreeNode, key int) bool {
	if n == nil {
		return false
	}
	if key < n.key {
		return search(n.left, key)
	}
	if key > n.key {
		return search(n.right, key)
	}
	return true
}

// Remove removes the string with key `key` from the tree
func (bst *DualTree) Remove(key int) {
	bst.lock.Lock()
	defer bst.lock.Unlock()
	remove(bst.root, key)
}

// internal recursive function to remove an item
func remove(node *DualTreeNode, key int) *DualTreeNode {
	if node == nil {
		return nil
	}
	if key < node.key {
		node.left = remove(node.left, key)
		return node
	}
	if key > node.key {
		node.right = remove(node.right, key)
		return node
	}
	// key == node.key
	if node.left == nil && node.right == nil {
		node = nil
		return nil
	}
	if node.left == nil {
		node = node.right
		return node
	}
	if node.right == nil {
		node = node.left
		return node
	}
	leftmostrightside := node.right
	for {
		//find smallest value on the right side
		if leftmostrightside != nil && leftmostrightside.left != nil {
			leftmostrightside = leftmostrightside.left
		} else {
			break
		}
	}
	node.key, node.value = leftmostrightside.key, leftmostrightside.value
	node.right = remove(node.right, node.key)
	return node
}

// String prints a visual representation of the tree
func (bst *DualTree) String() {
	bst.lock.Lock()
	defer bst.lock.Unlock()
	fmt.Println("------------------------------------------------")
	stringify(bst.root, 0)
	fmt.Println("------------------------------------------------")
}

// internal recursive function to print a tree
func stringify(n *DualTreeNode, level int) {
	if n != nil {
		format := ""
		for i := 0; i < level; i++ {
			format += "       "
		}
		format += "---[ "
		level++
		stringify(n.left, level)
		fmt.Printf(format+"%d\n", n.key)
		stringify(n.right, level)
	}
}