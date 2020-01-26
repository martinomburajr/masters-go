package evolution

import (
	"fmt"
	"math/rand"
)

func Mutate(outgoingParents []*Individual, children []*Individual, kind int,
	opts EvolutionParams) (parents []*Individual, childs []*Individual, err error) {
	if kind == IndividualAntagonist {
		for i := 0; i < (len(outgoingParents)); i++ {
			probabilityOfMutation := rand.Float64()
			if probabilityOfMutation < opts.Reproduction.ProbabilityOfMutation {
				err := outgoingParents[i].Mutate(opts.Strategies.AntagonistAvailableStrategies)
				if err != nil {
					return nil, nil, err
				}
			}
		}
		// childs
		for i := 0; i < (len(children)); i++ {
			probabilityOfMutation := rand.Float64()
			if probabilityOfMutation < opts.Reproduction.ProbabilityOfMutation {
				err := children[i].Mutate(opts.Strategies.AntagonistAvailableStrategies)
				if err != nil {
					return nil, nil, err
				}
			}
		}
	} else if kind == IndividualProtagonist {
		for i := 0; i < (len(outgoingParents)); i++ {
			probabilityOfMutation := rand.Float64()
			if probabilityOfMutation < opts.Reproduction.ProbabilityOfMutation {
				err := outgoingParents[i].Mutate(opts.Strategies.ProtagonistAvailableStrategies)
				if err != nil {
					return nil, nil, err
				}
			}
		}
		// childs
		for i := 0; i < (len(children)); i++ {
			probabilityOfMutation := rand.Float64()
			if probabilityOfMutation < opts.Reproduction.ProbabilityOfMutation {
				err := children[i].Mutate(opts.Strategies.ProtagonistAvailableStrategies)
				if err != nil {
					return nil, nil, err
				}
			}
		}
	} else {
		return nil, nil, fmt.Errorf("Judgement Day | Invalid kind")
	}
	return outgoingParents, children, nil
}