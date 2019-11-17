package evolution

import (
	"reflect"
	"testing"
)

func TestGenerateRandomStrategy(t *testing.T) {
	type args struct {
		number              int
		availableStrategies []Strategy
	}
	tests := []struct {
		name string
		args args
		want []Strategy
	}{
		{"small number", args{0, nil}, []Strategy{}},
		{"small Strategy length", args{4, nil}, []Strategy{}},
		{"ok", args{4, []Strategy{StrategyMutateTerminal, StrategyAddRandomSubTree}}, []Strategy{StrategyMutateTerminal,
			StrategyAddRandomSubTree, StrategyMutateTerminal, StrategyAddRandomSubTree}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []Strategy
			if got = GenerateRandomStrategy(tt.args.number, tt.args.availableStrategies); len(got) != len(tt.want) {
				t.Errorf("GenerateRandomStrategy() = %v, isEqual %v", got, tt.want)
			}
		})
	}
}

//func TestGenerateRandomIndividuals(t *testing.T) {
//	type depth struct {
//		number                int
//		idTemplate            string
//		Kind                  int
//		strategyLength        int
//		maxNumberOfStrategies int
//		availableStrategies   []Strategy
//	}
//	tests := []struct {
//		name    string
//		depth    depth
//		want    []*Individual
//		wantErr bool
//	}{
//		{"small number", depth{0, "", IndividualAntagonist, 0, 0, nil}, nil, true},
//		{"empty Kind", depth{10, "sometemplate", IndividualAntagonist, 10, 10, nil}, nil, true},
//		{"small Strategy len", depth{10, "sometemplate", IndividualAntagonist, 0, 10, nil}, nil, true},
//		{"small max number of strategies", depth{10, "sometemplate", IndividualAntagonist, 10, 0, nil}, nil, true},
//		{"nil available strategies", depth{10, "sometemplate", IndividualAntagonist, 10, 10, nil}, nil, true},
//		{"empty strategies", depth{10, "sometemplate", IndividualAntagonist, 10, 10, []Strategy{}}, nil, true},
//		{"empty Id-template", depth{10, "", IndividualAntagonist, 10, 10, []Strategy{StrategyAddRandomSubTree}}, nil, true},
//		{"ok", depth{2, "bugs", IndividualAntagonist, 10, 10, []Strategy{StrategyAddRandomSubTree}},
//			[]Individual{IndividualProg0Kind1, IndividualProg0Kind1}, false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			var got []Individual
//			var err error
//			got, err = GenerateRandomIndividuals(tt.depth.number, tt.depth.idTemplate, tt.depth.Kind, tt.depth.strategyLength, tt.depth.maxNumberOfStrategies, tt.depth.availableStrategies)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GenerateRandomIndividuals() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if err == nil {
//				if len(got) != len(tt.want) {
//					t.Errorf("GenerateRandomIndividuals() = %v, isEqual %v", got, tt.want)
//				}
//			}
//		})
//	}
//}

