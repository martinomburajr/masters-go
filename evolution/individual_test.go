package evolution

import (
	"reflect"
	"testing"
)

func TestGenerateRandomStrategy(t *testing.T) {
	type args struct {
		number              int
		strategyLength      int
		availableStrategies []Strategy
	}
	tests := []struct {
		name string
		args args
		want []Strategy
	}{
		{"small number", args{0, 0, nil}, []Strategy{}},
		{"small strategy length", args{4, 0, nil}, []Strategy{}},
		{"ok", args{4, 12, []Strategy{StrategyMutateNode, StrategyAddSubTree}}, []Strategy{StrategyMutateNode,
			StrategyAddSubTree, StrategyMutateNode, StrategyAddSubTree}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []Strategy
			if got = GenerateRandomStrategy(tt.args.number, tt.args.strategyLength, tt.args.availableStrategies); len(got) != len(tt.want) {
				t.Errorf("GenerateRandomStrategy() = %v, isEqual %v", got, tt.want)
			}
		})
	}
}

//func TestGenerateRandomIndividuals(t *testing.T) {
//	type depth struct {
//		number                int
//		idTemplate            string
//		kind                  int
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
//		{"empty kind", depth{10, "sometemplate", IndividualAntagonist, 10, 10, nil}, nil, true},
//		{"small strategy len", depth{10, "sometemplate", IndividualAntagonist, 0, 10, nil}, nil, true},
//		{"small max number of strategies", depth{10, "sometemplate", IndividualAntagonist, 10, 0, nil}, nil, true},
//		{"nil available strategies", depth{10, "sometemplate", IndividualAntagonist, 10, 10, nil}, nil, true},
//		{"empty strategies", depth{10, "sometemplate", IndividualAntagonist, 10, 10, []Strategy{}}, nil, true},
//		{"empty id-template", depth{10, "", IndividualAntagonist, 10, 10, []Strategy{StrategyAddSubTree}}, nil, true},
//		{"ok", depth{2, "bugs", IndividualAntagonist, 10, 10, []Strategy{StrategyAddSubTree}},
//			[]Individual{IndividualProg0Kind1, IndividualProg0Kind1}, false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			var got []Individual
//			var err error
//			got, err = GenerateRandomIndividuals(tt.depth.number, tt.depth.idTemplate, tt.depth.kind, tt.depth.strategyLength, tt.depth.maxNumberOfStrategies, tt.depth.availableStrategies)
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

