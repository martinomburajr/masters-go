package evolution

import (
	"fmt"
)

const (
	MinAllowableGenerationsToTerminate  = 9
	TopologyRoundRobin                  = "TopologyRoundRobin"
	TopologyKRandom                     = "TopologyKRandom"
	TopologyHallOfFame                  = "TopologyHallOfFame"
	TopologySingleEliminationTournament = "TopologySET"
)

type ITopology interface {
	Topology(currentGeneration *Generation, params EvolutionParams) (*Generation, error)
}

type IEvolve interface {
	Evolve(params EvolutionParams, topology ITopology) (*EvolutionResult, error)
}

func (engine *EvolutionEngine) Evolve(params EvolutionParams) (*EvolutionResult, error) {
	switch engine.Parameters.Topology.Type {
	case TopologyHallOfFame:
		hallOfFame := &HallOfFame{Engine: engine}
		evolutionResult, err := hallOfFame.Evolve(params, hallOfFame)
		if err != nil {
			return nil, err
		}
		return evolutionResult, nil
	case TopologyKRandom:
		kRandom := &KRandom{Engine: engine}
		evolutionResult, err := kRandom.Evolve(params, kRandom)
		if err != nil {
			return nil, err
		}
		return evolutionResult, nil
	case TopologyRoundRobin:
		roundRobin := &RoundRobin{Engine: engine}
		evolutionResult, err := roundRobin.Evolve(params, roundRobin)
		if err != nil {
			return nil, err
		}
		return evolutionResult, nil
	case TopologySingleEliminationTournament:
		singleEliminationTournament := &SingleEliminationTournamentTopology{Engine: engine}
		evolutionResult, err := singleEliminationTournament.Evolve(params, singleEliminationTournament)
		if err != nil {
			return nil, err
		}
		return evolutionResult, nil
	default:
		return nil, fmt.Errorf("Compete | invalid Evolutionary Topology set")
	}
}