//func TestCrossover(t *testing.T) {
//	type args struct {
//		individual  Individual
//		individual2 Individual
//		params      EvolutionParams
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    Individual
//		want1   Individual
//		wantErr bool
//	}{
//		//{"individual Id = empty", args{individual:Individual{}, individual2:Individual{}, params:EvolutionParams{}},
//		//	Individual{}, Individual{}, true},
//		//{"individual Strategy nil", args{individual:Individual{Id: GenerateProgramID(0)}, individual2:Individual{},
//		//	params:EvolutionParams{}},
//		//	Individual{}, Individual{}, true},
//		//{"individual Strategy empty", args{individual:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{}},
//		//	individual2:Individual{},
//		//	params:EvolutionParams{}},
//		//	Individual{}, Individual{}, true},
//		//{"individual has not calculated Fitness", args{individual:Individual{Id: GenerateProgramID(0),
//		//	Strategy: []Strategy{StrategyAddRandomSubTree}}, individual2:Individual{},
//		//	params:EvolutionParams{}},
//		//	Individual{}, Individual{}, true},
//		//{"individual has not applied Strategy", args{individual:Individual{Id: GenerateProgramID(0),
//		//	Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true}, individual2:Individual{},
//		//	params:EvolutionParams{}},
//		//	Individual{}, Individual{}, true},
//		//{"individual has not applied Strategy", args{individual:Individual{Id: GenerateProgramID(0),
//		//	Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true}, individual2:Individual{},
//		//	params:EvolutionParams{}},
//		//	Individual{}, Individual{}, true},
//		//{"individual nil PRogram", args{individual:Individual{Id: GenerateProgramID(0),
//		//	Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true, Program:nil},
//		//	individual2:Individual{},
//		//	params:EvolutionParams{}},
//		//	Individual{}, Individual{}, true},
//		//{"individual nil Tree", args{individual:Individual{Id: GenerateProgramID(0),
//		//	Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//	Program:&ProgNil},
//		//	individual2:Individual{},
//		//	params:EvolutionParams{}},
//		//	Individual{}, Individual{}, true},
//		//
//		//{"individual2 Id = empty", args{individual:Individual{Id: GenerateProgramID(0),
//		//	Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//	Program:&ProgX}, individual2:Individual{},
//		//	params:EvolutionParams{}},
//		//	Individual{}, Individual{}, true},
//		//{"individual2 starategy = nil", args{individual:Individual{Id: GenerateProgramID(0),
//		//	Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//	Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0)},
//		//	params:EvolutionParams{}},
//		//	Individual{}, Individual{}, true},
//		//{"individual2 Strategy=empty", args{individual:Individual{Id: GenerateProgramID(0),
//		//	Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//	Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{}},
//		//	params:EvolutionParams{}},
//		//	Individual{}, Individual{}, true},
//		//{"individual2 Strategy=empty", args{individual:Individual{Id: GenerateProgramID(0),
//		//	Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//	Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree}},
//		//	params:EvolutionParams{}},
//		//	Individual{}, Individual{}, true},
//		//{"individual2 has not calculated Fitness", args{individual:Individual{Id: GenerateProgramID(0),
//		//	Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//	Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree}},
//		//	params:EvolutionParams{}},
//		//	Individual{}, Individual{}, true},
//		//{"individual2 has not applied Strategy", args{individual:Individual{Id: GenerateProgramID(0),
//		//	Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//	Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree},
//		//		HasCalculatedFitness:true},
//		//	params:EvolutionParams{}},
//		//	Individual{}, Individual{}, true},
//		//{"individual2 has not applied Strategy", args{individual:Individual{Id: GenerateProgramID(0),
//		//	Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//	Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree},
//		//		HasCalculatedFitness:true},
//		//	params:EvolutionParams{}},
//		//	Individual{}, Individual{}, true},
//		//{"individual2 has not applied Strategy", args{individual:Individual{Id: GenerateProgramID(0),
//		//	Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//	Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree},
//		//		HasCalculatedFitness:true},
//		//	params:EvolutionParams{}},
//		//	Individual{}, Individual{}, true},
//		//{"individual2 nil Program", args{individual:Individual{Id: GenerateProgramID(0),
//		//	Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//	Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree},
//		//		HasCalculatedFitness:true, Program:&ProgNil},
//		//	params:EvolutionParams{}},
//		//	Individual{}, Individual{}, true},
//		//{"individual2 nil Tree", args{individual:Individual{Id: GenerateProgramID(0),
//		//	Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//	Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree},
//		//		HasCalculatedFitness:true, Program:&ProgX},
//		//	params:EvolutionParams{}},
//		//	Individual{}, Individual{}, true},
//		//{"params Strategy length limit", args{individual:Individual{Id: GenerateProgramID(0),
//		//	Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//	Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree},
//		//		HasCalculatedFitness:true, HasAppliedStrategy:true, Program:&ProgX},
//		//	params:EvolutionParams{StrategyLengthLimit: 0}},
//		//	Individual{}, Individual{}, true},
//
//		//	WORKING
//		//{"params Strategy length limit20", args{
//		//	individual:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//		Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyDeleteSubTree}, HasCalculatedFitness:true,
//		//		HasAppliedStrategy:true, Program:&ProgX},
//		//	params:EvolutionParams{StrategyLengthLimit: 20}},
//		//	Individual{}, Individual{}, false},
//		//{"params Strategy length crossover% = 1", args{
//		//	individual:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//		Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyDeleteSubTree}, HasCalculatedFitness:true,
//		//		HasAppliedStrategy:true, Program:&ProgX},
//		//	params:EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage:1}},
//		//	Individual{}, Individual{}, false},
//
//		//{"params Strategy length crossover% = 0.5", args{
//		//	individual:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//		Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyDeleteSubTree}, HasCalculatedFitness:true,
//		//		HasAppliedStrategy:true, Program:&ProgX},
//		//	params:EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage:0.5,
//		//		MaintainCrossoverGeneTransferEquality:true}},
//		//	Individual{}, Individual{}, false},
//		//
//		//{"params Strategy length crossover% = 0.5", args{
//		//	individual:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//		Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyDeleteSubTree, StrategyMutateTerminal},
//		//		HasCalculatedFitness:true,
//		//		HasAppliedStrategy:true, Program:&ProgX},
//		//	params:EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage:0.5,
//		//		MaintainCrossoverGeneTransferEquality:true}},
//		//	Individual{}, Individual{}, false},
//		//
//		//{"params Strategy length crossover% = 0.5", args{
//		//	individual:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree}, HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//		Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyDeleteSubTree,
//		//		StrategyMutateTerminal, StrategyDeleteMalicious, StrategyReplaceBranch},
//		//		HasCalculatedFitness:true,
//		//		HasAppliedStrategy:true, Program:&ProgX},
//		//	params:EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage:0.5,
//		//		MaintainCrossoverGeneTransferEquality:true}},
//		//	Individual{}, Individual{}, false},
//		//
//		//
//		//{"params Strategy length crossover% = 0.5", args{
//		//	individual:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree, StrategyDeleteTerminal},
//		//		HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//		Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyDeleteSubTree,
//		//		StrategyMutateTerminal, StrategyDeleteMalicious, StrategyReplaceBranch},
//		//		HasCalculatedFitness:true,
//		//		HasAppliedStrategy:true, Program:&ProgX},
//		//	params:EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage:0.5,
//		//		MaintainCrossoverGeneTransferEquality:true}},
//		//	Individual{}, Individual{}, false},
//
//		//{"params Strategy length crossover% = 0.5 | No Maintain Crossover Equality", args{
//		//	individual:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree, StrategyDeleteTerminal},
//		//		HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//		Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyDeleteSubTree,
//		//		StrategyMutateTerminal, StrategyDeleteMalicious, StrategyReplaceBranch},
//		//		HasCalculatedFitness:true,
//		//		HasAppliedStrategy:true, Program:&ProgX},
//		//	params:EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage:0.5,
//		//		MaintainCrossoverGeneTransferEquality:false}},
//		//	Individual{}, Individual{}, false},
//		//
//		//{"params Strategy length crossover% = 0.5 | No Maintain Crossover Equality", args{
//		//	individual:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree,
//		//		StrategyDeleteTerminal, StrategyDeleteMalicious, StrategyReplaceBranch},
//		//		HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//		Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyDeleteSubTree,
//		//		StrategyMutateTerminal},
//		//		HasCalculatedFitness:true,
//		//		HasAppliedStrategy:true, Program:&ProgX},
//		//	params:EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage:0.5,
//		//		MaintainCrossoverGeneTransferEquality:true}},
//		//	Individual{}, Individual{}, false},
//		//{"params Strategy length crossover% = 0.5 | No Maintain Crossover Equality", args{
//		//	individual:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree,
//		//		StrategyDeleteTerminal, StrategyDeleteMalicious, StrategyReplaceBranch},
//		//		HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//		Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyDeleteSubTree},
//		//		HasCalculatedFitness:true,
//		//		HasAppliedStrategy:true, Program:&ProgX},
//		//	params:EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage:0.5,
//		//		MaintainCrossoverGeneTransferEquality:true}},
//		//	Individual{}, Individual{}, false},
//		//{"params Strategy length crossover% = 0.5 | No Maintain Crossover Equality", args{
//		//	individual:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree,
//		//		StrategyDeleteTerminal, StrategyDeleteMalicious, StrategyReplaceBranch},
//		//		HasCalculatedFitness: true, HasAppliedStrategy:true,
//		//		Program:&ProgX},
//		//	individual2:Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyDeleteSubTree,
//		//		StrategyMutateTerminal},
//		//		HasCalculatedFitness:true,
//		//		HasAppliedStrategy:true, Program:&ProgX},
//		//	params:EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage:0.5,
//		//		MaintainCrossoverGeneTransferEquality:false}},
//		//	Individual{}, Individual{}, false},
//		{"params Strategy length crossover% = 0.5 | No Maintain Crossover Equality", args{
//			individual: Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyAddRandomSubTree,
//				StrategyDeleteTerminal, StrategyDeleteMalicious, StrategyReplaceBranch},
//				HasCalculatedFitness: true, HasAppliedStrategy: true,
//				Program: &ProgX},
//			individual2: Individual{Id: GenerateProgramID(0), Strategy: []Strategy{StrategyDeleteSubTree,
//				StrategyMutateTerminal, StrategyDeleteMalicious, StrategyReplaceBranch},
//				HasCalculatedFitness: true,
//				HasAppliedStrategy:   true, Program: &ProgX},
//			params: EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage: 0.5,
//				MaintainCrossoverGeneTransferEquality: false}},
//			Individual{}, Individual{}, false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1, err := Crossover(tt.args.individual, tt.args.individual2, tt.args.params)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Crossover() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			got.Program.Fitness()
//			got1.Program.Fitness()
//			//if !reflect.DeepEqual(got, tt.want) {
//			//	t.Errorf("Crossover() got = %v, want %v", got, tt.want)
//			//}
//			//if !reflect.DeepEqual(got1, tt.want1) {
//			//	t.Errorf("Crossover() got1 = %v, want %v", got1, tt.want1)
//			//}
//		})
//	}
//}