func TestCrossover(t *testing.T) {
	type args struct {
		individual  Individual
		individual2 Individual
		params      EvolutionParams
	}
	tests := []struct {
		name    string
		args    args
		want    Individual
		want1   Individual
		wantErr bool
	}{
		//{"individual id = empty", args{individual:Individual{}, individual2:Individual{}, params:EvolutionParams{}},
		//	Individual{}, Individual{}, true},
		//{"individual strategy nil", args{individual:Individual{id: GenerateProgramID(0)}, individual2:Individual{},
		//	params:EvolutionParams{}},
		//	Individual{}, Individual{}, true},
		//{"individual strategy empty", args{individual:Individual{id: GenerateProgramID(0), strategy: []Strategy{}},
		//	individual2:Individual{},
		//	params:EvolutionParams{}},
		//	Individual{}, Individual{}, true},
		//{"individual has not calculated fitness", args{individual:Individual{id: GenerateProgramID(0),
		//	strategy: []Strategy{StrategyAddSubTree}}, individual2:Individual{},
		//	params:EvolutionParams{}},
		//	Individual{}, Individual{}, true},
		//{"individual has not applied strategy", args{individual:Individual{id: GenerateProgramID(0),
		//	strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true}, individual2:Individual{},
		//	params:EvolutionParams{}},
		//	Individual{}, Individual{}, true},
		//{"individual has not applied strategy", args{individual:Individual{id: GenerateProgramID(0),
		//	strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true}, individual2:Individual{},
		//	params:EvolutionParams{}},
		//	Individual{}, Individual{}, true},
		//{"individual nil PRogram", args{individual:Individual{id: GenerateProgramID(0),
		//	strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true, Program:nil},
		//	individual2:Individual{},
		//	params:EvolutionParams{}},
		//	Individual{}, Individual{}, true},
		//{"individual nil Tree", args{individual:Individual{id: GenerateProgramID(0),
		//	strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true,
		//	Program:&ProgNil},
		//	individual2:Individual{},
		//	params:EvolutionParams{}},
		//	Individual{}, Individual{}, true},
		//
		//{"individual2 id = empty", args{individual:Individual{id: GenerateProgramID(0),
		//	strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true,
		//	Program:&Prog0}, individual2:Individual{},
		//	params:EvolutionParams{}},
		//	Individual{}, Individual{}, true},
		//{"individual2 starategy = nil", args{individual:Individual{id: GenerateProgramID(0),
		//	strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true,
		//	Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0)},
		//	params:EvolutionParams{}},
		//	Individual{}, Individual{}, true},
		//{"individual2 strategy=empty", args{individual:Individual{id: GenerateProgramID(0),
		//	strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true,
		//	Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{}},
		//	params:EvolutionParams{}},
		//	Individual{}, Individual{}, true},
		//{"individual2 strategy=empty", args{individual:Individual{id: GenerateProgramID(0),
		//	strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true,
		//	Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree}},
		//	params:EvolutionParams{}},
		//	Individual{}, Individual{}, true},
		//{"individual2 has not calculated fitness", args{individual:Individual{id: GenerateProgramID(0),
		//	strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true,
		//	Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree}},
		//	params:EvolutionParams{}},
		//	Individual{}, Individual{}, true},
		//{"individual2 has not applied strategy", args{individual:Individual{id: GenerateProgramID(0),
		//	strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true,
		//	Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree},
		//		hasCalculatedFitness:true},
		//	params:EvolutionParams{}},
		//	Individual{}, Individual{}, true},
		//{"individual2 has not applied strategy", args{individual:Individual{id: GenerateProgramID(0),
		//	strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true,
		//	Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree},
		//		hasCalculatedFitness:true},
		//	params:EvolutionParams{}},
		//	Individual{}, Individual{}, true},
		//{"individual2 has not applied strategy", args{individual:Individual{id: GenerateProgramID(0),
		//	strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true,
		//	Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree},
		//		hasCalculatedFitness:true},
		//	params:EvolutionParams{}},
		//	Individual{}, Individual{}, true},
		//{"individual2 nil Program", args{individual:Individual{id: GenerateProgramID(0),
		//	strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true,
		//	Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree},
		//		hasCalculatedFitness:true, Program:&ProgNil},
		//	params:EvolutionParams{}},
		//	Individual{}, Individual{}, true},
		//{"individual2 nil Tree", args{individual:Individual{id: GenerateProgramID(0),
		//	strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true,
		//	Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree},
		//		hasCalculatedFitness:true, Program:&Prog0},
		//	params:EvolutionParams{}},
		//	Individual{}, Individual{}, true},
		//{"params strategy length limit", args{individual:Individual{id: GenerateProgramID(0),
		//	strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true,
		//	Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree},
		//		hasCalculatedFitness:true, hasAppliedStrategy:true, Program:&Prog0},
		//	params:EvolutionParams{StrategyLengthLimit: 0}},
		//	Individual{}, Individual{}, true},

		//	WORKING
		//{"params strategy length limit20", args{
		//	individual:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true,
		//		Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyDeleteSubTree}, hasCalculatedFitness:true,
		//		hasAppliedStrategy:true, Program:&Prog0},
		//	params:EvolutionParams{StrategyLengthLimit: 20}},
		//	Individual{}, Individual{}, false},
		//{"params strategy length crossover% = 1", args{
		//	individual:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true,
		//		Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyDeleteSubTree}, hasCalculatedFitness:true,
		//		hasAppliedStrategy:true, Program:&Prog0},
		//	params:EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage:1}},
		//	Individual{}, Individual{}, false},

		//{"params strategy length crossover% = 0.5", args{
		//	individual:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true,
		//		Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyDeleteSubTree}, hasCalculatedFitness:true,
		//		hasAppliedStrategy:true, Program:&Prog0},
		//	params:EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage:0.5,
		//		MaintainCrossoverGeneTransferEquality:true}},
		//	Individual{}, Individual{}, false},
		//
		//{"params strategy length crossover% = 0.5", args{
		//	individual:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true,
		//		Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyDeleteSubTree, StrategyMutateNode},
		//		hasCalculatedFitness:true,
		//		hasAppliedStrategy:true, Program:&Prog0},
		//	params:EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage:0.5,
		//		MaintainCrossoverGeneTransferEquality:true}},
		//	Individual{}, Individual{}, false},
		//
		//{"params strategy length crossover% = 0.5", args{
		//	individual:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree}, hasCalculatedFitness: true, hasAppliedStrategy:true,
		//		Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyDeleteSubTree,
		//		StrategyMutateNode, StrategyDeleteMalicious, StrategySwapSubTree},
		//		hasCalculatedFitness:true,
		//		hasAppliedStrategy:true, Program:&Prog0},
		//	params:EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage:0.5,
		//		MaintainCrossoverGeneTransferEquality:true}},
		//	Individual{}, Individual{}, false},
		//
		//
		//{"params strategy length crossover% = 0.5", args{
		//	individual:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree, StrategySoftDeleteSubTree},
		//		hasCalculatedFitness: true, hasAppliedStrategy:true,
		//		Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyDeleteSubTree,
		//		StrategyMutateNode, StrategyDeleteMalicious, StrategySwapSubTree},
		//		hasCalculatedFitness:true,
		//		hasAppliedStrategy:true, Program:&Prog0},
		//	params:EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage:0.5,
		//		MaintainCrossoverGeneTransferEquality:true}},
		//	Individual{}, Individual{}, false},

		//{"params strategy length crossover% = 0.5 | No Maintain Crossover Equality", args{
		//	individual:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree, StrategySoftDeleteSubTree},
		//		hasCalculatedFitness: true, hasAppliedStrategy:true,
		//		Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyDeleteSubTree,
		//		StrategyMutateNode, StrategyDeleteMalicious, StrategySwapSubTree},
		//		hasCalculatedFitness:true,
		//		hasAppliedStrategy:true, Program:&Prog0},
		//	params:EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage:0.5,
		//		MaintainCrossoverGeneTransferEquality:false}},
		//	Individual{}, Individual{}, false},
		//
		//{"params strategy length crossover% = 0.5 | No Maintain Crossover Equality", args{
		//	individual:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree,
		//		StrategySoftDeleteSubTree, StrategyDeleteMalicious, StrategySwapSubTree},
		//		hasCalculatedFitness: true, hasAppliedStrategy:true,
		//		Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyDeleteSubTree,
		//		StrategyMutateNode},
		//		hasCalculatedFitness:true,
		//		hasAppliedStrategy:true, Program:&Prog0},
		//	params:EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage:0.5,
		//		MaintainCrossoverGeneTransferEquality:true}},
		//	Individual{}, Individual{}, false},
		//{"params strategy length crossover% = 0.5 | No Maintain Crossover Equality", args{
		//	individual:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree,
		//		StrategySoftDeleteSubTree, StrategyDeleteMalicious, StrategySwapSubTree},
		//		hasCalculatedFitness: true, hasAppliedStrategy:true,
		//		Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyDeleteSubTree},
		//		hasCalculatedFitness:true,
		//		hasAppliedStrategy:true, Program:&Prog0},
		//	params:EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage:0.5,
		//		MaintainCrossoverGeneTransferEquality:true}},
		//	Individual{}, Individual{}, false},
		//{"params strategy length crossover% = 0.5 | No Maintain Crossover Equality", args{
		//	individual:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree,
		//		StrategySoftDeleteSubTree, StrategyDeleteMalicious, StrategySwapSubTree},
		//		hasCalculatedFitness: true, hasAppliedStrategy:true,
		//		Program:&Prog0},
		//	individual2:Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyDeleteSubTree,
		//		StrategyMutateNode},
		//		hasCalculatedFitness:true,
		//		hasAppliedStrategy:true, Program:&Prog0},
		//	params:EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage:0.5,
		//		MaintainCrossoverGeneTransferEquality:false}},
		//	Individual{}, Individual{}, false},
		{"params strategy length crossover% = 0.5 | No Maintain Crossover Equality", args{
			individual: Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyAddSubTree,
				StrategySoftDeleteSubTree, StrategyDeleteMalicious, StrategySwapSubTree},
				hasCalculatedFitness: true, hasAppliedStrategy: true,
				Program: &Prog0},
			individual2: Individual{id: GenerateProgramID(0), strategy: []Strategy{StrategyDeleteSubTree,
				StrategyMutateNode, StrategyDeleteMalicious, StrategySwapSubTree},
				hasCalculatedFitness: true,
				hasAppliedStrategy:   true, Program: &Prog0},
			params: EvolutionParams{StrategyLengthLimit: 20, CrossoverPercentage: 0.5,
				MaintainCrossoverGeneTransferEquality: false}},
			Individual{}, Individual{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := Crossover(tt.args.individual, tt.args.individual2, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Crossover() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got.Program.Fitness()
			got1.Program.Fitness()
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("Crossover() got = %v, want %v", got, tt.want)
			//}
			//if !reflect.DeepEqual(got1, tt.want1) {
			//	t.Errorf("Crossover() got1 = %v, want %v", got1, tt.want1)
			//}
		})
	}
}

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
		{"swapLength < 1", args{a: []Strategy{StrategySwapSubTree}, b: []Strategy{StrategyMutateSubTree}, swapLength: -1},
			[]Strategy{StrategySwapSubTree}, []Strategy{StrategyMutateSubTree}},
		{"startIndex < 1", args{a: []Strategy{StrategySwapSubTree}, b: []Strategy{StrategyMutateSubTree}, swapLength: 1,
			startIndex: -1}, []Strategy{StrategyMutateSubTree}, []Strategy{StrategySwapSubTree}},
		{"a1 : b1 : sl=0 : si=0", args{
			a:          []Strategy{StrategySwapSubTree},
			b:          []Strategy{StrategySoftDeleteSubTree},
			swapLength: 0,
			startIndex: 0},
			[]Strategy{StrategySwapSubTree},
			[]Strategy{StrategySoftDeleteSubTree},
		},
		{"a1 : b1 : sl=2 : si=0", args{
			a:          []Strategy{StrategySwapSubTree},
			b:          []Strategy{StrategySoftDeleteSubTree},
			swapLength: 2,
			startIndex: 0},
			[]Strategy{StrategySoftDeleteSubTree},
			[]Strategy{StrategySwapSubTree},
		},
		{"a2 : b1 : sl=0 : si=0", args{
			a:          []Strategy{StrategySwapSubTree, StrategyDeleteMalicious},
			b:          []Strategy{StrategySoftDeleteSubTree},
			swapLength: 0,
			startIndex: 0},
			[]Strategy{StrategySwapSubTree, StrategyDeleteMalicious},
			[]Strategy{StrategySoftDeleteSubTree},
		},
		{"a2 : b1 : sl=1 : si=0", args{
			a:          []Strategy{StrategySwapSubTree, StrategyDeleteMalicious},
			b:          []Strategy{StrategySoftDeleteSubTree},
			swapLength: 1,
			startIndex: 0},
			[]Strategy{StrategySoftDeleteSubTree, StrategyDeleteMalicious},
			[]Strategy{StrategySwapSubTree},
		},
		{"a2 : b2 : sl=1 : si=0", args{
			a:          []Strategy{StrategySwapSubTree, StrategyDeleteMalicious},
			b:          []Strategy{StrategySoftDeleteSubTree, StrategyAddSubTree},
			swapLength: 1,
			startIndex: 0},
			[]Strategy{StrategySoftDeleteSubTree, StrategyDeleteMalicious},
			[]Strategy{StrategySwapSubTree, StrategyAddSubTree},
		},
		{"a1 : b2 : sl=1 : si=0", args{
			a:          []Strategy{StrategySwapSubTree},
			b:          []Strategy{StrategySoftDeleteSubTree, StrategyAddSubTree},
			swapLength: 1,
			startIndex: 0},
			[]Strategy{StrategySoftDeleteSubTree},
			[]Strategy{StrategySwapSubTree, StrategyAddSubTree},
		},
		{"a1 : b2 : sl=2 : si=1", args{
			a:          []Strategy{StrategySwapSubTree},
			b:          []Strategy{StrategySoftDeleteSubTree, StrategyAddSubTree},
			swapLength: 2,
			startIndex: 1},
			[]Strategy{StrategySoftDeleteSubTree},
			[]Strategy{StrategySwapSubTree, StrategyAddSubTree},
		},
		{"a2 : b2 : sl=1 : si=1", args{
			a:          []Strategy{StrategySwapSubTree, StrategyDeleteMalicious},
			b:          []Strategy{StrategySoftDeleteSubTree, StrategyAddSubTree},
			swapLength: 1,
			startIndex: 1},
			[]Strategy{StrategySwapSubTree, StrategyAddSubTree},
			[]Strategy{StrategySoftDeleteSubTree, StrategyDeleteMalicious},
		},
		{"a2 : b2 : sl=1 : si=5", args{
			a:          []Strategy{StrategySwapSubTree, StrategyDeleteMalicious},
			b:          []Strategy{StrategySoftDeleteSubTree, StrategyAddSubTree},
			swapLength: 1,
			startIndex: 5},
			[]Strategy{StrategySoftDeleteSubTree, StrategyDeleteMalicious},
			[]Strategy{StrategySwapSubTree, StrategyAddSubTree},
		},
		{"a3 : b2 : sl=2 : si=5", args{
			a:          []Strategy{StrategySwapSubTree, StrategyDeleteMalicious, StrategyMutateSubTree},
			b:          []Strategy{StrategySoftDeleteSubTree, StrategyAddSubTree},
			swapLength: 2,
			startIndex: 5},
			[]Strategy{StrategySoftDeleteSubTree, StrategyAddSubTree, StrategyMutateSubTree},
			[]Strategy{StrategySwapSubTree, StrategyDeleteMalicious},
		},
		{"a3 : b2 : sl=2 : si=1", args{
			a:          []Strategy{StrategySwapSubTree, StrategyDeleteMalicious, StrategyMutateSubTree},
			b:          []Strategy{StrategySoftDeleteSubTree, StrategyAddSubTree},
			swapLength: 2,
			startIndex: 1},
			[]Strategy{StrategySoftDeleteSubTree, StrategyAddSubTree, StrategyMutateSubTree},
			[]Strategy{StrategySwapSubTree, StrategyDeleteMalicious},
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
		//	a:          []Strategy{StrategySwapSubTree},
		//	b:          []Strategy{StrategySoftDeleteSubTree},
		//	swapLength: 0,
		//	startIndex: 0},
		//	[]Strategy{StrategySwapSubTree},
		//	[]Strategy{StrategySoftDeleteSubTree},
		//},
		//{"a1 : b1 : sl=1 : si=0", args{
		//	a:          []Strategy{StrategySwapSubTree},
		//	b:          []Strategy{StrategySoftDeleteSubTree},
		//	swapLength: 1,
		//	startIndex: 0},
		//	[]Strategy{StrategySoftDeleteSubTree},
		//	[]Strategy{StrategySwapSubTree},
		//},
		//{"a1 : b1 : sl=1 : si=1", args{
		//	a:          []Strategy{StrategySwapSubTree},
		//	b:          []Strategy{StrategySoftDeleteSubTree},
		//	swapLength: 1,
		//	startIndex: 1},
		//	[]Strategy{StrategySoftDeleteSubTree},
		//	[]Strategy{StrategySwapSubTree},
		//},
		{"a2 : b1 : sl=1 : si=0", args{
			a:          []Strategy{StrategySwapSubTree, StrategyDeleteMalicious},
			b:          []Strategy{StrategySoftDeleteSubTree},
			swapLength: 1,
			startIndex: 0},
			[]Strategy{StrategySoftDeleteSubTree, StrategyDeleteMalicious},
			[]Strategy{StrategySwapSubTree},
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
