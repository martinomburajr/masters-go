package utils

import "math/rand"

type NumberGenerator struct {
	Min int
	Max int
}

func (g *NumberGenerator) GenerateRandomInt() int {
	if g.Min == 0 && g.Max == 0 {
		return 0
	}
	if g.Min == 0 && g.Max != 0 {
		return rand.Intn(g.Max)
	}
	return rand.Intn(g.Max-g.Min) + g.Min
}

func (g *NumberGenerator) GenerateRandomFloat() float32 {
	if g.Max != 0 {
		return rand.Float32() * float32(g.Max)
	}
	return rand.Float32()
}