func TestStrategySwapper(t *testing.T) {
	type args struct {
		a          []Strategy
		b          []Strategy
		swapLength int
		startIndex int
	}
	tests := []struct {
		name  string
		args  args
		want  []Strategy
		want1 []Strategy
	}{
		{"nil", args{a: nil, b: nil, swapLength: 0, startIndex: 0}, nil, nil},
		{"size-empty", args{a: []Strategy{}, b: []Strategy{}, swapLength: 0}, nil, nil},
		{"swapLength < 1", args{a: []Strategy{StrategyReplaceBranch}, b: []Strategy{StrategyMutateNonTerminal}, swapLength: -1},
			[]Strategy{StrategyReplaceBranch}, []Strategy{StrategyMutateNonTerminal}},
		{"startIndex < 1", args{a: []Strategy{StrategyReplaceBranch}, b: []Strategy{StrategyMutateNonTerminal}, swapLength: 1,
			startIndex: -1}, []Strategy{StrategyMutateNonTerminal}, []Strategy{StrategyReplaceBranch}},
		{"a1 : b1 : sl=0 : si=0", args{
			a:          []Strategy{StrategyReplaceBranch},
			b:          []Strategy{StrategyDeleteTerminal},
			swapLength: 0,
			startIndex: 0},
			[]Strategy{StrategyReplaceBranch},
			[]Strategy{StrategyDeleteTerminal},
		},
		{"a1 : b1 : sl=2 : si=0", args{
			a:          []Strategy{StrategyReplaceBranch},
			b:          []Strategy{StrategyDeleteTerminal},
			swapLength: 2,
			startIndex: 0},
			[]Strategy{StrategyDeleteTerminal},
			[]Strategy{StrategyReplaceBranch},
		},
		{"a2 : b1 : sl=0 : si=0", args{
			a:          []Strategy{StrategyReplaceBranch, StrategyDeleteMalicious},
			b:          []Strategy{StrategyDeleteTerminal},
			swapLength: 0,
			startIndex: 0},
			[]Strategy{StrategyReplaceBranch, StrategyDeleteMalicious},
			[]Strategy{StrategyDeleteTerminal},
		},
		{"a2 : b1 : sl=1 : si=0", args{
			a:          []Strategy{StrategyReplaceBranch, StrategyDeleteMalicious},
			b:          []Strategy{StrategyDeleteTerminal},
			swapLength: 1,
			startIndex: 0},
			[]Strategy{StrategyDeleteTerminal, StrategyDeleteMalicious},
			[]Strategy{StrategyReplaceBranch},
		},
		{"a2 : b2 : sl=1 : si=0", args{
			a:          []Strategy{StrategyReplaceBranch, StrategyDeleteMalicious},
			b:          []Strategy{StrategyDeleteTerminal, StrategyAddRandomSubTree},
			swapLength: 1,
			startIndex: 0},
			[]Strategy{StrategyDeleteTerminal, StrategyDeleteMalicious},
			[]Strategy{StrategyReplaceBranch, StrategyAddRandomSubTree},
		},
		{"a1 : b2 : sl=1 : si=0", args{
			a:          []Strategy{StrategyReplaceBranch},
			b:          []Strategy{StrategyDeleteTerminal, StrategyAddRandomSubTree},
			swapLength: 1,
			startIndex: 0},
			[]Strategy{StrategyDeleteTerminal},
			[]Strategy{StrategyReplaceBranch, StrategyAddRandomSubTree},
		},
		{"a1 : b2 : sl=2 : si=1", args{
			a:          []Strategy{StrategyReplaceBranch},
			b:          []Strategy{StrategyDeleteTerminal, StrategyAddRandomSubTree},
			swapLength: 2,
			startIndex: 1},
			[]Strategy{StrategyDeleteTerminal},
			[]Strategy{StrategyReplaceBranch, StrategyAddRandomSubTree},
		},
		{"a2 : b2 : sl=1 : si=1", args{
			a:          []Strategy{StrategyReplaceBranch, StrategyDeleteMalicious},
			b:          []Strategy{StrategyDeleteTerminal, StrategyAddRandomSubTree},
			swapLength: 1,
			startIndex: 1},
			[]Strategy{StrategyReplaceBranch, StrategyAddRandomSubTree},
			[]Strategy{StrategyDeleteTerminal, StrategyDeleteMalicious},
		},
		{"a2 : b2 : sl=1 : si=5", args{
			a:          []Strategy{StrategyReplaceBranch, StrategyDeleteMalicious},
			b:          []Strategy{StrategyDeleteTerminal, StrategyAddRandomSubTree},
			swapLength: 1,
			startIndex: 5},
			[]Strategy{StrategyDeleteTerminal, StrategyDeleteMalicious},
			[]Strategy{StrategyReplaceBranch, StrategyAddRandomSubTree},
		},
		{"a3 : b2 : sl=2 : si=5", args{
			a:          []Strategy{StrategyReplaceBranch, StrategyDeleteMalicious, StrategyMutateNonTerminal},
			b:          []Strategy{StrategyDeleteTerminal, StrategyAddRandomSubTree},
			swapLength: 2,
			startIndex: 5},
			[]Strategy{StrategyDeleteTerminal, StrategyAddRandomSubTree, StrategyMutateNonTerminal},
			[]Strategy{StrategyReplaceBranch, StrategyDeleteMalicious},
		},
		{"a3 : b2 : sl=2 : si=1", args{
			a:          []Strategy{StrategyReplaceBranch, StrategyDeleteMalicious, StrategyMutateNonTerminal},
			b:          []Strategy{StrategyDeleteTerminal, StrategyAddRandomSubTree},
			swapLength: 2,
			startIndex: 1},
			[]Strategy{StrategyDeleteTerminal, StrategyAddRandomSubTree, StrategyMutateNonTerminal},
			[]Strategy{StrategyReplaceBranch, StrategyDeleteMalicious},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := StrategySwapper(tt.args.a, tt.args.b, tt.args.swapLength, tt.args.startIndex)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StrategySwapper() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("StrategySwapper() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestStrategySwapperIgnorant(t *testing.T) {
	type args struct {
		a          []Strategy
		b          []Strategy
		swapLength int
		startIndex int
	}
	tests := []struct {
		name  string
		args  args
		want  []Strategy
		want1 []Strategy
	}{
		//{"nil", args{a: nil, b: nil, swapLength: 0, startIndex: 0}, nil, nil},
		//{"size-empty", args{a: []Strategy{}, b: []Strategy{}, swapLength: 0}, nil, nil},
		//{"a1 : b1 : sl=0 : si=0", args{
		//	a:          []Strategy{StrategyReplaceBranch},
		//	b:          []Strategy{StrategyDeleteTerminal},
		//	swapLength: 0,
		//	startIndex: 0},
		//	[]Strategy{StrategyReplaceBranch},
		//	[]Strategy{StrategyDeleteTerminal},
		//},
		//{"a1 : b1 : sl=1 : si=0", args{
		//	a:          []Strategy{StrategyReplaceBranch},
		//	b:          []Strategy{StrategyDeleteTerminal},
		//	swapLength: 1,
		//	startIndex: 0},
		//	[]Strategy{StrategyDeleteTerminal},
		//	[]Strategy{StrategyReplaceBranch},
		//},
		//{"a1 : b1 : sl=1 : si=1", args{
		//	a:          []Strategy{StrategyReplaceBranch},
		//	b:          []Strategy{StrategyDeleteTerminal},
		//	swapLength: 1,
		//	startIndex: 1},
		//	[]Strategy{StrategyDeleteTerminal},
		//	[]Strategy{StrategyReplaceBranch},
		//},
		{"a2 : b1 : sl=1 : si=0", args{
			a:          []Strategy{StrategyReplaceBranch, StrategyDeleteMalicious},
			b:          []Strategy{StrategyDeleteTerminal},
			swapLength: 1,
			startIndex: 0},
			[]Strategy{StrategyDeleteTerminal, StrategyDeleteMalicious},
			[]Strategy{StrategyReplaceBranch},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := StrategySwapperIgnorant(tt.args.a, tt.args.b, tt.args.swapLength, tt.args.startIndex)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StrategySwapperIgnorant() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("StrategySwapperIgnorant() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestIndividual_Clone(t *testing.T) {
	prog1 := &Prog1
	tests := []struct {
		name    string
		fields  Individual
		want    Individual
		wantErr bool
	}{
		{"",
			Individual{Strategy:[]Strategy{"1","2","3"}, AverageFitness: 10},
			Individual{Strategy:[]Strategy{"1","2", "3"}, AverageFitness: 10},
			false},
		{"",
			Individual{Strategy:[]Strategy{"1","2","3"}, AverageFitness: 10, BirthGen: 3},
			Individual{Strategy:[]Strategy{"1","2", "3"}, AverageFitness: 10, BirthGen: 3},
			false},
		{"",
			Individual{Strategy:[]Strategy{"1","2","3"}, AverageFitness: 10, BirthGen: 3, Program:prog1},
			Individual{Strategy:[]Strategy{"1","2","3"}, AverageFitness: 10, BirthGen: 3, Program:prog1},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Clone()
			//fmt.Printf("got: %p, want: %p", got.Strategy, &tt.want.Strategy)
			if (err != nil) != tt.wantErr {
				t.Errorf("Individual.Clone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				containsSubTree, err := got.Program.T.ContainsSubTree(tt.want.Program.T)
				if treeEqual, _ := containsSubTree, err; !treeEqual {
					t.Errorf("Individual.Clone() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestIndividual_CloneWithTree(t *testing.T) {
	type fields struct {
		Id                       string
		Strategy                 []Strategy
		Fitness                  []float64
		HasAppliedStrategy       bool
		HasCalculatedFitness     bool
		FitnessCalculationMethod int
		Kind                     int
		Age                      int
		TotalFitness             float64
		Program                  *Program
	}
	type args struct {
		tree DualTree
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Individual
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Individual{
				Id:                   tt.fields.Id,
				Strategy:             tt.fields.Strategy,
				Fitness:              tt.fields.Fitness,
				HasAppliedStrategy:   tt.fields.HasAppliedStrategy,
				HasCalculatedFitness: tt.fields.HasCalculatedFitness,
				//FitnessCalculationMethod: tt.fields.FitnessCalculationMethod,
				Kind:           tt.fields.Kind,
				Age:            tt.fields.Age,
				AverageFitness: tt.fields.TotalFitness,
				Program:        tt.fields.Program,
			}
			if got := i.CloneWithTree(tt.args.tree); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Individual.CloneWithTree() = %v, want %v", got, tt.want)
			}
		})
	}
}
