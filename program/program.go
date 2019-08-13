package program

import (
	"github.com/martinomburajr/masters-go/program/tree/dualtree"
)

// TODO generate AST tree from polynomial expression
type Program struct {
	ID                   string
	T                    *dualtree.DualTree
	Strategies           []Strategable
	hasAppliedStrategies bool
	generation           *Generation
}

func (p *Program) ApplyStrategy() {
	return
}

func (p *Program) Fitness() float32 {
	return p.generation.engine.Fitness()
}

func (p *Program) Eval() float32 {
	return 0
}


func (p *Program) Terminals() []*dualtree.Terminal {
	return nil
}

func (p *Program) NonTerminals() []*dualtree.NodeType {
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
	T    *dualtree.DualTree
	spec Spec
}

func (p *InitialProgram) Spec(spec Spec) *InitialProgram {
	p.spec = spec
	return p
}

func (p *InitialProgram) Validate() error {

}