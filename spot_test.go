package gospot

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func testAny(q float64, level float64, maxExcess uint64, generator func(size uint64) []float64) (A, E, N int) {
	trainingSize := uint64(float64(maxExcess) / (1 - level))
	testSize := 10 * trainingSize
	data := generator(trainingSize)

	s, _ := NewSpot(q, false, true, level, maxExcess)
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

func BenchmarkSpot(b *testing.B) {
	s, err := NewSpot(1e-5, false, true, 0.99, 2000)
	if err != nil {
		panic(err)
	}
	data := gaussian(10000)
	s.Fit(data)

	A := 0
	E := 0
	N := 0

	b.ResetTimer()
	for _, x := range gaussian(uint64(b.N)) {
		switch s.Step(x) {
		case ANOMALY:
			A++
		case EXCESS:
			E++
		default:
			N++
		}
	}
}

func Example() {
	s, err := NewSpot(1e-5, false, true, 0.99, 2000)
	if err != nil {
		panic(err)
	}
	trainingSize := 10_000
	testingSize := 1_000_000
	data := make([]float64, trainingSize)
	for i := 0; i < trainingSize; i++ {
		data[i] = rand.NormFloat64()
	}
	s.Fit(data)

	A := 0
	E := 0
	N := 0

	for i := 0; i < testingSize; i++ {
		x := rand.NormFloat64()
		switch s.Step(x) {
		case ANOMALY:
			A++
		case EXCESS:
			E++
		default:
			N++
		}
	}

	fmt.Printf("ANOMALY:%d EXCESS:%d NORMAL:%d\n", A, E, N)
}
