package evolution

/**
	Any strategy operation below will ensure the tree remains in a valid state.
Worst case being a single terminal with value 0.
*/

type Strategable interface{ Apply(t *DualTree) }

type Strategy string

const (

	// ############################# NON DETERMINISTIC ############################################
	// #############################  Delete Strategies ############################################
	// All delete operations will still allow the tree to remain in a valid state.
	// Worst case scenario the resulting tree will have a root of terminal value 0.

	// StrategyDeleteNonTerminal will select a non-root non-terminal element from a given tree and delete it by
	// setting it to 0
	StrategyDeleteNonTerminal = "DeleteNonTerminalR"
	// StrategyDeleteMalicious randomly selects any element of a tree (
	// including the root) and convert it to a value of 0 potentially deleting all
	// genetic material if the root is selected
	StrategyDeleteMalicious = "DeleteMaliciousR"
	// StrategyDeleteTerminal will convert a terminal node to 0.
	StrategyDeleteTerminal = "DeleteTerminalR"
	// StrategyMutateNode randomly selects a non-terminal in a tree and changes its value to one of the available
	// nonterminals in the parameter list.
	// If the tree only contains a root that is a terminal it will ignore it.
	StrategyMutateNonTerminal = "MutateNonTerminalR"
	// StrategyMutateTerminal randomly selects a terminal in a tree and changes its value to one of the available
	// terminals in the parameter list.
	// If the tree only contains a root that is a terminal it will ignore it.
	StrategyMutateTerminal = "MutateTerminalR"
	// StrategyReplaceBranch takes a given tree and randomly selects a branch i.
	// e non-terminal and will swap it with a randomly generated tree of variable depth
	StrategyReplaceBranch  = "ReplaceBranchR"
	StrategyReplaceBranchX = "ReplaceBranchXR"
	//StrategyAddRandomSubTree is a generic strategy that adds a randomly generated subtree anywhere on a given tree
	//  If an add strategy encounters a 0 at the root, it will replace the 0.
	StrategyAddRandomSubTree = "AddRandomSubTreeR"
	// StrategyAddToLeaf is similar to AddSubTree,
	// however the SubTree will only be placed on a randomly selected leaf. It will not replace a non-terminal.
	// It can replace a root
	StrategyAddToLeaf  = "AddToLeafR"
	StrategyAddToLeafX = "AddToLeafX"
	// StrategyAddTreeWithMult will add a subTree with a root of multiplication to a given leaf node
	StrategyAddTreeWithMult = "AddTreeWithMult"
	// StrategyAddTreeWithMult will add a subTree with a root of subtract to a given leaf node
	StrategyAddTreeWithSub = "AddTreeWithSub"
	// StrategyAddTreeWithMult will add a subTree with a root of add to a given leaf node
	StrategyAddTreeWithAdd = "AddTreeWithAdd"
	StrategyAddTreeWithDiv = "AddTreeWithDiv"

	// ####################################################### DETERMINISTIC STRATEGIES #############################
	// StrategySkip performs no operations on the given subtree.
	StrategySkip = "SkipD"
	// StrategyFellTree destroys the tree and sets its root to 0 and kills it all.
	StrategyFellTree = "FellTreeD"
	StrategyMultXD   = "MultXD"
	StrategyAddXD    = "AddXD"
	StrategySubXD    = "SubXD"
	StrategyDivXD    = "DivXD"

	//Strategy
)
