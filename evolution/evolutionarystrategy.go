package evolution

var (
	EvolutionaryStrategy = evolutionaryStrategy{}
)

type evolutionaryStrategy struct {
	Tournament   bool
	Rank         bool
	Minimization bool
}
