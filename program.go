package main

type Tree struct {
	rootNode *Node
	Print string
}

type Node struct {
	left *Node
	right *Node
	item string
}

type Terminal Node
type NonTerminal Node


type Program struct {
	ID string
	T *Tree
	Arity int
	Strategies []*Strategy
	Evaluable
	Fitnessable
	ApplyStrategeable
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