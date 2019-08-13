package evolution

import (
	"github.com/martinomburajr/masters-go/program/tree/dualtree"
)

type Strategable interface{ Apply(t *dualtree.DualTree) }

type Strategy struct {
	Name   string
	Action func(program *Program) Program
}

// NewStrategy creates a new strategy.
func NewStrategy(name string, action func(program *Program) Program) Strategy {
	return Strategy{name, action}
}
