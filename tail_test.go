package gospot

import (
	"fmt"
	"math"
	"sort"
	"testing"
)

func testFit(data []float64, size uint64, check func(float64) bool) bool {
	tail := NewTail(size)

	initial := sort.Float64Slice(data)
	initial.Sort()
	index := uint64(len(initial)) - size
	u := initial[index:]

	for _, ui := range u {
		tail.Push(ui - initial[index-1])
	}

	tail.Fit()
	return check(tail.Gamma)
}

func TestFitUniform(t *testing.T) {
	var size uint64 = 1000
	check := func(g float64) bool { return (g < -0.15) }

	N := 50
	s := 0
	for i := 0; i < N; i++ {
		data := uniformX(100*size, 15.)
		if success := testFit(data, size, check); success {
			s++
		}
	}
	result := float64(s) / float64(N)
	if result < 0.50 {
		t.Errorf("Success rate: %f%%", 100*result)
	} else {
		t.Logf("Success rate: %f%%", 100*result)
	}
	fmt.Println(result)
}

func TestFitGaussian(t *testing.T) {
	var size uint64 = 1000
	check := func(g float64) bool { return (math.Abs(g) < 0.15) }

	N := 50
	s := 0
	for i := 0; i < N; i++ {
		data := gaussian(100 * size)
		// name := fmt.Sprintf("FitUniform %d/%d", i, N)
		if success := testFit(data, size, check); success {
			s++
		}
	}
	result := float64(s) / float64(N)
	if result < 0.90 {
		t.Errorf("Success rate: %f%%", 100*result)
	} else {
		t.Logf("Success rate: %f%%", 100*result)
	}
}

func TestFitLogGaussian(t *testing.T) {
	var size uint64 = 1000
	check := func(g float64) bool { return (g > 0.15) }

	N := 50
	s := 0
	for i := 0; i < N; i++ {
		data := logGaussian(100 * size)
		// name := fmt.Sprintf("FitUniform %d/%d", i, N)
		if success := testFit(data, size, check); success {
			s++
		}
	}
	result := float64(s) / float64(N)
	if result < 0.95 {
		t.Errorf("Success rate: %f%%", 100*result)
	} else {
		t.Logf("Success rate: %f%%", 100*result)
	}
}

func newTail(data []float64, size uint64) (*Tail, float64) {
	tail := NewTail(size)

	initial := sort.Float64Slice(data)
	initial.Sort()
	index := uint64(len(initial)) - size
	u := initial[index:]

	th := initial[index-1]
	for _, ui := range u {
		tail.Push(ui - th)
	}

	tail.Fit()
	return tail, th
}

func testProbabilityQuantile(data []float64, size uint64, pz float64, z float64, ztol float64) bool {
	tail, th := newTail(data, size)
	s := float64(size) / float64(len(data))

	d := z - th
	p := tail.Probability(s, d)
	q := th + tail.Quantile(s, 1-pz)

	// fmt.Println(p, 1-pz, math.Abs(p-(1-pz)))
	// fmt.Println(q, z, math.Abs(q-z)/z)
	return math.Abs(p-(1-pz)) <= 0.02 && math.Abs(q-z)/z <= ztol
}

func TestProbabilityQuantileUniform(t *testing.T) {
	N := 50
	var size uint64 = 1000
	s := 0

	for i := 0; i < N; i++ {
		data := uniform(100 * size)
		if success := testProbabilityQuantile(data, size, 0.999, 0.999, 0.02); success {
			s++
		}
	}
	result := float64(s) / float64(N)
	if result < 0.95 {
		t.Errorf("Success rate: %f%%", 100*result)
	} else {
		t.Logf("Success rate: %f%%", 100*result)
	}
}

func TestProbabilityQuantileGaussian(t *testing.T) {
	N := 50
	var size uint64 = 1000
	s := 0

	for i := 0; i < N; i++ {
		data := gaussian(100 * size)
		if success := testProbabilityQuantile(data, size, 0.999, 3.090232306167813, 0.02); success {
			s++
		}
	}
	result := float64(s) / float64(N)
	if result < 0.90 {
		t.Errorf("Success rate: %f%%", 100*result)
	} else {
		t.Logf("Success rate: %f%%", 100*result)
	}
}

func TestProbabilityQuantileLogGaussian(t *testing.T) {
	N := 50
	var size uint64 = 1000
	s := 0

	for i := 0; i < N; i++ {
		data := logGaussian(100 * size)
		if success := testProbabilityQuantile(data, size, 0.999, 21.982183979582828, 0.05); success {
			s++
		}
	}
	result := float64(s) / float64(N)
	if result < 0.90 {
		t.Errorf("Success rate: %f%%", 100*result)
	} else {
		t.Logf("Success rate: %f%%", 100*result)
	}
}
