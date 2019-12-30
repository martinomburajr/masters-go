package main

import (
	"testing"
)

func TestEvolution1(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func BenchmarkEvolution1(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		paramsFolder := "_params"
		parallelism := true
		logging := false

		Scheduler(paramsFolder, parallelism, logging, false)
	}
}
