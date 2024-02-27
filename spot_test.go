package gospot

import (
	"fmt"
	"math"
	"testing"
)

func testAny(q float64, level float64, maxExcess uint64, generator func(size uint64) []float64) (A, E, N int) {
	trainingSize := uint64(float64(maxExcess) / (1 - level))
	testSize := 10 * trainingSize
	data := generator(trainingSize)

	s, _ := NewSpot(1e-5, false, true, level, maxExcess)
	s.Fit(data)

	A = 0
	E = 0
	N = 0

	for _, x := range generator(testSize) {
		switch s.Step(x) {
		case ANOMALY:
			A++
		case EXCESS:
			E++
		default:
			N++
		}
	}

	return
}

func TestUniform(t *testing.T) {
	q := 1e-5
	level := 0.98
	maxExcess := uint64(2000)

	A, E, N := testAny(q, level, maxExcess, uniform)
	r := float64(A) / float64(A+E+N)
	if math.Abs(r-q) > 2*q {
		t.Errorf("Anomaly ratio: %E (A:%d, E:%d, N:%d)", r, A, E, N)
	}
	fmt.Println(r, q)
}

func TestGaussian(t *testing.T) {
	q := 5e-5
	level := 0.98
	maxExcess := uint64(2000)

	A, E, N := testAny(q, level, maxExcess, gaussian)
	r := float64(A) / float64(A+E+N)
	if math.Abs(r-q) > 2*q {
		t.Errorf("Anomaly ratio: %E (A:%d, E:%d, N:%d)", r, A, E, N)
	}
	fmt.Println(r, q)
}
