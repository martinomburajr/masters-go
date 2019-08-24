package evolution

import "testing"

func BenchmarkDualTree_FromSymbolicExpressionSet_1(b *testing.B) {
	expressionSet1 := GenerateRandomSymbolicExpressionSet(1)
	tree1 := DualTree{}
	for i := 0; i < b.N; i++ {
		tree1.FromSymbolicExpressionSet(expressionSet1)
	}
}

func BenchmarkDualTree_FromSymbolicExpressionSet2_1(b *testing.B) {
	expressionSet1 := GenerateRandomSymbolicExpressionSet(1)
	tree1 := DualTree{}
	for i := 0; i < b.N; i++ {
		tree1.FromSymbolicExpressionSet2(expressionSet1)
	}
}

func BenchmarkDualTree_FromSymbolicExpressionSet_10(b *testing.B) {
	expressionSet1 := GenerateRandomSymbolicExpressionSet(10)
	tree1 := DualTree{}
	for i := 0; i < b.N; i++ {
		tree1.FromSymbolicExpressionSet(expressionSet1)
	}
}

func BenchmarkDualTree_FromSymbolicExpressionSet2_10(b *testing.B) {
	expressionSet1 := GenerateRandomSymbolicExpressionSet(10)
	tree1 := DualTree{}
	for i := 0; i < b.N; i++ {
		tree1.FromSymbolicExpressionSet2(expressionSet1)
	}
}

func BenchmarkDualTree_FromSymbolicExpressionSet_1000(b *testing.B) {
	expressionSet1 := GenerateRandomSymbolicExpressionSet(1000)
	tree1 := DualTree{}
	for i := 0; i < b.N; i++ {
		tree1.FromSymbolicExpressionSet(expressionSet1)
	}
}

func BenchmarkDualTree_FromSymbolicExpressionSet2_1000(b *testing.B) {
	expressionSet1 := GenerateRandomSymbolicExpressionSet(1000)
	tree1 := DualTree{}
	for i := 0; i < b.N; i++ {
		tree1.FromSymbolicExpressionSet2(expressionSet1)
	}
}

func BenchmarkDualTree_FromSymbolicExpressionSet_100000(b *testing.B) {
	expressionSet1 := GenerateRandomSymbolicExpressionSet(100000)
	tree100000 := DualTree{}
	for i := 0; i < b.N; i++ {
		tree100000.FromSymbolicExpressionSet(expressionSet1)
	}
}

func BenchmarkDualTree_FromSymbolicExpressionSet2_100000(b *testing.B) {
	expressionSet1 := GenerateRandomSymbolicExpressionSet(100000)
	tree1 := DualTree{}
	for i := 0; i < b.N; i++ {
		tree1.FromSymbolicExpressionSet2(expressionSet1)
	}
}



func BenchmarkDualTree_FromSymbolicExpressionSet_1000000(b *testing.B) {
	expressionSet1 := GenerateRandomSymbolicExpressionSet(1000000)
	tree1 := DualTree{}
	for i := 0; i < b.N; i++ {
		tree1.FromSymbolicExpressionSet(expressionSet1)
	}
}


func BenchmarkDualTree_FromSymbolicExpressionSet2_1000000(b *testing.B) {
	expressionSet1 := GenerateRandomSymbolicExpressionSet(1000000)
	tree1 := DualTree{}
	for i := 0; i < b.N; i++ {
		tree1.FromSymbolicExpressionSet2(expressionSet1)
	}
}

func BenchmarkGenerateRandomTree_1 (b *testing.B) {
	depth := 1
	terminals :=[]SymbolicExpression{X1, Const0, Const1, Const2, Const3, Const4, Const5, Const6,
		Const7, Const8, Const9}
	nonTerminals := []SymbolicExpression{Add, Mult, Sub}
	for i:=0; i < b.N; i++ {
		GenerateRandomTree(depth, terminals, nonTerminals)
	}
}
func BenchmarkGenerateRandomTree_10 (b *testing.B) {
	depth := 10
	terminals :=[]SymbolicExpression{X1, Const0, Const1, Const2, Const3, Const4, Const5, Const6,
		Const7, Const8, Const9}
	nonTerminals := []SymbolicExpression{Add, Mult, Sub}
	for i:=0; i < b.N; i++ {
		GenerateRandomTree(depth, terminals, nonTerminals)
	}
}
func BenchmarkGenerateRandomTree_20 (b *testing.B) {
	depth := 20
	terminals :=[]SymbolicExpression{X1, Const0, Const1, Const2, Const3, Const4, Const5, Const6,
		Const7, Const8, Const9}
	nonTerminals := []SymbolicExpression{Add, Mult, Sub}
	for i:=0; i < b.N; i++ {
		GenerateRandomTree(depth, terminals, nonTerminals)
	}
}