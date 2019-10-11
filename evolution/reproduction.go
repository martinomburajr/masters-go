package evolution

import (
	"fmt"
	"math/rand"
)

// CrossoverTree is a evolutionary technique used to take two parents swap their genetic material and form two new children.
func CrossoverTree(individual1 *Individual, individual2 *Individual, maxDepth int, params EvolutionParams) (child1 Individual,
	child2 Individual,
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
	cloneA, err := individual1.Clone()
	if err != nil {
		return Individual{}, Individual{}, err
	}
	cloneB, err := individual2.Clone()
	if err != nil {
		return Individual{}, Individual{}, err
	}

	cloneATree := cloneA.Program.T
	cloneBTree := cloneB.Program.T

	cloneADepth, err := cloneATree.Depth()
	if err != nil {
		return Individual{}, Individual{}, err
	}

	cloneBDepth, err := cloneBTree.Depth()
	if err != nil {
		return Individual{}, Individual{}, err
	}

	// Check Depths for Swap

	// 1. If depths < 1 in case it is just a Tree with only a root
	if cloneADepth < 1 {
		nodeB, _, err := cloneBTree.RandomLeafAware()
		if err != nil {
			return Individual{}, Individual{}, err
		}
		hoboA, _, err := cloneATree.Replace(cloneATree.root, *nodeB)
		if err != nil {
			return Individual{}, Individual{}, err
		}
		_, _, err = cloneBTree.Replace(nodeB, hoboA)
		if err != nil {
			return Individual{}, Individual{}, err
		}
		return cloneA, cloneB, err
	}

	if cloneBDepth < 1 {
		nodeA, _, err := cloneATree.RandomLeafAware()
		if err != nil {
			return Individual{}, Individual{}, err
		}
		hoboB, _, err := cloneBTree.Replace(cloneATree.root, *nodeA)
		if err != nil {
			return Individual{}, Individual{}, err
		}
		_, _, err = cloneATree.Replace(nodeA, hoboB)
		if err != nil {
			return Individual{}, Individual{}, err
		}
		return cloneA, cloneB, err
	}

	shortestNodeA, _, shortestDepthA, err := cloneATree.GetShortestBranch(maxDepth / 2)
	if err != nil {
		return Individual{}, Individual{}, nil
	}
	shortestNodeB, _, shortestDepthB, err := cloneBTree.GetShortestBranch(maxDepth / 2)
	if err != nil {
		return Individual{}, Individual{}, nil
	}

	if cloneADepth > maxDepth {
		if shortestDepthA >= maxDepth {
			// Penalize Parent
			penalty := int(params.DepthPenaltyStrategyPenalization) * (shortestDepthA / maxDepth)
			if individual1.HasCalculatedFitness {
				return Individual{}, Individual{}, fmt.Errorf("cannot be penalized | Fitness uncalculated")
			}
			individual1.TotalFitness = individual1.TotalFitness + int(penalty)
		}
	}
	if cloneBDepth > maxDepth {
		if shortestDepthB >= maxDepth {
			// Penalize Parent
			penalty := int(params.DepthPenaltyStrategyPenalization) * (shortestDepthB / maxDepth)
			if individual2.HasCalculatedFitness {
				return Individual{}, Individual{}, fmt.Errorf("cannot be penalized | Fitness uncalculated")
			}
			individual2.TotalFitness = individual2.TotalFitness + int(penalty)
		}
	}

	if shortestDepthA <= shortestDepthB {
		subTreeBAtDepth, err := cloneBTree.GetRandomSubTreeAtDepthAware(cloneBDepth) // confirm
		if err != nil {
			return Individual{}, Individual{}, err
		}
		hoboA, _, err := cloneATree.Replace(shortestNodeA, *subTreeBAtDepth.root)
		if err != nil {
			return Individual{}, Individual{}, err
		}
		_, _, err = cloneBTree.Replace(subTreeBAtDepth.root, hoboA)
		if err != nil {
			return Individual{}, Individual{}, err
		}
	} else {
		subTreeAAtDepth, err := cloneATree.GetRandomSubTreeAtDepthAware(cloneADepth) // confirm
		if err != nil {
			return Individual{}, Individual{}, err
		}
		hoboB, _, err := cloneBTree.Replace(shortestNodeB, *subTreeAAtDepth.root)
		if err != nil {
			return Individual{}, Individual{}, err
		}
		_, _, err = cloneATree.Replace(subTreeAAtDepth.root, hoboB)
		if err != nil {
			return Individual{}, Individual{}, err
		}
	}

	cloneA.Program.T = cloneATree
	cloneB.Program.T = cloneBTree
	return cloneA, cloneB, err
}

