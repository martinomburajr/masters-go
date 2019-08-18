package evolution

import (
	"fmt"
	"github.com/martinomburajr/masters-go/utils"
	"math"
	"math/rand"
	"strings"
	"sync"
)

// DualTree the binary search tree of Items
type DualTree struct {
	root *DualTreeNode
	lock sync.RWMutex
}

// RandomLeaf locates a random leaf within a tree and returns the ref to the node.
func (bst *DualTree) RandomLeaf() (*DualTreeNode, error) {
	if bst.root == nil {
		return nil, fmt.Errorf("root cannot be nil")
	}
	node := bst.root
	if node.left == nil && node.right == nil {
		return node, nil
	}

	nodes, err := bst.Leafs()
	if err != nil {
		return nil, err
	}

	randIndex := rand.Intn(len(nodes))
	return nodes[randIndex], nil
}

// Leafs returns all the leaves in a given tree
func (d *DualTree) Leafs() ([]*DualTreeNode, error) {
	nodes := make([]*DualTreeNode, 0)
	if d.root == nil {
		return nil, fmt.Errorf("tree root cannot be nil")
	}

	if d.root.right == nil && d.root.left == nil {
		nodes = append(nodes, d.root)
		return *(&nodes), nil
	}

	leaf(d.root, &nodes)
	return *(&nodes), nil
}

// leaf recursively adds terminal nodes to the nodes slice
func leaf(node *DualTreeNode, nodes *[]*DualTreeNode) {
	if node.left != nil {
		leaf(node.left, nodes)
		if node.right != nil {
			leaf(node.right, nodes)
		}
	}
	if node.left == nil {
		*nodes = append(*nodes, node)
		return
	}
}

// Branches returns a list of non-terminal nodes
func (d *DualTree) Branches() ([]*DualTreeNode, error) {
	nodes := make([]*DualTreeNode, 0)
	if d.root == nil {
		return nil, fmt.Errorf("tree root cannot be nil")
	}

	if d.root.right == nil && d.root.left == nil {
		nodes = append(nodes, d.root)
		return nil, fmt.Errorf("tree has size (1) root is not a nonterminal")
	}

	branch(d.root, &nodes)
	return *(&nodes), nil
}

// branch recursively adds non-terminal nodes to the nodes slice
func branch(node *DualTreeNode, nodes *[]*DualTreeNode) {
	if node.left != nil {
		branch(node.left, nodes)
		*nodes = append(*nodes, node)
		branch(node.right, nodes)
	}

	return
}


func (bst *DualTree) AddSubTree(subTree DualTree) error {
	return nil
}

func (bst *DualTree) DeleteSubTree() error {
	return nil
}

func (bst *DualTree) SoftDeleteSubTree() error {
	return nil
}

func (bst *DualTree) SwapSubTrees() error {
	return nil
}

func (bst *DualTree) Mutate() error {
	return nil
}

func (bst *DualTree) MutateDelete() error {
	return nil
}

func (bst *DualTree) GetRandomSubTree(depth int) error {
	return nil
}

func (bst *DualTree) Count() (int) {
	count := 0
	bst.InOrderTraverse(func(node *DualTreeNode) {
		count++
	})
	return count
}



