package evolution

import (
	"math/rand"
	"time"
)

// Crossover is a evolutionary technique used to take two parents swap their genetic material and form two new children.
func Crossover(individual1 Individual, individual2 Individual, maxDepth int) (child1 Individual, child2 Individual,
	err error) {

	// 1. Depth Information
		// 1.a prog1Depth = Get Depth from prog1
		// 1.b Get Depth from  prog2

	prog1Depth, err := individual1.Program.T.Depth()
	if err != nil {
		return Individual{}, Individual{}, nil
	}

	prog2Depth, err := individual2.Program.T.Depth()
	if err != nil {
		return Individual{}, Individual{}, nil
	}

	// 2. Calculate depth remainders
		// 2.a Get remainder for maxDepth - depth(prog1)
		// 2.b Get remainder for maxDepth - depth(prog2)

	prog1DepthRem := maxDepth - prog1Depth
	prog2DepthRem := maxDepth - prog2Depth

	// 3. Calculate a random depth in the tree to extract information from
		// 3.a Get randomDepth for prog1
		// 3.b Get randomDepth for prog2

	rand.Seed(time.Now().UnixNano())
	randDepthProg1 := rand.Intn(prog1DepthRem)

	rand.Seed(time.Now().UnixNano())
	randDepthProg2 := rand.Intn(prog2DepthRem)

	// 4. Get Random SubTree based on depth of each program
		// 4.a subTreeInd1 = getRandomSubTree(prog1Depth)
		// 4.b subTreeInd2 = getRandomSubTree(prog2Depth)

	individual1.Program.T.GetRandomSubTreeAtDepth(randDepthProg1)
	individual2.Program.T.GetRandomSubTreeAtDepth(randDepthProg2)

	// 5. Get depth of randomly selected subTrees
		//	5.a subTreeInd1Depth = subTreeInd1.depth
		//	5.b subTreeInd2Depth = subTreeInd2.depth

	// 6. Select Nodes with A Given Depth
		// 6.a nodesInd1SubTree = SelectNodesWithDepth(Ind1, maxDepth - subTreeInd1Depth -1)
		// 6.b nodesInd2SubTree = SelectNodesWithDepth(Ind2, maxDepth - subTreeInd2Depth -1)

	// 7. Get random node in subTree
	// 7.a node1 = rand(nodesInd1SubTree)
	// 7.b node2 = rand(nodesInd2SubTree)






	return Individual{}, Individual{}, nil

}