// getRandomDepthTargetLocation obtains a random depth for each individual that the crossover will target.
// For example if the depth of individual1 is 10,
// and the max depth is 15. The remainder will be 5.
// This function takes the remainder and gets a random number between 0 and 5 of individual2
func getRandomDepthTargetLocation(individual1DepthRemainderFromMaX int, individual2DepthRemainderFromMax int) (int, int) {
	var randDepthOfInd1FromRem, randDepthOfInd2FromRem int
	if individual1DepthRemainderFromMaX != 0 {

		randDepthOfInd1FromRem = rand.Intn(individual1DepthRemainderFromMaX)
	}
	if individual2DepthRemainderFromMax != 0 {

		randDepthOfInd2FromRem = rand.Intn(individual2DepthRemainderFromMax)
	}
	return randDepthOfInd1FromRem, randDepthOfInd2FromRem
}

// calculateRemainderDepths returns the remainders. It does no checking to ensure individuals are correct.
// This has to be done by the user.
func calculateRemainderDepths(individual1 *Individual, individual2 *Individual, maxDepth int,
	params EvolutionParams) (int, int, error) {

	individual1Depth, err := individual1.Program.T.Depth()
	if err != nil {
		return -1, -1, err
	}

	individual2Depth, err := individual2.Program.T.Depth()
	if err != nil {
		return -1, -1, err
	}

	if params.DepthPenaltyStrategy == DepthPenaltyStrategyIgnore {
		i, i2 := depthPenaltyIgnore(maxDepth, individual1Depth, individual2Depth)
		return i, i2, nil
	}
	if params.DepthPenaltyStrategy == DepthPenaltyStrategyPenalize {
		return depthPenaltyPenalization(individual1, individual2, individual1Depth, individual2Depth, maxDepth,
			params.DepthPenaltyStrategyPenalization)
	}
	i, i2 := depthPenaltyIgnore(maxDepth, individual1Depth, individual2Depth)
	return i, i2, nil
}

func depthPenaltyIgnore(maxDepth int, individual1Depth int, individual2Depth int) (int, int) {
	if maxDepth < 0 {
		maxDepth = 0
	}
	var individual1DepthRemainderFromMaX, individual2DepthRemainderFromMax int
	if individual1Depth >= maxDepth {
		individual1DepthRemainderFromMaX = 0
	} else {
		individual1DepthRemainderFromMaX = maxDepth - individual1Depth
	}
	if individual2Depth >= maxDepth {
		individual2DepthRemainderFromMax = 0
	} else {
		individual2DepthRemainderFromMax = maxDepth - individual2Depth
	}
	return individual1DepthRemainderFromMaX, individual2DepthRemainderFromMax
}

// depthPenaltyPenalization applies a penalty to an individual whose depth exceeds maxDepth.
// Ensure that the individual has calculated its Fitness
func depthPenaltyPenalization(individual1 *Individual, individual2 *Individual, individual1Depth int,
	individual2Depth int, maxDepth int,
	penalization float64) (int, int, error) {
	if maxDepth < 0 {
		maxDepth = 0
	}
	var individual1DepthRemainderFromMaX, individual2DepthRemainderFromMax int
	if individual1Depth >= maxDepth {
		if individual1.HasCalculatedFitness {
			individual1.TotalFitness = individual1.TotalFitness + int(penalization)
		} else {
			return -1, -1, fmt.Errorf("crossover | depthPenalty | Fitness of individual %s has not been calculated"+
				" before crossover", individual1.Id)
		}
	} else {
		individual1DepthRemainderFromMaX = maxDepth - individual1Depth
	}

	if individual2Depth >= maxDepth {
		if individual2.HasCalculatedFitness {
			individual2.TotalFitness = individual2.TotalFitness + int(penalization)
		} else {
			return -1, -1, fmt.Errorf("crossover | depthPenalty | Fitness of individual %s has not been calculated"+
				" before crossover", individual1.Id)
		}
	} else {
		individual2DepthRemainderFromMax = maxDepth - individual2Depth
	}
	return individual1DepthRemainderFromMaX, individual2DepthRemainderFromMax, nil
}

//func depthPenaltyTrim(individual1 *Individual, individual2 *Individual, individual1Depth int,
//	individual2Depth int, maxDepth int,
//	penalization float64)(int, int, error) {
//	if maxDepth < 0 {
//		maxDepth = 0
//	}
//	var individual1DepthRemainderFromMaX, individual2DepthRemainderFromMax int
//	if individual1Depth >= maxDepth {
//		if individual1.HasCalculatedFitness {
//			individual1.
//		}else {
//			return -1, -1, fmt.Errorf("crossover | depthPenalty | Fitness of individual %s has not been calculated" +
//				" before crossover", individual1.Id)
//		}
//	}else {
//		individual1DepthRemainderFromMaX = maxDepth - individual1Depth
//	}
//
//	if individual2Depth >= maxDepth {
//		if individual2.HasCalculatedFitness {
//			individual2.TotalFitness = individual2.TotalFitness + int(penalization)
//		}else {
//			return -1, -1, fmt.Errorf("crossover | depthPenalty | Fitness of individual %s has not been calculated" +
//				" before crossover", individual1.Id)
//		}
//	}else {
//		individual2DepthRemainderFromMax = maxDepth - individual2Depth
//	}
//	return individual1DepthRemainderFromMaX, individual2DepthRemainderFromMax, nil
//}

