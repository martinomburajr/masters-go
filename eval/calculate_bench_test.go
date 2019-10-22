package eval

import (
	"github.com/Knetic/govaluate"
	"github.com/PaesslerAG/gval"
	"github.com/martinomburajr/masters-go/utils"
	"github.com/soniah/evaler"
	"testing"
)

var ans1 float64
var expression = "10.9*9.8*8.7*7.6*6.5*5.4*4.3*3.2*2.1*1"
var expressionVar = "x*10.9*9.8*8.7*7.6*6.5*5.4*4.3*3.2*2.1*1"
var expressionManyVar = "x*x*x*x*x*x*x*x*x*x*x*x*x"
var expressionManyVarXY = "x*y*x*y*x*y*x*y*x*y"

func BenchmarkCalculate(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ans, err := Calculate(expression)
		if err != nil {
			b.Error(err)
		}
		ans1 = ans
	}
}

var ans = 0.0
func BenchmarkCalculateV2(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N ; i++  {
		ans, _ =CalculateV2(expression)
	}
}

func BenchmarkCalcEvaler(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = evaler.Eval(expression)
	}
}

func BenchmarkGoValuate(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ex, _ := govaluate.NewEvaluableExpression(expression)
		evaluate, _ := ex.Evaluate(nil)
		ans, _ = utils.ConvertToFloat64(evaluate)
	}
}

func BenchmarkGVal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ans, err := gval.Evaluate(expression, nil)
		if err != nil {
			b.Error(err)
		}
		ans1, err = utils.ConvertToFloat64(ans)
		if err != nil {
			b.Error(err)
		}
	}
}
//
//func BenchmarkCalculateWithVarX(b *testing.B) {
//	b.ReportAllocs()
//	for i := 0; i < b.N; i++ {
//		ans, err := CalculateWithVar(expressionVar, map[string]float64{"x": 10})
//		if err != nil {
//			b.Error(err)
//		}
//		ans1 = ans
//	}
//}
//
//func BenchmarkGValWithVarX(b *testing.B) {
//	b.ReportAllocs()
//	for i := 0; i < b.N; i++ {
//		ans, err := gval.Evaluate(expressionVar, map[string]float64{"x": 10})
//		if err != nil {
//			b.Error(err)
//		}
//		ans1, err = utils.ConvertToFloat64(ans)
//		if err != nil {
//			b.Error(err)
//		}
//	}
//}
//
//func BenchmarkCalculateWithManyVarX(b *testing.B) {
//	b.ReportAllocs()
//	for i := 0; i < b.N; i++ {
//		ans, err := CalculateWithVar(expressionManyVar, map[string]float64{"x": 10})
//		if err != nil {
//			b.Error(err)
//		}
//		ans1 = ans
//	}
//}
//
//func BenchmarkGValWithManyVarX(b *testing.B) {
//	b.ReportAllocs()
//	for i := 0; i < b.N; i++ {
//		ans, err := gval.Evaluate(expressionManyVar, map[string]float64{"x": 10})
//		if err != nil {
//			b.Error(err)
//		}
//		ans1, err = utils.ConvertToFloat64(ans)
//		if err != nil {
//			b.Error(err)
//		}
//	}
//}
//
//func BenchmarkCalculateWithManyVarXY(b *testing.B) {
//	b.ReportAllocs()
//	for i := 0; i < b.N; i++ {
//		ans, err := CalculateWithVar(expressionManyVar, map[string]float64{"x": 10, "y": 5})
//		if err != nil {
//			b.Error(err)
//		}
//		ans1 = ans
//	}
//}
//
//func BenchmarkGValWithManyVarXY(b *testing.B) {
//	b.ReportAllocs()
//	for i := 0; i < b.N; i++ {
//		ans, err := gval.Evaluate(expressionManyVar, map[string]float64{"x": 10, "y": 5})
//		if err != nil {
//			b.Error(err)
//		}
//		ans1, err = utils.ConvertToFloat64(ans)
//		if err != nil {
//			b.Error(err)
//		}
//	}
//}
//
//var negativeNumber = "0--1+2--2--3"
//var negativeNumberAns = ""
//
//func BenchmarkNegativeNumberParser(b *testing.B) {
//	b.ReportAllocs()
//	for i := 0; i < b.N; i++ {
//		g := NegativeNumberParser(negativeNumber)
//		negativeNumberAns = g
//	}
//}
//
//var replacerAns = ""
//
//func BenchmarkMartinsReplace(b *testing.B) {
//	b.ReportAllocs()
//	for i := 0; i < b.N; i++ {
//		replacerAns = MartinsReplace("x * x * x", " ", "")
//	}
//}
//
//func BenchmarkReplace(b *testing.B) {
//	b.ReportAllocs()
//	for i := 0; i < b.N; i++ {
//		replacerAns = strings.ReplaceAll("x * x * x", " ", "")
//	}
//}
//
//var replacerLong = "x                       *                                  x                              *                                            x"
//
//func BenchmarkMartinsReplaceLong(b *testing.B) {
//	b.ReportAllocs()
//	for i := 0; i < b.N; i++ {
//		replacerAns = MartinsReplace(replacerLong, " ", "")
//	}
//}
//
//func BenchmarkReplaceLong(b *testing.B) {
//	b.ReportAllocs()
//	for i := 0; i < b.N; i++ {
//		replacerAns = strings.ReplaceAll(replacerLong, " ", "")
//	}
//}


