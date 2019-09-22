package evolution

// IDualTree represents a complete behavior for a treeNode
type IDualTree interface {
	IDualTreeInsertable
	Get(index int) (*DualTreeNode, error)
	GetFirst(node DualTreeNode) (*DualTreeNode, error)
	Pop(index int) (*DualTreeNode, error)
	Delete(index int) error
	DeleteFirst(node DualTreeNode) error
	Swap(index int, node DualTreeNode, newNode DualTreeNode) error
	Traverse(traversalMethod string) []*DualTreeNode
}

type IDualTreeInsertable interface {
	Insert(node DualTreeNode, index int) error
}

type IDualTreeGettable interface {
	Get(index int) (*DualTreeNode, error)
}

type IDualTreeGetFirstable interface {
	GetFirst(node DualTreeNode) (*DualTreeNode, error)
}

type IDualTreePoppable interface {
	Pop(index int) (*DualTreeNode, error)
}

// IDualTreeTraversable manages traversal behaviours
type IDualTreeTraversable interface {
	Traverse(traversalMethod string) []*DualTreeNode
}

// IDualTreeDeletable manages deletion behaviour
type IDualTreeDeletable interface {
	Delete(index int) error
}
