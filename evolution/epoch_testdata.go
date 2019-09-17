package evolution

var EpochNil = Epoch{}

var Epoch0 = Epoch{
	id:                               "epoch0",
	program:                          Prog0,
	generation:                       &Generation{},
	protagonist:                      &IndividualProg0Kind1,
	antagonist:                       &IndividualProg0Kind0,
	probabilityOfNonTerminalMutation: 0.1,
	probabilityOfMutation:            0.2,
	hasProtagonistApplied:            false,
	hasAntagonistApplied:             false,
	nonTerminalSet:                   SymbolicExpressionSet{Mult},
	terminalSet:                      SymbolicExpressionSet{X1},
	protagonistBegins:                true,
	isComplete:                       false,
}

var Epoch1 = Epoch{
	id:                               "epoch1",
	program:                          ProgTreeT_NT_T_0,
	generation:                       &Generation{},
	protagonist:                      &IndividualProg0Kind1,
	antagonist:                       &IndividualProg0Kind0,
	probabilityOfNonTerminalMutation: 0.1,
	probabilityOfMutation:            0.2,
	hasProtagonistApplied:            false,
	hasAntagonistApplied:             false,
	nonTerminalSet:                   SymbolicExpressionSet{Mult},
	terminalSet:                      SymbolicExpressionSet{X1},
	protagonistBegins:                true,
	isComplete:                       false,
}