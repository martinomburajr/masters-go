package evolution

type Strategable interface{ Apply(t *DualTree) }

type Strategy string

const (
	StrategyAddSubTree        = "StrategyAddSubTree"
	StrategyDeleteSubTree     = "StrategyDeleteSubTree"
	StrategyDeleteMalicious   = "StrategyDeleteMalicious"
	StrategySoftDeleteSubTree = "StrategySoftDeleteSubTree"
	StrategySwapSubTree       = "StrategySwapSubTree"
	StrategyMutateNode        = "StrategyMutateNode"
	StrategyMutateSubTree     = "StrategyMutateSubTree"
)
