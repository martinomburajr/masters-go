package tree

import "github.com/martinomburajr/masters-go/program/tree/dualtree"

// IDualTree represents a complete behavior for a tree
type IDualTree interface {
	IDualTreeInsertable
	Get(index int) (*dualtree.DualTreeNode, error)
	GetFirst(node dualtree.DualTreeNode) (*dualtree.DualTreeNode, error)
	Pop(index int) (*dualtree.DualTreeNode, error)
	Delete(index int) error
	DeleteFirst(node dualtree.DualTreeNode) error
	Swap(index int, node dualtree.DualTreeNode, newNode dualtree.DualTreeNode) error
	Traverse(traversalMethod string) []*dualtree.DualTreeNode
}

type IDualTreeInsertable interface {
	Insert(node dualtree.DualTreeNode, index int) error
}

type IDualTreeGettable interface {
	Get(index int) (*dualtree.DualTreeNode, error)
}

type IDualTreeGetFirstable interface {
	GetFirst(node dualtree.DualTreeNode) (*dualtree.DualTreeNode, error)
}

type IDualTreePoppable interface {
	Pop(index int) (*dualtree.DualTreeNode, error)
}

// IDualTreeTraversable manages traversal behaviours
type IDualTreeTraversable interface {
	Traverse(traversalMethod string) []*dualtree.DualTreeNode
}

// IDualTreeDeletable manages deletion behaviour
type IDualTreeDeletable interface {
	Delete(index int) error
}
