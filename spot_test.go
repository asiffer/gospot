package gospot

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

// test to check if we can pipe options (more idiomatic go) but it creates other issue
// like parameter inconsitency. There is also the problem with the tail that must be allocated

func defaultSpot() *Spot {
	s, err := NewSpot(1e-5, false, true, 0.98, 1000)
	if err != nil {
		panic(err)
	}
	return s
}

func (s *Spot) WithQ(q float64) *Spot {
	if q >= (1.0-s.Level) || q <= 0.0 {
		panic("q must be in (0, 1-level)")
	}
	s.Q = q
	return s
}

func (s *Spot) LowerTail() *Spot {
	s.Low = true
	return s
}

func (s *Spot) UpperTail() *Spot {
	s.Low = false
	return s
}

func (s *Spot) WithLevel(level float64) *Spot {
	if level < 0.0 || level >= 1.0 {
		panic("level must be between 0 and 1")
	}
	if (1 - level) < s.Q {
		panic("1-level must be greater than q")
	}
	s.Level = level
	return s
}

func (s *Spot) WithMaxExcess(maxExcess uint64) *Spot {
	s.Reset() // we must reset everything in this case
	s.Tail = NewTail(maxExcess)
	return s
}

func testAny(s *Spot, generator func(size uint64) []float64) (A, E, N int) {
	trainingSize := uint64(float64(len(s.Tail.Peaks.Container.Data)) / (1 - s.Level))
	testSize := 10 * trainingSize
	data := generator(trainingSize)

	// s, _ := NewSpot(q, false, true, level, maxExcess)
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
	s := defaultSpot().WithQ(q).WithLevel(level).WithMaxExcess(maxExcess)

	A, E, N := testAny(s, uniform)
	r := float64(A) / float64(A+E+N)
	if math.Abs(r-q) > 2*q {
		t.Errorf("Anomaly ratio: %E (A:%d, E:%d, N:%d)", r, A, E, N)
	}
}

func TestGaussian(t *testing.T) {
	q := 5e-5
	level := 0.98
	maxExcess := uint64(2000)
	s := defaultSpot().WithQ(q).WithLevel(level).WithMaxExcess(maxExcess)

	A, E, N := testAny(s, gaussian)
	r := float64(A) / float64(A+E+N)
	if math.Abs(r-q) > 2*q {
		t.Errorf("Anomaly ratio: %E (A:%d, E:%d, N:%d)", r, A, E, N)
	}
}

func TestLowGaussian(t *testing.T) {
	q := 5e-5
	level := 0.98
	maxExcess := uint64(2000)
	s := defaultSpot().LowerTail().WithQ(q).WithLevel(level).WithMaxExcess(maxExcess)

	A, E, N := testAny(s, gaussian)
	r := float64(A) / float64(A+E+N)
	if math.Abs(r-q) > 2*q {
		t.Errorf("Anomaly ratio: %E (A:%d, E:%d, N:%d)", r, A, E, N)
	}
	fmt.Println(r, q)
}

func TestReset(t *testing.T) {
	q := 5e-5
	level := 0.98
	maxExcess := uint64(2000)

	s, err := NewSpot(q, false, true, level, maxExcess)
	if err != nil {
		t.Error(err)
	}
	s.Fit(gaussian(10 * maxExcess))
	if s.N == 0 {
		t.Errorf("N parameter not updated by fit")
	}
	if s.Nt == 0 {
		t.Errorf("no excess found")
	}
	if math.IsNaN(s.ExcessThreshold) {
		t.Errorf("excess threshold not computed")
	}
	if math.IsNaN(s.AnomalyThreshold) {
		t.Errorf("anomaly threshold not computed")
	}

	// reset the struct
	s.Reset()

	if s.N != 0 {
		t.Errorf("N attribute not reset")
	}
	if s.Nt != 0 {
		t.Errorf("Nt attribute not reset")
	}
	if !math.IsNaN(s.ExcessThreshold) {
		t.Errorf("excess threshold not reset")
	}
	if !math.IsNaN(s.AnomalyThreshold) {
		t.Errorf("anomaly threshold not reset")
	}
	if s.Tail.Peaks.Size() > 0 {
		t.Errorf("peaks container not reset")
	}
}

func TestBadNew(t *testing.T) {
	_, err := NewSpot(1e-5, false, true, -0.5, 1000)
	if err == nil {
		t.Errorf("must return an error because level < 0")
	}

	_, err = NewSpot(1e-5, false, true, 1., 1000)
	if err == nil {
		t.Errorf("must return an error because level >= 1.")
	}

	_, err = NewSpot(0.98, false, true, 0.99, 1000)
	if err == nil {
		t.Errorf("must return an error because q < level")
	}
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
