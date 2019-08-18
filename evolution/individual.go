package evolution

type Individual struct {
	id                 string
	strategy           []*Strategy
	fitness            []float32
	hasAppliedStrategy bool
	kind               string
}

func (i *Individual) CalculateFitness() (float32, error) {
	return -1, nil
}

type Antagonist Individual
type Protagonist Individual
