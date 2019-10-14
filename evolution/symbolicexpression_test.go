package evolution

//func TestParseString(t *testing.T) {
//	type args struct {
//		expression        string
//		validNonTerminals []string
//	}
//	tests := []struct {
//		name             string
//		args             args
//		wantTerminals    []SymbolicExpression
//		wantNonTerminals []SymbolicExpression
//		wantExpression   []SymbolicExpression
//		wantErr          bool
//	}{
//		{"nil", args{"", nil}, nil, nil, nil, true},
//		{"exp", args{"x", nil}, nil, nil, nil, true},
//		{"exp | valid NT | x", args{"x", []string{"*"}},
//			[]SymbolicExpression{{value: "x", kind: 0, arity: 0}},
//			[]SymbolicExpression{},
//			[]SymbolicExpression{{value: "x", kind: 0, arity: 0}},
//			false},
//		{"exp | valid NT | x + ", args{"x +   ", []string{"+"}},
//			[]SymbolicExpression{{value: "x", kind: 0, arity: 0}, {value: "1", kind: 0, arity: 0}},
//			[]SymbolicExpression{{value: "+", kind: 1, arity: 2}},
//			[]SymbolicExpression{{value: "x", kind: 0, arity: 0}, {value: "+", kind: 1, arity: 2}, {value: "1", kind: 0, arity: 0}},
//			true},
//		{"exp | valid NT | x + 1 ", args{"x + 1", []string{"+"}},
//			[]SymbolicExpression{{value: "x", kind: 0, arity: 0}, {value: "1", kind: 0, arity: 0}},
//			[]SymbolicExpression{{value: "+", kind: 1, arity: 2}},
//			[]SymbolicExpression{{value: "x", kind: 0, arity: 0}, {value: "+", kind: 1, arity: 2}, {value: "1", kind: 0, arity: 0}},
//			false},
//		{"exp | missing valid NT | x + 1 ", args{"x + 1", []string{"-"}},
//			nil, nil, nil,
//			true},
//		{"exp | valid NT | x + 1 - 2", args{"x + 1 - 2", []string{"+", "-"}},
//			[]SymbolicExpression{{value: "x", kind: 0, arity: 0}, {value: "1", kind: 0, arity: 0}, {value: "2", kind: 0, arity: 0}},
//			[]SymbolicExpression{{value: "+", kind: 1, arity: 2}, {value: "-", kind: 1, arity: 2}},
//			[]SymbolicExpression{{value: "x", kind: 0, arity: 0}, {value: "+", kind: 1, arity: 2}, {value: "1", kind: 0, arity: 0},
//				{value: "-", kind: 1, arity: 2}, {value: "2", kind: 0, arity: 0}},
//			false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotTerminals, gotNonTerminals, gotExpression, err := ParseString(tt.args.expression,
//				tt.args.validNonTerminals)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("ParseString() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if err == nil {
//				if !reflect.DeepEqual(gotTerminals, tt.wantTerminals) {
//					t.Errorf("ParseString() gotTerminals = %v, want %v", gotTerminals, tt.wantTerminals)
//				}
//				if !reflect.DeepEqual(gotNonTerminals, tt.wantNonTerminals) {
//					t.Errorf("ParseString() gotNonTerminals = %v, want %v", gotNonTerminals, tt.wantNonTerminals)
//				}
//				if !reflect.DeepEqual(gotExpression, tt.wantExpression) {
//					t.Errorf("ParseString() gotExpression = %v, want %v", gotNonTerminals, tt.wantNonTerminals)
//				}
//			}
//		})
//	}
//}
