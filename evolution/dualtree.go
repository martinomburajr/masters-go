package evolution

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
)

// TODO enable concurrent safe access for DualTree Methods

// DualTree the binary search treeNode of Items
type DualTree struct {
	root *DualTreeNode
	lock sync.RWMutex
}

// RandomTerminal locates a random leaf within a treeNode and returns the ref to the node.
func (bst *DualTree) RandomTerminal() (*DualTreeNode, error) {
	if bst.root == nil {
		return nil, fmt.Errorf("root cannot be nil")
	}
	node := bst.root
	if node.left == nil && node.right == nil {
		return node, nil
	}

	nodes, err := bst.Terminals()
	if err != nil {
		return nil, err
	}

	randIndex := rand.Intn(len(nodes))
	return nodes[randIndex], nil
}

// RandomTerminalAware returns a random leaf along with their parent. If the Tree is of depth 0 -> Parent is always nil.
// It is the clients responsibility to check for this.
// The node value can never be nil on a valid tree if the error is nil
func (bst *DualTree) RandomTerminalAware() (node *DualTreeNode, parent *DualTreeNode, err error) {
	if bst.root == nil {
		return nil, nil, fmt.Errorf("root cannot be nil")
	}

	trees, err := bst.GetTerminalsAware()
	if err != nil {
		return nil, nil, err
	}

	randIndex := rand.Intn(len(trees))
	node = trees[randIndex].node
	parent = trees[randIndex].parent

	return node, parent, nil
}

// GetTerminalsAware returns a slice of AwareTrees.
// The AwareTree will always have at least one value if no error is returned.
func (bst *DualTree) GetTerminalsAware() ([]AwareTree, error) {
	if bst.root == nil {
		return nil, fmt.Errorf("GetNonTerminalsAware | root cannot be nil")
	}
	awareTrees := make([]AwareTree, 0)
	if bst.root.left == nil && bst.root.right == nil {
		awareTrees = append(awareTrees, AwareTree{node: bst.root, parent: nil})
		return awareTrees, nil
	}

	bst.InOrderTraverseAware(func(n *DualTreeNode, parentNode *DualTreeNode) {
		if n.IsLeaf() {
			awareTrees = append(awareTrees, AwareTree{node: n, parent: parentNode})
		}
	})
	return awareTrees, nil
}

// RandomNonTerminal locates a random non-terminal within a tree and returns the ref to the node.
func (bst *DualTree) RandomNonTerminal() (*DualTreeNode, error) {
	if bst.root == nil {
		return nil, fmt.Errorf("root cannot be nil")
	}
	node := bst.root
	if node.left == nil && node.right == nil {
		return nil, fmt.Errorf("invalid treeNode, cannot only contain non-terminal")
	}

	nodes, err := bst.NonTerminals()
	if err != nil {
		return nil, err
	}

	randIndex := rand.Intn(len(nodes))
	return nodes[randIndex], nil
}

// RandomNonTerminalAware returns a random leaf along with their parent. If the Tree is of depth 0 -> Parent is always nil.
// It is the RandomNonTerminalAware responsibility to check for this.
func (bst *DualTree) RandomNonTerminalAware() (node *DualTreeNode, parent *DualTreeNode, err error) {
	if bst.root == nil {
		return nil, nil, fmt.Errorf("RandomNonTerminalAware | root cannot be nil")
	}

	trees, err := bst.GetNonTerminalsAware()
	if err != nil {
		return nil, nil, err
	}

	if len(trees) > 0 {
		randIndex := rand.Intn(len(trees))
		node = trees[randIndex].node
		parent = trees[randIndex].parent
	}

	return node, parent, nil
}

