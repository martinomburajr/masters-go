package program

import "github.com/martinomburajr/masters-go/program/tree"

type Terminal tree.DualTreeNode
type NonTerminal tree.DualTreeNode

// TODO generate AST tree from polynomial expression
type Program struct {
	ID                   string
	T                    *tree.DualTree
	Strategies           []Strategable
	hasAppliedStrategies bool
	generation           *Generation
	Evaluable
	Fitnessable
	ApplyStrategeable
}

func (p *Program) ApplyStrategy() {
	return
}

func (p *Program) Fitness() float32 {
	return p.generation.engine.Fitness()
}

func (p *Program) Eval() float32 {

}

func (p *Program) Terminals() []*Terminal {
	return nil
}

func (p *Program) NonTerminals() []*NonTerminal {
	return nil
}

func (p *Program) Mutate() {

}

func (p *Program) Recombine() {

}

func (p *Program) Validate() error {
	return nil
}

type Bug *Program
type Test *Program

type InitialProgram struct {
	ID   string
	T    *tree.DualTree
	spec Spec
}

func (p *InitialProgram) Spec(spec Spec) *InitialProgram {
	p.spec = spec
	return p
}
