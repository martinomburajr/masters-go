package evolution

type Strategable interface{ Apply(t *DualTree) }

type Strategy string

const (
	StrategyAddSubTree        = "StrategyAddSubTree"
	StrategyDeleteSubTree     = "StrategyDeleteSubTree"
	StrategySoftDeleteSubTree = "StrategySoftDeleteSubTree"
	StrategySwapSubTree       = "StrategySwapSubTree"
	StrategyMutateNode        = "StrategyMutateNode"
	StrategyMutateSubTree     = "StrategyMutateSubTree"
)
