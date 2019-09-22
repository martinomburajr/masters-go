package evolution

import (
	"fmt"
	"math/rand"
	"time"
)

// Crossover is a evolutionary technique used to take two parents swap their genetic material and form two new children.
func Crossover(individual1 *Individual, individual2 *Individual, maxDepth int) (child1 Individual, child2 Individual,
	err error) {
		// Requirements
	if individual1 == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover: individual 1 cannot be nil")
	}
	if individual1.Program == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover: individual1.Program != nil")
	}
	if individual1.Program.T == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover: individual1.Program.T != nil")
	}
	if individual1.Program.T.root == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover: individual1.Program.T.root != nil")
	}
	if individual2 == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover: individual2 cannot be nil")
	}
	if individual2.Program == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover: individual2.Program != nil")
	}
	if individual2.Program.T == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover: individual2.Program.T != nil")
	}
	if individual2.Program.T.root == nil {
		return Individual{}, Individual{}, fmt.Errorf("crossover: individual2.Program.T.root != nil")
	}
	if maxDepth < 0 {
		return Individual{}, Individual{}, fmt.Errorf("crossover: max depth cannot be less than 0")
	}

	// DO!

	// 1. Depth Information
	prog1Depth, err := individual1.Program.T.Depth()
	if err != nil {
		return Individual{}, Individual{}, nil
	}

	prog2Depth, err := individual2.Program.T.Depth()
	if err != nil {
		return Individual{}, Individual{}, nil
	}

	// 2. Calculate depth remainders
	prog1DepthRem := maxDepth - prog1Depth
	prog2DepthRem := maxDepth - prog2Depth

	// 3. Calculate a random depth in the treeNode to extract information from
	var randDepthProg1, randDepthProg2 int
	if prog1DepthRem != 0 {
		rand.Seed(time.Now().UnixNano())
		randDepthProg1 = rand.Intn(prog1DepthRem)
	}
	if prog1DepthRem != 0 {
		rand.Seed(time.Now().UnixNano())
		randDepthProg2 = rand.Intn(prog2DepthRem)
	}

	// 4. Get Random SubTree based on depth of each program
	subTreeAtDepthProg1, err := individual1.Program.T.GetRandomSubTreeAtDepth(randDepthProg1)
	if err != nil {
		return Individual{}, Individual{}, err
	}
	subTreeAtDepthProg2, err := individual2.Program.T.GetRandomSubTreeAtDepth(randDepthProg2)
	if err != nil {
		return Individual{}, Individual{}, err
	}

	// 5. Get depth of randomly selected subTrees
	subTreeInd1Depth, err := subTreeAtDepthProg1.Depth()
	if err != nil {
		return Individual{}, Individual{}, err
	}
	subTreeInd2Depth, err := subTreeAtDepthProg2.Depth()
	if err != nil {
		return Individual{}, Individual{}, err
	}

	// 6. Select Nodes with A Given Depth
	nodesProg1, err := subTreeAtDepthProg1.DepthAt(maxDepth - subTreeInd1Depth)
	if err != nil {
		return Individual{}, Individual{}, err
	}
	nodesProg2, err := subTreeAtDepthProg2.DepthAt(maxDepth - subTreeInd2Depth)
	if err != nil {
		return Individual{}, Individual{}, err
	}

	rand.Seed(time.Now().UnixNano())
	nodesProg1RandomIndex := rand.Intn(len(nodesProg1))
	rand.Seed(time.Now().UnixNano())
	nodesProg2RandomIndex := rand.Intn(len(nodesProg2))

	// 7. Get random node in subTree
	// 7.a node1 = rand(nodesInd1SubTree)
	// 7.b node2 = rand(nodesInd2SubTree)
	nodeProg1Random := nodesProg1[nodesProg1RandomIndex]
	nodeProg2Random := nodesProg2[nodesProg2RandomIndex]

	nodeProg1Random.ToDualTree()
	nodeProg2Random.ToDualTree()

	return Individual{}, Individual{}, nil
}


