package evolution

/**
	Any strategy operation below will ensure the tree remains in a valid state.
Worst case being a single terminal with value 0.
*/

type Strategable interface{ Apply(t *DualTree) }

type Strategy string



const (


	// #############################  Delete Strategies ############################################
	// All delete operations will still allow the tree to remain in a valid state.
	// Worst case scenario the resulting tree will have a root of terminal value 0.

	// StrategyDeleteNonTerminal will select a non-root non-terminal element from a given tree and delete it by
	// setting it to 0
	StrategyDeleteNonTerminal     = "StrategyDeleteNonTerminal"
	// StrategyDeleteMalicious selects any element of a tree (
	// including the root) and convert it to a value of 0 potentially deleting all
	// genetic material.
	StrategyDeleteMalicious   = "StrategyDeleteMalicious"
	// StrategyDeleteTerminal will convert a terminal node to 0.
	StrategyDeleteTerminal = "StrategySoftDeleteSubTree"

	// StrategyFellTree destroys the tree and sets its root to 0 and kills it all.
	StrategyFellTree = "StrategyFellTree"

	// #############################  Mutate Strategies ############################################

	// StrategyMutateNode randomly selects a non-terminal in a tree and changes its value to one of the available
	// nonterminals in the parameter list.
	// If the tree only contains a root that is a terminal it will ignore it.
	StrategyMutateNonTerminal       = "StrategyMutateNonTerminal"

	// StrategyMutateTerminal randomly selects a terminal in a tree and changes its value to one of the available
	// terminals in the parameter list.
	// If the tree only contains a root that is a terminal it will ignore it.
	StrategyMutateTerminal = "StrategyMutateTerminal"

	// #############################  Replace Strategies ############################################

	// StrategyReplaceBranch takes a given tree and randomly selects a branch i.
	// e non-terminal and will swap it with a randomly generated tree of variable depth
	StrategyReplaceBranch = "StrategyReplaceBranch"

	// #############################  Add Strategies ############################################
	// If an add strategy encounters a 0 at the root, it will replace the 0.

	//StrategyAddSubTree is a generic strategy that adds a randomly generated subtree anywhere on a given tree
	StrategyAddSubTree        = "StrategyAddSubTree"
	// StrategyAddToLeaf is similar to AddSubTree,
	// however the SubTree will only be placed on a randomly selected leaf. It will not replace a non-terminal.
	// It can replace a root
	StrategyAddToLeaf = "StrategyAddToLeaf"
	// StrategyAddMult will add a subTree with a root of multiplication to a given leaf node
	StrategyAddMult = "StrategyAddMult"
	// StrategyAddMult will add a subTree with a root of subtract to a given leaf node
	StrategyAddSub = "StrategyAddSub"
	// StrategyAddMult will add a subTree with a root of add to a given leaf node
	StrategyAddAdd = "StrategyAddAdd"

)