// GetNonTerminalsAware is a utility function that returns a set of non-terminal nodes.
// In the event that the tree contains a lone terminal as its root. It will return an empty awareTree array.
// It is for the caller to understand that an empty awareTree array means there are no non-terminals within the tree
// and act appropriately
func (bst *DualTree) GetNonTerminalsAware() ([]AwareTree, error) {
	if bst.root == nil {
		return nil, fmt.Errorf("GetNonTerminalsAware | root cannot be nil")
	}
	awareTrees := make([]AwareTree, 0)
	if bst.root.left == nil && bst.root.right == nil {
		return awareTrees, nil
	}

	bst.InOrderTraverseAware(func(n *DualTreeNode, parentNode *DualTreeNode) {
		if !n.IsLeaf() {
			if n.IsEqual(bst.root) {
				awareTrees = append(awareTrees, AwareTree{node: n, parent: nil})
			} else {
				awareTrees = append(awareTrees, AwareTree{node: n, parent: parentNode})
			}
		}
	})
	return awareTrees, nil
}

type AwareTree struct {
	node   *DualTreeNode
	parent *DualTreeNode
}

// Terminals returns all the leaves in a given tree
func (d *DualTree) Terminals() ([]*DualTreeNode, error) {
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

// NonTerminals returns a list of non-terminal nodes
func (d *DualTree) NonTerminals() ([]*DualTreeNode, error) {
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

// AddEmptyToTree is a conservative means of performing add and delete operations on trees that end up as single
// nodes. This function will add a 0 or subtract a 0 whenever an add operation encounters a single node
func (bst *DualTree) AddEmptyToTreeRoot(subTree *DualTree) error {
	if bst.root == nil {
		return fmt.Errorf("treeNode you are deleting to has nil root")
	}
	if bst.root.right == nil && bst.root.left == nil {
		addNode := SymbolicExpression{value: "+", arity: 2}
		treeNode := bst.root.Clone()
		nodePlus := addNode.ToDualTreeNode(RandString(5))
		bst.root = nodePlus
		bst.root.left = &treeNode
		bst.root.right = subTree.root
		return nil
	}
	return fmt.Errorf("AddEmptyToTreeRoot | is not Tree root")
}

// DeleteSubTree will delete a random branch or node of a subTree.
// Depending on the deletion Strategy 0 - DeleteSafe 1 - DeleteMalicious the Tree will be deleted in either
// conservative or aggressive ways
// If the Tree is a lone terminal - DeleteSafe will append a subtraction as well as copy of the value itself e.g.
// if the lone node is 3, the new Tree will become 3 - 3. Similarly if the node is x the new Tree will become x - x.
// This only occurs for lone nodes. If the Tree is not a lone terminal, but is of size 1
// delete safe will turn the value of a terminal to zero. E.g. If a Tree represents  3 * 4,
// one of the terminals will be set to 0. The resulting may be 0 * 3 or 0 * 4.
// Trees that have a greater size than 1 can have actual branches deleted.
// A Tree like 3 + 4 * x may become 4 * x but never be converted to 0. This may have little impact on + or - operators.
// DeleteMalicious will follow the same semantics as deleteSafe except
// If the Tree is already a terminal it will convert the lone-terminal to a zero.
// DeleteMalicious can also reduce an entire Tree to zero if it chooses the root as a target to delete.
// This function will never return a nil Tree or a Tree with a nil root.
func (bst *DualTree) DeleteSubTree(deletionStrategy int) error {
	if bst.root == nil {
		return fmt.Errorf("treeNode you are deleting to has nil root")
	}
	if deletionStrategy < 0 && deletionStrategy > 1 {
		deletionStrategy = 0
	}

	if bst.root.left == nil && bst.root.right == nil {
		if deletionStrategy == 1 { // MaliciousDelete
			const0 := SymbolicExpression{arity: 0, value: "0"}
			const0Node := const0.ToDualTreeNode(RandString(5))
			bst.root = nil
			bst.root = const0Node
			return nil
		} else { // Conservative Delete
			addNode := SymbolicExpression{value: "-", arity: 2}
			treeNode := bst.root.Clone()
			treeNode2 := bst.root.Clone()
			nodePlus := addNode.ToDualTreeNode(RandString(5))
			bst.root = nodePlus
			bst.root.left = &treeNode
			bst.root.right = &treeNode2
			return nil
		}
	}

	// SafeDelete will ensure the depth is greater than 1,
	// if the Tree is of depth one. It will randomly select a node and turn the node to a zero
	if deletionStrategy == 0 {
		depth, err := bst.Depth()
		if err != nil {
			return err
		}
		if depth < 2 {
			node, err := bst.RandomTerminal()
			if err != nil {
				return err
			}
			node.arity = 0
			node.value = "0"
			return nil
		}
	}

	// Malicious Delete may turn the whole Tree to zero if it selects the root as a viable place to initiate the
	// delete. If not, it will convert a subtree to 0.
	node, err := bst.RandomNonTerminal()
	if err != nil {
		return err
	}

	remove2(node)
	node.arity = 0
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

	nodes, err := bst.NonTerminals()
	if err != nil {
		return err
	}

	nonTerminalIndex0 := 0
	nonTerminalIndex1 := 0

	for nonTerminalIndex0 == nonTerminalIndex1 {

		nonTerminalIndex0 = rand.Intn(len(nodes))

		nonTerminalIndex1 = rand.Intn(len(nodes))
	}
	// once they are different

	//nodes[nonTerminalIndex0]
	return nil
}

// Replace replaces a node with replacer node. It DOES NOT Check to see if they are the same type.
// For simplicity a parent will never be nil if there is no error.
func (bst *DualTree) Replace(node *DualTreeNode, replacer DualTreeNode) (hobo DualTreeNode, parent *DualTreeNode,
	err error) {
	if bst.root == nil {
		return DualTreeNode{}, nil, fmt.Errorf("replace | treeNode you are swapping to has nil root")
	}
	if node == nil {
		return DualTreeNode{}, nil, fmt.Errorf("replace | cannot swap nil node")
	}
	if replacer.value == "" {
		return DualTreeNode{}, nil, fmt.Errorf("replace | replacer cannot have an empty value")
	}

	treeNode, parent, err := bst.Search(node.key)
	if err != nil {
		return DualTreeNode{}, nil, err
	}
	if treeNode == nil {
		return DualTreeNode{}, nil, fmt.Errorf("replace | cannot find the node to replace in Tree")
	}

	// DO
	// Apply New Keys for Each new node to be added.
	replacerTree, err := replacer.ToDualTree()
	if err != nil {
		return DualTreeNode{}, nil, err
	}
	replacerTree.InOrderTraverse(func(replacerNodes *DualTreeNode) {
		replacerNodes.key = RandString(5) // Give it a new key
	})

	if parent == nil {
		bst.root = &replacer
		return node.Clone(), bst.root, nil
	}

	if parent.left.key == node.key {
		parent.left = &replacer
	} else {
		parent.right = &replacer
	}

	return node.Clone(), parent, nil
}

// Replace replaces a node with replacer node. It DOES NOT Check to see if they are the same type,
func (bst *DualTree) ReplaceStrict(node *DualTreeNode, replacer DualTreeNode) (hobo DualTreeNode, parent *DualTreeNode,
	err error) {
	if bst.root == nil {
		return DualTreeNode{}, nil, fmt.Errorf("replace | treeNode you are swapping to has nil root")
	}
	if node == nil {
		return DualTreeNode{}, nil, fmt.Errorf("replace | cannot swap nil node")
	}
	if replacer.value == "" {
		return DualTreeNode{}, nil, fmt.Errorf("replace | replacer cannot have an empty value")
	}

	treeNode, parent, err := bst.Search(node.key)
	if err != nil {
		return DualTreeNode{}, nil, err
	}
	if treeNode == nil {
		return DualTreeNode{}, nil, fmt.Errorf("replace | cannot find the node to replace in Tree")
	}

	//childRight := treeNode.right
	//childLeft := treeNode.left

	// DO
	replacerKey := RandString(5) // Give it a new key
	replacer.key = replacerKey

	if parent == nil {
		bst.root = &replacer
		return node.Clone(), bst.root, nil
	}

	//if replacer.arity >= node.arity {
	if parent.left.key == node.key {
		parent.left = &replacer
		//if replacer.arity == 0 {
		//	parent.left.left = childLeft
		//} else if replacer.arity == 1 {
		//	parent.left.left = childLeft
		//	parent.left.right = childRight
		//} else  {
		//	parent.left.left = childLeft
		//	parent.left.right = childRight
		//}
	} else {
		parent.right = &replacer
		//if replacer.arity == 0 {
		//	parent.right.left = childLeft
		//} else if replacer.arity == 1 {
		//	parent.right.left = childLeft
		//} else  {
		//	parent.right.left = childLeft
		//	parent.right.right = childRight
		//}
	}
	//}
	//if replacer.arity < node.arity {
	//	if parent.left.key == node.key {
	//		parent.left = &replacer
	//	}else {
	//		parent.right = &replacer
	//	}
	//}

	//if replacer.arity
	//if node.arity != replacer.arity {
	//	return DualTreeNode{}, nil,
	//	fmt.Errorf("replace | replacer and node must have same arity")
	//}

	return node.Clone(), parent, nil
}

// Search will use a node Id and linearly traverse the Tree using Inorder Depth First Search until it comes across
// the correct node. It will also return the parent of the given node.
// If the Tree only contains a root and the search key matches,
// then the node is set to the root and the parent is set to nil. In the event the given key is not found,
// no error will be returned,
// therefore the onus is on the user to verify that the returned node is not nil over and above typical error handling.
func (bst *DualTree) Search(key string) (node *DualTreeNode, parent *DualTreeNode, err error) {
	if bst.root == nil {
		return nil, nil, fmt.Errorf("search | Tree root cannot be nil")
	}
	inOrderTraverseAware(bst.root, parent, func(n *DualTreeNode, parentNode *DualTreeNode) {
		if n != nil {
			if n.key == key && node == nil {
				node = n
				parent = nil
			}
			if n.left != nil {
				if n.left.key == key {
					parent = n
					node = n.left
					return
				}
			}
			if n.right != nil {
				if n.right.key == key {
					parent = n
					node = n.right
					return
				}
			}
		}
	})
	return node, parent, nil
}

func (bst *DualTree) GetRandomSubTree() (*DualTree, error) {
	if bst.root == nil {
		return nil, fmt.Errorf("treeNode you are adding to has nil root")
	}
	if bst.root.left == nil && bst.root.right == nil {
		return nil, fmt.Errorf("treeNode you are adding to is a lone terminal")
	}
	//node, err := bst.RandomNonTerminal()
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
		} else {
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

// Clone will perform an O(N) deep clone of a treeNode and its items and return its copy.
func (bst DualTree) Clone() (DualTree, error) {
	x := bst
	expressions, err := x.ToSymbolicExpressionSet()
	if err != nil {
		return DualTree{}, err
	}
	symbolicExpressionSet := expressions

	if len(symbolicExpressionSet) < 1 {
		return DualTree{}, nil
	}
	t := DualTree{}
	err = t.FromSymbolicExpressionSet2(symbolicExpressionSet)
	if err != nil {
		return DualTree{}, nil
	}

	return t, nil
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
		if tree[i].IsEqual(subTreeSlice[0]) {
			count := 0
			for j := 0; j < len(subTreeSlice); j++ {
				if !tree[i+j].IsEqual(subTreeSlice[j]) {
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
func (bst *DualTree) FromSymbolicExpressionSet2(terminalSet []SymbolicExpression) error {
	//EdgeCases
	if terminalSet == nil {
		return fmt.Errorf("terminalSet cannot be nil")
	}
	if len(terminalSet) < 1 {
		return fmt.Errorf("terminalSet cannot be empty i.e size 0")
	}
	if terminalSet[0].kind >= 1 {
		return fmt.Errorf("terminalSet cannot start with type nonterminal i.e SymbolicExpression.Kind > 1")
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
	return currentDepth - 1, nil
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
		diveTo(node.left, nodes, currDepth, maxDepth)
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
		diveUntil(node.left, nodes, currDepth, maxDepth)
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

// InOrderTraverse visits all nodes with in-order traversing but remembers its parent. (A good child :D)
func (bst *DualTree) InOrderTraverseDepthAware(f func(node *DualTreeNode, parentNode *DualTreeNode, depth *int,
	shouldReturn *bool)) {
	bst.lock.RLock()
	defer bst.lock.RUnlock()
	depth := -1
	shouldReturn := false
	inOrderTraverseAwareDepth(bst.root, bst.root, &depth, &shouldReturn, f)
}

// internal recursive function to traverse in order
func inOrderTraverseAwareDepth(n *DualTreeNode, parent *DualTreeNode, depth *int, shouldReturn *bool,
	f func(node *DualTreeNode,
		parentNode *DualTreeNode, depth *int, shouldReturn *bool)) {

	if *shouldReturn {
		return
	}
	if n != nil {
		if *shouldReturn {
			return
		}
		*depth++
		inOrderTraverseAwareDepth(n.left, n, depth, shouldReturn, f)

		f(n, parent, depth, shouldReturn)

		if *shouldReturn {
			return
		}
		inOrderTraverseAwareDepth(n.right, n, depth, shouldReturn, f)
		*depth--
	}
}

// InOrderTraverse visits all nodes with in-order traversing
func (bst *DualTree) ToSymbolicExpressionSet() ([]SymbolicExpression, error) {
	symbSet := make([]SymbolicExpression, 0)
	var err1 error
	bst.InOrderTraverse(func(node *DualTreeNode) {
		symbolicExpression, err := node.ToSymbolicExpression()
		err1 = err
		if err != nil {
			return
		}
		symbSet = append(symbSet, symbolicExpression)
	})
	return symbSet, err1
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
		inOrderTraverseAware(n.right, n, f)
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

// Print prints a visual representation of the treeNode
func (bst *DualTree) ToString() strings.Builder {
	bst.lock.Lock()
	defer bst.lock.Unlock()
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("------------------------------------------------\n"))
	stringifyBuilder(bst.root, sb, 0)
	sb.WriteString(fmt.Sprintf("------------------------------------------------\n"))

	return *sb
}

// internal recursive function to print a treeNode
func stringifyBuilder(n *DualTreeNode, sb *strings.Builder, level int) {
	if n != nil {
		format := ""
		for i := 0; i < level; i++ {
			format += "       "
			//sb.WriteString(format)
		}
		format += "---[ "
		//sb.WriteString(format)
		level++
		stringifyBuilder(n.right, sb, level)
		sb.WriteString(fmt.Sprintf(format+"%s\n", n.value))
		stringifyBuilder(n.left, sb, level)
	}

}

// ToMathematicalString returns a mathematical representation of the treeNode after reading it using Inorder DFS
func (d *DualTree) ToMathematicalString() (string, error) {
	if d.root == nil {
		return "", fmt.Errorf("treeNode root is nil cannot compute mathematical expression")
	}

	var err error = nil

	sb := &strings.Builder{}
	sb.WriteString("(")
	MathPreorder(d.root, sb)
	sb.WriteString(")")

	return sb.String(), err
}

func MathPreorder(node *DualTreeNode, sb *strings.Builder) {
	if node == nil {
		sb.WriteString("")
		return
	}

	if node.right != nil && node.left != nil {
		sb.WriteString("(")
		MathPreorder(node.left, sb)
		sb.WriteString(")")

		sb.WriteString(node.value)

		sb.WriteString("(")
		MathPreorder(node.right, sb)
		sb.WriteString(")")
		return
	}
	sb.WriteString(node.value)

	return
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

		randTerminalIndex := rand.Intn(len(terminals))
		randTerminals[i] = terminals[randTerminalIndex]
	}

	randNonTerminals := make([]SymbolicExpression, nonTerminalCount)
	for i := 0; i < nonTerminalCount; i++ {

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

	err = tree.Validate()
	return &tree, err
}

// GenerateRandomTree generates a given treeNode of a depth between 0 (i.e) root and (inclusive of) the depth specified.
// Assuming a binary structured treeNode. The number of terminals (T) is equal to 2^D where D is the depth.
// The number of NonTerminals (NT) is equal to 2^D - 1
func GenerateRandomTreeEnforceIndependentVariable(depth int, independentVar SymbolicExpression,
	terminals []SymbolicExpression,
	nonTerminals []SymbolicExpression) (*DualTree, error) {

	if depth < 0 {
		return nil, fmt.Errorf("depth cannot be less than 0")
	}
	if independentVar.value == "" {
		return nil, fmt.Errorf("independentVar cannot be empty")
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
		tree := &DualTree{}
		tree.root = independentVar.ToDualTreeNode(RandString(5))
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
	randTerminals[0] = independentVar
	for i := 1; i < terminalCount; i++ {

		randTerminalIndex := rand.Intn(len(terminals))
		randTerminals[i] = terminals[randTerminalIndex]
	}

	randNonTerminals := make([]SymbolicExpression, nonTerminalCount)
	for i := 0; i < nonTerminalCount; i++ {

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
	err = tree.Validate()
	return &tree, err
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

func (bst *DualTree) hasDiverseNonTerminalSet() (bool, error) {
	branches, err := bst.NonTerminals()
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

	randomDepth := rand.Intn(depth + 1)

	nodes, err := d.DepthTo(randomDepth)
	if err != nil {
		return DualTree{}, err
	}

	randomNodeIndex := rand.Intn(len(nodes))
	randomNode := nodes[randomNodeIndex]

	tree, err := randomNode.ToDualTree()
	if err != nil {
		return DualTree{}, err
	}

	return tree, nil
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

	randomDepth := 0
	if depth < 2 {
		randomDepth = depth
	} else {

		randomDepth = rand.Intn(depth)
	}

	nodes, err := d.DepthAt(randomDepth)
	if err != nil {
		return DualTree{}, err
	}

	randomNodeIndex := rand.Intn(len(nodes))
	randomNode := nodes[randomNodeIndex]

	tree, err := randomNode.ToDualTree()
	if err != nil {
		return DualTree{}, err
	}

	return tree, nil
}

// GetShortestBranch returns the shortest branch in a given Tree.
// The minAcceptableDepth is used to accept a node on a Tree with a small enough acceptable depth that it can be
// used. This will prevent having to check each node. As with any Inorder DFS traversal,
// nodes placed furthest right are checked last.
// If the parent is nil and there is a nil error, assume the Tree itself only contains the root.
// You have to explicitly check this. This heavily skews to nodes on the left
func (bst *DualTree) GetShortestBranch(minAcceptableDepth int) (shortestNode *DualTreeNode,
	shortestNodeParent *DualTreeNode, shortestDepth int, err error) {
	if bst.root == nil {
		return nil, nil, -1, fmt.Errorf("cannot get depth - treeNode nil")
	}
	if minAcceptableDepth < 0 {
		return nil, nil, -1, fmt.Errorf("getShortestBranch | minAcceptableDepth cannot be negative")
	}
	if minAcceptableDepth == 0 {
		minAcceptableDepth = 1
	}

	nodeDepth := struct {
		parent *DualTreeNode
		node   *DualTreeNode
		depth  int
	}{}
	depth := -1
	shouldReturn := false
	inOrderTraverseAwareDepth(bst.root, bst.root, &depth, &shouldReturn, func(n *DualTreeNode, p *DualTreeNode, d *int, shouldReturn *bool) {
		if *d <= minAcceptableDepth {
			if *shouldReturn {
				return
			}
			nodeDepth.depth = *d
			if n.IsEqual(p) {
				nodeDepth.node = n
				nodeDepth.parent = nil
				*shouldReturn = true
			} else {
				nodeDepth.node = n
				nodeDepth.parent = p
				*shouldReturn = true
			}
			return
		}
		return
	})

	return nodeDepth.node, nodeDepth.parent, nodeDepth.depth, nil
}

//func Swapper(tree1 DualTree, tree2 DualTree, subTree1 DualTree, subTree2 DualTree) (p1 DualTree, p2 DualTree,
//	c1 DualTree, c2 DualTree) {
//
//}

/**
STRATEGY IMPLEMENTORS
*/

// DeleteNonTerminal will select a non-root non-terminal element from a given tree and delete it by
// setting it to 0. If the tree only contains a root it will ignore it.
func (bst *DualTree) DeleteNonTerminal() error {
	if bst.root == nil {
		return fmt.Errorf(" DeleteNonTerminal | treeNode you are swapping to has nil root")
	}

	branches, err := bst.NonTerminals()
	if err != nil {
		return err
	}

	if len(branches) > 0 {
		randIndex := rand.Intn(len(branches))
		randomLeaf := branches[randIndex]

		randomLeaf.value = "0"
		randomLeaf.arity = 0
		randomLeaf.right = nil
		randomLeaf.left = nil
	}
	// Ignore if tree only contains terminal at root
	return nil
}

// DeleteMalicious selects any element of a tree (
// including the root) and convert it to a value of 0 potentially deleting all
// genetic material.
func (bst *DualTree) DeleteMalicious() error {
	if bst.root == nil {
		return fmt.Errorf(" DeleteMalicious | treeNode you are swapping to has nil root")
	}
	leafs, err := bst.Terminals()
	if err != nil {
		return err
	}

	randIndex := rand.Intn(len(leafs))
	randomLeaf := leafs[randIndex]

	randomLeaf.value = "0"
	randomLeaf.arity = 0
	randomLeaf.right = nil
	randomLeaf.left = nil

	return nil
}

// FellTree destroys the tree and sets its root to 0 and kills it all.
func (bst *DualTree) FellTree() error {
	if bst.root == nil {
		return fmt.Errorf(" DeleteMalicious | treeNode you are swapping to has nil root")
	}

	bst.root.value = "0"
	bst.root.arity = 0
	bst.root.right = nil
	bst.root.left = nil

	return nil
}

// DeleteTerminal will select a non-root non-terminal element from a given tree and delete it by
// setting it to 0. If the tree only contains a root it will ignore it.
func (bst *DualTree) DeleteTerminal() error {
	if bst.root == nil {
		return fmt.Errorf(" DeleteTerminal | treeNode you are swapping to has nil root")
	}
	if bst.root.right == nil && bst.root.left == nil {
		// Skip root
		return nil
	}

	leafs, err := bst.Terminals()
	if err != nil {
		return err
	}
	// ensures that at least another terminal is chosen other than the root. Hence > 1
	if len(leafs) > 1 {
		randIndex := rand.Intn(len(leafs))
		if randIndex == 0 {
			randIndex = 1
		}

		randomLeaf := leafs[randIndex]
		randomLeaf.value = "0"
		randomLeaf.arity = 0
		randomLeaf.right = nil
		randomLeaf.left = nil
	}
	return nil
}

// MutateTerminal will mutate a terminal to another valid terminal if the terminalSet only contains a single item
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

	nodes, err := bst.Terminals()
	if err != nil {
		return err
	}

	nodeValue := ""
	itemFromSet := ""

	for nodeValue == itemFromSet {

		nonTerminalIndex0 := rand.Intn(len(nodes))

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

// MutateNonTerminal will mutate a terminal to another valid nonTerminal.
// Ensure set is nonTerminal set only otherwise arities will break. If the tree is a lone terminal at the root,
// it will be ignored and the program will exit
// NOTE ensure nonTerminalSet contains no duplicates
func (bst *DualTree) MutateNonTerminal(nonTerminalSet []SymbolicExpression) error {
	if bst.root == nil {
		return fmt.Errorf("MutateNonTerminal | treeNode you are swapping to has nil root")
	}
	if bst.root.left == nil && bst.root.right == nil {
		//log.Printf("MutateNonTerminal | treeNode you are swapping to is a lone terminal not a non-terminal")
		return nil
	}
	if nonTerminalSet == nil {
		return fmt.Errorf("MutateNonTerminal | nonTerminalSet set cannot be nil")
	}
	if len(nonTerminalSet) < 1 {
		return fmt.Errorf("MutateNonTerminal | nonTerminalSet set cannot be empty")
	}

	nodes, err := bst.NonTerminals()
	if err != nil {
		return err
	}

	nodeValue := ""
	fromSetValue := ""
	counterLimit := 10

	counter := 0
	for nodeValue == fromSetValue && len(nonTerminalSet) >= 1 && counter < counterLimit { //pray for no duplicates.
		// Counter is a failsafe to prevent infinite looping

		nonTerminalIndex := rand.Intn(len(nodes))
		nonTerminalSetIndex := rand.Intn(len(nonTerminalSet))

		nodeValue = nodes[nonTerminalIndex].value
		fromSetValue = nonTerminalSet[nonTerminalSetIndex].value

		if nodeValue == fromSetValue {
			if counter == counterLimit-1 {
				break
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

// ReplaceBranch takes a given tree and randomly selects a branch i.
// e non-terminal and will swap it with a randomly generated tree of variable depth. This includes the root
func (bst *DualTree) ReplaceBranch(tree DualTree) error {
	if bst.root == nil {
		return fmt.Errorf(" ReplaceBranch | treeNode you are swapping to has nil root")
	}
	if tree.root == nil {
		return fmt.Errorf(" ReplaceBranch | treeNode you are swapping with has nil root")
	}

	node, parent, err := bst.RandomNonTerminalAware()
	if err != nil {
		return err
	}
	if node == nil { // If the tree is a lone terminal. Swap it out with the incoming tree
		bst.root = tree.root
	}
	if parent == nil { // If the tree is a terminal
		bst.root = tree.root
	} else {
		if parent.left.key == node.key {
			parent.left = tree.root
		} else if parent.right.key == node.key {
			parent.right = tree.root
		} else {
			return fmt.Errorf("ReplaceBranch | Parent is not a parent of the node %#v", node)
		}
	}

	return nil
}

// AddToLeaf is similar to AddSubTree, however the SubTree will only be placed on a randomly selected leaf. It will not replace a non-terminal
func (bst *DualTree) AddToLeaf(tree DualTree) error {
	if bst.root == nil {
		return fmt.Errorf(" AddToLeaf | treeNode you are swapping to has nil root")
	}
	if tree.root == nil {
		return fmt.Errorf(" AddToLeaf | treeNode you are swapping with has nil root")
	}
	if tree.root.left == nil && tree.root.right == nil {
		return fmt.Errorf(" AddToLeaf | subTree cannot be composed of a single terminal - no operation to add it to" +
			" the treeNode.")
	}

	node, parent, err := bst.RandomTerminalAware()
	if err != nil {
		return err
	}
	if parent == nil { // If the tree is a lone terminal
		bst.root = tree.root
	} else {
		if parent.left.key == node.key {
			parent.left = tree.root
		} else if parent.right.key == node.key {
			parent.right = tree.root
		} else {
			return fmt.Errorf("AddToLeaf | Parent is not a parent of the node %#v", node)
		}
	}

	return nil
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
		err := bst.AddEmptyToTreeRoot(subTree)
		return err
	}

	node, err := bst.RandomNonTerminal()
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

	return nil
}
