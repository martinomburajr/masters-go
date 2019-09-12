package evolution

// JudgementDay represents a moment where all individuals have completed their epoch phase and are waiting a decision
// onto who proceeds to the next generation. Judgement day is a compound function or abstraction that includes the
// following processes.
// 1. Parent Selection
// 2. Reproduction (via Crossover)
// 3. Mutation (low probability)
// 4. Survivor Selection
// 5. Statistical Output
// 6. FinalPopulation configuration (incrementing age, clearing fitness values for old worthy individuals)
func JudgementDay(incomingPopulation []Individual, opts EvolutionParams) ([]Individual, error) {
	// Parent Selection
		// Tournament Selection

	// Reproduction
		// Crossover

	// Reproduction
		// Mutation

	// Survivor Selection

	// Statistical Output

	// Anointing Final Population and Return

	return nil, nil
}