/**
FromNodeTypes Creates a Tree from a list of NodeTypes
*/
func (bst *DualTree) FromSymbolicExpressionSet(terminalSet []SymbolicExpression) error {
	//EdgeCases
	if terminalSet == nil {
		return fmt.Errorf("terminalSet cannot be nil")
	}
	if len(terminalSet) < 1 {
		return fmt.Errorf("terminalSet cannot be empty i.e size 0")
	}
	if terminalSet[0].kind >= 1 {
		return fmt.Errorf("terminalSet cannot start with type nonterminal i.e SymbolicExpression.kind > 1")
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
		rem := bst.root.ArityRemainder()
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

	rem := bst.root.ArityRemainder()
	if rem != 0 {
		return fmt.Errorf("invalid tree - arity remainder is %d for root", rem)
	}

	return nil

}


func (bst *DualTree) Random(terminalSet []SymbolicExpression, maxDepth int) error {
	return nil
}

// Insert inserts the string t in the tree. Node must already contain the key and value
func (bst *DualTree) Insert(node *DualTreeNode) {
	bst.lock.Lock()
	defer bst.lock.Unlock()

	if bst.root == nil {
		if utils.TypeOf(node) == "SymbolicExpression" {
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
func (bst *DualTree) InOrderTraverse(f func(node *DualTreeNode)) {
	bst.lock.RLock()
	defer bst.lock.RUnlock()
	inOrderTraverse(bst.root, f)
}

// internal recursive function to traverse in order
func inOrderTraverse(n *DualTreeNode, f func(node *DualTreeNode)) {
	if n != nil {
		inOrderTraverse(n.left, f)
		f(n)
		inOrderTraverse(n.right, f)
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

// ToMathematicalString returns a mathematical representation of the tree after reading it using Inorder DFS
func (d *DualTree) ToMathematicalString() (string, error) {
	if d.root == nil {
		return "", fmt.Errorf("tree root is nil cannot compute mathematical expression")
	}

	var err error = nil

	sb := strings.Builder{}
	d.InOrderTraverse(func(node *DualTreeNode) {
		if node.arity == 1 && node.left == nil {
			err = fmt.Errorf("invalid tree structure, "+
				"unable to convert to mathematical expression: see node: %d", node.key)
			return
		}

		if node.arity > 1 && (node.left == nil || node.right == nil) {
			err = fmt.Errorf("invalid tree structure to convert to mathematical expression: see node: %d", node.key)
			return
		}
		sb.WriteString(node.value + " ")
	})

	if err != nil {
		return "", err
	}
	return strings.Trim(sb.String(), " "), err
}

func (d *DualTree) Validate() error {
	if d.root == nil {
		return fmt.Errorf("error: tree root is nil")
	}

	var err error = nil

	d.InOrderTraverse(func(node *DualTreeNode) {
		if node.arity == 1 && node.left == nil {
			err = fmt.Errorf("invalid tree structure, "+
				"unable to convert to mathematical expression: see node: %d", node.key)
			return
		}

		if node.arity > 1 && (node.left == nil || node.right == nil) {
			err = fmt.Errorf("invalid tree structure to convert to mathematical expression: see node: %d", node.key)
			return
		}
	})

	return err
}

// GenerateRandomTree generates a given tree of a depth between 0 (i.e) root and (inclusive of) the depth specified.
// Assuming a binary structured tree. The number of terminals (T) is equal to 2^D where D is the depth.
// The number of NonTerminals (NT) is equal to 2^D - 1
func GenerateRandomTree(maxDepth int, terminals []SymbolicExpression,
	nonTerminals []SymbolicExpression) (*DualTree, error) {

	if maxDepth < 0 {
		return nil, fmt.Errorf("maxDepth cannot be less than 0")
	}
	if terminals == nil {
		return nil, fmt.Errorf("terminal expression set cannot be nil")
	}
	if nonTerminals == nil {
		return nil, fmt.Errorf("nonterminal expression set cannot be nil")
	}
	if len(terminals) < 1 {
		return nil, fmt.Errorf("terminal expression set cannot be empty")
	}
	if len(nonTerminals) < 1 {
		return nil, fmt.Errorf("nonterminal expression set cannot be empty")
	}

	tree := DualTree{}

	depth := rand.Intn(maxDepth-0) + 0

	terminalCount := math.Pow(2, float64(depth))
	nonTerminalCount := math.Pow(2, float64(depth)) - 1

	randTerminals := make([]SymbolicExpression, int(terminalCount))
	for i := 0; i < int(terminalCount); i++ {
		randTerminalIndex := rand.Intn(len(terminals))
		randTerminals[i] = terminals[randTerminalIndex]
	}

	randNonTerminals := make([]SymbolicExpression, int(terminalCount))
	for i := 0; i < int(nonTerminalCount); i++ {
		index := rand.Intn(len(nonTerminals))
		randNonTerminals[i] = nonTerminals[index]
	}

	combinedArr := append(randTerminals, randTerminals...)
	err := tree.FromSymbolicExpressionSet(combinedArr)
	if err != nil {
		return nil, fmt.Errorf("error creating random tree | %s", err.Error())
	}
	return &tree, nil
}

// DualTreeNode a single node that composes the tree
type DualTreeNode struct {
	key   int
	value string
	left  *DualTreeNode //left
	right *DualTreeNode //right
	arity int
}

// SymbolicExpressionSet represents a mathematical expression broken into symbolic expressions.
// For Example x+1 will be broken into a SymbolicExpressionSet of size 3,
// containing both terminal and non terminal information
type SymbolicExpressionSet []SymbolicExpression

type SymbolicExpression struct {
	arity int
	value string
	kind  int //0 terminal >0 non-terminal
}

func (n *SymbolicExpression) CreateNonTerminal(arity int, value string) {
	n.arity = arity
	n.value = value
	n.kind = 1
}

func (n *SymbolicExpression) CreateTerminal(arity int, value string) {
	n.arity = arity
	n.value = value
	n.kind = 0
}

func (n *SymbolicExpression) ToDualTreeNode(key int) *DualTreeNode {
	return &DualTreeNode{
		value: n.value,
		arity: n.arity,
		left:  nil,
		right: nil,
		key:   key,
	}
}

// GenerateN generates a random SymbolicExpressionSet representing a valid mathematical expression.
// If size is less than 0, it reverts it to 0
func GenerateRandomSymbolicExpressionSet(size int) []SymbolicExpression {
	if size < 0 {
		size = 0
	}
	symbolicExpressions := make([]SymbolicExpression, 1)
	symbolicExpressions[0] = X1

	if size < 3 {
		return symbolicExpressions
	}
	for i := 1; i < size; i += 2 {
		if i%2 == 1 && i < (size-1) {
			symbolicExpressions = append(symbolicExpressions, Add)
			symbolicExpressions = append(symbolicExpressions, X1)
		}
	}

	return symbolicExpressions
}
