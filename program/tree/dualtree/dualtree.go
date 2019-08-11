package dualtree

import (
	"fmt"
	"github.com/martinomburajr/masters-go/utils"
	"strings"
	"sync"
)

// DualTree the binary search tree of Items
type DualTree struct {
	root *DualTreeNode
	lock sync.RWMutex
}

// GetEquationSlice returns a slice array containing the various equation items for the given string
func (bst *DualTree) GetEquationSlice(equationString string) ([]string, error) {
	if equationString == "" {
		return nil, fmt.Errorf("empty equation string")
	}
	return strings.Split(equationString, ","), nil
}

func (bst *DualTree) FromString(equationStrings string) error {
	_, err := bst.GetEquationSlice(equationStrings)
	if err != nil {
		return err
	}
	return nil
}

/**
FromNodeTypes Creates a Tree from a list of NodeTypes
*/
func (bst *DualTree) FromNodeTypes(E []NodeType) error {
	if E == nil {
		return fmt.Errorf("nodeType passed in is nil")
	}
	if len(E) < 1 {
		return fmt.Errorf("nodeType passed in cannot be empty")
	}
	if len(E) == 1 {
		if E[0].kind <= 0 {
			bst.root = E[0].ToDualTreeNode(0)
			return nil
		} else {
			return fmt.Errorf("invalid, tree cannot contain non-terminal only: size 1")
		}
	}
	terminalArr := make([]NodeType, 0)

	for i := range E {
		if E[i].kind <= 0 { //Terminal
			if bst.root == nil {
				terminalArr = append(terminalArr, E[i])
			} else {
				dtn := E[i].ToDualTreeNode(i)
				bst.root.right = dtn
			}
			continue
		}

		if E[i].kind > 0 { //NonTerminal
			dtn := E[i].ToDualTreeNode(i)

			if bst.root == nil {
				bst.root = dtn
			} else {
				oldRoot := bst.root
				dtn.left = oldRoot
				bst.root = dtn
			}

			if len(terminalArr) > 0 {
				bst.root.left = terminalArr[0].ToDualTreeNode(i)
				if len(terminalArr) == 1 {
					terminalArr = make([]NodeType, 0) //pop
				} else {
					_, terminalArr = terminalArr[0], terminalArr[1:] //pop
				}
			}
		}
	}

	return nil
}

/**
FromNodeTypes Creates a Tree from a list of NodeTypes
*/
func (bst *DualTree) FromTerminalSet(terminalSet []NodeType) error {
	//EdgeCases
	if terminalSet == nil {
		return fmt.Errorf("terminalSet cannot be nil")
	}
	if len(terminalSet) < 1 {
		return fmt.Errorf("terminalSet cannot be empty i.e size 0")
	}
	if terminalSet[0].kind >= 1 {
		return fmt.Errorf("terminalSet cannot start with type nonterminal i.e NodeType.kind > 1")
	}
	if len(terminalSet) == 1 && terminalSet[0].kind < 1 {
		bst.root = terminalSet[0].ToDualTreeNode(0)
		return nil
	}

	//MainCase setup  -  SetupRoot and First Child
	bst.root = terminalSet[1].ToDualTreeNode(0)
	bst.root.left = terminalSet[0].ToDualTreeNode(1)

	if terminalSet[0].kind < 1 && terminalSet[1].kind < 1 {
		return fmt.Errorf("cannot have adjacent terminals got %#v %#v", bst.root, bst.root.left)
	}

	//MainCases

	for i := 2; i < len(terminalSet); i++ {
		rem := arityRemainder(bst.root)
		if rem == 0 {
			if terminalSet[i].kind >= 1 {
				dtn := terminalSet[i].ToDualTreeNode(i)
				oldRoot := bst.root
				dtn.left = oldRoot
				bst.root = dtn
			} else {
				return fmt.Errorf("expected non-terminal at index: %d got terminal %s", i, terminalSet[i].value)
			}
		} else {
			if terminalSet[i].kind >= 1 {
				return fmt.Errorf("expected terminal at index: %d | got non-terminal %s", i, terminalSet[i].value)
			} else {
				dtn := terminalSet[i].ToDualTreeNode(i)
				bst.root.right = dtn
			}
		}
	}

	rem := arityRemainder(bst.root)
	if rem != 0 {
		return fmt.Errorf("invalid tree - arity remainder is %d for root", rem)
	}

	return nil

}

// arityRemainder calculates the remaining arity for a given root node.
// This is used to balance the NonTerminals and the Terminals depending on their requirements.
func arityRemainder(root *DualTreeNode) int {
	available := root.arity
	if root.arity == 2 {
		if root.right != nil {
			available--
		}
		if root.left != nil {
			available--
		}
		return available
	} else if root.arity == 1 {
		if root.left != nil {
			available--
		}
		return available
	}
	return 0
}

// Insert inserts the string t in the tree. Node must already contain the key and value
func (bst *DualTree) Insert(node *DualTreeNode) {
	bst.lock.Lock()
	defer bst.lock.Unlock()

	if bst.root == nil {
		if utils.TypeOf(node) == "NodeType" {
			bst.root = node
		}
	} else {
		insertNode(bst.root, node)
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
		stringify(n.right, level)
		fmt.Printf(format+"%s\n", n.value)
		stringify(n.left, level)
	}
}

// DualTreeNode a single node that composes the tree
type DualTreeNode struct {
	key   int
	value string
	left  *DualTreeNode //left
	right *DualTreeNode //right
	arity int
}

type NodeType struct {
	arity int
	value string
	kind  int //0 terminal >0 non-terminal
}

func (n *NodeType) CreateNonTerminal(arity int, value string) {
	n.arity = arity
	n.value = value
	n.kind = 1
}

func (n *NodeType) CreateTerminal(arity int, value string) {
	n.arity = arity
	n.value = value
	n.kind = 0
}

func (n *NodeType) ToDualTreeNode(key int) *DualTreeNode {
	return &DualTreeNode{
		value: n.value,
		arity: n.arity,
		left:  nil,
		right: nil,
		key:   key,
	}
}
