package evolution

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// TODO enable concurrent safe access for DualTree Methods

// DualTree the binary search treeNode of Items
type DualTree struct {
	root *DualTreeNode
	lock sync.RWMutex
}

// RandomLeaf locates a random leaf within a treeNode and returns the ref to the node.
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

// RandomBranch locates a random branch within a treeNode and returns the ref to the node.
func (bst *DualTree) RandomBranch() (*DualTreeNode, error) {
	if bst.root == nil {
		return nil, fmt.Errorf("root cannot be nil")
	}
	node := bst.root
	if node.left == nil && node.right == nil {
		return nil, fmt.Errorf("invalid treeNode, cannot only contain non-terminal")
	}

	nodes, err := bst.Branches()
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(len(nodes))
	return nodes[randIndex], nil
}

// Leafs returns all the leaves in a given treeNode
func (d *DualTree) Leafs() ([]*DualTreeNode, error) {
	nodes := make([]*DualTreeNode, 0)
	if d.root == nil {
		return nil, fmt.Errorf("treeNode root cannot be nil")
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
		return nil, fmt.Errorf("treeNode root cannot be nil")
	}

	if d.root.right == nil && d.root.left == nil {
		nodes = append(nodes, d.root)
		return nil, fmt.Errorf("treeNode has size (1) root is not a nonterminal")
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

// StrategyAddSubTree adds a given subtree to a treeNode.
func (bst *DualTree) AddSubTree(subTree *DualTree) error {
	if subTree == nil {
		return fmt.Errorf("cannot add a nil subTree")
	}
	if subTree.root == nil {
		return fmt.Errorf("cannot add a subTree with a nil root")
	}
	if subTree.root.left == nil && subTree.root.right == nil {
		return fmt.Errorf("subTree cannot be composed of a single terminal - no operation to add it to the treeNode.")
	}

	if bst.root == nil {
		return fmt.Errorf("treeNode you are adding to has nil root")
	}
	if bst.root.left == nil && bst.root.right == nil {
		return fmt.Errorf("treeNode you are adding to is a lone terminal")
	}

	node, err := bst.RandomBranch()
	if err != nil {
		return err
	}

	// Can check for arity
	intn := rand.Intn(2)
	if intn == 0 {
		node.right = subTree.root
	} else {
		node.left = subTree.root
	}

	log.Print("Tree after operation")
	bst.Print()
	return nil
}

// InsertSubTree will insert a subTree at a given index
func (bst *DualTree) InsertSubTree(index int, subTree *DualTree) error {
	return nil
}

// StrategyDeleteSubTree locates a random non-terminal and sets its value to 0,
// deleting its associated child nodes by setting them to nil
func (bst *DualTree) DeleteSubTree() error {
	if bst.root == nil {
		return fmt.Errorf("treeNode you are deleting to has nil root")
	}
	if bst.root.left == nil && bst.root.right == nil {
		return fmt.Errorf("treeNode you are deleting to is a lone terminal")
	}

	node, err := bst.RandomBranch()
	if err != nil {
		return err
	}

	remove2(node)
	node.arity = 1
	node.value = "0"
	return nil
}

func remove2(node *DualTreeNode) {
	if node.left != nil {
		remove2(node.left)
		node.left = nil
	}
	if node.right != nil {
		remove2(node.right)
		node.right = nil
	}
	if node.left == nil && node.right == nil {
		node = nil
		return
	}
}

func (bst *DualTree) SoftDeleteSubTree() error {
	return nil
}

// SwapSubTrees swaps a set of subtrees in a given treeNode. It is a bit expensive as the parent needs to be obtained
// TODO Create Efficient Way of Locating Parent of NonTerminal Node
func (bst *DualTree) SwapSubTrees() error {
	if bst.root == nil {
		return fmt.Errorf("treeNode you are swapping to has nil root")
	}
	if bst.root.left == nil && bst.root.right == nil {
		return fmt.Errorf("treeNode you are swapping to is a lone terminal")
	}

	nodes, err := bst.Branches()
	if err != nil {
		return err
	}

	nonTerminalIndex0 := 0
	nonTerminalIndex1 := 0

	for nonTerminalIndex0 == nonTerminalIndex1 {
		rand.Seed(time.Now().UnixNano())
		nonTerminalIndex0 = rand.Intn(len(nodes))
		rand.Seed(time.Now().UnixNano())
		nonTerminalIndex1 = rand.Intn(len(nodes))
	}
	// once they are different

	//nodes[nonTerminalIndex0]
	return nil
}

// MutateTerminal will mutate a terminal to another valid terminal. If the terminalSet only contains a single item,
// that is already in the treeNode and that treeNode element is of size 1 (root only).
// If both these elements are identical no change will occur
func (bst *DualTree) MutateTerminal(terminalSet []SymbolicExpression) error {
	if bst.root == nil {
		return fmt.Errorf("treeNode you are swapping to has nil root")
	}
	if terminalSet == nil {
		return fmt.Errorf("terminal set cannot be nil")
	}
	if len(terminalSet) < 1 {
		return fmt.Errorf("terminal set cannot be empty")
	}

	nodes, err := bst.Leafs()
	if err != nil {
		return err
	}

	nodeValue := ""
	itemFromSet := ""

	for nodeValue == itemFromSet {
		rand.Seed(time.Now().UnixNano())
		nonTerminalIndex0 := rand.Intn(len(nodes))

		rand.Seed(time.Now().UnixNano())
		itemFromTSet := terminalSet[rand.Intn(len(terminalSet))]
		nodeValue = nodes[nonTerminalIndex0].value
		itemFromSet = itemFromTSet.value

		if nodeValue == itemFromSet {
			if len(terminalSet) < 2 {
				return nil
			}
			continue
		} else {
			nodes[nonTerminalIndex0].value = itemFromSet
		}
	}

	return nil
}



// MutateNonTerminal will mutate a terminal to another valid nonTerminal. Ensure set is nonTerminal set only,
// otherwise arities will break
// NOTE ensure nonTerminalSet contains no duplicates
func (bst *DualTree) MutateNonTerminal(nonTerminalSet []SymbolicExpression) error {
	if bst.root == nil {
		return fmt.Errorf("treeNode you are swapping to has nil root")
	}
	if bst.root.left == nil && bst.root.right == nil {
		return fmt.Errorf("treeNode you are swapping to is a lone terminal")
	}
	if nonTerminalSet == nil {
		return fmt.Errorf("nonTerminalSet set cannot be nil")
	}
	if len(nonTerminalSet) < 1 {
		return fmt.Errorf("nonTerminalSet set cannot be empty")
	}

	nodes, err := bst.Branches()
	if err != nil {
		return err
	}

	nodeValue := ""
	fromSetValue := ""

	counter := 0
	for nodeValue == fromSetValue && len(nonTerminalSet) >= 1 && counter < 20 { //pray for no duplicates.
		// Counter is a failsafe to prevent infinite looping
		rand.Seed(time.Now().UnixNano())
		nonTerminalIndex := rand.Intn(len(nodes))

		rand.Seed(time.Now().UnixNano())
		nonTerminalSetIndex := rand.Intn(len(nonTerminalSet))

		nodeValue = nodes[nonTerminalIndex].value
		fromSetValue = nonTerminalSet[nonTerminalSetIndex].value

		if nodeValue == fromSetValue {
			if len(nonTerminalSet) == 1 {
				hasDiverseNonTerminalSet, err := bst.hasDiverseNonTerminalSet()
				if err != nil {
					return err
				}
				if !hasDiverseNonTerminalSet {
					// If the terminal set only has an item (
					// which will always get chosen since its just one item) and the treeNode has no set of diverse terminals
					// i.e. all the nonterminals have the same value, just return as no useful work can be done here.
					return nil
				} else {
					// if the treeNode has a diverse set of non terminals,
					// then the set should at least be able to replace one of those differing Non Terminals
					continue
				}
			}
			if len(nonTerminalSet) > 1 {
				continue
			}
		} else {
			nodes[nonTerminalIndex].value = fromSetValue
			return nil
		}
		counter++
	}

	return nil
}

func (bst *DualTree) GetRandomSubTree() (*DualTree, error) {
	if bst.root == nil {
		return nil, fmt.Errorf("treeNode you are adding to has nil root")
	}
	if bst.root.left == nil && bst.root.right == nil {
		return nil, fmt.Errorf("treeNode you are adding to is a lone terminal")
	}
	//node, err := bst.RandomBranch()
	//if err != nil {
	//	return nil, err
	//}

	return nil, nil
}

// GetNode returns the first node it encounters with the given value.
// It uses Inorder DFS to iterate through the treeNode, if it cannot locate an object it returns an error,
// if the treeNode is of size 1 i.e only containing the root, both parent and node will point to the root,
// in all other cases where size > 1, node and parent will be different granted the value can be found.
func (bst *DualTree) GetNode(value string) (node *DualTreeNode, parent *DualTreeNode, err error) {
	// require
	if bst.root == nil {
		return nil, nil, fmt.Errorf("treeNode root cannot be nil")
	}
	if value == "" {
		return nil, nil, fmt.Errorf("value cannot be empty")
	}
	if bst.root.left == nil && bst.root.right == nil {
		if bst.root.value != value {
			return nil, nil, fmt.Errorf("treeNode node not found")
		}else {
			return bst.root, bst.root, nil
		}
	}

	err = fmt.Errorf("treeNode node not found")

	bst.InOrderTraverseAware(func(n *DualTreeNode, p *DualTreeNode) {
		if n.value == value {
			node = n
			parent = p
			err = nil
			return
		}
	})

	return node, parent, err
}

// SelectNodesUpToDepth returns a list of nodes that do not exceed the indicated depth.
// This will start from the root of the treeNode and follow and InorderDFS
func (bst *DualTree) SelectNodesUpToDepth(depth int) ([]*DualTreeNode, error) {
	if bst.root == nil {
		return nil, fmt.Errorf("cannot find depth in nil treeNode | root == nil")
	}

	nodes := make([]*DualTreeNode, 0)
	nodes = diveMaxDepth(bst.root, depth, nodes)

	return nodes, nil
}

func diveMaxDepth(node *DualTreeNode, maxDepth int, nodes []*DualTreeNode) []*DualTreeNode {
	if node == nil {
		return nodes
	}
	lDepth := dive(node.left)
	rDepth := dive(node.right)

	if lDepth < maxDepth {
		nodes = append(nodes, node)
	}
	if rDepth < maxDepth {
		nodes = append(nodes, node)
	}
	return nodes
}

// CalculateMinDepth calculates the minimum depth of a given treeNode
//func CalculateMinDepth(nodeCount int) int {
//	if nodeCount < 1 {
//		return 0
//	}
//}

// Clone will perform an O(N) deep clone of a treeNode and its items and return its copy.
func (bst DualTree) Clone() DualTree {
	return bst
}

func (bst *DualTree) Size() int {
	count := 0
	bst.InOrderTraverse(func(node *DualTreeNode) {
		count++
	})
	return count
}

// ContainsSubTree checks to see if a treeNode contains part of a subTree
func (bst *DualTree) ContainsSubTree(subTree *DualTree) (bool, error) {
	if subTree == nil {
		return false, fmt.Errorf("cannot add a nil subTree")
	}
	if subTree.root == nil {
		return false, fmt.Errorf("cannot add a subTree with a nil root")
	}
	if bst.root == nil {
		return false, fmt.Errorf("treeNode you are adding to has nil root")
	}

	subTreeSlice := make([]*DualTreeNode, 0)
	subTree.InOrderTraverse(func(node *DualTreeNode) {
		subTreeSlice = append(subTreeSlice, node)
	})

	tree := make([]*DualTreeNode, 0)
	bst.InOrderTraverse(func(node *DualTreeNode) {
		tree = append(tree, node)
	})

	if len(subTreeSlice) > len(tree) {
		return false, nil
	}

	for i := range tree {
		if tree[i].IsValEqual(subTreeSlice[0]) {
			count := 0
			for j := 0; j < len(subTreeSlice); j++ {
				if !tree[i+j].IsValEqual(subTreeSlice[j]) {
					break
				}
				count++
				if count == len(subTreeSlice) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// ContainsNode checks to see if a treeNode contains a given node
func (bst *DualTree) ContainsNode(treeNode *DualTreeNode) (bool, error) {
	if bst.root == nil {
		return false, fmt.Errorf("treeNode has nil root")
	}
	if treeNode == nil {
		return false, fmt.Errorf("cannot search for a nil treeNode")
	}

	found := false
	bst.InOrderTraverse(func(node *DualTreeNode) {
		if treeNode.IsEqual(node) {
			found = true
		}
		return
	})

	return found, nil
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
		bst.root = terminalSet[0].ToDualTreeNode(RandString(5))
		return nil
	}

	//MainCase setup  -  SetupRoot and First Child
	bst.root = terminalSet[1].ToDualTreeNode(RandString(5))
	bst.root.left = terminalSet[0].ToDualTreeNode(RandString(5))

	if terminalSet[0].kind < 1 && terminalSet[1].kind < 1 {
		return fmt.Errorf("cannot have adjacent terminals got %#v %#v", bst.root, bst.root.left)
	}

	//MainCases

	for i := 2; i < len(terminalSet); i++ {
		rem := bst.root.ArityRemainder()
		if rem == 0 {
			if terminalSet[i].kind >= 1 {
				dtn := terminalSet[i].ToDualTreeNode(RandString(5))
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
				dtn := terminalSet[i].ToDualTreeNode(RandString(5))
				bst.root.right = dtn
			}
		}
	}

	rem := bst.root.ArityRemainder()
	if rem != 0 {
		return fmt.Errorf("invalid treeNode - arity remainder is %d for root", rem)
	}

	return nil

}

/**
FromNodeTypes Creates a Tree from a list of NodeTypes
*/
func (bst *DualTree) FromSymbolicExpressionSet2(terminalSet []SymbolicExpression) error {
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
		bst.root = terminalSet[0].ToDualTreeNode(RandString(5))
		return nil
	}

	nodes, err := Splitter(terminalSet)
	if err != nil {
		return err
	}

	if len(nodes) == 1 {
		bst.root = nodes[0]
		return nil
	}

	bst.root = combinatorArr(nodes[0:len(nodes)/2+1], nodes[len(nodes)/2+1:], &DualTreeNode{}, &DualTreeNode{})

	return nil
}

func combinatorArr(left, right []*DualTreeNode, x, y *DualTreeNode) *DualTreeNode {
	if len(left) > 2 {
		x = combinatorArr(left[0:len(left)/2], left[len(left)/2:], x, y)
	}
	if len(right) > 2 {
		y = combinatorArr(right[0:len(right)/2], right[len(right)/2:], x, y)
	}
	if len(left) <= 2 {
		if len(left) == 2 {
			x = combinator(left[0], left[1])
		} else if len(left) == 1 {
			x = combinator(left[0], nil)
		}
	}
	if len(right) <= 2 {
		if len(right) == 2 {
			y = combinator(right[0], right[1])
		} else if len(right) == 1 {
			y = combinator(right[0], nil)
		}
	}
	return combinator(x, y)
}

/**
	Splitter takes a set of symbolic expressions and breaks them out to a set of other symbolic expressions with the
remainder being passed back as the symbolicExpression.
This will not check for empty expressionSets or expressionSets of len less than 3.
*/
func Splitter(expressionSet []SymbolicExpression) ([]*DualTreeNode, error) {
	if len(expressionSet)%2 == 0 {
		return nil, fmt.Errorf("expression set must have odd numbered values")
	}

	nodeSet := make([]*DualTreeNode, len(expressionSet))
	for e := range expressionSet {
		nodeSet[e] = expressionSet[e].ToDualTreeNode(RandString(5))
	}

	initialTrees := make([]*DualTreeNode, 0)
	for i := 0; i < len(nodeSet)-1; i += 2 {
		nodeSet[i+1].left = nodeSet[i]
		initialTrees = append(initialTrees, nodeSet[i+1])
	}
	initialTrees[len(initialTrees)-1].right = expressionSet[len(expressionSet)-1].ToDualTreeNode(RandString(5))
	return initialTrees, nil
}

func combinator(node0, node1 *DualTreeNode) *DualTreeNode {
	if node0 == nil {
		return node0
	}
	if node1 == nil {
		return node0
	}

	if node0.right == nil {
		if node1.right == nil {
			node0.right = node1.left
			node1.left = node0
			return node1
		} else {
			if node1.right.ArityRemainder() == 0 {
				node0.right = node1
				combinator(node0, nil)
			} else {

			}
		}
	}
	return node0
}

// Depth calculates the height of the treeNode. A treeNode with a nil root returns -1.
func (d *DualTree) Depth() (int, error) {
	if d.root == nil {
		return -1, fmt.Errorf("cannot find depth in nil treeNode | root == nil")
	}

	currentDepth := dive(d.root)
	return currentDepth-1, nil
}

func dive(node *DualTreeNode) int {
	if node == nil {
		return 0
	}
	lDepth := dive(node.left)
	rDepth := dive(node.right)

	if lDepth > rDepth {
		return lDepth + 1
	} else {
		return rDepth + 1
	}
}

// DepthAt returns a group of nodes AT the specified depth. No more no less.
func (d *DualTree) DepthAt(depth int) ([]*DualTreeNode, error) {
	if d.root == nil {
		return nil, fmt.Errorf("cannot find depth in nil treeNode | root == nil")
	}
	if depth < 0 {
		return nil, fmt.Errorf("cannot find depth when supplied depth is negative")
	}

	nodes := make([]*DualTreeNode, 0)
	startDepth := 0
	nodes = diveTo(d.root, &nodes, &startDepth, depth)
	return nodes, nil
}

func diveTo(node *DualTreeNode, nodes *[]*DualTreeNode, currDepth *int, maxDepth int) []*DualTreeNode {
	if node == nil {
		return *nodes
	}
	if *currDepth <= maxDepth {
		if node.left == nil && node.right == nil {
			if *currDepth == maxDepth {
				*nodes = append(*nodes, node)
				return *nodes
			}
			return *nodes
		}
		cDIncr := *currDepth + 1
		if cDIncr > maxDepth {
			*nodes = append(*nodes, node)
			return *nodes
		}
		*currDepth++
		diveTo(node.left, nodes, currDepth,  maxDepth)
		diveTo(node.right, nodes, currDepth, maxDepth)
		*currDepth--
	}
	return *nodes
}



// DepthTo calculates the depth of the treeNode until the given set of nodes.
// All nodes traversed will be returned until that point (Not including).
// A depth value greater than the size of the treeNode will return the entire treeNode's nodes.
func (d *DualTree) DepthTo(depth int) ([]*DualTreeNode, error) {
	if d.root == nil {
		return nil, fmt.Errorf("cannot find depth in nil treeNode | root == nil")
	}
	if depth < 0 {
		return nil, fmt.Errorf("cannot find depth when supplied depth is negative")
	}

	nodes := make([]*DualTreeNode, 0)
	startDepth := 0
	nodes = diveUntil(d.root, &nodes, &startDepth, depth)
	return nodes, nil
}


func diveUntil(node *DualTreeNode, nodes *[]*DualTreeNode, currDepth *int, maxDepth int) []*DualTreeNode {
	if node == nil {
		return *nodes
	}
	if *currDepth <= maxDepth {
		if node.left == nil && node.right == nil {
			*nodes = append(*nodes, node)
			return *nodes
		}
		*nodes = append(*nodes, node)
		cDIncr := *currDepth + 1
		if cDIncr > maxDepth {
			return *nodes
		}
		*currDepth++
		diveUntil(node.left, nodes, currDepth,  maxDepth)
		diveUntil(node.right, nodes, currDepth, maxDepth)
		*currDepth--
	}
	return *nodes
}

func (bst *DualTree) Random(terminalSet []SymbolicExpression, maxDepth int) error {
	return nil
}

// InOrderTraverse visits all nodes with in-order traversing
func (bst *DualTree) InOrderTraverse(f func(node *DualTreeNode)) {
	bst.lock.RLock()
	defer bst.lock.RUnlock()
	inOrderTraverse(bst.root, f)
}

// InOrderTraverse visits all nodes with in-order traversing but remembers its parent. (A good child :D)
func (bst *DualTree) InOrderTraverseAware(f func(node *DualTreeNode, parentNode *DualTreeNode)) {
	bst.lock.RLock()
	defer bst.lock.RUnlock()
	inOrderTraverseAware(bst.root, bst.root, f)
}

// InOrderTraverse visits all nodes with in-order traversing
func (bst *DualTree) ToSymbolicExpressionSet() []SymbolicExpression {
	symbSet := make([]SymbolicExpression, 0)
	bst.InOrderTraverse(func(node *DualTreeNode) {
		symbSet = append(symbSet, node.ToSymbolicExpression())
	})
	return symbSet
}

// internal recursive function to traverse in order
func inOrderTraverse(n *DualTreeNode, f func(node *DualTreeNode)) {
	if n != nil {
		inOrderTraverse(n.left, f)
		f(n)
		inOrderTraverse(n.right, f)
	}
}

// internal recursive function to traverse in order
func inOrderTraverseAware(n *DualTreeNode, parent *DualTreeNode, f func(node *DualTreeNode, parentNode *DualTreeNode)) {
	if n != nil {
		inOrderTraverseAware(n.left, n, f)
		f(n, parent)
		inOrderTraverseAware(n.right,n, f)
	}
}

// Print prints a visual representation of the treeNode
func (bst *DualTree) Print() {
	bst.lock.Lock()
	defer bst.lock.Unlock()
	fmt.Println("------------------------------------------------")
	stringify(bst.root, 0)
	fmt.Println("------------------------------------------------")
}

// internal recursive function to print a treeNode
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

// ToMathematicalString returns a mathematical representation of the treeNode after reading it using Inorder DFS
func (d *DualTree) ToMathematicalString() (string, error) {
	if d.root == nil {
		return "", fmt.Errorf("treeNode root is nil cannot compute mathematical expression")
	}

	var err error = nil

	sb := strings.Builder{}
	d.InOrderTraverse(func(node *DualTreeNode) {
		if node.arity == 1 && node.left == nil {
			err = fmt.Errorf("invalid treeNode structure, "+
				"unable to convert to mathematical expression: see node: %s", node.key)
			return
		}

		if node.arity > 1 && (node.left == nil || node.right == nil) {
			err = fmt.Errorf("invalid treeNode structure to convert to mathematical expression: see node: %s", node.key)
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
		return fmt.Errorf("error: treeNode root is nil")
	}

	var err error = nil

	d.InOrderTraverse(func(node *DualTreeNode) {
		if node.arity == 1 && node.left == nil {
			err = fmt.Errorf("invalid treeNode structure, "+
				"unable to convert to mathematical expression: see node: %s", node.key)
			return
		}

		if node.arity > 1 && (node.left == nil || node.right == nil) {
			err = fmt.Errorf("invalid treeNode structure to convert to mathematical expression: see node: %s", node.key)
			return
		}
	})

	return err
}

// GenerateRandomTree generates a given treeNode of a depth between 0 (i.e) root and (inclusive of) the depth specified.
// Assuming a binary structured treeNode. The number of terminals (T) is equal to 2^D where D is the depth.
// The number of NonTerminals (NT) is equal to 2^D - 1
func GenerateRandomTree(depth int, terminals []SymbolicExpression,
	nonTerminals []SymbolicExpression) (*DualTree, error) {

	if depth < 0 {
		return nil, fmt.Errorf("depth cannot be less than 0")
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
	if depth > 0 && len(nonTerminals) < 1 {
		return nil, fmt.Errorf("non terminal expression set cannot be empty if depth > 0")
	}
	if len(nonTerminals) < 1 {
		rand.Seed(time.Now().UnixNano())
		tree := &DualTree{}
		tree.root = terminals[rand.Intn(len(terminals))].ToDualTreeNode(RandString(5))
		return tree, nil
	}

	tree := DualTree{}
	terminalCount := 2
	if depth > 1 {
		terminalCount = int(math.Pow(2, float64(depth)))
	}
	nonTerminalCount := 1
	if depth > 1 {
		nonTerminalCount = int(math.Pow(2, float64(depth)) - 1)
	}

	randTerminals := make([]SymbolicExpression, terminalCount)
	for i := 0; i < terminalCount; i++ {
		rand.Seed(time.Now().UnixNano())
		randTerminalIndex := rand.Intn(len(terminals))
		randTerminals[i] = terminals[randTerminalIndex]
	}

	randNonTerminals := make([]SymbolicExpression, nonTerminalCount)
	for i := 0; i < nonTerminalCount; i++ {
		rand.Seed(time.Now().UnixNano())
		index := rand.Intn(len(nonTerminals))
		randNonTerminals[i] = nonTerminals[index]
	}

	if (len(randTerminals)+len(randNonTerminals))%2 != 1 {
		return nil, fmt.Errorf("bad pairing of terminals and non-terminals")
	}

	combinedArr := weaver(randTerminals, randNonTerminals)

	err := tree.FromSymbolicExpressionSet2(combinedArr)
	if err != nil {
		return nil, fmt.Errorf("error creating random treeNode | %s", err.Error())
	}
	return &tree, nil
}

func weaver(terminals, nonTerminals []SymbolicExpression) []SymbolicExpression {
	if len(terminals) < 1 {
		return []SymbolicExpression{}
	}
	if len(terminals) > 0 {
		if len(nonTerminals) < 1 {
			return []SymbolicExpression{terminals[0]}
		}
	}

	combined := make([]SymbolicExpression, len(terminals)+len(nonTerminals))

	count := 0
	for i := 0; i < len(combined); i += 2 {
		combined[i] = terminals[count]
		count++
	}
	count = 0
	for i := 0; i < len(combined)-1; i += 2 {
		combined[(i + 1)] = nonTerminals[count]
		count++
	}
	return combined
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

func (n *SymbolicExpression) ToDualTreeNode(key string) *DualTreeNode {
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


func (bst *DualTree) hasDiverseNonTerminalSet() (bool, error) {
	branches, err := bst.Branches()
	if err != nil {
		return false, err
	}

	holder := branches[0]
	for i := range branches {
		if !branches[i].IsValEqual(holder) {
			return true, nil
		}
	}
	return false, nil
}

// GetRandomSubTreeAtDepth will obtain a random subTree from a given treeNode.
// It assumes the depth you provide is the appropriate range as it WILL NOT check the depth and panic in case of a
// depth out of bounds. (This is done to prevent an extra redundant call to the depth method of treeNode if the user has
// already called it.
func (d *DualTree) GetRandomSubTreeAtDepth(depth int) (DualTree, error) {
	if d.root == nil {
		return DualTree{}, fmt.Errorf("cannot get depth - treeNode nil")
	}
	if depth < 0 {
		return DualTree{}, fmt.Errorf("cannot get depth - depth is less than 0")
	}

	rand.Seed(time.Now().UnixNano())
	randomDepth := rand.Intn(depth+1)

	nodes, err := d.DepthTo(randomDepth)
	if err != nil {
		return DualTree{}, err
	}

	rand.Seed(time.Now().UnixNano())
	randomNodeIndex := rand.Intn(len(nodes))
	randomNode := nodes[randomNodeIndex]

	return randomNode.ToDualTree(), nil
}

// GetRandomSubTreeAtDepth will obtain a random subTree from a given treeNode.
// It assumes the depth you provide is the appropriate range as it WILL NOT check the depth and panic in case of a
// depth out of bounds. (This is done to prevent an extra redundant call to the depth method of treeNode if the user has
// already called it.
func (d *DualTree) GetRandomSubTreeAtDepthAware(depth int) (DualTree, error) {
	if d.root == nil {
		return DualTree{}, fmt.Errorf("cannot get depth - treeNode nil")
	}
	if depth < 0 {
		return DualTree{}, fmt.Errorf("cannot get depth - depth is less than 0")
	}

	rand.Seed(time.Now().UnixNano())
	randomDepth := rand.Intn(depth+1)

	nodes, err := d.DepthTo(randomDepth)
	if err != nil {
		return DualTree{}, err
	}

	rand.Seed(time.Now().UnixNano())
	randomNodeIndex := rand.Intn(len(nodes))
	randomNode := nodes[randomNodeIndex]

	return randomNode.ToDualTree(), nil
}

//func Swapper(tree1 DualTree, tree2 DualTree, subTree1 DualTree, subTree2 DualTree) (p1 DualTree, p2 DualTree,
//	c1 DualTree, c2 DualTree) {
//
//}