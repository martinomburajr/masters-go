package evolution

type Strategable interface{ Apply(t *DualTree) }

type Strategy struct {
	Kind   string
	Action func(program *Program) *Program
}

// NewStrategy creates a new strategy.
func NewStrategy(kind string, action func(program *Program) *Program) Strategy {
	return Strategy{kind, action}
}

const (
	AddSubTree        = "AddSubTree"
	DeleteSubTree     = "DeleteSubTree"
	SoftDeleteSubTree = "SoftDeleteSubTree"
	SwapSubTree       = "SwapSubTree"
	MutateNode        = "MutateNode"
	MutateSubTree     = "MutateSubTree"
)
