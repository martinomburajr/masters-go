package main

import (
	"github.com/martinomburajr/masters-go/program"
)

type Strategable interface{ Apply(t *program.DualTree) }

type Strategy struct {
	Name   string
	Action func(program *program.Program) program.Program
}

// NewStrategy creates a new strategy.
func NewStrategy(name string, action func(program *program.Program) program.Program) Strategy {
	return Strategy{name, action}
}