// 1. Depth and Remainder Information
//individual1DepthRemainderFromMaX, individual2DepthRemainderFromMax, err := calculateRemainderDepths(individual1, individual2, maxDepth, params)
//if err != nil {
//	return Individual{}, Individual{}, err
//}
//
//// 3. Calculate a random depth in the treeNode to extract information from
//randDepthOfInd1FromRem, randDepthOfInd2FromRem := getRandomDepthTargetLocation(individual1DepthRemainderFromMaX, individual2DepthRemainderFromMax)
//
//// 4. Get Random SubTree based on depth of each program
//subTreeAtDepthProg1, err := individual1.Program.T.GetRandomSubTreeAtDepth(randDepthOfInd1FromRem)
//if err != nil {
//	return Individual{}, Individual{}, err
//}
//subTreeAtDepthProg2, err := individual2.Program.T.GetRandomSubTreeAtDepth(randDepthOfInd2FromRem)
//if err != nil {
//	return Individual{}, Individual{}, err
//}
//
//// 5. Get depth of randomly selected subTrees
//subTreeInd1Depth, err := subTreeAtDepthProg1.Depth()
//if err != nil {
//	return Individual{}, Individual{}, err
//}
//subTreeInd2Depth, err := subTreeAtDepthProg2.Depth()
//if err != nil {
//	return Individual{}, Individual{}, err
//}
//
//// 6. Select Nodes with A Given Depth
//nodesProg1, err := subTreeAtDepthProg1.DepthAt(subTreeInd1Depth)
//if err != nil {
//	return Individual{}, Individual{}, err
//}
//nodesProg2, err := subTreeAtDepthProg2.DepthAt(subTreeInd2Depth)
//if err != nil {
//	return Individual{}, Individual{}, err
//}
//
//
//var nodesProg1RandomIndex, nodesProg2RandomIndex int
//if len(nodesProg1) > 0 {
//
//	nodesProg1RandomIndex = rand.Intn(len(nodesProg1))
//}
//if len(nodesProg2) > 0 {
//
//	nodesProg2RandomIndex = rand.Intn(len(nodesProg2))
//}
//
//// 7. Get random node in subTree
//// 7.a node1 = rand(nodesInd1SubTree)
//// 7.b node2 = rand(nodesInd2SubTree)
//nodeProg1Random := nodesProg1[nodesProg1RandomIndex]
//nodeProg2Random := nodesProg2[nodesProg2RandomIndex]
//
//nodeProg1PieceTree := nodeProg1Random.ToDualTree()
//nodeProg2PieceTree := nodeProg2Random.ToDualTree()
//
//child1 = individual1.Clone()
//child2 = individual2.Clone()
//
//// Child 1
//node, parent, err := child1.Program.T.Search(nodeProg1PieceTree.root.key)
//if err != nil {
//	return Individual{}, Individual{}, err
//}
//if node == nil {
//	return Individual{}, Individual{}, fmt.Errorf("crossover | failed to locate node in the Tree it came from..." +
//		" weird error")
//}
//if parent == nil {
//	return child1, child2, nil
//}
//if parent.right.key == nodeProg1PieceTree.root.key {
//	parent.right = nodeProg1PieceTree.root
//} else if parent.left.key == nodeProg1PieceTree.root.key {
//	parent.left = nodeProg1PieceTree.root
//}
//
//// Child 2
//node2, parent2, err := child2.Program.T.Search(nodeProg2PieceTree.root.key)
//if err != nil {
//	return Individual{}, Individual{}, err
//}
//if node2 == nil {
//	return Individual{}, Individual{}, fmt.Errorf("crossover | failed to locate node in the Tree it came from..." +
//		" weird error")
//}
//if parent2 == nil {
//	return Individual{}, Individual{}, fmt.Errorf("crossover | failed to locate parent in the Tree it came from.." +
//		"." +
//		" weird error")
//}
//if parent2.right.key == nodeProg2PieceTree.root.key {
//	parent2.right = nodeProg2PieceTree.root
//} else if parent.left.key == nodeProg2PieceTree.root.key {
//	parent2.left = nodeProg2PieceTree.root
//}
//
