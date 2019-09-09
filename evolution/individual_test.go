package evolution

import (
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
		{"ok", args{4, 12, []Strategy{{Kind: StrategyMutateNode}, {Kind: StrategyAddSubTree}}}, []Strategy{{Kind: StrategyMutateNode},
			{Kind: StrategyAddSubTree}, {Kind: StrategyMutateNode}, {Kind: StrategyAddSubTree}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []Strategy
			if got = GenerateRandomStrategy(tt.args.number, tt.args.strategyLength, tt.args.availableStrategies); len(got) != len(tt.want) {
				t.Errorf("GenerateRandomStrategy() = %v, wantAntagonist %v", got, tt.want)
			}
		})
	}
}

func TestGenerateRandomIndividuals(t *testing.T) {
	type args struct {
		number                int
		idTemplate            string
		kind                  int
		strategyLength        int
		maxNumberOfStrategies int
		availableStrategies   []Strategy
	}
	tests := []struct {
		name    string
		args    args
		want    []*Individual
		wantErr bool
	}{
		{"small number", args{0, "", IndividualAntagonist, 0, 0, nil}, nil, true},
		{"empty kind", args{10, "sometemplate", IndividualAntagonist, 10, 10, nil}, nil, true},
		{"small strategy len", args{10, "sometemplate", IndividualAntagonist, 0, 10, nil}, nil, true},
		{"small max number of strategies", args{10, "sometemplate", IndividualAntagonist, 10, 0, nil}, nil, true},
		{"nil available strategies", args{10, "sometemplate", IndividualAntagonist, 10, 10, nil}, nil, true},
		{"empty strategies", args{10, "sometemplate", IndividualAntagonist, 10, 10, []Strategy{}}, nil, true},
		{"empty id-template", args{10, "", IndividualAntagonist, 10, 10, []Strategy{{Kind: StrategyAddSubTree}}}, nil, true},
		{"ok", args{2, "bugs", IndividualAntagonist, 10, 10, []Strategy{{Kind: StrategyAddSubTree}}}, []Individual{{}, {}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []Individual
			var err error
			got, err = GenerateRandomIndividuals(tt.args.number, tt.args.idTemplate, tt.args.kind, tt.args.strategyLength, tt.args.maxNumberOfStrategies, tt.args.availableStrategies)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRandomIndividuals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if len(got) != len(tt.want) {
					t.Errorf("GenerateRandomIndividuals() = %v, wantAntagonist %v", got, tt.want)
				}
			}
		})
	}
}
